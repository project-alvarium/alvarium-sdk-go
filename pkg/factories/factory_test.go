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
	"github.com/DyrellC/alvarium-sdk-go/pkg/config"
	"github.com/DyrellC/alvarium-sdk-go/pkg/contracts"
	"github.com/DyrellC/alvarium-sdk-go/test"
	logConfig "github.com/project-alvarium/provider-logging/pkg/config"
	logFactory "github.com/project-alvarium/provider-logging/pkg/factories"
	"github.com/project-alvarium/provider-logging/pkg/logging"
	"testing"
)

func TestStreamProviderFactory(t *testing.T) {
	logInfo := logConfig.LoggingInfo{MinLogLevel: logging.InfoLevel}
	logger := logFactory.NewLogger(logInfo)

	pass := config.StreamInfo{
		Type:   contracts.IotaStream,
		Config: config.IotaStreamConfig{},
	}

	pass2 := config.StreamInfo{
		Type:   contracts.MockStream,
		Config: config.IotaStreamConfig{},
	}

	pass3 := config.StreamInfo{
		Type:   contracts.MqttStream,
		Config: config.MqttConfig{},
	}

	fail := config.StreamInfo{
		Type:   "invalid",
		Config: config.IotaStreamConfig{},
	}

	tests := []struct {
		name         string
		providerType config.StreamInfo
		expectError  bool
	}{
		{"valid iota type", pass, false},
		{"valid mock type", pass2, false},
		{"valid mqtt type", pass3, false},
		{"invalid random type", fail, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewStreamProvider(tt.providerType, logger)
			test.CheckError(err, tt.expectError, tt.name, t)
		})
	}
}

func TestAnnotatorFactory(t *testing.T) {
	cfg := config.SdkInfo{}

	tests := []struct {
		name        string
		cfg         config.SdkInfo
		key         contracts.AnnotationType
		expectError bool
	}{
		{"valid pki type", cfg, contracts.AnnotationPKI, false},
		{"valid tpm type", cfg, contracts.AnnotationTPM, false},
		{"invalid annotator type", cfg, "invalid", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewAnnotator(tt.key, tt.cfg)
			test.CheckError(err, tt.expectError, tt.name, t)
		})
	}
}
