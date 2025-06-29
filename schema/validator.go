/*
 *     Copyright 2025 The CNAI Authors
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

package schema

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	v1 "github.com/modelpack/model-spec/specs-go/v1"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

// Validator wraps a media type string identifier and implements validation against a JSON schema.
type Validator string

// Validate validates the given reader against the schema of the wrapped media type.
func (v Validator) Validate(src io.Reader) error {
	// run the media type specific validation
	if fn, ok := validateByMediaType[v]; ok {
		if fn == nil {
			return fmt.Errorf("internal error: mapValidate is nil for %s", string(v))
		}
		// buffer the src so the media type validation and the schema validation can both read it
		buf, err := io.ReadAll(src)
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		src = bytes.NewReader(buf)
		err = fn(buf)
		if err != nil {
			return err
		}
	}

	// json schema validation
	return v.validateSchema(src)
}

func (v Validator) validateSchema(src io.Reader) error {
	if _, ok := specs[v]; !ok {
		return fmt.Errorf("no validator available for %s", string(v))
	}

	c := jsonschema.NewCompiler()

	// load the schema files from the embedded FS
	dir, err := specFS.ReadDir(".")
	if err != nil {
		return fmt.Errorf("spec embedded directory could not be loaded: %w", err)
	}
	for _, file := range dir {
		if file.IsDir() {
			continue
		}
		specBuf, err := specFS.ReadFile(file.Name())
		if err != nil {
			return fmt.Errorf("could not read spec file %s: %w", file.Name(), err)
		}
		err = c.AddResource(file.Name(), bytes.NewReader(specBuf))
		if err != nil {
			return fmt.Errorf("failed to add spec file %s: %w", file.Name(), err)
		}
		if len(specURLs[file.Name()]) == 0 {
			// this would be a bug in the validation code itself, add any missing entry to schema.go
			return fmt.Errorf("spec file has no aliases: %s", file.Name())
		}
		for _, specURL := range specURLs[file.Name()] {
			err = c.AddResource(specURL, bytes.NewReader(specBuf))
			if err != nil {
				return fmt.Errorf("failed to add spec file %s as url %s: %w", file.Name(), specURL, err)
			}
		}
	}

	// compile based on the type of validator
	schema, err := c.Compile(specs[v])
	if err != nil {
		return fmt.Errorf("failed to compile schema %s: %w", string(v), err)
	}

	// read in the user input and validate
	var input interface{}
	err = json.NewDecoder(src).Decode(&input)
	if err != nil {
		return fmt.Errorf("unable to parse json to validate: %w", err)
	}
	err = schema.Validate(input)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	return nil
}

type validateFunc func([]byte) error

var validateByMediaType = map[Validator]validateFunc{
	ValidatorMediaTypeModelConfig: validateConfig,
}

func validateConfig(buf []byte) error {
	mc := v1.ModelConfig{}

	err := json.Unmarshal(buf, &mc)
	if err != nil {
		return fmt.Errorf("config format mismatch: %w", err)
	}

	return nil
}
