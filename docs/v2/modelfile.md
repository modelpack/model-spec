# Introduction to Modelfile

A Modelfile is a text file containing all commands, in order, needed to build a given model image. It automates the process of building model images.

## Modelfile Instructions

| **Instruction** | **Description** |
| --- | --- |
| CREATE | Create a new model image |
| FROM | Specify the base model image to use |
| NAME | Specify model name |
| FAMILY | Specify model family |
| ARCHITECTURE | Specify model architecture |
| LICENSE | Specify the legal license under which the model is used |
| CONFIG | Specify model configuration file |
| WEIGHTS | Specify model weights file |
| FORMAT | Specify model weights format |
| TOKENIZER | Specify tokenizer configuration |

## Modelfile Example

```plain
CREATE registry.cnai.com/sys/gemma-2b:latest

# Model Information

NAME         gemma-2b
FAMILY       gemma
ARCHITECTURE transformer
FORMAT       safetensors

# Model License

LICENSE      examples/huggingface/gemma-2b/LICENSE

# Model Configuration

CONFIG       examples/huggingface/gemma-2b/config.json
CONFIG       examples/huggingface/gemma-2b/generation_config.json

# Model Tokenizer

TOKENIZER    examples/huggingface/gemma-2b/tokenizer.json

# Model Weights

WEIGHTS      examples/huggingface/gemma-2b/model.safetensors.index.json
WEIGHTS      examples/huggingface/gemma-2b/model-00001-of-00002.safetensors
WEIGHTS      examples/huggingface/gemma-2b/model-00002-of-00002.safetensors

```

## Management tool

We propose a model management tool, which is a command-line tool for building, managing, and running AI models.

### build

We can use Modelfile to build model images.

```plain
mdctl build -f ./Modelfile
```

### list

We can list all the model images that have been pushed.

```plain
mdctl list
```

### push

We can push the built model image to a model repository.

```plain
mdctl push <model-image>
```

### pull

We can pull the model image from the model repository to local storage.

```plain
mdctl pull <model-image>
```

### unpack

We can pull the model image to local storage and then use mdctl to run the model.

```plain
mdctl unpack <model-image>
```
