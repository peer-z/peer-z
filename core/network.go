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
	"bytes"
	"encoding/gob"
	"errors"
	upnp "gitlab.com/NebulousLabs/go-upnp"
	"io"
	"log"
	"net"
)

const (
	MAGIC_NUMBER = 0x0368
	RELAY_IP = "127.0.0.1"
	RELAY_PORT = 33200
)

type netData struct {
	addr net.Addr
	data []byte
}

type peerData struct {
	peer    PeerInfo
	message peerMessage
}

var externalIP string
var igd *upnp.IGD

var (
	messagesIn  = make(chan netData, MAX_MESSAGES_IN)
	messagesOut = make(chan peerData, MAX_MESSAGES_OUT)
)

func OpenPort(port uint16, portName string) {
	// forward a port
	err := igd.Forward(port, portName)
	if err != nil {
		log.Fatal(err)
	}
}

// Listen is handling connections from other peers
func Listen() {
	//we're listening on all interfaces on chosen port
	connString := ":" + Me.PortString()
	udpAddress, err := net.ResolveUDPAddr("udp", connString)
	if err != nil {
		Logln("Incorrect connection string", err)
		panic(err)
	}
	conn, err := net.ListenUDP("udp", udpAddress)
	conn.SetReadBuffer(8192)
	if err != nil {
		Logln("Cannot start listening", err)
		panic(err)
	}
	Logln("Starting listening on port", Me.PortString())
	for i := 0; i < MAX_MESSAGE_HANDLERS; i++ {
		go messageHandler()
	}
	for i := 0; i < MAX_MESSAGE_SENDERS; i++ {
		go messageSender()
	}
	Me.SetState(STATE_READY)
	RequestPeering(&Peer{
		address: Address{1, 1, 1, 1},
		ip:      RELAY_IP,
		port:    RELAY_PORT,
	})
	for {
		buffer := make([]byte, 4096)
		count, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			Logln("Connection error", err)
			break
		}
		if count > 0 {
			messagesIn <- netData{
				addr: addr,
				data: buffer[0:count],
			}
		}
		//time.Sleep(1 * time.Second)
	}
	conn.Close()
}

func (peer Peer) MarshalPeerData(writer io.Writer, message peerMessage) error {
	encoder := gob.NewEncoder(writer)
	if err := encoder.Encode(MAGIC_NUMBER); err != nil {
		return err
	}
	peerData := &peerData{
		peer:    peer.Info(),
		message: message,
	}

	return encoder.Encode(peerData)
}

func UnmarshalPeerData(reader io.Reader) (*peerMessage, *PeerInfo, error) {
	decoder := gob.NewDecoder(reader)
	var magic int
	var peerData peerData
	if err := decoder.Decode(&magic); err != nil {
		return nil, nil, err
	}
	if magic != MAGIC_NUMBER {
		return nil, nil, errors.New("Wrong data format.")
	}
	if err := decoder.Decode(&peerData); err != nil {
		return nil, nil, err
	}
	return &peerData.message, &peerData.peer, nil
}

func messageHandler() {
	for {
		select {
		case netMessage := <-messagesIn:
			buffer := bytes.NewBuffer(netMessage.data)
			var message *peerMessage
			var peer *PeerInfo
			var err error
			if message, peer, err = UnmarshalPeerData(buffer); err != nil {
				Logln(err)
				continue
			}
			Logln("Receiving from", message.source, "delivered by", *peer)
			if message.destination != Me.address {
				message.route()
				break
			}
			if message.flags&MSG_ACK_REQUEST > 0 {
				reply := message
				reply.flags &= ^MSG_ACK_REQUEST
				reply.flags |= MSG_ACK
				reply.content = []byte{}
				reply.size = 0
				reply.source = Me.Info()
				reply.destination = message.source.Address
				reply.ttl = MAX_TTL
				reply.counter = 0
				reply.send()
			}
		}
	}
}

//func (ifaces Interfaces) String() (result string) {
//	for _, iface := range ifaces {
//		addrs, _ := iface.GetAddrs()
//		result += fmt.Sprintf("%v\n", addrs)
//	}
//	return
//}
//
//func GetConnectionString(address *Address) string {
//	return net.JoinHostPort(string(address.GetExternal()), address.GetPortString())
//}

func messageSender() {
	for {
		select {
		case peerData := <-messagesOut:
			peer := peerData.peer
			connString := peer.getConnectionString()
			udpAddr, err := net.ResolveUDPAddr("udp", connString)
			if err != nil {
				Logln("Can't build address for", peer)
				continue
			}
			var conn *net.UDPConn
			conn, err = net.DialUDP("udp", nil, udpAddr)
			if err != nil {
				Logln("Can't connect to peer", peer)
				continue
			}
			var count int
			buf := make([]byte, 4096)

			buffer := bytes.NewBuffer(buf)

			if err := Me.MarshalPeerData(buffer, peerData.message); err != nil {
				Logln("Can't marshall message")
				_ = conn.Close()
				continue
			}
			count, err = conn.Write(buf)
			if err != nil {
				Logln("Can't write to peer", peer)
				_ = conn.Close()
				continue
			}
			Logln("Send", count, "bytes...")
		}
	}
}
