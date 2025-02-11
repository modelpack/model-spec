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
	// MediaTypeModelManifest specifies the media type for a model manifest.
	MediaTypeModelManifest = "application/vnd.cnai.model.manifest.v1+json"

	// MediaTypeModelConfig specifies the media type for a model configuration.
	MediaTypeModelConfig = "application/vnd.cnai.model.config.v1+json"

	// MediaTypeModelLayer is the media type used for layers referenced by the manifest.
	MediaTypeModelLayer = "application/vnd.cnai.model.layer.v1.tar"

	// MediaTypeModelLayerGzip is the media type used for gzipped layers
	// referenced by the manifest.
	MediaTypeModelLayerGzip = "application/vnd.cnai.model.layer.v1.tar+gzip"

	// MediaTypeModelDoc specifies the media type for model documentation, including documentation files like `README.md`, `LICENSE`, etc.
	MediaTypeModelDoc = "application/vnd.cnai.model.doc.v1.tar"

	// MediaTypeModelCode specifies the media type for a model code, includes code artifacts like scripts, code files etc.
	MediaTypeModelCode = "application/vnd.cnai.model.code.v1.tar"

	// MediaTypeModelDataset specifies the media type for a model dataset, includes datasets that may be needed for the lifecycle of AI/ML models.
	MediaTypeModelDataset = "application/vnd.cnai.model.dataset.v1.tar"
)
