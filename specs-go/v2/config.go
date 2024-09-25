package v2

import (
	oci "github.com/opencontainers/image-spec/specs-go/v1"
)

// Config represents the JSON structure that encapsulates essential metadata and configuration details of a machine learning model.
type Config struct {
	// Name specifies the unique identifier or title of the model.
	Name string `json:"name"`

	// Family indicates the broader category or lineage of the model, such as 'GPT', 'LLAMA', or 'QWEN'.
	// This helps in grouping related models or identifying their general capabilities.
	Family string `json:"family"`

	// Architecture defines the fundamental structure or design of the model,
	// such as 'transformer', 'CNN' (Convolutional Neural Network), 'RNN' (Recurrent Neural Network), etc.
	// This information is crucial for understanding the model's underlying principles and potential applications.
	Architecture string `json:"architecture"`

	// Description provides detailed information about the model's purpose, capabilities, and usage.
	// It is represented as an array of Descriptors, allowing for rich, structured content.
	Description []oci.Descriptor `json:"description,omitempty"`

	// License contains the legal and usage terms associated with the model.
	// It includes policies and grants that govern how the model can be used, distributed, or modified.
	// Represented as an array of Descriptors to accommodate multiple or complex licensing terms.
	License []oci.Descriptor `json:"license,omitempty"`

	// Extensions allows for the inclusion of additional, model-specific configuration details.
	// Each extension is represented by a Descriptor, enabling flexible and extensible metadata.
	// This field accommodates unique requirements or features of different model types, such as:
	// - Generation configuration: Parameters for text generation in language models
	// - Quantization configuration: Details about model weight quantization
	// - Transformer configuration: Specific architectural details for transformer models
	// - Domain-specific settings: Configurations relevant to particular application domains
	// The use of Descriptors ensures that each extension can be properly identified and processed,
	// allowing for seamless integration of diverse model configurations within a unified structure.
	Extensions []oci.Descriptor `json:"extensions,omitempty"`
}
