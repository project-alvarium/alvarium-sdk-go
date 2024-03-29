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
	"encoding/json"
	"github.com/project-alvarium/alvarium-sdk-go/internal/hashprovider/md5"
	"github.com/project-alvarium/alvarium-sdk-go/internal/hashprovider/none"
	sha2562 "github.com/project-alvarium/alvarium-sdk-go/internal/hashprovider/sha256"
	"github.com/project-alvarium/alvarium-sdk-go/internal/signprovider/ed25519"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/config"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/contracts"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/interfaces"
	"github.com/project-alvarium/alvarium-sdk-go/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeriveHash(t *testing.T) {
	noneInput := test.FactoryRandomFixedLengthString(64, test.AlphanumericCharset)
	noneOutput := noneInput
	noHash := none.New()

	md5Input := []byte("foo")
	md5Output := "acbd18db4cc2f85cedef654fccc4a4d8"
	md5Hash := md5.New()

	sha256Input := []byte("bar")
	sha256Output := "fcde2b2edba56bf408601fb721fe9b5c338d10ee429ea04fae5511b68fbf8fb9"
	sha256Hash := sha2562.New()

	tests := []struct {
		name         string
		hashType     contracts.HashType
		hashProvider interfaces.HashProvider
		input        []byte
		output       string
	}{
		{"derive no hash", contracts.NoHash, noHash, []byte(noneInput), noneOutput},
		{"derive md5 hash", contracts.MD5Hash, md5Hash, md5Input, md5Output},
		{"derive sha256 hash", contracts.SHA256Hash, sha256Hash, sha256Input, sha256Output},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.hashProvider.Derive(tt.input)
			assert.Equal(t, tt.output, result)
		})
	}
}

func TestSignAnnotation(t *testing.T) {
	private := config.KeyInfo{
		Type: contracts.KeyEd25519,
		Path: "../../test/keys/ed25519/private.key",
	}

	//I'm using a JSON representation here b/c calling the Annotation constructor will populate different ID and Timestamp values each time.
	//Thus the resulting signature will be different for each test run if I don't use something static.
	var a contracts.Annotation
	sample := "{\"id\":\"01F9MS7QVH8Z3KMW757RGFKCBG\",\"key\":\"dummyKey\",\"hash\":\"none\",\"host\":\"ubuntu\",\"type\":\"tpm\",\"timestamp\":\"2021-07-02T18:35:36.561920812-05:00\"}"
	json.Unmarshal([]byte(sample), &a)

	tests := []struct {
		name        string
		cfg         config.KeyInfo
		annotation  contracts.Annotation
		signature   string
		sigProvider interfaces.SignatureProvider
		expectError bool
	}{
		{"valid ed25519 signature", private, a,
			"dafcee0a1844b9c5c0db87252067afc06853afefa14b1711a7c24a6eabc6f6b76b91f8aae2ce54f74b54500dbaf303fbb23dc550e151cfe03cef68ae26a22306",
			ed25519.New(), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := json.Marshal(tt.annotation)
			if err != nil {
				t.Error(err)
			}
			result, err := tt.sigProvider.Sign(tt.cfg, b)

			if err == nil {
				assert.Equal(t, tt.signature, result)
			}
		})
	}
}
