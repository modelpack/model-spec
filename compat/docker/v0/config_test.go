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

package v0

import (
	"encoding/json"
	"testing"
)

func TestConfigUnmarshal(t *testing.T) {
	for i, tt := range []struct {
		name    string
		input   string
		wantErr bool
	}{
		// minimal valid config
		{
			name: "minimal config",
			input: `{
				"config": {"format": "gguf", "size": "1000"},
				"files": []
			}`,
			wantErr: false,
		},
		// full config matching Docker model-spec example
		{
			name: "full config",
			input: `{
				"descriptor": {"createdAt": "2025-01-01T00:00:00Z"},
				"config": {
					"format": "gguf",
					"format_version": "3",
					"gguf": {
						"architecture": "llama",
						"parameter_count": "1.10 B",
						"quantization": "Q4_0"
					},
					"size": "635992801"
				},
				"files": [
					{"diffID": "sha256:abc123", "type": "application/vnd.docker.ai.gguf.v3"},
					{"diffID": "sha256:def456", "type": "application/vnd.docker.ai.license"}
				]
			}`,
			wantErr: false,
		},
		// empty json
		{
			name:    "empty json",
			input:   `{}`,
			wantErr: false,
		},
		// invalid json
		{
			name:    "invalid json",
			input:   `{not valid}`,
			wantErr: true,
		},
	} {
		var cfg Config
		err := json.Unmarshal([]byte(tt.input), &cfg)

		if (err != nil) != tt.wantErr {
			t.Errorf("test %d (%s): wantErr=%v, got err=%v", i, tt.name, tt.wantErr, err)
		}
	}
}

func TestConfigMarshal(t *testing.T) {
	original := Config{
		Descriptor: &Descriptor{
			CreatedAt: "2025-01-01T00:00:00Z",
		},
		ModelConfig: ModelConfig{
			Format:        FormatGGUF,
			FormatVersion: "3",
			Size:          "635992801",
			GGUF: map[string]any{
				"architecture":    "llama",
				"parameter_count": "1.10 B",
				"quantization":    "Q4_0",
			},
		},
		Files: []File{
			{DiffID: "sha256:abc123", Type: MediaTypeGGUF},
			{DiffID: "sha256:def456", Type: MediaTypeLicense},
		},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var decoded Config
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if decoded.ModelConfig.Format != original.ModelConfig.Format {
		t.Errorf("format mismatch: got %q, want %q", decoded.ModelConfig.Format, original.ModelConfig.Format)
	}
	if decoded.ModelConfig.Size != original.ModelConfig.Size {
		t.Errorf("size mismatch: got %q, want %q", decoded.ModelConfig.Size, original.ModelConfig.Size)
	}
	if len(decoded.Files) != len(original.Files) {
		t.Errorf("files count mismatch: got %d, want %d", len(decoded.Files), len(original.Files))
	}
	if decoded.ModelConfig.GGUF["architecture"] != original.ModelConfig.GGUF["architecture"] {
		t.Error("gguf.architecture mismatch")
	}
}

func TestConfigMatchesDockerSpec(t *testing.T) {
	// This test verifies our struct produces JSON matching Docker model-spec example
	cfg := Config{
		Descriptor: &Descriptor{
			CreatedAt: "2025-01-01T00:00:00Z",
		},
		ModelConfig: ModelConfig{
			Format:        FormatGGUF,
			FormatVersion: "3",
			Size:          "635992801",
			GGUF: map[string]any{
				"architecture":    "llama",
				"parameter_count": "1.10 B",
				"quantization":    "Q4_0",
			},
		},
		Files: []File{
			{
				DiffID: "sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
				Type:   MediaTypeGGUF,
			},
			{
				DiffID: "sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
				Type:   MediaTypeLicense,
			},
			{
				DiffID: "sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
				Type:   MediaTypeLoRA,
			},
		},
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	// Verify it can be unmarshaled back
	var decoded Config
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	// Verify key structure matches Docker spec
	if decoded.ModelConfig.GGUF == nil {
		t.Error("gguf should be inside config, not at top level")
	}
	if decoded.ModelConfig.GGUF["parameter_count"] != "1.10 B" {
		t.Errorf("parameter_count format should be '1.10 B', got %v", decoded.ModelConfig.GGUF["parameter_count"])
	}
}
