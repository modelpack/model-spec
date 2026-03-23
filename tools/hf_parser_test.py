#!/usr/bin/env python3
"""Tests for the HuggingFace config parser."""

from __future__ import annotations

import pytest

from hf_parser import NEEDS_REVIEW, parse_hf_config, parse_vision_config


# Minimal config.json samples based on real HuggingFace models.
MISTRAL_7B_CONFIG = {
    "model_type": "mistral",
    "vocab_size": 32000,
    "hidden_size": 4096,
    "num_hidden_layers": 32,
    "num_attention_heads": 32,
    "num_key_value_heads": 8,
    "intermediate_size": 14336,
    "max_position_embeddings": 32768,
    "rope_theta": 10000.0,
    "rms_norm_eps": 1e-5,
    "hidden_act": "silu",
    "sliding_window": 4096,
    "tie_word_embeddings": False,
    "attention_bias": False,
}

MIXTRAL_8X7B_CONFIG = {
    "model_type": "mixtral",
    "vocab_size": 32000,
    "hidden_size": 4096,
    "num_hidden_layers": 32,
    "num_attention_heads": 32,
    "num_key_value_heads": 8,
    "intermediate_size": 14336,
    "max_position_embeddings": 32768,
    "rope_theta": 1000000.0,
    "rms_norm_eps": 1e-5,
    "hidden_act": "silu",
    "num_local_experts": 8,
    "num_experts_per_tok": 2,
    "tie_word_embeddings": False,
}

QWEN2_7B_CONFIG = {
    "model_type": "qwen2",
    "vocab_size": 152064,
    "hidden_size": 3584,
    "num_hidden_layers": 28,
    "num_attention_heads": 28,
    "num_key_value_heads": 4,
    "intermediate_size": 18944,
    "max_position_embeddings": 131072,
    "rope_theta": 1000000.0,
    "rms_norm_eps": 1e-6,
    "hidden_act": "silu",
    "tie_word_embeddings": False,
    "attention_bias": True,
    "sliding_window": 131072,
}

GPT2_CONFIG = {
    "model_type": "gpt2",
    "vocab_size": 50257,
    "hidden_size": 768,
    "num_hidden_layers": 12,
    "num_attention_heads": 12,
    "intermediate_size": 3072,
    "max_position_embeddings": 1024,
    "layer_norm_eps": 1e-5,
    "activation_function": "gelu_new",
    "tie_word_embeddings": True,
}

DEEPSEEK_V2_LITE_CONFIG = {
    "model_type": "deepseek_v2",
    "vocab_size": 102400,
    "hidden_size": 2048,
    "num_hidden_layers": 27,
    "num_attention_heads": 16,
    "num_key_value_heads": 16,
    "intermediate_size": 10944,
    "max_position_embeddings": 163840,
    "rope_theta": 10000.0,
    "rms_norm_eps": 1e-6,
    "hidden_act": "silu",
    "kv_lora_rank": 512,
    "q_lora_rank": 1536,
    "qk_nope_head_dim": 128,
    "qk_rope_head_dim": 64,
    "v_head_dim": 128,
    "n_routed_experts": 64,
    "num_experts_per_tok": 6,
    "first_k_dense_replace": 1,
    "moe_layer_freq": 1,
    "num_shared_experts": 2,
    "routed_scaling_factor": 1.0,
    "topk_method": "group_limited_greedy",
    "norm_topk_prob": False,
    "tie_word_embeddings": False,
}

UNKNOWN_CONFIG = {
    "model_type": "some_new_model",
    "vocab_size": 65536,
    "hidden_size": 2048,
    "num_hidden_layers": 24,
    "num_attention_heads": 16,
}


class TestMistral:
    def test_basic_fields(self):
        spec = parse_hf_config(MISTRAL_7B_CONFIG)
        assert spec["vocabulary_size"] == 32000
        assert spec["hidden_size"] == 4096
        assert spec["num_layers"] == 32
        assert spec["type"] == "decoder"

    def test_attention(self):
        spec = parse_hf_config(MISTRAL_7B_CONFIG)
        attn = spec["attention"]
        assert attn["type"] == "gqa"
        assert attn["num_attention_heads"] == 32
        assert attn["num_key_value_heads"] == 8
        assert attn["head_dim"] == 128  # 4096 / 32
        assert attn["is_causal"] is True
        assert attn["sliding_window"] == 4096

    def test_ffn(self):
        spec = parse_hf_config(MISTRAL_7B_CONFIG)
        assert spec["ffn_type"] == "mlp"
        assert spec["mlp"]["intermediate_size"] == 14336
        assert spec["mlp"]["activation"] == "silu"
        assert spec["mlp"]["use_gated_activation"] is True

    def test_norm(self):
        spec = parse_hf_config(MISTRAL_7B_CONFIG)
        assert spec["norm"]["type"] == "rmsnorm"
        assert spec["norm"]["epsilon"] == 1e-5

    def test_position_embedding(self):
        spec = parse_hf_config(MISTRAL_7B_CONFIG)
        pe = spec["position_embedding"]
        assert pe["type"] == "rope"
        assert pe["rope_theta"] == 10000.0
        assert pe["max_position_embeddings"] == 32768


