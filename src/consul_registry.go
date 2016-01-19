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

func newConsulRegistry() {
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
}
