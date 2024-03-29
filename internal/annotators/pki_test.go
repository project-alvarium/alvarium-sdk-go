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
	hash256 "github.com/project-alvarium/alvarium-sdk-go/internal/hashprovider/sha256"
	"github.com/project-alvarium/alvarium-sdk-go/internal/signprovider/ed25519"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/config"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/interfaces"
	"github.com/project-alvarium/alvarium-sdk-go/test"
	"os"
	"testing"
)

func TestPkiAnnotator_Do(t *testing.T) {
	b, err := os.ReadFile("../../test/res/config.json")
	if err != nil {
		t.Fatalf(err.Error())
	}

	var cfg config.SdkInfo
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		t.Fatalf(err.Error())
	}

	badKeyType := cfg
	badKeyType.Signature.PublicKey.Type = "invalid"

	keyNotFound := cfg
	keyNotFound.Signature.PublicKey.Path = "/dev/null/private.key"

	// Set up example signed data type for test purposes
	type testData struct {
		Seed      string
		Signature string
	}
	t1 := testData{
		Seed: test.FactoryRandomFixedLengthString(64, test.AlphanumericCharset),
	}
	signer := ed25519.New()
	t1.Signature, err = signer.Sign(cfg.Signature.PrivateKey, []byte(t1.Seed))
	if err != nil {
		t.Fatalf(err.Error())
	}
	// end of basic example type setup

	t2 := t1
	t2.Signature = ""

	t3 := t1
	t3.Seed = "invalid"

	h := hash256.New()
	tests := []struct {
		name        string
		data        testData
		cfg         config.SdkInfo
		h           interfaces.HashProvider
		s           interfaces.SignatureProvider
		expectError bool
	}{
		{"pki annotation OK", t1, cfg, h, signer, false},
		{"pki key not found", t1, keyNotFound, h, signer, true},
		{"pki empty signature", t2, cfg, h, signer, false},
		{"pki invalid signature", t3, cfg, h, signer, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tpm := NewPkiAnnotator(tt.cfg, tt.h, tt.s)
			b, _ := json.Marshal(tt.data)
			anno, err := tpm.Do(context.Background(), b)
			test.CheckError(err, tt.expectError, tt.name, t)
			if err == nil {
				result, err := VerifySignature(tt.cfg.Signature.PublicKey, tt.s, anno)
				if err != nil {
					t.Error(err.Error())
				} else if !result {
					t.Error("signature not verified")
				}
				if tt.name == "pki empty signature" || tt.name == "pki invalid signature" {
					if anno.IsSatisfied {
						t.Errorf("satisfied should be false")
					}
				}
			}
		})
	}
}
