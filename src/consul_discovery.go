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
    "fmt"
    "strconv"
    "time"

    consul "github.com/hashicorp/consul/api"
    "github.com/hashicorp/consul/watch"
)

func (cr *ConsulRegistry) registerConsulWatch() (Watcher, error) {
    cw := &ConsulWatcher{
        Watchers: make(map[string]*watch.WatchPlan),
    }

    wp, err := watch.Parse(map[string]interface{}{"type": "services"})
    if err != nil {
        return nil, err
    }

    wp.Handler = cw.ServiceHandle
    go wp.Run(cr.Address)
    cw.WatchPlan = wp

    return cw, nil
}

func (cw *ConsulWatcher) serviceHandler(idx uint64, data interface{}) {
    entries, ok := data.([]*consul.ServiceEntry)
    fmt.Println("\n=================================================================")
    fmt.Printf("### %v :: SERVICE HANDLE ###\n", time.Now())
    if !ok {
        return
    }

    cs := &ConsulService{}

    for k, e := range entries {
        fmt.Println(k)

        fmt.Printf(" [Node] Node: %v (%v)\n", e.Node.Node, e.Node.Address)
        fmt.Printf(" [Service] ID: %v, Service: %v, Tags: %v, Port: %v, Address: %v\n",
            e.Service.ID, e.Service.Service, e.Service.Tags, e.Service.Port, e.Service.Address)
        for _, chk := range e.Checks {
            fmt.Printf(" [Check] Node: %v, CheckID: %v, Name: %v, Status: %v, Output: %v\n",
                chk.Node, chk.CheckID, chk.Name, chk.Status, chk.Output)
        }

        cs.Name = e.Service.Service
        cs.Nodes = append(cs.Nodes, &ServiceNode{
            Id: e.Service.ID,
            Address: e.Node.Address,
            Port: strconv.Itoa(e.Service.Port),
        })
        cs.Checks = e.Checks
    }

    consulRegistry.Lock()
    consulRegistry.Services[cs.Name] = cs
    consulRegistry.Unlock()
}

func (cw *ConsulWatcher) ServiceHandle(idx uint64, data interface{}) {
    services, ok := data.(map[string][]string)
    fmt.Printf(" %v -- SERVICES HANDLE\n", time.Now())
    fmt.Println(services)

    if !ok {
        return
    }

    // add new watchers
    for service, _ := range services {
        fmt.Printf("***** svc: %v *****\n", service)
        if _, ok := cw.Watchers[service]; ok {
            continue
        }

        wp, err := watch.Parse(map[string]interface{}{
            "type": "service",
            "service": service,
        })

        if err == nil {
            wp.Handler = cw.serviceHandler
            go wp.Run(consulRegistry.Address)
            cw.Watchers[service] = wp
        }
    }

    consulRegistry.RLock()
    rservices := consulRegistry.Services
    consulRegistry.RUnlock()

    // remove unknown services from registry
    for s, _ := range rservices {
        if _, ok := services[s]; !ok {
            consulRegistry.Lock()
            delete(consulRegistry.Services, s)
            consulRegistry.Unlock()
        }
    }

    // remove unknown services from watchers
    for s, w := range cw.Watchers {
        if _, ok := services[s]; !ok {
            w.Stop()
            delete(cw.Watchers, s)
        }
    }
}

//func (cw *ConsulWatcher) NodeHandle(idx uint64, data interface{}) {
//    nodes, ok := data.([]*consul.Node)
//    fmt.Printf(" %v -- NODES HANDLE\n", time.Now())
//    fmt.Println(nodes)
//    for _, n := range nodes {
//        node := &ClientNode{n.Node, n.Address}
//        fmt.Printf(" --> Node: %v (%v)\n", node.Name, node.Address)
//    }
//
//    if !ok {
//        return
//    }
//
//    //consulRegistry.RLock()
//    //rnodes := consulRegistry.Nodes
//    //consulRegistry.RUnlock()
//
//    // remove unknown nodes from registry
//    //for n, _ := range rnodes {
//    //    if _, ok := nodes[n]; !ok {
//    //        consulRegistry.Lock()
//    //        delete(consulRegistry.Nodes, n)
//    //        consulRegistry.Unlock()
//    //    }
//    //}
//
//    //// remove unknown nodes from watchers
//    //for n, w := range cw.Watchers {
//    //    if _, ok := nodes[n]; !ok {
//    //        w.Stop()
//    //        delete(cw.Watchers, n)
//    //    }
//    //}
//}

func (cw *ConsulWatcher) Stop() {
    if cw.WatchPlan == nil {
        return
    }
    cw.WatchPlan.Stop()
}
