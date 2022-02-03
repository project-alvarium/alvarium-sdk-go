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

type specialtyComponent string

const (
	method        specialtyComponent = "@method"
	authority     specialtyComponent = "@authority"
	scheme        specialtyComponent = "@scheme"
	requestTarget specialtyComponent = "@request-target"
	path          specialtyComponent = "@path"
	query         specialtyComponent = "@query"
	queryParams   specialtyComponent = "@query-params"
)

const (
	contentLength string = "Content-Length"
	contentType   string = "Content-Type"
	testRequest   string = "testRequest"
)

func (s specialtyComponent) Validate() bool {
	if s == method || s == authority || s == scheme || s == requestTarget || s == path || s == query || s == queryParams {
		return true
	}
	return false
}
