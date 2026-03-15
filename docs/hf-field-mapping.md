# HuggingFace to ModelPack Architecture Field Mapping

This document maps HuggingFace `config.json` field names to the ModelPack architecture
specification fields defined in [architecture.md](architecture.md) (PR #111).

## Purpose

Different model families use different field names and conventions in their HuggingFace
`config.json` files. This mapping is needed to:

1. Validate that the ModelPack architecture vocabulary covers real-world models
2. Guide the development of auto-generation tooling (issue #164)
3. Document where fields are directly mappable vs. need derivation or inference

## Models Analyzed

| Model | HuggingFace ID | Architecture Highlights |
|---|---|---|
| Llama 3.1 8B | `meta-llama/Llama-3.1-8B` | GQA, dense MLP, RoPE with Llama3 scaling |
| Mistral 7B v0.3 | `mistralai/Mistral-7B-v0.3` | GQA, no sliding window (disabled), RoPE |
| Mixtral 8x7B | `mistralai/Mixtral-8x7B-v0.1` | GQA, MoE (8 experts, top-2), RoPE |
| DeepSeek-V2-Lite | `deepseek-ai/DeepSeek-V2-Lite` | MLA, MoE with shared experts (64 routed + 2 shared, top-6), mixed layers, YaRN RoPE |

## Mapping Type Legend

| Type | Meaning |
|---|---|
| **direct** | HF field has the same name and semantics as the ModelPack field |
| **renamed** | Different name but same meaning; trivial 1:1 substitution |
| **derived** | Computed from other HF fields; formula is given in Notes |
| **inferred** | Not in HF config at all; must be determined from model family knowledge or code |
| **model-specific** | Only present in certain HF model families |
| **not-in-hf** | ModelPack spec-only field with no HF equivalent |

---

## Field Mapping Table

### 1. Top-Level Fields

| ModelPack Field | Llama 3.1 8B | Mistral 7B v0.3 | Mixtral 8x7B | DeepSeek-V2-Lite | Mapping Type | Notes |
|---|---|---|---|---|---|---|
| `architecture_version` | — | — | — | — | not-in-hf | ModelPack schema version; no HF equivalent |
| `type` | `architectures: ["LlamaForCausalLM"]` | `architectures: ["MistralForCausalLM"]` | `architectures: ["MixtralForCausalLM"]` | `architectures: ["DeepseekV2ForCausalLM"]` | derived | Extract: if class ends in `ForCausalLM` → `"decoder"` |
| `vocabulary_size` | `vocab_size: 128256` | `vocab_size: 32768` | `vocab_size: 32000` | `vocab_size: 102400` | renamed | `vocab_size` → `vocabulary_size` |
| `hidden_size` | `hidden_size: 4096` | `hidden_size: 4096` | `hidden_size: 4096` | `hidden_size: 2048` | direct | Identical field name and meaning |

---

### 2. Tokenizer

| ModelPack Field | Llama 3.1 8B | Mistral 7B v0.3 | Mixtral 8x7B | DeepSeek-V2-Lite | Mapping Type | Notes |
|---|---|---|---|---|---|---|
| `type` | — | — | — | — | inferred | Not in config.json; determined from `tokenizer_config.json` or model family knowledge. All four use BPE → `"bpe"` |
| `library` | — | — | — | — | inferred | Always `"huggingface"` for HF-hosted models; not recorded in config.json |
| `revision` | — | — | — | — | not-in-hf | ModelPack uses `"main"` as default; no HF equivalent in config.json |

---

### 3. Token Embedding

| ModelPack Field | Llama 3.1 8B | Mistral 7B v0.3 | Mixtral 8x7B | DeepSeek-V2-Lite | Mapping Type | Notes |
|---|---|---|---|---|---|---|
| `has_bias` | — | — | — | — | inferred | No HF field for embedding bias. Standard transformer embeddings have no bias → `false` for all four |
| `has_norm` | — | — | — | — | inferred | No HF field. None of the four apply a norm after the embedding → `false` |
| `shared_embedding` | `tie_word_embeddings: false` | `tie_word_embeddings: false` | `tie_word_embeddings: false` | `tie_word_embeddings: false` | renamed | `tie_word_embeddings` → `shared_embedding`; value is the same boolean |

---

### 4. Position Embedding

| ModelPack Field | Llama 3.1 8B | Mistral 7B v0.3 | Mixtral 8x7B | DeepSeek-V2-Lite | Mapping Type | Notes |
|---|---|---|---|---|---|---|
| `type` | — | — | — | — | inferred | No explicit field; presence of `rope_theta` strongly implies `"rope"`. Must be confirmed from model class |
| `max_position_embeddings` | `max_position_embeddings: 131072` | `max_position_embeddings: 32768` | `max_position_embeddings: 32768` | `max_position_embeddings: 163840` | direct | Identical field name. Llama's effective length is extended by `rope_scaling.factor`; DeepSeek uses YaRN scaling |
| `rope_theta` | `rope_theta: 500000.0` | `rope_theta: 1000000.0` | `rope_theta: 1000000.0` | `rope_theta: 10000` | direct | Identical field name |
| `rope_scaling` | `rope_scaling: {factor: 8.0, low_freq_factor: 1.0, high_freq_factor: 4.0, original_max_position_embeddings: 8192, rope_type: "llama3"}` | — (absent) | — (absent) | `rope_scaling: {type: "yarn", factor: 40, beta_fast: 32, beta_slow: 1, mscale: 0.707, mscale_all_dim: 0.707, original_max_position_embeddings: 4096}` | model-specific | Llama and DeepSeek both use this field but with different schemas. Llama uses `rope_type: "llama3"`; DeepSeek uses `type: "yarn"` (YaRN). The nested schema is not standardized across model families |

---

### 5. Attention

| ModelPack Field | Llama 3.1 8B | Mistral 7B v0.3 | Mixtral 8x7B | DeepSeek-V2-Lite | Mapping Type | Notes |
|---|---|---|---|---|---|---|
| `type` | `num_attention_heads: 32` vs `num_key_value_heads: 8` | same | same | `kv_lora_rank: 512` present | derived | For Llama/Mistral/Mixtral: `num_key_value_heads < num_attention_heads` → `"gqa"`. For DeepSeek: head counts are equal (16 == 16), which would incorrectly imply MHA; the correct signal is presence of `kv_lora_rank` → `"mla"`. The derivation rule must be updated to check for `kv_lora_rank` before applying the head-count comparison |
| `is_causal` | `architectures: ["LlamaForCausalLM"]` | `architectures: ["MistralForCausalLM"]` | `architectures: ["MixtralForCausalLM"]` | `architectures: ["DeepseekV2ForCausalLM"]` | inferred | Presence of `ForCausalLM` suffix → `true`. Not an explicit field |
| `num_attention_heads` | `num_attention_heads: 32` | `num_attention_heads: 32` | `num_attention_heads: 32` | `num_attention_heads: 16` | direct | Identical field name |
| `num_key_value_heads` | `num_key_value_heads: 8` | `num_key_value_heads: 8` | `num_key_value_heads: 8` | `num_key_value_heads: 16` | direct | Identical field name. For DeepSeek (MLA) this equals `num_attention_heads`; traditional KV-head semantics do not apply — MLA compresses K and V via a shared low-rank latent vector (`kv_lora_rank`) |
| `head_dim` | — (absent) | — (absent) | — (absent) | — (absent; split into three fields) | derived / model-specific | Llama/Mistral/Mixtral: derived as `hidden_size / num_attention_heads` = 4096 / 32 = **128**. DeepSeek (MLA): no single head dimension; Q head is split into `qk_nope_head_dim: 128` (non-positional) + `qk_rope_head_dim: 64` (RoPE) = 192 total; V head is `v_head_dim: 128`. The ModelPack single `head_dim` field cannot represent MLA — **spec gap** |
| `is_qkv_merged` | — | — | — | — | inferred | No HF field. Separate Q/K/V projection weights are standard for all four → `false` |
| `has_qkv_bias` | `attention_bias: false` | — (absent) | — (absent) | `attention_bias: false` | renamed / inferred | Llama and DeepSeek: `attention_bias` covers Q, K, V, and O projections. Mistral/Mixtral omit this field entirely → must infer `false` |
| `has_output_bias` | `attention_bias: false` | — (absent) | — (absent) | `attention_bias: false` | renamed / inferred | Same field as `has_qkv_bias` in Llama/DeepSeek (one flag covers all projections). No dedicated HF field |
| `has_pre_norm` | — | — | — | — | inferred | Pre-norm is architectural convention for all four (RMSNorm before each sub-layer) → `true`. Not explicit in config |
| `has_post_norm` | — | — | — | — | inferred | None of the four apply a second norm after the residual → `false` |
| `has_residual` | — | — | — | — | inferred | Residual connections are always present in standard transformers → `true`. Not an explicit HF field |

---

### 6. Feed-Forward (MLP)

| ModelPack Field | Llama 3.1 8B | Mistral 7B v0.3 | Mixtral 8x7B | DeepSeek-V2-Lite | Mapping Type | Notes |
|---|---|---|---|---|---|---|
| `intermediate_size` | `intermediate_size: 14336` | `intermediate_size: 14336` | `intermediate_size: 14336` | `intermediate_size: 10944` | direct | Identical field name. For DeepSeek this applies only to dense MLP layers; MoE expert width is separate (see MoE section) |
| `activation` | `hidden_act: "silu"` | `hidden_act: "silu"` | `hidden_act: "silu"` | `hidden_act: "silu"` | renamed | `hidden_act` → `activation`; value maps directly |
| `use_gated_activation` | — | — | — | — | inferred | No explicit HF field. `"silu"` activation in all four architectures implies SwiGLU (gated) → `true`. Must be derived from model class knowledge |
| `is_mlp_merged` | — | — | — | — | inferred | No HF field. Up and gate projections are separate weight matrices in all four → `false` |
| `has_bias` | `mlp_bias: false` | — (absent) | — (absent) | — (absent) | renamed / inferred | Llama: `mlp_bias` covers up, gate, and down projections. Mistral/Mixtral/DeepSeek omit this field → infer `false` |
| `has_residual` | — | — | — | — | inferred | Residual connections always present → `true` |
| `has_pre_norm` | — | — | — | — | inferred | Same pre-norm applied to MLP as to attention in all four → `true` |
| `has_post_norm` | — | — | — | — | inferred | Not used → `false` |

---

### 7. Mixture of Experts (MoE) — Mixtral 8x7B and DeepSeek-V2-Lite

| ModelPack Field | Mixtral 8x7B | DeepSeek-V2-Lite | Mapping Type | Notes |
|---|---|---|---|---|
| `num_experts` | `num_local_experts: 8` | `n_routed_experts: 64` | renamed | Different source field names across families: `num_local_experts` (Mixtral) vs `n_routed_experts` (DeepSeek). Auto-generation tooling needs a per-`model_type` lookup |
| `top_k` | `num_experts_per_tok: 2` | `num_experts_per_tok: 6` | renamed | `num_experts_per_tok` → `top_k`; consistent field name across both families |
| `moe_intermediate_size` | `intermediate_size: 14336` | `moe_intermediate_size: 1408` | renamed / direct | Mixtral: renamed from `intermediate_size` (shared field, no distinction). DeepSeek: dedicated `moe_intermediate_size` field distinct from `intermediate_size` (10944 for dense layers). DeepSeek's config is cleaner |
| `num_shared_experts` | — (absent) | `n_shared_experts: 2` | model-specific / renamed | Shared (always-active) experts are absent in Mixtral. DeepSeek: `n_shared_experts` → `num_shared_experts` |
| `shared_expert_intermediate_size` | — (absent) | — (absent; derived) | model-specific / derived | Not in either HF config. DeepSeek shared experts use the same width as routed experts: derived as `moe_intermediate_size: 1408` |
| `scoring_function` | — (absent) | `scoring_func: "softmax"` | inferred / renamed | Mixtral: not explicit → inferred as `"softmax"`. DeepSeek: `scoring_func` → `scoring_function` |
| `norm_topk_prob` | — (absent) | `norm_topk_prob: false` | inferred / direct | Mixtral: not present → inferred `false`. DeepSeek: explicit `norm_topk_prob` field; same name as ModelPack spec |
| `activation` | `hidden_act: "silu"` | `hidden_act: "silu"` | renamed | `hidden_act` → `activation`; consistent across both |
| `use_gated_activation` | — (absent) | — (absent) | inferred | SwiGLU implied by `"silu"` activation in both → `true` |
| `has_bias` | — (absent) | — (absent) | inferred | No bias in MoE expert projections → `false` for both |

---

### 8. Normalization

| ModelPack Field | Llama 3.1 8B | Mistral 7B v0.3 | Mixtral 8x7B | DeepSeek-V2-Lite | Mapping Type | Notes |
|---|---|---|---|---|---|---|
| `type` | `rms_norm_eps: 1e-05` | `rms_norm_eps: 1e-05` | `rms_norm_eps: 1e-05` | `rms_norm_eps: 1e-06` | derived | No explicit `norm_type` field. Presence of `rms_norm_eps` (vs `layer_norm_eps`) implies `"rmsnorm"` for all four |
| `epsilon` | `rms_norm_eps: 1e-05` | `rms_norm_eps: 1e-05` | `rms_norm_eps: 1e-05` | `rms_norm_eps: 1e-06` | renamed | `rms_norm_eps` → `epsilon`. DeepSeek uses a tighter value (1e-06) than the other three (1e-05) |

---

### 9. Layer Configuration

| ModelPack Field | Llama 3.1 8B | Mistral 7B v0.3 | Mixtral 8x7B | DeepSeek-V2-Lite | Mapping Type | Notes |
|---|---|---|---|---|---|---|
| `uniform_layers` vs `mixed_layers` | `uniform_layers` | `uniform_layers` | `uniform_layers` | `mixed_layers` | inferred / derived | Llama/Mistral/Mixtral: all layers identical → `uniform_layers` (inferred from absence of per-layer variation). DeepSeek: `first_k_dense_replace: 1` + `moe_layer_freq: 1` → layer 0 is dense MLP, layers 1–26 are MoE → `mixed_layers` (derived from both fields together) |
| `num_layers` | `num_hidden_layers: 32` | `num_hidden_layers: 32` | `num_hidden_layers: 32` | `num_hidden_layers: 27` | renamed | `num_hidden_layers` → `num_layers` |
| `mlp_layers[]` _(mixed_layers only)_ | — | — | — | derived from `first_k_dense_replace: 1` → `[0]` | derived | Not a field in Llama/Mistral/Mixtral (they use uniform layers). For DeepSeek: all layer indices `< first_k_dense_replace` are dense MLP. With value 1, only layer index 0 is dense → `mlp_layers: [0]` |
| `moe_frequency` _(mixed_layers only)_ | — | — | — | `moe_layer_freq: 1` | renamed | Not applicable to Llama/Mistral/Mixtral. DeepSeek: `moe_layer_freq` → `moe_frequency`; value 1 means every layer after the initial dense block is MoE |

---

## HuggingFace Fields Not Covered by the ModelPack Architecture Spec

These HF config.json fields have no equivalent in `architecture.md` and may represent gaps worth noting for issue #164:

| HF Field | Models Present | Meaning | Comment |
|---|---|---|---|
| `attention_dropout` | All four (= 0.0) | Dropout rate applied to attention weights | Training hyperparameter; arguably out of scope for an inference-focused spec |
| `sliding_window` | Mistral, Mixtral (= null here) | Sliding window attention size | Present but disabled for these versions; Mistral 7B v0.1 used 4096. The spec has no field for this attention pattern |
| `bos_token_id` / `eos_token_id` / `pad_token_id` | All four | Special token IDs | Tokenizer metadata; the spec records tokenizer `type` and `library` but not specific token IDs |
| `initializer_range` | All four | Std dev for weight initialization | Training-only; not needed for inference |
| `pretraining_tp` | Llama, DeepSeek | Tensor parallelism degree during pretraining | Training artifact; not relevant to the architecture spec |
| `output_router_logits` | Mixtral only | Whether to output MoE router logits | Training flag; no inference-time relevance |
| `router_aux_loss_coef` | Mixtral only | Auxiliary loss coefficient for router balancing | Training hyperparameter |
| `use_cache` | All four (= true) | Whether to use KV cache | Runtime toggle, not an architectural property |
| `torch_dtype` | All four | Default weight dtype | Partially overlaps with ModelPack `config.precision`, but at a different granularity |
| `transformers_version` | All four | Library version used | Provenance metadata; not architectural |
| `kv_lora_rank` | DeepSeek only | Rank of the compressed KV latent in MLA | Core MLA parameter with no ModelPack field. The spec defines `attention.type: "mla"` but does not capture the compression rank — **spec gap for MLA** |
| `q_lora_rank` | DeepSeek only | Rank of Q compression in MLA (null = uncompressed) | Optional Q-side compression in MLA; no ModelPack field |
| `qk_nope_head_dim` | DeepSeek only | Non-positional (non-RoPE) component of Q/K head dimension | MLA splits head_dim into RoPE and non-RoPE parts; the single ModelPack `head_dim` field cannot represent this |
| `qk_rope_head_dim` | DeepSeek only | RoPE component of Q/K head dimension | Same issue as above |
| `v_head_dim` | DeepSeek only | V projection head dimension (may differ from Q/K) | In MLA, V head dim can differ from Q/K head dim; no dedicated ModelPack field |
| `n_group` / `topk_group` | DeepSeek only | Expert grouping and group-level top-k routing | DeepSeek's group-constrained routing strategy; not captured in the spec |
| `topk_method` | DeepSeek only | Method for selecting top-k experts (e.g., `"greedy"`) | Distinct from `scoring_function` (which scores experts); the selection method is also not covered |
| `routed_scaling_factor` | DeepSeek only | Scaling factor applied to routed expert outputs | No ModelPack equivalent |
| `seq_aux` | DeepSeek only | Whether auxiliary load-balancing loss is computed at sequence level | Training hyperparameter |
| `aux_loss_alpha` | DeepSeek only | Auxiliary loss weight for MoE load balancing | Training hyperparameter |
| `auto_map` | DeepSeek only | Custom class mappings for AutoModel loading | HF-specific loading mechanism; not architectural |

---

## Key Observations

### Parameter Consistency and Divergence

- The three dense/standard models (Llama, Mistral, Mixtral) are highly consistent, sharing identical values for:
  - `hidden_size` (4096), `num_attention_heads` (32), `num_key_value_heads` (8), `num_hidden_layers` (32), `intermediate_size` (14336)
  - `rms_norm_eps` (1e-05)
- DeepSeek-V2-Lite differs from this group on all of the above parameters (e.g., `hidden_size: 2048`, `num_hidden_layers: 27`, `rms_norm_eps: 1e-06`).
- However, some parameters are consistent across all four models:
  - `hidden_act: "silu"`
  - `tie_word_embeddings: false`

### Naming inconsistencies between model families

- **`attention_bias` and `mlp_bias`**: Present explicitly in Llama, absent in Mistral/Mixtral. Auto-generation tooling must handle the missing-field case by defaulting to `false`.
- **`num_local_experts` vs `n_routed_experts`**: Mixtral and DeepSeek use different field names for the number of routed experts. Auto-generation tooling needs a per-`model_type` normalization layer.
- **`intermediate_size` dual role**: In dense models it is the MLP width; in Mixtral it is the per-expert MoE width. DeepSeek correctly separates these into `intermediate_size` (dense) and `moe_intermediate_size` (MoE). A converter must know which HF field to use based on whether the model is dense or MoE.
- **`scoring_func` vs `scoring_function`**: DeepSeek uses `scoring_func`; the ModelPack spec uses `scoring_function`. Mixtral lacks the field entirely.
- **`rope_scaling` schema inconsistency**: Llama uses `rope_type: "llama3"` inside `rope_scaling`; DeepSeek uses `type: "yarn"`. The key name for the scaling algorithm type differs between families.

### MLA field mapping complexity

DeepSeek's Multi-Head Latent Attention (MLA) is fundamentally incompatible with the single-`head_dim` abstraction in the current spec:

- The Q head has two components: `qk_nope_head_dim: 128` (content) + `qk_rope_head_dim: 64` (positional) = 192 total.
- The V head has a separate dimension: `v_head_dim: 128`.
- The KV cache is compressed via a low-rank latent of size `kv_lora_rank: 512`, which is the key MLA innovation.

None of `kv_lora_rank`, `q_lora_rank`, `qk_nope_head_dim`, `qk_rope_head_dim`, or `v_head_dim` have ModelPack equivalents. The spec would need MLA-specific sub-fields under the attention object to fully describe this architecture.

Additionally, the existing GQA detection heuristic (`num_key_value_heads < num_attention_heads` → `"gqa"`) silently misclassifies DeepSeek as MHA (since both head counts are 16). The correct detection is: if `kv_lora_rank` is present → `"mla"`. This check must take precedence over the head-count comparison.

### DeepSeek is the first model to exercise three spec features

- **`mixed_layers`**: DeepSeek-V2-Lite is the first model in this analysis to use a hybrid layer topology (1 dense MLP + 26 MoE). The `first_k_dense_replace` and `moe_layer_freq` fields together encode the pattern.
- **`num_shared_experts`**: The shared (always-active) expert mechanism is a DeepSeek-V2/V3 innovation. `n_shared_experts: 2` is the first real-world mapping for this ModelPack field.
- **`attention.type: "mla"`**: DeepSeek is the only model in this analysis that uses MLA, validating that the spec's three-way attention type enum (`mha`/`gqa`/`mla`) covers real-world diversity — while also revealing that the spec lacks MLA's internal parameters.

### ModelPack fields that cannot be determined from `config.json` alone

These require either architectural knowledge (hardcoded per `model_type`) or reading source code:

- `tokenizer.type` — requires reading `tokenizer_config.json` (`"tokenizer_class"` field)
- `position_embedding.type` — implied by `rope_theta` presence, but not explicit
- `attention.type` (mha/gqa/mla) — partially derivable: check for `kv_lora_rank` first (→ `"mla"`), then compare head counts (→ `"gqa"` or `"mha"`). The two-step rule must be coded explicitly
- `attention.is_causal` — inferred from the model class name suffix (`ForCausalLM`)
- `attention.is_qkv_merged`, `mlp.is_mlp_merged` — requires reading model implementation
- `mlp.use_gated_activation` — requires knowing that `silu` in these architectures implies SwiGLU
- `normalization.type` — inferred from which eps field is present (`rms_norm_eps` vs `layer_norm_eps`)
- All `has_pre_norm`, `has_post_norm`, `has_residual` flags — architectural conventions not encoded in config.json
- `uniform_layers` vs `mixed_layers` — for Llama/Mistral/Mixtral, inferred from absence of per-layer variation; for DeepSeek, derived from `first_k_dense_replace` and `moe_layer_freq`
