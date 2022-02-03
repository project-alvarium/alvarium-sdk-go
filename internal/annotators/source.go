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
	"github.com/project-alvarium/alvarium-sdk-go/pkg/config"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/contracts"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/interfaces"
	"os"
)

// SourceAnnotator is used to provide lineage from one version of data to another as the result of a change or transformation.
type SourceAnnotator struct {
	hash contracts.HashType
	kind contracts.AnnotationType
	sign config.SignatureInfo
}

func NewSourceAnnotator(cfg config.SdkInfo) interfaces.Annotator {
	a := SourceAnnotator{}
	a.hash = cfg.Hash.Type
	a.kind = contracts.AnnotationSource
	a.sign = cfg.Signature
	return &a
}

func (a *SourceAnnotator) Do(ctx context.Context, data []byte) (contracts.Annotation, error) {
	key := DeriveHash(a.hash, data)
	hostname, _ := os.Hostname()

	annotation := contracts.NewAnnotation(key, a.hash, hostname, a.kind, true)
	sig, err := SignAnnotation(a.sign.PrivateKey, annotation)
	if err != nil {
		return contracts.Annotation{}, err
	}
	annotation.Signature = string(sig)
	return annotation, nil
}
