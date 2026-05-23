# Vision Encoder Specification

This document describes the vision encoder architecture fields for multimodal models that process image and video inputs. It extends the model configuration defined in [config.md](./config.md) to cover the architectural details of how visual inputs are processed.

## Background

The current ModelPack specification supports declaring image modality via `capabilities.inputTypes: ["image"]`, but provides no architectural description of how images are processed. Every major model family now has a vision variant (LLaVA, Qwen2-VL, LLaMA-3.2 Vision, Gemma 2 VL), and inference engines need structured metadata about the vision encoder to correctly configure image preprocessing, patch embedding, and vision-language fusion.

## Architecture Overview

Vision-language models follow a common pattern:

```text
Input Image → Vision Encoder → Projector → Language Model → Text Output
                    ↓
           Visual token embeddings
```

The **vision encoder** converts raw images into a sequence of visual tokens using a Vision Transformer (ViT) or CLIP-ViT architecture. A **projector** module maps these visual tokens into the language model's embedding space. The **fusion type** determines how visual and textual tokens interact inside the language model.

## Properties

- **type** _string_, REQUIRED

  The vision encoder architecture type. Supported values:

  | Value | Description |
  |-------|-------------|
  | `"vit"` | Standard Vision Transformer |
  | `"clip_vit"` | CLIP-pretrained Vision Transformer |
  | `"other"` | Other vision encoder architecture |

- **hidden_size** _integer_, REQUIRED

  The hidden size (embedding dimension) of the vision encoder.

- **patch_size** _integer_, REQUIRED

  The spatial patch size in pixels. For example, `14` means the image is divided into 14×14 pixel patches. Each patch becomes one visual token.

- **image_size** _integer_, REQUIRED

  The default input image resolution in pixels.

- **num_layers** _integer_, REQUIRED

  The number of transformer layers in the vision encoder.

- **num_attention_heads** _integer_, REQUIRED

  The number of attention heads in the vision encoder.

- **intermediate_size** _integer_, OPTIONAL

  The FFN intermediate size in the vision encoder.

- **in_channels** _integer_, OPTIONAL

  The number of input image channels. Defaults to `3` (RGB).

- **activation** _string_, OPTIONAL

  The activation function used in the vision encoder, such as `"quick_gelu"`, `"gelu"`, or `"silu"`.

- **norm** _object_, OPTIONAL

  Normalization configuration for the vision encoder.

  - **type** _string_, OPTIONAL

    The normalization type. Supported values: `"layernorm"`, `"rmsnorm"`.

  - **epsilon** _number_, OPTIONAL

    The epsilon value for normalization.

- **projector** _object_, OPTIONAL

  The multimodal projector that maps vision encoder outputs to the language model embedding space.

  - **type** _string_, OPTIONAL

    The projector architecture type. Supported values:

    | Value | Description |
    |-------|-------------|
    | `"mlp"` | Multi-layer perceptron (e.g., LLaVA 1.5 uses 2-layer MLP with GELU) |
    | `"linear"` | Single linear projection |
    | `"cross_attention"` | Cross-attention layers (e.g., LLaMA-3.2 Vision) |
    | `"perceiver"` | Perceiver-style resampler |
    | `"other"` | Other projector architecture |

  - **num_layers** _integer_, OPTIONAL

    The number of layers in the projector (for MLP or cross-attention type projectors).

  - **activation** _string_, OPTIONAL

    The activation function in the projector, such as `"gelu"`.

- **special_tokens** _object_, OPTIONAL

  Special token IDs used for image and video inputs in the tokenizer.

  - **image_token_id** _integer_, OPTIONAL

    The token ID used as a placeholder for image input in the text sequence.

  - **vision_start_token_id** _integer_, OPTIONAL

    The token ID marking the start of a vision region (used by models like Qwen2-VL).

  - **vision_end_token_id** _integer_, OPTIONAL

    The token ID marking the end of a vision region.

  - **vision_token_id** _integer_, OPTIONAL

    The token ID for a generic vision placeholder (used by models like Qwen2-VL).

  - **video_token_id** _integer_, OPTIONAL

    The token ID for video frame placeholders.

- **dynamic_resolution** _object_, OPTIONAL

  Dynamic image resolution support, where the model can handle variable-resolution inputs.

  - **enabled** _boolean_, OPTIONAL

    Whether dynamic resolution is enabled.

  - **min_pixels** _integer_, OPTIONAL

    The minimum number of visual tokens.

  - **max_pixels** _integer_, OPTIONAL

    The maximum number of visual tokens.

  - **spatial_merge_size** _integer_, OPTIONAL

    The spatial merging stride for reducing visual token count.

- **temporal_patch_size** _integer_, OPTIONAL

  The temporal patch size for video understanding. Specifies how many frames are grouped into one temporal patch.

- **fusion_type** _string_, OPTIONAL

  How vision and language modalities are fused. Supported values:

  | Value | Description |
  |-------|-------------|
  | `"early"` | Visual tokens are concatenated with text tokens before the first transformer layer (e.g., Qwen2-VL) |
  | `"late"` | Visual tokens are injected after separate encoding (e.g., LLaVA) |
  | `"cross_attention"` | Dedicated cross-attention layers between vision and language (e.g., LLaMA-3.2 Vision) |

- **position_embedding** _object_, OPTIONAL

  Position embedding configuration for the vision encoder.

  - **type** _string_, OPTIONAL

    The type of position embedding. Supported values: `"learned"`, `"rope"`, `"mrope"`, `"sinusoidal"`.

  - **mrope_sections** _array of integers_, OPTIONAL

    Per-modality RoPE dimension sections. Only applicable when type is `"mrope"` (e.g., Qwen2-VL uses `[16, 24, 24]` for temporal, height, width dimensions).

## Model Coverage

| Model | Encoder | Patch Size | Image Size | Projector | Fusion | Special Features |
|-------|---------|-----------|------------|-----------|--------|------------------|
| LLaVA 1.5 | CLIP-ViT-L/14 | 14 | 336 | 2-layer MLP | late | — |
| Qwen2-VL | ViT | 14 | dynamic | — | early | mRoPE, dynamic resolution, video |
| LLaMA-3.2 Vision | CLIP-ViT | 14 | 560 | cross-attention | cross_attention | Gated cross-attention |
| Gemma 2 VL | SigLIP | 14 | 224 | linear | late | — |

## Example

```json
{
  "type": "clip_vit",
  "hidden_size": 1024,
  "patch_size": 14,
  "image_size": 336,
  "num_layers": 24,
  "num_attention_heads": 16,
  "intermediate_size": 4096,
  "in_channels": 3,
  "activation": "quick_gelu",
  "norm": {
    "type": "layernorm",
    "epsilon": 1e-5
  },
  "projector": {
    "type": "mlp",
    "num_layers": 2,
    "activation": "gelu"
  },
  "special_tokens": {
    "image_token_id": 32000
  },
  "fusion_type": "late",
  "position_embedding": {
    "type": "learned"
  }
}
```
