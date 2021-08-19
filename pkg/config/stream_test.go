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
package config

import (
	"encoding/json"
	"github.com/DyrellC/alvarium-sdk-go/pkg/contracts"
	"github.com/DyrellC/alvarium-sdk-go/test"
	"testing"
)

func TestStreamInfoUnmarshal(t *testing.T) {
	srv := ServiceInfo{
		Host:     "localhost",
		Protocol: "http",
		Port:     8080,
	}

	tangle := ServiceInfo{
		Host:     "localhost",
		Protocol: "http",
		Port:     14265,
	}

	mqtt := ServiceInfo{
		Host:     "localhost",
		Protocol: "tcp",
		Port:     1883,
	}

	stream := IotaStreamConfig{
		Provider:   srv,
		TangleNode: tangle,
		Encoding:   "utf-8",
	}

	streamMqtt := MqttConfig{
		ClientId: "sdk-test",
		Provider: mqtt,
	}

	pass := StreamInfo{
		Type:   contracts.IotaStream,
		Config: stream,
	}

	pass2 := StreamInfo{
		Type:   contracts.MockStream,
		Config: stream,
	}

	pass3 := StreamInfo{
		Type:   contracts.MqttStream,
		Config: streamMqtt,
	}

	fail := StreamInfo{
		Type:   "invalid",
		Config: stream,
	}

	fail2 := StreamInfo{
		Type:   contracts.PravegaStream,
		Config: stream,
	}

	a, _ := json.Marshal(&pass)
	b, _ := json.Marshal(&pass2)
	c, _ := json.Marshal(&pass3)
	d, _ := json.Marshal(&fail)
	e, _ := json.Marshal(&fail2)

	tests := []struct {
		name        string
		bytes       []byte
		expectError bool
	}{
		{"valid StreamInfo type #1", a, false},
		{"valid StreamInfo type #2", b, false},
		{"valid StreamInfo type #3", c, false},
		{"invalid StreamInfo type", d, true},
		{"unhandled StreamInfo type", e, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := StreamInfo{}
			err := s.UnmarshalJSON(tt.bytes)
			test.CheckError(err, tt.expectError, tt.name, t)
			if err == nil { // successful unmarshal
				if s.Type == contracts.IotaStream {
					cfg := s.Config.(IotaStreamConfig)
					if cfg.Provider.Uri() != "http://localhost:8080" {
						t.Errorf("unexpected provider Uri value %s", cfg.Provider.Uri())
					}
					if cfg.TangleNode.Uri() != "http://localhost:14265" {
						t.Errorf("unexpected tangle Uri value %s", cfg.TangleNode.Uri())
					}
				} else if s.Type == contracts.MqttStream {
					cfg := s.Config.(MqttConfig)
					if cfg.Provider.Uri() != "tcp://localhost:1883" {
						t.Errorf("unexpected provider Uri value %s", cfg.Provider.Uri())
					}
				}
			}
		})
	}
}
