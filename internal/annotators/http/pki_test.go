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
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/project-alvarium/alvarium-sdk-go/internal/annotators"
	"github.com/project-alvarium/alvarium-sdk-go/internal/signprovider/ed25519"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/contracts"
	"github.com/project-alvarium/alvarium-sdk-go/test"

	"github.com/project-alvarium/alvarium-sdk-go/pkg/config"
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

	badKeyType := cfg
	badKeyType.Signature.PublicKey.Type = "invalid"

	keyNotFound := cfg
	keyNotFound.Signature.PublicKey.Path = "/dev/null/private.key"

	req, data, err := buildRequest(cfg.Signature)
	if err != nil {
		t.Fatalf(err.Error())
	}

	ctx := context.WithValue(req.Context(), testRequest, req)

	//Building a new request with an empty signature
	req2 := req.Clone(req.Context())
	req2.Header.Set("Signature", "")
	ctx2 := context.WithValue(req2.Context(), testRequest, req2)

	tests := []struct {
		name        string
		ctx         context.Context
		cfg         config.SdkInfo
		expectError bool
	}{
		{"pki annotation OK", ctx, cfg, false},
		{"pki bad key type", ctx, badKeyType, true},
		{"pki key not found", ctx, keyNotFound, true},
		{"pki empty signature", ctx2, cfg, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pki := NewHttpPkiAnnotator(tt.cfg)
			anno, err := pki.Do(tt.ctx, data)
			test.CheckError(err, tt.expectError, tt.name, t)
			if err == nil {
				result, err := annotators.VerifySignature(tt.cfg.Signature.PublicKey, anno)
				if err != nil {
					t.Error(err.Error())
				} else if !result {
					t.Error("signature not verified")
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
	req.Header.Set("Host", "example.com")

	ticks := time.Now()
	now := ticks.String()
	req.Header.Set("Date", now)
	req.Header.Set(contentType, string(contracts.ContentTypeJSON))
	req.Header.Set(contentLength, strconv.FormatInt(req.ContentLength, 10))

	fields := []string{string(method), string(path), string(authority), contentType, contentLength}
	headerValue, signature, err := signRequest(ticks, fields, keys, req)

	req.Header.Set("Signature-Input", headerValue)
	req.Header.Set("Signature", signature)

	return req, b, err
}

func signRequest(ticks time.Time, fields []string, keys config.SignatureInfo, req *http.Request) (string, string, error) {
	headerValue := "" //This will be the value returned for populating the Signature-Input header
	inputValue := ""  //This will be the value used as input for the signature
	sigParams := ""   //This will be used to build the signatureParams field

	for i, f := range fields {
		headerValue += fmt.Sprintf("\"%s\"", f)
		switch f {
		case contentType:
			inputValue += fmt.Sprintf("\"%s\" %s", f, req.Header.Get(contentType))
		case contentLength:
			inputValue += fmt.Sprintf("\"%s\" %s", f, strconv.FormatInt(req.ContentLength, 10))
		case string(method):
			inputValue += fmt.Sprintf("\"%s\" %s", f, req.Method)
		case string(authority):
			inputValue += fmt.Sprintf("\"%s\" %s", f, req.Host)
		case string(scheme):
			scheme := strings.ToLower(strings.Split(req.Proto, "/")[0])
			inputValue += fmt.Sprintf("\"%s\" %s", f, scheme)
		case string(requestTarget):
			inputValue += fmt.Sprintf("\"%s\" %s", f, req.RequestURI)
		case string(path):
			inputValue += fmt.Sprintf("\"%s\" %s", f, req.URL.Path)
		case string(query):
			var query string = "?" + req.URL.RawQuery
			inputValue += fmt.Sprintf("\"%s\" %s", f, query)
		case string(queryParams):
			queryParamsRawMap := req.URL.Query()
			var queryParams []string
			for key, value := range queryParamsRawMap {
				b := new(bytes.Buffer)
				fmt.Fprintf(b, ";name=\"%s\": %s", key, value[0])
				queryParams = append(queryParams, b.String())
			}

			inputValue += fmt.Sprintf("\"%s\" %s", f, query)
		}

		sigParams += fmt.Sprintf("\"%s\"", f)
		inputValue += "\n"
		if i < len(fields)-1 {
			headerValue += " "
			sigParams += " "
		}
	}

	// fmt.Println("------------->", headerValue)

	//TODO
	/*tail := fmt.Sprintf("\"%s\": (%s);created=%s;keyid=\"%s\";alg=\"%s\";", string(signatureParams), sigParams, strconv.FormatInt(ticks.Unix(),10),
		 filepath.Base(keys.PublicKey.Path), keys.PublicKey.Type)
	 //headerValue += tail
	 inputValue += tail*/

	// fmt.Println("HEADER: " + headerValue)
	// fmt.Println("INPUT: " + inputValue)
	signer := ed25519.New()
	prv, err := ioutil.ReadFile(keys.PrivateKey.Path)
	if err != nil {
		return "", "", err
	}

	signature := signer.Sign(prv, []byte(inputValue))
	//fmt.Println("SIGNATURE: " + signature)
	return headerValue, signature, nil
}
