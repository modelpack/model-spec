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

"""ModelPack Python SDK - CNCF standard for packaging and distributing AI models."""

from modelpack.v1.config import (
    Model,
    ModelCapabilities,
    ModelConfig,
    ModelDescriptor,
    ModelFS,
    Modality,
)
from modelpack.v1.annotations import (
    ANNOTATION_FILEPATH,
    ANNOTATION_FILE_METADATA,
    ANNOTATION_MEDIA_TYPE_UNTESTED,
    FileMetadata,
)
from modelpack.v1.mediatype import (
    ARTIFACT_TYPE_MODEL_MANIFEST,
    MEDIA_TYPE_MODEL_CONFIG,
    MEDIA_TYPE_MODEL_WEIGHT_RAW,
    MEDIA_TYPE_MODEL_WEIGHT,
    MEDIA_TYPE_MODEL_WEIGHT_GZIP,
    MEDIA_TYPE_MODEL_WEIGHT_ZSTD,
    MEDIA_TYPE_MODEL_WEIGHT_CONFIG_RAW,
    MEDIA_TYPE_MODEL_WEIGHT_CONFIG,
    MEDIA_TYPE_MODEL_WEIGHT_CONFIG_GZIP,
    MEDIA_TYPE_MODEL_WEIGHT_CONFIG_ZSTD,
    MEDIA_TYPE_MODEL_DOC_RAW,
    MEDIA_TYPE_MODEL_DOC,
    MEDIA_TYPE_MODEL_DOC_GZIP,
    MEDIA_TYPE_MODEL_DOC_ZSTD,
    MEDIA_TYPE_MODEL_CODE_RAW,
    MEDIA_TYPE_MODEL_CODE,
    MEDIA_TYPE_MODEL_CODE_GZIP,
    MEDIA_TYPE_MODEL_CODE_ZSTD,
    MEDIA_TYPE_MODEL_DATASET_RAW,
    MEDIA_TYPE_MODEL_DATASET,
    MEDIA_TYPE_MODEL_DATASET_GZIP,
    MEDIA_TYPE_MODEL_DATASET_ZSTD,
)
from modelpack.v1.validator import validate_config

__all__ = [
    "Model",
    "ModelCapabilities",
    "ModelConfig",
    "ModelDescriptor",
    "ModelFS",
    "Modality",
    "FileMetadata",
    "ANNOTATION_FILEPATH",
    "ANNOTATION_FILE_METADATA",
    "ANNOTATION_MEDIA_TYPE_UNTESTED",
    "ARTIFACT_TYPE_MODEL_MANIFEST",
    "MEDIA_TYPE_MODEL_CONFIG",
    "MEDIA_TYPE_MODEL_WEIGHT_RAW",
    "MEDIA_TYPE_MODEL_WEIGHT",
    "MEDIA_TYPE_MODEL_WEIGHT_GZIP",
    "MEDIA_TYPE_MODEL_WEIGHT_ZSTD",
    "MEDIA_TYPE_MODEL_WEIGHT_CONFIG_RAW",
    "MEDIA_TYPE_MODEL_WEIGHT_CONFIG",
    "MEDIA_TYPE_MODEL_WEIGHT_CONFIG_GZIP",
    "MEDIA_TYPE_MODEL_WEIGHT_CONFIG_ZSTD",
    "MEDIA_TYPE_MODEL_DOC_RAW",
    "MEDIA_TYPE_MODEL_DOC",
    "MEDIA_TYPE_MODEL_DOC_GZIP",
    "MEDIA_TYPE_MODEL_DOC_ZSTD",
    "MEDIA_TYPE_MODEL_CODE_RAW",
    "MEDIA_TYPE_MODEL_CODE",
    "MEDIA_TYPE_MODEL_CODE_GZIP",
    "MEDIA_TYPE_MODEL_CODE_ZSTD",
    "MEDIA_TYPE_MODEL_DATASET_RAW",
    "MEDIA_TYPE_MODEL_DATASET",
    "MEDIA_TYPE_MODEL_DATASET_GZIP",
    "MEDIA_TYPE_MODEL_DATASET_ZSTD",
    "validate_config",
]
