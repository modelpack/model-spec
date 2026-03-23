#!/usr/bin/env python3
"""Parse HuggingFace model config.json into ModelPack transformer spec format.

This tool maps HuggingFace Transformers config.json fields to the ModelPack
unified transformer specification vocabulary defined in PR #111
(docs/architecture.md by @aftersnow).

Usage:
    python tools/hf_parser.py meta-llama/Meta-Llama-3-8B
    python tools/hf_parser.py mistralai/Mistral-7B-v0.3
    python tools/hf_parser.py --file path/to/config.json

The output is a YAML spec file following the ModelPack transformer spec format.
Fields that cannot be reliably inferred from config.json are marked as
NEEDS_REVIEW for human verification.
"""

from __future__ import annotations

import argparse
import json
import sys
from pathlib import Path

NEEDS_REVIEW = "__NEEDS_REVIEW__"

# Maps HuggingFace config.json field names to ModelPack transformer spec paths.
# Based on PR #111's field vocabulary (docs/architecture.md).
FIELD_MAP = {
    # Top-level transformer fields
    "vocab_size": "vocabulary_size",
    "hidden_size": "hidden_size",
    # Position embedding
    "max_position_embeddings": "position_embedding.max_position_embeddings",
    "rope_theta": "position_embedding.rope_theta",
    "rope_scaling": "position_embedding.rope_scaling",
    # Attention
    "num_attention_heads": "attention.num_attention_heads",
    "num_key_value_heads": "attention.num_key_value_heads",
    "head_dim": "attention.head_dim",
    # FFN / MLP
    "intermediate_size": "mlp.intermediate_size",
    # Transformer layers
    "num_hidden_layers": "num_layers",
    # Normalization
    "rms_norm_eps": "norm.epsilon",
    # MoE fields
    "num_local_experts": "moe.num_experts",
    "num_experts_per_tok": "moe.top_k",
    "num_experts": "moe.num_experts",
    "n_routed_experts": "moe.num_experts",  # DeepSeek naming variant
    # MLA fields (DeepSeek)
    "kv_lora_rank": "attention.kv_lora_rank",
    "q_lora_rank": "attention.q_lora_rank",
    "qk_nope_head_dim": "attention.qk_nope_head_dim",
    "qk_rope_head_dim": "attention.qk_rope_head_dim",
    "v_head_dim": "attention.v_head_dim",
}

# Known model type → attention type mapping
ATTENTION_TYPE_MAP = {
    "llama": "gqa",
    "mistral": "gqa",
    "mixtral": "gqa",
    "qwen2": "gqa",
    "qwen2_moe": "gqa",
    "gemma": "gqa",
    "gemma2": "gqa",
    "phi3": "gqa",
    "deepseek_v2": "mla",
    "deepseek_v3": "mla",
    "gpt2": "mha",
    "gpt_neo": "mha",
    "gpt_neox": "mha",
    "falcon": "mha",
}

# Known model type → FFN type mapping
FFN_TYPE_MAP = {
    "llama": "mlp",
    "mistral": "mlp",
    "mixtral": "moe",
    "qwen2": "mlp",
    "qwen2_moe": "moe",
    "gemma": "mlp",
    "gemma2": "mlp",
    "phi3": "mlp",
    "deepseek_v2": "moe",
    "deepseek_v3": "moe",
    "gpt2": "mlp",
    "gpt_neo": "mlp",
    "gpt_neox": "mlp",
    "falcon": "mlp",
}

# Known model type → activation function mapping
ACTIVATION_MAP = {
    "llama": "silu",
    "mistral": "silu",
    "mixtral": "silu",
    "qwen2": "silu",
    "qwen2_moe": "silu",
    "gemma": "gelu",
    "gemma2": "gelu",
    "phi3": "silu",
    "gpt2": "gelu",
    "gpt_neo": "gelu",
    "gpt_neox": "gelu",
    "falcon": "gelu",
}


def _set_nested(d: dict, path: str, value) -> None:
    """Set a value in a nested dict using dot-separated path."""
    keys = path.split(".")
    for key in keys[:-1]:
        d = d.setdefault(key, {})
    d[keys[-1]] = value


def _get_nested(d: dict, path: str, default=None):
    """Get a value from a nested dict using dot-separated path."""
    keys = path.split(".")
    for key in keys:
        if not isinstance(d, dict) or key not in d:
            return default
        d = d[key]
    return d


