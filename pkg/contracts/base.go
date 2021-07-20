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
package contracts

import (
	"github.com/oklog/ulid/v2"
	"math/rand"
	"time"
)

// NewULID is a convenience function for generating ULIDs where necessary.
// It is based on the example showing the "quick" method for initializing the algorithm's entropy parameter.
// https://github.com/oklog/ulid/blob/c6bb9e1d94a82e71dfd7ff279aa6cea7c52779bb/cmd/ulid/main.go#L67
// As described on that page, quick means "when generating, use non-crypto-grade entropy".
func NewULID() ulid.ULID {
	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)
	entropy := rand.New(source)

	id, _ := ulid.New(ulid.Timestamp(time.Now()), entropy)

	return id
}
