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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpPkiAnnotator_RequestParser(t *testing.T) {
	req := httptest.NewRequest("POST", "/foo?var1=&var2=2", nil)
	req.Header = http.Header{
		"Host":            []string{"example.com"},
		"Date":            []string{"Tue, 20 Apr 2021 02:07:55 GMT"},
		"Content-Type":    []string{"application/json"},
		"Content-Length":  []string{"18"},
		"Signature-Input": []string{""},
		"Signature":       []string{"whatever"},
	}

	seedTests := []struct {
		name           string
		signatureInput string
		expectedSeed   string
		expectError    bool
	}{
		{"testing integeration of all Signature-Input fields",
			"\"date\" \"@method\" \"@path\" \"@authority\" \"content-type\" \"content-length\" \"@query-params\" \"@query\";created=1644758607;keyid=\"public.key\";alg=\"ed25519\";",
			"\"date\" Tue, 20 Apr 2021 02:07:55 GMT\n\"@method\" POST\n\"@path\" /foo\n\"@authority\" example.com\n\"content-type\" application/json\n\"content-length\" 18\n\"@query-params\";name=\"var1\": \n\"@query-params\";name=\"var2\": 2\n\"@query\" ?var1=&var2=2\n;created=1644758607;keyid=\"public.key\";alg=\"ed25519\";", false},

		{"testing @method", "\"@method\";", "\"@method\" POST\n;", false},
		{"testing @authority", "\"@authority\";", "\"@authority\" example.com\n;", false},

		{"testing @scheme", "\"@scheme\";", "\"@scheme\" http\n;", false},
		{"testing @request-target", "\"@request-target\";", "\"@request-target\" /foo?var1=&var2=2\n;", false},

		{"testing @path", "\"@path\";", "\"@path\" /foo\n;", false},

		{"testing @query", "\"@query\";", "\"@query\" ?var1=&var2=2\n;", false},
		{"testing @query-params", "\"@query-params\";", "\"@query-params\";name=\"var1\": \n\"@query-params\";name=\"var2\": 2\n;", false},

		{"testing non-existant derived component", "\"@x-test\";", "", true},
		{"testing non-existant header field", "\"x-test\";", "", true},
	}

	for _, tt := range seedTests {
		req.Header.Set("Signature-Input", tt.signatureInput)

		t.Run(tt.name, func(t *testing.T) {
			signatureInfo, err := ParseSignature(req)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.expectedSeed, signatureInfo.Seed)
			}
		})

	}

	req.Header.Set("Signature-Input", "\"@query\";created=1644758607;keyid=\"public.key\";alg=\"ed25519\";")
	parsed, err := ParseSignature(req)
	if err != nil {
		t.Error(err.Error())
	}

	t.Run("testing signature", func(t *testing.T) {
		assert.Equal(t, "whatever", parsed.Signature)
	})

	t.Run("testing keyid", func(t *testing.T) {
		assert.Equal(t, "public.key", parsed.Keyid)
	})

	t.Run("testing algorithm", func(t *testing.T) {
		assert.Equal(t, "ed25519", parsed.Algorithm)
	})
}
