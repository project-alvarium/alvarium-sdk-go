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
package md5

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
			expected: "ACBD18DB4CC2F85CEDEF654FCCC4A4D8",
		},
		{
			name:     "text variation 2",
			data:     []byte("bar"),
			expected: "37B51D194A7513E45B56F6524F2D51F2",
		},
		{
			name:     "text variation 3",
			data:     []byte("baz"),
			expected: "73FEFFA4B7F6BB68E44CF984C85F6E88",
		},
		{
			name:     "byte sequence",
			data:     []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
			expected: "7F63CB6D067972C3F34F094BB7E776A8",
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
