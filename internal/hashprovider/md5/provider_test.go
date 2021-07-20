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
package md5

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
			expected: "acbd18db4cc2f85cedef654fccc4a4d8",
		},
		{
			name:     "text variation 2",
			data:     []byte("bar"),
			expected: "37b51d194a7513e45b56f6524f2d51f2",
		},
		{
			name:     "text variation 3",
			data:     []byte("baz"),
			expected: "73feffa4b7f6bb68e44cf984c85f6e88",
		},
		{
			name:     "byte sequence",
			data:     []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
			expected: "7f63cb6d067972c3f34f094bb7e776a8",
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
