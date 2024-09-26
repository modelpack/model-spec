package models

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/opencontainers/go-digest"
	oci "github.com/opencontainers/image-spec/specs-go/v1"
)

const (
	TAR  = 1
	ZSTD = 2
)

func FastDescriptor(srcPath string, media string) (*oci.Descriptor, error) {
	var err error
	srcPath, err = filepath.Abs(srcPath)
	if err != nil {
		return nil, fmt.Errorf("get real path failed: %w", err)
	}

	bin, err := os.Open(srcPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open source file: %w", err)
	}
	defer bin.Close()

	sha256sum := sha256.New()
	n, err := io.Copy(sha256sum, bin)
	if err != nil {
		return nil, fmt.Errorf("failed to read for sha256: %w", err)
	}

	return &oci.Descriptor{
		MediaType: media,
		Digest:    digest.Digest(fmt.Sprintf("sha256:%x", sha256sum.Sum(nil))),
		Size:      n,
		Annotations: map[string]string{
			"temp_name": srcPath,
		},
	}, nil
}

func BuildDescriptor(method int, src, media, newName string) (*oci.Descriptor, error) {
	var err error
	src, err = filepath.Abs(src)
	if err != nil {
		return nil, fmt.Errorf("get real path failed: %w", err)
	}

	blobs, err := GetBlobsPath("")
	if err != nil {
		return nil, fmt.Errorf("failed to get blobs path: %w", err)
	}

	delimiter := ":"
	pattern := strings.Join([]string{"sha256", "*-partial"}, delimiter)
	temp, err := os.CreateTemp(blobs, pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer temp.Close()

	switch method {
	case TAR:
		err = Tar(src, temp.Name(), newName)
	case ZSTD:
		err = Compress(src, temp.Name())
	default:
		err = Tar(src, temp.Name(), newName)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to build descriptor: %w", err)
	}

	bin, err := os.Open(temp.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to open source file: %w", err)
	}
	defer bin.Close()

	sha256sum := sha256.New()
	n, err := io.Copy(sha256sum, bin)
	if err != nil {
		return nil, fmt.Errorf("failed to read for sha256: %w", err)
	}

	return &oci.Descriptor{
		MediaType: media,
		Digest:    digest.Digest(fmt.Sprintf("sha256:%x", sha256sum.Sum(nil))),
		Size:      n,
		Annotations: map[string]string{
			"temp_name": temp.Name(),
		},
	}, nil
}

func NewDescriptor(r io.Reader, mediatype string) (*oci.Descriptor, error) {
	blobs, err := GetBlobsPath("")
	if err != nil {
		return nil, fmt.Errorf("failed to get blobs path: %w", err)
	}

	delimiter := ":"
	pattern := strings.Join([]string{"sha256", "*-partial"}, delimiter)
	temp, err := os.CreateTemp(blobs, pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer temp.Close()

	sha256sum := sha256.New()
	n, err := io.Copy(io.MultiWriter(temp, sha256sum), r)
	if err != nil {
		return nil, fmt.Errorf("failed to read for sha256: %w", err)
	}

	return &oci.Descriptor{
		MediaType: mediatype,
		//Digest:    fmt.Sprintf("sha256:%x", sha256sum.Sum(nil)),
		Digest: digest.Digest(fmt.Sprintf("sha256:%x", sha256sum.Sum(nil))),
		Size:   n,
		Annotations: map[string]string{
			"temp_name": temp.Name(),
		},
	}, nil
}

func Commit(l oci.Descriptor) (bool, error) {
	tempFileName := l.Annotations["temp_name"]
	if tempFileName == "" {
		return false, fmt.Errorf("temp file name is empty")
	}

	// always remove temp
	defer os.Remove(tempFileName)
	defer delete(l.Annotations, "temp_name")

	blob, err := GetBlobsPath(l.Digest.String())
	if err != nil {
		return false, fmt.Errorf("failed to get blobs path: %w", err)
	}

	if _, err := os.Stat(blob); err != nil {
		return true, os.Rename(tempFileName, blob)
	}

	return false, nil
}
