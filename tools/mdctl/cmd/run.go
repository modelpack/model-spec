package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/CloudNativeAI/model-spec/tools/mdctl/models"
)

const (
	DOT_GITS_DIR    = ".gits"
	DOT_VOLUMES_DIR = ".volumes"
	MODEL_DIR       = "model"
	DATASET_DIR     = "dataset"
	SOURCE_DIR      = "source"
	TASK_DIR        = "task"
	ENTRYPOINT      = "run.py"
	SETUP           = "setup.sh"
	CONFIG          = "config.json"
	INFO            = "info.json"
	LICENSE         = "LICENSE"
)

func RunModel(name string) error {
	mp := models.ParseModelPath(name)
	manifest, _, err := models.GetManifest(mp)
	if err != nil {
		return fmt.Errorf("failed to get manifest: %w", err)
	}

	root, err := filepath.Abs(mp.Name + ":" + mp.Tag)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	if err := os.MkdirAll(root, 0o755); err != nil {
		return fmt.Errorf("failed to mkdir: %w", err)
	}

	for _, layer := range manifest.Weights.File {
		filename, err := models.GetBlobsPath(layer.Digest.String())
		if err != nil {
			return fmt.Errorf("failed to get blobs path: %w", err)
		}
		if err := models.UnTar(filename, root); err != nil {
			return fmt.Errorf("failed to untar: %w", err)
		}
	}

	for _, layer := range manifest.Processor {
		filename, err := models.GetBlobsPath(layer.Digest.String())
		if err != nil {
			return fmt.Errorf("failed to get blobs path: %w", err)
		}
		if err := models.UnTar(filename, root); err != nil {
			return fmt.Errorf("failed to untar: %w", err)
		}
	}

	for _, layer := range manifest.Config.Description {
		filename, err := models.GetBlobsPath(layer.Digest.String())
		if err != nil {
			return fmt.Errorf("failed to get blobs path: %w", err)
		}
		if err := models.UnTar(filename, root); err != nil {
			return fmt.Errorf("failed to untar: %w", err)
		}
	}

	for _, layer := range manifest.Config.License {
		filename, err := models.GetBlobsPath(layer.Digest.String())
		if err != nil {
			return fmt.Errorf("failed to get blobs path: %w", err)
		}
		if err := models.UnTar(filename, root); err != nil {
			return fmt.Errorf("failed to untar: %w", err)
		}
	}

	for _, layer := range manifest.Config.Extensions {
		filename, err := models.GetBlobsPath(layer.Digest.String())
		if err != nil {
			return fmt.Errorf("failed to get blobs path: %w", err)
		}
		if err := models.UnTar(filename, root); err != nil {
			return fmt.Errorf("failed to untar: %w", err)
		}
	}

	err = os.Chdir(root)
	if err != nil {
		return fmt.Errorf("failed to change workdir: %w", err)
	}

	// stdout, stderr, err = executeScript(entrypoint, []string{})
	// if err != nil {
	// 	fmt.Printf("Error: %v\n", err)
	// 	if stderr != "" {
	// 		fmt.Printf("Stderr: %v\n", stderr)
	// 	}
	// 	return err
	// }
	// fmt.Printf("Stdout: %v\n", stdout)

	return nil
}

func executeBinary(binaryPath string, args []string) (stdout string, stderr string, err error) {
	cmd := exec.Command(binaryPath, args...)

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err = cmd.Run()
	if err != nil {
		return "", "", fmt.Errorf("failed to execute binary: %w", err)
	}

	stdout = outBuf.String()
	stderr = errBuf.String()

	return stdout, stderr, nil
}

func executeScript(scriptPath string, args []string) (stdout string, stderr string, err error) {
	var cmd *exec.Cmd
	if bytes.HasSuffix([]byte(scriptPath), []byte(".sh")) {
		cmd = exec.Command("bash", append([]string{scriptPath}, args...)...)
	} else if bytes.HasSuffix([]byte(scriptPath), []byte(".py")) {
		cmd = exec.Command("python3", append([]string{scriptPath}, args...)...)
	} else {
		return "", "", fmt.Errorf("unsupported script type: %s", scriptPath)
	}

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err = cmd.Run()
	if err != nil {
		return "", "", fmt.Errorf("failed to execute script: %w", err)
	}

	stdout = outBuf.String()
	stderr = errBuf.String()

	return stdout, stderr, nil
}
