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

"""Validation tests mirroring the Go test cases in schema/config_test.go.

Each test case matches the corresponding Go test to ensure
consistent validation behavior between the Go and Python SDKs.
"""

import json

import pytest
from jsonschema import ValidationError

from modelpack.v1.validator import validate_config

# A valid base config used across tests.
VALID_CONFIG = {
    "descriptor": {
        "name": "xyz-3-8B-Instruct",
        "version": "3.1",
    },
    "config": {
        "paramSize": "8b",
    },
    "modelfs": {
        "type": "layers",
        "diffIds": [
            "sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
        ],
    },
}


class TestValidConfigCases:
    """Tests that valid configs pass validation."""

    def test_minimal_valid(self):
        validate_config(VALID_CONFIG)

    def test_valid_as_json_string(self):
        validate_config(json.dumps(VALID_CONFIG))

    def test_valid_with_all_fields(self):
        config = {
            "descriptor": {
                "name": "llama3-8b-instruct",
                "version": "3.1",
                "family": "llama3",
                "authors": ["Meta"],
                "vendor": "Meta",
                "licenses": ["Apache-2.0"],
                "title": "Llama 3 8B Instruct",
                "description": "An instruction-tuned LLM",
                "createdAt": "2025-01-01T00:00:00Z",
                "docURL": "https://example.com/docs",
                "sourceURL": "https://github.com/meta/llama3",
                "datasetsURL": ["https://example.com/data"],
                "revision": "abc123",
            },
            "config": {
                "architecture": "transformer",
                "format": "safetensors",
                "paramSize": "8b",
                "precision": "float16",
                "quantization": "awq",
                "capabilities": {
                    "inputTypes": ["text"],
                    "outputTypes": ["text"],
                    "knowledgeCutoff": "2025-01-01T00:00:00Z",
                    "reasoning": True,
                    "toolUsage": True,
                    "reward": False,
                    "languages": ["en", "fr"],
                },
            },
            "modelfs": {
                "type": "layers",
                "diffIds": [
                    "sha256:abcdef1234567890abcdef1234567890"
                    "abcdef1234567890abcdef1234567890"
                ],
            },
        }
        validate_config(config)


