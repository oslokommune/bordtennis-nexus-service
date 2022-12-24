// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/oslokommune/bordtennis-nexus-service/pkg/core"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	// The websocket connection.
	Conn *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan []byte

	UnregisterFn func()
	BroadcastFn  func(message []byte)
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) ReadPump() {
	defer func() {
		c.UnregisterFn()

		err := c.Conn.Close()
		if err != nil {
			log.Printf("error: %v", err)
		}
	}()

	c.Conn.SetReadLimit(maxMessageSize)

	err := c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		log.Printf("error: %v", err)
	}

	c.Conn.SetPongHandler(func(string) error {
		err = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			log.Printf("error: %v", err)
		}

		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}

			break
		}

		err = validateMessage(message)
		if err != nil {
			log.Printf("invalid message: %v", err)

			continue
		}

		message = bytes.TrimSpace(bytes.ReplaceAll(message, newline, space))

		c.BroadcastFn(message)
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()

		err := c.Conn.Close()
		if err != nil {
			log.Printf("error: %v", err)
		}
	}()

	for {
		select {
		case message, ok := <-c.Send:
			err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Printf("error: %v", err)
			}

			if !ok {
				// The hub closed the channel.
				err = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.Printf("error: %v", err)
				}

				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			_, err = w.Write(message)
			if err != nil {
				log.Printf("error: %v", err)
			}

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				_, err = w.Write(newline)
				if err != nil {
					log.Printf("error: %v", err)
				}

				_, err = w.Write(<-c.Send)
				if err != nil {
					log.Printf("error: %v", err)
				}
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Printf("error: %v", err)
			}

			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func validateMessage(message []byte) error {
	var item core.Message

	err := json.Unmarshal(message, &item)
	if err != nil {
		return fmt.Errorf("unmarshalling: %w", err)
	}

	err = item.Validate()
	if err != nil {
		return fmt.Errorf("validating: %w", err)
	}

	return nil
}
