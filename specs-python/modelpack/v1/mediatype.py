#     Copyright 2025 The CNCF ModelPack Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""Media type constants matching specs-go/v1/mediatype.go."""

# Artifact type for a model manifest.
ARTIFACT_TYPE_MODEL_MANIFEST = "application/vnd.cncf.model.manifest.v1+json"

# Media type for a model configuration.
MEDIA_TYPE_MODEL_CONFIG = "application/vnd.cncf.model.config.v1+json"

# Model weight media types.
MEDIA_TYPE_MODEL_WEIGHT_RAW = "application/vnd.cncf.model.weight.v1.raw"
MEDIA_TYPE_MODEL_WEIGHT = "application/vnd.cncf.model.weight.v1.tar"
MEDIA_TYPE_MODEL_WEIGHT_GZIP = "application/vnd.cncf.model.weight.v1.tar+gzip"
MEDIA_TYPE_MODEL_WEIGHT_ZSTD = "application/vnd.cncf.model.weight.v1.tar+zstd"

# Model weight config media types.
MEDIA_TYPE_MODEL_WEIGHT_CONFIG_RAW = "application/vnd.cncf.model.weight.config.v1.raw"
MEDIA_TYPE_MODEL_WEIGHT_CONFIG = "application/vnd.cncf.model.weight.config.v1.tar"
MEDIA_TYPE_MODEL_WEIGHT_CONFIG_GZIP = (
    "application/vnd.cncf.model.weight.config.v1.tar+gzip"
)
MEDIA_TYPE_MODEL_WEIGHT_CONFIG_ZSTD = (
    "application/vnd.cncf.model.weight.config.v1.tar+zstd"
)

# Model documentation media types.
MEDIA_TYPE_MODEL_DOC_RAW = "application/vnd.cncf.model.doc.v1.raw"
MEDIA_TYPE_MODEL_DOC = "application/vnd.cncf.model.doc.v1.tar"
MEDIA_TYPE_MODEL_DOC_GZIP = "application/vnd.cncf.model.doc.v1.tar+gzip"
MEDIA_TYPE_MODEL_DOC_ZSTD = "application/vnd.cncf.model.doc.v1.tar+zstd"

# Model code media types.
MEDIA_TYPE_MODEL_CODE_RAW = "application/vnd.cncf.model.code.v1.raw"
MEDIA_TYPE_MODEL_CODE = "application/vnd.cncf.model.code.v1.tar"
MEDIA_TYPE_MODEL_CODE_GZIP = "application/vnd.cncf.model.code.v1.tar+gzip"
MEDIA_TYPE_MODEL_CODE_ZSTD = "application/vnd.cncf.model.code.v1.tar+zstd"

# Model dataset media types.
MEDIA_TYPE_MODEL_DATASET_RAW = "application/vnd.cncf.model.dataset.v1.raw"
MEDIA_TYPE_MODEL_DATASET = "application/vnd.cncf.model.dataset.v1.tar"
MEDIA_TYPE_MODEL_DATASET_GZIP = "application/vnd.cncf.model.dataset.v1.tar+gzip"
MEDIA_TYPE_MODEL_DATASET_ZSTD = "application/vnd.cncf.model.dataset.v1.tar+zstd"