class TestFailureConfigCases:
    """Tests mirroring Go config_test.go failure cases.

    Each test corresponds to a numbered test case in the Go file.
    """

    def test_config_missing(self):
        """Go test 0: config is missing."""
        config = {
            "descriptor": {
                "name": "xyz-3-8B-Instruct",
                "version": "3.1",
            },
            "modelfs": {
                "type": "layers",
                "diffIds": [
                    "sha256:1234567890abcdef1234567890abcdef"
                    "1234567890abcdef1234567890abcdef"
                ],
            },
        }
        with pytest.raises(ValidationError):
            validate_config(config)

    def test_version_is_number(self):
        """Go test 1: version is a number."""
        config = {
            "descriptor": {
                "name": "xyz-3-8B-Instruct",
                "version": 3.1,
            },
            "config": {"paramSize": "8b"},
            "modelfs": {
                "type": "layers",
                "diffIds": [
                    "sha256:1234567890abcdef1234567890abcdef"
                    "1234567890abcdef1234567890abcdef"
                ],
            },
        }
        with pytest.raises(ValidationError):
            validate_config(config)

    def test_revision_is_number(self):
        """Go test 2: revision is a number."""
        config = {
            "descriptor": {
                "name": "xyz-3-8B-Instruct",
                "version": "3.1",
                "revision": 1234567890,
            },
            "config": {"paramSize": "8b"},
            "modelfs": {
                "type": "layers",
                "diffIds": [
                    "sha256:1234567890abcdef1234567890abcdef"
                    "1234567890abcdef1234567890abcdef"
                ],
            },
        }
        with pytest.raises(ValidationError):
            validate_config(config)

    def test_created_at_not_rfc3339(self):
        """Go test 3: createdAt is not RFC3339 format."""
        config = {
            "descriptor": {
                "name": "xyz-3-8B-Instruct",
                "version": "3.1",
                "createdAt": "2025/01/01T00:00:00Z",
            },
            "config": {"paramSize": "8b"},
            "modelfs": {
                "type": "layers",
                "diffIds": [
                    "sha256:1234567890abcdef1234567890abcdef"
                    "1234567890abcdef1234567890abcdef"
                ],
            },
        }
        with pytest.raises(ValidationError):
            validate_config(config)

    def test_authors_not_array(self):
        """Go test 4: authors is not an array."""
        config = {
            "descriptor": {
                "name": "xyz-3-8B-Instruct",
                "version": "3.1",
                "authors": "John Doe",
            },
            "config": {"paramSize": "8b"},
            "modelfs": {
                "type": "layers",
                "diffIds": [
                    "sha256:1234567890abcdef1234567890abcdef"
                    "1234567890abcdef1234567890abcdef"
                ],
            },
        }
        with pytest.raises(ValidationError):
            validate_config(config)

    def test_licenses_not_array(self):
        """Go test 5: licenses is not an array."""
        config = {
            "descriptor": {
                "name": "xyz-3-8B-Instruct",
                "version": "3.1",
                "licenses": "Apache-2.0",
            },
            "config": {"paramSize": "8b"},
            "modelfs": {
                "type": "layers",
                "diffIds": [
                    "sha256:1234567890abcdef1234567890abcdef"
                    "1234567890abcdef1234567890abcdef"
                ],
            },
        }
        with pytest.raises(ValidationError):
            validate_config(config)

    def test_doc_url_is_array(self):
        """Go test 6: docURL is an array."""
        config = {
            "descriptor": {
                "name": "xyz-3-8B-Instruct",
                "version": "3.1",
                "docURL": ["https://example.com/doc"],
            },
            "config": {"paramSize": "8b"},
            "modelfs": {
                "type": "layers",
                "diffIds": [
                    "sha256:1234567890abcdef1234567890abcdef"
                    "1234567890abcdef1234567890abcdef"
                ],
            },
        }
        with pytest.raises(ValidationError):
            validate_config(config)

    def test_source_url_is_array(self):
        """Go test 7: sourceURL is an array."""
        config = {
            "descriptor": {
                "name": "xyz-3-8B-Instruct",
                "version": "3.1",
                "sourceURL": ["https://github.com/xyz/xyz3"],
            },
            "config": {"paramSize": "8b"},
            "modelfs": {
                "type": "layers",
                "diffIds": [
                    "sha256:1234567890abcdef1234567890abcdef"
                    "1234567890abcdef1234567890abcdef"
                ],
            },
        }
        with pytest.raises(ValidationError):
            validate_config(config)

    def test_datasets_url_not_array(self):
        """Go test 8: datasetsURL is not an array."""
        config = {
            "descriptor": {
                "name": "xyz-3-8B-Instruct",
                "version": "3.1",
                "sourceURL": "https://github.com/xyz/xyz3",
                "datasetsURL": "https://example.com/dataset",
            },
            "config": {"paramSize": "8b"},
            "modelfs": {
                "type": "layers",
                "diffIds": [
                    "sha256:1234567890abcdef1234567890abcdef"
                    "1234567890abcdef1234567890abcdef"
                ],
            },
        }
        with pytest.raises(ValidationError):
            validate_config(config)

    def test_param_size_is_number(self):
        """Go test 9: paramSize is a number."""
        config = {
            "descriptor": {
                "name": "xyz-3-8B-Instruct",
                "version": "3.1",
            },
            "config": {"paramSize": 8000000},
            "modelfs": {
                "type": "layers",
                "diffIds": [
                    "sha256:1234567890abcdef1234567890abcdef"
                    "1234567890abcdef1234567890abcdef"
                ],
            },
        }
        with pytest.raises(ValidationError):
            validate_config(config)

    def test_precision_is_number(self):
        """Go test 10: precision is a number."""
        config = {
            "descriptor": {
                "name": "xyz-3-8B-Instruct",
                "version": "3.1",
            },
            "config": {"precision": 16},
            "modelfs": {
                "type": "layers",
                "diffIds": [
                    "sha256:1234567890abcdef1234567890abcdef"
                    "1234567890abcdef1234567890abcdef"
                ],
            },
        }
        with pytest.raises(ValidationError):
            validate_config(config)

    def test_type_not_layers(self):
        """Go test 11: type is not 'layers'."""
        config = {
            "descriptor": {
                "name": "xyz-3-8B-Instruct",
                "version": "3.1",
            },
            "config": {"paramSize": "8b"},
            "modelfs": {
                "type": "layer",
                "diffIds": [
                    "sha256:1234567890abcdef1234567890abcdef"
                    "1234567890abcdef1234567890abcdef"
                ],
            },
        }
        with pytest.raises(ValidationError):
            validate_config(config)

    def test_diff_ids_not_array(self):
        """Go test 12: diffIds is not an array."""
        config = {
            "descriptor": {
                "name": "xyz-3-8B-Instruct",
                "version": "3.1",
            },
            "config": {"paramSize": "8b"},
            "modelfs": {
                "type": "layers",
                "diffIds": (
                    "sha256:1234567890abcdef1234567890abcdef"
                    "1234567890abcdef1234567890abcdef"
                ),
            },
        }
        with pytest.raises(ValidationError):
            validate_config(config)

    def test_diff_ids_empty(self):
        """Go test 13: diffIds is empty."""
        config = {
            "descriptor": {
                "name": "xyz-3-8B-Instruct",
                "version": "3.1",
            },
            "config": {"paramSize": "8b"},
            "modelfs": {
                "type": "layers",
                "diffIds": [],
            },
        }
        with pytest.raises(ValidationError):
            validate_config(config)

    def test_input_types_not_array(self):
        """Go test 14: inputTypes is not an array."""
        config = {
            "descriptor": {
                "name": "xyz-3-8B-Instruct",
                "version": "3.1",
            },
            "config": {
                "paramSize": "8b",
                "capabilities": {"inputTypes": "text"},
            },
            "modelfs": {
                "type": "layers",
                "diffIds": [
                    "sha256:1234567890abcdef1234567890abcdef"
                    "1234567890abcdef1234567890abcdef"
                ],
            },
        }
        with pytest.raises(ValidationError):
            validate_config(config)

    def test_output_types_not_array(self):
        """Go test 15: outputTypes is not an array."""
        config = {
            "descriptor": {
                "name": "xyz-3-8B-Instruct",
                "version": "3.1",
            },
            "config": {
                "paramSize": "8b",
                "capabilities": {"outputTypes": "text"},
            },
            "modelfs": {
                "type": "layers",
                "diffIds": [
                    "sha256:1234567890abcdef1234567890abcdef"
                    "1234567890abcdef1234567890abcdef"
                ],
            },
        }
        with pytest.raises(ValidationError):
            validate_config(config)

    def test_invalid_modality(self):
        """Go test 16: invalid modality value."""
        config = {
            "descriptor": {
                "name": "xyz-3-8B-Instruct",
                "version": "3.1",
            },
            "config": {
                "paramSize": "8b",
                "capabilities": {"inputTypes": ["img"]},
            },
            "modelfs": {
                "type": "layers",
                "diffIds": [
                    "sha256:1234567890abcdef1234567890abcdef"
                    "1234567890abcdef1234567890abcdef"
                ],
            },
        }
        with pytest.raises(ValidationError):
            validate_config(config)

    def test_knowledge_cutoff_not_rfc3339(self):
        """Go test 17: knowledgeCutoff is not RFC3339 format."""
        config = {
            "descriptor": {
                "name": "xyz-3-8B-Instruct",
                "version": "3.1",
            },
            "config": {
                "paramSize": "8b",
                "capabilities": {
                    "inputTypes": ["text"],
                    "outputTypes": ["text"],
                    "knowledgeCutoff": "2025-01-01",
                },
            },
            "modelfs": {
                "type": "layers",
                "diffIds": [
                    "sha256:1234567890abcdef1234567890abcdef"
                    "1234567890abcdef1234567890abcdef"
                ],
            },
        }
        with pytest.raises(ValidationError):
            validate_config(config)

    def test_reasoning_not_boolean(self):
        """Go test 18: reasoning is not boolean."""
        config = {
            "descriptor": {
                "name": "xyz-3-8B-Instruct",
                "version": "3.1",
            },
            "config": {
                "paramSize": "8b",
                "capabilities": {
                    "inputTypes": ["text"],
                    "outputTypes": ["text"],
                    "reasoning": "true",
                },
            },
            "modelfs": {
                "type": "layers",
                "diffIds": [
                    "sha256:1234567890abcdef1234567890abcdef"
                    "1234567890abcdef1234567890abcdef"
                ],
            },
        }
        with pytest.raises(ValidationError):
            validate_config(config)

    def test_tool_usage_not_boolean(self):
        """Go test 19: toolUsage is not boolean."""
        config = {
            "descriptor": {
                "name": "xyz-3-8B-Instruct",
                "version": "3.1",
            },
            "config": {
                "paramSize": "8b",
                "capabilities": {
                    "inputTypes": ["text"],
                    "outputTypes": ["text"],
                    "toolUsage": "true",
                },
            },
            "modelfs": {
                "type": "layers",
                "diffIds": [
                    "sha256:1234567890abcdef1234567890abcdef"
                    "1234567890abcdef1234567890abcdef"
                ],
            },
        }
        with pytest.raises(ValidationError):
            validate_config(config)

    def test_reward_not_boolean(self):
        """Go test 20: reward is not boolean."""
        config = {
            "descriptor": {
                "name": "xyz-3-8B-Instruct",
                "version": "3.1",
            },
            "config": {
                "paramSize": "8b",
                "capabilities": {
                    "inputTypes": ["text"],
                    "outputTypes": ["text"],
                    "reward": "true",
                },
            },
            "modelfs": {
                "type": "layers",
                "diffIds": [
                    "sha256:1234567890abcdef1234567890abcdef"
                    "1234567890abcdef1234567890abcdef"
                ],
            },
        }
        with pytest.raises(ValidationError):
            validate_config(config)

    def test_languages_not_array(self):
        """Go test 21: languages is not an array."""
        config = {
            "descriptor": {
                "name": "xyz-3-8B-Instruct",
                "version": "3.1",
            },
            "config": {
                "paramSize": "8b",
                "capabilities": {
                    "inputTypes": ["text"],
                    "outputTypes": ["text"],
                    "languages": "en",
                },
            },
            "modelfs": {
                "type": "layers",
                "diffIds": [
                    "sha256:1234567890abcdef1234567890abcdef"
                    "1234567890abcdef1234567890abcdef"
                ],
            },
        }
        with pytest.raises(ValidationError):
            validate_config(config)

    def test_language_code_not_iso639(self):
        """Go test 22: language code is not a two-letter ISO 639 code."""
        config = {
            "descriptor": {
                "name": "xyz-3-8B-Instruct",
                "version": "3.1",
            },
            "config": {
                "paramSize": "8b",
                "capabilities": {
                    "inputTypes": ["text"],
                    "outputTypes": ["text"],
                    "languages": ["fra"],
                },
            },
            "modelfs": {
                "type": "layers",
                "diffIds": [
                    "sha256:1234567890abcdef1234567890abcdef"
                    "1234567890abcdef1234567890abcdef"
                ],
            },
        }
        with pytest.raises(ValidationError):
            validate_config(config)

    def test_unknown_field_in_capabilities(self):
        """Go test 23: unknown field in capabilities."""
        config = {
            "descriptor": {
                "name": "xyz-3-8B-Instruct",
                "version": "3.1",
            },
            "config": {
                "paramSize": "8b",
                "capabilities": {
                    "inputTypes": ["text"],
                    "unknownField": True,
                },
            },
            "modelfs": {
                "type": "layers",
                "diffIds": [
                    "sha256:1234567890abcdef1234567890abcdef"
                    "1234567890abcdef1234567890abcdef"
                ],
            },
        }
        with pytest.raises(ValidationError):
            validate_config(config)


