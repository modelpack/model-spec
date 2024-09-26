package models

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	modelspec "github.com/CloudNativeAI/model-spec/specs-go/v2"
)

func WriteManifest(name string, manifest *modelspec.Manifest) error {
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(manifest); err != nil {
		return fmt.Errorf("failed to encode manifest: %w", err)
	}

	modelPath := ParseModelPath(name)
	// modelPath := ParseModelPath("")
	manifestPath, err := modelPath.GetManifestPath()
	if err != nil {
		return fmt.Errorf("failed to get manifest path: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(manifestPath), 0755); err != nil {
		return fmt.Errorf("failed to mkdir all: %w", err)
	}

	return os.WriteFile(manifestPath, b.Bytes(), 0644)
}

func GetManifest(mp ModelPath) (*modelspec.Manifest, string, error) {
	fp, err := mp.GetManifestPath()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get manifest path: %w", err)
	}

	if _, err = os.Stat(fp); err != nil {
		return nil, "", fmt.Errorf("failed to stat manifest: %w", err)
	}

	var manifest *modelspec.Manifest

	bts, err := os.ReadFile(fp)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read manifest: %w", err)
	}

	shaSum := sha256.Sum256(bts)
	shaStr := hex.EncodeToString(shaSum[:])

	if err := json.Unmarshal(bts, &manifest); err != nil {
		return nil, "", fmt.Errorf("failed to unmarshal manifest: %w", err)
	}

	return manifest, shaStr, nil
}
