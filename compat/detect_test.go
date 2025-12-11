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

package compat

import "testing"

func TestDetectFormat(t *testing.T) {
	for i, tt := range []struct {
		mediaType string
		want      Format
	}{
		// ModelPack 格式
		{"application/vnd.cncf.model.config.v1+json", FormatModelPack},
		{"application/vnd.cncf.model.manifest.v1+json", FormatModelPack},
		{"application/vnd.cncf.model.weight.v1.raw", FormatModelPack},

		// Docker 格式
		{"application/vnd.docker.ai.model.config.v0.1+json", FormatDocker},
		{"application/vnd.docker.ai.gguf.v3", FormatDocker},
		{"application/vnd.docker.ai.license", FormatDocker},

		// 未知格式
		{"application/json", FormatUnknown},
		{"application/octet-stream", FormatUnknown},
		{"", FormatUnknown},
	} {
		got := DetectFormat(tt.mediaType)
		if got != tt.want {
			t.Errorf("test %d: DetectFormat(%q) = %v, want %v", i, tt.mediaType, got, tt.want)
		}
	}
}

func TestFormatString(t *testing.T) {
	tests := []struct {
		f    Format
		want string
	}{
		{FormatModelPack, "modelpack"},
		{FormatDocker, "docker"},
		{FormatUnknown, "unknown"},
	}

	for _, tt := range tests {
		if got := tt.f.String(); got != tt.want {
			t.Errorf("Format(%d).String() = %q, want %q", tt.f, got, tt.want)
		}
	}
}
