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

"""Model configuration types matching the Go structs in specs-go/v1/config.go."""

from __future__ import annotations

import json
from dataclasses import dataclass, field
from datetime import datetime
from enum import Enum
from typing import Optional


def _format_datetime(dt: datetime) -> str:
    """Format a datetime as RFC 3339 with 'Z' suffix for UTC, matching Go."""
    s = dt.isoformat()
    if s.endswith("+00:00"):
        s = s[:-6] + "Z"
    return s


class Modality(str, Enum):
    """Defines the input and output types of the model.

    Mirrors the Go Modality type in specs-go/v1/config.go.
    """

    TEXT = "text"
    IMAGE = "image"
    AUDIO = "audio"
    VIDEO = "video"
    EMBEDDING = "embedding"
    OTHER = "other"


@dataclass
class ModelCapabilities:
    """Defines the special capabilities that the model supports.

    Mirrors the Go ModelCapabilities struct in specs-go/v1/config.go.
    """

    input_types: Optional[list[Modality]] = None
    output_types: Optional[list[Modality]] = None
    knowledge_cutoff: Optional[datetime] = None
    reasoning: Optional[bool] = None
    tool_usage: Optional[bool] = None
    reward: Optional[bool] = None
    languages: Optional[list[str]] = None

    def to_dict(self) -> dict:
        """Serialize to a dict matching the JSON schema field names."""
        d: dict = {}
        if self.input_types is not None:
            d["inputTypes"] = [m.value for m in self.input_types]
        if self.output_types is not None:
            d["outputTypes"] = [m.value for m in self.output_types]
        if self.knowledge_cutoff is not None:
            d["knowledgeCutoff"] = _format_datetime(self.knowledge_cutoff)
        if self.reasoning is not None:
            d["reasoning"] = self.reasoning
        if self.tool_usage is not None:
            d["toolUsage"] = self.tool_usage
        if self.reward is not None:
            d["reward"] = self.reward
        if self.languages is not None:
            d["languages"] = self.languages
        return d

    @classmethod
    def from_dict(cls, data: dict) -> ModelCapabilities:
        """Deserialize from a dict with JSON schema field names."""
        kwargs: dict = {}
        if "inputTypes" in data:
            kwargs["input_types"] = [Modality(v) for v in data["inputTypes"]]
        if "outputTypes" in data:
            kwargs["output_types"] = [Modality(v) for v in data["outputTypes"]]
        if "knowledgeCutoff" in data and data["knowledgeCutoff"]:
            kwargs["knowledge_cutoff"] = datetime.fromisoformat(
                data["knowledgeCutoff"].replace("Z", "+00:00")
            )
        if "reasoning" in data:
            kwargs["reasoning"] = data["reasoning"]
        if "toolUsage" in data:
            kwargs["tool_usage"] = data["toolUsage"]
        if "reward" in data:
            kwargs["reward"] = data["reward"]
        if "languages" in data:
            kwargs["languages"] = data["languages"]
        return cls(**kwargs)


@dataclass
class ModelConfig:
    """Defines the execution parameters for running a model.

    Mirrors the Go ModelConfig struct in specs-go/v1/config.go.
    """

    architecture: str = ""
    format: str = ""
    param_size: str = ""
    precision: str = ""
    quantization: str = ""
    capabilities: Optional[ModelCapabilities] = None

    def to_dict(self) -> dict:
        """Serialize to a dict matching the JSON schema field names."""
        d: dict = {}
        if self.architecture:
            d["architecture"] = self.architecture
        if self.format:
            d["format"] = self.format
        if self.param_size:
            d["paramSize"] = self.param_size
        if self.precision:
            d["precision"] = self.precision
        if self.quantization:
            d["quantization"] = self.quantization
        if self.capabilities is not None:
            d["capabilities"] = self.capabilities.to_dict()
        return d

    @classmethod
    def from_dict(cls, data: dict) -> ModelConfig:
        """Deserialize from a dict with JSON schema field names."""
        caps = None
        if "capabilities" in data:
            caps = ModelCapabilities.from_dict(data["capabilities"])
        return cls(
            architecture=data.get("architecture", ""),
            format=data.get("format", ""),
            param_size=data.get("paramSize", ""),
            precision=data.get("precision", ""),
            quantization=data.get("quantization", ""),
            capabilities=caps,
        )


