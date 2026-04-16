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
	"strings"
	"testing"
)

func TestMediaTypeConstants(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{"ArtifactTypeModelManifest", ArtifactTypeModelManifest, "application/vnd.cncf.model.manifest.v1+json"},
		{"MediaTypeModelConfig", MediaTypeModelConfig, "application/vnd.cncf.model.config.v1+json"},
		{"MediaTypeModelWeightRaw", MediaTypeModelWeightRaw, "application/vnd.cncf.model.weight.v1.raw"},
		{"MediaTypeModelWeight", MediaTypeModelWeight, "application/vnd.cncf.model.weight.v1.tar"},
		{"MediaTypeModelWeightGzip", MediaTypeModelWeightGzip, "application/vnd.cncf.model.weight.v1.tar+gzip"},
		{"MediaTypeModelWeightZstd", MediaTypeModelWeightZstd, "application/vnd.cncf.model.weight.v1.tar+zstd"},
		{"MediaTypeModelWeightConfigRaw", MediaTypeModelWeightConfigRaw, "application/vnd.cncf.model.weight.config.v1.raw"},
		{"MediaTypeModelWeightConfig", MediaTypeModelWeightConfig, "application/vnd.cncf.model.weight.config.v1.tar"},
		{"MediaTypeModelWeightConfigGzip", MediaTypeModelWeightConfigGzip, "application/vnd.cncf.model.weight.config.v1.tar+gzip"},
		{"MediaTypeModelWeightConfigZstd", MediaTypeModelWeightConfigZstd, "application/vnd.cncf.model.weight.config.v1.tar+zstd"},
		{"MediaTypeModelDocRaw", MediaTypeModelDocRaw, "application/vnd.cncf.model.doc.v1.raw"},
		{"MediaTypeModelDoc", MediaTypeModelDoc, "application/vnd.cncf.model.doc.v1.tar"},
		{"MediaTypeModelDocGzip", MediaTypeModelDocGzip, "application/vnd.cncf.model.doc.v1.tar+gzip"},
		{"MediaTypeModelDocZstd", MediaTypeModelDocZstd, "application/vnd.cncf.model.doc.v1.tar+zstd"},
		{"MediaTypeModelCodeRaw", MediaTypeModelCodeRaw, "application/vnd.cncf.model.code.v1.raw"},
		{"MediaTypeModelCode", MediaTypeModelCode, "application/vnd.cncf.model.code.v1.tar"},
		{"MediaTypeModelCodeGzip", MediaTypeModelCodeGzip, "application/vnd.cncf.model.code.v1.tar+gzip"},
		{"MediaTypeModelCodeZstd", MediaTypeModelCodeZstd, "application/vnd.cncf.model.code.v1.tar+zstd"},
		{"MediaTypeModelDatasetRaw", MediaTypeModelDatasetRaw, "application/vnd.cncf.model.dataset.v1.raw"},
		{"MediaTypeModelDataset", MediaTypeModelDataset, "application/vnd.cncf.model.dataset.v1.tar"},
		{"MediaTypeModelDatasetGzip", MediaTypeModelDatasetGzip, "application/vnd.cncf.model.dataset.v1.tar+gzip"},
		{"MediaTypeModelDatasetZstd", MediaTypeModelDatasetZstd, "application/vnd.cncf.model.dataset.v1.tar+zstd"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != tt.want {
				t.Errorf("%s = %q, want %q", tt.name, tt.value, tt.want)
			}
		})
	}
}

func TestMediaTypePrefix(t *testing.T) {
	mediaTypes := []string{
		MediaTypeModelConfig,
		MediaTypeModelWeight, MediaTypeModelWeightRaw, MediaTypeModelWeightGzip, MediaTypeModelWeightZstd,
		MediaTypeModelWeightConfig, MediaTypeModelWeightConfigRaw, MediaTypeModelWeightConfigGzip, MediaTypeModelWeightConfigZstd,
		MediaTypeModelDoc, MediaTypeModelDocRaw, MediaTypeModelDocGzip, MediaTypeModelDocZstd,
		MediaTypeModelCode, MediaTypeModelCodeRaw, MediaTypeModelCodeGzip, MediaTypeModelCodeZstd,
		MediaTypeModelDataset, MediaTypeModelDatasetRaw, MediaTypeModelDatasetGzip, MediaTypeModelDatasetZstd,
	}

	for _, mt := range mediaTypes {
		if !strings.HasPrefix(mt, "application/vnd.cncf.model.") {
			t.Errorf("media type %q does not have expected prefix", mt)
		}
	}
}

func TestMediaTypeCompressionSuffixes(t *testing.T) {
	gzipTypes := []string{
		MediaTypeModelWeightGzip,
		MediaTypeModelWeightConfigGzip,
		MediaTypeModelDocGzip,
		MediaTypeModelCodeGzip,
		MediaTypeModelDatasetGzip,
	}
	for _, mt := range gzipTypes {
		if !strings.HasSuffix(mt, "+gzip") {
			t.Errorf("gzip media type %q does not end with +gzip", mt)
		}
	}

	zstdTypes := []string{
		MediaTypeModelWeightZstd,
		MediaTypeModelWeightConfigZstd,
		MediaTypeModelDocZstd,
		MediaTypeModelCodeZstd,
		MediaTypeModelDatasetZstd,
	}
	for _, mt := range zstdTypes {
		if !strings.HasSuffix(mt, "+zstd") {
			t.Errorf("zstd media type %q does not end with +zstd", mt)
		}
	}

	rawTypes := []string{
		MediaTypeModelWeightRaw,
		MediaTypeModelWeightConfigRaw,
		MediaTypeModelDocRaw,
		MediaTypeModelCodeRaw,
		MediaTypeModelDatasetRaw,
	}
	for _, mt := range rawTypes {
		if !strings.HasSuffix(mt, ".raw") {
			t.Errorf("raw media type %q does not end with .raw", mt)
		}
	}

	tarTypes := []string{
		MediaTypeModelWeight,
		MediaTypeModelWeightConfig,
		MediaTypeModelDoc,
		MediaTypeModelCode,
		MediaTypeModelDataset,
	}
	for _, mt := range tarTypes {
		if !strings.HasSuffix(mt, ".tar") {
			t.Errorf("tar media type %q does not end with .tar", mt)
		}
	}
}

func TestMediaTypeUniqueness(t *testing.T) {
	allTypes := []string{
		ArtifactTypeModelManifest,
		MediaTypeModelConfig,
		MediaTypeModelWeightRaw, MediaTypeModelWeight, MediaTypeModelWeightGzip, MediaTypeModelWeightZstd,
		MediaTypeModelWeightConfigRaw, MediaTypeModelWeightConfig, MediaTypeModelWeightConfigGzip, MediaTypeModelWeightConfigZstd,
		MediaTypeModelDocRaw, MediaTypeModelDoc, MediaTypeModelDocGzip, MediaTypeModelDocZstd,
		MediaTypeModelCodeRaw, MediaTypeModelCode, MediaTypeModelCodeGzip, MediaTypeModelCodeZstd,
		MediaTypeModelDatasetRaw, MediaTypeModelDataset, MediaTypeModelDatasetGzip, MediaTypeModelDatasetZstd,
	}

	seen := make(map[string]bool)
	for _, mt := range allTypes {
		if seen[mt] {
			t.Errorf("duplicate media type: %q", mt)
		}
		seen[mt] = true
	}
}
