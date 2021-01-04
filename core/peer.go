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
    name "github.com/dustinkirkland/golang-petname"
    "gorm.io/gorm"
    "net"
    "sort"
    "strconv"
    "strings"
    "time"
)

// Peer state
const (
    STATE_OFF         = iota // peer unreachable
    STATE_INITIALIZED        // peer started
    STATE_READY              // peer ready to receive
    STATE_CONNECTED          // peer connected
    STATE_PR_SENT            // peer request sent
    STATE_PR_REFUSED         // peer request refused
    STATE_PR_WAIT            // peer request accepted
    STATE_PL_SENT            // peerlist sent
    STATE_PL_WAIT            // peerlist replied
    STATE_RP_SENT            // reroute protocol sent
    STATE_STANDBY            // waiting for connection confirmation
    STATE_RR_PENDING
)

type State int

var STATES = map[State]string{
    STATE_OFF:         "OFF",         // peer unreachable
    STATE_INITIALIZED: "Initialized", // peer started
    STATE_READY:       "Ready",       // peer ready to receive
    STATE_CONNECTED:   "Connected",   // peer connected
    STATE_PR_SENT:     "PR Sent",     // peer request sent
    STATE_PR_REFUSED:  "PR Refused",  // peer request refused
    STATE_PR_WAIT:     "PR Wait",     // peer request accepted
    STATE_PL_SENT:     "PL Sent",     // peerlist sent
    STATE_PL_WAIT:     "PL Wait",     // peerlist replied
    STATE_RP_SENT:     "RP Sent",     // reroute protocol sent
    STATE_STANDBY:     "Standby",     // waiting for connection confirmation
    STATE_RR_PENDING:  "RR Pending",  // Reroutelist proposal pending
}

type Peer struct {
    gorm.Model
    name        string       `json:"name"`
    ip          string       `json:"ip"`
    localIp     string       `json:"localIp"`
    port        int          `json:"port"`
    address     Address      `json:"address"`
    routes      []Route      `json:"routes"`
    peers       Peers        `json:"peers"`
    state       State        `json: "state"`
    stateInfo   string       // state explanation
    refreshed   time.Time    `json:"refreshed"`
    lastMsgId   uint64       // last message id
    score       float32      // efficiency score
    lastSent    time.Time    // connection watchdog
    lastRcvd    time.Time    // connection watchdog
    peerList    *PeerList    // list of peers (addresses only)
    rerouteList *RerouteList // rerouting proposal to peer
}

type PeerInfo struct {
    Name      string   `json:"name"`
    Address   Address  `json:"address"`
    IP        string   `json:"ip"`
    LocalIp   string   `json:"localIp"`
    Port      int      `json:"port"`
    Distance  Distance `json:"distance"`
    PeerCount int      `json:"peers"`
}

type ApiPeerInfo struct {
    Name      string   `json:"name"`
    Address   string   `json:"address"`
    IP        string   `json:"ip"`
    LocalIP   string   `json:"localIp"`
    Port      int      `json:"port"`
    Distance  Distance `json:"distance"`
    PeerCount int      `json:"peers"`
}

type Peers []Peer
type PeerList Peers

type PeersByDistance Peers
type PeersByName Peers
type PeersByAddress Peers

var Me *Peer

func NewPeer(name string, address Address, ip string, port int) *Peer {
    Logln("Peer created:", name, address, ip, port)
    return &Peer{name: name, ip: ip, address: address, port: port, peers: []Peer{}}
}

func (peer Peer) Distance(peer2 Peer) Distance {
    return peer.address.Distance(peer2.address)
}

func (state State) String() string {
    return STATES[state]
}

//send a protocol message to peer
func (peer Peer) send(message peerMessage) error {
    message.Destination = peer.Info()
    Logln("Messaging", peer, "with", message)
    return message.send()
}

//
// --- getter and setters
//

// return peer address
func (peer *Peer) Address() Address {
    return peer.address
}

// Update peer address
func (peer *Peer) SetAddress(address Address) *Peer {
    peer.address = address
    return peer
}

// Return peer name
func (peer *Peer) Name() string {
    return peer.name
}

// Set peer name
func (peer *Peer) SetName(name string) *Peer {
    peer.name = name
    return peer
}

// Return IP used by peer
func (peer *Peer) Ip() string {
    return peer.ip
}

func (peer *Peer) SetIp(ip string) *Peer {
    peer.ip = ip
    return peer
}

