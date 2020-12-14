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
	"crypto/rsa"
	"net/http"
)

const (
	SERVICE_NEW = 1 << iota
	SERVICE_RENAMED
	SERVICE_UPDATED
	SERVICE_ABANDONED
)

const (
	SVC_VERSION = 0x0100
)

type ServiceStarter interface {
	StartService()
}

type serviceManager interface {
	ServiceStarter
}

type ApiDefinition struct {
	path string
	handler func(http.ResponseWriter, *http.Request)
}

type ServiceInfo struct {
	Id          int64  `json:"id"`
	Version     uint16 `json:"version"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Port        int    `json:"port"`
	address     Address
	flags       uint32
}

type Service struct {
	Info ServiceInfo `json:"info"`
	serviceManager
	api  ApiDefinition
	init func()
	key  rsa.PrivateKey
}

func (service Service) StartService() {
	Loglnf("Starting service %s", service.Info.Name)
	serviceHandlers.AddHandler(service.api)
	if service.init != nil {
		service.init()
	}
}

type Directory struct {
	services []Service
}

//func (service Service) Info() (info ServiceInfo) {
//	return service.info
//}
