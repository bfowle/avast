package main

import (
    "fmt"
    "net/http"
    "time"

    "github.com/docker/engine-api/types"
    "github.com/dustin/go-humanize"
    "github.com/gorilla/mux"
)

type ImageNode struct {
    ID          string              `json:"id"`
    ParentID    string              `json:"parentId"`
    RepoTags    []string            `json:"repoTags"`
    RepoDigests []string            `json:"repoDigests"`
    Created     string              `json:"created"`
    Size        string              `json:"size"`
    VirtualSize string              `json:"virtualSize"`
    Labels      map[string]string   `json:"labels"`
    Children    []*ImageNode        `json:"children"`
}

func (node *ImageNode) add(parent string, nodes []*ImageNode) {
    for _, n := range nodes {
        if n.ParentID == parent {
            node.Children = append(node.Children, n)
            n.add(n.ID, nodes)
        }
    }
}

func (dc *DockerClient) ImagesHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
    options := types.ImageListOptions{All: true}
    images, err := dc.ImageList(options)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(fmt.Sprintf("Docker engine endpoint failed: %v", err)))
        return nil, nil
    }

    nodes := make([]*ImageNode, len(images))
    for i, img := range images {
        nodes[i] = &ImageNode{
            img.ID,
            img.ParentID,
            img.RepoTags,
            img.RepoDigests,
            humanize.Time(time.Unix(img.Created, 0)),
            humanize.Bytes(uint64(img.Size)),
            humanize.Bytes(uint64(img.VirtualSize)),
            img.Labels,
            nil,
        }
    }

    root := &ImageNode{}
    root.add("", nodes)

    return root, nil
}

func (dc *DockerClient) HistoryHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
    vars := mux.Vars(r)

    history, err := dc.ImageHistory(vars["id"])
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(fmt.Sprintf("Docker engine endpoint failed: %v", err)))
        return nil, nil
    }

    return history, nil
}
