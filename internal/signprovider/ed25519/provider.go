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
package ed25519

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
)

// provider is a receiver that encapsulates required dependencies.
type provider struct{}

// New is a factory function that returns an initialized provider.
func New() *provider {
	return &provider{}
}

func (p *provider) Sign(key, content []byte) string {
	keyDecoded := make([]byte, hex.DecodedLen(len(key)))
	hex.Decode(keyDecoded, key)
	signed := ed25519.Sign(keyDecoded, content)
	return fmt.Sprintf("%x", signed)
}

func (p *provider) Verify(key, content, signature []byte) bool {
	keyDecoded := make([]byte, hex.DecodedLen(len(key)))
	hex.Decode(keyDecoded, key)

	sigDecoded := make([]byte, hex.DecodedLen(len(signature)))
	hex.Decode(sigDecoded, signature)
	return ed25519.Verify(keyDecoded, content, sigDecoded)
}
