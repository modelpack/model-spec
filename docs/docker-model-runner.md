# Using Docker Model Runner with ModelPack

This guide shows you how to use [Docker Model Runner](https://docs.docker.com/desktop/features/model-runner/) to pull and run AI models packaged using the ModelPack specification. Note that runtime compatibility depends on the model format and inference engine combinations that Docker Model Runner supports (currently GGUF models).

## What is Docker Model Runner?

Docker Model Runner is a built-in feature of Docker Desktop that enables pulling, managing, and running AI models directly from OCI registries. Since [v1.0.19](https://github.com/docker/model-runner/releases/tag/v1.0.19), it can detect and convert ModelPack-formatted OCI artifacts, allowing you to use ModelPack-packaged models that are in a supported format (e.g., GGUF) without any additional tools.

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

Once pulled, you can run inference using the model. The model must be pulled before running — `docker model run` does not pull automatically:

```bash
# First pull the model (required before running)
docker model pull myregistry.com/mymodel:v1.0

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

When Docker Model Runner pulls a ModelPack model, it converts the config to Docker's internal format. Some fields are remapped to different names or locations:

**Fields that are renamed (breaking changes if modified):**

| ModelPack Field | Docker Field | Notes |
|---|---|---|
| `descriptor.createdAt` | `created` (top-level) | Moved out of `descriptor` and renamed |
| `config.paramSize` | `parameters` (top-level) | Moved out of `config` and renamed |
| `modelfs` (top-level object) | `rootfs` (top-level object) | Renamed at top level |

**Fields that are passed through unchanged (within their parent objects):**

| ModelPack Field | Docker Field | Notes |
|---|---|---|
| `descriptor.name` | `descriptor.name` | Kept in `descriptor` object |
| `descriptor.family` | `descriptor.family` | Kept in `descriptor` object |
| `descriptor.description` | `descriptor.description` | Kept in `descriptor` object |
| `descriptor.licenses` | `descriptor.licenses` | Kept in `descriptor` object |
| `config.format` | `config.format` | Kept in `config` object |
| `config.quantization` | `config.quantization` | Kept in `config` object |
| `config.architecture` | `config.architecture` | Kept in `config` object |

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
