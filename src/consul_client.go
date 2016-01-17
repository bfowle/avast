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
    "sync"

    consulapi "github.com/hashicorp/consul/api"
)

var consulRegistry *ConsulRegistry

type ClientNode struct {
    Name        string  `json:"name"`
    Address     string  `json:"address"`
}

type ConsulNode struct {
    *ClientNode
    Services  []*consulapi.AgentService  `json:"services"`
    Checks    []*consulapi.HealthCheck   `json:"checks"`
}

type ServiceNode struct {
    Id      string  `json:"id"`
    Address string  `json:"address"`
    Port    string  `json:"port"`
}

type ConsulRegistry struct {
    Address     string
    Client      *consulapi.Client
    Agent       *consulapi.Agent
    Catalog     *consulapi.Catalog
    Health      *consulapi.Health
    Services    map[string]*consulapi.ServiceEntry
    Nodes       []*ConsulNode
    sync.RWMutex
}

func registerConsul() {
    config := consulapi.DefaultConfig()
    c, err := consulapi.NewClient(config)
    if err != nil {
    }

    consulRegistry = &ConsulRegistry{
        Address:    config.Address,
        Client:     c,
        Agent:      c.Agent(),
        Catalog:    c.Catalog(),
        Health:     c.Health(),
        Services:   make(map[string]*consulapi.ServiceEntry),
    }

    // Watchers: key, keyprefix, services, nodes, service, checks, event
    // - https://github.com/hashicorp/consul/blob/master/watch/funcs.go
    // Checks: {status: passing|warning|failing|critical}
    // @TODO: on watch results, add to "trend" to show stats over time
}
