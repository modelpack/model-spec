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
	"time"

	v1 "github.com/modelpack/model-spec/specs-go/v1"
	"github.com/opencontainers/go-digest"
)

func TestFromModelPack(t *testing.T) {
	now := time.Now()

	for i, tt := range []struct {
		name    string
		input   v1.Model
		wantErr bool
		checkFn func(t *testing.T, got *Config)
	}{
		// basic conversion with only required fields
		{
			name: "basic model",
			input: v1.Model{
				Descriptor: v1.ModelDescriptor{},
				Config:     v1.ModelConfig{Format: FormatGGUF},
				ModelFS: v1.ModelFS{
					Type:    "layers",
					DiffIDs: []digest.Digest{"sha256:abc123"},
				},
			},
			wantErr: false,
			checkFn: func(t *testing.T, got *Config) {
				if got.ModelConfig.Format != FormatGGUF {
					t.Errorf("format: got %q, want %q", got.ModelConfig.Format, FormatGGUF)
				}
				if len(got.Files) != 1 {
					t.Errorf("files count: got %d, want 1", len(got.Files))
				}
				if got.ModelConfig.Size != "0" {
					t.Errorf("size: got %q, want %q", got.ModelConfig.Size, "0")
				}
			},
		},
		// conversion with timestamp
		{
			name: "with timestamp",
			input: v1.Model{
				Descriptor: v1.ModelDescriptor{
					CreatedAt: &now,
				},
				Config: v1.ModelConfig{Format: FormatGGUF},
				ModelFS: v1.ModelFS{
					Type:    "layers",
					DiffIDs: []digest.Digest{"sha256:def456"},
				},
			},
			wantErr: false,
			checkFn: func(t *testing.T, got *Config) {
				if got.Descriptor == nil {
					t.Fatal("descriptor should not be nil")
				}
				if got.Descriptor.CreatedAt == "" {
					t.Error("createdAt should not be empty")
				}
			},
		},
		// conversion with multiple layers
		{
			name: "multiple layers",
			input: v1.Model{
				Config: v1.ModelConfig{Format: FormatGGUF},
				ModelFS: v1.ModelFS{
					Type: "layers",
					DiffIDs: []digest.Digest{
						"sha256:layer1",
						"sha256:layer2",
						"sha256:layer3",
					},
				},
			},
			wantErr: false,
			checkFn: func(t *testing.T, got *Config) {
				if len(got.Files) != 3 {
					t.Errorf("files count: got %d, want 3", len(got.Files))
				}
			},
		},
		// conversion with paramSize
		{
			name: "with paramSize",
			input: v1.Model{
				Config: v1.ModelConfig{
					Format:    FormatGGUF,
					ParamSize: "8b",
				},
				ModelFS: v1.ModelFS{
					Type:    "layers",
					DiffIDs: []digest.Digest{"sha256:abc123"},
				},
			},
			wantErr: false,
			checkFn: func(t *testing.T, got *Config) {
				if got.ModelConfig.GGUF == nil {
					t.Fatal("gguf should not be nil")
				}
				paramCount, ok := got.ModelConfig.GGUF["parameter_count"]
				if !ok {
					t.Fatal("parameter_count should exist in gguf")
				}
				if paramCount != "8 B" {
					t.Errorf("parameter_count: got %v, want %q", paramCount, "8 B")
				}
			},
		},
		// conversion with architecture and quantization
		{
			name: "with architecture and quantization",
			input: v1.Model{
				Config: v1.ModelConfig{
					Format:       FormatGGUF,
					Architecture: "llama",
					Quantization: "Q4_0",
				},
				ModelFS: v1.ModelFS{
					Type:    "layers",
					DiffIDs: []digest.Digest{"sha256:abc123"},
				},
			},
			wantErr: false,
			checkFn: func(t *testing.T, got *Config) {
				if got.ModelConfig.GGUF == nil {
					t.Fatal("gguf should not be nil")
				}
				if got.ModelConfig.GGUF["architecture"] != "llama" {
					t.Errorf("architecture: got %v, want %q", got.ModelConfig.GGUF["architecture"], "llama")
				}
				if got.ModelConfig.GGUF["quantization"] != "Q4_0" {
					t.Errorf("quantization: got %v, want %q", got.ModelConfig.GGUF["quantization"], "Q4_0")
				}
			},
		},
	} {
		got, err := FromModelPack(tt.input)

		if (err != nil) != tt.wantErr {
			t.Errorf("test %d (%s): wantErr=%v, got err=%v", i, tt.name, tt.wantErr, err)
			continue
		}

		if !tt.wantErr && tt.checkFn != nil {
			tt.checkFn(t, got)
		}
	}
}

