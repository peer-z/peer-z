/*
 * Copyright 2020 PeerGum
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package layer

import (
	"errors"
)

const (
	M_PING = iota
	M_ROUTING
	M_REROUTING
)

const (
	MSG_ACK_REQUEST = uint32(1) << iota
	MSG_ACK
	MSG_READ_REQUEST
	MSG_ROUTING_REQUEST
	MSG_TTL_EXPIRED
)

const (
	TTL_EXPIRED = "Message TTL expired"
)

const (
	MSG_FMT_VERSION = 0x100
)

type Content []byte

type peerMessage struct {
	id          uint64
	version     uint16
	source      PeerInfo
	destination Address
	msgType     int
	ref         uint64
	content     Content
	size        uint64
	flags       uint32
	ttl         int
	counter     int
	signature   []byte
}

//type MessageMarshaler interface {
//	MarshalMessage(writer io.Writer, message peerMessage) error
//}
//
//type MessageUnmarshaler interface {
//	UnmarshalMessage(reader io.Reader) (peerMessage, error)
//}

func (message peerMessage) send() error {
	nextHop, err := message.getNextHop()
	if err == nil {
		messagesOut <- peerData{
			peer:    nextHop.Info(),
			message: message,
		}
		return nil
	}
	return err
}

func (message *peerMessage) route() error {
	message.ttl--
	message.counter++
	if message.ttl > 0 {
		return message.send()
	}
	return errors.New(TTL_EXPIRED)
}

func (message *peerMessage) getNextHop() (Peer, error) {
	var distance Distance = DistMax
	var nextHop Peer
	for _, peer := range Me.peers {
		d := message.destination.Distance(peer.address)
		if d < distance && (peer.address != message.destination || message.counter >= MIN_HOPS) {
			distance = d
			nextHop = peer
		}
	}
	return nextHop, nil
}

func (message *peerMessage) ack(destination Address) (*peerMessage, error) {
	return message, nil
}

func (message *peerMessage) nak(destination Address) (*peerMessage, error) {
	return message, nil
}

func (message *peerMessage) encrypt(destination Address, key Key) (*peerMessage, error) {
	var err error
	message.signature, err = message.content.sign()
	if err != nil {
		return nil, err
	}
	message.content, err = message.content.encrypt(key, []byte("MSG"))
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (message peerMessage) setMsgType(command int) peerMessage {
	message.msgType = command
	return message
}

func (message peerMessage) setRef(messageId uint64) peerMessage {
	message.ref = messageId
	return message
}

func (message peerMessage) setStateInfo(reason string) peerMessage {
	message.setStateInfo(reason)
	return message
}

func (message peerMessage) setContent(content []byte) peerMessage {
	message.content = content
	return message
}
