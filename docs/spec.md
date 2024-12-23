# Model Format Specification

The specification defines an open standard for packaging and distribution Artificial Intelligence models as OCI artifacts, adhering to [the OCI image specification](https://github.com/opencontainers/image-spec/blob/main/spec.md#image-format-specification).

The goal of this specification is to outline a blueprint and enable the creation of interoperable solutions for packaging and retrieving AI/ML models by leveraging the existing OCI ecosystem, thereby facilitating efficient model management, deployment and serving in cloud-native environments.

## Use Cases

* An OCI Registry could storage and manage AI/ML model artifacts with model versions, metadata, and parameters retrievable and displayable.
* A Data Scientist can package models together with their metadata (e.g., format, precision) and upload them to a registry, facilitating collaboration with MLOps Engineers while streamlining the deployment process to efficiently deliver models into production.
* A model serving/deployment platform can read model metadata (e.g., format, precision) from a registry to understand the AI/ML model details, identify the required server runtime
  (as well as startup parameters, necessary resources, etc.), and serve the model in Kubernetes by [mounting it directly as a volume source](https://kubernetes.io/blog/2024/08/16/kubernetes-1-31-image-volume-source/)
  without needing to pre-download it in an init-container or bundle it within the server runtime container.

## Overview

At a high level, the Model Format Specification is based on the [OCI Image Format Specification](https://github.com/opencontainers/image-spec/blob/main/spec.md#image-format-specification) and incorporates [all its components](https://github.com/opencontainers/image-spec/blob/main/spec.md#understanding-the-specification). The key distinction lies in extending the [OCI Image Manifest Specification](https://github.com/opencontainers/image-spec/blob/main/manifest.md) to accommodate artifact usage specifically tailored for AI/ML models.

### Extended OCI Image Manifest Specification For Model Artifacts

The image manifest of model artifacts follows the [OCI Image Manifest Specification](https://github.com/opencontainers/image-spec/blob/main/manifest.md) and adheres to the [guidelines for artifacts usage](https://github.com/opencontainers/image-spec/blob/main/manifest.md#guidelines-for-artifact-usage). Specifically, it leverages the extensible `artifactType` and `annotations` properties to define attributes specific to model artifacts.

![manifest](./img/manifest.svg)


- **`artifactType`** _string_

  This REQUIRED property MUST be `application/vnd.cnai.model.manifest.v1+json`.

- **`layers`** _array of objects_

  - **`mediaType`** _string_

    This REQUIRED property MUST be one of the [OCI Image Media Types](https://github.com/opencontainers/image-spec/blob/main/media-types.md) designated for [layers](https://github.com/opencontainers/image-spec/blob/main/layer.md).
    Otherwise, it will not be compatible with the container runtime.

  - **`artifactType`** _string_

    This REQUIRED property MUST be at least the following media types:

    - `application/vnd.cnai.model.layer.v1.tar`: The layer is a [tar archive](https://en.wikipedia.org/wiki/Tar_(computing)) that contains the model weight file. If the model has multiple weight files, they SHOULD be packaged into separate layers.
    - `application/vnd.cnai.model.layer.v1.tar+gzip`: The layer is a [tar archive](https://en.wikipedia.org/wiki/Tar_(computing)) compressed with [gzip](https://datatracker.ietf.org/doc/html/rfc1952) that contains the model weight file.
      If the model has multiple weight files, they SHOULD be packaged in separate layers.
      
      _Implementers note_: It is recommended to package weight files without compression to avoid unnecessary overhead of decompression by the container runtime as model weight files are typically already compressed.
    - `application/vnd.cnai.model.doc.v1.tar`: The layer is a [tar archive](https://en.wikipedia.org/wiki/Tar_(computing)) that includes documentation files like `README.md`, `LICENSE`, etc.
    - `application/vnd.cnai.model.config.v1.tar`: The layer is a [tar archive](https://en.wikipedia.org/wiki/Tar_(computing)) that includes additional configuration files such as `config.json`，`tokenizer.json`, `generation_config.json`, etc.

  - **`annotations`** _string-string map_

    This OPTIONAL property contains arbitrary attributes for the layer. For metadata specific to models, implementations SHOULD use the predefined annotation keys as outlined in the [Layer Annotation Keys](./annotations.md#layer-annotation-keys).

## Workflow

As the model format specification conforms to the [OCI Image Specification](https://github.com/opencontainers/image-spec/blob/main/layer.md), it naturally aligns with the standard [OCI distribution workflow](https://github.com/opencontainers/distribution-spec/blob/main/spec.md).

This section outlines the typical workflow for a model OCI artifact, which consists of two main stages: `BUILD & PUSH` and `PULL & SERVE`.

### BUILD & PUSH

Build tools can package required resources into an OCI artifact following the model format specification.

The generated artifact can then be pushed to OCI registries (e.g., Harbor, DockerHub) for storage and management.

![build-push](./img/build-and-push.png)

### PULL & SERVE

Once the model artifact is stored in an OCI registry, the container runtime (e.g., containerd, CRI-O) can pull it from the OCI registry and mount it as a read-only volume during the model serving process, if required.

![pull-serve](./img/pull-and-serve.png)
