// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hub

import (
	"encoding/json"

	"github.com/oslokommune/bordtennis-nexus-service/pkg/client"
	"github.com/oslokommune/bordtennis-nexus-service/pkg/core"
	"github.com/rs/zerolog/log"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	lobby      string
	gameStatus status

	// Registered clients.
	clients map[*client.Client]bool

	// Inbound messages from the clients.
	Broadcast chan []byte

	// Register requests from the clients.
	Register chan *client.Client

	// Unregister requests from clients.
	Unregister chan *client.Client
}

func New(lobby string) *Hub {
	hub := &Hub{
		lobby:      lobby,
		gameStatus: status{},
		Register:   make(chan *client.Client),
		Unregister: make(chan *client.Client),
		Broadcast:  make(chan []byte),
		clients:    make(map[*client.Client]bool),
	}

	hub.gameStatus.Reset()

	return hub
}

func (h *Hub) Run() {
	for {
		logEvent := log.Debug()

		select {
		case client := <-h.Register:
			logEvent.Str("event", "register")

			h.clients[client] = true

			msg := core.Message{
				Origin:  "server",
				Type:    core.TypeStatus,
				Payload: h.gameStatus.Serialize(),
			}

			rawMessage, err := json.Marshal(msg)
			if err != nil {
				logEvent.Err(err).Msg("failed to marshal message")

				continue
			}

			client.Send <- rawMessage
		case client := <-h.Unregister:
			logEvent.Str("event", "unregister")

			if _, ok := h.clients[client]; ok {
				logEvent.Str("event", "unregister:ok")

				delete(h.clients, client)
				close(client.Send)
			}
		case message := <-h.Broadcast:
			msg := core.Message{}

			err := json.Unmarshal(message, &msg)
			if err != nil {
				logEvent.Err(err).Msg("json.Unmarshal")

				continue
			}

			err = msg.Validate()
			if err != nil {
				logEvent.Err(err).Msg("msg.Validate")

				continue
			}

			logEvent.Str("event", "broadcast")
			logEvent.RawJSON("message", message)
			h.registerMessage(msg)

			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
		}

		logEvent.Msg("hub:Run()")
	}
}
