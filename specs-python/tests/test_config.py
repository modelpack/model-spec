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

"""Tests for auto-generated Pydantic model types."""

import json
from datetime import datetime, timezone

import pytest
from pydantic import ValidationError

from modelpack.v1.models import (
    Language,
    Modality,
    Model,
    ModelCapabilities,
    ModelConfig,
    ModelDescriptor,
    ModelFS,
)


class TestModality:
    """Tests for the Modality RootModel."""

    def test_all_values(self):
        for val in ("text", "image", "audio", "video", "embedding", "other"):
            m = Modality(root=val)
            assert m.root == val

    def test_from_string(self):
        m = Modality.model_validate("text")
        assert m.root == "text"

    def test_invalid_value(self):
        with pytest.raises(ValidationError):
            Modality.model_validate("invalid")


class TestLanguage:
    """Tests for the Language RootModel."""

    def test_valid(self):
        lang = Language.model_validate("en")
        assert lang.root == "en"

    def test_invalid_three_letter(self):
        with pytest.raises(ValidationError):
            Language.model_validate("fra")

    def test_invalid_uppercase(self):
        with pytest.raises(ValidationError):
            Language.model_validate("EN")


class TestModelCapabilities:
    """Tests for ModelCapabilities Pydantic model."""

    def test_empty(self):
        caps = ModelCapabilities()
        d = caps.model_dump(exclude_none=True)
        assert d == {}

    def test_round_trip(self):
        caps = ModelCapabilities(
            inputTypes=[Modality(root="text"), Modality(root="image")],
            outputTypes=[Modality(root="text")],
            reasoning=True,
            toolUsage=False,
            reward=True,
            languages=[Language(root="en"), Language(root="fr")],
        )
        d = caps.model_dump(exclude_none=True)
        assert d["inputTypes"] == ["text", "image"]
        assert d["outputTypes"] == ["text"]
        assert d["reasoning"] is True
        assert d["toolUsage"] is False
        assert d["reward"] is True
        assert d["languages"] == ["en", "fr"]

        restored = ModelCapabilities.model_validate(d)
        assert restored.inputTypes[0].root == "text"
        assert restored.reasoning is True
        assert restored.toolUsage is False

    def test_knowledge_cutoff(self):
        dt = datetime(2025, 1, 1, tzinfo=timezone.utc)
        caps = ModelCapabilities(knowledgeCutoff=dt)
        d = caps.model_dump(exclude_none=True, mode="json")
        assert "knowledgeCutoff" in d

        restored = ModelCapabilities.model_validate(d)
        assert restored.knowledgeCutoff is not None


class TestModelConfig:
    """Tests for ModelConfig Pydantic model."""

    def test_empty(self):
        cfg = ModelConfig()
        d = cfg.model_dump(exclude_none=True)
        assert d == {}

    def test_round_trip(self):
        cfg = ModelConfig(
            architecture="transformer",
            format="safetensors",
            paramSize="8b",
            precision="float16",
            quantization="awq",
        )
        d = cfg.model_dump(exclude_none=True)
        assert d["architecture"] == "transformer"
        assert d["paramSize"] == "8b"

        restored = ModelConfig.model_validate(d)
        assert restored.architecture == "transformer"
        assert restored.paramSize == "8b"

    def test_with_capabilities(self):
        cfg = ModelConfig(
            paramSize="8b",
            capabilities=ModelCapabilities(
                inputTypes=[Modality(root="text")],
                outputTypes=[Modality(root="text")],
            ),
        )
        d = cfg.model_dump(exclude_none=True)
        assert "capabilities" in d
        assert d["capabilities"]["inputTypes"] == ["text"]


class TestModelFS:
    """Tests for ModelFS Pydantic model."""

    def test_round_trip(self):
        fs = ModelFS(
            type="layers",
            diffIds=["sha256:abc123"],
        )
        d = fs.model_dump()
        assert d["type"] == "layers"
        assert d["diffIds"] == ["sha256:abc123"]

        restored = ModelFS.model_validate(d)
        assert restored.type == "layers"
        assert restored.diffIds == ["sha256:abc123"]

    def test_invalid_type(self):
        with pytest.raises(ValidationError):
            ModelFS(type="invalid", diffIds=["sha256:abc"])

    def test_empty_diff_ids(self):
        with pytest.raises(ValidationError):
            ModelFS(type="layers", diffIds=[])


