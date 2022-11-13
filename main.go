package main

import (
	"log"
	"net/http"
)

const defaultPort = "3000"

func main() {
	lobbies := make(map[string]*Hub)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)

		lobby := r.URL.Path[1:]

		if r.URL.Path == "/" {
			lobby = "default"
		}

		if _, ok := lobbies[lobby]; !ok {
			lobbies[lobby] = newHub()
			go lobbies[lobby].run()
		}

		serveWs(lobbies[lobby], w, r)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write([]byte("OK")); err != nil {
			log.Println(err)
		}
	})

	log.Println("Listening on port", defaultPort)

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
