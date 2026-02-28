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

"""Tests for model config types - serialization and deserialization."""

import json
from datetime import datetime, timezone

import pytest

from modelpack.v1.config import (
    Model,
    ModelCapabilities,
    ModelConfig,
    ModelDescriptor,
    ModelFS,
    Modality,
)


class TestModality:
    """Tests for the Modality enum."""

    def test_all_values(self):
        assert Modality.TEXT.value == "text"
        assert Modality.IMAGE.value == "image"
        assert Modality.AUDIO.value == "audio"
        assert Modality.VIDEO.value == "video"
        assert Modality.EMBEDDING.value == "embedding"
        assert Modality.OTHER.value == "other"

    def test_from_string(self):
        assert Modality("text") == Modality.TEXT
        assert Modality("image") == Modality.IMAGE

    def test_invalid_value(self):
        with pytest.raises(ValueError):
            Modality("invalid")


class TestModelCapabilities:
    """Tests for ModelCapabilities serialization."""

    def test_empty(self):
        caps = ModelCapabilities()
        assert caps.to_dict() == {}

    def test_round_trip(self):
        caps = ModelCapabilities(
            input_types=[Modality.TEXT, Modality.IMAGE],
            output_types=[Modality.TEXT],
            reasoning=True,
            tool_usage=False,
            reward=True,
            languages=["en", "fr"],
        )
        d = caps.to_dict()
        assert d["inputTypes"] == ["text", "image"]
        assert d["outputTypes"] == ["text"]
        assert d["reasoning"] is True
        assert d["toolUsage"] is False
        assert d["reward"] is True
        assert d["languages"] == ["en", "fr"]

        restored = ModelCapabilities.from_dict(d)
        assert restored.input_types == [Modality.TEXT, Modality.IMAGE]
        assert restored.reasoning is True
        assert restored.tool_usage is False

    def test_knowledge_cutoff(self):
        dt = datetime(2025, 1, 1, tzinfo=timezone.utc)
        caps = ModelCapabilities(knowledge_cutoff=dt)
        d = caps.to_dict()
        assert "knowledgeCutoff" in d

        restored = ModelCapabilities.from_dict(d)
        assert restored.knowledge_cutoff is not None


class TestModelConfig:
    """Tests for ModelConfig serialization."""

    def test_empty(self):
        cfg = ModelConfig()
        assert cfg.to_dict() == {}

    def test_round_trip(self):
        cfg = ModelConfig(
            architecture="transformer",
            format="safetensors",
            param_size="8b",
            precision="float16",
            quantization="awq",
        )
        d = cfg.to_dict()
        assert d["architecture"] == "transformer"
        assert d["paramSize"] == "8b"

        restored = ModelConfig.from_dict(d)
        assert restored.architecture == "transformer"
        assert restored.param_size == "8b"

    def test_with_capabilities(self):
        cfg = ModelConfig(
            param_size="8b",
            capabilities=ModelCapabilities(
                input_types=[Modality.TEXT],
                output_types=[Modality.TEXT],
            ),
        )
        d = cfg.to_dict()
        assert "capabilities" in d
        assert d["capabilities"]["inputTypes"] == ["text"]


class TestModelFS:
    """Tests for ModelFS serialization."""

    def test_round_trip(self):
        fs = ModelFS(
            type="layers",
            diff_ids=["sha256:abc123"],
        )
        d = fs.to_dict()
        assert d["type"] == "layers"
        assert d["diffIds"] == ["sha256:abc123"]

        restored = ModelFS.from_dict(d)
        assert restored.type == "layers"
        assert restored.diff_ids == ["sha256:abc123"]


class TestModelDescriptor:
    """Tests for ModelDescriptor serialization."""

    def test_empty(self):
        desc = ModelDescriptor()
        assert desc.to_dict() == {}

    def test_round_trip(self):
        desc = ModelDescriptor(
            name="llama3-8b-instruct",
            version="3.1",
            family="llama3",
            authors=["Meta"],
            licenses=["Apache-2.0"],
        )
        d = desc.to_dict()
        assert d["name"] == "llama3-8b-instruct"
        assert d["version"] == "3.1"

        restored = ModelDescriptor.from_dict(d)
        assert restored.name == "llama3-8b-instruct"
        assert restored.authors == ["Meta"]

    def test_created_at(self):
        dt = datetime(2025, 6, 15, 10, 30, 0, tzinfo=timezone.utc)
        desc = ModelDescriptor(name="test", created_at=dt)
        d = desc.to_dict()
        assert "createdAt" in d

        restored = ModelDescriptor.from_dict(d)
        assert restored.created_at is not None


class TestModel:
    """Tests for Model serialization."""

    def test_minimal(self):
        model = Model(
            descriptor=ModelDescriptor(name="test-model"),
            modelfs=ModelFS(type="layers", diff_ids=["sha256:abc"]),
            config=ModelConfig(param_size="8b"),
        )
        d = model.to_dict()
        assert d["descriptor"]["name"] == "test-model"
        assert d["modelfs"]["type"] == "layers"
        assert d["config"]["paramSize"] == "8b"

    def test_json_round_trip(self):
        model = Model(
            descriptor=ModelDescriptor(
                name="llama3-8b-instruct",
                version="3.1",
            ),
            modelfs=ModelFS(
                type="layers",
                diff_ids=[
                    "sha256:1234567890abcdef1234567890abcdef"
                    "1234567890abcdef1234567890abcdef"
                ],
            ),
            config=ModelConfig(param_size="8b"),
        )
        json_str = model.to_json()
        restored = Model.from_json(json_str)
        assert restored.descriptor.name == "llama3-8b-instruct"
        assert restored.modelfs.type == "layers"
        assert restored.config.param_size == "8b"

    def test_from_json_string(self):
        data = json.dumps(
            {
                "descriptor": {"name": "test"},
                "modelfs": {"type": "layers", "diffIds": ["sha256:abc"]},
                "config": {"paramSize": "1b"},
            }
        )
        model = Model.from_json(data)
        assert model.descriptor.name == "test"
        assert model.config.param_size == "1b"

    def test_full_model(self):
        model = Model(
            descriptor=ModelDescriptor(
                name="qwen2-vl-72b-instruct",
                version="2.0",
                family="qwen2",
                vendor="Alibaba",
                authors=["Qwen Team"],
                licenses=["Apache-2.0"],
                title="Qwen2 VL 72B Instruct",
                description="A vision-language model",
                doc_url="https://example.com/docs",
                source_url="https://github.com/example/qwen2",
                datasets_url=["https://example.com/dataset"],
            ),
            modelfs=ModelFS(
                type="layers",
                diff_ids=["sha256:aabbcc", "sha256:ddeeff"],
            ),
            config=ModelConfig(
                architecture="transformer",
                format="safetensors",
                param_size="72b",
                precision="bfloat16",
                capabilities=ModelCapabilities(
                    input_types=[Modality.TEXT, Modality.IMAGE],
                    output_types=[Modality.TEXT],
                    reasoning=True,
                    tool_usage=True,
                    languages=["en", "zh"],
                ),
            ),
        )
        d = model.to_dict()
        assert d["descriptor"]["vendor"] == "Alibaba"
        assert d["config"]["capabilities"]["inputTypes"] == ["text", "image"]
        assert d["config"]["capabilities"]["languages"] == ["en", "zh"]

        json_str = model.to_json()
        restored = Model.from_json(json_str)
        assert restored.config.capabilities.input_types == [
            Modality.TEXT,
            Modality.IMAGE,
        ]
        assert restored.config.capabilities.languages == ["en", "zh"]