class TestModelDescriptor:
    """Tests for ModelDescriptor Pydantic model."""

    def test_empty(self):
        desc = ModelDescriptor()
        d = desc.model_dump(exclude_none=True)
        assert d == {}

    def test_round_trip(self):
        desc = ModelDescriptor(
            name="llama3-8b-instruct",
            version="3.1",
            family="llama3",
            authors=["Meta"],
            licenses=["Apache-2.0"],
        )
        d = desc.model_dump(exclude_none=True)
        assert d["name"] == "llama3-8b-instruct"
        assert d["version"] == "3.1"

        restored = ModelDescriptor.model_validate(d)
        assert restored.name == "llama3-8b-instruct"
        assert restored.authors == ["Meta"]

    def test_created_at(self):
        dt = datetime(2025, 6, 15, 10, 30, 0, tzinfo=timezone.utc)
        desc = ModelDescriptor(name="test", createdAt=dt)
        d = desc.model_dump(exclude_none=True, mode="json")
        assert "createdAt" in d

        restored = ModelDescriptor.model_validate(d)
        assert restored.createdAt is not None

    def test_empty_name_rejected(self):
        with pytest.raises(ValidationError):
            ModelDescriptor(name="")

    def test_extra_field_rejected(self):
        with pytest.raises(ValidationError):
            ModelDescriptor.model_validate({"name": "test", "unknownField": "value"})


class TestModel:
    """Tests for Model Pydantic model."""

    def test_minimal(self):
        model = Model(
            descriptor=ModelDescriptor(name="test-model"),
            modelfs=ModelFS(type="layers", diffIds=["sha256:abc"]),
            config=ModelConfig(paramSize="8b"),
        )
        d = model.model_dump(exclude_none=True)
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
                diffIds=[
                    "sha256:1234567890abcdef1234567890abcdef"
                    "1234567890abcdef1234567890abcdef"
                ],
            ),
            config=ModelConfig(paramSize="8b"),
        )
        json_str = model.model_dump_json()
        restored = Model.model_validate_json(json_str)
        assert restored.descriptor.name == "llama3-8b-instruct"
        assert restored.modelfs.type == "layers"
        assert restored.config.paramSize == "8b"

    def test_from_json_string(self):
        data = json.dumps(
            {
                "descriptor": {"name": "test"},
                "modelfs": {"type": "layers", "diffIds": ["sha256:abc"]},
                "config": {"paramSize": "1b"},
            }
        )
        model = Model.model_validate_json(data)
        assert model.descriptor.name == "test"
        assert model.config.paramSize == "1b"

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
                docURL="https://example.com/docs",
                sourceURL="https://github.com/example/qwen2",
                datasetsURL=["https://example.com/dataset"],
            ),
            modelfs=ModelFS(
                type="layers",
                diffIds=["sha256:aabbcc", "sha256:ddeeff"],
            ),
            config=ModelConfig(
                architecture="transformer",
                format="safetensors",
                paramSize="72b",
                precision="bfloat16",
                capabilities=ModelCapabilities(
                    inputTypes=[Modality(root="text"), Modality(root="image")],
                    outputTypes=[Modality(root="text")],
                    reasoning=True,
                    toolUsage=True,
                    languages=[Language(root="en"), Language(root="zh")],
                ),
            ),
        )
        d = model.model_dump(exclude_none=True)
        assert d["descriptor"]["vendor"] == "Alibaba"
        assert d["config"]["capabilities"]["inputTypes"] == ["text", "image"]
        assert d["config"]["capabilities"]["languages"] == ["en", "zh"]

        json_str = model.model_dump_json()
        restored = Model.model_validate_json(json_str)
        assert restored.config.capabilities.inputTypes[0].root == "text"
        assert restored.config.capabilities.inputTypes[1].root == "image"
        assert restored.config.capabilities.languages[0].root == "en"
        assert restored.config.capabilities.languages[1].root == "zh"

    def test_missing_required_fields(self):
        with pytest.raises(ValidationError):
            Model.model_validate({})

    def test_extra_field_at_root(self):
        with pytest.raises(ValidationError):
            Model.model_validate(
                {
                    "descriptor": {"name": "test"},
                    "modelfs": {"type": "layers", "diffIds": ["sha256:abc"]},
                    "config": {"paramSize": "8b"},
                    "extraField": "should fail",
                }
            )
