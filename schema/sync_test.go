/*
 *     Copyright 2025 The CNCF ModelPack Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package schema_test

import (
	"bufio"
	"encoding/json"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	v1 "github.com/modelpack/model-spec/specs-go/v1"
	digest "github.com/opencontainers/go-digest"
)

const schemaTypeArray = "array"

// schemaProperty represents a single property in a JSON schema definition.
type schemaProperty struct {
	Type   string                    `json:"type"`
	Ref    string                    `json:"$ref"`
	Format string                    `json:"format"`
	Items  *schemaProperty           `json:"items"`
	Props  map[string]schemaProperty `json:"properties"`
}

// schemaDef represents a JSON schema definition block.
type schemaDef struct {
	Type  string                    `json:"type"`
	Props map[string]schemaProperty `json:"properties"`
}

// schemaRoot represents the top-level JSON schema.
type schemaRoot struct {
	Props map[string]schemaProperty `json:"properties"`
	Defs  map[string]schemaDef      `json:"$defs"`
}

// expectedType describes what a Go type should look like in JSON schema.
type expectedType struct {
	typ    string        // "string", "boolean", schemaTypeArray, "object"
	format string        // "date-time" or ""
	ref    string        // "$ref" value or ""
	items  *expectedType // for array element type
}

// mapGoType converts a Go reflect.Type to the expected JSON schema type.
func mapGoType(t reflect.Type) expectedType {
	// unwrap pointer
	if t.Kind() == reflect.Ptr {
		inner := t.Elem()
		if inner == reflect.TypeOf(time.Time{}) {
			return expectedType{typ: "string", format: "date-time"}
		}
		if inner.Kind() == reflect.Bool {
			return expectedType{typ: "boolean"}
		}
		return expectedType{ref: "#/$defs/" + inner.Name()}
	}

	// slice types
	if t.Kind() == reflect.Slice {
		elem := mapGoType(t.Elem())
		return expectedType{typ: schemaTypeArray, items: &elem}
	}

	// named string types (Modality, digest.Digest)
	if t.Kind() == reflect.String && t.Name() != "string" {
		if t == reflect.TypeOf(digest.Digest("")) {
			return expectedType{typ: "string"}
		}
		return expectedType{ref: "#/$defs/" + t.Name()}
	}

	// non-pointer struct
	if t.Kind() == reflect.Struct {
		return expectedType{ref: "#/$defs/" + t.Name()}
	}

	switch t.Kind() {
	case reflect.String:
		return expectedType{typ: "string"}
	case reflect.Bool:
		return expectedType{typ: "boolean"}
	default:
		return expectedType{typ: "unknown"}
	}
}

// schemaTypeString returns a human-readable description of a schema property.
func schemaTypeString(p schemaProperty) string {
	if p.Ref != "" {
		return p.Ref
	}
	s := p.Type
	if p.Format != "" {
		s += " (format: " + p.Format + ")"
	}
	if p.Type == schemaTypeArray && p.Items != nil {
		s += " of " + schemaTypeString(*p.Items)
	}
	return s
}

// expectedTypeString returns a human-readable description of an expected type.
func expectedTypeString(e expectedType) string {
	if e.ref != "" {
		return e.ref
	}
	s := e.typ
	if e.format != "" {
		s += " (format: " + e.format + ")"
	}
	if e.typ == schemaTypeArray && e.items != nil {
		s += " of " + expectedTypeString(*e.items)
	}
	return s
}

// matchesExpected checks if a JSON schema property matches the expected type from Go reflection.
func matchesExpected(p schemaProperty, e expectedType) bool {
	if e.ref != "" || p.Ref != "" {
		return p.Ref == e.ref
	}
	if p.Type != e.typ || p.Format != e.format {
		return false
	}
	if e.typ == schemaTypeArray {
		if (p.Items == nil) != (e.items == nil) {
			return false
		}
		if p.Items != nil {
			return matchesExpected(*p.Items, *e.items)
		}
	}
	return true
}

// structFields extracts JSON-tagged fields from a Go struct type.
func structFields(t reflect.Type) map[string]reflect.StructField {
	fields := make(map[string]reflect.StructField)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		jsonName := strings.Split(tag, ",")[0]
		fields[jsonName] = f
	}
	return fields
}

func loadSchema(t *testing.T) schemaRoot {
	t.Helper()
	data, err := os.ReadFile("config-schema.json")
	if err != nil {
		t.Fatalf("failed to read config-schema.json: %v", err)
	}
	var root schemaRoot
	if err := json.Unmarshal(data, &root); err != nil {
		t.Fatalf("failed to parse config-schema.json: %v", err)
	}
	return root
}

// =============================================================================
// Layer 1: Go structs <-> JSON Schema
// =============================================================================

// TestStructSchemaSync verifies that Go struct fields and JSON schema properties
// are in sync. This prevents drift where fields are added to Go structs but
// not to the schema, or vice versa.
func TestStructSchemaSync(t *testing.T) {
	root := loadSchema(t)

	checks := []struct {
		goType     reflect.Type
		schemaName string // empty = root-level properties
	}{
		{reflect.TypeOf(v1.Model{}), ""},
		{reflect.TypeOf(v1.ModelDescriptor{}), "ModelDescriptor"},
		{reflect.TypeOf(v1.ModelConfig{}), "ModelConfig"},
		{reflect.TypeOf(v1.ModelCapabilities{}), "ModelCapabilities"},
		{reflect.TypeOf(v1.ModelFS{}), "ModelFS"},
	}

	for _, check := range checks {
		name := check.goType.Name()
		t.Run(name, func(t *testing.T) {
			var schemaProps map[string]schemaProperty
			if check.schemaName == "" {
				schemaProps = root.Props
			} else {
				def, ok := root.Defs[check.schemaName]
				if !ok {
					t.Fatalf("schema definition %q not found in $defs", check.schemaName)
				}
				schemaProps = def.Props
			}

			goFields := structFields(check.goType)

			// Go fields missing from schema
			for jsonName, field := range goFields {
				prop, exists := schemaProps[jsonName]
				if !exists {
					expected := mapGoType(field.Type)
					t.Errorf("Go struct field %q (type: %s, expected schema: %s) missing from JSON schema",
						jsonName, field.Type, expectedTypeString(expected))
					continue
				}
				expected := mapGoType(field.Type)
				if !matchesExpected(prop, expected) {
					t.Errorf("type mismatch for field %q: Go expects %s, schema has %s",
						jsonName, expectedTypeString(expected), schemaTypeString(prop))
				}
			}

			// Schema properties missing from Go
			for propName := range schemaProps {
				if _, exists := goFields[propName]; !exists {
					t.Errorf("schema property %q has no corresponding Go struct field in %s",
						propName, name)
				}
			}
		})
	}
}

// =============================================================================
// Layer 2: JSON Schema <-> Doc example (docs/config.md)
// =============================================================================

// TestDocExampleCoversSchemaFields verifies that the JSON example in
// docs/config.md includes every field defined in the JSON schema.
// This catches drift where a field is added to the schema but never
// shown in documentation, or vice versa.
func TestDocExampleCoversSchemaFields(t *testing.T) {
	root := loadSchema(t)

	// Extract the JSON example from docs/config.md
	example := extractJSONExample(t, "../docs/config.md")

	// Parse the example JSON
	var exampleDoc map[string]interface{}
	if err := json.Unmarshal([]byte(example), &exampleDoc); err != nil {
		t.Fatalf("failed to parse JSON example from docs/config.md: %v", err)
	}

	// Check root-level properties
	t.Run("root", func(t *testing.T) {
		checkExampleCoversSchema(t, root.Props, exampleDoc, "root")
	})

	// Check each $defs section against the corresponding nested object in the example
	defToPath := map[string]string{
		"ModelDescriptor":   "descriptor",
		"ModelConfig":       "config",
		"ModelFS":           "modelfs",
		"ModelCapabilities": "config.capabilities",
	}

	for defName, jsonPath := range defToPath {
		def, ok := root.Defs[defName]
		if !ok {
			continue
		}
		t.Run(defName, func(t *testing.T) {
			nested := navigateToPath(exampleDoc, jsonPath)
			if nested == nil {
				t.Errorf("doc example missing top-level object %q (for schema def %s)", jsonPath, defName)
				return
			}
			nestedMap, ok := nested.(map[string]interface{})
			if !ok {
				t.Errorf("doc example %q is not an object", jsonPath)
				return
			}
			checkExampleCoversSchema(t, def.Props, nestedMap, defName)
		})
	}
}

// checkExampleCoversSchema verifies bidirectional coverage between schema and example.
func checkExampleCoversSchema(t *testing.T, schemaProps map[string]schemaProperty, example map[string]interface{}, context string) {
	t.Helper()

	// Schema fields missing from example
	for propName := range schemaProps {
		if _, exists := example[propName]; !exists {
			t.Errorf("schema field %q (%s) not demonstrated in docs/config.md example", propName, context)
		}
	}

	// Example fields not in schema
	for exKey := range example {
		if _, exists := schemaProps[exKey]; !exists {
			// Check if it's a $ref to a known def (e.g., "capabilities" is under ModelConfig)
			if !isRefField(schemaProps, exKey) {
				t.Errorf("docs/config.md example has field %q (%s) not defined in JSON schema", exKey, context)
			}
		}
	}
}

// isRefField checks if a field name exists as a $ref target in the schema properties.
func isRefField(props map[string]schemaProperty, name string) bool {
	for _, p := range props {
		if p.Ref != "" && strings.HasSuffix(p.Ref, "/"+strings.ToUpper(name[:1])+name[1:]) {
			return true
		}
	}
	return false
}

// navigateToPath navigates a nested map by dot-separated or slash-separated path.
func navigateToPath(obj map[string]interface{}, path string) interface{} {
	parts := strings.Split(path, ".")
	var current interface{} = obj
	for _, part := range parts {
		m, ok := current.(map[string]interface{})
		if !ok {
			return nil
		}
		current, ok = m[part]
		if !ok {
			return nil
		}
	}
	return current
}

// extractJSONExample extracts the first JSON code block tagged with a mediatype
// from a markdown file.
func extractJSONExample(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read %s: %v", path, err)
	}
	content := string(data)

	// Find JSON code blocks - look for ```json blocks
	inBlock := false
	var blockLines []string
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "```json") && !inBlock {
			inBlock = true
			blockLines = nil
			continue
		}
		if strings.HasPrefix(line, "```") && inBlock {
			inBlock = false
			// Return the first complete JSON block
			block := strings.Join(blockLines, "\n")
			if strings.Contains(block, "\"descriptor\"") {
				return block
			}
		}
		if inBlock {
			blockLines = append(blockLines, line)
		}
	}
	t.Fatal("no JSON example with 'descriptor' found in docs/config.md")
	return ""
}

// =============================================================================
// Layer 3: Doc field descriptions <-> JSON Schema
// =============================================================================

// TestDocDescriptionsMatchSchema parses the property descriptions in
// docs/config.md (lines like "- **fieldName** _type_, REQUIRED/OPTIONAL")
// and verifies every documented field exists in the schema and vice versa.
func TestDocDescriptionsMatchSchema(t *testing.T) {
	root := loadSchema(t)
	docFields := parseDocFields(t, "../docs/config.md")

	// Map doc sections to schema definitions
	sectionToDef := map[string]string{
		"descriptor":   "ModelDescriptor",
		"config":       "ModelConfig",
		"modelfs":      "ModelFS",
		"capabilities": "ModelCapabilities",
	}

	for section, defName := range sectionToDef {
		t.Run(defName, func(t *testing.T) {
			def, ok := root.Defs[defName]
			if !ok {
				t.Fatalf("schema definition %q not found", defName)
			}

			fields, ok := docFields[section]
			if !ok {
				t.Fatalf("docs/config.md has no section for %q", section)
			}

			// Doc fields missing from schema
			for _, fieldName := range fields {
				if _, exists := def.Props[fieldName]; !exists {
					t.Errorf("docs/config.md documents field %q in %s, but it is missing from JSON schema",
						fieldName, defName)
				}
			}

			// Schema fields missing from docs
			docSet := make(map[string]bool)
			for _, f := range fields {
				docSet[f] = true
			}
			for propName := range def.Props {
				if !docSet[propName] {
					t.Errorf("JSON schema has field %q in %s, but docs/config.md does not document it",
						propName, defName)
				}
			}
		})
	}
}

// parseDocFields extracts documented field names from docs/config.md.
// It parses lines like: "  - **fieldName** _type_, REQUIRED/OPTIONAL"
// and groups them by their parent section (descriptor, config, modelfs, capabilities).
func parseDocFields(t *testing.T, path string) map[string][]string {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("failed to open %s: %v", path, err)
	}
	defer file.Close()

	// Matches lines like: "- **fieldName** _type_, OPTIONAL" or "  - **fieldName** _type_, REQUIRED"
	fieldPattern := regexp.MustCompile(`^\s*-\s+\*\*(\w+)\*\*\s+_[^_]+_`)

	result := make(map[string][]string)
	currentSection := ""

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Detect top-level sections: "- **descriptor** _object_, REQUIRED"
		if matches := fieldPattern.FindStringSubmatch(line); matches != nil {
			fieldName := matches[1]
			indent := len(line) - len(strings.TrimLeft(line, " "))

			if indent <= 1 {
				// Top-level field — this is a section
				switch fieldName {
				case "descriptor", "config", "modelfs":
					currentSection = fieldName
				default:
					currentSection = ""
				}
			} else if indent >= 2 && currentSection != "" {
				// Nested field under a section
				if fieldName == "capabilities" {
					// capabilities is both a field of config AND a sub-section
					result[currentSection] = append(result[currentSection], fieldName)
					currentSection = "capabilities"
				} else {
					result[currentSection] = append(result[currentSection], fieldName)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("failed to scan %s: %v", path, err)
	}

	return result
}
