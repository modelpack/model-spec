# Model Artifact Configuration

Each model artifact has an associated JSON structure which describes some basic information about the model such as name and version, as well as technical metadata such as format, precision and quantization. This content is referred to as _Model Artifact Configuration_ and is identified by the [media type][oci-media-type] `application/vnd.cnai.model.config.v1+json`.

This section defines `application/vnd.cnai.model.config.v1+json` media type.

## Terminology

The following terms are used in this section:

- [Layer](./spec.md#guidance-on-layers)
- Layer DiffID

  A layer DiffID is the hash of the layer's uncompressed tar archive.

## Properties

- **descriptor** _object_, REQUIRED

  Contains the general information about the model.

  - **createdAt** _string_, OPTIONAL

    The date and time at which the model was created, formatted as defined by [RFC 3339, section 5.6][rfc3339-s5.6].

  - **authors** _array of strings_, OPTIONAL

    A list of contact details for the individuals or organizations responsible for the model (freeform string).

  - **vendor** _string_, OPTIONAL

    The name of the organization or company distributing the model.

  - **family** _string_, OPTIONAL

    The model family or lineage, such as "llama3", "gpt2", or "qwen2".

  - **name** _string_, OPTIONAL

    The name of the model.

  - **version** _string_, OPTIONAL

    The version of the model.

  - **title** _string_, OPTIONAL

    A human-readable title for the model.

  - **description** _string_, OPTIONAL

    A human-readable description of the model.

  - **docURL** _string_, OPTIONAL

    A URL to get more information or details about the model.

  - **sourceURL** _string_, OPTIONAL

    A URL to get the source code or resources needed to build or understand the model's implementation.

  - **revision** _string_, OPTIONAL

    The source control revision identifier for the model.

  - **licenses** _array of string_, OPTIONAL

    A list of licenses under which the model is distributed, represented as [SPDX License Expressions][spdx-license-expression].

- **config** _object_, REQUIRED

  Contains the technical metadata for the model.

  - **architecture** _string_, OPTIONAL

    The architecture of the model, such as "transformer", "cnn", or "rnn".

  - **format** _string_, OPTIONAL

    The format for the model, such as "onnx", "tensorflow", or "pytorch".

  - **parameterSize** _integer_, OPTIONAL

    The total number of parameters of the model parameters.

  - **precision** _string_, OPTIONAL

    The computational precision of the model, e.g., "bf16", "fp16", "int8", or "mixed".

  - **quantization** _string_, OPTIONAL

    Quantization technique applied to the model, such as "awq", or "gptq".

- **modelfs** _object_, REQUIRED

  Contains hashes of each uncompressed layer's content.

  - **type** _string_, REQUIRED

    Must be set to "layers".

  - **diff_ids** _array of strings_, REQUIRED

    An array of layer content hashes (`DiffIDs`), in order from first to last.

## Example

Here is an example model artifact configuration JSON document:

```json
{
    "descriptor": {
        "createdAt": "2025-01-01T00:00:00Z",
        "authors": ["xyz@xyz.com"],
        "vendor": "XYZ Corp.",
        "family": "xyz3",
        "name": "xyz-3-8B-Instruct",
        "version": "3.1",
        "title": "XYZ 3 8B Instruct",
        "description": "xyz is a large language model.",
        "docURL": "https://www.xyz.com/get-started/",
        "sourceURL": "https://github.com/xyz/xyz3",
        "revision": "1234567890",
        "licenses": ["Apache-2.0"]
    },
    "config": {
        "architecture": "transformer",
        "format": "pytorch",
        "parameterSize": "50000000000",
        "precision": "fp16",
        "quantization": "gptq"
    },
    "modelfs": {
        "type": "layers",
        "diff_ids": [
            "sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
            "sha256:abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
        ]
    }
}
```

[oci-media-type]: https://github.com/opencontainers/image-spec/blob/main/descriptor.md#properties
[rfc3339-s5.6]: https://tools.ietf.org/html/rfc3339#section-5.6
[spdx-license-expression]: https://spdx.github.io/spdx-spec/v2.3/SPDX-license-expressions/
