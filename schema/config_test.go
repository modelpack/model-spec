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

package schema_test

import (
	"strings"
	"testing"

	"github.com/modelpack/model-spec/schema"
)

func TestConfig(t *testing.T) {
	for i, tt := range []struct {
		config string
		fail   bool
	}{
		// expected failure: config is missing
		{
			config: `
{
  "descriptor": {
    "name": "xyz-3-8B-Instruct",
    "version": "3.1"
  },
  "modelfs": {
    "type": "layers",
    "diff_ids": [
       "sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
    ]
  }
}
`,
			fail: true,
		},
		// expected failure: version is a number
		{
			config: `
{
  "descriptor": {
    "name": "xyz-3-8B-Instruct",
    "version": 3.1
  },
  "config": {
     "paramSize": "8b"
  },
  "modelfs": {
    "type": "layers",
    "diff_ids": [
       "sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
    ]
  }
}
`,
			fail: true,
		},
		// expected failure: revision is a number
		{
			config: `
{
  "descriptor": {
    "name": "xyz-3-8B-Instruct",
    "version": "3.1",
    "revision": 1234567890
  },
  "config": {
     "paramSize": "8b"
  },
  "modelfs": {
    "type": "layers",
    "diff_ids": [
       "sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
    ]
  }
}
`,
			fail: true,
		},
		// expected failure: createdAt is not RFC3339 format
		{
			config: `
{
  "descriptor": {
    "name": "xyz-3-8B-Instruct",
    "version": "3.1",
    "createdAt": "2025/01/01T00:00:00Z"
  },
  "config": {
     "paramSize": "8b"
  },
  "modelfs": {
    "type": "layers",
    "diff_ids": [
       "sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
    ]
  }
}
`,
			fail: true,
		},
		// expected failure: authors is not an array
		{
			config: `
{
  "descriptor": {
    "name": "xyz-3-8B-Instruct",
    "version": "3.1",
	"authors": "John Doe"
  },
  "config": {
     "paramSize": "8b"
  },
  "modelfs": {
    "type": "layers",
    "diff_ids": [
       "sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
    ]
  }
}
`,
			fail: true,
		},
		// expected failure: licenses is not an array
		{
			config: `
{
  "descriptor": {
    "name": "xyz-3-8B-Instruct",
    "version": "3.1",
	"licenses": "Apache-2.0"
  },
  "config": {
     "paramSize": "8b"
  },
  "modelfs": {
    "type": "layers",
    "diff_ids": [
       "sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
    ]
  }
}
`,
			fail: true,
		},
		// expected failure: docURL is an array
		{
			config: `
{
  "descriptor": {
    "name": "xyz-3-8B-Instruct",
    "version": "3.1",
    "docURL": [
       "https://example.com/doc"
    ]
  },
  "config": {
     "paramSize": "8b"
  },
  "modelfs": {
    "type": "layers",
    "diff_ids": [
       "sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
    ]
  }
}
`,
			fail: true,
		},
		// expected failure: sourceURL is an array
		{
			config: `
{
  "descriptor": {
    "name": "xyz-3-8B-Instruct",
    "version": "3.1",
    "sourceURL": [
       "https://github.com/xyz/xyz3"
    ]
  },
  "config": {
     "paramSize": "8b"
  },
  "modelfs": {
    "type": "layers",
    "diff_ids": [
       "sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
    ]
  }
}
`,
			fail: true,
		},
		// expected failure: paramSize is a number
		{
			config: `
{
  "descriptor": {
    "name": "xyz-3-8B-Instruct",
    "version": "3.1"
  },
  "config": {
     "paramSize": 8000000
  },
  "modelfs": {
    "type": "layers",
    "diff_ids": [
       "sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
    ]
  }
}
`,
			fail: true,
		},
		// expected failure: precision is a number
		{
			config: `
{
  "descriptor": {
    "name": "xyz-3-8B-Instruct",
    "version": "3.1"
  },
  "config": {
     "precision": 16
  },
  "modelfs": {
    "type": "layers",
    "diff_ids": [
       "sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
    ]
  }
}
`,
			fail: true,
		},
		// expected failure: type is not "layers"
		{
			config: `
{
  "descriptor": {
    "name": "xyz-3-8B-Instruct",
    "version": "3.1"
  },
  "config": {
     "paramSize": "8b"
  },
  "modelfs": {
    "type": "layer",
    "diff_ids": [
       "sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
    ]
  }
}
`,
			fail: true,
		},
		// expected failure: diff_ids is not an array
		{
			config: `
{
  "descriptor": {
    "name": "xyz-3-8B-Instruct",
    "version": "3.1"
  },
  "config": {
     "paramSize": "8b"
  },
  "modelfs": {
    "type": "layers",
    "diff_ids": "sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
  }
}
`,
			fail: true,
		},
		// expected failure: diff_ids is empty
		{
			config: `
{
  "descriptor": {
    "name": "xyz-3-8B-Instruct",
    "version": "3.1"
  },
  "config": {
     "paramSize": "8b"
  },
  "modelfs": {
    "type": "layers",
    "diff_ids": []
  }
}
`,
			fail: true,
		},
	} {
		r := strings.NewReader(tt.config)
		err := schema.ValidatorMediaTypeModelConfig.Validate(r)

		if got := err != nil; tt.fail != got {
			t.Errorf("test %d: expected validation failure %t but got %t, err %v", i, tt.fail, got, err)
		}
	}
}
