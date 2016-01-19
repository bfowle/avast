package main

import (
    "fmt"
    "net/http"

    "github.com/gorilla/mux"
    consulapi "github.com/hashicorp/consul/api"
)

func (cr *ConsulRegistry) HealthHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
    vars := mux.Vars(r)

    options := &consulapi.QueryOptions{Datacenter: vars["dc"]}
    check, _, err := cr.health.Node(vars["name"], options)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(fmt.Sprintf("Consul endpoint failed: %v", err)))
        return nil, nil
    }

    return check, nil
}
