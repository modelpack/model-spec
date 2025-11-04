# Using modctl with ModelPack

This guide shows you how to use `modctl` to package, distribute, and manage AI models using the ModelPack specification.

## Installation

Follow the instructions to install `modctl` from the [modctl GitHub repository](https://github.com/modelpack/modctl/blob/main/docs/getting-started.md#installation) to install the CLI tool.

## Download A Model

To package a model, you need to download it to your local directory. The following example shows how to download a model from Huggingface.

```bash
export HF_MODEL="Qwen/Qwen3-0.6B"
export MODEL_PATH=my-model-directory

# Install the huggingface cli
pip install 'huggingface_hub'

# Login the huggingface cli
hf auth login --token <your-huggingface-token>

# Download a model
hf download $HF_MODEL --local-dir $MODEL_PATH
```

## Package Your First Model

The following script will walk through how to build a ModelPack format model artifact and push it to the model registry.

```bash
# Please modify the MODEL_REGISTRY environment variable to point to your OCI model registry
export MODEL_REGISTRY=myregistry.com

# If $MODEL_REGISTRY needs authentication, please login first
modctl login -u <username> -p <password> $MODEL_REGISTRY

# Generate a sample Modelfile, and edit the fields as needed
modctl modelfile generate $MODEL_PATH

# Build a model artifact from your model files
modctl build -t $MODEL_REGISTRY/mymodel:v1.0 $MODEL_PATH

# Push to an OCI registry
modctl push $MODEL_REGISTRY/mymodel:v1.0
```

## Next Steps

- **Explore more [modctl commands](https://github.com/modelpack/modctl)** for additional functionality
- **Learn about the [Model CSI Driver](https://github.com/modelpack/model-csi-driver)** for Kubernetes integration
- **Read the [full ModelPack specification](./spec.md)** for technical implementation details