def parse_hf_config(raw: dict) -> dict:
    """Parse a HuggingFace config.json dict into ModelPack transformer spec.

    Args:
        raw: The parsed config.json dict from HuggingFace.

    Returns:
        A dict following the ModelPack transformer spec format.
    """
    result: dict = {}
    model_type = raw.get("model_type", "").lower()

    # Map static fields
    for hf_key, mp_path in FIELD_MAP.items():
        if hf_key in raw and raw[hf_key] is not None:
            _set_nested(result, mp_path, raw[hf_key])

    # Derive head_dim if absent
    if "attention" in result and "head_dim" not in result.get("attention", {}):
        hidden = result.get("hidden_size")
        n_heads = _get_nested(result, "attention.num_attention_heads")
        if hidden and n_heads:
            _set_nested(result, "attention.head_dim", hidden // n_heads)

    # Set architecture type
    result["type"] = "decoder"
    result["architecture_version"] = "0.1.0"

    # Infer attention type from model_type
    attn_type = ATTENTION_TYPE_MAP.get(model_type, NEEDS_REVIEW)
    _set_nested(result, "attention.type", attn_type)
    _set_nested(result, "attention.is_causal", True)

    # Check for sliding window attention
    if raw.get("sliding_window") is not None:
        _set_nested(result, "attention.sliding_window", raw["sliding_window"])

    # Infer FFN type
    ffn_type = FFN_TYPE_MAP.get(model_type, NEEDS_REVIEW)
    result["ffn_type"] = ffn_type

    # Set activation function
    hf_activation = raw.get("hidden_act", raw.get("activation_function"))
    if hf_activation:
        activation = hf_activation.lower()
        if "silu" in activation or "swish" in activation:
            activation = "silu"
        elif "gelu" in activation:
            activation = "gelu"
        elif "relu" in activation:
            activation = "relu"
    else:
        activation = ACTIVATION_MAP.get(model_type, NEEDS_REVIEW)

    if ffn_type == "mlp":
        _set_nested(result, "mlp.activation", activation)
        # Most modern models use gated activation (SwiGLU, GeGLU)
        use_gated = model_type in (
            "llama", "mistral", "mixtral", "qwen2", "qwen2_moe", "phi3",
            "gemma", "gemma2", "deepseek_v2", "deepseek_v3",
        )
        _set_nested(result, "mlp.use_gated_activation", use_gated)
    elif ffn_type == "moe":
        _set_nested(result, "moe.activation", activation)
        # MoE-specific fields
        if "moe_intermediate_size" in raw:
            _set_nested(result, "moe.moe_intermediate_size", raw["moe_intermediate_size"])
        if "num_shared_experts" in raw:
            _set_nested(result, "moe.num_shared_experts", raw["num_shared_experts"])
        if "shared_expert_intermediate_size" in raw:
            _set_nested(
                result, "moe.shared_expert_intermediate_size",
                raw["shared_expert_intermediate_size"],
            )
        # DeepSeek MoE-specific fields (from PR #185 research)
        if "routed_scaling_factor" in raw:
            _set_nested(result, "moe.routed_scaling_factor", raw["routed_scaling_factor"])
        if "topk_method" in raw:
            _set_nested(result, "moe.topk_method", raw["topk_method"])
        if "norm_topk_prob" in raw:
            _set_nested(result, "moe.norm_topk_prob", raw["norm_topk_prob"])

    # Mixed layers support (DeepSeek uses dense layers before switching to MoE)
    if "first_k_dense_replace" in raw and "moe_layer_freq" in raw:
        result["layer_structure"] = "mixed"
        _set_nested(result, "mixed_layers.first_k_dense_replace", raw["first_k_dense_replace"])
        _set_nested(result, "mixed_layers.moe_layer_freq", raw["moe_layer_freq"])

    # Normalization
    norm_type = "rmsnorm"  # Most modern models use RMSNorm
    if model_type in ("gpt2", "gpt_neo"):
        norm_type = "layernorm"
    _set_nested(result, "norm.type", norm_type)

    if "layer_norm_eps" in raw:
        _set_nested(result, "norm.epsilon", raw["layer_norm_eps"])

    # Tokenizer
    _set_nested(result, "tokenizer.type", "bpe")
    _set_nested(result, "tokenizer.library", "huggingface")

    # Position embedding type
    if model_type in ("gpt2", "gpt_neo"):
        _set_nested(result, "position_embedding.type", "learned")
    else:
        _set_nested(result, "position_embedding.type", "rope")

    # Embedding
    tie_embeddings = raw.get("tie_word_embeddings", False)
    _set_nested(result, "token_embedding.shared_embedding", tie_embeddings)

    # Bias flags
    attn_bias = raw.get("attention_bias", False)
    _set_nested(result, "attention.has_qkv_bias", attn_bias)
    _set_nested(result, "attention.has_output_bias", attn_bias)

    mlp_bias = raw.get("mlp_bias", False)
    if ffn_type == "mlp":
        _set_nested(result, "mlp.has_bias", mlp_bias)

    # Vision encoder (multimodal models)
    vision = parse_vision_config(raw)
    if vision:
        result["vision_encoder"] = vision

    return result


# Maps HuggingFace vision config field names to ModelPack vision encoder paths.
VISION_FIELD_MAP = {
    "hidden_size": "hidden_size",
    "patch_size": "patch_size",
    "image_size": "image_size",
    "num_hidden_layers": "num_layers",
    "num_attention_heads": "num_attention_heads",
    "intermediate_size": "intermediate_size",
    "num_channels": "in_channels",
    "in_chans": "in_channels",
}

# Known vision model types
VISION_MODEL_TYPES = {
    "llava", "llava_next", "llava_onevision",
    "mllama",  # LLaMA-3.2 Vision
    "qwen2_vl",
    "paligemma", "idefics2", "idefics3",
}

# Known vision encoder type mapping
VISION_ENCODER_TYPE_MAP = {
    "clip_vision_model": "clip_vit",
    "siglip_vision_model": "clip_vit",
    "CLIPVisionConfig": "clip_vit",
    "SiglipVisionConfig": "clip_vit",
}

# Known projector type mapping
PROJECTOR_TYPE_MAP = {
    "llava": ("mlp", 2, "gelu"),
    "llava_next": ("mlp", 2, "gelu"),
    "llava_onevision": ("mlp", 2, "gelu"),
    "mllama": ("cross_attention", None, None),
    "paligemma": ("linear", 1, None),
}

# Known fusion type mapping
FUSION_TYPE_MAP = {
    "llava": "late",
    "llava_next": "late",
    "llava_onevision": "late",
    "mllama": "cross_attention",
    "qwen2_vl": "early",
    "paligemma": "late",
    "idefics2": "late",
}


def parse_vision_config(raw: dict) -> dict | None:
    """Parse vision encoder fields from a HuggingFace multimodal config.

    Args:
        raw: The parsed config.json dict from HuggingFace.

    Returns:
        A dict following the ModelPack vision encoder spec, or None if not
        a vision model.
    """
    model_type = raw.get("model_type", "").lower()

    # Extract the nested vision_config dict
    vcfg = raw.get("vision_config")
    if vcfg is None and model_type not in VISION_MODEL_TYPES:
        return None

    # Some models embed vision config as a flat dict, others as nested
    if isinstance(vcfg, dict):
        vision_raw = vcfg
    else:
        return None

    result: dict = {}

    # Map static vision fields
    for hf_key, mp_path in VISION_FIELD_MAP.items():
        if hf_key in vision_raw and vision_raw[hf_key] is not None:
            _set_nested(result, mp_path, vision_raw[hf_key])

    # Qwen2-VL uses different field names
    if model_type == "qwen2_vl":
        if "depth" in vision_raw:
            result["num_layers"] = vision_raw["depth"]
        if "embed_dim" in vision_raw:
            result["hidden_size"] = vision_raw["embed_dim"]
        if "num_heads" in vision_raw:
            result["num_attention_heads"] = vision_raw["num_heads"]
        if "spatial_patch_size" in vision_raw:
            result["patch_size"] = vision_raw["spatial_patch_size"]
        if "temporal_patch_size" in vision_raw:
            result["temporal_patch_size"] = vision_raw["temporal_patch_size"]
        if "spatial_merge_size" in vision_raw:
            _set_nested(result, "dynamic_resolution.enabled", True)
            _set_nested(
                result, "dynamic_resolution.spatial_merge_size",
                vision_raw["spatial_merge_size"],
            )

    # Infer encoder type
    vision_model_type = vision_raw.get("model_type", "")
    encoder_type = VISION_ENCODER_TYPE_MAP.get(vision_model_type)
    if encoder_type is None:
        # Check for CLIP-like configs
        if "projection_dim" in vision_raw or vision_model_type in ("clip_vision_model",):
            encoder_type = "clip_vit"
        else:
            encoder_type = "vit"
    result["type"] = encoder_type

    # Activation
    hf_act = vision_raw.get("hidden_act")
    if hf_act:
        act = hf_act.lower()
        if "quick_gelu" in act:
            result["activation"] = "quick_gelu"
        elif "gelu" in act:
            result["activation"] = "gelu"
        elif "silu" in act or "swish" in act:
            result["activation"] = "silu"

    # Normalization
    eps = vision_raw.get("layer_norm_eps")
    if eps is not None:
        _set_nested(result, "norm.type", "layernorm")
        _set_nested(result, "norm.epsilon", eps)

    # Projector
    proj_info = PROJECTOR_TYPE_MAP.get(model_type)
    if proj_info:
        proj_type, proj_layers, proj_act = proj_info
        _set_nested(result, "projector.type", proj_type)
        if proj_layers is not None:
            _set_nested(result, "projector.num_layers", proj_layers)
        if proj_act is not None:
            _set_nested(result, "projector.activation", proj_act)
    # Check for explicit projector config
    if "projector_hidden_act" in raw:
        act = raw["projector_hidden_act"].lower()
        _set_nested(result, "projector.activation", act)

    # Special tokens
    token_fields = {
        "image_token_id": "special_tokens.image_token_id",
        "image_token_index": "special_tokens.image_token_id",
        "vision_start_token_id": "special_tokens.vision_start_token_id",
        "vision_end_token_id": "special_tokens.vision_end_token_id",
        "video_token_id": "special_tokens.video_token_id",
    }
    for hf_key, mp_path in token_fields.items():
        val = raw.get(hf_key)
        if val is not None:
            _set_nested(result, mp_path, val)

    # Fusion type
    fusion = FUSION_TYPE_MAP.get(model_type, NEEDS_REVIEW)
    result["fusion_type"] = fusion

    # Cross-attention layers (LLaMA-3.2 Vision)
    if model_type == "mllama":
        cross_attn_layers = raw.get("cross_attention_layers")
        if cross_attn_layers is not None:
            result["num_cross_attention_layers"] = len(cross_attn_layers)

    # Position embedding for vision encoder
    if model_type == "qwen2_vl":
        rope_scaling = raw.get("rope_scaling", {})
        if rope_scaling.get("type") == "mrope":
            _set_nested(result, "position_embedding.type", "mrope")
            sections = rope_scaling.get("mrope_section")
            if sections:
                _set_nested(result, "position_embedding.mrope_sections", sections)
    else:
        _set_nested(result, "position_embedding.type", "learned")

    return result


def format_yaml(spec: dict, indent: int = 0) -> str:
    """Format a spec dict as YAML string."""
    lines = []
    prefix = "  " * indent
    for key, value in spec.items():
        if isinstance(value, dict):
            lines.append(f"{prefix}{key}:")
            lines.append(format_yaml(value, indent + 1))
        elif isinstance(value, bool):
            lines.append(f"{prefix}{key}: {str(value).lower()}")
        elif isinstance(value, str):
            if value == NEEDS_REVIEW:
                lines.append(f"{prefix}{key}: {value}  # requires human review")
            else:
                lines.append(f'{prefix}{key}: "{value}"')
        elif value is None:
            lines.append(f"{prefix}{key}: null")
        else:
            lines.append(f"{prefix}{key}: {value}")
    return "\n".join(lines)


def load_config(source: str) -> dict:
    """Load a config.json from a file path or HuggingFace model ID.

    Args:
        source: Either a local file path or a HuggingFace model ID.

    Returns:
        The parsed config.json dict.
    """
    path = Path(source)
    if path.is_file():
        with path.open(encoding="utf-8") as f:
            return json.load(f)

    # Try loading from HuggingFace Hub
    try:
        from huggingface_hub import hf_hub_download

        config_path = hf_hub_download(repo_id=source, filename="config.json")
        with open(config_path, encoding="utf-8") as f:
            return json.load(f)
    except ImportError:
        print(
            "error: huggingface_hub not installed. "
            "Install with: pip install huggingface_hub",
            file=sys.stderr,
        )
        sys.exit(1)
    except Exception as e:
        print(f"error: failed to load config from '{source}': {e}", file=sys.stderr)
        sys.exit(1)


def main() -> int:
    parser = argparse.ArgumentParser(
        description="Parse HuggingFace config.json into ModelPack transformer spec",
    )
    parser.add_argument(
        "model",
        help="HuggingFace model ID (e.g., meta-llama/Meta-Llama-3-8B) "
        "or path to config.json",
    )
    parser.add_argument(
        "--format",
        choices=["yaml", "json"],
        default="yaml",
        help="Output format (default: yaml)",
    )

    args = parser.parse_args()

    raw = load_config(args.model)
    spec = parse_hf_config(raw)

    model_type = raw.get("model_type", "unknown")
    model_name = raw.get("_name_or_path", args.model)

    if args.format == "json":
        print(json.dumps(spec, indent=2))
    else:
        print(f"# ModelPack Transformer Spec")
        print(f"# Generated from: {model_name}")
        print(f"# Model type: {model_type}")
        print(f"# NOTE: Fields marked NEEDS_REVIEW require human verification")
        print()
        print(format_yaml(spec))

    # Report coverage
    needs_review = []
    _find_needs_review(spec, "", needs_review)
    if needs_review:
        print(f"\n# --- Fields requiring review ({len(needs_review)}) ---")
        for field in needs_review:
            print(f"#   - {field}")

    return 0


def _find_needs_review(d: dict, prefix: str, result: list) -> None:
    """Recursively find all NEEDS_REVIEW fields."""
    for key, value in d.items():
        path = f"{prefix}.{key}" if prefix else key
        if isinstance(value, dict):
            _find_needs_review(value, path, result)
        elif value == NEEDS_REVIEW:
            result.append(path)


if __name__ == "__main__":
    raise SystemExit(main())
