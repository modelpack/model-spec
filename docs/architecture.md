# Model Architecture Configuration

Each model artifact has an associated optional architecture configuration that describes the detailed structure and components of the model. Currently, only decoder-type transformer architectures are supported. Future extensions will include:

- Multi-modal language models
- State Space Models
- Diffusion Models

## Terminology

The transformer is the most popular architecture for LLMs. It consists of a stack of structured layers, where each layer contains a self-attention block and a feed-forward network, with normalization layers and residual connections. The complete architecture includes a tokenizer, input embedding layer, position embedding layer, transformer layers, and output embedding layer. The transformer architecture has remained relatively stable since [Attention is all you need][attention-paper]. As shown in the table below, current open-weight model architectures are converging, making it feasible to define a common abstraction.

| Model                           | Tokenizer | PE         | Self-Attention | Norm       | Feed-Forward | Residual |
|---------------------------------|-----------|------------|----------------|------------|--------------|----------|
| [GPT2][gpt2-repo]               | BPE       | Sinusoidal | MHA            | Layer Norm | MLP          | Yes      |
| [Llama2][llama2-paper]          | BPE       | RoPE       | GQA            | RMS Norm   | MLP          | Yes      |
| [Llama3][llama3-paper]          | BPE       | RoPE       | GQA            | RMS Norm   | MLP          | Yes      |
| [Qwen2][qwen2-paper]            | BPE       | RoPE       | GQA            | RMS Norm   | MoE          | Yes      |
| [Qwen3][qwen3-paper]            | BPE       | RoPE       | GQA            | RMS Norm   | MoE          | Yes      |
| [Gemma2][gemma2-paper]          | BPE       | RoPE       | GQA            | RMS Norm   | MLP          | Yes      |
| [Gemma3][gemma3-paper]          | BPE       | RoPE       | GQA            | RMS Norm   | MLP          | Yes      |
| [Mixtral][mixtral-paper]        | BPE       | RoPE       | SWA            | RMS Norm   | MoE          | Yes      |
| [DeepseekV2][deepseek-v2-paper] | BPE       | RoPE       | MLA            | RMS Norm   | MoE          | Yes      |
| [DeepseekV3][deepseek-v3-paper] | BPE       | RoPE       | MLA            | RMS Norm   | MoE          | Yes      |
| [Kimi-K2][kimi-k2-paper]        | BPE       | RoPE       | MLA            | RMS Norm   | MoE          | Yes      |

*Note: Each model represents the largest variant within its respective series.*


## Properties

