package main

import (
    "fmt"
    "net/http"

    "github.com/docker/engine-api/types"
    "github.com/gorilla/mux"
)

func (dc *DockerClient) ContainersHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
    options := types.ContainerListOptions{All: true}
    containers, err := dc.ContainerList(options)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(fmt.Sprintf("Docker engine endpoint failed: %v", err)))
        return nil, nil
    }

    return containers, nil
}

func (dc *DockerClient) ContainerHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
    vars := mux.Vars(r)

    container, err := dc.ContainerInspect(vars["name"])
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(fmt.Sprintf("Docker engine endpoint failed: %v", err)))
        return nil, nil
    }

    return container, nil
}
