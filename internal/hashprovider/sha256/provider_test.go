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
package sha256

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
			expected: "2c26b46b68ffc68ff99b453c1d30413413422d706483bfa0f98a5e886266e7ae",
		},
		{
			name:     "text variation 2",
			data:     []byte("bar"),
			expected: "fcde2b2edba56bf408601fb721fe9b5c338d10ee429ea04fae5511b68fbf8fb9",
		},
		{
			name:     "text variation 3",
			data:     []byte("baz"),
			expected: "baa5a0964d3320fbc0c6a922140453c8513ea24ab8fd0577034804a967248096",
		},
		{
			name:     "byte sequence",
			data:     []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
			expected: "9a89c68c4c5e28b8c4a5567673d462fff515db46116f9900624d09c474f593fb",
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