class TestMixtral:
    def test_moe_detection(self):
        spec = parse_hf_config(MIXTRAL_8X7B_CONFIG)
        assert spec["ffn_type"] == "moe"
        assert spec["moe"]["num_experts"] == 8
        assert spec["moe"]["top_k"] == 2

    def test_attention_still_gqa(self):
        spec = parse_hf_config(MIXTRAL_8X7B_CONFIG)
        assert spec["attention"]["type"] == "gqa"
        assert spec["attention"]["num_key_value_heads"] == 8


class TestQwen2:
    def test_attention_bias(self):
        spec = parse_hf_config(QWEN2_7B_CONFIG)
        assert spec["attention"]["has_qkv_bias"] is True
        assert spec["attention"]["has_output_bias"] is True

    def test_rope_theta(self):
        spec = parse_hf_config(QWEN2_7B_CONFIG)
        assert spec["position_embedding"]["rope_theta"] == 1000000.0


class TestGPT2:
    def test_mha_attention(self):
        spec = parse_hf_config(GPT2_CONFIG)
        assert spec["attention"]["type"] == "mha"

    def test_layernorm(self):
        spec = parse_hf_config(GPT2_CONFIG)
        assert spec["norm"]["type"] == "layernorm"
        assert spec["norm"]["epsilon"] == 1e-5

    def test_tied_embeddings(self):
        spec = parse_hf_config(GPT2_CONFIG)
        assert spec["token_embedding"]["shared_embedding"] is True

    def test_no_gated_activation(self):
        spec = parse_hf_config(GPT2_CONFIG)
        assert spec["mlp"]["use_gated_activation"] is False

    def test_gelu_activation(self):
        spec = parse_hf_config(GPT2_CONFIG)
        assert spec["mlp"]["activation"] == "gelu"


class TestDeepSeekV2:
    def test_mla_attention(self):
        spec = parse_hf_config(DEEPSEEK_V2_LITE_CONFIG)
        assert spec["attention"]["type"] == "mla"

    def test_mla_fields(self):
        spec = parse_hf_config(DEEPSEEK_V2_LITE_CONFIG)
        attn = spec["attention"]
        assert attn["kv_lora_rank"] == 512
        assert attn["q_lora_rank"] == 1536
        assert attn["qk_nope_head_dim"] == 128
        assert attn["qk_rope_head_dim"] == 64
        assert attn["v_head_dim"] == 128

    def test_moe_with_n_routed_experts(self):
        spec = parse_hf_config(DEEPSEEK_V2_LITE_CONFIG)
        assert spec["ffn_type"] == "moe"
        assert spec["moe"]["num_experts"] == 64
        assert spec["moe"]["top_k"] == 6

    def test_shared_experts(self):
        spec = parse_hf_config(DEEPSEEK_V2_LITE_CONFIG)
        assert spec["moe"]["num_shared_experts"] == 2

    def test_moe_routing_fields(self):
        spec = parse_hf_config(DEEPSEEK_V2_LITE_CONFIG)
        assert spec["moe"]["routed_scaling_factor"] == 1.0
        assert spec["moe"]["topk_method"] == "group_limited_greedy"
        assert spec["moe"]["norm_topk_prob"] is False

    def test_mixed_layers(self):
        spec = parse_hf_config(DEEPSEEK_V2_LITE_CONFIG)
        assert spec["layer_structure"] == "mixed"
        assert spec["mixed_layers"]["first_k_dense_replace"] == 1
        assert spec["mixed_layers"]["moe_layer_freq"] == 1


class TestGPT2PositionEmbedding:
    def test_learned_position_embedding(self):
        spec = parse_hf_config(GPT2_CONFIG)
        assert spec["position_embedding"]["type"] == "learned"


class TestUnknownModel:
    def test_needs_review_flags(self):
        spec = parse_hf_config(UNKNOWN_CONFIG)
        assert spec["attention"]["type"] == NEEDS_REVIEW
        assert spec["ffn_type"] == NEEDS_REVIEW

    def test_static_fields_still_parsed(self):
        spec = parse_hf_config(UNKNOWN_CONFIG)
        assert spec["vocabulary_size"] == 65536
        assert spec["hidden_size"] == 2048
        assert spec["num_layers"] == 24

    def test_head_dim_derived(self):
        spec = parse_hf_config(UNKNOWN_CONFIG)
        assert spec["attention"]["head_dim"] == 128  # 2048 / 16


