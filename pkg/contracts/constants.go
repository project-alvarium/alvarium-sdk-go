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
package contracts

type ContentType string

const (
	ContentTypeJSON ContentType = "application/json"
)

type HashType string

const (
	MD5Hash    HashType = "md5"
	SHA256Hash HashType = "sha256"
	NoHash     HashType = "none"
)

func (t HashType) Validate() bool {
	if t == MD5Hash || t == SHA256Hash || t == NoHash {
		return true
	}
	return false
}

type KeyAlgorithm string

const (
	KeyEd25519 KeyAlgorithm = "ed25519"
)

func (k KeyAlgorithm) Validate() bool {
	if k == KeyEd25519 {
		return true
	}
	return false
}

type StreamType string

const (
	IotaStream    StreamType = "iota"
	MockStream    StreamType = "mock"
	MqttStream    StreamType = "mqtt"
	PravegaStream StreamType = "pravega" // Currently unsupported but indicating extension point
)

func (t StreamType) Validate() bool {
	if t == IotaStream || t == MockStream || t == MqttStream || t == PravegaStream {
		return true
	}
	return false
}

type AnnotationType string

const (
	AnnotationPKI     AnnotationType = "pki"
	AnnotationPKIHttp AnnotationType = "pki-http"
	AnnotationSource  AnnotationType = "src"
	AnnotationTLS     AnnotationType = "tls"
	AnnotationTPM     AnnotationType = "tpm"
)

func (t AnnotationType) Validate() bool {
	if t == AnnotationPKI || t == AnnotationTLS || t == AnnotationTPM || t == AnnotationSource {
		return true
	}
	return false
}

type SpecialtyComponent string

const (
	Method        SpecialtyComponent = "@method"
	Authority     SpecialtyComponent = "@authority"
	Scheme        SpecialtyComponent = "@scheme"
	RequestTarget SpecialtyComponent = "@request-target"
	Path          SpecialtyComponent = "@path"
	Query         SpecialtyComponent = "@query"
	QueryParams   SpecialtyComponent = "@query-params"
)

const (
	// HttpRequestKey is the key used to reference the value within the incoming Context that corresponds to the request we need to validate.
	HttpRequestKey  string = "HttpRequestKey"
	ContentLength   string = "Content-Length"
	HttpContentType string = "Content-Type"
)

func (s SpecialtyComponent) Validate() bool {
	if s == Method || s == Authority || s == Scheme || s == RequestTarget || s == Path || s == Query || s == QueryParams {
		return true
	}
	return false
}
