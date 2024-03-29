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

// Default path to a TPM 2.0 device (https://wiki.archlinux.org/title/Trusted_Platform_Module)
const tpmPath string = "/dev/tpm0"

// TpmAnnotator is used to attest whether or not the host machine has TPM capability for managing secrets
type TpmAnnotator struct {
	hash      interfaces.HashProvider
	hashType  contracts.HashType
	kind      contracts.AnnotationType
	signature interfaces.SignatureProvider
	privKey   config.KeyInfo
	layer     contracts.LayerType
}

func NewTpmAnnotator(cfg config.SdkInfo, hash interfaces.HashProvider, sign interfaces.SignatureProvider) interfaces.Annotator {
	a := TpmAnnotator{}
	a.hash = hash
	a.hashType = cfg.Hash.Type
	a.kind = contracts.AnnotationTPM
	a.signature = sign
	a.privKey = cfg.Signature.PrivateKey
	a.layer = cfg.Layer
	return &a
}

func (a *TpmAnnotator) Do(ctx context.Context, data []byte) (contracts.Annotation, error) {
	key := a.hash.Derive(data)
	hostname, _ := os.Hostname()
	isSatisfied := false

	// If mounted path exists, check that it's either a device or a socket (emulator)
	// This logic based on code found in gp-tpm module.
	// https://github.com/google/go-tpm/blob/b3942ee5b15a7bd19e6419d5903e6e64fbb3d4ba/tpmutil/run_other.go#L29
	fi, err := os.Stat(tpmPath)
	if err == nil {
		// TPM mounted at default path
		if fi.Mode()&os.ModeDevice != 0 || fi.Mode()&os.ModeSocket != 0 {
			isSatisfied = true
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
