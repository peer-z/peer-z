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

package core

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

var messageIdCounter uint64

type peerMessage struct {
	Id          uint64
	Version     uint16
	Source      PeerInfo
	Destination PeerInfo
	MsgType     int
	Ref         uint64
	Content     Content
	Size        uint64
	Flags       uint32
	Ttl         int
	Counter     int
	Signature   []byte
}

//type MessageMarshaler interface {
//	MarshalMessage(writer io.Writer, message peerMessage) error
//}
//
//type MessageUnmarshaler interface {
//	UnmarshalMessage(reader io.Reader) (peerMessage, error)
//}

func (message peerMessage) send() error {
	var nextHopInfo PeerInfo
	if len(message.Destination.IP)>0 {
		nextHopInfo =message.Destination
	} else {
		nextHop, err := message.getNextHop()
		if err == nil {
			nextHopInfo = nextHop.Info()
		} else {
			return err
		}
	}
	messagesOut <- peerData{
		Peer:    nextHopInfo,
		Message: message,
	}
	return nil
}

func (message *peerMessage) route() error {
	message.Ttl--
	message.Counter++
	if message.Ttl > 0 {
		return message.send()
	}
	return errors.New(TTL_EXPIRED)
}

func (message *peerMessage) getNextHop() (Peer, error) {
	var distance Distance = DistMax
	var nextHop Peer
	if message.Destination.Address.isBroadcast() {
		return *Me.peers.searchByAddress(message.Destination.Address), nil
	}
	for _, peer := range Me.peers {
		if peer.address == message.Destination.Address && message.Counter == -1 {
			// direct message to a peer
			return peer, nil
		}
		d := message.Destination.Address.Distance(peer.address)
		if d < distance && (peer.address != message.Destination.Address || message.Counter >= MIN_HOPS) {
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
	message.Signature, err = message.Content.sign()
	if err != nil {
		return nil, err
	}
	message.Content, err = message.Content.encrypt(key, []byte("MSG"))
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (message peerMessage) setMsgType(command int) peerMessage {
	message.MsgType = command
	return message
}

func (message peerMessage) setRef(messageId uint64) peerMessage {
	message.Ref = messageId
	return message
}

func (message peerMessage) setStateInfo(reason string) peerMessage {
	message.setStateInfo(reason)
	return message
}

func (message peerMessage) setContent(content []byte) peerMessage {
	message.Content = content
	return message
}
