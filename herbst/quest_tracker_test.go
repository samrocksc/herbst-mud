package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Test handleQuestsCommand - verifies quests display
func TestHandleQuestsCommand(t *testing.T) {
	// Create mock server for quests endpoint
	questsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Mock quest data - in real implementation this would come from DB
		json.NewEncoder(w).Encode(map[string]interface{}{
			"quests": []map[string]interface{}{
				{
					"id":          "quest_prove_yourself",
					"name":        "Prove Yourself",
					"description": "Kill 3 Scrap Rats to earn your place.",
					"status":      "in_progress",
					"objectives": []map[string]interface{}{
						{
							"description": "Kill Scrap Rat",
							"current":     2,
							"total":       3,
						},
					},
					"giver":   "Guard Marco",
					"rewards": "10 coins",
				},
				{
					"id":          "quest_ooze_samples",
					"name":        "Ooze Samples",
					"description": "Collect 5 glowing goo for Jane.",
					"status":      "available",
					"objectives": []map[string]interface{}{
						{
							"description": "Collect glowing goo",
							"current":     0,
							"total":       5,
						},
					},
					"giver":   "Scavenger Jane",
					"rewards": "repair_kit",
				},
			},
		})
	})
	server := httptest.NewServer(questsHandler)
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

	m.handleQuestsCommand("quests")

	if m.messageType != "info" {
		t.Errorf("Expected messageType 'info', got '%s'", m.messageType)
	}

	if m.message == "" {
		t.Error("Expected non-empty message")
	}

	// Verify quest names are displayed
	if !strings.Contains(m.message, "Prove Yourself") {
		t.Error("Expected 'Prove Yourself' quest to be displayed")
	}

	// Verify quest progress is shown
	if !strings.Contains(m.message, "2/3") {
		t.Error("Expected quest progress '2/3' to be displayed")
	}

	// Verify status is shown
	if !strings.Contains(m.message, "In Progress") && !strings.Contains(m.message, "in progress") {
		t.Error("Expected quest status to be displayed")
	}
}

// Test handleQuestsCommandWithNoQuests - verifies empty state
func TestHandleQuestsCommandWithNoQuests(t *testing.T) {
	// Create mock server returning empty quests
	questsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"quests": []map[string]interface{}{},
		})
	})
	server := httptest.NewServer(questsHandler)
	defer server.Close()

	oldBase := RESTAPIBase
	RESTAPIBase = server.URL
	defer func() { RESTAPIBase = oldBase }()

	m := &model{
		currentCharacterID: 1,
		message:            "",
		messageType:        "",
	}

	m.handleQuestsCommand("quests")

	// Should show message about no quests
	if m.messageType != "info" && m.messageType != "error" {
		t.Errorf("Expected messageType 'info' or 'error', got '%s'", m.messageType)
	}
}

// Test handleQuestsCommandAPIError - verifies error handling
func TestHandleQuestsCommandAPIError(t *testing.T) {
	// Create mock server that returns error
	questsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	})
	server := httptest.NewServer(questsHandler)
	defer server.Close()

	oldBase := RESTAPIBase
	RESTAPIBase = server.URL
	defer func() { RESTAPIBase = oldBase }()

	m := &model{
		currentCharacterID: 1,
		message:            "",
		messageType:        "",
	}

	m.handleQuestsCommand("quests")

	if m.messageType != "error" {
		t.Errorf("Expected messageType 'error', got '%s'", m.messageType)
	}
}

// Test questPanelStyles - verifies Lip Gloss styling
func TestQuestPanelStyles(t *testing.T) {
	// Verify quest panel colors are defined
	if questTitleColor == "" {
		t.Error("Expected questTitleColor to be defined")
	}

	if questProgressColor == "" {
		t.Error("Expected questProgressColor to be defined")
	}

	if questCompletedColor == "" {
		t.Error("Expected questCompletedColor to be defined")
	}

	// Verify style objects are not nil
	if questTitleStyle == nil {
		t.Error("Expected questTitleStyle to be defined")
	}

	if questBoxStyle == nil {
		t.Error("Expected questBoxStyle to be defined")
	}
}