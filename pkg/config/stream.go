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

	if a.Type == contracts.IotaStream {
		type iotaAlias struct {
			Type   contracts.StreamType `json:"type,omitempty"`
			Config IotaStreamConfig     `json:"config,omitempty"`
		}

		i := iotaAlias{}
		// Error with unmarshaling
		if err = json.Unmarshal(data, &i); err != nil {
			return err
		}
		s.Type = i.Type
		s.Config = i.Config
	} else if a.Type == contracts.MqttStream {
		type mqttAlias struct {
			Type   contracts.StreamType `json:"type,omitempty"`
			Config MqttConfig           `json:"config,omitempty"`
		}

		m := mqttAlias{}
		//Error with unmarshaling
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
		//Error with unmarshaling
		if err = json.Unmarshal(data, &m); err != nil {
			return err
		}
		s.Type = m.Type
		s.Config = m.Config
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

	if a.Type == contracts.IotaStream {
		type iotaAlias struct {
			Type   contracts.StreamType `yaml:"type"`
			Config IotaStreamConfig     `yaml:"config"`
		}

		i := iotaAlias{}
		// Error with unmarshaling
		if err = data.Decode(&i); err != nil {
			return err
		}
		s.Type = i.Type
		s.Config = i.Config
	} else if a.Type == contracts.MqttStream {
		type mqttAlias struct {
			Type   contracts.StreamType `yaml:"type"`
			Config MqttConfig           `yaml:"config"`
		}

		m := mqttAlias{}
		//Error with unmarshaling
		if err = data.Decode(&m); err != nil {
			return err
		}
		s.Type = m.Type
		s.Config = m.Config
	} else if a.Type == contracts.MockStream {
		type mockAlias struct {
			Type   contracts.StreamType `yaml:"type"`
			Config MockStreamConfig     `yaml:"config"`
		}
		m := mockAlias{}
		//Error with unmarshaling
		if err = data.Decode(&m); err != nil {
			return err
		}
		s.Type = m.Type
		s.Config = m.Config
	} else {
		return fmt.Errorf("unhandled StreamInfo.Type value %s", a.Type)
	}

	return nil
}

// IotaStreamConfig exposes properties relevant to connecting to an existing IOTA Stream and accompanying Tangle node
type IotaStreamConfig struct {
	Provider   ServiceInfo `json:"provider,omitempty" yaml:"provider"` // Provider is the endpoint from which the Streams subscription is obtained
	TangleNode ServiceInfo `json:"tangle,omitempty" yaml:"tangle"`     // TangleNode is the endpoint of the local Hornet instance. Transactions are written here.
	Encoding   string      `json:"encoding,omitempty" yaml:"encoding"` // Encoding specifies the encoding of transaction messages
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
