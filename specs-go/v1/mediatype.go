/*
 *     Copyright 2024 The CNAI Authors
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

const (
	// ArtifactTypeModelManifest specifies the media type for a model manifest.
	ArtifactTypeModelManifest = "application/vnd.cnai.model.manifest.v1+json"
)

const (
	// ArtifactTypeModelLayer is the media type used for layers referenced by the manifest.
	ArtifactTypeModelLayer = "application/vnd.cnai.model.layer.v1.tar"

	// ArtifactTypeModelLayerGzip is the media type used for gzipped layers
	// referenced by the manifest.
	ArtifactTypeModelLayerGzip = "application/vnd.cnai.model.layer.v1.tar+gzip"
)
