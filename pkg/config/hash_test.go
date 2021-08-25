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
	"github.com/project-alvarium/alvarium-sdk-go/pkg/contracts"
	"github.com/project-alvarium/alvarium-sdk-go/test"
	"testing"
)

func TestHashTypeValues(t *testing.T) {
	tests := []struct {
		name         string
		value        contracts.HashType
		expectResult bool
	}{
		{"valid1", contracts.MD5Hash, true},
		{"valid2", contracts.SHA256Hash, true},
		{"valid3", contracts.NoHash, true},
		{"valid4", "invalid", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value.Validate() != tt.expectResult {
				t.Errorf("unexpcted failure for value %s", tt.value)
			}
		})
	}
}

func TestHashTypeUnmarshal(t *testing.T) {
	validMD5 := HashInfo{Type: contracts.MD5Hash}
	validSHA256 := HashInfo{Type: contracts.SHA256Hash}
	validNone := HashInfo{Type: contracts.NoHash}
	invalid := HashInfo{Type: "invalid"}

	tests := []struct {
		name        string
		info        HashInfo
		expectError bool
	}{
		{"validMD5", validMD5, false},
		{"validSHA256", validSHA256, false},
		{"validNone", validNone, false},
		{"invalid hash type", invalid, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := json.Marshal(tt.info)
			var x HashInfo
			err := json.Unmarshal(b, &x)
			test.CheckError(err, tt.expectError, tt.name, t)
			if err == nil {
				if x.Type != tt.info.Type {
					t.Errorf("HashInfo type mismatch expected %s received %s", tt.info.Type, x.Type)
				}
			}
		})
	}
}