class TestHeadDimDerivation:
    def test_explicit_head_dim(self):
        config = {
            "model_type": "mistral",
            "vocab_size": 32000,
            "hidden_size": 4096,
            "num_attention_heads": 32,
            "head_dim": 64,  # explicit, not derived
        }
        spec = parse_hf_config(config)
        assert spec["attention"]["head_dim"] == 64  # uses explicit value

    def test_derived_head_dim(self):
        config = {
            "model_type": "llama",
            "vocab_size": 32000,
            "hidden_size": 4096,
            "num_attention_heads": 32,
        }
        spec = parse_hf_config(config)
        assert spec["attention"]["head_dim"] == 128  # 4096 / 32


# ============================================================================
# Vision model configs
# ============================================================================

LLAVA_15_CONFIG = {
    "model_type": "llava",
    "vocab_size": 32064,
    "hidden_size": 4096,
    "num_hidden_layers": 32,
    "num_attention_heads": 32,
    "num_key_value_heads": 32,
    "intermediate_size": 11008,
    "max_position_embeddings": 4096,
    "rms_norm_eps": 1e-5,
    "hidden_act": "silu",
    "image_token_index": 32000,
    "projector_hidden_act": "gelu",
    "vision_config": {
        "model_type": "clip_vision_model",
        "hidden_size": 1024,
        "patch_size": 14,
        "image_size": 336,
        "num_hidden_layers": 24,
        "num_attention_heads": 16,
        "intermediate_size": 4096,
        "num_channels": 3,
        "hidden_act": "quick_gelu",
        "layer_norm_eps": 1e-5,
        "projection_dim": 768,
    },
}

QWEN2_VL_CONFIG = {
    "model_type": "qwen2_vl",
    "vocab_size": 152064,
    "hidden_size": 3584,
    "num_hidden_layers": 28,
    "num_attention_heads": 28,
    "num_key_value_heads": 4,
    "intermediate_size": 18944,
    "max_position_embeddings": 32768,
    "rms_norm_eps": 1e-6,
    "hidden_act": "silu",
    "vision_start_token_id": 151652,
    "vision_end_token_id": 151653,
    "vision_token_id": 151654,
    "image_token_id": 151655,
    "video_token_id": 151656,
    "rope_scaling": {
        "type": "mrope",
        "mrope_section": [16, 24, 24],
    },
    "vision_config": {
        "depth": 32,
        "embed_dim": 1280,
        "num_heads": 16,
        "in_chans": 3,
        "spatial_patch_size": 14,
        "spatial_merge_size": 2,
        "temporal_patch_size": 2,
        "hidden_act": "quick_gelu",
    },
}

LLAMA_32_VISION_CONFIG = {
    "model_type": "mllama",
    "vocab_size": 128256,
    "hidden_size": 4096,
    "num_hidden_layers": 32,
    "num_attention_heads": 32,
    "num_key_value_heads": 8,
    "intermediate_size": 14336,
    "max_position_embeddings": 131072,
    "rms_norm_eps": 1e-5,
    "hidden_act": "silu",
    "image_token_index": 128256,
    "cross_attention_layers": [3, 8, 13, 18, 23, 28, 33, 38],
    "vision_config": {
        "model_type": "clip_vision_model",
        "hidden_size": 1280,
        "patch_size": 14,
        "image_size": 560,
        "num_hidden_layers": 32,
        "num_attention_heads": 16,
        "intermediate_size": 5120,
        "num_channels": 3,
        "hidden_act": "gelu",
        "layer_norm_eps": 1e-5,
    },
}

TEXT_ONLY_CONFIG = {
    "model_type": "llama",
    "vocab_size": 32000,
    "hidden_size": 4096,
    "num_hidden_layers": 32,
    "num_attention_heads": 32,
}


# ============================================================================
# Vision model tests
# ============================================================================

