package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const dummyUUID = "5e0f39f5-74d9-4ddd-9c8f-e624c0a6b057"

func TestMessageValidation(t *testing.T) {
	testCases := []struct {
		name        string
		withType    string
		withPayload string
		expectError bool
	}{
		{
			name:        "bump-team should pass if payload is 1",
			withType:    TypeBumpTeam,
			withPayload: "1",
			expectError: false,
		},
		{
			name:        "bump-team should pass if payload is 2",
			withType:    TypeBumpTeam,
			withPayload: "2",
			expectError: false,
		},
		{
			name:        "Should fail if payload is 3",
			withType:    TypeBumpTeam,
			withPayload: "3",
			expectError: true,
		},
		{
			name:        "bump-team should fail if payload is 0",
			withType:    TypeBumpTeam,
			withPayload: "0",
			expectError: true,
		},

		{
			name:        "change-server should pass if payload is 1",
			withType:    TypeChangeServer,
			withPayload: "1",
			expectError: false,
		},
		{
			name:        "change-server should pass if payload is 2",
			withType:    TypeChangeServer,
			withPayload: "2",
			expectError: false,
		},
		{
			name:        "change-server should fail if payload is 3",
			withType:    TypeChangeServer,
			withPayload: "3",
			expectError: true,
		},
		{
			name:        "change-server should fail if payload is 0",
			withType:    TypeChangeServer,
			withPayload: "0",
			expectError: true,
		},

		{
			name:        "reset should pass if payload is empty",
			withType:    TypeReset,
			withPayload: "",
			expectError: false,
		},
		{
			name:        "reset should fail if payload is not empty",
			withType:    TypeReset,
			withPayload: "1",
			expectError: true,
		},

		{
			name:        "status should pass if payload is 0,0,1",
			withType:    TypeStatus,
			withPayload: "0,0,1",
			expectError: false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			msg := Message{
				Origin:  dummyUUID,
				Type:    tc.withType,
				Payload: tc.withPayload,
			}

			err := msg.Validate()

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
