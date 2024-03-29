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
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/project-alvarium/alvarium-sdk-go/pkg/config"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/contracts"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/interfaces"
)

type TlsAnnotator struct {
	hash      interfaces.HashProvider
	hashType  contracts.HashType
	kind      contracts.AnnotationType
	signature interfaces.SignatureProvider
	privKey   config.KeyInfo
	layer     contracts.LayerType
}

func NewTlsAnnotator(cfg config.SdkInfo, hash interfaces.HashProvider, sign interfaces.SignatureProvider) interfaces.Annotator {
	a := TlsAnnotator{}
	a.hash = hash
	a.hashType = cfg.Hash.Type
	a.kind = contracts.AnnotationTLS
	a.signature = sign
	a.privKey = cfg.Signature.PrivateKey
	a.layer = cfg.Layer
	return &a
}

func (a *TlsAnnotator) Do(ctx context.Context, data []byte) (contracts.Annotation, error) {
	key := a.hash.Derive(data)
	hostname, _ := os.Hostname()
	isSatisfied := false

	// Currently this annotator should only be used in the context of HTTP. TLS is also applicable to pub/sub but
	// a different approach may be required in that scenario to annotate based on the connection rather than per
	// message. More thought required.
	//
	// The methodology below is also very suitable to HTTP requests given that the tls.ConnectionState is readily
	// available off the incoming request whereas pub/sub connection providers may only expose the tls.Config (see
	// https://pkg.go.dev/crypto/tls#Config) requiring a function implementation for
	// VerifyConnection func(ConnectionState) error
	val := ctx.Value(contracts.AnnotationTLS)
	if val != nil {
		tls, ok := val.(*tls.ConnectionState)
		if !ok {
			return contracts.Annotation{}, errors.New(fmt.Sprintf("unexpected type %T", tls))
		}
		if tls != nil {
			isSatisfied = tls.HandshakeComplete
		}
	}
	annotation := contracts.NewAnnotation(key, a.hashType, hostname, a.layer, a.kind, isSatisfied)

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
