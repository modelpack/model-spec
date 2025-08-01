# Model Artifact Configuration

Each model artifact has an associated JSON structure which describes some basic information about the model such as name and version, as well as technical metadata such as format, precision and quantization. This content is referred to as _Model Artifact Configuration_ and is identified by the [media type][oci-media-type] `application/vnd.cncf.model.config.v1+json`.

This section defines `application/vnd.cncf.model.config.v1+json` media type.

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

  - **architecture_config** _object_, OPTIONAL

    The configuration of the architecture. The details are defined in the [Architecture](./architectures.md) file.

  - **format** _string_, OPTIONAL

    The format for the model, such as "onnx", "safetensors", "gguf", or "pt"(pytorch format).

  - **paramSize** _string_, OPTIONAL

    The model size is represented as a combination of a decimal `count` and a single-letter `scale-prefix` in the format of `<count><scale-prefix>`, which together specify the total number of parameters in the model.

    - `count`:
      A numeric value representing the base parameter count before scaling. This value may include up to one digit after the decimal point to allow for partial scaling precision. For example: `6.7`.

    - `scale-prefix`:
      A single letter indicating the order of magnitude multiplier applied to the count. The prefix is case-insensitive and must be one of the following:
      - `Q` or `q` (Quadrillion)
      - `T` or `t` (Trillion)
      - `B` or `b` (Billion)
      - `M` or `m` (Million)
      - `K` or `k` (Thousand)

    Some examples: `6.7B`(6.7 Billion parameters), `1.0t`(1 Trillion parameters), `100m`(100 Million parameters).

  - **precision** _string_, OPTIONAL

    The computational precision of the model. Supported values include:

    | Precision | Description |
    |-----------|-------------|
    | `"float32"` | 32-bit floating point |
    | `"float64"` | 64-bit floating point |
    | `"float16"` | 16-bit floating point. Uses 1 sign, 5 exponent, and 10 significand bits. |
    | `"bfloat16"` | 16-bit brain floating point. Uses 1 sign, 8 exponent and 7 significand bits. |
    | `"float8_e4m3"` | 8-bit floating point, e4m3 format. Uses 1 sign, 4 exponent, and 3 significand bits. |
    | `"float8_e5m2"` | 8-bit floating point, e5m2 format. Uses 1 sign, 5 exponent, and 2 significand bits. |
    | `"complex32"` | 32-bit complex |
    | `"complex64"` | 64-bit complex |
    | `"complex128"` | 128-bit complex |
    | `"int8"` | 8-bit signed integer |
    | `"int16"` | 16-bit signed integer |
    | `"int32"` | 32-bit signed integer |
    | `"int64"` | 64-bit signed integer |
    | `"uint8"` | 8-bit unsigned integer |
    | `"uint16"` | 16-bit unsigned integer |
    | `"uint32"` | 32-bit unsigned integer |
    | `"uint64"` | 64-bit unsigned integer |
    | `"bool"` | Boolean |

    If multiple precisions are used, they should be separated by commas. For example, if the model uses float16 and float8_e4m3, the precision should be set to `"float16,float8_e4m3"`.

  - **quantization** _string_, OPTIONAL

    Quantization technique applied to the model, such as "awq", or "gptq".

  - **capabilities** _object_, OPTIONAL

    Special capabilities that the model supports, such as reasoning, toolusage, etc.

- **modelfs** _object_, REQUIRED

  Contains hashes of each uncompressed layer's content.

  - **type** _string_, REQUIRED

    Must be set to "layers".

  - **diffIds** _array of strings_, REQUIRED

    An array of layer content hashes (`DiffIDs`), in order from first to last.

- **capabilities** _object_, OPTIONAL

    Special capabilities that the model supports, such as reasoning, toolusage, etc.

  - **inputTypes** _array of string_, OPTIONAL

    Input types that the model supports, such as "text", "image", "audio", "video", etc.

  - **outputTypes** _array of string_, OPTIONAL

    Output types that the model supports, such as "text", "image", "audio", "video", etc.

  - **knowledgeCutoff** _string_, OPTIONAL

    The date and time of the datasets that the model was trained on, formatted as defined by [RFC 3339, section 5.6][rfc3339-s5.6].

  - **reasoning** _boolean_, OPTIONAL

    Whether the model can perform reasoning tasks.

  - **toolUsage** _boolean_, OPTIONAL

    Whether the model can use external tools or APIs to perform tasks.

  - **embedding** _boolean_, OPTIONAL

    Whether the model can perform embedding tasks.

  - **reward** _boolean_, OPTIONAL

    Whether the model is a reward model.

## Example

Here is an example model artifact configuration JSON document:

```json,title=Model%20Config%20JSON&mediatype=application/vnd.cncf.model.config.v1%2Bjson
{
  "descriptor": {
    "createdAt": "2025-01-01T00:00:00Z",
    "authors": [
      "xyz@xyz.com"
    ],
    "vendor": "XYZ Corp.",
    "family": "xyz3",
    "name": "xyz-3-8B-Instruct",
    "version": "3.1",
    "title": "XYZ 3 8B Instruct",
    "description": "xyz is a large language model.",
    "docURL": "https://www.xyz.com/get-started/",
    "sourceURL": "https://github.com/xyz/xyz3",
    "revision": "1234567890",
    "licenses": [
      "Apache-2.0"
    ]
  },
  "config": {
    "architecture": "transformer",
    "format": "pytorch",
    "paramSize": "8b",
    "precision": "fp16",
    "quantization": "gptq",
    "capabilities": {
      "inputTypes": [
        "text"
      ],
      "outputTypes": [
        "text",
        "image"
      ],
      "knowledgeCutoff": "2024-05-21T00:00:00Z",
      "reasoning": true,
      "toolUsage": false,
      "embedding": false,
      "reward": false
    }
  },
  "modelfs": {
    "type": "layers",
    "diffIds": [
      "sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
      "sha256:abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
    ]
  }
}
```

[oci-media-type]: https://github.com/opencontainers/image-spec/blob/main/descriptor.md#properties
[rfc3339-s5.6]: https://tools.ietf.org/html/rfc3339#section-5.6
[spdx-license-expression]: https://spdx.github.io/spdx-spec/v2.3/SPDX-license-expressions/
