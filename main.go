package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/oslokommune/bordtennis-nexus-service/pkg/router"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const defaultPort = "3000"

var allowedHosts = strings.Split(os.Getenv("ALLOWED_HOSTS"), ";")

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log.Info().
		Str("Port", defaultPort).
		Str("AllowedHosts", fmt.Sprintf("%v+", allowedHosts)).
		Msg("Starting server")

	router := router.New(allowedHosts)

	err := http.ListenAndServe(":3000", router) // #nosec: G114
	if err != nil {
		log.Fatal().Err(err).Msg("error")
	}
}
