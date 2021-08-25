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
	"context"
	"encoding/json"
	"fmt"
	"github.com/project-alvarium/alvarium-sdk-go/internal/signprovider"
	"github.com/project-alvarium/alvarium-sdk-go/internal/signprovider/ed25519"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/config"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/contracts"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/interfaces"
	"io/ioutil"
	"os"
)

// PkiAnnotator is used to validate whether the signature on a given piece of data is valid
type PkiAnnotator struct {
	hash contracts.HashType
	kind contracts.AnnotationType
	sign config.SignatureInfo
}

func NewPkiAnnotator(cfg config.SdkInfo) interfaces.Annotator {
	a := PkiAnnotator{}
	a.hash = cfg.Hash.Type
	a.kind = contracts.AnnotationPKI
	a.sign = cfg.Signature
	return &a
}

func (a *PkiAnnotator) Do(ctx context.Context, data []byte) (contracts.Annotation, error) {
	key := deriveHash(a.hash, data)
	hostname, _ := os.Hostname()

	var sig signable
	err := json.Unmarshal(data, &sig)
	if err != nil {
		return contracts.Annotation{}, err
	}

	ok, err := sig.verifySignature(a.sign.PublicKey)
	if err != nil {
		return contracts.Annotation{}, err
	}
	annotation := contracts.NewAnnotation(string(key), a.hash, hostname, a.kind, ok)
	signed, err := signAnnotation(a.sign.PrivateKey, annotation)
	if err != nil {
		return contracts.Annotation{}, err
	}
	annotation.Signature = string(signed)
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
type signable struct {
	Seed      string `json:"seed,omitempty"`
	Signature string `json:"signature,omitempty"`
}

func (s *signable) verifySignature(key config.KeyInfo) (bool, error) {
	if len(s.Signature) == 0 { // no signature detected
		return false, nil
	}
	var p signprovider.Provider
	switch key.Type {
	case contracts.KeyEd25519:
		p = ed25519.New()
	default:
		return false, fmt.Errorf("unrecognized key type %s", key.Type)
	}

	pub, err := ioutil.ReadFile(key.Path)
	if err != nil {
		return false, err
	}

	ok := p.Verify(pub, []byte(s.Seed), []byte(s.Signature))
	return ok, nil
}
