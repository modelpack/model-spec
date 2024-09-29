package v2

// manifest
const (
	// MediaTypeModelManifest specifies the media type for a models manifest.
	MediaTypeModelManifest = "application/vnd.cnai.model.manifest.v2+json"
)

// configs
const (
	// MediaTypeModelConfig specifies the media type for model configuration.
	MediaTypeModelConfig = "application/vnd.cnai.model.config.v2+json"

	// MediaTypeModelLicense specifies the media type for model license.
	MediaTypeModelLicense = "application/vnd.cnai.model.license.v2+plaintext"

	// MediaTypeModelDescription specifies the media type for model description.
	MediaTypeModelDescription = "application/vnd.cnai.model.description.v2+plaintext"

	// MediaTypeModelExtension specifies the media type for model configuration extension.
	MediaTypeModelExtension = "application/vnd.cnai.model.extension.v2+json"
)

// processors
const (
	// MediaTypeModelProcessorText specifies the media type for text processors.
	// This includes tokenizers like sentencepiece, used for processing textual input.
	MediaTypeModelProcessorText = "application/vnd.cnai.model.processor.text.v2.tar"

	// MediaTypeModelProcessorAudio specifies the media type for audio processors.
	// These are used for processing audio input, such as speech-to-text models.
	MediaTypeModelProcessorAudio = "application/vnd.cnai.model.processor.audio.v2.tar"

	// MediaTypeModelProcessorImage specifies the media type for image processors.
	// These are used for processing image input, such as in computer vision models.
	MediaTypeModelProcessorImage = "application/vnd.cnai.model.processor.image.v2.tar"

	// MediaTypeModelProcessorMultiModal specifies the media type for multi-modal processors.
	// These are used for models that can process multiple types of input (e.g., text and images).
	MediaTypeModelProcessorMultiModal = "application/vnd.cnai.model.processor.multimodal.v2.tar"
)

// weights
const (
	// MediaTypeModelWeights specifies the media type for model weights.
	MediaTypeModelWeights = "application/vnd.cnai.model.weights.v2.tar"
)

// engine
const (
	// MediaTypeModelEngine specifies the media type for model engine.
	MediaTypeModelEngine = "application/vnd.cnai.model.engine.v2.tar"
)

// transformer architecture
const (
	// MediaTypeModelArchitectureTransformer specifies the media type for model architecture.
	MediaTypeModelArchitectureTransformer = "application/vnd.cnai.model.architecture.transformer.v2.tar"
)
