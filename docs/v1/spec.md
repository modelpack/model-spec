# Model Specification Version 1

The specification defines an open standard Artifacial Intelegence model. It is defined through the artifact extension based on [the OCI image specification](https://github.com/opencontainers/image-spec/blob/main/spec.md#image-format-specification), and extends model features through `artifactType` and `annotations`. Model storage and distribution can be optimized based on artifact extension.

The goal of this specification is to package models in an OCI artifact to take advantage of OCI distribution and ensure efficient model deployment.

The model specification needs to consider two factors:

1. The model needs to be stored in the OCI registry and display the parameters of the model. So that the model should use
   the [artifact extension](https://github.com/opencontainers/image-spec/blob/main/artifacts-guidance.md) to
   packaging content other than OCI image specification.
2. The model needs to be mounted by the container runtime as
   [read only volumes based on the OCI Artifacts in Kubernetes 1.31+](https://kubernetes.io/blog/2024/08/16/kubernetes-1-31-image-volume-source/).
   Container runtimes can only pull OCI artifact that follows the OCI image specification.

Therefore, the model specification must be defined through the artifact extension based on the [OCI image specification](https://github.com/opencontainers/image-spec/blob/main/spec.md#image-format-specification). It can be better compatible with the kubernetes ecosystem.

## Overview

The model specification is defined through the artifact extension based on the OCI image specification, and extend model features through `artifactType` and `annotations`. Model storage and distribution can be optimized based on artifact extension.

![manifest](../img/v1/manifest.svg)

## Workflow

The model specification running workflow is divided into two stages: `BUILD & PUSH` and `PULL & SERVE`.

### BUILD & PUSH

Use tools(ORAS, Ollama, etc.) to build required resources in the model repository into artifact based on the model specification. Note that the model layer MUST NOT be compressed, because the files of model weight has been compressed. If the model layer is compressed, the container runtime will cost long time to decompress the model layer. Therefore, it's recommended to use the `application/vnd.oci.image.layer.v1.tar` format for the model layer to avoid compression

Next push the artifact to the OCI registry(Harbor, Docker Hub, etc.), and use the functionalities of the OCI registry to manage the model artifact.

![build-push](../img/v1/build-and-push.png)

### PULL & SERVE

The container runtime(containerd, cri-o, etc) pulls the model artifact from the OCI registry, and mounts the model artifact as a read-only volume. Therefore, distributed model can use the P2P technology(Dragonfly, Kraken, etc) to reduce the pressure on the registry and preheat the model artifact into each node. If the model artifact is already present on the node, the container runtime can reuse the model artifact to mount different containers in the same node.

![pull-serve](../img/v1/pull-and-serve.png)

## Understanding the Specification

The model specification is based on the [OCI image specification](https://github.com/opencontainers/image-spec/blob/main/spec.md) and focuses on defining the artifact extension according to the [Artifacts Guidance](https://github.com/opencontainers/image-spec/blob/main/artifacts-guidance.md).

### Image Manifest Extension Properties

- **`artifactType`** _string_

  This REQUIRED property MUST contain the media type `application/vnd.cnai.model.manifest.v1+json`.

- **`layers`** _array of objects_

  - **`artifactType`** _string_

    Implementations MUST support at least the following media types:

    - `application/vnd.cnai.model.layer.v1.tar`
    - `application/vnd.cnai.model.layer.v1.tar+gzip`

    If `mediaType` is `application/vnd.oci.image.layer.v1.tar`, the `artifactType` MUST be `application/vnd.cnai.model.layer.v1.tar`. If `mediaType` is `application/vnd.oci.image.layer.v1.tar+gzip`, the `artifactType` MUST be `application/vnd.cnai.model.layer.v1.tar+gzip`. The `mediaType` and `artifactType` MUST be consistent, for detailed definitions of Filesystem Layers, please refer to the [Image Layer Filesystem Changeset](https://github.com/opencontainers/image-spec/blob/main/layer.md).

  - **`annotations`** _string-string map_

    This OPTIONAL property contains arbitrary metadata for the layer. For model specification, SHOULD set the pre-defined annotation keys, refer to the [Layer Annotation Keys](./annotations.md#layer-annotation-keys).

- **`annotations`** _string-string map_

  This OPTIONAL property contains arbitrary metadata for the image manifest. For model specification, SHOULD set the pre-defined annotation keys, refer to the [Manifest Annotation Keys](./annotations.md#manifest-annotation-keys).
