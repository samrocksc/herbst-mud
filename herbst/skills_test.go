package main

import (
	"strings"
	"testing"

	"herbst/db"
)

// TestSkillCommands verifies the skill command handlers work
func TestSkillCommands(t *testing.T) {
	tests := []struct {
		name     string
		cmd      string
		wantMsg  string
	}{
		{
			name:    "skills command initializes combat skill state",
			cmd:     "skills",
			wantMsg: "Combat Skills",
		},
		{
			name:    "swap-skill without args shows usage",
			cmd:     "/swap-skill",
			wantMsg: "Usage: /swap-skill",
		},
		{
			name:    "swap-skill with args without db shows connection error",
			cmd:     "/swap-skill slash fireball",
			wantMsg: "Database not connected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &model{
				client: nil,
			}

			m.processCommand(tt.cmd)

			if !strings.Contains(m.message, tt.wantMsg) {
				t.Errorf("processCommand(%q) = %q, want to contain %q", tt.cmd, m.message, tt.wantMsg)
			}
		})
	}
}

// TestAbilityDbTypes verifies the db package has the Ability type
func TestAbilityDbTypes(t *testing.T) {
	var _ *db.Ability
	var _ *db.AbilityCreate
	var _ *db.AbilityQuery
}