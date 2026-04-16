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

package v1

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	digest "github.com/opencontainers/go-digest"
)

func boolPtr(b bool) *bool {
	return &b
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func TestModelMarshalJSON(t *testing.T) {
	created := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
	m := Model{
		Descriptor: ModelDescriptor{
			Name:      "llama3-8b-instruct",
			Family:    "llama3",
			Version:   "1.0.0",
			CreatedAt: timePtr(created),
			Authors:   []string{"Meta"},
			Vendor:    "Meta",
		},
		ModelFS: ModelFS{
			Type:    "layers",
			DiffIDs: []digest.Digest{"sha256:abc123"},
		},
		Config: ModelConfig{
			Architecture: "transformer",
			Format:       "pytorch",
			ParamSize:    "8b",
			Precision:    "bf16",
		},
	}

	data, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("failed to marshal Model: %v", err)
	}

	var got Model
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("failed to unmarshal Model: %v", err)
	}

	if !reflect.DeepEqual(got, m) {
		t.Errorf("unmarshaled Model does not match original.\ngot:  %+v\nwant: %+v", got, m)
	}
}

func TestModelJSONFieldNames(t *testing.T) {
	m := Model{
		Descriptor: ModelDescriptor{Name: "test"},
		ModelFS:    ModelFS{Type: "layers", DiffIDs: []digest.Digest{"sha256:abc"}},
		Config:     ModelConfig{Architecture: "transformer"},
	}

	data, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("failed to unmarshal to map: %v", err)
	}

	for _, key := range []string{"descriptor", "modelfs", "config"} {
		if _, ok := raw[key]; !ok {
			t.Errorf("expected top-level JSON key %q, not found", key)
		}
	}
}

func TestModelOmitEmptyFields(t *testing.T) {
	m := Model{
		Descriptor: ModelDescriptor{Name: "test"},
		ModelFS:    ModelFS{Type: "layers", DiffIDs: []digest.Digest{"sha256:abc"}},
	}

	data, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("failed to unmarshal to map: %v", err)
	}

	// config has omitempty, but since it's a struct (not pointer), it will still appear
	// with its zero-value fields omitted inside it
	if configRaw, ok := raw["config"]; ok {
		var configMap map[string]json.RawMessage
		if err := json.Unmarshal(configRaw, &configMap); err != nil {
			t.Fatalf("failed to unmarshal config: %v", err)
		}
		if _, ok := configMap["architecture"]; ok {
			t.Error("empty architecture should be omitted")
		}
	}
}

func TestModelDescriptorMarshalJSON(t *testing.T) {
	created := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	d := ModelDescriptor{
		CreatedAt:   timePtr(created),
		Authors:     []string{"Alice", "Bob"},
		Family:      "qwen2",
		Name:        "qwen2-vl-72b-instruct",
		DocURL:      "https://example.com/docs",
		SourceURL:   "https://example.com/source",
		DatasetsURL: []string{"https://example.com/dataset1"},
		Version:     "2.0.0",
		Revision:    "abc123",
		Vendor:      "Alibaba",
		Licenses:    []string{"Apache-2.0"},
		Title:       "Qwen2 VL 72B Instruct",
		Description: "A multimodal model",
	}

	data, err := json.Marshal(d)
	if err != nil {
		t.Fatalf("failed to marshal ModelDescriptor: %v", err)
	}

	var got ModelDescriptor
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("failed to unmarshal ModelDescriptor: %v", err)
	}

	if !reflect.DeepEqual(got, d) {
		t.Errorf("unmarshaled ModelDescriptor does not match original.\ngot:  %+v\nwant: %+v", got, d)
	}
}

