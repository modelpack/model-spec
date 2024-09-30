module github.com/CloudNativeAI/model-spec/tools/mdctl

go 1.23.1

replace github.com/CloudNativeAI/model-spec/specs-go => ../../specs-go/

require (
	github.com/CloudNativeAI/model-spec/specs-go v0.0.0-00010101000000-000000000000
	github.com/klauspost/compress v1.17.10
	github.com/opencontainers/go-digest v1.0.0
	github.com/opencontainers/image-spec v1.1.0
	github.com/spf13/cobra v1.8.1
	golang.org/x/term v0.24.0
	oras.land/oras-go/v2 v2.5.0
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sync v0.6.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
)
