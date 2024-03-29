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

// SourceAnnotator is used to provide lineage from one version of data to another as the result of a change or transformation.
type SourceAnnotator struct {
	hash      interfaces.HashProvider
	hashType  contracts.HashType
	kind      contracts.AnnotationType
	signature interfaces.SignatureProvider
	privKey   config.KeyInfo
	layer     contracts.LayerType
}

func NewSourceAnnotator(cfg config.SdkInfo, hash interfaces.HashProvider, sign interfaces.SignatureProvider) interfaces.Annotator {
	a := SourceAnnotator{}
	a.hash = hash
	a.hashType = cfg.Hash.Type
	a.kind = contracts.AnnotationSource
	a.signature = sign
	a.privKey = cfg.Signature.PrivateKey
	a.layer = cfg.Layer
	return &a
}

func (a *SourceAnnotator) Do(ctx context.Context, data []byte) (contracts.Annotation, error) {
	key := a.hash.Derive(data)
	hostname, _ := os.Hostname()

	annotation := contracts.NewAnnotation(key, a.hashType, hostname, a.layer, a.kind, true)
	b, err := json.Marshal(annotation)
	if err != nil {
		return contracts.Annotation{}, err
	}
	sig, err := a.signature.Sign(a.privKey, b)
	if err != nil {
		return contracts.Annotation{}, err
	}
	annotation.Signature = sig
	return annotation, nil
}