class TestLLaVA:
    def test_vision_encoder_present(self):
        spec = parse_hf_config(LLAVA_15_CONFIG)
        assert "vision_encoder" in spec

    def test_encoder_type(self):
        spec = parse_hf_config(LLAVA_15_CONFIG)
        ve = spec["vision_encoder"]
        assert ve["type"] == "clip_vit"

    def test_encoder_fields(self):
        spec = parse_hf_config(LLAVA_15_CONFIG)
        ve = spec["vision_encoder"]
        assert ve["hidden_size"] == 1024
        assert ve["patch_size"] == 14
        assert ve["image_size"] == 336
        assert ve["num_layers"] == 24
        assert ve["num_attention_heads"] == 16
        assert ve["intermediate_size"] == 4096
        assert ve["in_channels"] == 3

    def test_activation(self):
        spec = parse_hf_config(LLAVA_15_CONFIG)
        assert spec["vision_encoder"]["activation"] == "quick_gelu"

    def test_norm(self):
        spec = parse_hf_config(LLAVA_15_CONFIG)
        ve = spec["vision_encoder"]
        assert ve["norm"]["type"] == "layernorm"
        assert ve["norm"]["epsilon"] == 1e-5

    def test_projector(self):
        spec = parse_hf_config(LLAVA_15_CONFIG)
        proj = spec["vision_encoder"]["projector"]
        assert proj["type"] == "mlp"
        assert proj["num_layers"] == 2
        assert proj["activation"] == "gelu"

    def test_special_tokens(self):
        spec = parse_hf_config(LLAVA_15_CONFIG)
        assert spec["vision_encoder"]["special_tokens"]["image_token_id"] == 32000

    def test_fusion_type(self):
        spec = parse_hf_config(LLAVA_15_CONFIG)
        assert spec["vision_encoder"]["fusion_type"] == "late"

    def test_position_embedding(self):
        spec = parse_hf_config(LLAVA_15_CONFIG)
        assert spec["vision_encoder"]["position_embedding"]["type"] == "learned"


class TestQwen2VL:
    def test_vision_encoder_present(self):
        spec = parse_hf_config(QWEN2_VL_CONFIG)
        assert "vision_encoder" in spec

    def test_qwen_specific_fields(self):
        spec = parse_hf_config(QWEN2_VL_CONFIG)
        ve = spec["vision_encoder"]
        assert ve["num_layers"] == 32
        assert ve["hidden_size"] == 1280
        assert ve["num_attention_heads"] == 16
        assert ve["patch_size"] == 14

    def test_temporal_patch(self):
        spec = parse_hf_config(QWEN2_VL_CONFIG)
        assert spec["vision_encoder"]["temporal_patch_size"] == 2

    def test_dynamic_resolution(self):
        spec = parse_hf_config(QWEN2_VL_CONFIG)
        dr = spec["vision_encoder"]["dynamic_resolution"]
        assert dr["enabled"] is True
        assert dr["spatial_merge_size"] == 2

    def test_special_tokens(self):
        spec = parse_hf_config(QWEN2_VL_CONFIG)
        tokens = spec["vision_encoder"]["special_tokens"]
        assert tokens["image_token_id"] == 151655
        assert tokens["vision_start_token_id"] == 151652
        assert tokens["vision_end_token_id"] == 151653
        assert tokens["video_token_id"] == 151656

    def test_mrope_position_embedding(self):
        spec = parse_hf_config(QWEN2_VL_CONFIG)
        pe = spec["vision_encoder"]["position_embedding"]
        assert pe["type"] == "mrope"
        assert pe["mrope_sections"] == [16, 24, 24]

    def test_fusion_type(self):
        spec = parse_hf_config(QWEN2_VL_CONFIG)
        assert spec["vision_encoder"]["fusion_type"] == "early"


class TestLLaMA32Vision:
    def test_vision_encoder_present(self):
        spec = parse_hf_config(LLAMA_32_VISION_CONFIG)
        assert "vision_encoder" in spec

    def test_encoder_fields(self):
        spec = parse_hf_config(LLAMA_32_VISION_CONFIG)
        ve = spec["vision_encoder"]
        assert ve["hidden_size"] == 1280
        assert ve["patch_size"] == 14
        assert ve["image_size"] == 560
        assert ve["num_layers"] == 32
        assert ve["num_attention_heads"] == 16

    def test_cross_attention_projector(self):
        spec = parse_hf_config(LLAMA_32_VISION_CONFIG)
        proj = spec["vision_encoder"]["projector"]
        assert proj["type"] == "cross_attention"

    def test_cross_attention_layers_count(self):
        spec = parse_hf_config(LLAMA_32_VISION_CONFIG)
        assert spec["vision_encoder"]["num_cross_attention_layers"] == 8

    def test_fusion_type(self):
        spec = parse_hf_config(LLAMA_32_VISION_CONFIG)
        assert spec["vision_encoder"]["fusion_type"] == "cross_attention"

    def test_special_tokens(self):
        spec = parse_hf_config(LLAMA_32_VISION_CONFIG)
        assert spec["vision_encoder"]["special_tokens"]["image_token_id"] == 128256


class TestTextOnlyModel:
    def test_no_vision_encoder(self):
        spec = parse_hf_config(TEXT_ONLY_CONFIG)
        assert "vision_encoder" not in spec

    def test_no_vision_config_returns_none(self):
        result = parse_vision_config(TEXT_ONLY_CONFIG)
        assert result is None
