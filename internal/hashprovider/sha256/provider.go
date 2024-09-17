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
package sha256

import (
	crypto "crypto/sha256"
	"encoding/hex"
	"strings"
)

// provider is a receiver that encapsulates required dependencies.
type provider struct{}

// New is a factory function that returns an initialized provider.
func New() *provider {
	return &provider{}
}

// Derive converts data to an identity value.
func (*provider) Derive(data []byte) string {
	h := crypto.Sum256(data)
	hashEncoded := make([]byte, hex.EncodedLen(len(h)))
	hex.Encode(hashEncoded, h[:])
	return strings.ToUpper(string(hashEncoded))
}
