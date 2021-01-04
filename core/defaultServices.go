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
    "log"
)

//
// Default Service Definitions --- DO NOT CHANGE
//

const (
    unusedServiceID = iota
    messagingServiceID
    profileServiceID
    webServiceID
    chatServiceID
    storageServiceID
    checkerServiceID
    chessServiceID
)

var (
    messagingServiceAddress = Address{0, 0, 0, messagingServiceID}
    profileServiceAddress   = Address{0, 0, 0, profileServiceID}
    webServiceAddress       = Address{0, 0, 0, webServiceID}
    chatServiceAddress      = Address{0, 0, 0, chatServiceID}
    storageServiceAddress   = Address{0, 0, 0, storageServiceID}
    checkerServiceAddress   = Address{0, 0, 0, checkerServiceID}
    chessServiceAddress     = Address{0, 0, 0, chessServiceID}
)

var defaultServices = []Service{
    messagingService,
    profileService,
    webService,
    chatService,
    storageService,
    checkerService,
    chessService,
}

// Internal function to register initial services
func registerDefaultServices() {
    for _, service := range defaultServices {
        if err := directory.Register(service); err != nil {
            log.Fatalf("Failed registering service %d [%s]", service.Info.Id, service.Info.Name)
        }
        Loglnf("Registering service %d [%s]", service.Info.Id, service.Info.Name)
        service.StartService()
    }
}

//
// ---
//
