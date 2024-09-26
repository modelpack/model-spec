package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/CloudNativeAI/model-spec/tools/mdctl/models"
)

func ListModel() error {
	dir, err := models.GetManifestRoot()
	if err != nil {
		return err
	}
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk model: %w", err)
		}
		if !info.IsDir() {
			term := dir + string(os.PathSeparator)
			name := strings.TrimPrefix(path, term)
			lastSeparatorIndex := strings.LastIndex(name, string(os.PathSeparator))
			if lastSeparatorIndex != -1 {
				name = name[:lastSeparatorIndex] + ":" + name[lastSeparatorIndex+1:]
			}
			fmt.Println(name)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to list model: %w", err)
	}
	return nil
}
