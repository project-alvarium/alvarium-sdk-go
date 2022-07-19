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

package factories

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/project-alvarium/alvarium-sdk-go/internal/annotators"
	httpAnnotators "github.com/project-alvarium/alvarium-sdk-go/internal/annotators/http"
	handler "github.com/project-alvarium/alvarium-sdk-go/internal/annotators/http/handler"
	"github.com/project-alvarium/alvarium-sdk-go/internal/iota"
	"github.com/project-alvarium/alvarium-sdk-go/internal/mock"
	"github.com/project-alvarium/alvarium-sdk-go/internal/mqtt"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/config"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/contracts"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/interfaces"
	logInterface "github.com/project-alvarium/provider-logging/pkg/interfaces"
)

func NewStreamProvider(cfg config.StreamInfo, logger logInterface.Logger) (interfaces.StreamProvider, error) {
	switch cfg.Type {
	case contracts.IotaStream:
		info, ok := cfg.Config.(config.IotaStreamConfig)
		if !ok {
			return nil, errors.New("invalid cast for IotaStream")
		}
		return iota.NewIotaPublisher(info, logger)
	case contracts.MockStream:
		info, ok := cfg.Config.(config.IotaStreamConfig)
		if !ok {
			return nil, errors.New("invalid cast for MockStream")
		}
		return mock.NewMockPublisher(info, logger)
	case contracts.MqttStream:
		info, ok := cfg.Config.(config.MqttConfig)
		if !ok {
			return nil, errors.New("invalid cast for MqttStream")
		}
		return mqtt.NewMqttPublisher(info, logger), nil
	default:
		return nil, fmt.Errorf("unrecognized config Type value %s", cfg.Type)
	}
}

func NewAnnotator(kind contracts.AnnotationType, cfg config.SdkInfo) (interfaces.Annotator, error) {
	var a interfaces.Annotator
	switch kind {
	case contracts.AnnotationSource:
		a = annotators.NewSourceAnnotator(cfg)
	case contracts.AnnotationTPM:
		a = annotators.NewTpmAnnotator(cfg)
	case contracts.AnnotationPKI:
		a = annotators.NewPkiAnnotator(cfg)
	case contracts.AnnotationPKIHttp:
		a = httpAnnotators.NewHttpPkiAnnotator(cfg)
	case contracts.AnnotationTLS:
		a = annotators.NewTlsAnnotator(cfg)
	default:
		return nil, fmt.Errorf("unrecognized AnnotationType %s", kind)
	}
	return a, nil
}

func NewRequestHandler(request *http.Request, keys config.SignatureInfo) (interfaces.RequestHandler, error) {
	var r interfaces.RequestHandler

	switch keys.PrivateKey.Type {
	case contracts.KeyEd25519:
		r = handler.NewEd25519RequestHandler(request)
	default:
		return nil, fmt.Errorf("unrecognized Key Type %s", keys.PrivateKey.Type)
	}
	return r, nil
}
