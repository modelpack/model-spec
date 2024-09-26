# How to use mdctl examples

## Download model

Download model from huggingface:

```plain
git lfs install
git clone https://huggingface.co/gemma-ai/gemma-2b
```

## Build model image

Put the modelfile to the model directory and build model image:

```plain
mdctl build -f Modelfile
```
