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

// Package compat provides compatibility utilities between ModelPack and other formats.
package compat

import "strings"

// Format represents the format type of a model artifact.
type Format int

const (
	FormatUnknown Format = iota
	FormatModelPack
	FormatDocker
)

// String returns the format name.
func (f Format) String() string {
	switch f {
	case FormatModelPack:
		return "modelpack"
	case FormatDocker:
		return "docker"
	default:
		return "unknown"
	}
}

// DetectFormat determines the artifact format based on its media type.
func DetectFormat(mediaType string) Format {
	switch {
	case strings.HasPrefix(mediaType, "application/vnd.cncf.model."):
		return FormatModelPack
	case strings.HasPrefix(mediaType, "application/vnd.docker.ai."):
		return FormatDocker
	default:
		return FormatUnknown
	}
}
