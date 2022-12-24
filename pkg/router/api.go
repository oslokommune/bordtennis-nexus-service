package router

import (
	"net/http"

	"github.com/oslokommune/bordtennis-nexus-service/pkg/hub"
	"github.com/rs/zerolog/log"
)

func New(allowedHosts []string) *http.ServeMux {
	router := http.NewServeMux()
	lobbies := make(map[string]*hub.Hub)

	router.HandleFunc("/health", healthCheckHandler)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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

		serveWebsocket(lobbies[lobby], w, r, allowedHosts)
	})

	return router
}
