/*
 *     Copyright 2025 The CNCF ModelPack Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1

import (
	"time"

	digest "github.com/opencontainers/go-digest"
)

// ModelConfig defines the execution parameters
// which should be used as a base when running a model using an inference engine.
type ModelConfig struct {
	// The model architecture, such as transformer, cnn, rnn, etc.
	Architecture string `json:"architecture,omitempty"`

	// The model format, such as onnx, tensorflow, pytorch, etc.
	Format string `json:"format,omitempty"`

	// The size of the model parameters, such as "8b", "16b", "32b", etc.
	ParamSize string `json:"paramSize,omitempty"`

	// The model precision, such as bf16, fp16, int8, mixed etc.
	Precision string `json:"precision,omitempty"`

	// The model quantization, such as awq, gptq, etc
	Quantization string `json:"quantization,omitempty"`

	// TransformerConfig contains Transformer-specific architectural parameters.
	// Populated only when Architecture is "transformer".
	TransformerConfig *TransformerConfig `json:"transformerConfig,omitempty"`

	// Special capabilities that the model supports
	Capabilities *ModelCapabilities `json:"capabilities,omitempty"`
}

// ModelFS describes a layer content addresses
type ModelFS struct {
	// Type is the type of the rootfs. MUST be set to "layers".
	Type string `json:"type"`

	// DiffIDs is an array of layer content hashes (DiffIDs), in order from bottom-most to top-most.
	DiffIDs []digest.Digest `json:"diffIds"`
}

// ModelDescriptor defines the general information of a model
type ModelDescriptor struct {
	// Date and time on which the model was built
	CreatedAt *time.Time `json:"createdAt,omitempty"`

	// The contact details of the people or organization responsible for the model
	Authors []string `json:"authors,omitempty"`

	// The model family, such as llama3, gpt2, qwen2, etc.
	Family string `json:"family,omitempty"`

	// The model name, such as llama3-8b-instruct, gpt2-xl, qwen2-vl-72b-instruct, etc.
	Name string `json:"name,omitempty"`

	// The URL to get documentation on the model
	DocURL string `json:"docURL,omitempty"`

	// The URL to get source code for building the model
	SourceURL string `json:"sourceURL,omitempty"`

	// The URLs to reference to datasets that the model was trained upon.
	DatasetsURL []string `json:"datasetsURL,omitempty"`

	// The version of the packaged software
	Version string `json:"version,omitempty"`

	// The source control revision identifier for the packaged software
	Revision string `json:"revision,omitempty"`

	// The name of the distributing entity, organization or individual
	Vendor string `json:"vendor,omitempty"`

	// The license(s) under which contained software is distributed as an SPDX License Expression
	Licenses []string `json:"licenses,omitempty"`

	// The human-readable title of the model
	Title string `json:"title,omitempty"`

	// The human-readable description of the software packaged in the model
	Description string `json:"description,omitempty"`
}

// AttentionType defines the attention mechanism used in a Transformer model.
type AttentionType string

const (
	// MultiHeadAttention is the standard multi-head attention (MHA) from "Attention Is All You Need".
	MultiHeadAttention AttentionType = "mha"

	// GroupedQueryAttention uses fewer KV heads than query heads, reducing KV cache size (e.g. LLaMA 3, Mistral).
	GroupedQueryAttention AttentionType = "gqa"

	// MultiLatentAttention uses low-rank key-value compression to minimize KV cache (e.g. DeepSeek-V2).
	MultiLatentAttention AttentionType = "mla"
)

// MLPType defines the feed-forward / MLP layer variant used in a Transformer model.
type MLPType string

const (
	// DenseFF is a standard dense feed-forward layer.
	DenseFF MLPType = "dense"

	// MixtureOfExperts routes tokens to a subset of expert FFN layers (e.g. Mixtral, DeepSeek-V3).
	MixtureOfExperts MLPType = "moe"
)

// TransformerConfig captures the key architectural parameters of a Transformer model.
// It is only populated when ModelConfig.Architecture is "transformer".
type TransformerConfig struct {
	// AttentionType is the attention mechanism variant (mha, gqa, mla).
	AttentionType AttentionType `json:"attentionType,omitempty"`

	// MLPType is the feed-forward layer variant (dense, moe).
	MLPType MLPType `json:"mlpType,omitempty"`

	// NumLayers is the total number of transformer layers (blocks).
	NumLayers *int `json:"numLayers,omitempty"`

	// NumAttentionHeads is the number of query attention heads.
	NumAttentionHeads *int `json:"numAttentionHeads,omitempty"`

	// NumKVHeads is the number of key/value heads. For GQA this is less than NumAttentionHeads.
	// Omitted or equal to NumAttentionHeads implies full MHA.
	NumKVHeads *int `json:"numKVHeads,omitempty"`

	// HiddenSize is the model's hidden dimension (d_model).
	HiddenSize *int `json:"hiddenSize,omitempty"`

	// IntermediateSize is the inner dimension of the feed-forward layer.
	IntermediateSize *int `json:"intermediateSize,omitempty"`
}

// Modality defines the input and output types of the model
// such as text, image, audio, video, etc.
// It is used to define the input and output types of the model.
type Modality string

const (
	// TextModality indicates that the model can process or generate text.
	// If present in InputTypes, the model accepts text as input.
	// If present in OutputTypes, the model produces text as output.
	TextModality Modality = "text"

	// ImageModality indicates that the model can process or generate images.
	// If present in InputTypes, the model accepts images as input.
	// If present in OutputTypes, the model produces images as output.
	ImageModality Modality = "image"

	// AudioModality indicates that the model can process or generate audio.
	// If present in InputTypes, the model accepts audio as input.
	// If present in OutputTypes, the model produces audio as output.
	AudioModality Modality = "audio"

	// VideoModality indicates that the model can process or generate video.
	// If present in InputTypes, the model accepts video as input.
	// If present in OutputTypes, the model produces video as output.
	VideoModality Modality = "video"

	// EmbeddingModality indicates that the model can process or generate embeddings.
	// If present in InputTypes, the model accepts embeddings as input.
	// If present in OutputTypes, the model produces embeddings as output.
	// For a dedicated "embedding model", EmbeddingModality should be present in its OutputTypes.
	EmbeddingModality Modality = "embedding"

	// OtherModality indicates that the model supports a modality not explicitly listed.
	OtherModality Modality = "other"
)

// ModelCapabilities defines the special capabilities that the model supports
type ModelCapabilities struct {
	// The model supports the following input types
	InputTypes []Modality `json:"inputTypes,omitempty"`

	// The model supports the following output types
	OutputTypes []Modality `json:"outputTypes,omitempty"`

	// KnowledgeCutoff is the date of the datasets that the model was trained on, formatted as defined by RFC 3339
	KnowledgeCutoff *time.Time `json:"knowledgeCutoff,omitempty"`

	// Reasoning indicates whether the model can perform reasoning tasks
	Reasoning *bool `json:"reasoning,omitempty"`

	// ToolUsage indicates whether the model can use external tools
	// such as a calculator, a search engine, etc.
	ToolUsage *bool `json:"toolUsage,omitempty"`

	// Reward indicates whether the model is a reward model
	Reward *bool `json:"reward,omitempty"`

	// Language indicates the languages that the model can speak.
	// Encoded as ISO 639 two letter codes. For example, ["en", "fr", "zh"].
	Languages []string `json:"languages,omitempty"`
}

// Model defines the basic information of a model.
// It provides the `application/vnd.cncf.model.config.v1+json` mediatype when marshalled to JSON.
type Model struct {
	// The model descriptor
	Descriptor ModelDescriptor `json:"descriptor"`

	// The model describes a layer content addresses
	ModelFS ModelFS `json:"modelfs"`

	// Config defines the execution parameters which should be used as a base when running a model using an inference engine.
	Config ModelConfig `json:"config,omitempty"`
}
