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
package iota

/*
#cgo CFLAGS: -I./include -DIOTA_STREAMS_CHANNELS_CLIENT
#cgo LDFLAGS: -L./include -liota_streams_c
#include <channels.h>
*/
import "C"
import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/config"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/interfaces"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/message"
	logInterface "github.com/project-alvarium/provider-logging/pkg/interfaces"
	"github.com/project-alvarium/provider-logging/pkg/logging"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

// For randomized seed generation
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type iotaPublisher struct {
	cfg        config.IotaStreamConfig
	logger     logInterface.Logger
	keyload    *C.message_links_t // The Keyload indicates a key needed by the publisher to send messages to the stream
	subscriber *C.subscriber_t    // The publisher is actually subscribed to the stream
	seed       string
}

func NewIotaPublisher(cfg config.IotaStreamConfig, logger logInterface.Logger) (interfaces.StreamProvider, error) {
	bytes := make([]byte, 64)
	rand.Seed(time.Now().UnixNano())
	for i := range bytes {
		bytes[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	seed := string(bytes)
	logger.Write(logging.DebugLevel, fmt.Sprintf("generated streams seed %s", seed))
	return &iotaPublisher{
		cfg:    cfg,
		logger: logger,
		seed:   seed,
	}, nil
}

func (p *iotaPublisher) Connect() error {
	// Generate Transport client
	transport := C.transport_client_new_from_url(C.CString(p.cfg.TangleNode.Uri()))
	p.logger.Write(logging.DebugLevel, fmt.Sprintf("transport established %s", p.cfg.TangleNode.Uri()))

	// Generate Subscriber instance
	cErr := C.sub_new(&p.subscriber, C.CString(p.seed), transport)
	p.logger.Write(logging.DebugLevel, fmt.Sprintf(get_error(cErr)))
	p.logger.Write(logging.DebugLevel, fmt.Sprintf("subscriber established seed=%s", p.seed))

	// Process announcement message
	rawId, err := p.getAnnouncementId(p.cfg.Provider.Uri())
	if err != nil {
		return err
	}

	address := C.address_from_string(C.CString(rawId))
	cErr = C.sub_receive_announce(p.subscriber, address)
	p.logger.Write(logging.DebugLevel, fmt.Sprintf(get_error(cErr)))

	// Fetch sub link and pk for subscription
	var subLink *C.address_t
	var subPk *C.public_key_t

	cErr = C.sub_send_subscribe(&subLink, p.subscriber, address)
	p.logger.Write(logging.DebugLevel, fmt.Sprintf(get_error(cErr)))
	cErr = C.sub_get_public_key(&subPk, p.subscriber)
	p.logger.Write(logging.DebugLevel, fmt.Sprintf(get_error(cErr)))

	subIdStr := C.get_address_id_str(subLink)
	subPkStr := C.public_key_to_string(subPk)

	p.logger.Write(logging.DebugLevel, fmt.Sprintf("send subscription request %s", C.GoString(subIdStr)))
	r := subscriptionRequest{
		MsgId: C.GoString(subIdStr),
		Pk:    C.GoString(subPkStr),
	}
	body, _ := json.Marshal(&r)
	sendSubscriptionIdToAuthor(p.cfg.Provider.Uri(), body)
	p.logger.Write(logging.DebugLevel, "subscription request sent")

	// Obtain key for publishing messages
	p.keyload = p.awaitKeyLoad()

	// Free generated c strings from mem
	C.drop_str(subIdStr)
	C.drop_str(subPkStr)

	return nil
}

func (p *iotaPublisher) Publish(msg message.PublishWrapper) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	messageBytes := C.CBytes(b)
	messageLen := len(string(b))

	p.logger.Write(logging.DebugLevel, fmt.Sprintf("attempting to publish %s", string(b)))
	var msgLinks C.message_links_t
	cErr := C.sub_send_signed_packet(
		&msgLinks,
		p.subscriber,
		*p.keyload,
		nil, 0,
		(*C.uchar)(messageBytes), C.size_t(messageLen))
	p.logger.Write(logging.DebugLevel, fmt.Sprintf(get_error(cErr)))

	var addrLink *C.address_t
	addrLink = msgLinks.msg_link

	inst := C.get_address_inst_str(addrLink)
	id := C.get_address_id_str(addrLink)
	p.logger.Write(logging.DebugLevel, fmt.Sprintf("publish complete %s:%s", C.GoString(inst), C.GoString(id)))

	C.drop_str(inst)
	C.drop_str(id)
	C.drop_links(msgLinks)
	return nil
}

func (p *iotaPublisher) Close() error {
	C.sub_drop(p.subscriber)
	return nil
}

func (p *iotaPublisher) awaitKeyLoad() *C.message_links_t {
	var keyload *C.message_links_t
	var msgIds *C.next_msg_ids_t
	for { // TODO: This should timeout after a configurable period
		// Gen next message ids to look for existing messages
		cErr := C.sub_gen_next_msg_ids(&msgIds, p.subscriber)
		p.logger.Write(logging.DebugLevel, fmt.Sprintf(get_error(cErr)))

		// Search for keyload message from these ids and try to process it
		cErr = C.sub_receive_keyload_from_ids(keyload, p.subscriber, msgIds)
		p.logger.Write(logging.DebugLevel, fmt.Sprintf(get_error(cErr)))

		// Free memory for c msgids object
		C.drop_next_msg_ids(msgIds)

		if cErr != C.ERR_OK {
			p.logger.Write(logging.DebugLevel, "Keyload not found yet... Checking again...")
			// Loop until keyload is found and processed
			time.Sleep(1000 * time.Millisecond)
		} else {
			p.logger.Write(logging.DebugLevel, "obtained keyload successfully")
			break
		}
	}
	return keyload
}

func (p *iotaPublisher) getAnnouncementId(url string) (string, error) {
	type announcementResponse struct {
		AnnouncementId string `json:"announcement_id"`
	}

	p.logger.Write(logging.DebugLevel, fmt.Sprintf("GET %s/get_announcement_id", url))
	resp, err := http.Get(url + "/get_announcement_id")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	p.logger.Write(logging.DebugLevel, fmt.Sprintf("announcement response - %s", string(bodyBytes)))
	var annResp announcementResponse
	if err := json.Unmarshal(bodyBytes, &annResp); err != nil {
		return "", err
	}
	return annResp.AnnouncementId, nil
}

func sendSubscriptionIdToAuthor(url string, body []byte) error {
	client := http.Client{}
	data := bytes.NewReader(body)
	req, err := http.NewRequest("POST", url+"/subscribe", data)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

type subscriptionRequest struct {
	MsgId string `json:"msgid"`
	Pk    string `json:"pk"`
}

func get_error(err C.err_t) string {
	var e = "Unknown Error"
	switch err {
	case C.ERR_OK:
		e = "Operation completed successfully"
	case C.ERR_OPERATION_FAILED:
		e = "The operation failed to complete"
	case C.ERR_NULL_ARGUMENT:
		e = "The function was passed a null argument"
	case C.ERR_BAD_ARGUMENT:
		e = "The function was passed a bad argument"
	}
	return e
}
