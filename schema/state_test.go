/*
 *     Copyright 2026 The CNCF ModelPack Authors
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

package schema_test

import (
	"strings"
	"testing"

	"github.com/modelpack/model-spec/schema"
)

func TestState(t *testing.T) {
	for i, tt := range []struct {
		state string
		fail  bool
	}{
		// valid: state with all fields
		{
			state: `{"format": "llama.cpp", "tokens": [1, 887, 526], "precision": "float16"}`,
			fail:  false,
		},
		// valid: state without precision
		{
			state: `{"format": "vllm", "tokens": [1, 887, 526]}`,
			fail:  false,
		},
		// expected failure: state is null
		{
			state: `null`,
			fail:  true,
		},
		// expected failure: format is missing
		{
			state: `{"tokens": [1, 887, 526]}`,
			fail:  true,
		},
		// expected failure: format is empty
		{
			state: `{"format": "", "tokens": [1, 887, 526]}`,
			fail:  true,
		},
		// expected failure: format is null
		{
			state: `{"format": null, "tokens": [1, 887, 526]}`,
			fail:  true,
		},
		// expected failure: tokens is missing
		{
			state: `{"format": "vllm"}`,
			fail:  true,
		},
		// expected failure: tokens is null
		{
			state: `{"format": "vllm", "tokens": null}`,
			fail:  true,
		},
		// expected failure: token ID is null
		{
			state: `{"format": "vllm", "tokens": [1, null, 526]}`,
			fail:  true,
		},
		// expected failure: precision is null
		{
			state: `{"format": "vllm", "tokens": [1, 887, 526], "precision": null}`,
			fail:  true,
		},
		// expected failure: token ID is negative
		{
			state: `{"format": "vllm", "tokens": [1, -887, 526]}`,
			fail:  true,
		},
		// expected failure: token ID is not an integer
		{
			state: `{"format": "vllm", "tokens": [1, 887.5, 526]}`,
			fail:  true,
		},
		// expected failure: unknown property
		{
			state: `{"format": "vllm", "tokens": [1, 887, 526], "foo": "bar"}`,
			fail:  true,
		},
		// expected failure: trailing garbage
		{
			state: `{"format": "vllm", "tokens": [1, 887, 526]}garbage`,
			fail:  true,
		},
		// expected failure: a second JSON value
		{
			state: `{"format": "vllm", "tokens": [1, 887, 526]}{}`,
			fail:  true,
		},
		// expected failure: a fractional token ID rounds to an integer as float64
		{
			state: `{"format": "vllm", "tokens": [1.0000000000000001]}`,
			fail:  true,
		},
	} {
		r := strings.NewReader(tt.state)
		err := schema.ValidatorMediaTypeModelStateConfig.Validate(r)

		if got := err != nil; tt.fail != got {
			t.Errorf("test %d: expected validation failure %t but got %t, err %v", i, tt.fail, got, err)
		}
	}
}
