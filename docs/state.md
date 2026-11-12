# Model State Artifacts

A model state artifact packages runtime state, such as a prompt cache or saved conversation. Its manifest's `subject` field references the model artifact it was built from.

State contents are runtime-specific; [`format`](#properties) identifies the encoding. States built from images, audio or other non-token input, or with an adapter applied, are out of scope.

## Use Cases

- A long system prompt or shared document is processed once when the state is built, and replicas, including new ones during autoscaling, load the state instead of processing it again.
- A saved conversation resumes on another machine.

## Loading a State

A runtime that loads a state SHOULD produce the same output as it would by processing the state's [`tokens`](#properties) without the state.

A state MUST NOT be loaded for a model other than the one its `subject` references. A runtime MUST NOT load a state whose `format` it does not support.

A runtime cannot check that a state's layers match its `tokens` without processing the tokens, which is the work the state saves. Runtimes SHOULD load only states they trust, for example states with a verified signature or states referenced by digest, and SHOULD NOT load a state only because it is listed as a referrer of the model.

## Image Manifest

The image manifest of a model state artifact follows the [OCI Image Manifest Specification][image-manifest].

- **`mediaType`** _string_

  This REQUIRED property MUST be `application/vnd.oci.image.manifest.v1+json`.

- **`artifactType`** _string_

  This REQUIRED property MUST be `application/vnd.cncf.model.state.manifest.v1+json`.

- **`subject`** _descriptor_

  This REQUIRED property references the manifest of the [model artifact](./spec.md) the state was built from.

- **`config`** _[descriptor](#configuration)_

  This REQUIRED property references the configuration of the state, by digest.

  - **`mediaType`** _string_

    This REQUIRED property MUST be `application/vnd.cncf.model.state.config.v1+json`.

- **`layers`** _array of objects_

  - **`mediaType`** _string_

    Implementations SHOULD support the following media types:

    - `application/vnd.cncf.model.state.v1.raw`: The layer is an unarchived, uncompressed state file.

    - `application/vnd.cncf.model.state.v1.tar`: The layer is a [tar archive][tar-archive] that contains state files.

    - `application/vnd.cncf.model.state.v1.tar+gzip`: The layer is a [tar archive][tar-archive] that contains state files. The archive is compressed with [gzip][rfc1952_2].

    - `application/vnd.cncf.model.state.v1.tar+zstd`: The layer is a [tar archive][tar-archive] that contains state files. The archive is compressed with [zstd][rfc8878].

  - **`annotations`** _string-string map_

    This OPTIONAL property contains arbitrary attributes for the layer. For metadata specific to models, implementations SHOULD use the predefined annotation keys as outlined in the [Layer Annotation Keys](./annotations.md#layer-annotation-keys).

  Layers follow the [Guidance on Layers](./spec.md#guidance-on-layers).

### Example Image Manifest

```json
{
    "schemaVersion": 2,
    "mediaType": "application/vnd.oci.image.manifest.v1+json",
    "artifactType": "application/vnd.cncf.model.state.manifest.v1+json",
    "subject": {
        "mediaType": "application/vnd.oci.image.manifest.v1+json",
        "digest": "sha256:f7db0c7f34e97395f5de3060198f05f701aa7c6f601674c738698bdea7282146",
        "size": 1043
    },
    "config": {
        "mediaType": "application/vnd.cncf.model.state.config.v1+json",
        "digest": "sha256:968e736d678755ecbb8b493177aee26bc93e25e9fc4c9fa8e614db4a59e3a0bc",
        "size": 124
    },
    "layers": [
        {
            "mediaType": "application/vnd.cncf.model.state.v1.tar",
            "digest": "sha256:813cee68413fd5a3dd00aafcdb49308f4b58553b991686c87b42532cb3aebabf",
            "size": 52428800
        }
    ]
}
```

## Configuration

The configuration of a model state artifact is a JSON object identified by the media type `application/vnd.cncf.model.state.config.v1+json`.

### Properties

- **format** _string_, REQUIRED

  The format of the state's layers and its version, such as "example.session/v1".

- **tokens** _array of integers_, REQUIRED

  The token IDs, from the model's vocabulary, that the state was built from, in order. Anyone who can pull a state can read the input it was built from, through `tokens` and often its layers.

- **precision** _string_, OPTIONAL

  The precision of the state's tensors, as one of the values defined for the model's [`precision`](./config.md#properties) or as a format-specific type such as `q8_0`. It can differ from the model's precision: a `bfloat16` model may use a `float8_e4m3` key-value cache.

### Example Configuration

```json,title=Model%20State%20Config%20JSON&mediatype=application/vnd.cncf.model.state.config.v1%2Bjson
{
    "format": "example.session/v1",
    "tokens": [1, 887, 526, 263, 8444, 20255, 29889, 13],
    "precision": "float16"
}
```

## Discovery

Clients find the states of a model by listing the [referrers][referrers-api] of its manifest, filtered by the artifact type `application/vnd.cncf.model.state.manifest.v1+json`. For registries that do not support the referrers API, clients use the [referrers tag schema][referrers-tag-schema].

The link through `subject` is advisory. Pulling a state does not pull its model, copying a model copies its states only with tools that also copy referrers, and some registries delete a model's states along with the model.

[image-manifest]: https://github.com/opencontainers/image-spec/blob/main/manifest.md
[referrers-api]: https://github.com/opencontainers/distribution-spec/blob/main/spec.md#listing-referrers
[referrers-tag-schema]: https://github.com/opencontainers/distribution-spec/blob/main/spec.md#referrers-tag-schema
[tar-archive]: https://en.wikipedia.org/wiki/Tar_(computing)
[rfc1952_2]: https://tools.ietf.org/html/rfc1952
[rfc8878]: https://www.rfc-editor.org/rfc/rfc8878
