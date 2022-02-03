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

/*
	TODO:


*/

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

	base := httptest.NewRequest("POST", "/foo", nil)

	base.Header.Set("Host", "example.com")
	base.Header.Set("Date", "Tue, 20 Apr 2021 02:07:55 GMT")
	base.Header.Set("Content-Type", "application/json")
	base.Header.Set("Content-Length", "18")
	base.Header.Set("Signature", "whatever")

	req2 := base.Clone(base.Context())

	req2.Header.Set("Signature-Input", "\"date\" \"@method\" \"@path\" \"@authority\" \"content-type\" \"content-length\"")
	seed, _ := requestParser(req2)

	t.Run("parser_test_should_succeed_if_equal", func(t *testing.T) {
		expectedSeed := "\"date\" Tue, 20 Apr 2021 02:07:55 GMT\n\"@method\" POST\n\"@path\" /foo\n\"@authority\" example.com\n\"content-type\" application/json\n\"content-length\" 18\n"
		assert.Equal(t, expectedSeed, seed)
	})

	t.Run("parser_test_should_fail_if_not_equal", func(t *testing.T) {
		expectedSeed := "\"date\" GMT\n\"@method\" POST\n\"@path\" /foo\n\"@authority\" example.com\n\"content-type\" application/json\n\"content-length\" 18\n"
		assert.NotEqual(t, expectedSeed, seed)
	})

	req3 := base.Clone(base.Context())
	req3.Header.Set("Signature-Input", "\"@method\"")

	seed, _ = requestParser(req3)

	t.Run("parser_test_@method_should_succeed", func(t *testing.T) {
		expectedSeed := "\"@method\" POST\n"
		assert.Equal(t, expectedSeed, seed)
	})

	req4 := base.Clone(base.Context())
	req4.Header.Set("Signature-Input", "\"@authority\"")

	seed, _ = requestParser(req4)

	t.Run("parser_test_@authority_should_succeed", func(t *testing.T) {
		expectedSeed := "\"@authority\" example.com\n"
		assert.Equal(t, expectedSeed, seed)
	})

	req5 := base.Clone(base.Context())
	req5.Header.Set("Signature-Input", "\"@x-test\"")

	seed, err = requestParser(req5)

	t.Run("parser_test_should_fail_unhandled_speciality_component", func(t *testing.T) {
		assert.Error(t, err)
	})

	req6 := base.Clone(base.Context())
	req6.Header.Set("Signature-Input", "\"x-test\"")

	seed, err = requestParser(req6)

	t.Run("parser_test_should_fail_unhandled_header_field", func(t *testing.T) {
		assert.Error(t, err)
	})

	// TODO:  the rest of the methods

}
