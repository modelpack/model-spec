# Compatibility

This document describes the compatibility support between ModelPack and other model packaging formats.

## Docker Model Spec

The `compat` package provides conversion utilities between CNCF ModelPack and [Docker Model Spec](https://github.com/docker/model-spec) formats.

### Format Detection

Use `compat.DetectFormat()` to identify the format of a model artifact by its media type:

```go
import "github.com/modelpack/model-spec/compat"

format := compat.DetectFormat("application/vnd.cncf.model.config.v1+json")
// Returns: compat.FormatModelPack

format = compat.DetectFormat("application/vnd.docker.ai.model.config.v0.1+json")
// Returns: compat.FormatDocker
```

### Media Types

The `dockerv0` package provides constants for Docker model-spec media types:

```go
import dockerv0 "github.com/modelpack/model-spec/compat/docker/v0"

dockerv0.MediaTypeConfig        // application/vnd.docker.ai.model.config.v0.1+json
dockerv0.MediaTypeGGUF          // application/vnd.docker.ai.gguf.v3
dockerv0.MediaTypeLoRA          // application/vnd.docker.ai.gguf.v3.lora
dockerv0.MediaTypeMMProj        // application/vnd.docker.ai.gguf.v3.mmproj
dockerv0.MediaTypeLicense       // application/vnd.docker.ai.license
dockerv0.MediaTypeChatTemplate  // application/vnd.docker.ai.chat.template.jinja
```

### Converting ModelPack to Docker Format

```go
import (
    v1 "github.com/modelpack/model-spec/specs-go/v1"
    dockerv0 "github.com/modelpack/model-spec/compat/docker/v0"
)

model := v1.Model{
    Config: v1.ModelConfig{
        Format:       "gguf",
        ParamSize:    "8b",
        Architecture: "llama",
        Quantization: "Q4_0",
    },
    ModelFS: v1.ModelFS{
        Type:    "layers",
        DiffIDs: []digest.Digest{"sha256:abc123"},
    },
}

dockerCfg, err := dockerv0.FromModelPack(model)
// dockerCfg.ModelConfig.GGUF["parameter_count"] == "8 B"
// dockerCfg.ModelConfig.GGUF["architecture"] == "llama"
// dockerCfg.ModelConfig.GGUF["quantization"] == "Q4_0"
```

### Converting Docker Format to ModelPack

```go
import dockerv0 "github.com/modelpack/model-spec/compat/docker/v0"

dockerCfg := dockerv0.Config{
    ModelConfig: dockerv0.ModelConfig{
        Format: "gguf",
        Size:   "635992801",
        GGUF: map[string]any{
            "parameter_count": "8 B",
            "architecture":    "llama",
            "quantization":    "Q4_0",
        },
    },
    Files: []dockerv0.File{
        {DiffID: "sha256:abc123", Type: dockerv0.MediaTypeGGUF},
    },
}

model, err := dockerv0.ToModelPack(dockerCfg)
// model.Config.ParamSize == "8b"
// model.Config.Architecture == "llama"
// model.Config.Quantization == "Q4_0"
```

### Field Mapping

| ModelPack | Docker | Notes |
| --------- | ------ | ----- |
| `descriptor.createdAt` | `descriptor.createdAt` | RFC3339 format |
| `config.format` | `config.format` | Direct mapping |
| `config.paramSize` | `config.gguf.parameter_count` | Format conversion (e.g., "8b" ↔ "8 B") |
| `config.architecture` | `config.gguf.architecture` | Direct mapping |
| `config.quantization` | `config.gguf.quantization` | Direct mapping |
| `modelfs.diffIds` | `files[].diffID` | Structure conversion |

### Limitations

**Size field:** The Docker `config.size` field (total model size in bytes) cannot be derived from ModelPack format. When converting ModelPack to Docker, this field is set to "0". Callers should update this field if the actual size is known.

**Fields lost when converting ModelPack to Docker:**

- `descriptor.authors`
- `descriptor.name`, `descriptor.version`
- `descriptor.vendor`, `descriptor.licenses`
- `descriptor.family`, `descriptor.title`, `descriptor.description`
- `config.precision`
- `config.capabilities.*`

**Fields lost when converting Docker to ModelPack:**

- `config.format_version`
- `config.gguf.*` (except `parameter_count`, `architecture`, `quantization`)
