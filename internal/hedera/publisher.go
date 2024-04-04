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

package hedera

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/hashgraph/hedera-sdk-go/v2"
	"github.com/project-alvarium/alvarium-sdk-go/internal/mqtt"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/config"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/contracts"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/interfaces"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/message"
)

type HederaPublisher struct {
	cfg             config.HederaConfig
	logger          interfaces.Logger
	hederaClient    *hedera.Client
	broadcastStream interfaces.StreamProvider
}

func NewHederaPublisher(
	cfg config.HederaConfig,
	logger interfaces.Logger,
) (interfaces.StreamProvider, error) {
	client, err := initHederaClient(cfg)
	if err != nil {
		return nil, err
	}

	p := HederaPublisher{
		cfg:          cfg,
		logger:       logger,
		hederaClient: client,
	}
	return &p, nil
}

// hedera client implicitly connects to the hedera net.
// no need for manual initiation. Instead, topics used to
// publish annotations will be broadcasted according to
// configuration
func (p *HederaPublisher) Connect() error {
	if p.cfg.ShouldBroadcastTopic {

		stream, err := initBroadcastStream(p.cfg, p.logger)
		if err != nil {
			return err
		}

		p.broadcastStream = stream
		for _, topic := range p.cfg.Topics {
			msg := message.PublishWrapper{
				Action:      message.ActionBroadcast,
				MessageType: fmt.Sprintf("%T", topic),
				Content:     []byte(topic),
			}
			err := p.broadcastStream.Publish(msg)
			if err != nil {
				return err
			}

		}
	}
	return nil
}

func (p *HederaPublisher) Publish(msg message.PublishWrapper) error {
	b, _ := json.Marshal(msg)

	// publish to all topic IDs
	for _, topic := range p.cfg.Topics {
		p.logger.Write(
			slog.LevelDebug,
			fmt.Sprintf("attempting publish, topic %s %s", topic, string(b)),
		)
		topicId, err := hedera.TopicIDFromString(topic)
		if err != nil {
			return err
		}

		// submit message to consensus service
		_, err = hedera.NewTopicMessageSubmitTransaction().
			SetMessage(b).
			SetTopicID(topicId).
			Execute(p.hederaClient)
		if err != nil {
			return err
		}
	}

	return nil
}

// Parties aware of Hedera topics are notified that the stream
// is closing
func (p *HederaPublisher) Close() error {
	if p.cfg.ShouldBroadcastTopic {
		for _, topic := range p.cfg.Topics {
			msg := message.PublishWrapper{
				Action:      message.ActionEndStream,
				MessageType: fmt.Sprintf("%T", topic),
				Content:     []byte(topic),
			}
			err := p.broadcastStream.Publish(msg)
			if err != nil {
				return err
			}

		}
	}
	return p.hederaClient.Close()
}

// Initialize a Hedera client with configuration driven values
// and default values
func initHederaClient(cfg config.HederaConfig) (*hedera.Client, error) {
	var client *hedera.Client
	switch netType := cfg.NetType; netType {
	case contracts.Mainnet:
		client = hedera.ClientForMainnet()
	case contracts.Testnet:
		client = hedera.ClientForTestnet()
	case contracts.Previewnet:
		client = hedera.ClientForPreviewnet()
	default:
		return nil, errors.New("nettype not valid")
	}

	accountId, err := hedera.AccountIDFromString(cfg.AccountId)
	if err != nil {
		return nil, err
	}

	privateKey, err := readPrivateKey(cfg)
	if err != nil {
		return nil, err
	}

	client.SetOperator(accountId, privateKey)

	err = setDefaultValues(client, cfg)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func readPrivateKey(cfg config.HederaConfig) (hedera.PrivateKey, error) {
	b, err := os.ReadFile(cfg.PrivateKeyPath)
	if err != nil {
		return hedera.PrivateKey{}, err
	}

	// It was reported by multiple parties that a `\n` character is
	// occasionally loaded into the private key byte array, and
	// other times it was not when the file was created using
	// different methods or saved by different editors
	//
	// This part will remove the newline character only if it exists
	privateKeyDER := string(b)
	if privateKeyDER[len(privateKeyDER)-1] == '\n' {
		privateKeyDER = privateKeyDER[:len(privateKeyDER)-1]
	}

	privateKey, err := hedera.PrivateKeyFromStringDer(privateKeyDER)
	if err != nil {
		return hedera.PrivateKey{}, err
	}

	return privateKey, nil
}

func setDefaultValues(client *hedera.Client, cfg config.HederaConfig) error {
	maxTxFeeInHbar := hedera.HbarFrom(cfg.DefaultMaxTxFee, hedera.HbarUnits.Hbar)
	maxQueryPaymentFee := hedera.HbarFrom(cfg.DefaultMaxQueryPayment, hedera.HbarUnits.Hbar)

	err := client.SetDefaultMaxTransactionFee(maxTxFeeInHbar)
	if err != nil {
		return err
	}

	err = client.SetDefaultMaxQueryPayment(maxQueryPaymentFee)
	if err != nil {
		return err
	}

	return nil
}

func initBroadcastStream(
	cfg config.HederaConfig,
	logger interfaces.Logger,
) (interfaces.StreamProvider, error) {
	stream := mqtt.NewMqttPublisher(cfg.BroadcastStream, logger)

	err := stream.Connect()
	if err != nil {
		return nil, err
	}
	return stream, nil
}
