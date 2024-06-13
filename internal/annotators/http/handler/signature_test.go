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

package http

import (
	"encoding/json"
	"fmt"
	"github.com/project-alvarium/alvarium-sdk-go/internal/signprovider/ecdsa/secp256k1"
	"github.com/project-alvarium/alvarium-sdk-go/internal/signprovider/ecdsa/x509"
	"github.com/project-alvarium/alvarium-sdk-go/internal/signprovider/ed25519"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/interfaces"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/project-alvarium/alvarium-sdk-go/pkg/config"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/contracts"
	"github.com/stretchr/testify/assert"
)

func TestHttpPkiAnnotator_Signature_AddSignatureHeaders(t *testing.T) {
	type sample struct {
		cfg       string
		signature interfaces.SignatureProvider
	}

	testCases := []sample{
		{
			cfg:       "./test/config.json",
			signature: ed25519.New(),
		},
		{
			cfg:       "./test/config-ecdsa-x509.json",
			signature: x509.New(),
		},
		{
			cfg:       "./test/config-ecdsa-sepc256k1.json",
			signature: secp256k1.New(),
		},
	}

	for _, tc := range testCases {
		b, err := os.ReadFile(tc.cfg)
		if err != nil {
			t.Fatalf(err.Error())
		}

		var cfg config.SdkInfo
		err = json.Unmarshal(b, &cfg)
		if err != nil {
			t.Fatalf(err.Error())
		}
		ticks := time.Now()
		now := ticks.String()
		req := httptest.NewRequest("POST", "http://www.example.com/foo?var1=&var2=2", nil)

		req.Header = http.Header{
			"Date":           []string{now},
			"Content-Type":   []string{string(contracts.ContentTypeJSON)},
			"Content-Length": []string{strconv.FormatInt(req.ContentLength, 10)},
		}

		fields := []string{string(contracts.Method), string(contracts.Path), string(contracts.Authority), contracts.HttpContentType, contracts.ContentLength}
		keys := cfg.Signature
		instance := NewSignatureRequestHandler(req, tc.signature)
		err = instance.AddSignatureHeaders(ticks, fields, keys)
		if err != nil {
			t.Error(err.Error())
		}

		t.Run("testing assembler signature input construction", func(t *testing.T) {
			expectedSignatureInput := fmt.Sprintf("\"@method\" \"@path\" \"@authority\" \"Content-Type\" \"Content-Length\";created=%s;keyid=\"%s\";alg=\"%s\";",
				strconv.FormatInt(ticks.Unix(), 10), filepath.Base(keys.PublicKey.Path), keys.PublicKey.Type)
			assert.Equal(t, expectedSignatureInput, req.Header.Get("Signature-Input"))
		})
	}
}
