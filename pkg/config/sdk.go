/*******************************************************************************
 * Copyright 2023 Dell Inc.
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
	"log/slog"

	"github.com/project-alvarium/alvarium-sdk-go/pkg/contracts"
	"gopkg.in/yaml.v3"
)

type SdkInfo struct {
	Annotators []contracts.AnnotationType `json:"annotators,omitempty" yaml:"annotators"`
	Hash       HashInfo                   `json:"hash,omitempty" yaml:"hash"`
	Signature  SignatureInfo              `json:"signature,omitempty" yaml:"signature"`
	Stream     StreamInfo                 `json:"stream,omitempty" yaml:"stream"`
}

type LoggingInfo struct {
	MinLogLevel slog.Level `json:"minLogLevel,omitempty"`
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

func (s *SdkInfo) UnmarshalYAML(data *yaml.Node) (err error) {
	type Alias struct {
		Annotators []contracts.AnnotationType `yaml:"annotators"`
		Hash       HashInfo                   `yaml:"hash"`
		Signature  SignatureInfo              `yaml:"signature"`
		Stream     StreamInfo                 `yaml:"stream"`
	}
	a := Alias{}
	// Error with unmarshaling
	if err = data.Decode(&a); err != nil {
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
