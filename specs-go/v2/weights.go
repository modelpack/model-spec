package v2

import (
	oci "github.com/opencontainers/image-spec/specs-go/v1"
)

// Weights represents the structure for the `application/vnd.cnai.models.weights.v0+json` mediatype when marshalled to JSON.
// It encapsulates the essential information about a model's weights, including their storage format, numerical precision, and file references.
type Weights struct {
	// File is an array of Descriptors referencing the inline files or directories containing the model weights.
	// These Descriptors provide details such as the file size, digest, and media type of each weight file or directory.
	File []oci.Descriptor `json:"file,omitempty"`

	// Format specifies the storage format of the weights. This field can include values such as:
	// - 'safetensors': A fast and safe format for storing tensors
	// - 'gguf': GPT-Generated Unified Format, used by some language models
	// - 'onnx': Open Neural Network Exchange format
	// - 'pytorch': PyTorch's native serialization format
	// The format information is crucial for correctly loading and interpreting the weight data.
	Format string `json:"format,omitempty"`

	// Precision indicates the numerical precision of the weights. This field can include values such as:
	// - 'bf16': Brain Floating Point (bfloat16)
	// - 'fp16': Half-precision floating-point
	// - 'fp32': Single-precision floating-point
	// - 'int8': 8-bit integer quantization
	// The precision information is essential for memory management and computational efficiency.
	Precision string `json:"precision,omitempty"`
}
