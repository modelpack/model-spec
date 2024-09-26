module github.com/CloudNativeAI/model-spec/tools/mdctl

go 1.22.4

require (
	github.com/klauspost/compress v1.17.7
	github.com/opencontainers/go-digest v1.0.0
	github.com/opencontainers/image-spec v1.1.0
	github.com/spf13/cobra v1.7.0
	oras.land/oras-go/v2 v2.4.0
)

require (
	github.com/CloudNativeAI/model-spec/specs-go v0.0.0-20240926032628-0609fae554bf // indirect
	golang.org/x/sync v0.6.0 // indirect
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/term v0.15.0
)

replace github.com/CloudNativeAI/model-spec/specs-go/ => ../../specs-go/
