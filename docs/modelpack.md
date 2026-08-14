# Running ModelPack Models with Docker Model Runner

This guide explains how to run CNCF ModelPack artifacts with [Docker Model Runner](https://github.com/docker/model-runner).

## What ModelPack compatibility means

Docker Model Runner includes ModelPack compatibility introduced in [docker/model-runner#516](https://github.com/docker/model-runner/pull/516).

This allows Docker Model Runner to pull and run model artifacts that follow the ModelPack specification while keeping the normal Docker Model Runner workflow.

## Run a ModelPack model

1. Package and push a model artifact using a ModelPack-compatible tool (for example, [modctl](https://github.com/modelpack/modctl) or [AIKit](https://kaito-project.github.io/aikit/docs/packaging)).
2. Use the artifact reference from your OCI registry with Docker Model Runner.
3. Run inference with Docker Model Runner.

## Example workflow

```bash
# Example ModelPack artifact reference in an OCI registry
MODEL_REF=registry.example.com/models/my-model:v1

# Run inference via Docker Model Runner
docker model run $MODEL_REF "Hello"
```

For additional CLI and API usage, see:

- [docker/model-runner repository](https://github.com/docker/model-runner)
- [Docker Model Runner docs](https://docs.docker.com/ai/model-runner/get-started/)

## Compatibility notes

- Use a Docker Model Runner version that includes [docker/model-runner#516](https://github.com/docker/model-runner/pull/516) or newer.
- The model artifact must be a valid OCI artifact following the ModelPack specification.
- Backend/runtime-specific behavior is defined by Docker Model Runner.
