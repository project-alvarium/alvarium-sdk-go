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
	"github.com/project-alvarium/alvarium-sdk-go/pkg/contracts"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/interfaces"
	"github.com/project-alvarium/alvarium-sdk-go/test"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestTlsAnnotator_Base(t *testing.T) {
	b, err := os.ReadFile("../../test/res/config.json")
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
	signer := ed25519.New()
	h := hash256.New()
	tests := []struct {
		name        string
		data        string
		cfg         config.SdkInfo
		h           interfaces.HashProvider
		s           interfaces.SignatureProvider
		expectError bool
	}{
		{"tls annotation OK", rndString, cfg, h, signer, false},
		{"tls bad hash type", rndString, badHashType, h, signer, false}, // returns "none" hash type
		{"tls key not found", rndString, keyNotFound, h, signer, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tls := NewTlsAnnotator(tt.cfg, tt.h, tt.s)
			anno, err := tls.Do(context.Background(), []byte(tt.data))
			test.CheckError(err, tt.expectError, tt.name, t)
			if err == nil {
				result, err := VerifySignature(tt.cfg.Signature.PublicKey, tt.s, anno)
				if err != nil {
					t.Error(err.Error())
				} else if !result {
					t.Error("signature not verified")
				}
			}
		})
	}
}

func TestTlsAnnotator_ServeTLS(t *testing.T) {
	b, err := os.ReadFile("../../test/res/config.json")
	if err != nil {
		t.Fatalf(err.Error())
	}

	var cfg config.SdkInfo
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		t.Fatalf(err.Error())
	}
	tls := NewTlsAnnotator(cfg, hash256.New(), ed25519.New())

	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), contracts.AnnotationTLS, r.TLS)
		anno, err := tls.Do(ctx, []byte(test.FactoryRandomFixedLengthString(1024, test.AlphanumericCharset)))
		if err != nil {
			t.Error(err)
		}
		b, _ := json.Marshal(anno)
		w.Write(b)
	}))
	defer ts.Close()

	client := ts.Client()
	result, err := client.Get(ts.URL)
	if err != nil {
		t.Error(err)
	} else {
		defer result.Body.Close()
		var a contracts.Annotation
		body, _ := io.ReadAll(result.Body)
		json.Unmarshal(body, &a)
		if !a.IsSatisfied {
			t.Error("annotation.IsSatisfied should be true")
		}
	}
}

func TestTlsAnnotator_Serve(t *testing.T) {
	b, err := os.ReadFile("../../test/res/config.json")
	if err != nil {
		t.Fatalf(err.Error())
	}

	var cfg config.SdkInfo
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		t.Fatalf(err.Error())
	}
	tls := NewTlsAnnotator(cfg, hash256.New(), ed25519.New())

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), contracts.AnnotationTLS, r.TLS)
		anno, err := tls.Do(ctx, []byte(test.FactoryRandomFixedLengthString(1024, test.AlphanumericCharset)))
		if err != nil {
			t.Error(err)
		}
		b, _ := json.Marshal(anno)
		w.Write(b)
	}))
	defer ts.Close()

	client := ts.Client()
	result, err := client.Get(ts.URL)
	if err != nil {
		t.Error(err)
	} else {
		defer result.Body.Close()
		var a contracts.Annotation
		body, _ := io.ReadAll(result.Body)
		json.Unmarshal(body, &a)
		if a.IsSatisfied {
			t.Error("annotation.IsSatisfied should be false")
		}
	}
}
