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
    "strings"
)

const (
    ALL = iota
    REPORT_BY_NAME
    REPORT_BY_DESCRIPTION
    REPORT_BY_ID
    REPORT_CATEGORIES
)

type ServiceReport [REPORT_CATEGORIES][]*Service

var directory Directory

//
// --- Handle services
//

// Register a new service
func (directory *Directory) Register(service Service) error {
    if _, count := directory.Search(service); count > 0 {
        return errors.New("Service already registered")
    }
    service.Info.flags |= SERVICE_NEW
    directory.services = append(directory.services, service)
    return nil
}

// Abandon a service
func (directory Directory) Deregister(service Service) error {
    var report ServiceReport
    var count int
    if report, count = directory.Search(service); count == 0 {
        return errors.New("Service not found")
    }
    report[REPORT_BY_ID][0].Info.address = AddressNone
    report[REPORT_BY_ID][0].Info.flags |= SERVICE_ABANDONED | SERVICE_UPDATED
    return nil
}

//
// --- Find services
//

func (directory Directory) Search(serviceSearched Service) (report ServiceReport, count int) {
    for _, service := range directory.services {
        if strings.Compare(serviceSearched.Info.Name, service.Info.Name) == 0 {
            report.add(REPORT_BY_NAME, &service)
        }
        if strings.Compare(serviceSearched.Info.Description, service.Info.Description) == 0 {
            report.add(REPORT_BY_DESCRIPTION, &service)
        }
        if serviceSearched.Info.Id == service.Info.Id {
            report.add(REPORT_BY_ID, &service)
            count++
        }
    }
    return report, count
}

func (directory Directory) SearchByID(id int64) (*Service, error) {
    for _, service := range directory.services {
        if service.Info.Id == id {
            return &service, nil
        }
    }
    return nil, errors.New("Service not found")
}

func (directory Directory) SearchByName(name string) (*Service, error) {
    for _, service := range directory.services {
        if strings.Compare(service.Info.Name, name) == 0 {
            return &service, nil
        }
    }
    return nil, errors.New("Service not found")
}

func (directory Directory) SearchByDescription(description string) (services []*Service) {
    for _, service := range directory.services {
        if strings.Contains(service.Info.Description, description) {
            services = append(services, &service)
        }
    }
    return services
}

//
// --- Reports
//

func (report *ServiceReport) add(category int, service *Service) {
    report[category] = append(report[category], service)
}
