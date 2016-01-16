// Copyright 2016 Brett Fowle
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var (
	router     *mux.Router
	addr       string
	apiVersion string
	indexTpl   = template.Must(template.ParseFiles("index.html"))
)

func setHeaders(w http.ResponseWriter, headers map[string]string) {
	for k, v := range headers {
		w.Header().Set(http.CanonicalHeaderKey(k), v)
	}
}

func processEnv() {
	if addr = os.Getenv("AVAST_ADDR"); addr == "" {
		addr = ":8080"
	}
	if apiVersion = os.Getenv("AVAST_API_VERSION"); apiVersion == "" {
		apiVersion = "v1"
	}
}

func startWebserver() {
	processEnv()

	router = mux.NewRouter()
	router.HandleFunc("/", wrap(indexHandler))
	router.HandleFunc("/ws", wrap(wsHandler))

	dockerRouter := router.PathPrefix(fmt.Sprintf("/api/%v/docker", apiVersion)).Subrouter()
	dockerRouter.HandleFunc("/containers", wrap(dockerContainersHandler))
	dockerRouter.HandleFunc("/container/{name}", wrap(dockerContainerHandler))
	dockerRouter.HandleFunc("/images", wrap(dockerImagesHandler))
	dockerRouter.HandleFunc("/history/{id}", wrap(dockerHistoryHandler))
	dockerRouter.HandleFunc("/info", wrap(dockerInfoHandler))

	consulRouter := router.PathPrefix(fmt.Sprintf("/api/%v/consul", apiVersion)).Subrouter()
	consulRouter.HandleFunc("/datacenters", wrap(consulRegistry.DatacentersHandler))
	consulRouter.HandleFunc("/nodes", wrap(consulRegistry.NodesHandler))
	consulRouter.HandleFunc("/nodes/{dc}", wrap(consulRegistry.NodesHandler))
	consulRouter.HandleFunc("/node/{name}", wrap(consulRegistry.NodeHandler))
	consulRouter.HandleFunc("/health/{name}", wrap(consulRegistry.HealthHandler))
	consulRouter.HandleFunc("/health/{name}/{dc}", wrap(consulRegistry.HealthHandler))

	http.Handle("/", router)
	loggedRouter := handlers.CombinedLoggingHandler(os.Stdout, router)
	panic(http.ListenAndServe(addr, handlers.CompressHandler(loggedRouter)))
}

func indexHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	indexTpl.Execute(w, r.Host)

	return nil, nil
}

func wrap(handler func(w http.ResponseWriter, r *http.Request) (interface{}, error)) func(w http.ResponseWriter, r *http.Request) {
	f := func(w http.ResponseWriter, r *http.Request) {
		setHeaders(w, map[string]string{
			"Content-Type":                "text/html; charset=utf8",
			"Access-Control-Allow-Origin": "*",
		})

		obj, err := handler(w, r)

	HAS_ERR:
		if err != nil {
			code := http.StatusInternalServerError
			errMsg := err.Error()
			if strings.Contains(errMsg, "Permission denied") {
				code = http.StatusForbidden
			}

			w.WriteHeader(code)
			w.Write([]byte(err.Error()))
			return
		}

		if obj != nil {
			buf, err := json.Marshal(obj)
			if err != nil {
				goto HAS_ERR
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(buf)
		}
	}

	return f
}
