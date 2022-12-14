package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/oslokommune/bordtennis-nexus-service/pkg/client"
	"github.com/oslokommune/bordtennis-nexus-service/pkg/hub"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const defaultPort = "3000"

var allowedHosts = strings.Split(os.Getenv("ALLOWED_HOSTS"), ";")

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	lobbies := make(map[string]*hub.Hub)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		event := log.Info()

		lobby := r.URL.Path[1:]

		if r.URL.Path == "/" {
			lobby = "default"
		}

		event.Str("lobby", lobby)

		if _, ok := lobbies[lobby]; !ok {
			event.Msg("Creating new lobby")

			lobbies[lobby] = hub.New()
			go lobbies[lobby].Run()
		}

		event.Msg("Serving existing lobby")

		client.ServeWebsocket(lobbies[lobby], w, r, allowedHosts)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write([]byte("OK")); err != nil {
			log.Error().Err(err).Msg("failed to write health response")
		}
	})

	log.Info().
		Str("Port", defaultPort).
		Str("AllowedHosts", fmt.Sprintf("%v+", allowedHosts)).
		Msg("Starting server")

	err := http.ListenAndServe(":3000", nil) // #nosec: G114
	if err != nil {
		log.Fatal().Err(err).Msg("error")
	}
}