- **transformer** _object_, REQUIRED

  Contains the transformer configuration parameters.

  - **architecture_version** _string_, REQUIRED

    The version of the transformer architecture configuration using semantic versioning. An independent version is required for future extensibility.

  - **type** _string_, REQUIRED

    The type of transformer architecture. Currently supported: `decoder`. The default is `decoder`.

  - **vocabulary_size** _uint64_, REQUIRED

    Vocabulary size of the model.

  - **hidden_size** _uint64_, REQUIRED

    The hidden size of the model.

  - **tokenizer** _object_, REQUIRED

    Contains the tokenizer configuration parameters.

    - **type** _string_, REQUIRED

      Tokenizer type. Currently supported: `bpe`. The default is `bpe`.

    - **library** _string_, REQUIRED

      The name or URL of the tokenizer library. Currently supported: `huggingface`. The default is `huggingface`.

    - **revision** _string_, OPTIONAL

      Revision of the tokenizer library. Can be a branch name, tag name, commit ID, or `main` (latest version). The default is `main`.

  - **token_embedding** _object_, REQUIRED

    Contains the token embedding configuration parameters.

    - **has_bias** _boolean_, REQUIRED

      Whether the embedding has a bias. The default is `false`.

    - **has_norm** _boolean_, REQUIRED

      Whether the embedding has a normalization. The default is `true`. The normalization configuration is defined in the normalization property.

    - **shared_embedding** _boolean_, REQUIRED

      Whether the embedding is shared with the model prediction head. The default is `false`.

  - **position_embedding** _object_, REQUIRED

    Contains the position embedding configuration parameters.

    - **type** _string_, REQUIRED

      Position embedding type. Currently supported: `rope`. The default is `rope`. For more details, see [RoPE][rope-paper] and its [PyTorch implementation][rope-pytorch].

    - **max_position_embeddings** _uint64_, REQUIRED

      The maximum number of position embeddings. The default is `1024`.

    - **rope_theta** _float_, REQUIRED

      The theta parameter in the RoPE position embedding. The default is `10000`.

    - **rope_scaling** _object_, OPTIONAL

      The scaling configuration for the RoPE embeddings. The default is `null`.

  - **transformer_layer** _object_, REQUIRED

    Contains the transformer layer configuration parameters. Must specify either uniform_layers or mixed_layers.

    - **uniform_layers** _object_, OPTIONAL

      Configuration for uniform layers where all layers have identical structure.

      - **num_layers** _uint64_, REQUIRED

        Number of transformer layers. The default is `0`.

      - **attention** _object_, REQUIRED

        Contains the attention configuration parameters.

        - **type** _string_, REQUIRED

          Attention mechanism type. Currently supported: [MHA][mha-paper], [GQA][gqa-paper], [MLA][mla-paper]. The default is `mha`.

        - **is_causal** _boolean_, REQUIRED

          Whether the attention is causal. The default is `true`.

        - **is_qkv_merged** _boolean_, REQUIRED

          Whether the QKV projection is merged. The default is `false`.

        - **num_attention_heads** _uint64_, REQUIRED

          Number of attention heads. The default is `0`.

        - **num_key_value_heads** _uint64_, REQUIRED

          Number of key-value heads. The default is `0`.

        - **head_dim** _uint64_, REQUIRED

          The attention head dimension. If 0, defaults to hidden_size / num_attention_heads. The default is `0`.

        - **has_residual** _boolean_, REQUIRED

          Whether the attention has a residual connection. The default is `true`.

        - **has_qkv_bias** _boolean_, REQUIRED

          Whether the QKV projection has a bias. The default is `false`.

        - **has_output_bias** _boolean_, REQUIRED

          Whether the output projection has a bias. The default is `false`.

        - **has_pre_norm** _boolean_, REQUIRED

          Whether the attention has a pre-normalization. The default is `false`.

        - **has_post_norm** _boolean_, REQUIRED

          Whether the attention has a post-normalization. The default is `false`.

      - **mlp** _object_, OPTIONAL

        MLP configuration parameters. Either mlp or moe must be specified.

        - **intermediate_size** _uint64_, REQUIRED

          The size of the intermediate layer. The default is `0`.

        - **activation** _string_, REQUIRED

          The activation function. The default is `gelu`.

        - **use_gated_activation** _boolean_, REQUIRED

          Whether to use gated activation. The default is `true`.

        - **has_residual** _boolean_, REQUIRED

          Whether the MLP has a residual connection. The default is `true`.

        - **has_bias** _boolean_, REQUIRED

          Whether the MLP has a bias. The default is `false`.

        - **has_pre_norm** _boolean_, REQUIRED

          Whether the MLP has a pre-normalization. The default is `false`.

        - **has_post_norm** _boolean_, REQUIRED

          Whether the MLP has a post-normalization. The default is `false`.

        - **is_mlp_merged** _boolean_, REQUIRED

          Whether the MLP projection is merged. The default is `false`.

      - **moe** _object_, OPTIONAL

        MoE configuration parameters.

        - **has_bias** _boolean_, REQUIRED

          Whether the MOE has a bias. The default is `false`.

        - **activation** _string_, REQUIRED

          The activation function. The default is `gelu`.

        - **use_gated_activation** _boolean_, REQUIRED

          Whether to use gated activation. The default is `true`.

        - **num_experts** _uint64_, REQUIRED

          Number of experts. The default is `0`.

        - **moe_intermediate_size** _uint64_, REQUIRED

          The size of the intermediate layer of the routed expert. The default is `0`.

        - **num_shared_experts** _uint64_, REQUIRED

          Number of shared experts. The default is `0`.

        - **shared_expert_intermediate_size** _uint64_, REQUIRED

          The size of the intermediate layer of the shared expert. The default is `0`.

        - **top_k** _uint64_, REQUIRED

          Top k experts to be used. The default is `0`.

        - **scoring_function** _string_, REQUIRED

          Method of computing expert weights. The default is `softmax`.

        - **norm_topk_prob** _boolean_, REQUIRED

          Whether to normalize the top k probabilities. The default is `false`.

    - **mixed_layers** _object_, OPTIONAL

      Configuration for mixed layers where layers have different structures.

      - **num_layers** _uint64_, REQUIRED

        Number of transformer layers. The default is `0`.

      - **mlp_layers** _array_, REQUIRED

        Layers that use MLP. If empty, moe_frequency determines sparsity. The default is `[]`.

      - **pre_norm_layers** _array_, OPTIONAL

        Layers that use pre-normalization. The default is `[]`.

      - **post_norm_layers** _array_, OPTIONAL

        Layers that use post-normalization. The default is `[]`.

      - **moe_frequency** _uint64_, REQUIRED

        Frequency of the MoE layer. The default is `0`.

      - **attention** _object_, REQUIRED

        Attention parameters (same structure as in uniform_layers).

      - **mlp** _object_, OPTIONAL

        MLP parameters (same structure as in uniform_layers).

      - **moe** _object_, OPTIONAL

        MoE parameters (same structure as in uniform_layers).

  - **normalization** _object_, REQUIRED

    Contains the normalization configuration parameters.

    - **type** _string_, REQUIRED

      Normalization type. Supported: [`RMSNorm`][rmsnorm-paper], [`LayerNorm`][layernorm-paper]. The default is `rmsnorm`.

    - **epsilon** _float_, REQUIRED

      Epsilon for the normalization. The default is `1e-5`.