func TestToModelPack(t *testing.T) {
	for i, tt := range []struct {
		name    string
		input   Config
		wantErr bool
		checkFn func(t *testing.T, got *v1.Model)
	}{
		// basic conversion
		{
			name: "basic config",
			input: Config{
				ModelConfig: ModelConfig{Format: FormatGGUF, Size: "1000"},
				Files: []File{
					{DiffID: "sha256:abc123", Type: "application/vnd.docker.ai.gguf.v3"},
				},
			},
			wantErr: false,
			checkFn: func(t *testing.T, got *v1.Model) {
				if got.Config.Format != FormatGGUF {
					t.Errorf("format: got %q, want %q", got.Config.Format, FormatGGUF)
				}
				if len(got.ModelFS.DiffIDs) != 1 {
					t.Errorf("diffIDs count: got %d, want 1", len(got.ModelFS.DiffIDs))
				}
			},
		},
		// conversion with descriptor
		{
			name: "with descriptor",
			input: Config{
				Descriptor:  &Descriptor{CreatedAt: "2025-01-01T00:00:00Z"},
				ModelConfig: ModelConfig{Format: FormatGGUF, Size: "1000"},
				Files:       []File{},
			},
			wantErr: false,
			checkFn: func(t *testing.T, got *v1.Model) {
				if got.Descriptor.CreatedAt == nil {
					t.Error("createdAt should not be nil")
				}
			},
		},
		// conversion with gguf metadata
		{
			name: "with gguf metadata",
			input: Config{
				ModelConfig: ModelConfig{
					Format: FormatGGUF,
					Size:   "1000",
					GGUF: map[string]any{
						"parameter_count": "8 B",
						"architecture":    "llama",
						"quantization":    "Q4_0",
					},
				},
				Files: []File{},
			},
			wantErr: false,
			checkFn: func(t *testing.T, got *v1.Model) {
				if got.Config.ParamSize != "8b" {
					t.Errorf("paramSize: got %q, want %q", got.Config.ParamSize, "8b")
				}
				if got.Config.Architecture != "llama" {
					t.Errorf("architecture: got %q, want %q", got.Config.Architecture, "llama")
				}
				if got.Config.Quantization != "Q4_0" {
					t.Errorf("quantization: got %q, want %q", got.Config.Quantization, "Q4_0")
				}
			},
		},
	} {
		got, err := ToModelPack(tt.input)

		if (err != nil) != tt.wantErr {
			t.Errorf("test %d (%s): wantErr=%v, got err=%v", i, tt.name, tt.wantErr, err)
			continue
		}

		if !tt.wantErr && tt.checkFn != nil {
			tt.checkFn(t, got)
		}
	}
}

func TestParseParamSize(t *testing.T) {
	for i, tt := range []struct {
		input string
		want  int64
	}{
		{"8b", 8_000_000_000},
		{"8B", 8_000_000_000},
		{"70b", 70_000_000_000},
		{"1.5b", 1_500_000_000},
		{"7m", 7_000_000},
		{"100k", 100_000},
		{"1t", 1_000_000_000_000},
		{"", 0},
		{"invalid", 0},
	} {
		got := parseParamSize(tt.input)
		if got != tt.want {
			t.Errorf("test %d: parseParamSize(%q) = %d, want %d", i, tt.input, got, tt.want)
		}
	}
}

func TestFormatParamSizeHuman(t *testing.T) {
	for i, tt := range []struct {
		input string
		want  string
	}{
		{"8b", "8 B"},
		{"8B", "8 B"},
		{"70b", "70 B"},
		{"1.5b", "1.5 B"},
		{"7m", "7 M"},
		{"100k", "100 K"},
		{"1t", "1 T"},
		{"", ""},
		{"invalid", "invalid"},
	} {
		got := formatParamSizeHuman(tt.input)
		if got != tt.want {
			t.Errorf("test %d: formatParamSizeHuman(%q) = %q, want %q", i, tt.input, got, tt.want)
		}
	}
}

func TestParseParamSizeHuman(t *testing.T) {
	for i, tt := range []struct {
		input string
		want  string
	}{
		{"8 B", "8b"},
		{"70 B", "70b"},
		{"1.5 B", "1.5b"},
		{"7 M", "7m"},
		{"100 K", "100k"},
		{"1 T", "1t"},
		{"", ""},
		{"8B", "8b"},
	} {
		got := parseParamSizeHuman(tt.input)
		if got != tt.want {
			t.Errorf("test %d: parseParamSizeHuman(%q) = %q, want %q", i, tt.input, got, tt.want)
		}
	}
}

func TestJSONSerialization(t *testing.T) {
	cfg := Config{
		Descriptor: &Descriptor{CreatedAt: "2025-01-01T00:00:00Z"},
		ModelConfig: ModelConfig{
			Format: FormatGGUF,
			Size:   "635992801",
			GGUF: map[string]any{
				"architecture":    "llama",
				"parameter_count": "1.10 B",
				"quantization":    "Q4_0",
			},
		},
		Files: []File{
			{
				DiffID: "sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
				Type:   "application/vnd.docker.ai.gguf.v3",
			},
		},
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var unmarshaled Config
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if unmarshaled.ModelConfig.Format != cfg.ModelConfig.Format {
		t.Errorf("format mismatch: got %q, want %q", unmarshaled.ModelConfig.Format, cfg.ModelConfig.Format)
	}
	if unmarshaled.ModelConfig.Size != cfg.ModelConfig.Size {
		t.Errorf("size mismatch: got %q, want %q", unmarshaled.ModelConfig.Size, cfg.ModelConfig.Size)
	}
	if unmarshaled.ModelConfig.GGUF["architecture"] != cfg.ModelConfig.GGUF["architecture"] {
		t.Errorf("architecture mismatch")
	}
}