func TestModelDescriptorJSONFieldNames(t *testing.T) {
	created := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	d := ModelDescriptor{
		CreatedAt:   timePtr(created),
		Authors:     []string{"author"},
		Family:      "llama",
		Name:        "llama-7b",
		DocURL:      "https://example.com",
		SourceURL:   "https://example.com",
		DatasetsURL: []string{"https://example.com"},
		Version:     "1.0",
		Revision:    "rev1",
		Vendor:      "Meta",
		Licenses:    []string{"MIT"},
		Title:       "Title",
		Description: "Desc",
	}

	data, err := json.Marshal(d)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("failed to unmarshal to map: %v", err)
	}

	expectedKeys := []string{
		"createdAt", "authors", "family", "name", "docURL",
		"sourceURL", "datasetsURL", "version", "revision",
		"vendor", "licenses", "title", "description",
	}
	for _, key := range expectedKeys {
		if _, ok := raw[key]; !ok {
			t.Errorf("expected JSON key %q in ModelDescriptor, not found", key)
		}
	}
}

func TestModelDescriptorOmitEmpty(t *testing.T) {
	d := ModelDescriptor{Name: "minimal"}

	data, err := json.Marshal(d)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("failed to unmarshal to map: %v", err)
	}

	if _, ok := raw["name"]; !ok {
		t.Error("name should be present")
	}
	for _, key := range []string{"createdAt", "authors", "family", "vendor", "licenses"} {
		if _, ok := raw[key]; ok {
			t.Errorf("empty field %q should be omitted", key)
		}
	}
}

func TestModelConfigMarshalJSON(t *testing.T) {
	c := ModelConfig{
		Architecture: "transformer",
		Format:       "onnx",
		ParamSize:    "16b",
		Precision:    "fp16",
		Quantization: "awq",
		Capabilities: &ModelCapabilities{
			InputTypes:  []Modality{TextModality},
			OutputTypes: []Modality{TextModality},
			Reasoning:   boolPtr(true),
		},
	}

	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("failed to marshal ModelConfig: %v", err)
	}

	var got ModelConfig
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("failed to unmarshal ModelConfig: %v", err)
	}

	if !reflect.DeepEqual(got, c) {
		t.Errorf("unmarshaled ModelConfig does not match original.\ngot:  %+v\nwant: %+v", got, c)
	}
}

func TestModelConfigNilCapabilities(t *testing.T) {
	c := ModelConfig{Architecture: "cnn"}

	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if _, ok := raw["capabilities"]; ok {
		t.Error("nil capabilities should be omitted from JSON")
	}
}

func TestModelCapabilitiesMarshalJSON(t *testing.T) {
	cutoff := time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC)
	caps := ModelCapabilities{
		InputTypes:      []Modality{TextModality, ImageModality},
		OutputTypes:     []Modality{TextModality},
		KnowledgeCutoff: timePtr(cutoff),
		Reasoning:       boolPtr(true),
		ToolUsage:       boolPtr(false),
		Reward:          boolPtr(false),
		Languages:       []string{"en", "fr", "zh"},
	}

	data, err := json.Marshal(caps)
	if err != nil {
		t.Fatalf("failed to marshal ModelCapabilities: %v", err)
	}

	var got ModelCapabilities
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("failed to unmarshal ModelCapabilities: %v", err)
	}

	if !reflect.DeepEqual(got, caps) {
		t.Errorf("unmarshaled ModelCapabilities does not match original.\ngot:  %+v\nwant: %+v", got, caps)
	}
}