@dataclass
class ModelFS:
    """Describes layer content addresses.

    Mirrors the Go ModelFS struct in specs-go/v1/config.go.
    """

    type: str = ""
    diff_ids: list[str] = field(default_factory=list)

    def to_dict(self) -> dict:
        """Serialize to a dict matching the JSON schema field names."""
        return {
            "type": self.type,
            "diffIds": self.diff_ids,
        }

    @classmethod
    def from_dict(cls, data: dict) -> ModelFS:
        """Deserialize from a dict with JSON schema field names."""
        return cls(
            type=data.get("type", ""),
            diff_ids=data.get("diffIds", []),
        )


@dataclass
class ModelDescriptor:
    """Defines the general information of a model.

    Mirrors the Go ModelDescriptor struct in specs-go/v1/config.go.
    """

    created_at: Optional[datetime] = None
    authors: Optional[list[str]] = None
    family: str = ""
    name: str = ""
    doc_url: str = ""
    source_url: str = ""
    datasets_url: Optional[list[str]] = None
    version: str = ""
    revision: str = ""
    vendor: str = ""
    licenses: Optional[list[str]] = None
    title: str = ""
    description: str = ""

    def to_dict(self) -> dict:
        """Serialize to a dict matching the JSON schema field names."""
        d: dict = {}
        if self.created_at is not None:
            d["createdAt"] = _format_datetime(self.created_at)
        if self.authors is not None:
            d["authors"] = self.authors
        if self.family:
            d["family"] = self.family
        if self.name:
            d["name"] = self.name
        if self.doc_url:
            d["docURL"] = self.doc_url
        if self.source_url:
            d["sourceURL"] = self.source_url
        if self.datasets_url is not None:
            d["datasetsURL"] = self.datasets_url
        if self.version:
            d["version"] = self.version
        if self.revision:
            d["revision"] = self.revision
        if self.vendor:
            d["vendor"] = self.vendor
        if self.licenses is not None:
            d["licenses"] = self.licenses
        if self.title:
            d["title"] = self.title
        if self.description:
            d["description"] = self.description
        return d

    @classmethod
    def from_dict(cls, data: dict) -> ModelDescriptor:
        """Deserialize from a dict with JSON schema field names."""
        created_at = None
        if "createdAt" in data:
            created_at = datetime.fromisoformat(
                data["createdAt"].replace("Z", "+00:00")
            )
        return cls(
            created_at=created_at,
            authors=data.get("authors"),
            family=data.get("family", ""),
            name=data.get("name", ""),
            doc_url=data.get("docURL", ""),
            source_url=data.get("sourceURL", ""),
            datasets_url=data.get("datasetsURL"),
            version=data.get("version", ""),
            revision=data.get("revision", ""),
            vendor=data.get("vendor", ""),
            licenses=data.get("licenses"),
            title=data.get("title", ""),
            description=data.get("description", ""),
        )


@dataclass
class Model:
    """Defines the basic information of a model.

    Provides the application/vnd.cncf.model.config.v1+json mediatype
    when marshalled to JSON.

    Mirrors the Go Model struct in specs-go/v1/config.go.
    """

    descriptor: ModelDescriptor = field(default_factory=ModelDescriptor)
    modelfs: ModelFS = field(default_factory=ModelFS)
    config: ModelConfig = field(default_factory=ModelConfig)

    def to_dict(self) -> dict:
        """Serialize to a dict matching the JSON schema field names."""
        return {
            "descriptor": self.descriptor.to_dict(),
            "modelfs": self.modelfs.to_dict(),
            "config": self.config.to_dict(),
        }

    def to_json(self, indent: Optional[int] = 2) -> str:
        """Serialize to a JSON string."""
        return json.dumps(self.to_dict(), indent=indent)

    @classmethod
    def from_dict(cls, data: dict) -> Model:
        """Deserialize from a dict with JSON schema field names."""
        return cls(
            descriptor=ModelDescriptor.from_dict(data.get("descriptor", {})),
            modelfs=ModelFS.from_dict(data.get("modelfs", {})),
            config=ModelConfig.from_dict(data.get("config", {})),
        )

    @classmethod
    def from_json(cls, json_str: str) -> Model:
        """Deserialize from a JSON string."""
        return cls.from_dict(json.loads(json_str))
