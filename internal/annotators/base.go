/*******************************************************************************
 * Copyright 2024 Dell Inc.
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
	"github.com/project-alvarium/alvarium-sdk-go/pkg/config"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/contracts"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/interfaces"
)

func SignAnnotation(key config.KeyInfo, signature interfaces.SignatureProvider, a contracts.Annotation) (string, error) {
	b, err := json.Marshal(a)
	if err != nil {
		return "", err
	}

	return signature.Sign(key, b)
}

// VerifySignature will validate the signature on an Annotation
//
// Currently (29-Mar-2024) this method is only used by unit tests
func VerifySignature(key config.KeyInfo, signature interfaces.SignatureProvider, src contracts.Annotation) (bool, error) {
	// Annotations are signed based on their JSON representation prior to populating the Signature property.
	// Thus we need to reflect that prior state by setting the Signature property to empty before marshalling
	verifiable := src.Signature
	src.Signature = ""
	b, err := json.Marshal(src)
	if err != nil {
		return false, err
	}

	return signature.Verify(key, b, []byte(verifiable))
}