// Return port used by peer
func (peer *Peer) Port() int {
    return peer.port
}

// Return string representation of peer port
func (peer *Peer) PortString() string {
    return strconv.Itoa(peer.port)
}

// Update peer port
func (peer *Peer) SetPort(port int) *Peer {
    peer.port = port
    return peer
}

// Return peer state
func (peer *Peer) State() State {
    return peer.state
}

//Set peer state
func (peer *Peer) SetState(state State) *Peer {
    peer.state = state
    return peer
}

//Set peer stateInfo
func (peer *Peer) SetStateInfo(info string) *Peer {
    peer.stateInfo = info
    return peer
}

func (peerInfo PeerInfo) getConnectionString() string {
    ip := peerInfo.IP
    if len(peerInfo.LocalIp) > 0 {
        ip = peerInfo.LocalIp
    }
    return net.JoinHostPort(ip, strconv.Itoa(peerInfo.Port))
}

func (peer Peer) adminUrl() string {
    return "http://127.0.0.1:" + strconv.Itoa(*mgmtPort)
}

//
// --- peer lists
//

func (peer Peer) GetList() Peers {
    return peer.peers
}

func (peers Peers) searchByInfo(info PeerInfo) *Peer {
    for _, peer := range peers {
        if peer.ip == info.IP && peer.address == info.Address {
            return &peer
        }
    }
    return nil
}

func (peers Peers) searchByAddress(address Address) *Peer {
    for _, peer := range peers {
        if peer.address == address {
            return &peer
        }
    }
    return nil
}

func (peer *Peer) CheckList(peer2 *Peer) RerouteList {
    rerouteList := RerouteList{}
    for _, aPeer := range peer.GetList().FilterMe() {
        if peer.Distance(aPeer) > peer2.Distance(aPeer) {
            rerouteList.Get = append(rerouteList.Get, aPeer.Info())
        }
    }
    for _, aPeer := range peer2.GetList().FilterMe() {
        if peer.Distance(aPeer) < peer2.Distance(aPeer) {
            rerouteList.Gimme = append(rerouteList.Gimme, aPeer.Info())
        }
    }
    // truncate lists if more than half total exchanged
    half := (len(rerouteList.Gimme) + len(rerouteList.Get)) / 2
    if len(rerouteList.Gimme) > half {
        rerouteList.Gimme = (rerouteList.Gimme)[:half]
    }
    if len(rerouteList.Get) > half {
        rerouteList.Get = (rerouteList.Get)[:half]
    }
    sort.Sort(PeerInfosByZero(rerouteList.Gimme))
    sort.Sort(PeerInfosByZero(rerouteList.Get))
    fmt.Println("Reroutelist", rerouteList)
    return rerouteList
}

//
// --- Sorting
//

func (peers PeersByDistance) Less(i, j int) bool {
    return Me.Distance(peers[i]) < Me.Distance(peers[j])
}

func (peers PeersByDistance) Swap(i, j int) {
    peers[i], peers[j] = peers[j], peers[i]
}

func (peers PeersByDistance) Len() int {
    return len(peers)
}

func (peers PeersByName) Less(i, j int) bool {
    return strings.ToLower(peers[i].name) < strings.ToLower(peers[j].name)
}

func (peers PeersByName) Swap(i, j int) {
    peers[i], peers[j] = peers[j], peers[i]
}

func (peers PeersByName) Len() int {
    return len(peers)
}

func (peers PeersByAddress) Len() int {
    return len(peers)
}

func (peers PeersByAddress) Swap(i, j int) {
    peers[i], peers[j] = peers[j], peers[i]
}

func (peers PeersByAddress) Less(i int, j int) bool {
    for p := range peers[i].address {
        if peers[i].address[p] > peers[i].address[p] {
            return false
        }
    }
    return true
}

type PeerInfosByZero PeerInfos

// -- ordering relative to address ZERO

func (peerInfos PeerInfosByZero) Less(i, j int) bool {
    zero := &Address{0, 0, 0, 0}
    return zero.Distance(peerInfos[i].Address) < zero.Distance(peerInfos[j].Address)
}

func (peerInfos PeerInfosByZero) Swap(i, j int) {
    peerInfos[i], peerInfos[j] = peerInfos[j], peerInfos[i]
}

func (peerInfos PeerInfosByZero) Len() int {
    return len(peerInfos)
}

//
// --- utility
//