func TestModelCapabilitiesBoolPointerSemantics(t *testing.T) {
	tests := []struct {
		name     string
		caps     ModelCapabilities
		wantKeys []string
		omitKeys []string
	}{
		{
			name:     "all nil booleans omitted",
			caps:     ModelCapabilities{},
			omitKeys: []string{"reasoning", "toolUsage", "reward"},
		},
		{
			name:     "false booleans present",
			caps:     ModelCapabilities{Reasoning: boolPtr(false), ToolUsage: boolPtr(false), Reward: boolPtr(false)},
			wantKeys: []string{"reasoning", "toolUsage", "reward"},
		},
		{
			name:     "true booleans present",
			caps:     ModelCapabilities{Reasoning: boolPtr(true), ToolUsage: boolPtr(true), Reward: boolPtr(true)},
			wantKeys: []string{"reasoning", "toolUsage", "reward"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.caps)
			if err != nil {
				t.Fatalf("failed to marshal: %v", err)
			}

			var raw map[string]json.RawMessage
			if err := json.Unmarshal(data, &raw); err != nil {
				t.Fatalf("failed to unmarshal: %v", err)
			}

			for _, key := range tt.wantKeys {
				if _, ok := raw[key]; !ok {
					t.Errorf("expected key %q to be present", key)
				}
			}
			for _, key := range tt.omitKeys {
				if _, ok := raw[key]; ok {
					t.Errorf("expected key %q to be omitted", key)
				}
			}
		})
	}
}

func TestModelCapabilitiesJSONFieldNames(t *testing.T) {
	cutoff := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	caps := ModelCapabilities{
		InputTypes:      []Modality{TextModality},
		OutputTypes:     []Modality{TextModality},
		KnowledgeCutoff: timePtr(cutoff),
		Reasoning:       boolPtr(true),
		ToolUsage:       boolPtr(true),
		Reward:          boolPtr(true),
		Languages:       []string{"en"},
	}

	data, err := json.Marshal(caps)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	expectedKeys := []string{
		"inputTypes", "outputTypes", "knowledgeCutoff",
		"reasoning", "toolUsage", "reward", "languages",
	}
	for _, key := range expectedKeys {
		if _, ok := raw[key]; !ok {
			t.Errorf("expected JSON key %q in ModelCapabilities, not found", key)
		}
	}
}

func TestModelFSMarshalJSON(t *testing.T) {
	fs := ModelFS{
		Type: "layers",
		DiffIDs: []digest.Digest{
			"sha256:abc123",
			"sha256:def456",
		},
	}

	data, err := json.Marshal(fs)
	if err != nil {
		t.Fatalf("failed to marshal ModelFS: %v", err)
	}

	var got ModelFS
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("failed to unmarshal ModelFS: %v", err)
	}

	if !reflect.DeepEqual(got, fs) {
		t.Errorf("unmarshaled ModelFS does not match original.\ngot:  %+v\nwant: %+v", got, fs)
	}
}

func TestModelFSJSONFieldNames(t *testing.T) {
	fs := ModelFS{
		Type:    "layers",
		DiffIDs: []digest.Digest{"sha256:abc"},
	}

	data, err := json.Marshal(fs)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	for _, key := range []string{"type", "diffIds"} {
		if _, ok := raw[key]; !ok {
			t.Errorf("expected JSON key %q in ModelFS, not found", key)
		}
	}
}

func TestModalityConstants(t *testing.T) {
	tests := []struct {
		modality Modality
		want     string
	}{
		{TextModality, "text"},
		{ImageModality, "image"},
		{AudioModality, "audio"},
		{VideoModality, "video"},
		{EmbeddingModality, "embedding"},
		{OtherModality, "other"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if string(tt.modality) != tt.want {
				t.Errorf("Modality = %q, want %q", tt.modality, tt.want)
			}
		})
	}
}

func TestModalityMarshalJSON(t *testing.T) {
	caps := ModelCapabilities{
		InputTypes:  []Modality{TextModality, ImageModality, AudioModality},
		OutputTypes: []Modality{TextModality, EmbeddingModality},
	}

	data, err := json.Marshal(caps)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var got ModelCapabilities
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if len(got.InputTypes) != 3 {
		t.Fatalf("inputTypes len = %d, want 3", len(got.InputTypes))
	}
	if got.InputTypes[0] != TextModality || got.InputTypes[1] != ImageModality || got.InputTypes[2] != AudioModality {
		t.Errorf("inputTypes = %v, want [text image audio]", got.InputTypes)
	}
	if len(got.OutputTypes) != 2 || got.OutputTypes[1] != EmbeddingModality {
		t.Errorf("outputTypes = %v, want [text embedding]", got.OutputTypes)
	}
}

