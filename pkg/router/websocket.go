package router

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/oslokommune/bordtennis-nexus-service/pkg/client"
	"github.com/oslokommune/bordtennis-nexus-service/pkg/hub"
)

func serveWebsocket(hub *hub.Hub, w http.ResponseWriter, r *http.Request, allowedHosts []string) {
	upgrader.CheckOrigin = originChecker(allowedHosts)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)

		return
	}

	currentClient := &client.Client{
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	hub.Register <- currentClient

	currentClient.UnregisterFn = func() {
		hub.Unregister <- currentClient
	}

	currentClient.BroadcastFn = func(message []byte) {
		hub.Broadcast <- message
	}

	go currentClient.WritePump()
	go currentClient.ReadPump()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func originChecker(allowedHosts []string) func(*http.Request) bool {
	return func(r *http.Request) bool {
		return contains(allowedHosts, r.Header.Get("Origin"))
	}
}
