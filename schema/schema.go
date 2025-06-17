/*
 *     Copyright 2025 The CNAI Authors
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

package schema

import (
	"embed"
	"net/http"

	v1 "github.com/modelpack/model-spec/specs-go/v1"
)

// Media types for the model-spec related formats
const (
	ValidatorMediaTypeModelConfig Validator = v1.MediaTypeModelConfig
)

var (
	// specFS stores the embedded http.FileSystem having the model-spec related JSON schema files in root "/".
	//go:embed *.json
	specFS embed.FS

	// specs maps model-spec schema media types to schema files.
	specs = map[Validator]string{
		ValidatorMediaTypeModelConfig: "config-schema.json",
	}

	// specURLs lists the various URLs a given spec may be known by.
	specURLs = map[string][]string{
		"config-schema.json": {
			"https://github.com/modelpack/model-spec/config",
		},
	}
)

// FileSystem returns an in-memory filesystem including the schema files.
// The schema files are located at the root directory.
func FileSystem() http.FileSystem {
	return http.FS(specFS)
}
