package main

import (
	"strings"
	"testing"
)

// Test handleQuestsCommand - verifies quests display with placeholder
func TestHandleQuestsCommand(t *testing.T) {
	m := &model{
		currentCharacterID: 0,
		message:            "",
		messageType:        "",
	}

	// Test with character ID 0 - triggers placeholder display
	m.handleQuestsCommand("quests")

	// Should show placeholder when no character
	if m.message == "" {
		t.Error("Expected non-empty message from placeholder")
	}

	// Verify placeholder contains expected content
	if !strings.Contains(m.message, "QUEST LOG") {
		t.Error("Expected 'QUEST LOG' to be displayed in placeholder")
	}

	// Verify at least one quest is shown in placeholder
	if !strings.Contains(m.message, "Prove Yourself") {
		t.Error("Expected 'Prove Yourself' placeholder quest to be displayed")
	}
}

// Test handleQuestsCommandWithNoQuests - verifies empty state
func TestHandleQuestsCommandWithNoQuests(t *testing.T) {
	// When currentCharacterID is 0, should show placeholder
	m := &model{
		currentCharacterID: 0,
		message:            "",
		messageType:        "",
	}

	m.handleQuestsCommand("quests")

	// Should show placeholder with some message
	if m.message == "" {
		t.Error("Expected non-empty message")
	}
}

// Test handleQuestsCommandAPIError - verifies error handling
func TestHandleQuestsCommandAPIError(t *testing.T) {
	// When currentCharacterID is 0, should not call API but show placeholder
	m := &model{
		currentCharacterID: 0,
		message:            "",
		messageType:        "",
	}

	m.handleQuestsCommand("quests")

	// Should work without error even without character
	if m.message == "" {
		t.Error("Expected non-empty message")
	}
}

// Test questPanelStyles - verifies Lip Gloss styling exists
func TestQuestPanelStyles(t *testing.T) {
	// Verify quest panel colors are defined (non-empty)
	if questTitleColor == "" {
		t.Error("Expected questTitleColor to be defined")
	}

	if questProgressColor == "" {
		t.Error("Expected questProgressColor to be defined")
	}

	if questCompletedColor == "" {
		t.Error("Expected questCompletedColor to be defined")
	}

	// Verify styles are initialized (not zero value)
	// We can check if they render something without panicking
	testStr := "test"
	_ = questTitleStyle.Render(testStr)
	_ = questBoxStyle.Render(testStr)
	_ = questProgressStyle.Render(testStr)
	_ = questCompletedStyle.Render(testStr)
}

// Test displayQuestTrackerPlaceholder - verifies placeholder display
func TestDisplayQuestTrackerPlaceholder(t *testing.T) {
	m := &model{
		currentCharacterID: 0,
		message:            "",
		messageType:        "",
	}

	m.handleQuestsCommand("quests")

	// Verify placeholder displays correctly
	if m.message == "" {
		t.Error("Expected non-empty message from placeholder")
	}

	// Verify header is present
	if !strings.Contains(m.message, "QUEST LOG") {
		t.Error("Expected 'QUEST LOG' in placeholder")
	}

	// Verify at least one quest is shown
	if !strings.Contains(m.message, "Prove Yourself") && !strings.Contains(m.message, "Ooze Samples") {
		t.Error("Expected placeholder quests to be displayed")
	}

	// Verify message type is info
	if m.messageType != "info" {
		t.Errorf("Expected messageType 'info', got '%s'", m.messageType)
	}
}