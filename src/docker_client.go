package main

import (
    "fmt"

    "github.com/docker/engine-api/client"
)

var dockerClient *DockerClient

type DockerClient struct {
    *client.Client
}

func newDockerClient() {
    cli, err := client.NewEnvClient()
    if err != nil {
        fmt.Println(err)
        return
    }
    dockerClient = &DockerClient{cli}
}
