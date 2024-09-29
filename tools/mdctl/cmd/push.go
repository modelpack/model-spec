package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	modelspec "github.com/CloudNativeAI/model-spec/specs-go/v2"
	"github.com/CloudNativeAI/model-spec/tools/mdctl/models"
	"github.com/CloudNativeAI/model-spec/tools/mdctl/registry"
	oci "github.com/opencontainers/image-spec/specs-go/v1"
)

func PushModel(name string) error {
	mp := models.ParseModelPath(name)

	manifestPath, err := mp.GetManifestPath()
	if err != nil {
		return fmt.Errorf("failed to get manifest path: %w", err)
	}

	manifestFile, err := os.Open(manifestPath)
	if err != nil {
		return fmt.Errorf("failed to open manifest file: %w", err)
	}
	defer manifestFile.Close()

	manifest := modelspec.Manifest{}
	if err := json.NewDecoder(manifestFile).Decode(&manifest); err != nil {
		return fmt.Errorf("failed to decode manifest: %w", err)
	}

	repo, err := registry.NewRepo(mp.Namespace, mp.Name)
	if err != nil {
		return fmt.Errorf("failed to new repo: %w", err)
	}
	ctx := context.Background()

	var layers []oci.Descriptor

	layerGroups := []struct {
		name   string
		layers []oci.Descriptor
	}{
		{"Description", manifest.Config.Description},
		{"License", manifest.Config.License},
		{"Extensions", manifest.Config.Extensions},
		{"Weights", manifest.Weights.File},
		{"Tokenizer", manifest.Processor},
	}

	for _, group := range layerGroups {
		if len(group.layers) == 0 {
			continue // Skip empty layer groups
		}
		layers = append(layers, group.layers...)

		for _, layer := range group.layers {
			fmt.Println("Push layer:", layer.Digest, layer.Size)
			_, err := registry.PushLayer(ctx, repo, &layer)
			if err != nil {
				return fmt.Errorf("failed to push layer: %w", err)
			}
		}
	}

	manifestDesc, err := registry.PushModelManifest(ctx, repo, manifestPath)
	if err != nil {
		return fmt.Errorf("failed to push model manifest: %w", err)
	}

	// push empty layer
	err = repo.Push(ctx, oci.DescriptorEmptyJSON, bytes.NewReader(oci.DescriptorEmptyJSON.Data))
	if err != nil {
		return fmt.Errorf("failed to push empty layer: %w", err)
	}

	// assemble descriptors and model manifest to a image manifest
	layers = append(layers, *manifestDesc)
	err = registry.PushModel(ctx, repo, mp.Tag, layers)
	if err != nil {
		return fmt.Errorf("failed to push oci image manifest: %w", err)
	}

	fmt.Println("Push succeed")
	return nil
}
