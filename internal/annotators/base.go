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
package annotators

import (
	"encoding/json"
	"fmt"
	"github.com/project-alvarium/alvarium-sdk-go/internal/hashprovider"
	"github.com/project-alvarium/alvarium-sdk-go/internal/hashprovider/md5"
	"github.com/project-alvarium/alvarium-sdk-go/internal/hashprovider/none"
	"github.com/project-alvarium/alvarium-sdk-go/internal/hashprovider/sha256"
	"github.com/project-alvarium/alvarium-sdk-go/internal/signprovider"
	"github.com/project-alvarium/alvarium-sdk-go/internal/signprovider/ed25519"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/config"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/contracts"
	"io/ioutil"
)

func deriveHash(hash contracts.HashType, data []byte) string {
	var h hashprovider.Provider
	switch hash {
	case contracts.MD5Hash:
		h = md5.New()
	case contracts.SHA256Hash:
		h = sha256.New()
	default:
		h = none.New()
	}

	return h.Derive(data)
}

func signAnnotation(key config.KeyInfo, a contracts.Annotation) (string, error) {
	var s signprovider.Provider
	switch key.Type {
	case contracts.KeyEd25519:
		s = ed25519.New()
	default:
		return "", fmt.Errorf("unrecognized key type %s", key.Type)
	}

	b, err := json.Marshal(a)
	if err != nil {
		return "", err
	}

	prv, err := ioutil.ReadFile(key.Path)
	if err != nil {
		return "", err
	}

	signed := s.Sign(prv, b)
	return signed, nil
}

func verifySignature(key config.KeyInfo, src contracts.Annotation) (bool, error) {
	var s signprovider.Provider
	switch key.Type {
	case contracts.KeyEd25519:
		s = ed25519.New()
	default:
		return false, fmt.Errorf("unrecognized key type %s", key.Type)
	}

	// Annotations are signed based on their JSON representation prior to populating the Signature property.
	// Thus we need to reflect that prior state by setting the Signature property to empty before marshalling
	signature := src.Signature
	src.Signature = ""
	b, err := json.Marshal(src)
	if err != nil {
		return false, err
	}

	pub, err := ioutil.ReadFile(key.Path)
	if err != nil {
		return false, err
	}

	ok := s.Verify(pub, b, []byte(signature))
	return ok, nil
}
