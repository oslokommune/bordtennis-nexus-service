// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "github.com/rs/zerolog/log"

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		event := log.Debug()

		select {
		case client := <-h.register:
			event.Str("event", "register")

			h.clients[client] = true
		case client := <-h.unregister:
			event.Str("event", "unregister")

			if _, ok := h.clients[client]; ok {
				event.Str("event", "unregister:ok")

				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			event.Str("event", "broadcast")
			event.RawJSON("message", message)

			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}

		event.Msg("hub:run()")
	}
}
