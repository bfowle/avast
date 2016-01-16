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
    "fmt"
    "net/http"

    "github.com/gorilla/mux"
    consul "github.com/hashicorp/consul/api"
)

func (cr *ConsulRegistry) HealthHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
    vars := mux.Vars(r)
    options := &consul.QueryOptions{Datacenter: vars["dc"]}

    check, _, err := cr.Health.Node(vars["name"], options)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(fmt.Sprintf("Consul endpoint failed: %v", err)))
        return nil, nil
    }

    return check, nil
}
