package v2

// TransformerForCausalLLM represents the configuration of a transformer model for causal language modeling.
// It defines the architecture and hyperparameters of the model.
//
// Supported features:
// - Attention mechanisms: Multi-Head Attention (MHA) and Grouped Query Attention (GQA)
// - Activation functions: GELU, ReLU
// - Position embeddings: Rotary Position Embedding (RoPE)
// - Normalization: RMSNorm (Root Mean Square Layer Normalization)
//
// This structure is designed to be flexible and accommodate various transformer architectures
// used in state-of-the-art language models.
type TransformerForCausalLLM struct {
	// Version of the transformer architecture config
	Version string `json:"version"`

	// Vocabulary size of the model
	VocabularySize int `json:"vocabulary_size"`

	// The hidden size of the model, e.g. 768, 1024, 2048, etc.
	HiddenSize int `json:"hidden_size"`

	// The number of transformer layers of the model.
	NumHiddenLayers int `json:"num_hidden_layers"`

	// The number of attention heads, e.g. 12, 16, 32, etc.
	NumAttentionHeads int `json:"num_attention_heads"`

	// The number of key value heads, e.g. 1, 2, 4, etc.
	// Only used by GQA attention mechanism.
	NumKeyValueHeads int `json:"num_key_value_heads"`

	// The activation function used by the pointwise feed-forward layers, e.g. 'gelu', 'relu', 'tanh', etc.
	Activation string `json:"activation"`

	// The intermediate size in the feed-forward layers. The non-linearity is applied in this intermediate size.
	IntermediateSize int `json:"intermediate_size"`

	// The rms_norm parameter
	NormEpsilon float64 `json:"norm_epsilon"`

	// The position embedding type, for example 'rope', 'sinusoidal', 'alibi', etc.
	PositionEmbedding string `json:"position_embedding"`

	// The base in signifying the rotary embedding period.
	RotaryEmbeddingBase int `json:"rotary_embedding_base,omitempty"`

	// Fraction of hidden size to apply rotary embeddings to. Must be in [0,1].
	RotaryEmbeddingFraction float64 `json:"rotary_embedding_fraction,omitempty"`
}
