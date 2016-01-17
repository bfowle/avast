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

    consul "github.com/hashicorp/consul/api"
    "github.com/hashicorp/consul/watch"
)

var consulRegistry *ConsulRegistry
var consulWatcher *ConsulWatcher

type ClientNode struct {
    Name        string  `json:"name"`
    Address     string  `json:"address"`
}

type ConsulNode struct {
    *ClientNode
    Services  []*consul.AgentService  `json:"services"`
    Checks    []*consul.HealthCheck   `json:"checks"`
}

type ServiceNode struct {
    Id      string  `json:"id"`
    Address string  `json:"address"`
    Port    string  `json:"port"`
}

type ConsulService struct {
    Name    string
    Nodes   []*ServiceNode
    Checks  []*consul.HealthCheck
}

type ConsulWatcher struct {
    WatchPlan   *watch.WatchPlan
    Watchers    map[string]*watch.WatchPlan
}

type Watcher interface {
    Stop()
}

type ConsulRegistry struct {
    Address     string
    Client      *consul.Client
    Agent       *consul.Agent
    Catalog     *consul.Catalog
    Health      *consul.Health
    Services    map[string]*ConsulService
    Nodes       []*ConsulNode
    sync.RWMutex
}

func registerConsul() {
    config := consul.DefaultConfig()
    c, err := consul.NewClient(config)
    if err != nil {
    }

    consulRegistry = &ConsulRegistry{
        Address:    config.Address,
        Client:     c,
        Agent:      c.Agent(),
        Catalog:    c.Catalog(),
        Health:     c.Health(),
        Services:   make(map[string]*ConsulService),
    }

    // Watchers: key, keyprefix, services, nodes, service, checks, event
    // - https://github.com/hashicorp/consul/blob/master/watch/funcs.go
    // - https://github.com/hashicorp/consul/blob/master/watch/funcs_test.go
    //
    // Checks: {status: passing|warning|failing|critical}

    consulWatcher = &ConsulWatcher{
        Watchers: make(map[string]*watch.WatchPlan),
    }

    consulRegistry.registerConsulWatch("services")
    consulRegistry.registerConsulWatch("nodes")
    consulRegistry.registerConsulWatch("checks")
}
