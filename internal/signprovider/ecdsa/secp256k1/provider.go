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

package secp256k1

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"github.com/dustinxie/ecc"
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
	privateKeyBytes, err := hex.DecodeString(string(prvKeyBytes))
	if err != nil {
		return "", err
	}

	curve := ecc.P256k1()
	privKey := &ecdsa.PrivateKey{D: new(big.Int).SetBytes(privateKeyBytes)}
	privKey.Curve = curve

	hash := sha256.Sum256(content)
	sig, err := ecc.SignBytes(privKey, hash[:], ecc.Normal)
	if err != nil {
		return "", err
	}
	return string(sig), nil
}

func (p *provider) Verify(key config.KeyInfo, content, signature []byte) (bool, error) {
	pubKeyBytes, err := os.ReadFile(key.Path)
	if err != nil {
		return false, err
	}

	publicKeyBytes, err := hex.DecodeString(string(pubKeyBytes))
	if err != nil {
		return false, err
	}

	x, y := ecc.UnmarshalCompressed(ecc.P256k1(), publicKeyBytes)
	pubKey := &ecdsa.PublicKey{
		Curve: ecc.P256k1(),
		X:     x,
		Y:     y,
	}
	hash := sha256.Sum256(content)
	return ecc.VerifyBytes(pubKey, hash[:], signature, ecc.Normal), nil
}
