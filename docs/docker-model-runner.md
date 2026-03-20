# Using Docker Model Runner with ModelPack

This guide shows you how to use [Docker Model Runner](https://docs.docker.com/desktop/features/model-runner/) to pull and run AI models packaged using the ModelPack specification.

## What is Docker Model Runner?

Docker Model Runner is a built-in feature of Docker Desktop that enables pulling, managing, and running AI models directly from OCI registries. It natively supports the ModelPack specification format, allowing you to run ModelPack-packaged models without any additional tools.

## Prerequisites

- [Docker Desktop](https://docs.docker.com/get-docker/) 4.40 or later with Model Runner enabled
- A ModelPack-compatible model pushed to an OCI registry (see [modctl](./modctl.md) or [AIKit](./aikit.md) for packaging)

## Enable Docker Model Runner

Docker Model Runner is available through Docker Desktop. Enable it in Docker Desktop settings:

1. Open Docker Desktop
2. Go to **Settings** > **Features in development**
3. Enable **Docker Model Runner**

You can verify it is enabled by running:

```bash
docker model list
```

## Pull a ModelPack Model

Docker Model Runner can pull models directly from OCI registries. When pulling a ModelPack-formatted artifact, Docker automatically detects the ModelPack config format and converts it for local use.

```bash
# Pull a model from an OCI registry
docker model pull myregistry.com/mymodel:v1.0
```

## Run a Model

Once pulled, you can run inference using the model:

```bash
# Run a model interactively
docker model run myregistry.com/mymodel:v1.0

# Send a prompt to the model
docker model run myregistry.com/mymodel:v1.0 "Explain cloud-native computing"
```

## List and Manage Models

```bash
# List all downloaded models
docker model list

# Remove a model
docker model rm myregistry.com/mymodel:v1.0
```

## Use Models via the OpenAI-Compatible API

Docker Model Runner exposes an OpenAI-compatible API endpoint, enabling integration with existing tools and libraries:

```bash
curl http://localhost:12434/engines/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "myregistry.com/mymodel:v1.0",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'
```

## How ModelPack Format Is Detected

Docker Model Runner identifies a ModelPack artifact by checking the OCI config blob for any of the following fields:

- `config.paramSize` — the model parameter size
- `descriptor.createdAt` — the model creation timestamp
- `modelfs` — the model filesystem descriptor

If any of these fields are present, the artifact is recognized as a ModelPack-formatted model.

## Field Mapping: ModelPack to Docker

When Docker Model Runner pulls a ModelPack model, it converts the config fields to Docker's internal format:

| ModelPack Field | Docker Field | Description |
|---|---|---|
| `descriptor.createdAt` | `created` | Model creation timestamp |
| `descriptor.name` | `descriptor.name` | Model name |
| `descriptor.family` | `descriptor.family` | Model family |
| `descriptor.description` | `descriptor.description` | Model description |
| `descriptor.licenses` | `descriptor.licenses` | License information |
| `config.paramSize` | `parameters` | Model parameter count |
| `config.format` | `config.format` | Model format (e.g., GGUF) |
| `config.quantization` | `config.quantization` | Quantization method |
| `config.architecture` | `config.architecture` | Model architecture |
| `modelfs` | `rootfs` | Layer content addresses |

## Media Type Mapping

ModelPack media types are converted to Docker's internal media types:

| ModelPack Media Type | Docker Media Type |
|---|---|
| `application/vnd.cncf.model.weight.v1.raw` | Mapped based on file extension (e.g., `.gguf` → `application/vnd.docker.ai.gguf.v3`) |
| `application/vnd.cncf.model.weight.v1.tar+gzip` | `application/vnd.docker.ai.gguf.v3+gzip` |
| `application/vnd.cncf.model.weight.config.v1.raw` | `application/vnd.docker.ai.config` |
| `application/vnd.cncf.model.doc.v1.raw` | `application/vnd.docker.ai.doc` |

## Next Steps

- **Package models** using [modctl](./modctl.md) or [AIKit](./aikit.md) to create ModelPack artifacts
- **Learn about the [Model CSI Driver](https://github.com/modelpack/model-csi-driver)** for Kubernetes integration
- **Read the [full ModelPack specification](./spec.md)** for technical implementation details
