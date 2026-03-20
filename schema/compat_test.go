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
	"strings"
	"testing"
	"time"

	digest "github.com/opencontainers/go-digest"

	v1 "github.com/modelpack/model-spec/specs-go/v1"
)

// TestDownstreamDetectionFields verifies that a fully populated ModelPack config
// serializes to JSON containing the fields that downstream consumers (e.g.,
// Docker Model Runner) use to detect ModelPack artifacts:
//   - "config" containing "paramSize"
//   - "descriptor" containing "createdAt"
//   - "modelfs" as a top-level key
//
// If any of these field names change, downstream detection will break.
func TestDownstreamDetectionFields(t *testing.T) {
	now := time.Now().UTC()
	boolTrue := true

	model := v1.Model{
		Descriptor: v1.ModelDescriptor{
			CreatedAt: &now,
			Name:      "test-model",
			Family:    "llama3",
			Version:   "1.0",
		},
		ModelFS: v1.ModelFS{
			Type:    "layers",
			DiffIDs: []digest.Digest{"sha256:abc123"},
		},
		Config: v1.ModelConfig{
			Architecture: "transformer",
			Format:       "gguf",
			ParamSize:    "8b",
			Precision:    "fp16",
			Quantization: "q4_0",
			Capabilities: &v1.ModelCapabilities{
				InputTypes:  []v1.Modality{v1.TextModality},
				OutputTypes: []v1.Modality{v1.TextModality},
				Reasoning:   &boolTrue,
			},
		},
	}

	data, err := json.Marshal(model)
	if err != nil {
		t.Fatalf("failed to marshal Model: %v", err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("failed to unmarshal to map: %v", err)
	}

	// Top-level keys that downstream consumers depend on.
	for _, key := range []string{"descriptor", "modelfs", "config"} {
		if _, ok := raw[key]; !ok {
			t.Errorf("top-level key %q missing from serialized Model JSON — downstream detection will break", key)
		}
	}

	// Verify descriptor contains "createdAt".
	var desc map[string]json.RawMessage
	if err := json.Unmarshal(raw["descriptor"], &desc); err != nil {
		t.Fatalf("failed to unmarshal descriptor: %v", err)
	}
	if _, ok := desc["createdAt"]; !ok {
		t.Error("descriptor missing \"createdAt\" field — downstream detection will break")
	}

	// Verify config contains "paramSize".
	var cfg map[string]json.RawMessage
	if err := json.Unmarshal(raw["config"], &cfg); err != nil {
		t.Fatalf("failed to unmarshal config: %v", err)
	}
	if _, ok := cfg["paramSize"]; !ok {
		t.Error("config missing \"paramSize\" field — downstream detection will break")
	}
}

// TestDownstreamFieldMapping verifies that every config field used by
// downstream consumers maps to the expected JSON key name. These mappings
// MUST NOT change without coordinating with downstream projects.
func TestDownstreamFieldMapping(t *testing.T) {
	now := time.Now().UTC()
	boolTrue := true

	model := v1.Model{
		Descriptor: v1.ModelDescriptor{
			CreatedAt:   &now,
			Name:        "test-model",
			Family:      "llama3",
			Description: "A test model",
			Licenses:    []string{"Apache-2.0"},
		},
		ModelFS: v1.ModelFS{
			Type:    "layers",
			DiffIDs: []digest.Digest{"sha256:abc123"},
		},
		Config: v1.ModelConfig{
			ParamSize:    "8b",
			Format:       "gguf",
			Quantization: "q4_0",
			Architecture: "transformer",
			Capabilities: &v1.ModelCapabilities{
				InputTypes:  []v1.Modality{v1.TextModality},
				OutputTypes: []v1.Modality{v1.TextModality},
				Reasoning:   &boolTrue,
			},
		},
	}

	data, err := json.Marshal(model)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Parse into nested maps to check field names at each level.
	var full map[string]json.RawMessage
	if err := json.Unmarshal(data, &full); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	// Descriptor field mappings used by downstream consumers.
	var desc map[string]json.RawMessage
	if err := json.Unmarshal(full["descriptor"], &desc); err != nil {
		t.Fatalf("failed to unmarshal descriptor: %v", err)
	}
	descFields := []string{"createdAt", "name", "family", "description", "licenses"}
	for _, f := range descFields {
		if _, ok := desc[f]; !ok {
			t.Errorf("descriptor missing expected field %q", f)
		}
	}

	// Config field mappings used by downstream consumers.
	var cfg map[string]json.RawMessage
	if err := json.Unmarshal(full["config"], &cfg); err != nil {
		t.Fatalf("failed to unmarshal config: %v", err)
	}
	cfgFields := []string{"paramSize", "format", "quantization", "architecture", "capabilities"}
	for _, f := range cfgFields {
		if _, ok := cfg[f]; !ok {
			t.Errorf("config missing expected field %q", f)
		}
	}

	// ModelFS field mappings used by downstream consumers.
	var mfs map[string]json.RawMessage
	if err := json.Unmarshal(full["modelfs"], &mfs); err != nil {
		t.Fatalf("failed to unmarshal modelfs: %v", err)
	}
	mfsFields := []string{"type", "diffIds"}
	for _, f := range mfsFields {
		if _, ok := mfs[f]; !ok {
			t.Errorf("modelfs missing expected field %q", f)
		}
	}
}

// TestDownstreamMediaTypePrefixes verifies that the ModelPack media type
// constants use the expected prefix that downstream consumers rely on for
// layer type detection.
func TestDownstreamMediaTypePrefixes(t *testing.T) {
	prefix := "application/vnd.cncf.model."

	mediaTypes := []struct {
		name  string
		value string
	}{
		{"MediaTypeModelConfig", v1.MediaTypeModelConfig},
		{"MediaTypeModelWeightRaw", v1.MediaTypeModelWeightRaw},
		{"MediaTypeModelWeight", v1.MediaTypeModelWeight},
		{"MediaTypeModelWeightGzip", v1.MediaTypeModelWeightGzip},
		{"MediaTypeModelWeightZstd", v1.MediaTypeModelWeightZstd},
		{"MediaTypeModelWeightConfigRaw", v1.MediaTypeModelWeightConfigRaw},
		{"MediaTypeModelWeightConfig", v1.MediaTypeModelWeightConfig},
		{"MediaTypeModelWeightConfigGzip", v1.MediaTypeModelWeightConfigGzip},
		{"MediaTypeModelWeightConfigZstd", v1.MediaTypeModelWeightConfigZstd},
		{"MediaTypeModelDocRaw", v1.MediaTypeModelDocRaw},
		{"MediaTypeModelDoc", v1.MediaTypeModelDoc},
		{"MediaTypeModelDocGzip", v1.MediaTypeModelDocGzip},
		{"MediaTypeModelDocZstd", v1.MediaTypeModelDocZstd},
		{"MediaTypeModelCodeRaw", v1.MediaTypeModelCodeRaw},
		{"MediaTypeModelCode", v1.MediaTypeModelCode},
		{"MediaTypeModelCodeGzip", v1.MediaTypeModelCodeGzip},
		{"MediaTypeModelCodeZstd", v1.MediaTypeModelCodeZstd},
		{"MediaTypeModelDatasetRaw", v1.MediaTypeModelDatasetRaw},
		{"MediaTypeModelDataset", v1.MediaTypeModelDataset},
		{"MediaTypeModelDatasetGzip", v1.MediaTypeModelDatasetGzip},
		{"MediaTypeModelDatasetZstd", v1.MediaTypeModelDatasetZstd},
	}

	for _, mt := range mediaTypes {
		t.Run(mt.name, func(t *testing.T) {
			if !strings.HasPrefix(mt.value, prefix) {
				t.Errorf("%s = %q, does not have expected prefix %q", mt.name, mt.value, prefix)
			}
		})
	}
}

// TestDownstreamWeightMediaTypes verifies the exact media type strings for
// model weight layers that downstream consumers use for format detection
// and conversion.
func TestDownstreamWeightMediaTypes(t *testing.T) {
	tests := []struct {
		name     string
		got      string
		expected string
	}{
		{
			name:     "raw weight",
			got:      v1.MediaTypeModelWeightRaw,
			expected: "application/vnd.cncf.model.weight.v1.raw",
		},
		{
			name:     "tar weight",
			got:      v1.MediaTypeModelWeight,
			expected: "application/vnd.cncf.model.weight.v1.tar",
		},
		{
			name:     "gzip weight",
			got:      v1.MediaTypeModelWeightGzip,
			expected: "application/vnd.cncf.model.weight.v1.tar+gzip",
		},
		{
			name:     "zstd weight",
			got:      v1.MediaTypeModelWeightZstd,
			expected: "application/vnd.cncf.model.weight.v1.tar+zstd",
		},
		{
			name:     "raw weight config",
			got:      v1.MediaTypeModelWeightConfigRaw,
			expected: "application/vnd.cncf.model.weight.config.v1.raw",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("got %q, want %q — changing this value will break downstream media type mapping", tt.got, tt.expected)
			}
		})
	}
}

