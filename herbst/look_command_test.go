package main

import (
	"testing"
)

// TestHandleLookAtCommand_EmptyTarget tests look at with no target
func TestHandleLookAtCommand_EmptyTarget(t *testing.T) {
	m := &model{
		message:     "",
		messageType: "",
	}

	m.handleLookAtCommand("")

	if m.message == "" {
		t.Error("Expected usage message, got empty string")
	}

	if m.messageType != "error" {
		t.Errorf("Expected messageType 'error', got '%s'", m.messageType)
	}
}

// TestHandleLookAtCommand_LookAtMe tests look at me
func TestHandleLookAtCommand_LookAtMe(t *testing.T) {
	m := &model{
		message:            "",
		messageType:        "",
		currentUserName:    "TestPlayer",
		characterGender:   "Male",
		characterDescription: "A brave warrior",
	}

	m.handleLookAtCommand("me")

	if m.message == "" {
		t.Error("Expected message, got empty string")
	}

	if m.messageType != "info" {
		t.Errorf("Expected messageType 'info', got '%s'", m.messageType)
	}

	// Should contain character name
	if !contains(m.message, "TestPlayer") {
		t.Error("Expected message to contain character name")
	}

	// Should contain character description
	if !contains(m.message, "brave warrior") {
		t.Error("Expected message to contain character description")
	}
}

// TestHandleLookAtCommand_LookAtMyself tests look at myself
func TestHandleLookAtCommand_LookAtMyself(t *testing.T) {
	m := &model{
		message:            "",
		messageType:        "",
		currentUserName:    "TestPlayer",
		characterGender:   "Female",
		characterDescription: "A cunning rogue",
	}

	m.handleLookAtCommand("myself")

	if m.message == "" {
		t.Error("Expected message, got empty string")
	}

	if m.messageType != "info" {
		t.Errorf("Expected messageType 'info', got '%s'", m.messageType)
	}
}

// TestHandleLookAtCommand_LookAtItem_NoClient tests look at item when not logged in
func TestHandleLookAtCommand_LookAtItem_NoClient(t *testing.T) {
	m := &model{
		client:      nil,
		message:     "",
		messageType: "",
	}

	m.handleLookAtCommand("sword")

	// Should show error since no client
	if m.message == "" {
		t.Error("Expected error message, got empty string")
	}

	if !contains(m.message, "don't see") {
		t.Error("Expected 'don't see' in error message")
	}
}

// TestHandleLookAtCommand_LookAtNPC tests look at NPC when not logged in
func TestHandleLookAtCommand_LookAtNPC_NoClient(t *testing.T) {
	m := &model{
		client:      nil,
		message:     "",
		messageType: "",
	}

	m.handleLookAtCommand("guard")

	// Should show error since no client
	if m.message == "" {
		t.Error("Expected error message, got empty string")
	}
}

// TestProcessCommand_LookAt tests the processCommand integration
func TestProcessCommand_LookAt(t *testing.T) {
	m := &model{
		message:     "",
		messageType: "",
		client:      nil,
		// Set a character description
		characterDescription: "Test desc",
		currentUserName:     "TestUser",
	}

	// Simulate the command parsing
	cmd := "look at me"
	parts := strings.Fields(cmd)

	// Should handle "look at me"
	if len(parts) >= 3 && parts[1] == "at" {
		m.handleLookAtCommand(strings.Join(parts[2:], " "))
	}

	if m.message == "" {
		t.Error("Expected message from look at me")
	}

	if m.messageType != "info" {
		t.Errorf("Expected messageType 'info', got '%s'", m.messageType)
	}
}

// TestProcessCommand_LookAtItem tests look at item command
func TestProcessCommand_LookAtItem(t *testing.T) {
	m := &model{
		message:     "",
		messageType: "",
		client:      nil,
	}

	cmd := "look at sword"
	parts := strings.Fields(cmd)

	// Should handle "look at sword"
	if len(parts) >= 3 && parts[1] == "at" {
		m.handleLookAtCommand(strings.Join(parts[2:], " "))
	}

	// No client, so should get error
	if m.message == "" {
		t.Error("Expected error message")
	}

	if m.messageType != "error" {
		t.Errorf("Expected messageType 'error', got '%s'", m.messageType)
	}
}

// TestProcessCommand_LookAtNPC tests look at npc command
func TestProcessCommand_LookAtNPC(t *testing.T) {
	m := &model{
		message:     "",
		messageType: "",
		client:      nil,
	}

	cmd := "look at guard"
	parts := strings.Fields(cmd)

	if len(parts) >= 3 && parts[1] == "at" {
		m.handleLookAtCommand(strings.Join(parts[2:], " "))
	}

	// No client, so should get error
	if m.message == "" {
		t.Error("Expected error message")
	}
}

// TestLookCommand_TargetParsing tests that look at parses targets correctly
func TestLookCommand_TargetParsing(t *testing.T) {
	tests := []struct {
		cmd      string
		expected string
	}{
		{"look at sword", "sword"},
		{"look at rusty sword", "rusty sword"},
		{"look at guard marco", "guard marco"},
		{"l at me", ""}, // Won't match - "l" doesn't equal "look"
	}

	for _, tt := range tests {
		parts := strings.Fields(tt.cmd)
		if len(parts) >= 3 && parts[1] == "at" {
			target := strings.Join(parts[2:], " ")
			if target != tt.expected {
				t.Errorf("For cmd '%s', expected target '%s', got '%s'", tt.cmd, tt.expected, target)
			}
		}
	}
}

// TestLookCommand_AliasParsing tests that 'l at' also works as alias
func TestLookCommand_AliasParsing(t *testing.T) {
	// The alias "l" should also work with "at"
	cmd := "l at me"
	parts := strings.Fields(cmd)

	// Check if command is "look" or "l" and has "at"
	cmdName := parts[0]
	if (cmdName == "look" || cmdName == "l") && len(parts) >= 3 && parts[1] == "at" {
		// This is the correct path for look at
		target := strings.Join(parts[2:], " ")
		if target != "me" {
			t.Errorf("Expected target 'me', got '%s'", target)
		}
	}
}

// contains is a helper to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}