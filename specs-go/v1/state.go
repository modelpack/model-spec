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

package v1

// State defines the configuration of a model state artifact.
// It provides the `application/vnd.cncf.model.state.config.v1+json` mediatype when marshalled to JSON.
type State struct {
	// The format of the state's layers and its version, such as example.session/v1
	Format string `json:"format"`

	// The token IDs, from the model's vocabulary, that the state was built from, in order
	Tokens []int `json:"tokens"`

	// The precision of the state's tensors, such as bfloat16, float8_e4m3 or q8_0
	Precision string `json:"precision,omitempty"`
}
