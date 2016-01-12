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
    "net/http"
    "time"

    "github.com/gorilla/websocket"
)

const (
    // time allowed to write a message to the peer
    writeWait = 10 * time.Second
    // time allowed to read the next pong message from the peer
    pongWait = 60 * time.Second
    // send pings to peer with this period (must be less than pongWait)
    pingPeriod = (pongWait * 9) / 10
    // maximum message size allowed from peer
    maxMessageSize = 512
)

type hub struct {
    connections map[*connection]bool
    broadcast   chan []byte
    register    chan *connection
    unregister  chan *connection
}

var wsHub = hub{
    broadcast:   make(chan []byte),
    register:    make(chan *connection),
    unregister:  make(chan *connection),
    connections: make(map[*connection]bool),
}

func (h *hub) run() {
    for {
        select {
        case c := <-h.register:
            h.connections[c] = true
        case c := <-h.unregister:
            if _, ok := h.connections[c]; ok {
                delete(h.connections, c)
                close(c.send)
            }
        case m := <-h.broadcast:
            for c := range h.connections {
                select {
                case c.send <- m:
                default:
                    close(c.send)
                    delete(h.connections, c)
                }
            }
        }
    }
}

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

type connection struct {
    ws *websocket.Conn
    send chan []byte
}

// readPump pumps messages from the websocket connection to the hub
func (c *connection) readPump() {
    defer func() {
        wsHub.unregister <- c
        c.ws.Close()
    }()

    c.ws.SetReadLimit(maxMessageSize)
    c.ws.SetReadDeadline(time.Now().Add(pongWait))
    c.ws.SetPongHandler(func(string) error {
        c.ws.SetReadDeadline(time.Now().Add(pongWait))
        return nil
    })

    for {
        _, message, err := c.ws.ReadMessage()
        if err != nil {
            break
        }
        wsHub.broadcast <- message
    }
}

// write writes a message with the given message type and payload
func (c *connection) write(mt int, payload []byte) error {
    c.ws.SetWriteDeadline(time.Now().Add(writeWait))
    return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection
func (c *connection) writePump() {
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        ticker.Stop()
        c.ws.Close()
    }()

    for {
        select {
        case message, ok := <-c.send:
            if !ok {
                c.write(websocket.CloseMessage, []byte{})
                return
            }
            if err := c.write(websocket.TextMessage, message); err != nil {
                return
            }
        case <-ticker.C:
            if err := c.write(websocket.PingMessage, []byte{}); err != nil {
                return
            }
        }
    }
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
        return
    }

    c := &connection{ws, make(chan []byte, 256)}
    wsHub.register <- c

    go c.writePump()
    go c.echoEvents()
    c.readPump()
}
