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
    "encoding/json"
    "fmt"
    "time"

    consulapi "github.com/hashicorp/consul/api"
    "github.com/hashicorp/consul/watch"
)

type Watcher struct {
    addr        string
    wp          *watch.WatchPlan
    watchers    map[string]*watch.WatchPlan
}

type WatchEvent struct {
    Type        string      `json:"type"`
    Data        interface{} `json:"data"`
    Timestamp   time.Time   `json:"timestamp"`
}

func (cr *ConsulRegistry) registerWatcher(watchType string) error {
    w, err := newWatcher(cr.Address, watchType)
    if err != nil {
        fmt.Println(err)
        return err
    }
    go w.Run()

    return nil
}

func (w *Watcher) registerServiceWatcher(service string) error {
    wp, err := watch.Parse(map[string]interface{}{
        "type": "service",
        "service": service,
    })
    if err != nil {
        return err
    }

    wp.Handler = func(idx uint64, data interface{}) {
        switch d := data.(type) {
        // service
        case []*consulapi.ServiceEntry:
            for _, i := range d {
                fmt.Printf("[ %v ]\t%v\n", time.Now(), i)
                broadcastData("service", &i)

                consulRegistry.Lock()
                consulRegistry.Services[i.Service.Service] = i
                consulRegistry.Unlock()
            }
        }
    }

    go wp.Run(w.addr)
    w.watchers[service] = wp
    return nil
}

func newWatcher(addr string, watchType string) (*Watcher, error) {
    wp, err := watch.Parse(map[string]interface{}{
        "type": watchType,
    })
    if err != nil {
        return nil, err
    }

    w := &Watcher{
        addr,
        wp,
        make(map[string]*watch.WatchPlan),
    }

    wp.Handler = func(idx uint64, data interface{}) {
        // @TODO: type switch seems to convert back to interface{}
        // if applying multiple types on the case (e.g. case []*A, []*B
        // it would be nice to combine these case statements for similar
        // types; try using reflect.TypeOf instead, perhaps
        switch d := data.(type) {
        // nodes
        case []*consulapi.Node:
            for _, i := range d {
                fmt.Printf("[ %v ]\t%v\n", time.Now(), i)
                broadcastData(watchType, &i)
            }
        // checks
        case []*consulapi.HealthCheck:
            for _, i := range d {
                fmt.Printf("[ %v ]\t%v\n", time.Now(), i)
                broadcastData(watchType, &i)
            }
        // services
        case map[string][]string:
            for i, _ := range d {
                fmt.Printf("[ %v ]\t%v\n", time.Now(), i)
                if _, ok := w.watchers[i]; ok {
                    continue
                }
                w.registerServiceWatcher(i)
            }

            consulRegistry.RLock()
            rs := consulRegistry.Services
            consulRegistry.RUnlock()

            // remove unknown services from registry
            for s, _ := range rs {
                if _, ok := d[s]; !ok {
                    consulRegistry.Lock()
                    delete(consulRegistry.Services, s)
                    consulRegistry.Unlock()
                }
            }

            // remove unknown services from watchers
            for i, svc := range w.watchers {
                if _, ok := d[i]; !ok {
                    svc.Stop()
                    delete(w.watchers, i)
                }
            }
        }
    }

    return w, nil
}

func broadcastData(watchType string, data interface{}) {
    evt := &WatchEvent {
        watchType,
        data,
        time.Now(),
    }

    msg, err := json.Marshal(evt)
    if err != nil {
        fmt.Println(err)
    }
    wsHub.broadcast <- msg
}


func (w *Watcher) Run() {
    w.wp.Run(w.addr)
}

func (w *Watcher) Stop() {
    if w.wp == nil {
        return
    }
    w.wp.Stop()
}

func (c *connection) echoDiscovery() {
    consulRegistry.registerWatcher("nodes")
    consulRegistry.registerWatcher("checks")
    consulRegistry.registerWatcher("services")
}
