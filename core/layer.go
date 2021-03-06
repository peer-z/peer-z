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
    "flag"
    "gitlab.com/NebulousLabs/go-upnp"
    "log"
    "math/rand"
    "os"
    "os/exec"
    "os/signal"
    "strings"
    "testing"
)

const (
    VERSION              = "0.0.1"
    MAX_MESSAGES_IN      = 1000
    MAX_MESSAGES_OUT     = 1000
    MAX_MESSAGE_HANDLERS = 100
    MAX_MESSAGE_SENDERS  = 100
    MAX_PEERS            = 50
    MAX_ACTIVE_PEERS     = 500
    MGMT_PORT            = 80
    SERVICE_PORT         = 33301
    DEFAULT_PORT         = 33400
    MIN_HOPS             = 7
    MAX_TTL              = 30
    BASE_DIR             = ".peer-z"
)

var (
    debug        = flag.Bool("d", false, "Enable DEBUG mode")
    peerPort     = flag.Int("port", DEFAULT_PORT, "peer-z port to use when UPNP is unavailable")
    mgmtPort     = flag.Int("mgmt", MGMT_PORT, "peer-z port to use for peer management")
    relay        = flag.Bool("relay", false, "Run as a relay-only peer")
    maxPeers     = flag.Int("maxpeers", MAX_PEERS, "Max number of active peers")
    keyPath      = flag.String("key", "", "Path to private RSA key")
    keyChallenge = flag.String("keyChallenge", "", "Private RSA key challenge")
    //server       = flag.Bool("server", false, "Start as a server")
)

var baseDir string

var (
    signalChannel = make(chan os.Signal, 1)
)

//
// --- Initialization
//

func signalHandler() {
    for {
        select {
        case <-signalChannel:
            igd.Clear((uint16)(*peerPort))
            os.Exit(0)
        }
    }
}

func init() {
    func() {
        testing.Init()
    }()

    signal.Notify(signalChannel, os.Interrupt)

    messageIdCounter = uint64(rand.Int63())

    flag.Parse()

    if *relay {
        *peerPort = RELAY_PORT
    }

    configDir, err := os.UserConfigDir()
    if err != nil {
        log.Fatal("Houston we got a problem: Can't find config dir\n", err)
    }
    Logln("configdir: ", configDir)
    baseDir = configDir + "/" + BASE_DIR
    if err := os.Mkdir(baseDir, 0700); err != nil && !os.IsExist(err) {
        log.Fatal("Can't create app directory: ", err)
    }
    initDatabase()
    registerDefaultServices()

    log.Println("peer-z layer initializing, please wait...")

    if strings.Compare(*keyPath, "") == 0 {
        *keyPath = baseDir + "/id_peerz"
    }
    privKey, err = getPrivateKey(*keyPath)
    if err != nil {
        log.Fatal("Can't get private key: ", err)
    }

    // connect to router
    igd, err = upnp.Discover()
    if err != nil {
        log.Fatal(err)
    }

    go signalHandler()

    // discover external IP
    externalIP, err = igd.ExternalIP()
    if err != nil {
        log.Fatal(err)
    }

    if portIsForwarded, err := igd.IsForwardedUDP((uint16)(*peerPort)); err != nil || !portIsForwarded {
        // forward a port
        err = igd.Forward((uint16)(*peerPort), "peer-z")
        if err != nil {
            log.Fatal(err)
        }
    }

    address := generateAddress()

    if *relay {
        address = BroadcastAddress
        *mgmtPort = 33201
    }

    Me = NewPeer("me", address, externalIP, *peerPort)
    //GeneratePeers()
}

//
// --- Start it all
//

func Start() error {
    Me.addPeer(*NewPeer("relay", BroadcastAddress, "127.0.0.1", 33200))
    go adminServer()
    cmd := exec.Command("open", Me.adminUrl())
    cmd.Run()
    Logln("Listening to peers.")
    Listen()
    Logln("Terminating.")
    return nil
}

//
// --- Logging
//

func Logln(v ...interface{}) {
    if *debug {
        log.Println(v...)
    }
}

func Loglnf(f string, v ...interface{}) {
    Logf(f+"\n", v...)
}

func Logf(f string, v ...interface{}) {
    if *debug {
        log.Printf(f, v...)
    }
}
