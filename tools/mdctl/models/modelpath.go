package models

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const (
	DefaultProtocolScheme = "https"
	DefaultRegistry       = "registry.cnai.com"
	DefaultNamespace      = "sys"
	DefaultTag            = "latest"
)

var (
	ErrInvalidImageFormat = errors.New("invalid models format")
	ErrInvalidProtocol    = errors.New("invalid protocol scheme")
	ErrInsecureProtocol   = errors.New("insecure protocol http")
)

var errModelPathInvalid = errors.New("invalid models path")

// func realpath(mfDir, from string) string {
// 	abspath, err := filepath.Abs(from)
// 	if err != nil {
// 		return from
// 	}

// 	home, err := os.UserHomeDir()
// 	if err != nil {
// 		return abspath
// 	}

// 	if from == "~" {
// 		return home
// 	} else if strings.HasPrefix(from, "~/") {
// 		return filepath.Join(home, from[2:])
// 	}

// 	if _, err := os.Stat(filepath.Join(mfDir, from)); err == nil {
// 		// this is a file relative to the Modelfile
// 		return filepath.Join(mfDir, from)
// 	}

// 	return abspath
// }

func ParseModelPath(name string) ModelPath {
	mp := ModelPath{
		ProtocolScheme: DefaultProtocolScheme,
		Registry:       DefaultRegistry,
		Namespace:      DefaultNamespace,
		Name:           "",
		Tag:            DefaultTag,
	}

	before, after, found := strings.Cut(name, "://")
	if found {
		mp.ProtocolScheme = before
		name = after
	}

	parts := strings.Split(name, string(os.PathSeparator))
	switch len(parts) {
	case 3:
		mp.Registry = parts[0]
		mp.Namespace = parts[1]
		mp.Name = parts[2]
	case 2:
		mp.Namespace = parts[0]
		mp.Name = parts[1]
	case 1:
		mp.Name = parts[0]
	}

	if repo, tag, found := strings.Cut(mp.Name, ":"); found {
		mp.Name = repo
		mp.Tag = tag
	}

	return mp
}

// ModelDir returns the path to the models directory.
func ModelDir() (string, error) {
	if models, exists := os.LookupEnv("MODELS_DIR"); exists {
		return models, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home dir: %w", err)
	}
	return filepath.Join(home, ".models"), nil
}

func GetManifestRoot() (string, error) {
	dir, err := ModelDir()
	if err != nil {
		return "", fmt.Errorf("failed to get models dir: %w", err)
	}

	path := filepath.Join(dir, "manifests")
	if err := os.MkdirAll(path, 0o755); err != nil {
		return "", fmt.Errorf("failed to mkdir all: %w", err)
	}

	return path, nil
}

func GetBlobsPath(digest string) (string, error) {
	dir, err := ModelDir()
	if err != nil {
		return "", fmt.Errorf("failed to get models dir: %w", err)
	}

	path := filepath.Join(dir, "blobs", digest)
	dirPath := filepath.Dir(path)
	if digest == "" {
		dirPath = path
	}

	if err := os.MkdirAll(dirPath, 0o755); err != nil {
		return "", fmt.Errorf("failed to mkdir all: %w", err)
	}

	return path, nil
}

type ModelPath struct {
	ProtocolScheme string
	Registry       string
	Namespace      string
	Name           string
	Tag            string
}

func (mp ModelPath) Validate() error {
	if mp.Name == "" {
		return fmt.Errorf("%w: models repository Name is required", errModelPathInvalid)
	}

	if strings.Contains(mp.Tag, ":") {
		return fmt.Errorf("%w: ':' (colon) is not allowed in tag names", errModelPathInvalid)
	}

	return nil
}

func (mp ModelPath) GetNamespaceRepository() string {
	return fmt.Sprintf("%s/%s", mp.Namespace, mp.Name)
}

func (mp ModelPath) GetFullTagname() string {
	return fmt.Sprintf("%s/%s/%s:%s", mp.Registry, mp.Namespace, mp.Name, mp.Tag)
}

func (mp ModelPath) GetShortTagname() string {
	if mp.Registry == DefaultRegistry {
		if mp.Namespace == DefaultNamespace {
			return fmt.Sprintf("%s:%s", mp.Name, mp.Tag)
		}
		return fmt.Sprintf("%s/%s:%s", mp.Namespace, mp.Name, mp.Tag)
	}
	return fmt.Sprintf("%s/%s/%s:%s", mp.Registry, mp.Namespace, mp.Name, mp.Tag)
}

// GetManifestRoot returns the path to the manifest file for the given models path,
// it is up to the caller to create the directory if it does not exist.
func (mp ModelPath) GetManifestPath() (string, error) {
	dir, err := ModelDir()
	if err != nil {
		return "", fmt.Errorf("failed to get models dir: %w", err)
	}

	return filepath.Join(dir, "manifests", mp.Registry, mp.Namespace, mp.Name, mp.Tag), nil
}

func (mp ModelPath) BaseURL() *url.URL {
	return &url.URL{
		Scheme: mp.ProtocolScheme,
		Host:   mp.Registry,
	}
}
