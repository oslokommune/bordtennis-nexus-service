package hub

import (
	"fmt"
	"strconv"

	"github.com/oslokommune/bordtennis-nexus-service/pkg/core"
	"github.com/rs/zerolog/log"
)

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
	case core.TypeChangeServer:
		if h.gameStatus.TeamOne+h.gameStatus.TeamTwo != 0 {
			log.Warn().Msg("cannot change server after first point")

			break
		}

		var err error

		h.gameStatus.InitialServe, err = strconv.Atoi(msg.Payload)
		if err != nil {
			log.Warn().Err(err).Msg(fmt.Sprintf("invalid payload %s, ignoring", msg.Payload))
		}
	case core.TypeReset:
		h.gameStatus.Reset()
	}

	log.Debug().
		Str("lobby", h.lobby).
		Str("type", msg.Type).
		Str("payload", msg.Payload).
		Str("status", h.gameStatus.Serialize()).
		Msg("status update")
}

type status struct {
	TeamOne      int `json:"teamOne"`
	TeamTwo      int `json:"teamTwo"`
	InitialServe int `json:"initialServe"`
}

func (s *status) Serialize() string {
	return fmt.Sprintf("%d;%d;%d", s.TeamOne, s.TeamTwo, s.InitialServe)
}

func (s *status) Reset() {
	s.TeamOne = 0
	s.TeamTwo = 0
	s.InitialServe = 1
}
