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

"""Tests for annotation constants and FileMetadata type."""

from datetime import datetime, timezone

from modelpack.v1.annotations import (
    ANNOTATION_FILE_METADATA,
    ANNOTATION_FILEPATH,
    ANNOTATION_MEDIA_TYPE_UNTESTED,
    FileMetadata,
)


class TestAnnotationConstants:
    """Verify annotation constants match Go definitions exactly."""

    def test_filepath(self):
        assert ANNOTATION_FILEPATH == "org.cncf.model.filepath"

    def test_file_metadata(self):
        assert ANNOTATION_FILE_METADATA == "org.cncf.model.file.metadata+json"

    def test_media_type_untested(self):
        assert (
            ANNOTATION_MEDIA_TYPE_UNTESTED == "org.cncf.model.file.mediatype.untested"
        )


class TestFileMetadata:
    """Tests for FileMetadata serialization."""

    def test_round_trip(self):
        dt = datetime(2025, 1, 1, 12, 0, 0, tzinfo=timezone.utc)
        meta = FileMetadata(
            name="model.bin",
            mode=0o644,
            uid=1000,
            gid=1000,
            size=1024,
            mod_time=dt,
            typeflag=0,
        )
        d = meta.to_dict()
        assert d["name"] == "model.bin"
        assert d["mode"] == 0o644
        assert d["size"] == 1024
        assert "mtime" in d

        restored = FileMetadata.from_dict(d)
        assert restored.name == "model.bin"
        assert restored.mode == 0o644
        assert restored.size == 1024

    def test_empty(self):
        meta = FileMetadata()
        d = meta.to_dict()
        assert d["name"] == ""
        assert d["size"] == 0
        assert "mtime" in d
