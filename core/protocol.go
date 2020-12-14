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
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/alecthomas/binary"
	"reflect"
)

const (
	CMD_PR     = iota // peering request
	CMD_PR_ACK        // peering request acknowledgement
	CMD_PR_NAK        // peering requestrefusal
	CMD_PL            // peer list information
	CMD_PL_ACK        // peer listresponse
	CMD_RP            // re-routing proposal
	CMD_RP_ACK        // re-routing proposalacknowledgement
	CMD_RP_NAK        // re-routing proposalrefusal
	CMD_SR            // standby request
	CMD_SR_ACK        // standby request acknowledgement
	CMD_PR_OK         // confirm peering established
	CMD_RR            // Re-routing request
	CMD_RR_ACK        // Re-routing request acknowledgement
	CMD_RR_NAK        // Re-routing request refusal
	CMD_DR            // Disconnect request
	CMD_DR_ACK        // Disconnect request acknowledgement
)

type RerouteList struct {
	Gimme *PeerList
	Get   *PeerList
}

// --

//type Contact struct {
//	ip base.IPAddress
//	Port     int
//	Address  AddressValue
//}
//
// Initializer
func init() {
	// if any interface is used to be gob encoded/decoded, register here
	// e.g. gob.Register(struct)
}

// -- peering requests

// RequestPeering Starts a peering request with peer
func RequestPeering(peer *Peer) error {
	message := &peerMessage{
		msgType: CMD_PR,
		source:  Me.Info(),
	}
	err := peer.send(message)
	if err != nil {
		Logln("Failed Peering Request", err)
		return err
	}
	peer.SetState(STATE_PR_SENT)
	return nil
}

// accept the peering request
func acceptPeering(peer *Peer, messageId uint64) error {
	message := &peerMessage{
		msgType: CMD_PR_ACK,
		source:  Me.Info(),
	}
	err := peer.send(message)
	if err != nil {
		Logln("Peering Acceptation Failed", err)
		return err
	}
	peer.SetState(STATE_PR_WAIT).
		SetStateInfo("waiting for re-routing")
	Logln("Peering acceptance sent.")
	return nil
}

// refuse the peering request with reason
func refusePeering(peer *Peer, messageId uint64, reason string) error {
	message := &peerMessage{
		msgType: CMD_PR_NAK,
		ref:     messageId,
		source:  Me.Info(),
		content: []byte(reason),
	}
	err := peer.send(message)
	if err != nil {
		Logln("Peering Refusal Failed", err)
		return err
	}
	peer.SetState(STATE_PR_REFUSED).
		SetStateInfo(reason)
	Logln("Peering refusal sent.")
	return nil
}

// Sends the peerlist. if ack is true, this is the reply
func sendPeerList(peer *Peer, messageId uint64, peers Peers, reply bool) error {
	command := CMD_PL
	if reply {
		command = CMD_PL_ACK
	}
	content, err := binary.Marshal(peers)
	if err != nil {
		Logln("Error sending PeerList", err)
		return err
	}
	message := &peerMessage{
		msgType: command,
		ref:     messageId,
		source:  Me.Info(),
		content: content,
	}
	err = peer.send(message)
	if err != nil {
		Logln("Error sending PeerList", err)
		return err
	}
	if reply {
		peer.SetState(STATE_PL_WAIT)
		Logln("PeerList replied.", peers)
	} else {
		peer.SetState(STATE_PL_SENT)
		Logln("PeerList sent.", peers)
	}
	return nil
}

// sends the rerouting proposal
func sendRerouteProposal(peer *Peer, messageId uint64, rerouteList *RerouteList) error {
	content, err := binary.Marshal(rerouteList)
	if err != nil {
		Logln("Error sending RerouteList", err)
		return err
	}
	message := &peerMessage{
		msgType: CMD_RP,
		ref:     messageId,
		source:  Me.Info(),
		content: content,
	}
	err = peer.send(message)
	if err != nil {
		Logln("Error sending RerouteList", err)
		return err
	}
	peer.SetState(STATE_RP_SENT)
	Logln("RerouteList sent.")
	return nil
}

// accept the rerouting proposal
func acceptRerouting(peer *Peer, messageId uint64) error {
	message := &peerMessage{
		msgType: CMD_RP_ACK,
		source:  Me.Info(),
	}
	err := peer.send(message)
	if err != nil {
		Logln("Rerouting Acceptation Failed", err)
		return err
	}
	peer.SetState(STATE_STANDBY).
		SetStateInfo("waiting for re-routing")
	Logln("Rerouting acceptance sent.")
	return nil
}

