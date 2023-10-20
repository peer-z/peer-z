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
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
)

type ServiceHandlers []ApiDefinition

var serviceHandlers ServiceHandlers

func adminServer() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", serverHandler)
	r.Get("/js/*", http.StripPrefix("/js/", http.FileServer(http.Dir("core/admin/dist/js"))).ServeHTTP)
	r.Get("/css/*", http.StripPrefix("/css/", http.FileServer(http.Dir("core/admin/dist/css"))).ServeHTTP)
	r.Get("/images/*", http.StripPrefix("/images/", http.FileServer(http.Dir("core/admin/dist/images"))).ServeHTTP)
	r.Get("/api/{apiName}", apiHandler)
	r.Get("/services/{serviceId}", serviceHandler)
	for _, definition := range serviceHandlers {
		if definition.router != nil {
			r.Route(definition.path, definition.router)
		}
	}
	r.NotFound(notFoundHandler)

	Logln("Server Started.")
	if err := http.ListenAndServe(":"+strconv.Itoa(*mgmtPort), r); err != nil {
		log.Fatal(err)
	}
	Logln("Server Stopped.")
}

func serviceHandler(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(req, "serviceId"))
	if err != nil {
		log.Fatal(err)
	}
	service, err := directory.SearchByID(int64(id))
	if err != nil {
		log.Fatal(err)
	}
	Loglnf("Service [%s] requested", service.Info.Name)
	err = SendPage(w, fmt.Sprintf("services/%d.html", id))
	if err != nil {
		err := SendPage(w, "errors/500.html")
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			Logln("Page Error: ", err)
		}
	}
}

func notFoundHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "text/html")
	w.WriteHeader(404)
	err := SendPage(w, "errors/404.html")
	Logln("Page not found:", req.RequestURI)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		Logln("Page Error: ", err)
	}
}

func apiHandler(w http.ResponseWriter, req *http.Request) {
	apiName := chi.URLParam(req, "apiName")
	switch apiName {
	case "me":
		Logln("Info requested")
		err := SendJson(w, Me.Info().ApiInfo())
		if err != nil {
			err := SendPage(w, "errors/500.html")
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte(err.Error()))
				Logln("Page Error: ", err)
			}
		}
	case "peers":
		Logln("Peers requested")
		sort.Sort(PeersByDistance(Me.GetList()))
		err := SendJson(w, PeerInfos(Me.GetList().Info()).ApiPeerInfos())
		if err != nil {
			err := SendPage(w, "errors/500.html")
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte(err.Error()))
				Logln("Page Error: ", err)
			}
		}
	case "services":
		Logln("Services requested")
		err := SendJson(w, directory.services)
		if err != nil {
			err := SendPage(w, "errors/500.html")
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte(err.Error()))
				Logln("Page Error: ", err)
			}
		}
	default:
		w.WriteHeader(404)
		err := SendPage(w, "errors/404.html")
		Logln("Page not found:", req.RequestURI)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			Logln("Page Error: ", err)
		}
	}
}

func serverHandler(w http.ResponseWriter, req *http.Request) {
	//exp := regexp.MustCompile("[/]+([^/]+)")
	//var matches [][]string
	//matches = exp.FindAllStringSubmatch(req.RequestURI, -1)
	//if len(matches) == 0 {
	Logln("Homepage requested")
	err := SendPage(w, "index.html")
	if err != nil {
		err := SendPage(w, "errors/500.html")
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			Logln("Page Error: ", err)
		}
	}
	//} else {
	//	w.WriteHeader(404)
	//	err := SendPage(w, "errors/404.html")
	//	Logln("Page not found:", req.RequestURI)
	//	if err != nil {
	//		w.WriteHeader(500)
	//		w.Write([]byte(err.Error()))
	//		Logln("Page Error: ", err)
	//	}
	//}
}

func SendJson(w http.ResponseWriter, data interface{}) error {
	json, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.Header().Set("Content-type", "application/json")
	w.Write(json)
	return nil
}

func SendPage(w http.ResponseWriter, pageName string) error {
	r, err := getPageReader(pageName)
	if err != nil {
		return err
	}
	w.Header().Set("Content-type", "text/html")
	for {
		b, err := r.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return err
			}
			w.Write([]byte(b))
			break
		}
		w.Write([]byte(b))
	}
	return nil
}

func getPageReader(name string) (*bufio.Reader, error) {
	page, err := os.Open("core/admin/dist/" + name)
	if err != nil {
		return nil, err
	}
	return bufio.NewReader(page), nil
}

func (serviceHandlers *ServiceHandlers) AddHandler(definition ApiDefinition) {
	*serviceHandlers = append(*serviceHandlers, definition)
}
