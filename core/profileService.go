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
    "net/http"
    "regexp"
)

var profileService = Service{
    Info: ServiceInfo{
        Name:        "Profile",
        Id:          profileServiceID,
        Version:     0x0100, // version 1.00
        address:     profileServiceAddress,
        Description: "A basic profile service to introduce oneself",
        flags:       0,
    },
    init: initProfile,
    api: ApiDefinition{
        path:    "/profile",
        router: profileRouter,
    },
}

func initProfile() {
    Logln("Profile Started")
    //db.AutoMigrate(&Mailbox{}, &Message{})
    //result := db.Find(&messageDB)
    //if result.RowsAffected == 0 {
    //	var messages []Message
    //	for i:=0; i<10; i++ {
    //		messages = append(messages,Message{
    //			Subject: fmt.Sprintf("Welcome to your new mailbox - msg %d",i),
    //			Body:    "Your new mailboxes were created and are ready for you to enjoy.\nWe hope you'll have good fun.",
    //		})
    //	}
    //	mailboxes := []Mailbox{
    //		Mailbox{Name: "Inbox",
    //			Messages: messages,
    //		},
    //		Mailbox{Name: "Outbox"},
    //		Mailbox{Name: "Sent"},
    //		Mailbox{Name: "Trash"},
    //	}
    //	db.Create(mailboxes)
    //	db.Commit()
    //	_ = db.Find(&messageDB)
    //
    //} else {
    //	fmt.Println(inbox)
    //}
}

func profileRouter(r chi.Router) {
    r.Get("/", func(w http.ResponseWriter, req *http.Request) {
        w.Header().Set("Content-type", "text/html")

        exp := regexp.MustCompile("[/]+([^/]+)")
        var matches [][]string
        matches = exp.FindAllStringSubmatch(req.RequestURI, -1)

        fmt.Println(matches)

        //if strings.Compare(matches[0][1], "mailboxes") == 0 {
        //	if len(matches) == 1 {
        //		Logln("Mailboxes requested")
        //		err := SendJson(w, messageDB)
        //		if err != nil {
        //			err := SendPage(w, "errors/500.html")
        //			if err != nil {
        //				w.WriteHeader(500)
        //				w.Write([]byte(err.Error()))
        //				Logln("Page Error: ", err)
        //			}
        //		}
        //	} else {
        //		mailboxID, err := strconv.Atoi(matches[1][1])
        //		Loglnf("Mailbox %d requested", mailboxID)
        //		var mailbox Mailbox
        //		var messages []Message
        //		db.Find(&mailbox, mailboxID)
        //		db.Model(&mailbox).Association("Messages").Find(&messages)
        //		err = SendJson(w, messages)
        //		if err != nil {
        //			err := SendPage(w, "errors/500.html")
        //			if err != nil {
        //				w.WriteHeader(500)
        //				w.Write([]byte(err.Error()))
        //				Logln("Page Error: ", err)
        //			}
        //		}
        //	}
        //}
    })
}
