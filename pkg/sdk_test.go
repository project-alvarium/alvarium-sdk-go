/*******************************************************************************
 * Copyright 2023 Dell Inc.
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
	"log/slog"
	"os"
	"sync"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/project-alvarium/alvarium-sdk-go/pkg/config"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/contracts"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/factories"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/interfaces"
)

func TestNewSdkJson(t *testing.T) {
	logger := factories.NewLogger(config.LoggingInfo{MinLogLevel: slog.LevelInfo})

	b, err := os.ReadFile("../test/res/config.json")
	if err != nil {
		t.Fatalf(err.Error())
	}

	var cfg config.SdkInfo
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		t.Fatalf(err.Error())
	}

	pass := cfg

	fail := cfg
	fail.Stream.Type = "invalid"

	tests := []struct {
		name         string
		cfg          config.SdkInfo
		log          interfaces.Logger
		expectResult bool
	}{
		{"new sdk valid params", pass, logger, true},
		{"new sdk invalid params", fail, logger, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			annotator, err := factories.NewAnnotator(contracts.AnnotationTPM, tt.cfg)
			if err != nil {
				t.Fatalf(err.Error())
			}
			anno := []interfaces.Annotator{
				annotator,
			}
			instance := NewSdk(anno, tt.cfg, tt.log)
			var wg sync.WaitGroup
			wg.Add(1)
			defer wg.Done()

			result := instance.BootstrapHandler(context.Background(), &wg)
			if result != tt.expectResult {
				t.Errorf("unexpected result: %v", result)
			}
		})
	}
}

func TestNewSdkYaml(t *testing.T) {
	logger := factories.NewLogger(config.LoggingInfo{MinLogLevel: slog.LevelInfo})

	b, err := os.ReadFile("../test/res/config.yaml")
	if err != nil {
		t.Fatalf(err.Error())
	}

	var cfg config.SdkInfo
	err = yaml.Unmarshal(b, &cfg)
	if err != nil {
		t.Fatalf(err.Error())
	}

	pass := cfg

	fail := cfg
	fail.Stream.Type = "invalid"

	tests := []struct {
		name         string
		cfg          config.SdkInfo
		log          interfaces.Logger
		expectResult bool
	}{
		{"new sdk valid params", pass, logger, true},
		{"new sdk invalid params", fail, logger, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			annotator, err := factories.NewAnnotator(contracts.AnnotationTPM, tt.cfg)
			if err != nil {
				t.Fatalf(err.Error())
			}
			anno := []interfaces.Annotator{
				annotator,
			}
			instance := NewSdk(anno, tt.cfg, tt.log)
			var wg sync.WaitGroup
			wg.Add(1)
			defer wg.Done()

			result := instance.BootstrapHandler(context.Background(), &wg)
			if result != tt.expectResult {
				t.Errorf("unexpected result: %v", result)
			}
		})
	}
}
