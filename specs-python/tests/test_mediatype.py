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

"""Tests for media type constants matching specs-go/v1/mediatype.go."""

from modelpack.v1 import mediatype


class TestMediaTypes:
    """Verify media type constants match Go definitions exactly."""

    def test_artifact_type(self):
        assert (
            mediatype.ARTIFACT_TYPE_MODEL_MANIFEST
            == "application/vnd.cncf.model.manifest.v1+json"
        )

    def test_config(self):
        assert (
            mediatype.MEDIA_TYPE_MODEL_CONFIG
            == "application/vnd.cncf.model.config.v1+json"
        )

    def test_weight_types(self):
        assert (
            mediatype.MEDIA_TYPE_MODEL_WEIGHT_RAW
            == "application/vnd.cncf.model.weight.v1.raw"
        )
        assert (
            mediatype.MEDIA_TYPE_MODEL_WEIGHT
            == "application/vnd.cncf.model.weight.v1.tar"
        )
        assert (
            mediatype.MEDIA_TYPE_MODEL_WEIGHT_GZIP
            == "application/vnd.cncf.model.weight.v1.tar+gzip"
        )
        assert (
            mediatype.MEDIA_TYPE_MODEL_WEIGHT_ZSTD
            == "application/vnd.cncf.model.weight.v1.tar+zstd"
        )

    def test_weight_config_types(self):
        assert (
            mediatype.MEDIA_TYPE_MODEL_WEIGHT_CONFIG_RAW
            == "application/vnd.cncf.model.weight.config.v1.raw"
        )
        assert (
            mediatype.MEDIA_TYPE_MODEL_WEIGHT_CONFIG
            == "application/vnd.cncf.model.weight.config.v1.tar"
        )
        assert (
            mediatype.MEDIA_TYPE_MODEL_WEIGHT_CONFIG_GZIP
            == "application/vnd.cncf.model.weight.config.v1.tar+gzip"
        )
        assert (
            mediatype.MEDIA_TYPE_MODEL_WEIGHT_CONFIG_ZSTD
            == "application/vnd.cncf.model.weight.config.v1.tar+zstd"
        )

    def test_doc_types(self):
        assert (
            mediatype.MEDIA_TYPE_MODEL_DOC_RAW
            == "application/vnd.cncf.model.doc.v1.raw"
        )
        assert mediatype.MEDIA_TYPE_MODEL_DOC == "application/vnd.cncf.model.doc.v1.tar"
        assert (
            mediatype.MEDIA_TYPE_MODEL_DOC_GZIP
            == "application/vnd.cncf.model.doc.v1.tar+gzip"
        )
        assert (
            mediatype.MEDIA_TYPE_MODEL_DOC_ZSTD
            == "application/vnd.cncf.model.doc.v1.tar+zstd"
        )

    def test_code_types(self):
        assert (
            mediatype.MEDIA_TYPE_MODEL_CODE_RAW
            == "application/vnd.cncf.model.code.v1.raw"
        )
        assert (
            mediatype.MEDIA_TYPE_MODEL_CODE == "application/vnd.cncf.model.code.v1.tar"
        )
        assert (
            mediatype.MEDIA_TYPE_MODEL_CODE_GZIP
            == "application/vnd.cncf.model.code.v1.tar+gzip"
        )
        assert (
            mediatype.MEDIA_TYPE_MODEL_CODE_ZSTD
            == "application/vnd.cncf.model.code.v1.tar+zstd"
        )

    def test_dataset_types(self):
        assert (
            mediatype.MEDIA_TYPE_MODEL_DATASET_RAW
            == "application/vnd.cncf.model.dataset.v1.raw"
        )
        assert (
            mediatype.MEDIA_TYPE_MODEL_DATASET
            == "application/vnd.cncf.model.dataset.v1.tar"
        )
        assert (
            mediatype.MEDIA_TYPE_MODEL_DATASET_GZIP
            == "application/vnd.cncf.model.dataset.v1.tar+gzip"
        )
        assert (
            mediatype.MEDIA_TYPE_MODEL_DATASET_ZSTD
            == "application/vnd.cncf.model.dataset.v1.tar+zstd"
        )