// TestDownstreamArtifactType verifies the artifact type constant used in OCI
// manifests. Downstream consumers match on this to identify ModelPack manifests.
func TestDownstreamArtifactType(t *testing.T) {
	expected := "application/vnd.cncf.model.manifest.v1+json"
	if v1.ArtifactTypeModelManifest != expected {
		t.Errorf("ArtifactTypeModelManifest = %q, want %q — changing this will break downstream manifest detection", v1.ArtifactTypeModelManifest, expected)
	}
}

// TestDownstreamRoundTrip verifies that a ModelPack config can be marshalled
// and unmarshalled without losing any fields. This ensures that downstream
// consumers can reliably parse configs produced by ModelPack tooling.
func TestDownstreamRoundTrip(t *testing.T) {
	now := time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC)
	boolTrue := true
	boolFalse := false

	original := v1.Model{
		Descriptor: v1.ModelDescriptor{
			CreatedAt:   &now,
			Authors:     []string{"CNCF ModelPack Authors"},
			Family:      "llama3",
			Name:        "llama3-8b-instruct",
			DocURL:      "https://example.com/docs",
			SourceURL:   "https://example.com/source",
			DatasetsURL: []string{"https://example.com/dataset1"},
			Version:     "3.1",
			Revision:    "abc123",
			Vendor:      "Example Corp",
			Licenses:    []string{"Apache-2.0"},
			Title:       "Llama 3 8B Instruct",
			Description: "An instruction-tuned language model",
		},
		ModelFS: v1.ModelFS{
			Type: "layers",
			DiffIDs: []digest.Digest{
				"sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			},
		},
		Config: v1.ModelConfig{
			Architecture: "transformer",
			Format:       "gguf",
			ParamSize:    "8b",
			Precision:    "fp16",
			Quantization: "q4_0",
			Capabilities: &v1.ModelCapabilities{
				InputTypes:  []v1.Modality{v1.TextModality, v1.ImageModality},
				OutputTypes: []v1.Modality{v1.TextModality},
				Reasoning:   &boolTrue,
				ToolUsage:   &boolFalse,
				Reward:      &boolFalse,
				Languages:   []string{"en", "fr", "zh"},
			},
		},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var restored v1.Model
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	// Re-marshal to compare JSON output (avoids pointer comparison issues).
	data2, err := json.Marshal(restored)
	if err != nil {
		t.Fatalf("re-marshal failed: %v", err)
	}

	if string(data) != string(data2) {
		t.Errorf("round-trip JSON mismatch:\n  original:  %s\n  restored: %s", string(data), string(data2))
	}
}
