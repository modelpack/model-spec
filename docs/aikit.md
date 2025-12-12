# Using AIKit with ModelPack

This guide shows you how to use AIKit to package AI models as OCI artifacts using the ModelPack specification.

## What is AIKit?

[AIKit](https://kaito-project.github.io/aikit/docs/) is a comprehensive platform to quickly get started to host, deploy, build and fine-tune large language models (LLMs). AIKit also provides packaging models as OCI artifacts for distribution through any OCI-compliant registry.

## Prerequisites

- Docker with [BuildKit](https://docs.docker.com/build/buildkit/) support
- [ORAS](https://oras.land/docs/installation) or [Skopeo](https://github.com/containers/skopeo/blob/main/install.md) for pushing to registries

## Package a Model

AIKit uses Docker BuildKit to package models from various sources (local files, HTTP/HTTPS, or Hugging Face).

### Example: Package from Hugging Face

```bash
export HF_MODEL="Qwen/Qwen3-0.6B"
export MODEL_NAME="qwen3"
export OUTPUT_DIR="qwen"

docker buildx build \
  --build-arg BUILDKIT_SYNTAX=ghcr.io/kaito-project/aikit/aikit:latest \
  --target packager/modelpack \
  --build-arg source=huggingface://$HF_MODEL \
  --build-arg name=$MODEL_NAME \
  --output=$OUTPUT_DIR - <<< ""
```

For more packaging options including compression modes, layer categorization, and exclusions, see the [AIKit packaging documentation](https://kaito-project.github.io/aikit/docs/packaging).

## Push to a Registry

Use ORAS or Skopeo to push the OCI layout to a remote registry:

```bash
export REGISTRY="myregistry.com/mymodel:v1.0"

# Using ORAS
oras cp --from-oci-layout $OUTPUT_DIR/layout:$MODEL_NAME $REGISTRY

# Or using Skopeo
skopeo copy oci:$OUTPUT_DIR/layout docker://$REGISTRY
```

## Pull from a Registry

Pull models using ORAS or Skopeo:

```bash
export REGISTRY="myregistry.com/mymodel:v1.0"

# Using ORAS (preserves file names automatically)
oras pull $REGISTRY --output path/to/model/

# Or using Skopeo
skopeo copy docker://$REGISTRY dir://path/to/model/
# rename files based on annotations
(
  cd path/to/model/
  for digest in $(jq -r '.layers[].digest' manifest.json); do
    name=$(jq -r --arg digest "$digest" '.layers[] | select(.digest==$digest) | .annotations["org.cncf.model.filepath"]' manifest.json)
    if [ "$name" != "null" ]; then mv "${digest#sha256:}" "$name"; fi
  done
)
```

## Next Steps

- **See the [AIKit packaging documentation](https://kaito-project.github.io/aikit/docs/packaging)** for more information on packaging options
- **Learn about the [Model CSI Driver](https://github.com/modelpack/model-csi-driver)** for Kubernetes integration
- **Read the [full ModelPack specification](./spec.md)** for technical implementation details
