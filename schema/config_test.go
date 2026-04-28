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
    "diffIds": [
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
    "diffIds": [
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
    "diffIds": [
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
    "diffIds": [
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
    "diffIds": [
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
    "diffIds": [
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
    "diffIds": [
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
    "diffIds": [
       "sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
    ]
  }
}
`,
			fail: true,
		},
		// expected failure: datasetsURL is not an array
		{
			config: `
{
  "descriptor": {
    "name": "xyz-3-8B-Instruct",
    "version": "3.1",
    "sourceURL": "https://github.com/xyz/xyz3",
    "datasetsURL": "https://example.com/dataset"
  },
  "config": {
     "paramSize": "8b"
  },
  "modelfs": {
    "type": "layers",
    "diffIds": [
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
    "diffIds": [
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
    "diffIds": [
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
    "diffIds": [
       "sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
    ]
  }
}
`,
			fail: true,
		},
		// expected failure: diffIds is not an array
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
    "diffIds": "sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
  }
}
`,
			fail: true,
		},
		// expected failure: diffIds is empty
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
    "diffIds": []
  }
}
`,
			fail: true,
		},
		// expected failure: inputTypes is not an array
		{
			config: `
{
  "descriptor": {
    "name": "xyz-3-8B-Instruct",
    "version": "3.1"
  },
  "config": {
     "paramSize": "8b",
     "capabilities": {
        "inputTypes": "text"
     }
  },
  "modelfs": {
    "type": "layers",
    "diffIds": [
       "sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
    ]
  }
}
`,
			fail: true,
		},
		// expected failure: outputTypes is not an array
		{
			config: `
{
  "descriptor": {
    "name": "xyz-3-8B-Instruct",
    "version": "3.1"
  },
  "config": {
     "paramSize": "8b",
     "capabilities": {
        "outputTypes": "text"
     }
  },
  "modelfs": {
    "type": "layers",
    "diffIds": [
       "sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
    ]
  }
}
`,
			fail: true,
		},
		// expected failure: the element of inputTypes/outputTypes is not a valid type
		{
			config: `
{
  "descriptor": {
    "name": "xyz-3-8B-Instruct",
    "version": "3.1"
  },
  "config": {
     "paramSize": "8b",
     "capabilities": {
        "inputTypes": ["img"]
     }
  },
  "modelfs": {
    "type": "layers",
    "diffIds": [
       "sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
    ]
  }
}
`,
			fail: true,
		},
		// expected failure: knowledgeCutoff is not RFC3339 format
		{
			config: `
{
  "descriptor": {
    "name": "xyz-3-8B-Instruct",
    "version": "3.1"
  },
  "config": {
     "paramSize": "8b",
     "capabilities": {
        "inputTypes": ["text"],
        "outputTypes": ["text"],
        "knowledgeCutoff": "2025-01-01"
     }
  },
  "modelfs": {
    "type": "layers",
    "diffIds": [
       "sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
    ]
  }
}
`,
			fail: true,
		},
		// expected failure: reasoning is not boolean
		{
			config: `
{
  "descriptor": {
    "name": "xyz-3-8B-Instruct",
    "version": "3.1"
  },
  "config": {
     "paramSize": "8b",
     "capabilities": {
        "inputTypes": ["text"],
        "outputTypes": ["text"],
        "reasoning": "true"
     }
  },
  "modelfs": {
    "type": "layers",
    "diffIds": [
       "sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
    ]
  }
}
`,
			fail: true,
		},
		// expected failure: toolUsage is not boolean
		{
			config: `
{
  "descriptor": {
    "name": "xyz-3-8B-Instruct",
    "version": "3.1"
  },
  "config": {
     "paramSize": "8b",
     "capabilities": {
        "inputTypes": ["text"],
        "outputTypes": ["text"],
        "toolUsage": "true"
     }
  },
  "modelfs": {
    "type": "layers",
    "diffIds": [
       "sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
    ]
  }
}
`,
			fail: true,
		},
		// expected failure: reward is not boolean
		{
			config: `
{
  "descriptor": {
    "name": "xyz-3-8B-Instruct",
    "version": "3.1"
  },
  "config": {
     "paramSize": "8b",
     "capabilities": {
        "inputTypes": ["text"],
        "outputTypes": ["text"],
        "reward": "true"
     }
  },
  "modelfs": {
    "type": "layers",
    "diffIds": [
       "sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
    ]
  }
}
`,
			fail: true,
		},
		// expected failure: languages is not an array
		{
			config: `
{
  "descriptor": {
    "name": "xyz-3-8B-Instruct",
    "version": "3.1"
  },
  "config": {
     "paramSize": "8b",
     "capabilities": {
        "inputTypes": ["text"],
        "outputTypes": ["text"],
        "languages": "en"
     }
  },
  "modelfs": {
    "type": "layers",
    "diffIds": [
       "sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
    ]
  }
}
`,
			fail: true,
		},
		// expected failure: language code is not a two-letter ISO 639 code
		{
			config: `
{
  "descriptor": {
    "name": "xyz-3-8B-Instruct",
    "version": "3.1"
  },
  "config": {
     "paramSize": "8b",
     "capabilities": {
        "inputTypes": ["text"],
        "outputTypes": ["text"],
        "languages": ["fra"]
     }
  },
  "modelfs": {
    "type": "layers",
    "diffIds": [
       "sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
    ]
  }
}
`,
			fail: true,
		},
		// expected failure: unknown field in capabilities
		{
			config: `
{
  "descriptor": {
    "name": "xyz-3-8B-Instruct",
    "version": "3.1"
  },
  "config": {
     "paramSize": "8b",
     "capabilities": {
        "inputTypes": ["text"],
        "unknownField": true
     }
  },
  "modelfs": {
    "type": "layers",
    "diffIds": [
       "sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
    ]
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

func TestValidateConfigParsesModelNotModelConfig(t *testing.T) {
	// This test verifies that validateConfig correctly parses the full Model structure,
	// not just ModelConfig. Previously, validateConfig unmarshaled into ModelConfig,
	// which always succeeded because all fields are optional.

	// Test 1: Incomplete model with only config (should fail)
	invalidJSON := `{
		"config": {"paramSize": "8b"}
	}`

	err := schema.ValidatorMediaTypeModelConfig.Validate(strings.NewReader(invalidJSON))
	if err == nil {
		t.Fatalf("expected validation to fail for incomplete model")
	}

	// Test 2: Config-only JSON (should fail)
	configOnlyJSON := `{
		"paramSize": "8b",
		"architecture": "transformer"
	}`

	err = schema.ValidatorMediaTypeModelConfig.Validate(strings.NewReader(configOnlyJSON))
	if err == nil {
		t.Fatalf("expected failure for config-only JSON without descriptor/modelfs, but got nil")
	}

	// Test 3: Valid full Model (should pass)
	validJSON := `{
		"descriptor": {"name": "test-model"},
		"config": {"paramSize": "8b"},
		"modelfs": {"type": "layers", "diffIds": ["sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"]}
	}`

	err = schema.ValidatorMediaTypeModelConfig.Validate(strings.NewReader(validJSON))
	if err != nil {
		t.Fatalf("expected valid Model to pass, but got error: %v", err)
	}
}

func TestValidateDigestFormat(t *testing.T) {
	// Test that validateConfig rejects invalid OCI digest formats

	// Test 1: Invalid digest format (no algorithm:hex format)
	invalidDigestJSON := `{
		"descriptor": {"name": "test"},
		"config": {"paramSize": "8b"},
		"modelfs": {
			"type": "layers",
			"diffIds": ["invalid-digest"]
		}
	}`

	err := schema.ValidatorMediaTypeModelConfig.Validate(strings.NewReader(invalidDigestJSON))
	if err == nil {
		t.Fatalf("expected failure for invalid digest format")
	}

	// Test 2: Invalid hex in digest
	invalidHexJSON := `{
		"descriptor": {"name": "test"},
		"config": {"paramSize": "8b"},
		"modelfs": {
			"type": "layers",
			"diffIds": ["sha256:xyz"]
		}
	}`

	err = schema.ValidatorMediaTypeModelConfig.Validate(strings.NewReader(invalidHexJSON))
	if err == nil {
		t.Fatalf("expected failure for invalid hex in digest")
	}

	// Test 3: Empty hash
	emptyHashJSON := `{
		"descriptor": {"name": "test"},
		"config": {"paramSize": "8b"},
		"modelfs": {
			"type": "layers",
			"diffIds": ["sha256:"]
		}
	}`

	err = schema.ValidatorMediaTypeModelConfig.Validate(strings.NewReader(emptyHashJSON))
	if err == nil {
		t.Fatalf("expected failure for empty hash in digest")
	}

	// Test 4: Multiple invalid digests
	multipleInvalidJSON := `{
		"descriptor": {"name": "test"},
		"config": {"paramSize": "8b"},
		"modelfs": {
			"type": "layers",
			"diffIds": ["invalid", "also-invalid"]
		}
	}`

	err = schema.ValidatorMediaTypeModelConfig.Validate(strings.NewReader(multipleInvalidJSON))
	if err == nil {
		t.Fatalf("expected failure for multiple invalid digests")
	}

	// Test 5: Valid digest (should pass)
	validJSON := `{
		"descriptor": {"name": "test"},
		"config": {"paramSize": "8b"},
		"modelfs": {
			"type": "layers",
			"diffIds": ["sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"]
		}
	}`

	err = schema.ValidatorMediaTypeModelConfig.Validate(strings.NewReader(validJSON))
	if err != nil {
		t.Fatalf("expected valid digest to pass, got: %v", err)
	}

	// Test 6: Multiple valid digests (should pass)
	multipleValidJSON := `{
		"descriptor": {"name": "test"},
		"config": {"paramSize": "8b"},
		"modelfs": {
			"type": "layers",
			"diffIds": [
				"sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
				"sha256:abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
			]
		}
	}`

	err = schema.ValidatorMediaTypeModelConfig.Validate(strings.NewReader(multipleValidJSON))
	if err != nil {
		t.Fatalf("expected multiple valid digests to pass, got: %v", err)
	}
}
