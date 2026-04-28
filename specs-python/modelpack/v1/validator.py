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

"""JSON schema validation for ModelPack configs.

Loads config-schema.json from the repo root (schema/config-schema.json)
as the single source of truth, matching the Go validator.
"""

from __future__ import annotations

import json
from pathlib import Path

from jsonschema import Draft4Validator, FormatChecker


def _load_schema() -> dict:
    """Load and return the config JSON schema from the repo root."""
    schema_path = (
        Path(__file__).resolve().parent.parent.parent.parent
        / "schema"
        / "config-schema.json"
    )
    with schema_path.open(encoding="utf-8") as f:
        return json.load(f)


def validate_config(data: dict | str) -> None:
    """Validate a model config against the JSON schema.

    Args:
        data: Either a dict or a JSON string representing the model config.

    Raises:
        jsonschema.ValidationError: If the config is invalid.
        jsonschema.SchemaError: If the schema itself is invalid.
        json.JSONDecodeError: If data is a string that is not valid JSON.
    """
    if isinstance(data, str):
        data = json.loads(data)

    schema = _load_schema()
    format_checker = FormatChecker()
    Draft4Validator(schema, format_checker=format_checker).validate(data)
