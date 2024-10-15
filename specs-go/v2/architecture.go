package v2

// Architecture represents the architecture of the model.
type Architecture struct {
	// Transformer architecture
	Transformer Transformer `json:"transformer"`

	// TODO: Other architectures, like mamba, etc.
}

// Transformer represents the transformer architecture.
type Transformer struct {
	// Transformer for causal language modeling
	CausalLM TransformerForCausalLM `json:"causal_lm"`

	// multi-modal transformer
	// MultiModal TransformerForMultiModal `json:"multi_modal"`
}

// TransformerForCausalLM represents the transformer architecture for causal language modeling.
// This structure is designed to be flexible and accommodate various transformer architectures
// used in state-of-the-art language models.
type TransformerForCausalLM struct {
	// Version of the transformer architecture
	Version string `json:"version"`

	// Vocabulary size of the model
	VocabularySize int `json:"vocabulary_size"`

	// The hidden size of the model
	HiddenSize int `json:"hidden_size"`

	// embedding
	Embedding Embedding `json:"embedding"`

	// Position embedding type
	PositionEmbedding PositionEmbedding `json:"position_embedding"`

	// Number of transformer layers
	NumTransformerLayers int `json:"num_transformer_layers"`

	// Transformer layer
	TransformerLayer TransformerLayer `json:"transformer_layer"`

	// Normalization parameters
	Normalization Normalization `json:"normalization"`
}

// TransformerLayer represents the transformer layer parameters.
type TransformerLayer struct {
	// Attention parameters
	Attention Attention `json:"attention"`

	// MLP parameters
	MLP MLP `json:"mlp"`
}

// MLP represents the MLP (Multi-Layer Perceptron) parameters.
// TODO: Add support for other MLP architectures, such as MoE.
type MLP struct {
	// The size of the intermediate layer
	IntermediateSize int `json:"intermediate_size"`

	// Activation function
	Activation string `json:"activation"`

	// Whether to use gated activation
	UseGatedActivation bool `json:"use_gated_activation"`

	// Whether the MLP has a residual connection
	HasResidual bool `json:"has_residual"`

	// Whether the MLP has a bias
	HasBias bool `json:"has_bias"`

	// Whether the MLP has a pre-normalization
	HasPreNorm bool `json:"has_pre_norm"`

	// Whether the MLP has a post-normalization
	HasPostNorm bool `json:"has_post_norm"`
}

// Attention represents the parameters for various attention mechanisms.
//
// Supported attention types:
// - Multi-Head Attention (MHA):
// - Grouped Query Attention (GQA):
// - Multi-Query Attention (MQA):
//
// TODO: Add support for other attention mechanisms:
// - Sliding Window Attention (SWA)
// - Multi-Linear Attention (MLA)
type Attention struct {
	// Attention type
	AttentionType string `json:"attention_type"`

	// Whether the attention is causal
	IsCausal bool `json:"is_causal"`

	// Number of attention heads
	NumAttentionHeads int `json:"num_attention_heads"`

	// Number of key-value heads
	NumKeyValueHeads int `json:"num_key_value_heads"`

	// The attention head dimension. If 0, it will default to hidden_size / NumAttentionHeads
	HeadDim int `json:"head_dim"`

	// Whether the attention has a residual connection
	HasResidual bool `json:"has_residual"`

	// Whether the attention has a bias
	HasBias bool `json:"has_bias"`

	// Whether the attention has a pre-normalization
	HasPreNorm bool `json:"has_pre_norm"`

	// Whether the attention has a post-normalization
	HasPostNorm bool `json:"has_post_norm"`
}

// PositionEmbedding represents the position embedding type and parameters.
type PositionEmbedding struct {
	// Type of position embedding, e.g. 'rope', 'alibi', etc.
	Type string `json:"type"`

	// The maximum number of position embeddings
	MaxPositionEmbeddings int `json:"max_position_embeddings"`

	// Only used with 'RoPE'. The theta parameter in the RoPE position embedding.
	RotaryEmbeddingTheta float64 `json:"rope_theta,omitempty"`

	// Only used with 'RoPE'. The scaling configuration for the RoPE embeddings
	RotaryEmbeddingScaling RotaryEmbeddingScaling `json:"rope_scaling,omitempty"`
}

// RotaryEmbeddingScaling represents the scaling configuration for the RoPE embeddings.
type RotaryEmbeddingScaling struct {
	// Type of scaling, can be one of ['default', 'linear', 'dynamic', 'llama3'], with 'default' being the original RoPE implementation.
	Type string `json:"type"`

	// The scaling factor
	Factor float64 `json:"factor"`

	// The original max position used during pretraining.
	OriginalMaxPosition int `json:"original_max_position"`

	// Only used with 'llama3'. Scaling factor applied to low frequency components of the RoPE
	LowFreqFactor float64 `json:"low_freq_factor"`

	// Only used with 'llama3'. Scaling factor applied to high frequency components of the RoPE
	HighFreqFactor float64 `json:"high_freq_factor"`
}

// Normalization represents the normalization parameters.
type Normalization struct {
	// Type of normalization, e.g. 'layernorm', 'rmsnorm', etc.
	Type string `json:"type"`

	// Epsilon for the normalization
	Epsilon float64 `json:"epsilon"`
}

// Embedding represents the embedding parameters.
type Embedding struct {
	// Whether the embedding has a bias
	HasBias bool `json:"has_bias"`

	// Whether the embedding has a normalization
	HasNorm bool `json:"has_norm"`
}
