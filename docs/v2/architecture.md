# Architecture

## Tensor naming convention

[version].[vendor].[family].[name].[arch].[modality].[block_name].[layer_name].[tensor_name].[tensor_type]

The dot in the name should be replaced with a underscore.

### Naming Conventions

- **version**: The version of the naming convention.
- **vendor**: The vendor of the model.
- **family**: The family of the model.
- **name**: The name of the model.
- **arch**: The architecture of the model.
- **modality**: The modality of the model.
- **block_name**: The name of the block.
- **layer_name**: The name and 0-indexed layer number of the layer.
- **tensor_name**: The name of the tensor.
- **tensor_type**: Weight or bias of the tensor.

### Example

```plain
v1.meta.llama-3_2-1b.transformer.text.decoder.layers_0.embedding.projection.weight
```

```plain
v1.meta.llama-3_2-1b.transformer.text.decoder.layers_1.attention.query.weight
```