class TestEdgeCases:
    """Additional edge case tests."""

    def test_empty_dict(self):
        with pytest.raises(ValidationError):
            validate_config({})

    def test_invalid_json_string(self):
        with pytest.raises(Exception):
            validate_config("{invalid json")

    def test_empty_name(self):
        """Name with minLength: 1 should reject empty string."""
        config = {
            "descriptor": {"name": "", "version": "1.0"},
            "config": {"paramSize": "8b"},
            "modelfs": {
                "type": "layers",
                "diffIds": ["sha256:abc"],
            },
        }
        with pytest.raises(ValidationError):
            validate_config(config)

    def test_unknown_field_at_root(self):
        config = {
            "descriptor": {"name": "test", "version": "1.0"},
            "config": {"paramSize": "8b"},
            "modelfs": {
                "type": "layers",
                "diffIds": ["sha256:abc"],
            },
            "extraField": "should fail",
        }
        with pytest.raises(ValidationError):
            validate_config(config)

    def test_unknown_field_in_descriptor(self):
        config = {
            "descriptor": {
                "name": "test",
                "version": "1.0",
                "unknownField": "value",
            },
            "config": {"paramSize": "8b"},
            "modelfs": {
                "type": "layers",
                "diffIds": ["sha256:abc"],
            },
        }
        with pytest.raises(ValidationError):
            validate_config(config)

    def test_unknown_field_in_config(self):
        config = {
            "descriptor": {"name": "test", "version": "1.0"},
            "config": {"paramSize": "8b", "unknownField": "value"},
            "modelfs": {
                "type": "layers",
                "diffIds": ["sha256:abc"],
            },
        }
        with pytest.raises(ValidationError):
            validate_config(config)

    def test_modelfs_missing(self):
        config = {
            "descriptor": {"name": "test", "version": "1.0"},
            "config": {"paramSize": "8b"},
        }
        with pytest.raises(ValidationError):
            validate_config(config)

    def test_descriptor_missing(self):
        config = {
            "config": {"paramSize": "8b"},
            "modelfs": {
                "type": "layers",
                "diffIds": ["sha256:abc"],
            },
        }
        with pytest.raises(ValidationError):
            validate_config(config)
