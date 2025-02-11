# Annotations

This property contains arbitrary metadata, and SHOULD follow the rules of [OCI image annotations](https://github.com/opencontainers/image-spec/blob/main/annotations.md).

## Pre-defined Annotation Keys

### Manifest Annotation Keys

#### Model-common Annotation Keys

- **`org.cnai.model.created`**: Date and time on which the model was built, conforming to [RFC 3339](https://tools.ietf.org/html/rfc3339#section-5.6).
- **`org.cnai.model.authors`**: Contact details of the people or organization responsible for the model (freeform string).
- **`org.cnai.model.url`**: URL to find more information on the model (string).
- **`org.cnai.model.documentation`**: URL to get documentation on the model (string).
- **`org.cnai.model.source`**: URL to get source code for building the model (string).
- **`org.cnai.model.version`**: Version of the model.
- **`org.cnai.model.revision`**: Source control revision identifier for the model.
- **`org.cnai.model.vendor`**: Name of the distributing entity, organization or individual.
- **`org.cnai.model.licenses`**: License(s) under which contained software is distributed as an [SPDX License Expression](https://spdx.github.io/spdx-spec/v2.3/SPDX-license-expressions/).
- **`org.cnai.model.ref.name`**: Name of the reference for a target (string).
- **`org.cnai.model.title`**: Human-readable title of the model (string).
- **`org.cnai.model.description`**: Human-readable description of the software packaged in the model (string).

#### Model-specific Annotation Keys

- **`org.cnai.model.architecture`**: Model architecture (string), such as `transformer`, `cnn`, `rnn`, etc.
- **`org.cnai.model.family`**: Model family (string), such as `llama3`, `gpt2`, `qwen2`, etc.
- **`org.cnai.model.name`**: Model name (string), such as `llama3-8b-instruct`, `gpt2-xl`, `qwen2-vl-72b-instruct`, etc.
- **`org.cnai.model.format`**: Model format (string), such as `onnx`, `tensorflow`, `pytorch`, etc.
- **`org.cnai.model.param.size`**: Number of parameters in the model (integer).
- **`org.cnai.model.precision`**: Model precision (string), such as `bf16`, `fp16`, `int8`, etc.
- **`org.cnai.model.quantization`**: Model quantization (string), such as `awq`, `gptq`, etc.

### Layer Annotation Keys

- **`org.cnai.model.filepath`**: Specifies the file path of the layer (string).
