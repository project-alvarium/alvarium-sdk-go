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

package http

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/project-alvarium/alvarium-sdk-go/pkg/config"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/interfaces"
)

type sigRequestHandler struct {
	Request   *http.Request
	Signature interfaces.SignatureProvider
}

func NewSignatureRequestHandler(request *http.Request, signature interfaces.SignatureProvider) interfaces.RequestHandler {
	instance := sigRequestHandler{
		Request:   request,
		Signature: signature,
	}
	return &instance
}

func (h *sigRequestHandler) AddSignatureHeaders(ticks time.Time, fields []string, keys config.SignatureInfo) error {
	var headerValue strings.Builder //This will be the value returned for populating the Signature-Input header
	inputValue := ""                //This will be the value used as input for the signature

	for i, f := range fields {
		headerValue.WriteString(fmt.Sprintf("\"%s\"", f))
		if i < len(fields)-1 {
			headerValue.WriteString(" ")
		}
	}

	tail := fmt.Sprintf(";created=%s;keyid=\"%s\";alg=\"%s\";", strconv.FormatInt(ticks.Unix(), 10),
		filepath.Base(keys.PublicKey.Path), keys.PublicKey.Type)

	headerValue.WriteString(tail)

	h.Request.Header.Set("Signature-Input", headerValue.String())

	parsed, err := ParseSignature(h.Request)
	if err != nil {
		return err
	}
	inputValue = parsed.Seed

	signature, err := h.Signature.Sign(keys.PrivateKey, []byte(inputValue))

	h.Request.Header.Set("Signature", signature)
	return nil
}
