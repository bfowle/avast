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
    consulapi "github.com/hashicorp/consul/api"
)

func (cr *ConsulRegistry) NodesHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
    vars := mux.Vars(r)
    options := &consulapi.QueryOptions{Datacenter: vars["dc"]}

    nodes, _, err := cr.Catalog.Nodes(options)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(fmt.Sprintf("Consul endpoint failed: %v", err)))
        return nil, nil
    }

    var cnodes []*ClientNode
    for _, n := range nodes {
        cnodes = append(cnodes, &ClientNode{
            Name: n.Node,
            Address: n.Address,
        })
    }

    return cnodes, nil
}

func (cr *ConsulRegistry) NodeHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
    vars := mux.Vars(r)
    options := &consulapi.QueryOptions{Datacenter: vars["dc"]}

    node, _, err := cr.Catalog.Node(vars["name"], options)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(fmt.Sprintf("Consul endpoint failed: %v", err)))
        return nil, nil
    }

    health, _, err := cr.Health.Node(vars["name"], options)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(fmt.Sprintf("Consul endpoint failed: %v", err)))
        return nil, nil
    }

    services := make([]*consulapi.AgentService, 0, len(node.Services))
    for  _, s := range node.Services {
        services = append(services, s)
    }

    return &ConsulNode{
        ClientNode: &ClientNode{
            Name: node.Node.Node,
            Address: node.Node.Address,
        },
        Services:   services,
        Checks:     health,
    }, nil
}
