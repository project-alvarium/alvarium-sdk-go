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
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/project-alvarium/alvarium-sdk-go/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestHttpPkiAnnotator_RequestParser(t *testing.T) {
	b, err := ioutil.ReadFile("./test/config.json")
	if err != nil {
		t.Fatalf(err.Error())
	}

	var cfg config.SdkInfo
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		t.Fatalf(err.Error())
	}

	base := httptest.NewRequest("POST", "/foo?var1=1&var2=2", nil)

	base.Header.Set("Host", "example.com")
	base.Header.Set("Date", "Tue, 20 Apr 2021 02:07:55 GMT")
	base.Header.Set("Content-Type", "application/json")
	base.Header.Set("Content-Length", "18")
	base.Header.Set("Signature", "whatever")

	tests := []struct {
		name           string
		signatureInput string
		expectedSeed   string
		expectError    bool
	}{
		{"testing integeration of all Signature-Input fields",
			"\"date\" \"@method\" \"@path\" \"@authority\" \"content-type\" \"content-length\" \"@query-params\" \"@query\" ",
			"\"date\" Tue, 20 Apr 2021 02:07:55 GMT\n\"@method\" POST\n\"@path\" /foo\n\"@authority\" example.com\n\"content-type\" application/json\n\"content-length\" 18\n\"@query-params\";name=\"var1\": 1\n\"@query-params\";name=\"var2\": 2\n\"@query\" ?var1=1&var2=2\n", false},

		{"testing @method ", "\"@method\"", "\"@method\" POST\n", false},
		{"testing @authority", "\"@authority\"", "\"@authority\" example.com\n", false},

		{"testing @scheme", "\"@scheme\"", "\"@scheme\" http\n", false},
		{"testing @request-target", "\"@request-target\"", "\"@request-target\" /foo?var1=1&var2=2\n", false},

		{"testing @path", "\"@path\"", "\"@path\" /foo\n", false},

		{"testing @query", "\"@query\"", "\"@query\" ?var1=1&var2=2\n", false},
		{"testing @query-params", "\"@query-params\"", "\"@query-params\";name=\"var1\": 1\n\"@query-params\";name=\"var2\": 2\n", false},

		{"testing non-existant derived component", "\"@x-test\"", "", true},
		{"testing non-existant  header field", "\"x-test\"", "", true},
	}

	var seed string
	for _, tt := range tests {

		req := base.Clone(base.Context())
		req.Header.Set("Signature-Input", tt.signatureInput)

		t.Run(tt.name, func(t *testing.T) {
			seed, err = requestParser(req)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.expectedSeed, seed)
			}
		})

	}

}
