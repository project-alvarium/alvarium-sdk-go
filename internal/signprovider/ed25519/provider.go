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

package ed25519

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/config"
	"os"
)

// provider is a receiver that encapsulates required dependencies.
type provider struct{}

// New is a factory function that returns an initialized provider.
func New() *provider {
	return &provider{}
}

func (p *provider) Sign(key config.KeyInfo, content []byte) (string, error) {
	prv, err := os.ReadFile(key.Path)
	if err != nil {
		return "", err
	}

	keyDecoded := make([]byte, hex.DecodedLen(len(prv)))
	hex.Decode(keyDecoded, prv)
	signed := ed25519.Sign(keyDecoded, content)
	return fmt.Sprintf("%x", signed), nil
}

func (p *provider) Verify(key config.KeyInfo, content, signature []byte) (bool, error) {
	pub, err := os.ReadFile(key.Path)
	if err != nil {
		return false, err
	}

	keyDecoded := make([]byte, hex.DecodedLen(len(pub)))
	hex.Decode(keyDecoded, pub)

	sigDecoded := make([]byte, hex.DecodedLen(len(signature)))
	hex.Decode(sigDecoded, signature)
	return ed25519.Verify(keyDecoded, content, sigDecoded), nil
}
