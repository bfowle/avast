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

type ConsulRegistry struct {
    addr        string
    client      *consulapi.Client
    agent       *consulapi.Agent
    catalog     *consulapi.Catalog
    health      *consulapi.Health
    services    map[string]*consulapi.ServiceEntry
    nodes       []*ConsulNode
    sync.RWMutex
}

func registerConsul() {
    config := consulapi.DefaultConfig()
    c, err := consulapi.NewClient(config)
    if err != nil {
        fmt.Println(err)
        return
    }

    consulRegistry = &ConsulRegistry{
        addr:       config.Address,
        client:     c,
        agent:      c.Agent(),
        catalog:    c.Catalog(),
        health:     c.Health(),
        services:   make(map[string]*consulapi.ServiceEntry),
    }

    // Watchers: key, keyprefix, services, nodes, service, checks, event
    // - https://github.com/hashicorp/consul/blob/master/watch/funcs.go
    // Checks: {status: passing|warning|failing|critical}
    // @TODO: on watch results, add to "trend" to show stats over time
}
