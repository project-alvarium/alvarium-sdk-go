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
package none

import (
	"fmt"
	"strings"
	"testing"

	"github.com/project-alvarium/alvarium-sdk-go/test"
	"github.com/stretchr/testify/assert"
)

// newSUT returns a new system under test.
func newSUT() *provider {
	return New()
}

// TestProvider_Derive tests provider.Derive.
func TestProvider_Derive(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Run(
			"variation "+fmt.Sprint(i),
			func(t *testing.T) {
				data := test.FactoryRandomFixedLengthString(64, test.AlphanumericCharset)
				sut := newSUT()

				result := sut.Derive([]byte(data))

				assert.Equal(t, strings.ToUpper(data), result)
			},
		)
	}
}
