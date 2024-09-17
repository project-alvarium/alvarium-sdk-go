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
package sha256

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// newSUT returns a new system under test.
func newSUT() *provider {
	return New()
}

// TestProvider_Derive tests provider.Derive.
func TestProvider_Derive(t *testing.T) {
	cases := []struct {
		name     string
		data     []byte
		expected string
	}{
		{
			name:     "text variation 1",
			data:     []byte("foo"),
			expected: "2C26B46B68FFC68FF99B453C1D30413413422D706483BFA0F98A5E886266E7AE",
		},
		{
			name:     "text variation 2",
			data:     []byte("bar"),
			expected: "FCDE2B2EDBA56BF408601FB721FE9B5C338D10EE429EA04FAE5511B68FBF8FB9",
		},
		{
			name:     "text variation 3",
			data:     []byte("baz"),
			expected: "BAA5A0964D3320FBC0C6A922140453C8513EA24AB8FD0577034804A967248096",
		},
		{
			name:     "byte sequence",
			data:     []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
			expected: "9A89C68C4C5E28B8C4A5567673D462FFF515DB46116F9900624D09C474F593FB",
		},
	}

	for i := range cases {
		t.Run(
			cases[i].name,
			func(t *testing.T) {
				sut := newSUT()

				result := sut.Derive(cases[i].data)

				assert.Equal(t, cases[i].expected, result)
			},
		)
	}
}
