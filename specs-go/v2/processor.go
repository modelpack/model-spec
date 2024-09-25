package v2

import (
	oci "github.com/opencontainers/image-spec/specs-go/v1"
)

// TextProcessor represents the structure for the `application/vnd.cnai.models.tokenizer.v0+json` mediatype when marshalled to JSON.
// It encapsulates the essential components of a text tokenizer used in natural language processing models.
type TextProcessor struct {
	// TokenizerConfig is a Descriptor referencing the configuration file(s) for the tokenizer.
	// This can be a single file or multiple files containing essential information such as:
	// - Vocabulary: The set of tokens used by the tokenizer
	// - Settings: Parameters that control tokenization behavior
	// - Special tokens: Tokens with specific meanings or functions (e.g., [PAD], [CLS], [SEP])
	// Modern tokenizers often consolidate all configuration into a single file for simplicity,
	// while some may still use separate files for different components.
	TokenizerConfig oci.Descriptor `json:"tokenizer_config,omitempty"`

	// Algorithm is the tokenization algorithm used by the tokenizer, such as BPE, WordPiece, Unigram, etc.
	Algorithm string `json:"algorithm,omitempty"`

	// Library is the library used by the tokenizer, such as sentencepiece, tiktoken, huggingface tokenizers, etc.
	Library string `json:"library,omitempty"`
}

// AudioProcessor represents the structure for the `application/vnd.cnai.models.processor.audio.v2+json` mediatype when marshalled to JSON.
// It encapsulates the essential components of an audio processor used in audio processing models.
type AudioProcessor struct {
	// TODO: to be defined
}

// ImageProcessor represents the structure for the `application/vnd.cnai.models.processor.image.v2+json` mediatype when marshalled to JSON.
// It encapsulates the essential components of an image processor used in image processing models.
type ImageProcessor struct {
	// TODO: to be defined
}

// MultiModalProcessor represents the structure for the `application/vnd.cnai.models.processor.multimodal.v2+json` mediatype when marshalled to JSON.
// It encapsulates the essential components of a multi-modal processor used in multi-modal processing models.
type MultiModalProcessor struct {
	// TODO: to be defined
}
