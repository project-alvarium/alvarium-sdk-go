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
package contracts

import (
	"encoding/json"
	"fmt"
	"github.com/oklog/ulid/v2"
	"time"
)

// Annotation represents an individual criterion of evaluation in regard to a piece of data
type Annotation struct {
	Id        ulid.ULID      `json:"id,omitempty"`        // Id should probably be a ULID -- uniquely identifies the annotation itself
	Key       string         `json:"key,omitempty"`       // Key is the hash value of the data being annotated
	Hash      HashType       `json:"hash,omitempty"`      // Hash identifies which algorithm was used to construct the hash
	Host      string         `json:"host,omitempty"`      // Host is the hostname of the node making the annotation
	Kind      AnnotationType `json:"type,omitempty"`      // Kind indicates what kind of annotation this is
	Signature string         `json:"signature,omitempty"` // Signature contains the signature of the party making the annotation
	Satisfied bool           `json:"satisfied"`           // Satisfied indicates whether the criteria defining the annotation were fulfilled
	Timestamp time.Time      `json:"timestamp,omitempty"` // Timestamp indicates when the annotation was created
}

// AnnotationList is an envelope for zero to many annotations
type AnnotationList struct {
	Items []Annotation `json:"items,omitempty"` // Items contains 0-many annotations
}

// NewAnnotation is the constructor for an Annotation instance.
func NewAnnotation(key string, hash HashType, host string, kind AnnotationType, satisfied bool) Annotation {
	return Annotation{
		Id:        NewULID(),
		Key:       key,
		Hash:      hash,
		Host:      host,
		Kind:      kind,
		Satisfied: satisfied,
		Timestamp: time.Now(),
	}
}

func (a *Annotation) UnmarshalJSON(data []byte) (err error) {
	type Alias struct {
		Id        ulid.ULID
		Key       string
		Hash      HashType
		Host      string
		Kind      AnnotationType
		Signature string
		Satisfied bool
		Timestamp time.Time
	}
	x := Alias{}
	// Error with unmarshaling
	if err = json.Unmarshal(data, &x); err != nil {
		return err
	}

	if !x.Hash.Validate() {
		return fmt.Errorf("invalid HashType value provided %s", x.Hash)
	}

	if !x.Kind.Validate() {
		return fmt.Errorf("invalid AnnotationType value provided %s", x.Hash)
	}

	a.Id = x.Id
	a.Key = x.Key
	a.Hash = x.Hash
	a.Host = x.Host
	a.Kind = x.Kind
	a.Signature = x.Signature
	a.Satisfied = x.Satisfied
	a.Timestamp = x.Timestamp
	return nil
}