// GetPeers is used to obtain an initial list of peers
func GeneratePeers() Peers {
    // for now generate random peers
    for i := 0; i < *maxPeers; i++ {
        peer := &Peer{}
        peer.address = generateAddress()
        peer.name = fmt.Sprint(name.Generate(2, "-"))
        peer.SetIp("127.0.0.1")
        peer.port = DEFAULT_PORT + i + 1
        peer.state = STATE_READY // simulate already connected
        Me.peers = append(Me.peers, *peer)
        //peerCounter++
    }
    return Me.peers
}

func (peer *Peer) addPeer(newPeer Peer) (*Peer, error) {
    if len(peer.peers) > MAX_PEERS {
        return nil, errors.New("Peer limit reached")
    }
    peer.peers = append(peer.peers, newPeer)
    return peer, nil
}

//
// --- representation
//

func (peer Peer) String() string {
    //return fmt.Sprint([]interface{}{peer.name, peer.ip, peer.address, peer.state, peer.routes})
    return fmt.Sprintf("%s [%s] (%s|%s:%d) %s", peer.name, peer.address, peer.ip, peer.localIp, peer.port, peer.state)
}

func (peer Peer) Info() (info PeerInfo) {
    return PeerInfo{
        Name:      peer.name,
        Address:   peer.address,
        IP:        peer.ip,
        LocalIp:   peer.localIp,
        Port:      peer.port,
        Distance:  Me.Distance(peer),
        PeerCount: len(peer.peers),
    }
}

func (peerInfo PeerInfo) ApiInfo() ApiPeerInfo {
    return ApiPeerInfo{
        Name:      peerInfo.Name,
        Address:   peerInfo.Address.String(),
        IP:        peerInfo.IP,
        LocalIP:   peerInfo.LocalIp,
        Port:      peerInfo.Port,
        Distance:  peerInfo.Distance,
        PeerCount: peerInfo.PeerCount,
    }
}

func (peers Peers) Info() (infos PeerInfos) {
    for _, peer := range peers {
        infos = append(infos, peer.Info())
    }
    return infos
}

type PeerInfos []PeerInfo

func (infos PeerInfos) ApiPeerInfos() (apiPeerInfos []ApiPeerInfo) {
    for _, info := range infos {
        apiPeerInfos = append(apiPeerInfos, info.ApiInfo())
    }
    return apiPeerInfos
}

func (peers Peers) String() (output string) {
    for _, peer := range peers {
        output += fmt.Sprintf("- %s\n", peer)
    }
    return output
}

func equal(peer, peer2 Peer) bool {
    // only compare unvariant values
    if peer.address == peer2.address &&
        peer.ip == peer2.ip &&
        peer.port == peer2.port &&
        peer.name == peer2.name {
        return true
    }
    return false
}

//func notequal(peer, peer2 Peer) bool {
//}

func (peerInfo PeerInfo) String() string {
    return fmt.Sprintf("%s [%s] (%s|%s:%d)", peerInfo.Name, peerInfo.Address, peerInfo.IP, peerInfo.LocalIp, peerInfo.Port)
}

func (peerInfos PeerInfos) String() (output string) {
    for _, info := range peerInfos {
        output += fmt.Sprintf("- %s\n", info)
    }
    return output
}

func (peers Peers) FilterMe() (filteredPeers Peers) {
    return peers.Filter(Me.address)
}

func (peers Peers) Filter(address Address) (filteredPeers Peers) {
    for _, peer := range peers {
        if peer.address != address {
            filteredPeers = append(filteredPeers, peer)
        }
    }
    return filteredPeers
}

func (peers Peers) Bytes() []byte {
    buf := new(bytes.Buffer)
    encoder := gob.NewEncoder(buf)
    if err := encoder.Encode(peers); err != nil {
        return []byte("ERR")
    }
    fmt.Print("Peers encoded:", buf.Bytes())
    return buf.Bytes()
}

func (peerInfos PeerInfos) Bytes() []byte {
    buf := new(bytes.Buffer)
    encoder := gob.NewEncoder(buf)
    encoder.Encode(peerInfos)
    return buf.Bytes()
}

func (peerInfos PeerInfos) Peers() (peers Peers) {
    for _, info := range peerInfos {
        peers = append(peers, info.Peer())
    }
    return peers
}

func (info PeerInfo) Peer() (peer Peer) {
    peer = Peer{
        localIp:   info.LocalIp,
        name:      info.Name,
        ip:        info.IP,
        address:   info.Address,
        port:      info.Port,
        refreshed: time.Now(),
    }
    return peer
}
