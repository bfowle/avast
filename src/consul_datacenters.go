package main

import (
    "fmt"
    "net/http"
)

func (cr *ConsulRegistry) DatacentersHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
    datacenters, err := cr.catalog.Datacenters()
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(fmt.Sprintf("Consul endpoint failed: %v", err)))
        return nil, nil
    }

    return datacenters, nil
}
