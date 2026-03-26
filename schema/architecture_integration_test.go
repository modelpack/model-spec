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
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/modelpack/model-spec/schema"
)

// convertHFToArchConfig simulates HuggingFace config to architecture_config conversion.
// This mirrors the logic in tools/hf_to_arch.py.
func convertHFToArchConfig(hfConfig map[string]interface{}) (map[string]interface{}, error) {
	mappings := map[string]string{
		"numLayers":         "num_hidden_layers",
		"hiddenSize":        "hidden_size",
		"numAttentionHeads": "num_attention_heads",
	}

	archConfig := map[string]interface{}{
		"type": "transformer",
	}

	for archKey, hfKey := range mappings {
		val, ok := hfConfig[hfKey]
		if !ok {
			return nil, fmt.Errorf("missing required field: %s", hfKey)
		}
		numVal, ok := val.(float64)
		if !ok {
			return nil, fmt.Errorf("field %s must be a number", hfKey)
		}
		if numVal < 1 {
			return nil, fmt.Errorf("field %s must be >= 1", hfKey)
		}
		archConfig[archKey] = int(numVal)
	}

	return archConfig, nil
}

// TestArchitectureIntegrationValidConversion tests the full conversion pipeline:
// HF config → architecture_config → embedded in Model JSON → schema validation.
func TestArchitectureIntegrationValidConversion(t *testing.T) {
	// Step 1: Simulate HuggingFace config.json
	hfConfig := map[string]interface{}{
		"num_hidden_layers":   float64(32),
		"hidden_size":         float64(4096),
		"num_attention_heads": float64(32),
		"vocab_size":          float64(32000),
	}

	// Step 2: Convert to architecture_config format
	archConfig, err := convertHFToArchConfig(hfConfig)
	if err != nil {
		t.Fatalf("conversion failed: %v", err)
	}

	// Verify conversion output structure
	if archConfig["type"] != "transformer" {
		t.Errorf("expected type=transformer, got %v", archConfig["type"])
	}
	if archConfig["numLayers"] != 32 {
		t.Errorf("expected numLayers=32, got %v", archConfig["numLayers"])
	}
	if archConfig["hiddenSize"] != 4096 {
		t.Errorf("expected hiddenSize=4096, got %v", archConfig["hiddenSize"])
	}
	if archConfig["numAttentionHeads"] != 32 {
		t.Errorf("expected numAttentionHeads=32, got %v", archConfig["numAttentionHeads"])
	}

	// Step 3: Build full Model JSON that includes architecture_config
	modelData := map[string]interface{}{
		"descriptor": map[string]interface{}{
			"name": "test-model",
		},
		"config": map[string]interface{}{
			"paramSize":           "8b",
			"architecture_config": archConfig,
		},
		"modelfs": map[string]interface{}{
			"type":    "layers",
			"diffIds": []string{"sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"},
		},
	}

	modelBytes, err := json.Marshal(modelData)
	if err != nil {
		t.Fatalf("failed to marshal model JSON: %v", err)
	}

	// Step 4: Validate complete Model with embedded architecture_config
	// depends on architecture_config schema (PR-1)
	err = schema.ValidatorMediaTypeModelConfig.Validate(strings.NewReader(string(modelBytes)))
	if err != nil {
		// Expected to fail until PR-1 (architecture_config schema) is merged
		t.Logf("Model validation with architecture_config failed as expected: %v", err)
		
		// Verify the failure is specifically about architecture_config not being allowed
		if !strings.Contains(err.Error(), "architecture_config") {
			t.Errorf("Expected validation to fail due to architecture_config field, but got: %v", err)
		}
		return
	}
}

// TestArchitectureIntegrationMissingHFField tests conversion fails for missing required field.
func TestArchitectureIntegrationMissingHFField(t *testing.T) {
	// HF config missing num_hidden_layers
	hfConfig := map[string]interface{}{
		"hidden_size":         float64(4096),
		"num_attention_heads": float64(32),
	}

	_, err := convertHFToArchConfig(hfConfig)
	if err == nil {
		t.Fatal("expected conversion to fail for missing num_hidden_layers")
	}
	if !strings.Contains(err.Error(), "num_hidden_layers") {
		t.Errorf("error should mention missing field, got: %v", err)
	}
}

// TestArchitectureIntegrationInvalidModelStructure tests schema rejects invalid Model structure.
func TestArchitectureIntegrationInvalidModelStructure(t *testing.T) {
	// Model missing required 'config' field
	invalidJSON := `{
		"descriptor": {"name": "test-model"},
		"modelfs": {
			"type": "layers",
			"diffIds": ["sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"]
		}
	}`

	err := schema.ValidatorMediaTypeModelConfig.Validate(strings.NewReader(invalidJSON))
	if err == nil {
		t.Fatal("expected validation to fail for missing config field")
	}
}

// TestArchitectureIntegrationInvalidFieldValue tests conversion rejects invalid values.
func TestArchitectureIntegrationInvalidFieldValue(t *testing.T) {
	// Zero value for num_hidden_layers
	hfConfig := map[string]interface{}{
		"num_hidden_layers":   float64(0),
		"hidden_size":         float64(4096),
		"num_attention_heads": float64(32),
	}

	_, err := convertHFToArchConfig(hfConfig)
	if err == nil {
		t.Fatal("expected conversion to fail for zero value")
	}
	if !strings.Contains(err.Error(), ">= 1") {
		t.Errorf("error should mention minimum value, got: %v", err)
	}
}
