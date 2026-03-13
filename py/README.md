# Python ModelPack Types

This directory provides auto-generated Python data structures for the ModelPack specification.

The models are generated from the canonical JSON Schema at `schema/config-schema.json` and are intended for downstream projects that need importable spec-aligned types.

## Usage

```python
from model_spec.v1 import Model

model = Model.model_validate_json(json_payload)
print(model.descriptor.name)
```

## Regenerate

Run:

```bash
make generate-python-api
```

This executes `tools/generate_python_models.py`, which uses `datamodel-codegen` to regenerate `py/model_spec/v1/models.py`.

## Important

Do not edit generated models manually. Update the schema and regenerate instead.