func TestModelUnmarshalFromJSON(t *testing.T) {
	input := `{
		"descriptor": {
			"name": "gpt2-xl",
			"family": "gpt2",
			"version": "1.0",
			"authors": ["OpenAI"],
			"vendor": "OpenAI",
			"licenses": ["MIT"],
			"createdAt": "2025-01-15T10:30:00Z"
		},
		"modelfs": {
			"type": "layers",
			"diffIds": ["sha256:abcdef1234567890"]
		},
		"config": {
			"architecture": "transformer",
			"format": "pytorch",
			"paramSize": "1.5b",
			"precision": "fp32",
			"capabilities": {
				"inputTypes": ["text"],
				"outputTypes": ["text"],
				"reasoning": false,
				"toolUsage": false,
				"languages": ["en"]
			}
		}
	}`

	var m Model
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if m.Descriptor.Name != "gpt2-xl" {
		t.Errorf("name = %q, want %q", m.Descriptor.Name, "gpt2-xl")
	}
	if m.Descriptor.Family != "gpt2" {
		t.Errorf("family = %q, want %q", m.Descriptor.Family, "gpt2")
	}
	if m.Descriptor.CreatedAt == nil {
		t.Fatal("createdAt should not be nil")
	}
	if m.Descriptor.CreatedAt.Year() != 2025 {
		t.Errorf("createdAt year = %d, want 2025", m.Descriptor.CreatedAt.Year())
	}
	if m.Config.ParamSize != "1.5b" {
		t.Errorf("paramSize = %q, want %q", m.Config.ParamSize, "1.5b")
	}
	if m.Config.Capabilities == nil {
		t.Fatal("capabilities should not be nil")
	}
	if *m.Config.Capabilities.Reasoning != false {
		t.Error("reasoning should be false")
	}
	if len(m.Config.Capabilities.Languages) != 1 || m.Config.Capabilities.Languages[0] != "en" {
		t.Errorf("languages = %v, want [en]", m.Config.Capabilities.Languages)
	}
}

func TestModelRoundTrip(t *testing.T) {
	cutoff := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	created := time.Date(2025, 3, 1, 12, 0, 0, 0, time.UTC)
	original := Model{
		Descriptor: ModelDescriptor{
			CreatedAt:   timePtr(created),
			Authors:     []string{"Research Team"},
			Family:      "llama3",
			Name:        "llama3-70b",
			Version:     "3.0",
			Vendor:      "Meta",
			Licenses:    []string{"Llama-3"},
			Description: "Large language model",
		},
		ModelFS: ModelFS{
			Type: "layers",
			DiffIDs: []digest.Digest{
				"sha256:1111111111111111111111111111111111111111111111111111111111111111",
				"sha256:2222222222222222222222222222222222222222222222222222222222222222",
			},
		},
		Config: ModelConfig{
			Architecture: "transformer",
			Format:       "safetensors",
			ParamSize:    "70b",
			Precision:    "bf16",
			Quantization: "gptq",
			Capabilities: &ModelCapabilities{
				InputTypes:      []Modality{TextModality},
				OutputTypes:     []Modality{TextModality},
				KnowledgeCutoff: timePtr(cutoff),
				Reasoning:       boolPtr(true),
				ToolUsage:       boolPtr(true),
				Reward:          boolPtr(false),
				Languages:       []string{"en", "fr", "de", "zh"},
			},
		},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var roundTripped Model
	if err := json.Unmarshal(data, &roundTripped); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	// Re-marshal and compare JSON output
	data2, err := json.Marshal(roundTripped)
	if err != nil {
		t.Fatalf("failed to re-marshal: %v", err)
	}

	if string(data) != string(data2) {
		t.Errorf("round-trip JSON mismatch:\n  first:  %s\n  second: %s", data, data2)
	}
}
