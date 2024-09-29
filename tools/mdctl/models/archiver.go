package models

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/klauspost/compress/zstd"
)

func Tar(src, dst, newName string) error {
	fi, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to stat src: %w", err)
	}

	if fi.IsDir() {
		return tarDirectory(src, dst)
	}

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create tar file: %w", err)
	}
	defer out.Close()
	return tarFile(src, newName, out)
}

func tarFile(src, newName string, writers ...io.Writer) error {
	mw := io.MultiWriter(writers...)
	tw := tar.NewWriter(mw)
	defer tw.Close()

	file, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open src: %w", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat src: %w", err)
	}

	// todo: check the link arg?
	header, err := tar.FileInfoHeader(stat, stat.Name())
	if err != nil {
		return fmt.Errorf("failed to get file info header: %w", err)
	}
	if newName != "" {
		header.Name = newName
	}

	if err := tw.WriteHeader(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	if _, err := io.Copy(tw, file); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}

func tarDirectory(srcPath, tarPath string) error {
	tarFile, err := os.Create(tarPath)
	if err != nil {
		return fmt.Errorf("failed to create tar file: %w", err)
	}
	defer tarFile.Close()

	tw := tar.NewWriter(tarFile)
	defer tw.Close()

	var tarWalkFn filepath.WalkFunc = func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk file: %w", err)
		}
		relPath, err := filepath.Rel(srcPath, file)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}
		if relPath == "." {
			// Skip the root directory entry
			return nil
		}
		header, err := tar.FileInfoHeader(fi, relPath)
		if err != nil {
			return fmt.Errorf("failed to get file info header: %w", err)
		}
		header.Name = relPath
		if err := tw.WriteHeader(header); err != nil {
			return fmt.Errorf("failed to write header: %w", err)
		}
		if fi.IsDir() {
			return nil
		}
		data, err := os.Open(file)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer data.Close()
		if _, err := io.Copy(tw, data); err != nil {
			return fmt.Errorf("failed to copy file: %w", err)
		}
		return nil
	}

	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		return fmt.Errorf("failed to stat src: %w", err)
	}
	if srcInfo.IsDir() {
		// Walk the source path and tar each file and directory
		if err := filepath.Walk(srcPath, tarWalkFn); err != nil {
			return fmt.Errorf("failed to walk file: %w", err)
		}
	}

	return nil
}

func compressFile(srcPath, dstPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	encoder, err := zstd.NewWriter(dstFile)
	if err != nil {
		return fmt.Errorf("failed to create zstd encoder: %w", err)
	}
	defer encoder.Close()

	if _, err := io.Copy(encoder, srcFile); err != nil {
		return fmt.Errorf("failed to compress file: %w", err)
	}

	return nil
}

func decompressFile(srcPath, dstPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	decoder, err := zstd.NewReader(srcFile)
	if err != nil {
		return fmt.Errorf("failed to create zstd decoder: %w", err)
	}
	defer decoder.Close()

	if _, err := io.Copy(dstFile, decoder); err != nil {
		return fmt.Errorf("failed to decompress file: %w", err)
	}

	return nil
}

func untar(tarPath, dstPath string) error {
	fileName := filepath.Base(tarPath)
	fmt.Printf("Unpack layer: %s\n", fileName)

	tarFile, err := os.Open(tarPath)
	if err != nil {
		return fmt.Errorf("failed to open tar file: %w", err)
	}
	defer tarFile.Close()

	tr := tar.NewReader(tarFile)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar file: %w", err)
		}

		target := filepath.Join(dstPath, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
		case tar.TypeReg:
			file, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}
			if _, err := io.Copy(file, tr); err != nil {
				_ = file.Close()
				return fmt.Errorf("failed to write file: %w", err)
			}
			_ = file.Close()
		}
	}

	return nil
}

//func Tar(srcDir, destFilePath string) error {
//	return tarit(srcDir, destFilePath)
//}

func UnTar(srcFilePath, destDir string) error {
	return untar(srcFilePath, destDir)
}

func Compress(src, destFilePath string) error {
	tempTarPath := destFilePath + ".tar"
	if err := Tar(src, tempTarPath, ""); err != nil {
		return fmt.Errorf("failed to tar: %w", err)
	}
	defer os.Remove(tempTarPath)
	if err := compressFile(tempTarPath, destFilePath); err != nil {
		return fmt.Errorf("failed to compress file: %w", err)
	}

	return nil
}

func Decompress(srcFilePath, destDir string) error {
	tempTarPath := srcFilePath + ".tar"
	if err := decompressFile(srcFilePath, tempTarPath); err != nil {
		return fmt.Errorf("failed to decompress file: %w", err)
	}
	defer os.Remove(tempTarPath)

	if err := untar(tempTarPath, destDir); err != nil {
		return fmt.Errorf("failed to untar: %w", err)
	}

	return nil
}
