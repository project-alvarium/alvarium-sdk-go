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
	"context"
	"encoding/json"
	"os"

	"github.com/project-alvarium/alvarium-sdk-go/pkg/config"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/contracts"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/interfaces"
)

// PkiAnnotator is used to validate whether the signature on a given piece of data is valid
type PkiAnnotator struct {
	hash      interfaces.HashProvider
	hashType  contracts.HashType
	kind      contracts.AnnotationType
	signature interfaces.SignatureProvider
	privKey   config.KeyInfo
	pubKey    config.KeyInfo
	layer     contracts.LayerType
}

func NewPkiAnnotator(cfg config.SdkInfo, hash interfaces.HashProvider, sign interfaces.SignatureProvider) interfaces.Annotator {
	a := PkiAnnotator{}
	a.hash = hash
	a.hashType = cfg.Hash.Type
	a.kind = contracts.AnnotationPKI
	a.signature = sign
	a.privKey = cfg.Signature.PrivateKey
	a.pubKey = cfg.Signature.PublicKey
	a.layer = cfg.Layer
	return &a
}

func (a *PkiAnnotator) Do(ctx context.Context, data []byte) (contracts.Annotation, error) {
	key := a.hash.Derive(data)
	hostname, _ := os.Hostname()

	var sig signable
	err := json.Unmarshal(data, &sig)
	if err != nil {
		return contracts.Annotation{}, err
	}

	ok, err := sig.verifySignature(a.pubKey, a.signature)
	if err != nil {
		return contracts.Annotation{}, err
	}
	annotation := contracts.NewAnnotation(string(key), a.hashType, hostname, a.layer, a.kind, ok)

	b, err := json.Marshal(annotation)
	if err != nil {
		return contracts.Annotation{}, err
	}
	signed, err := a.signature.Sign(a.privKey, b)
	if err != nil {
		return contracts.Annotation{}, err
	}
	annotation.Signature = signed
	return annotation, nil
}

// The question of how/whether to validate signed data is tricky. We want this SDK to be as agnostic of the application data
// as possible.
//
// As to the "how", we need to require some minimal property support for the data being annotated by Alvarium. For
// example, the data type could have a property on it indicating the signature of the data creator, as well as a seed
// property that was used to generate the signature. But what are these properties called -- signature, sig, signed, etc?
// We would need to formally specify.
//
// Secondarily, it could be possible to pass the signature as a header on a pub/sub message or an HTTP call. I don't know
// the ultimate answer to where this should reside. For now I'm making the following assumptions.
// 1.) The incoming []byte is a JSON string
// 2.) That JSON can be unmarshaled into a type that has "Signature" and "Seed" properties of type string
//
// As to the "whether", the use case here is that the customer wants to validate a signature on the data itself and attest
// to that validation in flight. It is possible to validate a signature at a later stage at the point where all of the
// annotations are assessed to calculate the final score. The signature on the annotation itself, having been generated
// by the private key of the host machine and verified through a public key, could be enough to trust the associated data.
//
// It's likely that at some point this type should be placed in pkg/contracts
type signable struct {
	Seed      string `json:"seed,omitempty"`
	Signature string `json:"signature,omitempty"`
}

func (s *signable) verifySignature(key config.KeyInfo, signature interfaces.SignatureProvider) (bool, error) {
	return signature.Verify(key, []byte(s.Seed), []byte(s.Signature))
}
