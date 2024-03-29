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
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/project-alvarium/alvarium-sdk-go/pkg/config"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/contracts"
	"github.com/project-alvarium/alvarium-sdk-go/test"
)

func TestStreamProviderFactory(t *testing.T) {
	logger := NewLogger(config.LoggingInfo{MinLogLevel: slog.LevelDebug})

	pass := config.StreamInfo{
		Type:   contracts.MockStream,
		Config: config.MockStreamConfig{},
	}

	pass2 := config.StreamInfo{
		Type:   contracts.MqttStream,
		Config: config.MqttConfig{},
	}

	fail := config.StreamInfo{
		Type:   "invalid",
		Config: config.MqttConfig{},
	}

	fail2 := config.StreamInfo{
		Type:   "pravega",
		Config: config.MockStreamConfig{},
	}

	tests := []struct {
		name         string
		providerType config.StreamInfo
		expectError  bool
	}{
		{"valid mock type", pass, false},
		{"valid mqtt type", pass2, false},
		{"invalid random type", fail, true},
		{"unimplemented pravega type", fail2, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewStreamProvider(tt.providerType, logger)
			test.CheckError(err, tt.expectError, tt.name, t)
		})
	}
}

func TestHashProviderFactory(t *testing.T) {
	tests := []struct {
		name         string
		providerType contracts.HashType
		expectError  bool
	}{
		{"valid md5 type", contracts.MD5Hash, false},
		{"valid sha256 type", contracts.SHA256Hash, false},
		{"valid none type", contracts.NoHash, false},
		{"invalid hash type", "invalid", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewHashProvider(tt.providerType)
			test.CheckError(err, tt.expectError, tt.name, t)
		})
	}
}

func TestSignatureProviderFactory(t *testing.T) {
	tests := []struct {
		name         string
		providerType contracts.KeyAlgorithm
		expectError  bool
	}{
		{"valid ed25519 type", contracts.KeyEd25519, false},
		{"invalid hash type", "invalid", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewSignatureProvider(tt.providerType)
			test.CheckError(err, tt.expectError, tt.name, t)
		})
	}
}

func TestAnnotatorFactory(t *testing.T) {
	b, err := os.ReadFile("../../test/res/config.json")
	if err != nil {
		t.Fatalf(err.Error())
	}
	var cfg config.SdkInfo
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		t.Fatalf(err.Error())
	}

	tests := []struct {
		name        string
		cfg         config.SdkInfo
		key         contracts.AnnotationType
		expectError bool
	}{
		{"valid pki type", cfg, contracts.AnnotationPKI, false},
		{"valid httpPki type", cfg, contracts.AnnotationPKIHttp, false},
		{"valid src type", cfg, contracts.AnnotationSource, false},
		{"valid tpm type", cfg, contracts.AnnotationTPM, false},
		{"valid tls type", cfg, contracts.AnnotationTLS, false},
		{"invalid annotator type", cfg, "invalid", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewAnnotator(tt.key, tt.cfg)
			test.CheckError(err, tt.expectError, tt.name, t)
		})
	}
}

func TestRequestHandlerFactory(t *testing.T) {

	type sample struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	a := sample{Key: "keyA", Value: "This is some test data"}
	b, _ := json.Marshal(a)

	req := httptest.NewRequest("POST", "/foo?param=value&foo=bar&baz=batman", bytes.NewReader(b))

	cfg := config.SignatureInfo{}
	pass := cfg
	pass.PrivateKey.Type = contracts.KeyEd25519
	fail := cfg
	fail.PublicKey.Type = "invalid"

	tests := []struct {
		name        string
		cfg         config.SignatureInfo
		expectError bool
	}{
		{"valid ed25519 type", pass, false},
		{"invalid key type", fail, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewRequestHandler(req, tt.cfg)
			test.CheckError(err, tt.expectError, tt.name, t)
		})
	}
}
