// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hub

import (
	"encoding/json"
	"fmt"

	"github.com/oslokommune/bordtennis-nexus-service/pkg/client"
	"github.com/oslokommune/bordtennis-nexus-service/pkg/core"
	"github.com/oslokommune/bordtennis-nexus-service/pkg/status"
	"github.com/rs/zerolog/log"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	gameStatus status.Data

	// Registered clients.
	clients map[*client.Client]bool

	// Inbound messages from the clients.
	Broadcast chan []byte

	// Register requests from the clients.
	Register chan *client.Client

	// Unregister requests from clients.
	Unregister chan *client.Client
}

func New() *Hub {
	statusData := status.Data{}
	status.Reset(&statusData)

	return &Hub{
		gameStatus: statusData,
		Register:   make(chan *client.Client),
		Unregister: make(chan *client.Client),
		Broadcast:  make(chan []byte),
		clients:    make(map[*client.Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		logEvent := log.Debug()

		select {
		case client := <-h.Register:
			logEvent.Str("event", "register")

			h.clients[client] = true

			client.send <- []byte(status.Serialize(h.gameStatus))
		case client := <-h.Unregister:
			logEvent.Str("event", "unregister")

			if _, ok := h.clients[client]; ok {
				logEvent.Str("event", "unregister:ok")

				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.Broadcast:
			msg := core.Message{}

			err := json.Unmarshal(message, &msg)
			if err != nil {
				logEvent.Err(err).Msg("json.Unmarshal")

				continue
			}

			logEvent.Str("event", "broadcast")
			logEvent.RawJSON("message", message)
			h.registerMessage(msg)

			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}

		logEvent.Msg("hub:Run()")
	}
}

func (h *Hub) registerMessage(msg core.Message) {
	switch msg.Type {
	case core.TypeBumpTeam:
		if msg.Payload == "1" {
			h.gameStatus.TeamOne++
		} else if "2" == msg.Payload {
			h.gameStatus.TeamTwo++
		} else {
			log.Warn().Msg(fmt.Sprintf("invalid payload %s, ignoring", msg.Payload))
		}
	case core.TypeReset:
		status.Reset(&h.gameStatus)
	}
}
