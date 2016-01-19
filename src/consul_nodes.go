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
    nodes, _, err := cr.catalog.Nodes(options)
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
    node, _, err := cr.catalog.Node(vars["name"], options)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(fmt.Sprintf("Consul endpoint failed: %v", err)))
        return nil, nil
    }

    health, _, err := cr.health.Node(vars["name"], options)
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
        Services: services,
        Checks: health,
    }, nil
}
