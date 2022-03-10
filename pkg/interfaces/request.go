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
package interfaces

import (
	"time"

	"github.com/project-alvarium/alvarium-sdk-go/pkg/config"
)

type RequestHandler interface {
	// AddSignatureHeaders takes time of creation of request, the fields to be taken into consideration
	// for the SignatureInput header, and the keys to be used in signing the seed.
	// Assembles the SignatureInput and Signature fields, then adds them to the request as headers.
	AddSignatureHeaders(ticks time.Time, fields []string, keys config.SignatureInfo) error
}
