/*
 *     Copyright 2025 The CNCF ModelPack Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1

import (
	"encoding/json"
	"testing"
	"time"
)

func TestAnnotationConstants(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{"AnnotationFilepath", AnnotationFilepath, "org.cncf.model.filepath"},
		{"AnnotationFileMetadata", AnnotationFileMetadata, "org.cncf.model.file.metadata+json"},
		{"AnnotationMediaTypeUntested", AnnotationMediaTypeUntested, "org.cncf.model.file.mediatype.untested"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != tt.want {
				t.Errorf("%s = %q, want %q", tt.name, tt.value, tt.want)
			}
		})
	}
}

func TestFileMetadataMarshalJSON(t *testing.T) {
	mtime := time.Date(2025, 3, 15, 14, 30, 0, 0, time.UTC)
	fm := FileMetadata{
		Name:     "model.safetensors",
		Mode:     0644,
		Uid:      1000,
		Gid:      1000,
		Size:     1024000,
		ModTime:  mtime,
		Typeflag: '0',
	}

	data, err := json.Marshal(fm)
	if err != nil {
		t.Fatalf("failed to marshal FileMetadata: %v", err)
	}

	var got FileMetadata
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("failed to unmarshal FileMetadata: %v", err)
	}

	if got.Name != fm.Name {
		t.Errorf("name = %q, want %q", got.Name, fm.Name)
	}
	if got.Mode != fm.Mode {
		t.Errorf("mode = %d, want %d", got.Mode, fm.Mode)
	}
	if got.Uid != fm.Uid {
		t.Errorf("uid = %d, want %d", got.Uid, fm.Uid)
	}
	if got.Gid != fm.Gid {
		t.Errorf("gid = %d, want %d", got.Gid, fm.Gid)
	}
	if got.Size != fm.Size {
		t.Errorf("size = %d, want %d", got.Size, fm.Size)
	}
	if !got.ModTime.Equal(mtime) {
		t.Errorf("mtime = %v, want %v", got.ModTime, mtime)
	}
	if got.Typeflag != fm.Typeflag {
		t.Errorf("typeflag = %d, want %d", got.Typeflag, fm.Typeflag)
	}
}

func TestFileMetadataJSONFieldNames(t *testing.T) {
	fm := FileMetadata{
		Name:     "test.bin",
		Mode:     0755,
		Uid:      0,
		Gid:      0,
		Size:     100,
		ModTime:  time.Now(),
		Typeflag: '0',
	}

	data, err := json.Marshal(fm)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("failed to unmarshal to map: %v", err)
	}

	expectedKeys := []string{"name", "mode", "uid", "gid", "size", "mtime", "typeflag"}
	for _, key := range expectedKeys {
		if _, ok := raw[key]; !ok {
			t.Errorf("expected JSON key %q in FileMetadata, not found", key)
		}
	}
}

func TestFileMetadataRoundTrip(t *testing.T) {
	original := FileMetadata{
		Name:     "weights/layer1.bin",
		Mode:     0644,
		Uid:      1000,
		Gid:      1000,
		Size:     5242880,
		ModTime:  time.Date(2025, 2, 1, 8, 0, 0, 0, time.UTC),
		Typeflag: '0',
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var roundTripped FileMetadata
	if err := json.Unmarshal(data, &roundTripped); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	data2, err := json.Marshal(roundTripped)
	if err != nil {
		t.Fatalf("failed to re-marshal: %v", err)
	}

	if string(data) != string(data2) {
		t.Errorf("round-trip JSON mismatch:\n  first:  %s\n  second: %s", data, data2)
	}
}
