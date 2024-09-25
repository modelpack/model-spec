package v2

import (
	oci "github.com/opencontainers/image-spec/specs-go/v1"
)

// Manifest represents the structure for the `application/vnd.cncf.cnai.models.manifest.v2+json` mediatype when marshalled to JSON.
// It encapsulates all the essential components and metadata for a machine learning model.
type Manifest struct {
	// Version specifies the version of the manifest schema.
	Version string `json:"version"`

	// MediaType indicates the specific type of this document's data structure.
	// It should be set to `application/vnd.cnai.models.manifest.v2+json` or an applicable IANA media type.
	MediaType string `json:"mediaType"`

	// Config references the configuration object for the model.
	// This JSON blob contains essential setup information used by the runtime.
	Config Config `json:"config"`

	// Processor references the pre-processor object(s) by digest.
	// It's used for any data preparation or transformation required before model inference.
	Processor []oci.Descriptor `json:"processor"`

	// Weights references the model's weight object by digest.
	// These are typically binary blobs containing the trained parameters of the model.
	Weights Weights `json:"weights"`

	// Engine is an optional field that references the engine object by digest.
	// The engine structure contains information for setting up the runtime environment.
	Engine Engine `json:"engine,omitempty"`

	// Annotations is an optional map for storing arbitrary metadata related to the model manifest.
	// This can include information like creation date, author, or custom tags.
	Annotations map[string]string `json:"annotations,omitempty"`
}
