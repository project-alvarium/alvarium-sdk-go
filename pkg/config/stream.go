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
	"fmt"

	"github.com/project-alvarium/alvarium-sdk-go/pkg/contracts"
	"gopkg.in/yaml.v3"
)

// StreamInfo facilitates configuration of a given streaming platform that will receive annotations
type StreamInfo struct {
	Type   contracts.StreamType `json:"type,omitempty" yaml:"type"`
	Config interface{}          `json:"config,omitempty" yaml:"config"`
}

func (s *StreamInfo) UnmarshalJSON(data []byte) (err error) {
	type Alias struct {
		Type contracts.StreamType `json:"type,omitempty"`
	}
	a := Alias{}
	// Error with unmarshaling
	if err = json.Unmarshal(data, &a); err != nil {
		return err
	}

	if !a.Type.Validate() {
		return fmt.Errorf("invalid StreamType value provided %s", a.Type)
	}

	if a.Type == contracts.MqttStream {
		type mqttAlias struct {
			Type   contracts.StreamType `json:"type,omitempty"`
			Config MqttConfig           `json:"config,omitempty"`
		}

		m := mqttAlias{}
		// Error with unmarshaling
		if err = json.Unmarshal(data, &m); err != nil {
			return err
		}
		s.Type = m.Type
		s.Config = m.Config
	} else if a.Type == contracts.MockStream {
		type mockAlias struct {
			Type   contracts.StreamType `json:"type,omitempty"`
			Config MockStreamConfig     `json:"config,omitempty"`
		}
		m := mockAlias{}
		// Error with unmarshaling
		if err = json.Unmarshal(data, &m); err != nil {
			return err
		}
		s.Type = m.Type
		s.Config = m.Config
	} else if a.Type == contracts.ConsoleStream {
		type consoleAlias struct {
			Type   contracts.StreamType `json:"type,omitempty"`
			Config MockStreamConfig     `json:"config,omitempty"`
		}
		c := consoleAlias{}
		// Error with unmarshaling
		if err = json.Unmarshal(data, &c); err != nil {
			return err
		}
		s.Type = c.Type
		s.Config = MockStreamConfig{}
	} else if a.Type == contracts.HederaStream {
		type hederaAlias struct {
			Type   contracts.StreamType `json:"type,omitempty"`
			Config HederaConfig         `json:"config,omitempty"`
		}

		h := hederaAlias{}
		// Error with unmarshaling
		if err = json.Unmarshal(data, &h); err != nil {
			return err
		}
		s.Type = h.Type
		s.Config = h.Config
	} else {
		return fmt.Errorf("unhandled StreamInfo.Type value %s", a.Type)
	}

	return nil
}

func (s *StreamInfo) UnmarshalYAML(data *yaml.Node) (err error) {
	type Alias struct {
		Type contracts.StreamType `yaml:"type"`
	}
	a := Alias{}
	// Error with unmarshaling
	if err = data.Decode(&a); err != nil {
		return err
	}

	if !a.Type.Validate() {
		return fmt.Errorf("invalid StreamType value provided %s", a.Type)
	}

	if a.Type == contracts.MqttStream {
		type mqttAlias struct {
			Type   contracts.StreamType `yaml:"type"`
			Config MqttConfig           `yaml:"config"`
		}

		m := mqttAlias{}
		// Error with unmarshaling
		if err = data.Decode(&m); err != nil {
			return err
		}
		s.Type = m.Type
		s.Config = m.Config
	} else if a.Type == contracts.HederaStream {
		type hederaAlias struct {
			Type   contracts.StreamType `yaml:"type"`
			Config HederaConfig         `yaml:"config"`
		}

		h := hederaAlias{}
		// Error with unmarshaling
		if err = data.Decode(&h); err != nil {
			return err
		}
		s.Type = h.Type
		s.Config = h.Config
	} else if a.Type == contracts.MockStream {
		type mockAlias struct {
			Type   contracts.StreamType `yaml:"type"`
			Config MockStreamConfig     `yaml:"config"`
		}
		m := mockAlias{}
		// Error with unmarshaling
		if err = data.Decode(&m); err != nil {
			return err
		}
		s.Type = m.Type
		s.Config = m.Config
	} else if a.Type == contracts.ConsoleStream {
		type consoleAlias struct {
			Type   contracts.StreamType `json:"type,omitempty"`
			Config MockStreamConfig     `json:"config,omitempty"`
		}
		c := consoleAlias{}
		// Error with unmarshaling
		if err = data.Decode(&c); err != nil {
			return err
		}
		s.Type = c.Type
		s.Config = MockStreamConfig{}
	} else {
		return fmt.Errorf("unhandled StreamInfo.Type value %s", a.Type)
	}

	return nil
}

// MqttConfig exposes properties relevant to connecting to an existing MQTT broker
type MqttConfig struct {
	ClientId  string      `json:"clientId,omitempty" yaml:"clientId"`
	Qos       int         `json:"qos,omitempty" yaml:"qos"`
	User      string      `json:"user,omitempty" yaml:"user"`
	Password  string      `json:"password,omitempty" yaml:"password"`
	Provider  ServiceInfo `json:"provider,omitempty" yaml:"provider"`
	Cleanness bool        `json:"cleanness,omitempty" yaml:"cleanness"`
	Topics    []string    `json:"topics,omitempty" yaml:"topics"`
}

// MockStreamConfig exposes properties to simulate a stream connection for testing.
type MockStreamConfig struct {
	Provider ServiceInfo `json:"provider,omitempty" yaml:"provider"`
}

// HederaConfig provides configuartion required to init a Hedera client and connect to the consensus nodes

type HederaConfig struct {
	NetType                contracts.NetType `json:"netType,omitempty"    yaml:"netType"`
	Consensus              ServiceInfo       `json:"consensus,omitempty" yaml:"consensus"` // Only populated when NetType is "local"
	Mirror                 ServiceInfo       `json:"mirror,omitempty"   yaml:"mirror"`     // Only populated when NetType is "local"
	AccountId              string            `json:"accountId,omitempty"  yaml:"accountId"`
	PrivateKeyPath         string            `json:"privateKeyPath,omitempty" yaml:"privateKeyPath"`
	Topics                 []string          `json:"topics,omitempty"     yaml:"topics"`
	DefaultMaxTxFee        float64           `json:"defaultMaxTxFee,omitempty" yaml:"defaultMaxTxFee"`
	DefaultMaxQueryPayment float64           `json:"defaultMaxQueryPayment,omitempty" yaml:"defaultMaxQueryPayment"`
	ShouldBroadcastTopic   bool              `json:"shouldBroadcastTopic,omitempty" yaml:"shouldBroadcastTopic"`

	// TODO (Ali Amin): Add support for other providers
	BroadcastStream MqttConfig `json:"broadcastStream,omitempty" yaml:"broadcastStream"`
}

// ServiceInfo describes a service endpoint that the deployed service is a client of. Right now, this is implicitly
// an HTTP interaction
type ServiceInfo struct {
	Host     string `json:"host,omitempty" yaml:"host"`
	Port     int    `json:"port,omitempty" yaml:"port"`
	Protocol string `json:"protocol,omitempty" yaml:"protocol"`
}

// Uri constructs a string from the populated elements of the ServiceInfo
func (s ServiceInfo) Uri() string {
	return fmt.Sprintf("%s://%s:%v", s.Protocol, s.Host, s.Port)
}

// Address constructs a string representing the hostname/IP and port of a given endpoint
func (s ServiceInfo) Address() string {
	return fmt.Sprintf("%s:%v", s.Host, s.Port)
}
