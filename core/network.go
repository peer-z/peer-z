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
	"fmt"
	upnp "gitlab.com/NebulousLabs/go-upnp"
	"io"
	"log"
	"net"
	"strings"
)

const (
	MAGIC_NUMBER = 0x0368
	RELAY_IP     = "delphi.local"
	RELAY_PORT   = 33200
)

type netData struct {
	addr net.Addr
	data []byte
}

type peerData struct {
	Peer    PeerInfo
	Message peerMessage
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
	if Me.address != BroadcastAddress {
		RequestPeering(&Peer{
			address: BroadcastAddress,
			ip:      RELAY_IP,
			port:    RELAY_PORT,
		})
	}
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

	peerData := peerData{
		Peer:    peer.Info(),
		Message: message,
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
	return &peerData.Message, &peerData.Peer, nil
}

func messageHandler() {
	for {
		select {
		case netMessage := <-messagesIn:
			fmt.Println("Message from",netMessage.addr)
			buffer := bytes.NewBuffer(netMessage.data)
			var message *peerMessage
			var peer *PeerInfo
			var err error

			if message, peer, err = UnmarshalPeerData(buffer); err != nil {
				Logln("error: ", err)
				continue
			}
			Logln("Receiving msg type", message.MsgType, "from", message.Source, "delivered by", *peer)

			if message.Destination.Address != Me.address {
				Logln("Routing message...")
				message.route()
				continue
			}

			realDestinationIP := strings.Split(netMessage.addr.String(),":")[0]

			if message.Flags&MSG_ACK_REQUEST > 0 {
				reply := message
				reply.Flags &= ^MSG_ACK_REQUEST
				reply.Flags |= MSG_ACK
				reply.Content = []byte{}
				reply.Size = 0
				reply.Source = Me.Info()
				reply.Destination = message.Source
				reply.Destination.IP = realDestinationIP
				reply.Ttl = MAX_TTL
				reply.Counter = 0
				reply.send()
				Logln("Acked message...")
			}
			if err = message.handle(); err != nil {
				Logln("Can't handle message:", err)
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
			peer := peerData.Peer
			peerData.Message.Id = messageIdCounter

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
			//buf := make([]byte, 4096)
			//buffer := bytes.NewBuffer(buf)

			buf := new(bytes.Buffer)

			if err := Me.MarshalPeerData(buf, peerData.Message); err != nil {
				Logln("Can't marshall message: ", err)
				_ = conn.Close()
				continue
			}

			count, err = conn.Write(buf.Bytes())
			if err != nil {
				Logln("Can't write to peer", peer)
				_ = conn.Close()
				continue
			}
			Logln("Sent", count, "bytes...")
		}
	}
}

func dump(buffer []byte, count int) (output string) {
	output = "\n"
	for pos := 0; pos < count; pos += 16 {
		data := ""
		var i int
		for i = 0; i < 16; i++ {
			if pos+i >= count {
				break
			}
			output += fmt.Sprintf("%02x ", buffer[pos+i])
			fmt.Sprintf(data, "%s%c", buffer[i])
		}
		for j := i; j < 16; j++ {
			fmt.Sprintf(data, "   ")
		}
		output += fmt.Sprintf("%s\n", data)
	}
	return output
}
