/*******************************************************************************
 * Copyright 2022 Dell Inc.
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
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/project-alvarium/alvarium-sdk-go/internal/annotators"
	handler "github.com/project-alvarium/alvarium-sdk-go/internal/annotators/http/handler"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/config"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/contracts"
	"github.com/project-alvarium/alvarium-sdk-go/test"
)

func TestHttpPkiAnnotator_Do(t *testing.T) {
	b, err := ioutil.ReadFile("./test/config.json")
	if err != nil {
		t.Fatalf(err.Error())
	}

	var cfg config.SdkInfo
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		t.Fatalf(err.Error())
	}

	req, data, err := buildRequest(cfg.Signature)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Set up example signed data type for test purposes
	type testData struct {
		SignatureInput string
		Signature      string
	}

	// Tests the Signature signed by the assembler called by buildRequest
	t1 := testData{
		SignatureInput: req.Header.Get("Signature-Input"),
		Signature:      req.Header.Get("Signature"),
	}

	t2 := t1
	t2.SignatureInput = "\"@method\" \"@path\" \"@authority\" \"Content-Type\" \"Content-Length\";created=1646146637;keyid=\"public.key\";alg=\"invalid\""

	t3 := t1
	t3.SignatureInput = "\"@method\" \"@path\" \"@authority\" \"Content-Type\" \"Content-Length\";created=1646146637;keyid=\"invalid\";alg=\"ed25519\""

	t4 := t1
	t4.Signature = ""

	t5 := t1
	t5.Signature = "invalid"

	tests := []struct {
		name        string
		expectError bool
		data        testData
	}{
		{"pki annotation OK", false, t1},
		{"pki bad key type", true, t2},
		{"pki key not found", true, t3},
		{"pki empty signature", false, t4},
		{"pki invalid signature", false, t5},
	}

	for _, tt := range tests {
		req.Header.Set("Signature-Input", tt.data.SignatureInput)
		req.Header.Set("Signature", tt.data.Signature)

		ctx := context.WithValue(req.Context(), contracts.HttpRequestKey, req)

		t.Run(tt.name, func(t *testing.T) {
			pki := NewHttpPkiAnnotator(cfg)
			anno, err := pki.Do(ctx, data)
			test.CheckError(err, tt.expectError, tt.name, t)
			if err == nil {
				result, err := annotators.VerifySignature(cfg.Signature.PublicKey, anno)
				if err != nil {
					t.Error(err.Error())
				} else if !result {
					t.Error("signature not verified")
				}
				if tt.name == "pki empty signature" || tt.name == "pki invalid signature" {
					if anno.IsSatisfied {
						t.Errorf("satisfied should be false")
					}
				} else if tt.name == "pki annotation OK" {
					if !anno.IsSatisfied {
						t.Errorf("satisfied should be true")
					}
				}
			}
		})
	}
}

func buildRequest(keys config.SignatureInfo) (*http.Request, []byte, error) {
	type sample struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	t := sample{Key: "keyA", Value: "This is some test data"}
	b, _ := json.Marshal(t)

	req := httptest.NewRequest("POST", "/foo?param=value&foo=bar&baz=batman", bytes.NewReader(b))
	ticks := time.Now()
	now := ticks.String()
	req.Header = http.Header{
		"Host":           []string{"example.com"},
		"Date":           []string{now},
		"Content-Type":   []string{string(contracts.ContentTypeJSON)},
		"Content-Length": []string{strconv.FormatInt(req.ContentLength, 10)},
	}

	fields := []string{string(contracts.Method), string(contracts.Path), string(contracts.Authority), contracts.HttpContentType, contracts.ContentLength}
	handler := handler.NewEd25519RequestHandler(req)
	err := handler.AddSignatureHeaders(ticks, fields, keys)
	if err != nil {
		return nil, nil, err
	}
	return req, b, nil
}