// accept the rerouting proposal
func refuseRerouting(peer *Peer, messageId uint64, reason string) error {
	message := &peerMessage{
		msgType: CMD_RP_NAK,
		source:  Me.Info(),
		content: []byte(reason),
	}
	err := peer.send(message)
	if err != nil {
		Logln("Rerouting refusal failed", err)
		return err
	}
	peer.SetState(STATE_STANDBY).
		SetStateInfo(reason)
	Logln("Rerouting refusal sent.")
	return nil
}

// protocol finished, confirm connection
func confirmPeering(peer *Peer, messageId uint64) error {
	message := &peerMessage{
		msgType: CMD_PR_OK,
		source:  Me.Info(),
	}
	err := peer.send(message)
	if err != nil {
		Logln("Peering confirmation failed", err)
		return err
	}
	peer.
		SetState(STATE_CONNECTED).
		SetStateInfo("All Good")
	Logln("Connection confirmed.")
	Me.addPeer(*peer)
	Logln(Me.GetList())
	return nil
}

//
// -- set the protocol message
// 	these calls can be chained
//

func (message peerMessage) setPeerList(peerList PeerList) (peerMessage, error) {
	content, err := binary.Marshal(peerList)
	if err != nil {
		return message, err
	}
	message.setContent(content)
	return message, nil
}

func (message peerMessage) setRerouteList(rerouteList RerouteList) (peerMessage, error) {
	content, err := binary.Marshal(rerouteList)
	if err != nil {
		return message, err
	}
	message.setContent(content)
	return message, nil
}

//
// -- send and receive protocol messages
//

//func (message Message) send(peer *Peer) error {
//	message.MessageId = uint64(rand.Int63())
//	var buffer bytes.Buffer
//	encoder := gob.NewEncoder(&buffer)
//	message.encode(encoder)
//	return peer.send(buffer)
//}
//
//func (message *Message) receive(peer *Peer, data []byte) error {
//	buffer := bytes.NewBuffer(data)
//	decoder := gob.NewDecoder(buffer)
//	message.decode(decoder)
//	peer.lastRcvd = time.Now().UnixNano()
//	return message.handle()
//}

//
// -- Encode/Decode messages to/from gobs
//

// messageEncode encode messages to gobs
func (message peerMessage) encode(encoder *gob.Encoder) {
	err := encoder.Encode(&message)
	if err != nil {
		Logln("Can't encode", err)
	}
}

// messageDecode decode messages from gobs
func (message *peerMessage) decode(decoder *gob.Decoder) {
	err := decoder.Decode(&message)
	if err != nil {
		Logln("Can't encode", err)
	}
}

//
// -- handle protocol
//

func (message *peerMessage) handle() error {
	Logf("MessageId: %d\nCommand: %d\n", message.id, message.msgType)

	protocolFunctions := map[int]func(*peerMessage) error{
		CMD_PR:     handlePR,
		CMD_PR_ACK: sendPL,
		CMD_PR_NAK: cancelPR,
		CMD_PL:     handlePL,
		CMD_PL_ACK: sendRP,
		CMD_RP:     handleRP,
		CMD_RP_ACK: requestRRs,
		CMD_RP_NAK: concludePR,
		CMD_PR_OK:  finishPR,
	}

	if function, ok := protocolFunctions[message.msgType]; ok {
		return function(message)
	}
	Logln("Protocol command not currently supported.")
	return nil
}

// received a PR, handle it
func handlePR(message *peerMessage) error {
	source := message.source
	Logf("Peering Request from: %s\n", message.source)

	messageId := message.id
	peer := Me.peers.searchByInfo(source)
	if peer == nil {
		peer = NewPeer(source.Name, source.Address, source.IP, source.Port)
		if _, err := Me.addPeer(*peer); err != nil {
			refusePeering(peer, messageId, err.Error())
			return err
		}
	} else {
		Logln("Peer found:", peer)
	}
	acceptPeering(peer, messageId)
	return nil
}

// received a PR cancellation
func cancelPR(message *peerMessage) error {
	Logf("Peering Request refused by: %v\nReason:%s\n", message.source, string(message.content))

	peer := Me.peers.searchByInfo(message.source)
	if peer == nil {
		// contact not found... discard
		Logln("Can't find contact anymore... Weird!")
		return errors.New("Peer disappeared")
	}
	peer.
		SetState(STATE_PR_REFUSED).
		SetStateInfo(string(message.content))
	return nil
}

// received a PR confirmation, send PL
func sendPL(message *peerMessage) error {
	// search by IP since we don't know the peer address
	// at this point
	source := message.source
	Logf("Peering Request accepted by: %s\n", message.source)

	peer := Me.peers.searchByInfo(source)
	if peer == nil {
		// contact not found... discard
		Logln("Can't find contact anymore... Weird!")
		return errors.New("Peer disappeared")
	}
	Logln(source)
	peer.SetAddress(source.Address)
	Logln("Updated peer:", peer)
	peers := Me.GetList()
	sendPeerList(peer, message.id, peers, false)
	return nil
}

