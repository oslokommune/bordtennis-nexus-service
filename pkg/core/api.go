package core

import (
	"fmt"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

const (
	// TypeBumpTeam is used to increase the score of a team by 1. Bump team expects "1" or "2" as payload.
	TypeBumpTeam = "bump-team"

	// TypeChangeServer is used to change the initial server. Change server expects "1" or "2" as payload.
	TypeChangeServer = "change-server"

	// Reset is a special type that is used to reset the score to 0-0. Reset expects no payload
	TypeReset = "reset"

	// TypeStatus is used to get the current score
	TypeStatus = "status"
)

var allTypes = []string{TypeBumpTeam, TypeReset, TypeStatus, TypeChangeServer}

type Message struct {
	// Origin identifies the source of the event, usually an UUID
	Origin  string `json:"origin"`
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

var reStatusPayload = regexp.MustCompile(`^\d+,\d+,[1-2]$`)

func (m Message) Validate() error {
	err := validation.ValidateStruct(&m,
		validation.Field(&m.Origin, validation.Required, is.UUIDv4),
		validation.Field(&m.Type, validation.Required, validation.In(TypeBumpTeam, TypeReset, TypeStatus, TypeChangeServer)),
	)
	if err != nil {
		return fmt.Errorf("validating generic fields: %w", err)
	}

	var v validation.Rule

	switch m.Type {
	case TypeBumpTeam:
		v = validation.In("1", "2")
	case TypeChangeServer:
		v = validation.In("1", "2")
	case TypeReset:
		v = validation.In("")
	case TypeStatus:
		v = validation.Match(reStatusPayload)
	}

	return validation.ValidateStruct(&m, validation.Field(&m.Payload, v))
}
