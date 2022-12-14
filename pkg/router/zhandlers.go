package router

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte("OK")); err != nil {
		log.Error().Err(err).Msg("failed to write health response")
	}
}
