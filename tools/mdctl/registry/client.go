package registry

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	modelspec "github.com/CloudNativeAI/model-spec/specs-go/v2"
	"github.com/CloudNativeAI/model-spec/tools/mdctl/models"
	"github.com/opencontainers/image-spec/specs-go"
	"oras.land/oras-go/v2/content"

	oci "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
	"oras.land/oras-go/v2/registry/remote/retry"
)

const (
	ArtifactType = "application/vnd.test.artifact"
)

func NewRepo(ns, name string) (*remote.Repository, error) {
	user, exists := os.LookupEnv("MODEL_REGISTRY_USER")
	if !exists {
		return nil, fmt.Errorf("username not found")
	}
	password, exists := os.LookupEnv("MODEL_REGISTRY_PASSWORD")
	if !exists {
		return nil, fmt.Errorf("password not found")
	}
	reg, exists := os.LookupEnv("MODEL_REGISTRY_URL")
	if !exists {
		return nil, fmt.Errorf("registry url not found")
	}

	repo, err := remote.NewRepository(reg + "/" + ns + "/" + name)
	if err != nil {
		return nil, fmt.Errorf("failed to new repository: %w", err)
	}
	repo.PlainHTTP = true

	repo.Client = &auth.Client{
		Client: retry.DefaultClient,
		Cache:  auth.NewCache(),
		Credential: auth.StaticCredential(reg, auth.Credential{
			Username: user,
			Password: password,
		}),
	}

	return repo, nil
}

func NewOciManifest(layers []oci.Descriptor) ([]byte, error) {
	content := oci.Manifest{
		MediaType:    oci.MediaTypeImageManifest,
		ArtifactType: ArtifactType,
		Config:       oci.DescriptorEmptyJSON,
		Layers:       layers,
		Versioned:    specs.Versioned{SchemaVersion: 2},
	}
	return json.Marshal(content)
}

func PushLayer(repo *remote.Repository, ctx context.Context, descriptor *oci.Descriptor) (bool, error) {
	layerPath, err := models.GetBlobsPath(descriptor.Digest.String())
	if err != nil {
		return false, fmt.Errorf("failed to get blobs path: %w", err)
	}

	layerFile, err := os.Open(layerPath)
	if err != nil {
		return false, fmt.Errorf("failed to open layer file: %w", err)
	}
	defer layerFile.Close()

	exist, err := repo.Exists(ctx, *descriptor)
	if err != nil {
		return false, fmt.Errorf("failed to check if layer exists: %w", err)
	}
	if exist {
		return true, nil
	}

	return false, repo.Push(ctx, *descriptor, layerFile)
}

func PushModelManifest(repo *remote.Repository, ctx context.Context, modelManifestPath string) (*oci.Descriptor, error) {
	modelManifestFile, err := os.Open(modelManifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open model manifest file: %w", err)
	}
	defer modelManifestFile.Close()

	descriptor, err := models.FastDescriptor(modelManifestPath, modelspec.MediaTypeModelManifest)
	if err != nil {
		return nil, fmt.Errorf("failed to fast descriptor: %w", err)
	}

	fmt.Println("Push manifest:", descriptor.Digest, descriptor.Size)
	err = repo.Push(ctx, *descriptor, modelManifestFile)
	if err != nil {
		return nil, fmt.Errorf("failed to push model manifest: %w", err)
	}

	return descriptor, nil
}

func PushModel(repo *remote.Repository, tag string, ctx context.Context, layers []oci.Descriptor) error {
	manifestBlob, err := NewOciManifest(layers)
	if err != nil {
		return fmt.Errorf("failed to new oci manifest: %w", err)
	}
	manifestDesc := content.NewDescriptorFromBytes(oci.MediaTypeImageManifest, manifestBlob)

	err = repo.PushReference(ctx, manifestDesc, bytes.NewReader(manifestBlob), tag)
	if err != nil {
		return fmt.Errorf("failed to push model: %w", err)
	}

	return nil
}

func PullImageManifest(repo *remote.Repository, ctx context.Context, tag string) (*oci.Manifest, error) {
	descriptor, err := repo.Resolve(ctx, tag)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve image manifest: %w", err)
	}

	pulledBlob, err := content.FetchAll(ctx, repo, descriptor)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all: %w", err)
	}

	manifest := oci.Manifest{}
	if err := json.NewDecoder(bytes.NewReader(pulledBlob)).Decode(&manifest); err != nil {
		return nil, fmt.Errorf("failed to decode manifest: %w", err)
	}

	return &manifest, nil
}

func PullLayer(repo *remote.Repository, ctx context.Context, digest string, size int64, targetPath string) error {
	if fi, err := os.Stat(targetPath); err == nil && fi.Mode().IsRegular() && fi.Size() == size {
		return nil
	}

	descriptor, err := repo.Blobs().Resolve(ctx, digest)
	if err != nil {
		return fmt.Errorf("failed to resolve blob: %w", err)
	}

	rc, err := repo.Fetch(ctx, descriptor)
	if err != nil {
		return fmt.Errorf("failed to fetch blob: %w", err)
	}
	defer rc.Close()

	pulledBlob, err := content.ReadAll(rc, descriptor)
	if err != nil {
		return fmt.Errorf("failed to read all: %w", err)
	}

	blobs, err := models.GetBlobsPath("")
	if err != nil {
		return fmt.Errorf("failed to get blobs path: %w", err)
	}

	delimiter := ":"
	pattern := strings.Join([]string{"sha256", "*-downloading"}, delimiter)
	temp, err := os.CreateTemp(blobs, pattern)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer temp.Close()

	_, err = io.Copy(temp, bytes.NewReader(pulledBlob))
	if err != nil {
		return fmt.Errorf("failed to copy blob to temp file: %w", err)
	}

	err = os.Rename(temp.Name(), targetPath)
	if err != nil {
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}
