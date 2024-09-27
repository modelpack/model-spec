package v2

// Architecture represents the architecture of the model.
type Architecture struct {
	// Transformer architecture
	Transformer Transformer `json:"transformer"`

	// TODO: Other architectures
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
	Attention Attention `json:"attention"`
	MLP       MLP       `json:"mlp"`
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
	// Type of position embedding, e.g. 'rope', 'sinusoidal', 'alibi', etc.
	Type string `json:"type"`

	// The maximum number of position embeddings
	MaxPositionEmbeddings int `json:"max_position_embeddings"`

	// The base in signifying the rotary embedding period.
	RotaryEmbeddingBase int `json:"rotary_embedding_base,omitempty"`

	// Fraction of hidden size to apply rotary embeddings to. Must be in [0,1].
	RotaryEmbeddingFraction float64 `json:"rotary_embedding_fraction,omitempty"`
}

// Normalization represents the normalization parameters.
type Normalization struct {
	// Type of normalization, e.g. 'layernorm', 'rmsnorm', etc.
	Type string `json:"type"`

	// Epsilon for the normalization
	Epsilon float64 `json:"epsilon"`
}
