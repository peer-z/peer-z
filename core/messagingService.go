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
	"fmt"
	"github.com/go-chi/chi"
	"gorm.io/gorm"
	"net/http"
	"time"
)

//type MailDB struct {
//	Mailboxes []Mailbox
//}

type Mailbox struct {
	gorm.Model
	Name  string `json:"name"`
	Mails []Mail `gorm:"many2many:mailbox_messages",json:"messages"`
}

type Mail struct {
	gorm.Model
	//From      Address    `gorm:"type:string;default:''"`
	//To        Address    `gorm:"type:string;default:''"`
	DateSent     time.Time `json:"dateSent"`
	DateReceived time.Time `json:"dateReceived"`
	DateRead     time.Time `json:"dateRead"`
	Subject      string    `json:"subject"`
	Body         string    `json:"body"`
	Flags        uint64    `json:"flags"`
	Mailboxes    []Mailbox `gorm:"many2many:mailbox_messages",json:"mailboxes"`
}

var messagingService = Service{
	Info: ServiceInfo{
		Name:        "Messaging",
		Path:        "/messaging",
		Id:          messagingServiceID,
		Version:     0x0100, // version 1.00
		address:     messagingServiceAddress,
		Description: "A basic e-mail like messaging service for people to simply communicate",
		flags:       0,
		Component:   "messaging",
	},
	init: initMessaging,
	api: ApiDefinition{
		path:   "/mailboxes",
		router: messagingRouter,
	},
}

var mailDB []Mailbox
var inbox Mailbox

func initMessaging() {
	Logln("Messaging Started")
	db.AutoMigrate(&Mailbox{}, &Mail{})
	result := db.Find(&mailDB)
	if result.RowsAffected == 0 {
		var messages []Mail
		for i := 0; i < 10; i++ {
			messages = append(messages, Mail{
				Subject: fmt.Sprintf("Welcome to your new mailbox - msg %d", i),
				Body:    "Your new mailboxes were created and are ready for you to enjoy.\nWe hope you'll have good fun.",
			})
		}
		mailboxes := []Mailbox{
			Mailbox{Name: "Inbox",
				Mails: messages,
			},
			Mailbox{Name: "Outbox"},
			Mailbox{Name: "Sent"},
			Mailbox{Name: "Trash"},
		}
		db.Create(mailboxes)
		db.Commit()
		_ = db.Find(&mailDB)

	} else {
		fmt.Println(inbox)
	}
}

func messagingRouter(r chi.Router) {
	r.Get("/", listMailBoxes)
	r.Get("/{mailboxID}", listMail)
}

func listMailBoxes(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "text/html")

	Logln("Mailboxes requested")
	if err := SendJson(w, mailDB); err != nil {
		if err := SendPage(w, "errors/500.html"); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			Logln("Page Error: ", err)
		}

	}
}

func listMail(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "text/html")

	if mailboxID := chi.URLParam(req, "mailboxID"); mailboxID != "" {
		Loglnf("Mailbox %d requested", mailboxID)
		var mailbox Mailbox
		var messages []Mail
		db.Find(&mailbox, mailboxID)
		db.Model(&mailbox).Association("Mails").Find(&messages)
		if err := SendJson(w, messages); err != nil {
			err := SendPage(w, "errors/500.html")
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte(err.Error()))
				Logln("Page Error: ", err)
			}
		}
	} else {
		notFoundHandler(w, req)
	}
}
