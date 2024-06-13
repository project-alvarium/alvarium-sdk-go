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

package factories

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/project-alvarium/alvarium-sdk-go/internal/annotators"
	httpAnnotators "github.com/project-alvarium/alvarium-sdk-go/internal/annotators/http"
	handler "github.com/project-alvarium/alvarium-sdk-go/internal/annotators/http/handler"
	"github.com/project-alvarium/alvarium-sdk-go/internal/console"
	"github.com/project-alvarium/alvarium-sdk-go/internal/hashprovider/md5"
	"github.com/project-alvarium/alvarium-sdk-go/internal/hashprovider/none"
	"github.com/project-alvarium/alvarium-sdk-go/internal/hashprovider/sha256"
	"github.com/project-alvarium/alvarium-sdk-go/internal/hedera"
	"github.com/project-alvarium/alvarium-sdk-go/internal/mock"
	"github.com/project-alvarium/alvarium-sdk-go/internal/mqtt"
	"github.com/project-alvarium/alvarium-sdk-go/internal/signprovider/ecdsa/secp256k1"
	"github.com/project-alvarium/alvarium-sdk-go/internal/signprovider/ecdsa/x509"
	"github.com/project-alvarium/alvarium-sdk-go/internal/signprovider/ed25519"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/config"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/contracts"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/interfaces"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/logging"
)

func NewStreamProvider(cfg config.StreamInfo, logger interfaces.Logger) (interfaces.StreamProvider, error) {
	switch cfg.Type {
	case contracts.MockStream:
		info, ok := cfg.Config.(config.MockStreamConfig)
		if !ok {
			return nil, errors.New("invalid cast for MockStream")
		}
		return mock.NewMockPublisher(info, logger), nil
	case contracts.MqttStream:
		info, ok := cfg.Config.(config.MqttConfig)
		if !ok {
			return nil, errors.New("invalid cast for MqttStream")
		}
		return mqtt.NewMqttPublisher(info, logger), nil
	case contracts.ConsoleStream:
		return console.NewConsolePublisher(logger), nil
	case contracts.HederaStream:
		info, ok := cfg.Config.(config.HederaConfig)
		if !ok {
			return nil, errors.New("invalid cast for HederaStream")
		}
		return hedera.NewHederaPublisher(info, logger)
	default:
		return nil, fmt.Errorf("unrecognized config Type value %s", cfg.Type)
	}
}

func NewHashProvider(hash contracts.HashType) (interfaces.HashProvider, error) {
	switch hash {
	case contracts.MD5Hash:

		return md5.New(), nil
	case contracts.SHA256Hash:

		return sha256.New(), nil
	case contracts.NoHash:
		return none.New(), nil
	default:
		return nil, fmt.Errorf("unrecognized hash type value %s", hash)
	}
}

// NewSignatureProvider instantiates a signature provider based on the desired key algorithm
//
// The current working assumption is that all nodes within a Data Confidence Fabric will use the same algorithm
// to generate their identity keys. If later there's a good reason provided as to why this might be heterogeneous,
// the existing implementation around signatures will need to change
func NewSignatureProvider(k contracts.KeyAlgorithm) (interfaces.SignatureProvider, error) {
	switch k {
	case contracts.KeyEd25519:
		return ed25519.New(), nil
	case contracts.KeyEcdsaX509:
		return x509.New(), nil
	case contracts.KeyEcdsaSecp256k1:
		return secp256k1.New(), nil
	default:
		return nil, fmt.Errorf("unrecognized key algorithm value %s", k)
	}
}

func NewAnnotator(kind contracts.AnnotationType, cfg config.SdkInfo) (interfaces.Annotator, error) {
	h, err := NewHashProvider(cfg.Hash.Type)
	if err != nil {
		return nil, err
	}

	s, err := NewSignatureProvider(cfg.Signature.PrivateKey.Type)
	if err != nil {
		return nil, err
	}

	var a interfaces.Annotator
	switch kind {
	case contracts.AnnotationSource:
		a = annotators.NewSourceAnnotator(cfg, h, s)
	case contracts.AnnotationTPM:
		a = annotators.NewTpmAnnotator(cfg, h, s)
	case contracts.AnnotationPKI:
		a = annotators.NewPkiAnnotator(cfg, h, s)
	case contracts.AnnotationPKIHttp:
		a = httpAnnotators.NewHttpPkiAnnotator(cfg, h, s)
	case contracts.AnnotationTLS:
		a = annotators.NewTlsAnnotator(cfg, h, s)
	default:
		return nil, fmt.Errorf("unrecognized AnnotationType %s", kind)
	}
	return a, nil
}

func NewRequestHandler(request *http.Request, keys config.SignatureInfo) (interfaces.RequestHandler, error) {
	var r interfaces.RequestHandler

	switch keys.PrivateKey.Type {
	case contracts.KeyEd25519:
		r = handler.NewSignatureRequestHandler(request, ed25519.New())
	case contracts.KeyEcdsaX509:
		r = handler.NewSignatureRequestHandler(request, x509.New())
	case contracts.KeyEcdsaSecp256k1:
		r = handler.NewSignatureRequestHandler(request, secp256k1.New())
	default:
		return nil, fmt.Errorf("unrecognized Key Type %s", keys.PrivateKey.Type)
	}
	return r, nil
}

func NewLogger(cfg config.LoggingInfo) interfaces.Logger {
	return logging.NewConsoleLogger(cfg)
}
