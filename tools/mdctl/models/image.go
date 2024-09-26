package models

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
)

type RootFS struct {
	Type    string   `json:"type"`
	DiffIDs []string `json:"diff_ids"`
}

// GetSHA256Digest returns the SHA256 hash of a given buffer and returns it, and the size of buffer
func GetSHA256Digest(r io.Reader) (string, int64) {
	h := sha256.New()
	n, err := io.Copy(h, r)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("sha256:%x", h.Sum(nil)), n
}

func verifyBlob(digest string) error {
	fp, err := GetBlobsPath(digest)
	if err != nil {
		return fmt.Errorf("failed to get blobs path: %w", err)
	}

	f, err := os.Open(fp)
	if err != nil {
		return fmt.Errorf("failed to open blob: %w", err)
	}
	defer f.Close()

	fileDigest, _ := GetSHA256Digest(f)
	if digest != fileDigest {
		return fmt.Errorf("digest mismatch: want %s, got %s", digest, fileDigest)
	}

	return nil
}
