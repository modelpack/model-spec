# mdctl - Model Control Tool

`mdctl` is a command-line tool for building, managing, and running AI models.

## Installation

To install `mdctl`, clone the repository and build the binary:

```plain
git clone https://github.com/CloudNativeAI/mdctl.git
cd mdctl
go build
```

## Usage

To build a model, use the `build` command:

```plain
./mdctl build -f Modelfile
```

To list all models, use the `list` command:

```plain
./mdctl list
```

To push a model, use the `push` command. Before pushing, you need to set the model registry credentials:

```plain
export MODEL_REGISTRY_USER=<username>
export MODEL_REGISTRY_PASSWORD=<password>
export MODEL_REGISTRY_URL=<registry_url>
```

```plain
./mdctl push <model>
```

To pull a model, use the `pull` command:

```plain
./mdctl pull <model>
```

To run a model, use the `unpack` command:

```plain
./mdctl unpack -n <model>
```
