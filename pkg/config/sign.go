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
	"fmt"
	"github.com/DyrellC/alvarium-sdk-go/pkg/contracts"
)

type SignatureInfo struct {
	PublicKey  KeyInfo `json:"public,omitempty"`
	PrivateKey KeyInfo `json:"private,omitempty"`
}

type KeyInfo struct {
	Type contracts.KeyAlgorithm `json:"type,omitempty"` // Type indicates the algorithm used to generate the key
	Path string                 `json:"path,omitempty"` // Path indicates the filesystem path to the key.
	// Path will need to be extended later. Consider that keys may be sourced from difference locations -- file, TPM,
	// Vault, etc.
}

func (k *KeyInfo) UnmarshalJSON(data []byte) (err error) {
	type Alias struct {
		Type contracts.KeyAlgorithm `json:"type,omitempty"`
		Path string                 `json:"path,omitempty"`
	}
	a := Alias{}
	// Error with unmarshaling
	if err = json.Unmarshal(data, &a); err != nil {
		return err
	}

	if !a.Type.Validate() {
		return fmt.Errorf("invalid KeyAlgorithm value provided %s", a.Type)
	}
	k.Type = a.Type
	k.Path = a.Path

	return nil
}
