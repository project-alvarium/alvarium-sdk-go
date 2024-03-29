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

package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	handler "github.com/project-alvarium/alvarium-sdk-go/internal/annotators/http/handler"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/config"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/contracts"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/interfaces"
)

// HttpPkiAnnotator is used to validate whether the signature on a given piece of data is valid, both sent in the HTTP message
type HttpPkiAnnotator struct {
	hash      interfaces.HashProvider
	hashType  contracts.HashType
	kind      contracts.AnnotationType
	signature interfaces.SignatureProvider
	privKey   config.KeyInfo
	pubKey    config.KeyInfo
	layer     contracts.LayerType
}

func NewHttpPkiAnnotator(cfg config.SdkInfo, hash interfaces.HashProvider, sign interfaces.SignatureProvider) interfaces.Annotator {
	a := HttpPkiAnnotator{}
	a.hash = hash
	a.hashType = cfg.Hash.Type
	a.kind = contracts.AnnotationPKIHttp
	a.signature = sign
	a.privKey = cfg.Signature.PrivateKey
	a.pubKey = cfg.Signature.PublicKey
	a.layer = cfg.Layer
	return &a
}

func (a *HttpPkiAnnotator) Do(ctx context.Context, data []byte) (contracts.Annotation, error) {
	key := a.hash.Derive(data)
	hostname, _ := os.Hostname()

	//Call parser on request
	req := ctx.Value(contracts.HttpRequestKey)
	parsed, err := handler.ParseSignature(req.(*http.Request))

	if err != nil {
		return contracts.Annotation{}, err
	}
	var sig signable
	sig.Seed = parsed.Seed
	sig.Signature = parsed.Signature

	// Use the parsed request to obtain the key name and type we should use to validate the signature
	var k config.KeyInfo
	directory := filepath.Dir(a.pubKey.Path)
	k.Path = strings.Join([]string{directory, parsed.Keyid}, "/")
	k.Type = contracts.KeyAlgorithm(parsed.Algorithm)
	if !(k.Type.Validate()) {
		return contracts.Annotation{}, errors.New("invalid key type specified: " + parsed.Algorithm)
	}

	ok, err := sig.verifySignature(k, a.signature)
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

// TODO: At this point this type has converged with the one defined in annotators/pki.go. Eliminate duplicate definition
type signable struct {
	Seed      string
	Signature string
}

func (s *signable) verifySignature(key config.KeyInfo, signature interfaces.SignatureProvider) (bool, error) {
	return signature.Verify(key, []byte(s.Seed), []byte(s.Signature))
}
