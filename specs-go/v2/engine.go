package v2

import oci "github.com/opencontainers/image-spec/specs-go/v1"

// Engine provides the structure for the `application/vnd.cnai.models.engine.v0+json` mediatype when marshalled to JSON.
// It encapsulates the details necessary to describe and configure the execution environment for a model.
type Engine struct {
	// Name specifies the engine or framework used, such as 'transformers', 'tensorrt', or 'vllm'.
	// This field is crucial for identifying the runtime environment required for the model.
	Name string `json:"name,omitempty"`

	// Version indicates the specific version of the engine or framework.
	// Examples include '4.44.0', '8.10', '1.0', etc. This ensures compatibility and reproducibility.
	Version string `json:"version,omitempty"`

	// Dependencies lists the additional packages or libraries required by the engine.
	// This optional field is used to specify and install necessary components for the engine's operation.
	Dependencies []string `json:"dependencies,omitempty"`

	// Environment defines key-value pairs for environment variables.
	// These variables are used to configure the runtime environment for the engine executor.
	Environment map[string]string `json:"environment,omitempty"`

	// EntryPoint specifies the command or script to initiate the engine.
	// This optional field provides the starting point for executing the model within the engine.
	EntryPoint string `json:"entrypoint,omitempty"`

	// Extensions allows for additional, engine-specific configuration details.
	// Each extension is represented by a Descriptor, enabling flexible and extensible metadata
	// to accommodate unique requirements or features of different engine types.
	Extensions []oci.Descriptor `json:"extensions,omitempty"`
}
