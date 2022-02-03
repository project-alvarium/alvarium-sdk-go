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
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/project-alvarium/alvarium-sdk-go/pkg/config"
	"github.com/project-alvarium/alvarium-sdk-go/test"
)

func TestTpmAnnotator_Do(t *testing.T) {
	b, err := ioutil.ReadFile("../../test/res/config.json")
	if err != nil {
		t.Fatalf(err.Error())
	}

	var cfg config.SdkInfo
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		t.Fatalf(err.Error())
	}

	badHashType := cfg
	badHashType.Hash.Type = "invalid"

	badKeyType := cfg
	badKeyType.Signature.PrivateKey.Type = "invalid"

	keyNotFound := cfg
	keyNotFound.Signature.PrivateKey.Path = "/dev/null/private.key"

	rndString := test.FactoryRandomFixedLengthString(1024, test.AlphanumericCharset)
	tests := []struct {
		name        string
		data        string
		cfg         config.SdkInfo
		expectError bool
	}{
		{"tpm annotation OK", rndString, cfg, false},
		{"tpm bad hash type", rndString, badHashType, false}, // returns "none" hash type
		{"tpm bad key type", rndString, badKeyType, true},
		{"tpm key not found", rndString, keyNotFound, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tpm := NewTpmAnnotator(tt.cfg)
			anno, err := tpm.Do(context.Background(), []byte(tt.data))
			test.CheckError(err, tt.expectError, tt.name, t)
			if err == nil {
				result, err := VerifySignature(tt.cfg.Signature.PublicKey, anno)
				if err != nil {
					t.Error(err.Error())
				} else if !result {
					t.Error("signature not verified")
				}
			}
		})
	}
}
