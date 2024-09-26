package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	spec "github.com/CloudNativeAI/model-spec/specs-go/v2"
	"github.com/CloudNativeAI/model-spec/tools/mdctl/models"
	"github.com/CloudNativeAI/model-spec/tools/mdctl/registry"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

func FetchManifest(name string, manifest *spec.Manifest, config *spec.Config) (*v1.Manifest, error) {
	mp := models.ParseModelPath(name)
	repo, err := registry.NewRepo(mp.Namespace, mp.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to new repo: %w", err)
	}
	ctx := context.Background()
	imageManifest, err := registry.PullImageManifest(repo, ctx, mp.Tag)
	if err != nil {
		return nil, fmt.Errorf("failed to pull image manifest: %w", err)
	}

	// Fetch layers
	for _, layer := range imageManifest.Layers {
		switch layer.MediaType {
		case spec.MediaTypeModelManifest, spec.MediaTypeModelConfig:
		default:
			continue
		}

		// create temp file
		// TODO: use []byte
		tempRoot, err := models.GetBlobsPath("")
		if err != nil {
			return nil, fmt.Errorf("failed to get blobs path: %w", err)
		}
		delimiter := ":"
		pattern := strings.Join([]string{"sha256", "*-temp"}, delimiter)
		temp, err := os.CreateTemp(tempRoot, pattern)
		if err != nil {
			return nil, fmt.Errorf("failed to create temp file: %w", err)
		}
		defer temp.Close()

		//fmt.Println("Pull layer: ", layer.Digest, layer.Size)
		err = registry.PullLayer(repo, ctx, layer.Digest.String(), layer.Size, temp.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to pull layer: %w", err)
		}
		content, err := os.ReadFile(temp.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to read temp file: %w", err)
		}

		switch layer.MediaType {
		case spec.MediaTypeModelManifest:
			if err := json.Unmarshal(content, &manifest); err != nil {
				return nil, fmt.Errorf("failed to unmarshal manifest: %w", err)
			}
		case spec.MediaTypeModelConfig:
			if err := json.Unmarshal(content, &config); err != nil {
				return nil, fmt.Errorf("failed to unmarshal config: %w", err)
			}
		}
	}
	return imageManifest, nil
}

func PullModel(name string) error {
	mp := models.ParseModelPath(name)

	repo, err := registry.NewRepo(mp.Namespace, mp.Name)
	if err != nil {
		return fmt.Errorf("failed to new repo: %w", err)
	}
	ctx := context.Background()

	image_manifest, err := registry.PullImageManifest(repo, ctx, mp.Tag)
	if err != nil {
		return fmt.Errorf("failed to pull image manifest: %w", err)
	}

	for _, layer := range image_manifest.Layers {
		fmt.Println("Pull layer:", layer.Digest, layer.Size)
		digest := layer.Digest.String()

		var targetPath string
		if layer.MediaType == spec.MediaTypeModelManifest {
			targetPath, err = mp.GetManifestPath()
			targetDir := filepath.Dir(targetPath)
			if err := os.MkdirAll(targetDir, 0o755); err != nil {
				return fmt.Errorf("failed to mkdir: %w", err)
			}
		} else {
			targetPath, err = models.GetBlobsPath(digest)
		}
		if err != nil {
			return fmt.Errorf("failed to get blobs path: %w", err)
		}

		err = registry.PullLayer(repo, ctx, digest, layer.Size, targetPath)
		if err != nil {
			return fmt.Errorf("failed to pull layer: %w", err)
		}
	}
	fmt.Println("Pull succeed")

	return nil
}
