package core

import "fmt"

const (
	// TypeBumpTeam is used to increase the score of a team by 1. Bump team expects "team1" or "team2" as payload.
	TypeBumpTeam = "bump-team"

	// Reset is a special type that is used to reset the score to 0-0. Reset expects no payload
	TypeReset = "reset"

	// TypeStatus is used to get the current score
	TypeStatus = "status"
)

var allTypes = []string{TypeBumpTeam, TypeReset, TypeStatus}

type Message struct {
	// Origin identifies the source of the event, usually an UUID
	Origin  string `json:"origin"`
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

func (m *Message) Validate() error {
	if !contains(allTypes, m.Type) {
		return fmt.Errorf("invalid type: %s", m.Type)
	}

	return nil
}
