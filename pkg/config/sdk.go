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
	"github.com/project-alvarium/alvarium-sdk-go/pkg/contracts"
)

type SdkInfo struct {
	Annotators []contracts.AnnotationType `json:"annotators,omitempty"`
	Hash       HashInfo                   `json:"hash,omitempty"`
	Signature  SignatureInfo              `json:"signature,omitempty"`
	Stream     StreamInfo                 `json:"stream,omitempty"`
}

func (s *SdkInfo) UnmarshalJSON(data []byte) (err error) {
	type Alias struct {
		Annotators []contracts.AnnotationType
		Hash       HashInfo
		Signature  SignatureInfo
		Stream     StreamInfo
	}
	a := Alias{}
	// Error with unmarshaling
	if err = json.Unmarshal(data, &a); err != nil {
		return err
	}

	if len(a.Annotators) > 0 {
		for _, x := range a.Annotators {
			ok := x.Validate()
			if !ok {
				return fmt.Errorf("invalid AnnotationType received %s", x)
			}
		}
	}
	s.Annotators = a.Annotators
	s.Hash = a.Hash
	s.Signature = a.Signature
	s.Stream = a.Stream
	return nil
}
