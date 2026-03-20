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
	"encoding/json"
	"os"
	"reflect"
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
		// pointer to named struct
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

	// non-pointer struct (e.g. ModelDescriptor in Model)
	if t.Kind() == reflect.Struct {
		return expectedType{ref: "#/$defs/" + t.Name()}
	}

	// basic types
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
	// If either side uses a $ref, both must agree on the ref value.
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

// TestStructSchemaSync verifies that Go struct fields and JSON schema properties
// are in sync. This prevents drift where fields are added to Go structs but
// not to the schema, or vice versa.
func TestStructSchemaSync(t *testing.T) {
	data, err := os.ReadFile("config-schema.json")
	if err != nil {
		t.Fatalf("failed to read config-schema.json: %v", err)
	}

	var root schemaRoot
	if err := json.Unmarshal(data, &root); err != nil {
		t.Fatalf("failed to parse config-schema.json: %v", err)
	}

	// checks maps a Go struct to its corresponding JSON schema definition.
	// Empty schemaName means use the root-level properties.
	checks := []struct {
		goType     reflect.Type
		schemaName string
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
			// get schema properties
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

			// check: struct fields missing from schema
			for jsonName, field := range goFields {
				prop, exists := schemaProps[jsonName]
				if !exists {
					expected := mapGoType(field.Type)
					t.Errorf("struct field %q (type: %s, expected schema: %s) missing from JSON schema %q",
						jsonName, field.Type, expectedTypeString(expected), name)
					continue
				}

				// check type match
				expected := mapGoType(field.Type)
				if !matchesExpected(prop, expected) {
					t.Errorf("type mismatch for field %q in %q: Go expects %s, schema has %s",
						jsonName, name, expectedTypeString(expected), schemaTypeString(prop))
				}
			}

			// check: schema properties missing from struct
			for propName := range schemaProps {
				if _, exists := goFields[propName]; !exists {
					t.Errorf("schema property %q in %q has no corresponding Go struct field",
						propName, name)
				}
			}
		})
	}
}
