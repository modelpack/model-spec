# Model Specification Version 2

## Overview

The core of the v2 model specification is the definition of the model artifact, metadata and runtime environment.

The model artifact is a collection of files that represent the AI model. It consists of the model configuration, model weights, model tokenizer, and other model resources.

The model metadata is general information about the model, such as the model name, version, model family, description, author, license, and architecture. A model registry can parse the model metadata to display the model information.

The model runtime environment is the environment in which the model runs. It includes the inference engine information, such as verion, configuration, dependencies, and environment variables.

The model artifact, metadata and runtime environment are organized in a model manifest, which is a JSON file that describes the model. The model manifest is used to package and distribute the model, and can be stored in a model registry and downloaded by a model runtime.

With a proper defined model specification, we can package AI models of a model repository into a model image, and push the model image to the model registry. The model image can be pulled and run by the model runtime, either as a standalone package or as a readonly volume source in a container.

## Goals

The goals of developing the model specification are:

* To provide a way for developers to package and distribute AI models in a cloud native environment.
* To promote AI models as a first-class citizen and pave the way for the infrastructure to be organized around AI models.
* To define general model artifact, metadata, and runtime environment, so that the model can be easily understood and managed by any components of the infrastructure.
* To define a general model format description to allow easy integration of models with model runtimes.

## Non-Goals

* To build standard interfaces for model management tools to build, distribute, manage, and run AI models.

The model specification is designed to be a foundation for building standard interfaces to build, distribute, manage, and run AI models. But the model specification itself does not define such standard interfaces.

## Plans

The model specification is still pretty rough. It is a living document and will evolve over time. Future work includes:

* Figure out the details of AI model artifact, metadata, and runtime environment.
* Define a general transformer architecture abstraction to support build once and run everywhere of LLMs.
* Develop tools to build and save AI models in a model registry.
* Develop tools to pull and run AI models in a model runtime.
* Modify [vllm](https://github.com/vllm-project/vllm) to support the model specification and run any transformer architecture LLMs without modification.
