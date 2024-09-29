package models

import (
	"errors"
	"fmt"
	"os"

	oci "github.com/opencontainers/image-spec/specs-go/v1"
)

type Descriptors struct {
	Items []oci.Descriptor
}

func (ls *Descriptors) Add(layer *oci.Descriptor) {
	if layer.Size > 0 {
		ls.Items = append(ls.Items, *layer)
	}
}

func (ls *Descriptors) Replace(layer *oci.Descriptor) {
	if layer.Size > 0 {
		var newItems []oci.Descriptor
		for _, item := range ls.Items {
			if item.MediaType != layer.MediaType {
				newItems = append(newItems, item)
			}
		}
		ls.Items = append(newItems, *layer)
	}
}

func (ls *Descriptors) Delete(layer *oci.Descriptor) {
	var newItems []oci.Descriptor
	for _, item := range ls.Items {
		if item.MediaType != layer.MediaType {
			newItems = append(newItems, item)
		}
	}
	ls.Items = newItems
}

func (ls *Descriptors) AddFile(input string, media string) (*oci.Descriptor, error) {
	var err error
	fileInfo, err := os.Stat(input)
	if err != nil {
		return nil, fmt.Errorf("failed to stat input: %w", err)
	}
	if fileInfo.Mode().IsRegular() {
		bin, err := os.Open(input)
		if err != nil {
			return nil, fmt.Errorf("failed to open input: %w", err)
		}
		defer bin.Close()
		layer, err := NewDescriptor(bin, media)
		if err != nil {
			return nil, fmt.Errorf("failed to build descriptor: %w", err)
		}
		ls.Add(layer)
		return layer, nil
	}
	return nil, errors.New("not a regular file type")
}

func (ls *Descriptors) AddCompress(srcDir string, media string) (*oci.Descriptor, error) {
	layer, err := BuildDescriptor(ZSTD, srcDir, media, "")
	if err != nil {
		return nil, fmt.Errorf("failed to build descriptor: %w", err)
	}
	//fmt.Printf("add layer: %v\n", layer)
	ls.Add(layer)
	return layer, nil
}

func (ls *Descriptors) ReplaceCompress(srcDir string, media string) (*oci.Descriptor, error) {
	layer, err := BuildDescriptor(ZSTD, srcDir, media, "")
	if err != nil {
		return nil, fmt.Errorf("failed to build descriptor: %w", err)
	}
	//fmt.Printf("add layer: %v\n", layer)
	ls.Replace(layer)
	return layer, nil
}

func (ls *Descriptors) AddTar(srcDir string, media string) (*oci.Descriptor, error) {
	layer, err := BuildDescriptor(TAR, srcDir, media, "")
	if err != nil {
		return nil, fmt.Errorf("failed to build descriptor: %w", err)
	}
	ls.Add(layer)
	return layer, nil
}

func (ls *Descriptors) ReplaceTar(srcDir string, media string) (*oci.Descriptor, error) {
	layer, err := BuildDescriptor(TAR, srcDir, media, "")
	if err != nil {
		return nil, fmt.Errorf("failed to build descriptor: %w", err)
	}
	ls.Replace(layer)
	return layer, nil
}

func (ls *Descriptors) AddTarWithNewName(src, media, newName string) (*oci.Descriptor, error) {
	layer, err := BuildDescriptor(TAR, src, media, newName)
	if err != nil {
		return nil, fmt.Errorf("failed to build descriptor: %w", err)
	}
	ls.Add(layer)
	return layer, nil
}

func (ls *Descriptors) ReplaceTarWithNewName(src, media, newName string) (*oci.Descriptor, error) {
	layer, err := BuildDescriptor(TAR, src, media, newName)
	if err != nil {
		return nil, fmt.Errorf("failed to build descriptor: %w", err)
	}
	ls.Replace(layer)
	return layer, nil
}

//
//func (ls *Descriptors) CreateDescriptor(input string, mediatype string) (*spec.BuildDescriptor, error) {
//	var layer *spec.BuildDescriptor
//	var err error
//	fileInfo, err := os.Stat(input)
//	if err != nil { // Check error type
//		bin := strings.NewReader(input)
//		layer, err = NewDescriptor(bin, mediatype)
//		if err != nil {
//			return nil, err
//		}
//		ls.Add(layer)
//		return layer, nil
//	}
//	if fileInfo.Mode().IsRegular() {
//		bin, err := os.Open(input)
//		if err != nil {
//			return nil, err
//		}
//		layer, err := NewDescriptor(bin, mediatype)
//		if err != nil {
//			bin.Close()
//			return nil, err
//		}
//		ls.Add(layer)
//		bin.Close()
//	}
//	return layer, nil
//}

func (ls *Descriptors) Commit() error {
	// Commit every layer
	for _, layer := range ls.Items {
		committed, err := Commit(layer)
		if err != nil {
			return fmt.Errorf("failed to commit layer: %w", err)
		}
		status := "writing layer"
		if !committed {
			status = "layer already exists"
		}
		fmt.Printf("%s %s\n", status, layer.Digest)
	}
	return nil
}
