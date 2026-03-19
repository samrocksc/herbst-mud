package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test handleSkillsCommand - verifies skills display
func TestHandleSkillsCommand(t *testing.T) {
	// Create mock server for skills endpoint
	skillsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"skills": map[string]interface{}{
				"blades":  map[string]interface{}{"level": 15, "bonus": "+10%"},
				"staves":  map[string]interface{}{"level": 5, "bonus": "+0%"},
				"knives":  map[string]interface{}{"level": 20, "bonus": "+10%"},
				"martial": map[string]interface{}{"level": 0, "bonus": "+0%"},
			},
		})
	})
	server := httptest.NewServer(skillsHandler)
	defer server.Close()

	// Override RESTAPIBase for testing
	oldBase := RESTAPIBase
	RESTAPIBase = server.URL
	defer func() { RESTAPIBase = oldBase }()

	m := &model{
		currentCharacterID: 1,
		message:             "",
		messageType:         "",
	}

	m.handleSkillsCommand("skills")

	if m.messageType != "info" {
		t.Errorf("Expected messageType 'info', got '%s'", m.messageType)
	}

	if m.message == "" {
		t.Error("Expected non-empty message")
	}

	// Verify skills are displayed
	expectedSubstrings := []string{"blades", "15", "+10%"}
	for _, substr := range expectedSubstrings {
		if !containsSubstring(m.message, substr) {
			t.Errorf("Expected message to contain '%s', got: %s", substr, m.message)
		}
	}
}

// Test handleSkillsCommand without character
func TestHandleSkillsCommandNoCharacter(t *testing.T) {
	m := &model{
		currentCharacterID: 0,
		message:           "",
		messageType:       "",
	}

	m.handleSkillsCommand("skills")

	if m.messageType != "error" {
		t.Errorf("Expected messageType 'error', got '%s'", m.messageType)
	}

	if !containsSubstring(m.message, "need to be playing") {
		t.Errorf("Expected error about playing, got: %s", m.message)
	}
}

// Test handleTalentsCommand - verifies talents display
func TestHandleTalentsCommand(t *testing.T) {
	// Create mock server for talents endpoint
	talentsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Simulate slots array with JSON
		slots := [5]*struct {
			Slot       int    `json:"slot"`
			TalentID   int    `json:"talent_id"`
			Name       string `json:"name"`
			Description string `json:"description"`
		}{nil, nil, nil, nil, nil}

		// Slot 1 has a talent
		slots[1] = &struct {
			Slot       int    `json:"slot"`
			TalentID   int    `json:"talent_id"`
			Name       string `json:"name"`
			Description string `json:"description"`
		}{1, 1, "Warrior's Might", "Increase strength by 5"}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"character_id": 1,
			"slots":        slots,
		})
	})
	server := httptest.NewServer(talentsHandler)
	defer server.Close()

	oldBase := RESTAPIBase
	RESTAPIBase = server.URL
	defer func() { RESTAPIBase = oldBase }()

	m := &model{
		currentCharacterID: 1,
		message:            "",
		messageType:         "",
	}

	m.handleTalentsCommand("talents")

	if m.messageType != "info" {
		t.Errorf("Expected messageType 'info', got '%s'", m.messageType)
	}

	// Verify talent display
	if !containsSubstring(m.message, "Warrior's Might") {
		t.Errorf("Expected message to contain talent name, got: %s", m.message)
	}

	// Verify empty slots shown
	if !containsSubstring(m.message, "empty") {
		t.Errorf("Expected message to show empty slots, got: %s", m.message)
	}
}

// Test handleTalentsCommand without character
func TestHandleTalentsCommandNoCharacter(t *testing.T) {
	m := &model{
		currentCharacterID: 0,
		message:           "",
		messageType:       "",
	}

	m.handleTalentsCommand("talents")

	if m.messageType != "error" {
		t.Errorf("Expected messageType 'error', got '%s'", m.messageType)
	}
}

// Test handleTalentEquipCommand equip
func TestHandleTalentEquipCommand(t *testing.T) {
	equipCalled := false

	talentsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		equipCalled = true
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "slot": 2})
	})
	server := httptest.NewServer(talentsHandler)
	defer server.Close()

	oldBase := RESTAPIBase
	RESTAPIBase = server.URL
	defer func() { RESTAPIBase = oldBase }()

	m := &model{
		currentCharacterID: 1,
		message:           "",
		messageType:       "",
	}

	m.handleTalentEquipCommand("talent equip 1 2")

	if m.messageType != "success" {
		t.Errorf("Expected messageType 'success', got '%s'", m.messageType)
	}

	if !equipCalled {
		t.Error("Expected equip endpoint to be called")
	}
}

// Test handleTalentEquipCommand unequip
func TestHandleTalentEquipCommandUnequip(t *testing.T) {
	unequipCalled := false

	talentsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		unequipCalled = true
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "slot": 2})
	})
	server := httptest.NewServer(talentsHandler)
	defer server.Close()

	oldBase := RESTAPIBase
	RESTAPIBase = server.URL
	defer func() { RESTAPIBase = oldBase }()

	m := &model{
		currentCharacterID: 1,
		message:           "",
		messageType:       "",
	}

	m.handleTalentEquipCommand("talent unequip 2")

	if m.messageType != "success" {
		t.Errorf("Expected messageType 'success', got '%s'", m.messageType)
	}

	if !unequipCalled {
		t.Error("Expected unequip endpoint to be called")
	}
}

// Test handleTalentEquipCommand invalid slot
func TestHandleTalentEquipCommandInvalidSlot(t *testing.T) {
	m := &model{
		currentCharacterID: 1,
		message:           "",
		messageType:       "",
	}

	m.handleTalentEquipCommand("talent equip 1 5") // Invalid slot 5

	if m.messageType != "error" {
		t.Errorf("Expected messageType 'error', got '%s'", m.messageType)
	}
}

// Test handleTalentEquipCommand invalid usage
func TestHandleTalentEquipCommandInvalidUsage(t *testing.T) {
	m := &model{
		currentCharacterID: 1,
		message:           "",
		messageType:       "",
	}

	m.handleTalentEquipCommand("talent") // No subcommand

	if m.messageType != "error" {
		t.Errorf("Expected messageType 'error', got '%s'", m.messageType)
	}
}

// Test handleSkillEquipCommand
func TestHandleSkillEquipCommand(t *testing.T) {
	m := &model{
		currentCharacterID: 1,
		message:           "",
		messageType:       "",
	}

	m.handleSkillEquipCommand("skill equip 1")

	// Skills are always active, should give info message
	if m.messageType != "info" {
		t.Errorf("Expected messageType 'info', got '%s'", m.messageType)
	}

	if !containsSubstring(m.message, "always active") {
		t.Errorf("Expected message about skills being active, got: %s", m.message)
	}
}

// Helper function
func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && (s[:len(substr)] == substr || containsSubstring(s[1:], substr)))
}