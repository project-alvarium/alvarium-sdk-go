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
package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/project-alvarium/alvarium-sdk-go/internal/annotators"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/config"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/contracts"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/factories"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/interfaces"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/message"
	logInterface "github.com/project-alvarium/provider-logging/pkg/interfaces"
	"github.com/project-alvarium/provider-logging/pkg/logging"
	"sync"
)

type sdk struct {
	annotators []interfaces.Annotator
	cfg        config.SdkInfo
	stream     interfaces.StreamProvider
	logger     logInterface.Logger
}

func NewSdk(annotators []interfaces.Annotator, cfg config.SdkInfo, logger logInterface.Logger) interfaces.Sdk {
	instance := sdk{
		annotators: annotators,
		cfg:        cfg,
		logger:     logger,
	}
	return &instance
}

func (s *sdk) BootstrapHandler(ctx context.Context, wg *sync.WaitGroup) bool {
	stream, err := factories.NewStreamProvider(s.cfg.Stream, s.logger)
	if err != nil {
		s.logger.Error(err.Error())
		return false
	}
	s.stream = stream
	//Connect to stream provider
	err = s.stream.Connect()
	if err != nil {
		s.logger.Error(err.Error())
		return false
	}
	s.logger.Write(logging.DebugLevel, "stream provider connection successful")

	wg.Add(1)
	go func() { // Graceful shutdown
		defer wg.Done()

		<-ctx.Done()
		s.logger.Write(logging.InfoLevel, "shutdown received")
		s.stream.Close()
	}()
	return true
}

func (s *sdk) Create(ctx context.Context, data []byte) {
	var list contracts.AnnotationList

	for _, a := range s.annotators {
		annotation, err := a.Do(ctx, data)
		if err != nil {
			s.logger.Error(err.Error())
			return
		}
		list.Items = append(list.Items, annotation)
	}

	b, _ := json.Marshal(list)
	wrap := message.PublishWrapper{
		Action:      message.ActionCreate,
		MessageType: fmt.Sprintf("%T", list),
		Content:     b,
	}
	err := s.stream.Publish(wrap)
	if err != nil {
		s.logger.Error(err.Error())
	}
}

func (s *sdk) Mutate(ctx context.Context, old, new []byte) {
	src := annotators.NewSourceAnnotator(s.cfg)
	a, err := src.Do(ctx, old)

	var list contracts.AnnotationList
	list.Items = append(list.Items, a)

	for _, a := range s.annotators {
		annotation, err := a.Do(ctx, new)
		if err != nil {
			s.logger.Error(err.Error())
			return
		}
		if annotation.Kind != contracts.AnnotationTLS {
			list.Items = append(list.Items, annotation)
		}
	}

	b, _ := json.Marshal(list)
	wrap := message.PublishWrapper{
		Action:      message.ActionMutate,
		MessageType: fmt.Sprintf("%T", list),
		Content:     b,
	}
	err = s.stream.Publish(wrap)
	if err != nil {
		s.logger.Error(err.Error())
	}
}

func (s *sdk) Transit(ctx context.Context, data []byte) {
	var list contracts.AnnotationList

	for _, a := range s.annotators {
		annotation, err := a.Do(ctx, data)
		if err != nil {
			s.logger.Error(err.Error())
			return
		}
		list.Items = append(list.Items, annotation)
	}

	b, _ := json.Marshal(list)
	wrap := message.PublishWrapper{
		Action:      message.ActionTransit,
		MessageType: fmt.Sprintf("%T", list),
		Content:     b,
	}
	err := s.stream.Publish(wrap)
	if err != nil {
		s.logger.Error(err.Error())
	}
}

func (s *sdk) Publish(ctx context.Context, data []byte) {
	var list contracts.AnnotationList

	for _, a := range s.annotators {
		annotation, err := a.Do(ctx, data)
		if err != nil {
			s.logger.Error(err.Error())
			return
		}
		list.Items = append(list.Items, annotation)
	}

	b, _ := json.Marshal(list)
	wrap := message.PublishWrapper{
		Action:      message.ActionPublish,
		MessageType: fmt.Sprintf("%T", list),
		Content:     b,
	}
	err := s.stream.Publish(wrap)
	if err != nil {
		s.logger.Error(err.Error())
	}
}