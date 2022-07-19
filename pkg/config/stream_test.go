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
package config

import (
	"encoding/json"
	"testing"

	"github.com/project-alvarium/alvarium-sdk-go/pkg/contracts"
	"github.com/project-alvarium/alvarium-sdk-go/test"
)

func TestStreamInfoUnmarshal(t *testing.T) {
	srv := ServiceInfo{
		Host:     "localhost",
		Protocol: "http",
		Port:     8080,
	}

	mqtt := ServiceInfo{
		Host:     "localhost",
		Protocol: "tcp",
		Port:     1883,
	}

	streamMock := MockStreamConfig{
		Provider: srv,
	}

	streamMqtt := MqttConfig{
		ClientId: "sdk-test",
		Provider: mqtt,
	}

	streamHedera := HederaConfig{
		NetType:        contracts.Testnet,
		AccountId:      "testID",
		PrivateKeyPath: "testKey",
		Topics:         []string{"topic1", "topic2"},
	}

	pass := StreamInfo{
		Type:   contracts.MockStream,
		Config: streamMock,
	}

	pass2 := StreamInfo{
		Type:   contracts.MqttStream,
		Config: streamMqtt,
	}

	pass3 := StreamInfo{
		Type:   contracts.ConsoleStream,
		Config: nil,
	}

	pass4 := StreamInfo{
		Type:   contracts.HederaStream,
		Config: streamHedera,
	}

	fail := StreamInfo{
		Type:   "invalid",
		Config: streamMock,
	}

	fail2 := StreamInfo{
		Type:   contracts.PravegaStream,
		Config: streamMock,
	}

	a, _ := json.Marshal(&pass)
	b, _ := json.Marshal(&pass2)
	c, _ := json.Marshal(&pass3)
	d, _ := json.Marshal(&pass4)
	e, _ := json.Marshal(&fail)
	f, _ := json.Marshal(&fail2)

	tests := []struct {
		name        string
		bytes       []byte
		expectError bool
	}{
		{"valid StreamInfo type #1", a, false},
		{"valid StreamInfo type #2", b, false},
		{"valid StreamInfo type #3", c, false},
		{"valid StreamInfo type #4", d, false},
		{"invalid StreamInfo type", e, true},
		{"unhandled StreamInfo type", f, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := StreamInfo{}
			err := s.UnmarshalJSON(tt.bytes)
			test.CheckError(err, tt.expectError, tt.name, t)
			if err == nil { // successful unmarshal
				if s.Type == contracts.MqttStream {
					cfg := s.Config.(MqttConfig)
					if cfg.Provider.Uri() != "tcp://localhost:1883" {
						t.Errorf("unexpected provider Uri value %s", cfg.Provider.Uri())
					}
				} else if s.Type == contracts.HederaStream {
					cfg := s.Config.(HederaConfig)
					if cfg.AccountId != "testID" {
						t.Errorf("unexpected account ID value %s", cfg.AccountId)
					}
				} else if s.Type == contracts.MockStream {
					cfg := s.Config.(MockStreamConfig)
					if cfg.Provider.Uri() != "http://localhost:8080" {
						t.Errorf("unexpected provider Uri value %s", cfg.Provider.Uri())
					}
				}
			}
		})
	}
}
