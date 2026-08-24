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

"""Annotation constants and types matching specs-go/v1/annotations.go."""

from __future__ import annotations

from dataclasses import dataclass, field
from datetime import datetime, timezone

# Annotation key for the file path of the layer.
ANNOTATION_FILEPATH = "org.cncf.model.filepath"

# Annotation key for the file metadata of the layer.
ANNOTATION_FILE_METADATA = "org.cncf.model.file.metadata+json"

# Annotation key for file media type untested flag of the layer.
ANNOTATION_MEDIA_TYPE_UNTESTED = "org.cncf.model.file.mediatype.untested"


def _format_datetime(dt: datetime) -> str:
    """Format a datetime as RFC 3339 with 'Z' suffix for UTC, matching Go."""
    s = dt.isoformat()
    if s.endswith("+00:00"):
        s = s[:-6] + "Z"
    return s


@dataclass
class FileMetadata:
    """Represents the metadata of a file.

    Mirrors the Go FileMetadata struct in specs-go/v1/annotations.go.
    """

    name: str = ""
    mode: int = 0
    uid: int = 0
    gid: int = 0
    size: int = 0
    mod_time: datetime = field(
        default_factory=lambda: datetime(1, 1, 1, tzinfo=timezone.utc)
    )
    typeflag: int = 0

    def to_dict(self) -> dict:
        """Serialize to a dict matching the JSON field names.

        All fields are always present, matching Go's FileMetadata
        which has no omitempty tags.
        """
        return {
            "name": self.name,
            "mode": self.mode,
            "uid": self.uid,
            "gid": self.gid,
            "size": self.size,
            "mtime": _format_datetime(self.mod_time),
            "typeflag": self.typeflag,
        }

    @classmethod
    def from_dict(cls, data: dict) -> FileMetadata:
        """Deserialize from a dict with JSON field names."""
        mod_time = None
        if "mtime" in data:
            mod_time = datetime.fromisoformat(data["mtime"].replace("Z", "+00:00"))
        return cls(
            name=data.get("name", ""),
            mode=data.get("mode", 0),
            uid=data.get("uid", 0),
            gid=data.get("gid", 0),
            size=data.get("size", 0),
            mod_time=mod_time,
            typeflag=data.get("typeflag", 0),
        )
