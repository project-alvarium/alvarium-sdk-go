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

package x509

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/config"
	"math/big"
	"os"
)

// provider is a receiver that encapsulates required dependencies.
type provider struct{}

// New is a factory function that returns an initialized provider.
func New() *provider {
	return &provider{}
}

func (p *provider) Sign(key config.KeyInfo, content []byte) (string, error) {
	prvKeyBytes, err := os.ReadFile(key.Path)
	if err != nil {
		return "", err
	}

	keyDecoded, err := x509.ParseECPrivateKey(prvKeyBytes)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(content)

	b, err := ecdsa.SignASN1(rand.Reader, keyDecoded, hash[:])
	if err != nil {
		return "", err
	}
	return string(b), nil

}

func (p *provider) Verify(key config.KeyInfo, content, signature []byte) (bool, error) {
	pubKeyBytes, err := os.ReadFile(key.Path)
	if err != nil {
		return false, err
	}

	keyDecoded, err := x509.ParsePKIXPublicKey(pubKeyBytes)
	if err != nil {
		return false, err
	}

	hash := sha256.Sum256(content)

	r := new(big.Int).SetBytes(signature[:len(signature)/2])
	s := new(big.Int).SetBytes(signature[len(signature)/2:])
	return ecdsa.Verify(keyDecoded.(*ecdsa.PublicKey), hash[:], r, s), nil
}
