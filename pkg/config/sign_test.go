/*******************************************************************************
 * Copyright 2021 Dell Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/
package config

import (
	"encoding/json"
	"github.com/DyrellC/alvarium-sdk-go/pkg/contracts"
	"github.com/DyrellC/alvarium-sdk-go/test"
	"testing"
)

func TestKeyInfoUnmarshal(t *testing.T) {
	pass := KeyInfo{
		Type: contracts.KeyEd25519,
	}

	fail := KeyInfo{
		Type: "invalid",
	}

	tests := []struct {
		name        string
		k           KeyInfo
		expectError bool
	}{
		{"valid key ed25519", pass, false},
		{"invalid key", fail, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := json.Marshal(tt.k)
			var x KeyInfo
			err := json.Unmarshal(b, &x)
			test.CheckError(err, tt.expectError, tt.name, t)
		})
	}
}
