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

// Package v0 provides type definitions for Docker model-spec v0.1 format.
// See: https://github.com/docker/model-spec/blob/main/config.md
package v0

// FormatGGUF is the primary format supported by Docker model-spec.
const FormatGGUF = "gguf"

// Media types for Docker AI model artifacts.
// See: https://github.com/docker/model-spec/blob/main/spec.md
const (
	// MediaTypeConfig is the config media type.
	MediaTypeConfig = "application/vnd.docker.ai.model.config.v0.1+json"

	// MediaTypeGGUF is the GGUF model file media type.
	MediaTypeGGUF = "application/vnd.docker.ai.gguf.v3"

	// MediaTypeLoRA is the LoRA adapter media type.
	MediaTypeLoRA = "application/vnd.docker.ai.gguf.v3.lora"

	// MediaTypeMMProj is the multimodal projector media type.
	MediaTypeMMProj = "application/vnd.docker.ai.gguf.v3.mmproj"

	// MediaTypeLicense is the license file media type.
	MediaTypeLicense = "application/vnd.docker.ai.license"

	// MediaTypeChatTemplate is the Jinja chat template media type.
	MediaTypeChatTemplate = "application/vnd.docker.ai.chat.template.jinja"
)

// Config is the root structure for Docker AI model config.
// Media type: application/vnd.docker.ai.model.config.v0.1+json
type Config struct {
	Descriptor  *Descriptor `json:"descriptor,omitempty"`
	ModelConfig ModelConfig `json:"config"`
	Files       []File      `json:"files"`
}

// Descriptor contains provenance information about the artifact.
type Descriptor struct {
	CreatedAt string `json:"createdAt,omitempty"`
}

// ModelConfig contains technical metadata about the model.
type ModelConfig struct {
	// The packaging format (e.g., "gguf")
	Format string `json:"format"`

	// The packaging format version
	FormatVersion string `json:"format_version,omitempty"`

	// The total size of the model in bytes
	Size string `json:"size"`

	// Format-specific metadata (for GGUF models)
	GGUF map[string]any `json:"gguf,omitempty"`
}

// File describes a single file that makes up the model.
type File struct {
	// The file digest as <alg>:<hash>
	DiffID string `json:"diffID"`

	// The media type indicating how to interpret the file
	Type string `json:"type"`
}
