package main

import (
    "fmt"
    "net/http"
)

func (dc *DockerClient) InfoHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
    info, err := dc.Info()
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(fmt.Sprintf("Docker engine endpoint failed: %v", err)))
        return nil, nil
    }

    return info, nil
}