## Example

Here is an example transformer architecture configuration:

```json,title=Transformer%20Architecture%20Configuration&mediatype=application/vnd.cncf.model.architecture.v1%2Bjson
{
  "transformer": {
    "vocabulary_size": 32000,
    "hidden_size": 4096,
    "tokenizer": {
      "type": "bpe",
      "library": "huggingface",
      "revision": "main"
    },
    "token_embedding": {
      "has_bias": false,
      "has_norm": true,
      "shared_embedding": false
    },
    "position_embedding": {
      "type": "rope",
      "max_position_embeddings": 2048,
      "rope_theta": 10000.0,
      "rope_scaling": null
    },
    "transformer_layer": {
      "uniform_layers": {
        "num_layers": 32,
        "attention": {
          "type": "gqa",
          "is_causal": true,
          "is_qkv_merged": false,
          "num_attention_heads": 32,
          "num_key_value_heads": 8,
          "head_dim": 128,
          "has_residual": true,
          "has_qkv_bias": false,
          "has_output_bias": false,
          "has_pre_norm": true,
          "has_post_norm": false
        },
        "mlp": {
          "intermediate_size": 11008,
          "activation": "silu",
          "use_gated_activation": true,
          "has_residual": true,
          "has_bias": false,
          "has_pre_norm": false,
          "has_post_norm": true,
          "is_mlp_merged": false
        }
      }
    },
    "normalization": {
      "type": "rmsnorm",
      "epsilon": 1e-5
    }
  }
}
```

[attention-paper]: https://arxiv.org/abs/1706.03762
[gpt2-repo]: https://github.com/openai/gpt-2
[llama2-paper]: https://arxiv.org/abs/2307.09288
[llama3-paper]: https://arxiv.org/abs/2407.21783
[qwen2-paper]: https://arxiv.org/abs/2407.10671
[qwen3-paper]: https://arxiv.org/pdf/2505.09388
[gemma2-paper]: https://arxiv.org/abs/2408.00118
[gemma3-paper]: https://arxiv.org/pdf/2503.19786
[mixtral-paper]: https://arxiv.org/abs/2401.04088
[deepseek-v2-paper]: https://arxiv.org/abs/2405.04434
[deepseek-v3-paper]: https://arxiv.org/pdf/2412.19437
[kimi-k2-paper]: https://arxiv.org/pdf/2507.20534
[rope-paper]: https://arxiv.org/abs/2104.09864
[rope-pytorch]: https://pytorch.org/torchtune/stable/generated/torchtune.modules.RotaryPositionalEmbeddings.html
[mha-paper]: https://arxiv.org/abs/1706.03762
[gqa-paper]: https://arxiv.org/abs/2305.13245v3
[mla-paper]: https://arxiv.org/abs/2412.19437
[rmsnorm-paper]: https://arxiv.org/abs/1910.07467
[layernorm-paper]: https://arxiv.org/abs/1607.06450
