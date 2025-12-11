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
	"strconv"
	"strings"
	"time"

	v1 "github.com/modelpack/model-spec/specs-go/v1"
	"github.com/opencontainers/go-digest"
)

// FromModelPack converts a ModelPack model to Docker format.
// Note: The Size field cannot be derived from ModelPack and will be set to "0".
func FromModelPack(m v1.Model) (*Config, error) {
	cfg := &Config{
		ModelConfig: ModelConfig{
			Format: m.Config.Format,
			Size:   "0",
		},
		Files: make([]File, 0, len(m.ModelFS.DiffIDs)),
	}

	if m.Descriptor.CreatedAt != nil {
		cfg.Descriptor = &Descriptor{
			CreatedAt: m.Descriptor.CreatedAt.Format(time.RFC3339),
		}
	}

	for _, diffID := range m.ModelFS.DiffIDs {
		cfg.Files = append(cfg.Files, File{
			DiffID: string(diffID),
			Type:   mediaTypeForFormat(m.Config.Format),
		})
	}

	if m.Config.ParamSize != "" {
		cfg.ModelConfig.GGUF = map[string]any{
			"parameter_count": formatParamSizeHuman(m.Config.ParamSize),
		}
	}

	if m.Config.Architecture != "" {
		if cfg.ModelConfig.GGUF == nil {
			cfg.ModelConfig.GGUF = make(map[string]any)
		}
		cfg.ModelConfig.GGUF["architecture"] = m.Config.Architecture
	}

	if m.Config.Quantization != "" {
		if cfg.ModelConfig.GGUF == nil {
			cfg.ModelConfig.GGUF = make(map[string]any)
		}
		cfg.ModelConfig.GGUF["quantization"] = m.Config.Quantization
	}

	return cfg, nil
}

// ToModelPack converts a Docker format config to ModelPack format.
func ToModelPack(cfg Config) (*v1.Model, error) {
	m := &v1.Model{
		Config: v1.ModelConfig{
			Format: cfg.ModelConfig.Format,
		},
		ModelFS: v1.ModelFS{
			Type:    "layers",
			DiffIDs: make([]digest.Digest, 0, len(cfg.Files)),
		},
	}

	if cfg.Descriptor != nil && cfg.Descriptor.CreatedAt != "" {
		t, err := time.Parse(time.RFC3339, cfg.Descriptor.CreatedAt)
		if err == nil {
			m.Descriptor.CreatedAt = &t
		}
	}

	for _, f := range cfg.Files {
		m.ModelFS.DiffIDs = append(m.ModelFS.DiffIDs, digest.Digest(f.DiffID))
	}

	if cfg.ModelConfig.GGUF != nil {
		if paramCount, ok := cfg.ModelConfig.GGUF["parameter_count"]; ok {
			if s, ok := paramCount.(string); ok {
				m.Config.ParamSize = parseParamSizeHuman(s)
			}
		}
		if arch, ok := cfg.ModelConfig.GGUF["architecture"]; ok {
			if s, ok := arch.(string); ok {
				m.Config.Architecture = s
			}
		}
		if quant, ok := cfg.ModelConfig.GGUF["quantization"]; ok {
			if s, ok := quant.(string); ok {
				m.Config.Quantization = s
			}
		}
	}

	return m, nil
}

// mediaTypeForFormat returns the Docker media type for a given format.
func mediaTypeForFormat(format string) string {
	switch strings.ToLower(format) {
	case FormatGGUF:
		return MediaTypeGGUF
	default:
		return "application/octet-stream"
	}
}

// formatParamSizeHuman converts "8b" to Docker format "8 B".
func formatParamSizeHuman(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}

	var numPart string
	var unitPart string
	lower := strings.ToLower(s)
	for i, c := range lower {
		if (c >= '0' && c <= '9') || c == '.' {
			numPart = s[:i+1]
		} else {
			unitPart = lower[i:]
			break
		}
	}

	if numPart == "" {
		return s
	}

	switch unitPart {
	case "t":
		return numPart + " T"
	case "b":
		return numPart + " B"
	case "m":
		return numPart + " M"
	case "k":
		return numPart + " K"
	default:
		return s
	}
}

// parseParamSizeHuman converts Docker format "8 B" to ModelPack format "8b".
func parseParamSizeHuman(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}

	parts := strings.Fields(s)
	if len(parts) != 2 {
		return strings.ToLower(strings.ReplaceAll(s, " ", ""))
	}

	num := parts[0]
	unit := strings.ToLower(parts[1])

	return num + unit
}

// parseParamSize parses parameter size string, e.g., "8b" -> 8000000000.
func parseParamSize(s string) int64 {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" {
		return 0
	}

	var numPart string
	var unitPart string
	for i, c := range s {
		if (c >= '0' && c <= '9') || c == '.' {
			numPart = s[:i+1]
		} else {
			unitPart = s[i:]
			break
		}
	}

	if numPart == "" {
		numPart = s
	}

	num, err := strconv.ParseFloat(numPart, 64)
	if err != nil {
		return 0
	}

	var multiplier int64 = 1
	switch strings.ToLower(unitPart) {
	case "t":
		multiplier = 1_000_000_000_000
	case "b":
		multiplier = 1_000_000_000
	case "m":
		multiplier = 1_000_000
	case "k":
		multiplier = 1_000
	}

	return int64(num * float64(multiplier))
}

// formatParamSize converts a number to human-readable format, e.g., 8000000000 -> "8b".
func formatParamSize(n int64) string {
	if n <= 0 {
		return ""
	}

	switch {
	case n >= 1_000_000_000_000:
		return strconv.FormatFloat(float64(n)/1_000_000_000_000, 'f', -1, 64) + "t"
	case n >= 1_000_000_000:
		return strconv.FormatFloat(float64(n)/1_000_000_000, 'f', -1, 64) + "b"
	case n >= 1_000_000:
		return strconv.FormatFloat(float64(n)/1_000_000, 'f', -1, 64) + "m"
	case n >= 1_000:
		return strconv.FormatFloat(float64(n)/1_000, 'f', -1, 64) + "k"
	default:
		return strconv.FormatInt(n, 10)
	}
}
