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
    "flag"
    "net/http"
    "text/template"

    "github.com/gorilla/mux"
)

var addr = flag.String("addr", ":8080", "http service address")
var indexTpl = template.Must(template.ParseFiles("index.html"))

func webHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    indexTpl.Execute(w, r.Host)
}

func main() {
    flag.Parse()
    newClient()
    go wsHub.run()

    router := mux.NewRouter()
    router.HandleFunc("/", webHandler)
    router.HandleFunc("/ws", wsHandler)
    router.HandleFunc("/containers/list", containersHandler)
    router.HandleFunc("/container/{name}/inspect", containerHandler)
    router.HandleFunc("/images/list", imagesHandler)
    router.HandleFunc("/history/{id}", historyHandler)
    http.Handle("/", router)

    panic(http.ListenAndServe(*addr, nil))
}
