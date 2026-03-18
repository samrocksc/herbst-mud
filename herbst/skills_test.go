package main

import (
	"strings"
	"testing"

	"herbst/db"
)

// TestSkillAndTalentCommands verifies the skill/talent command handlers work
func TestSkillAndTalentCommands(t *testing.T) {
	// This test verifies the commands can be parsed and the handlers can be invoked
	// without a database connection (graceful fallback)
	
	tests := []struct {
		name     string
		cmd      string
		wantMsg  string
	}{
		{
			name:    "skills command shows message about database",
			cmd:     "skills",
			wantMsg: "Database not connected",
		},
		{
			name:    "talents command shows message about database",
			cmd:     "talents",
			wantMsg: "Database not connected",
		},
		{
			name:    "swap-skill without args shows usage",
			cmd:     "/swap-skill",
			wantMsg: "Usage: /swap-skill",
		},
		{
			name:    "swap-talent without args shows usage",
			cmd:     "/swap-talent",
			wantMsg: "Usage: /swap-talent",
		},
		{
			name:    "swap-skill with args without db shows connection error",
			cmd:     "/swap-skill slash fireball",
			wantMsg: "Database not connected",
		},
		{
			name:    "swap-talent with args without db shows connection error",
			cmd:     "/swap-talent warrior leader",
			wantMsg: "Database not connected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &model{
				client: nil, // No database connection
			}
			
			m.processCommand(tt.cmd)
			
			if !strings.Contains(m.message, tt.wantMsg) {
				t.Errorf("processCommand(%q) = %q, want to contain %q", tt.cmd, m.message, tt.wantMsg)
			}
		})
	}
}

// TestSkillTalentDbTypes verifies the db package has Skill and Talent types
func TestSkillTalentDbTypes(t *testing.T) {
	// Just verify the db package has the expected types
	var _ *db.Skill
	var _ *db.Talent
	var _ *db.SkillCreate
	var _ *db.TalentCreate
	var _ *db.SkillQuery
	var _ *db.TalentQuery
}