// received a PL, send ours
func handlePL(message *peerMessage) error {
	Logf("PeerList Received from: %v\n", message.source)

	peer := Me.peers.searchByInfo(message.source)
	if peer == nil {
		// contact not found... discard
		Logln("Can't find contact anymore... Weird!")
		return errors.New("Peer disappeared")
	}
	var peers Peers
	err := binary.Unmarshal(message.content, peers)
	if err!=nil {
		Logln("Failed decoding peer list")
		return errors.New("Failed decoding peer list")
	}
	Logln("peerList received:", peers)
	// save peer's list
	peer.peers = peers
	// get and send mine
	ourList := Me.GetList()
	sendPeerList(peer, message.id, ourList, true)
	return nil
}

// received a PL reply, calculcate and send re-routing proposal
func sendRP(message *peerMessage) error {
	Logf("PeerList Replied by: %s\n", message.source)
	peer := Me.peers.searchByInfo(message.source)
	if peer == nil {
		// contact not found... discard
		Logln("Can't find contact anymore... Weird!")
		return errors.New("Peer disappeared")
	}
	var peers Peers
	err := binary.Unmarshal(message.content,peers)
	if err!=nil {
		Logln("Failed decoding peer list")
		return errors.New("Failed decoding peer list")
	}

	peer.peers = peers
	rerouteList := Me.CheckList(peer)
	Logln("Sending RP:", rerouteList)
	sendRerouteProposal(peer, message.id, rerouteList)
	return nil
}

// receive a rerouting proposal, confirm or refuse
func handleRP(message *peerMessage) error {
	Logf("Rerouting List Replied by: %s\n", message.source)
	peer := Me.peers.searchByInfo(message.source)
	if peer == nil {
		// contact not found... discard
		Logln("Can't find contact anymore... Weird!")
		return errors.New("Peer disappeared")
	}
	var rerouteList RerouteList
	err := binary.Unmarshal(message.content,rerouteList)
	if err != nil {
		Logln("Failed decoding reroute list")
		return errors.New("Failed decoding reroute list")
	}
	if rerouteList.Gimme == nil {
		rerouteList.Gimme = &PeerList{}
	}
	if rerouteList.Get == nil {
		rerouteList.Get = &PeerList{}
	}
	myRerouteList := Me.CheckList(peer)
	for i, gimmePeer := range *rerouteList.Gimme {
		getPeer := (*myRerouteList.Get)[i]
		if reflect.DeepEqual(gimmePeer,getPeer) {
			err := errors.New("Proposals don't match")
			Logln(err, gimmePeer)
			refuseRerouting(peer, message.id, err.Error())
			return err
		}
	}
	for i, getPeer := range *rerouteList.Get {
		gimmePeer := (*myRerouteList.Get)[i]
		if reflect.DeepEqual(getPeer,gimmePeer) {
			err := errors.New("Proposals don't match")
			Logln(err, getPeer)
			refuseRerouting(peer, message.id, err.Error())
			return err
		}
	}
	peer.rerouteList = &rerouteList
	Logln("Proposals match.")
	acceptRerouting(peer, message.id)
	return nil
}

func requestRRs(message *peerMessage) error {
	Logf("Rerouting Protocol accepted by: %s\n", message.source)
	peer := Me.peers.searchByInfo(message.source)
	if peer == nil {
		// contact not found... discard
		Logln("Can't find contact anymore... Weird!")
		return errors.New("Peer disappeared")
	}
	confirmPeering(peer, message.id)
	return nil
}

func concludePR(message *peerMessage) error {
	Logf("Rerouting Protocol refused by: %s\n", message.source)
	peer := Me.peers.searchByInfo(message.source)
	if peer == nil {
		// contact not found... discard
		Logln("Can't find contact anymore... Weird!")
		return errors.New("Peer disappeared")
	}
	confirmPeering(peer, message.id)

	return nil
}

func finishPR(message *peerMessage) error {
	Logf("Peering confirmed by: %s\n", message.source)
	peer := Me.peers.searchByInfo(message.source)
	if peer == nil {
		// contact not found... discard
		Logln("Can't find contact anymore... Weird!")
		return errors.New("Peer disappeared")
	}
	peer.SetState(STATE_CONNECTED)
	Me.addPeer(*peer)
	Logln(Me.GetList())
	return nil
}

// properly display a contact
func (peerInfo PeerInfo) String() string {
	return fmt.Sprintf("%s [%s] (%s:%d)",peerInfo.Name,peerInfo.Address,peerInfo.IP,peerInfo.Port)
}

