package main

import (
    "encoding/json"
    "fmt"
    "io"
    "time"

    "github.com/docker/engine-api/types"
    "github.com/docker/engine-api/types/events"
)

type DockerEvent struct {
    From        string          `json:"from"`
    Type        string          `json:"type"`
    Data        *events.Message `json:"data"`
    Timestamp   time.Time       `json:"timestamp"`
}

func (dc *DockerClient) EchoEvents() {
    options := types.EventsOptions{}
    r, err := dc.Events(options)
    if err != nil {
        fmt.Println(err)
    }
    defer r.Close()

    d := json.NewDecoder(r)
    messages := make(chan *DockerEvent)

    go func() {
        for {
            var data events.Message
            if err := d.Decode(&data); err != nil {
                if err == io.EOF {
                    break
                }
                fmt.Println(err)
            }
            messages <- &DockerEvent{
                "docker",
                "event",
                &data,
                time.Now(),
            }
        }
    }()

    for {
        select {
        case msg := <-messages:
            data, err := json.Marshal(msg)
            if err != nil {
                fmt.Println(err)
            }
            wsHub.broadcast <- data
        }
    }
}
