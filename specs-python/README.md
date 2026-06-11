# Python ModelPack Types

This directory provides Python data structures for the CNCF ModelPack specification.

The core model types are **auto-generated** from the canonical JSON Schema at `schema/config-schema.json` using [datamodel-code-generator](https://github.com/koxudaxi/datamodel-code-generator), ensuring they stay in sync with the specification automatically.

## Requirements

- Python >= 3.10
- pydantic >= 2
- jsonschema >= 4.20.0

## Installation

```bash
cd specs-python
pip install -e .
```

For development:

```bash
pip install -e ".[dev]"
```

## Usage

```python
from modelpack.v1 import Model, ModelDescriptor, ModelFS, ModelConfig

# Create a model from a JSON payload
model = Model.model_validate_json(json_payload)
print(model.descriptor.name)

# Validate a config dict against the JSON schema
from modelpack.v1 import validate_config
validate_config(config_dict)
```

## Regenerate Models

If the schema changes, regenerate the Pydantic models:

```bash
pip install datamodel-code-generator
make generate-python-models
```

This runs `tools/generate_python_models.py`, which regenerates `specs-python/modelpack/v1/models.py`.

**Do not edit `models.py` manually.** Update the schema and regenerate instead.

## Run Tests

```bash
cd specs-python
pytest
```
