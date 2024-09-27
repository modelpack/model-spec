# Architecture

## Tensor naming convention

[vendor].[family].[name].[arch].[modality].[block_name].[layer_number].[tensor_name].[tensor_type]

### Naming Conventions

- **vendor**: The vendor of the model.
- **family**: The family of the model.
- **name**: The name of the model.
- **arch**: The architecture of the model.
- **modality**: The modality of the model.
- **block_name**: The name of the block.
- **layer_number**: The layer number.
- **tensor_name**: The name of the tensor.
- **tensor_type**: The type of the tensor.

### Example

```
meta.llama.llama3.2-1B.transformer.text.decoder.layer.0.self_attention.query.weight
```
