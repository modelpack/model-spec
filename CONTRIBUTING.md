# Contributors' Guide

## Getting started

Welcome to the Cloud Native AI Model Format Specification project! We are excited to have you contribute. Here are some steps to help you get started.

## Setting up your local environment

* **Clone the repository**:

```sh
git clone https://github.com/modelpack/model-spec.git
cd model-spec
```

* **Install dependencies**: Ensure you have [Go](https://go.dev/) installed, as the current spec implementation is written in Go. Follow the [official instructions to install Go](https://go.dev/doc/install).

## Where to put changes

Right now, we have a simple directory structure:

* `docs`: All detailed documents about the model spec.
* `docs/img`: Any referenced images in the documents should be put here.
* `specs-go`: A Go implementation of the model specification.

## Raise a pull request

* **Create a new branch**:

```sh
git checkout -b your-branch-name
```

* **Make your changes and commit them**:

```sh
git add .
git commit -s -m "Your descriptive commit message"
```

* **Push your changes to your fork**:

```sh
git push your-fork-repo your-branch-name
```

* **Open a pull request**: Go to the GitHub repository, compare your branch, and submit a pull request with a detailed description of your changes.

## Make sure pull request CI passes

Please check the CI status in your pull request and fix anything that fails. Here are some simple instructions to validate CI locally.

* **Install golangci-lint**: follow the [official installation guide](https://golangci-lint.run/welcome/install/#local-installation) to install golangci-lint.

* **Check for linting issues**:

```sh
golangci-lint run --verbose
```

We appreciate your contributions and look forward to working with you!
