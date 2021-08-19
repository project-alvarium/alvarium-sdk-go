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

type HashInfo struct {
	Type contracts.HashType `json:"type,omitempty"`
}

func (h *HashInfo) UnmarshalJSON(data []byte) (err error) {
	type Alias struct {
		Type contracts.HashType
	}

	a := Alias{}
	if err = json.Unmarshal(data, &a); err != nil {
		return err
	}

	if !a.Type.Validate() {
		return fmt.Errorf("invalid HashType value provided %s", a.Type)
	}
	h.Type = a.Type
	return nil
}
