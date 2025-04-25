# Annotations

This property contains arbitrary metadata, and SHOULD follow the rules of [OCI image annotations](https://github.com/opencontainers/image-spec/blob/main/annotations.md).

## Pre-defined Annotation Keys

### Layer Annotation Keys

- **`org.cnai.model.filepath`**: Specifies the file path of the layer (string).

- **`org.cnai.model.file.metadata+json`**: Specifies the metadata of the file (string), value is the JSON string of [File Metadata Annotation Value](#File-Metadata-Annotation-Value).

### Layer Annotation Values

#### File Metadata Annotation Value

```go
// FileMetadata represents the metadata of file, which is the value definition of AnnotationFileMetadata.
type FileMetadata struct {
	// File name
	Name string `json:"name"`

	// File permission mode (e.g., Unix permission bits)
	Mode uint32 `json:"mode"`

	// User ID (identifier of the file owner)
	Uid uint32 `json:"uid"`

	// Group ID (identifier of the file's group)
	Gid uint32 `json:"gid"`

	// File size (in bytes)
	Size int64 `json:"size"`

	// File last modification time
	ModTime time.Time `json:"mtime"`

	// File type flag (e.g., regular file, directory, etc.)
	Typeflag byte `json:"typeflag"`
}
```
