package main

import (
	"testing"
)

// TestDebugModeCommand tests the debug command handler
func TestDebugModeCommand(t *testing.T) {
	m := &model{
		debugMode: false,
	}

	// Test debug on
	m.handleDebugCommand("debug on")
	if !m.debugMode {
		t.Error("Expected debugMode to be true after 'debug on'")
	}

	// Test debug off
	m.handleDebugCommand("debug off")
	if m.debugMode {
		t.Error("Expected debugMode to be false after 'debug off'")
	}

	// Test debug toggle with different variations
	m.handleDebugCommand("debug true")
	if !m.debugMode {
		t.Error("Expected debugMode to be true after 'debug true'")
	}

	m.handleDebugCommand("debug false")
	if m.debugMode {
		t.Error("Expected debugMode to be false after 'debug false'")
	}

	// Test invalid subcommand
	m.handleDebugCommand("debug invalid")
	if m.debugMode {
		t.Error("Expected debugMode to remain false after invalid command")
	}
}

// TestDebugModeView tests that debug info appears in View when enabled
func TestDebugModeView(t *testing.T) {
	m := &model{
		debugMode:    true,
		currentRoom:  42,
		characterHP:         100,
		characterMaxHP:      100,
		characterStamina:    50,
		characterMaxStamina: 50,
		characterMana:       25,
		characterMaxMana:    25,
	}

	view := m.View()
	expectedRoomID := "Room: 42"
	if !stringContains(view, expectedRoomID) {
		t.Errorf("Expected View to contain '%s' when debugMode is true, got: %s", expectedRoomID, view)
	}

	// Test when debug is off - room ID should not appear
	m.debugMode = false
	view = m.View()
	if stringContains(view, "Room: 42") {
		t.Errorf("Expected View NOT to contain 'Room: 42' when debugMode is false, got: %s", view)
	}
}

// TestDebugModeStatus tests that debug mode status message is correct
func TestDebugModeStatus(t *testing.T) {
	m := &model{
		debugMode: true,
	}

	// Check debug on status message
	m.handleDebugCommand("debug on")
	if m.message != "Debug mode: ON (Room ID will show in status bar)" {
		t.Errorf("Unexpected message: %s", m.message)
	}

	m.handleDebugCommand("debug off")
	if m.message != "Debug mode: OFF" {
		t.Errorf("Unexpected message: %s", m.message)
	}

	// Check status when no argument provided
	m.debugMode = true
	m.handleDebugCommand("debug")
	if m.message != "Debug mode: ON (Room ID visible in status bar)" {
		t.Errorf("Unexpected message for status check: %s", m.message)
	}

	m.debugMode = false
	m.handleDebugCommand("debug")
	expectedMsg := "Debug mode: OFF\nUsage: debug on | debug off"
	if m.message != expectedMsg {
		t.Errorf("Expected message '%s', got: %s", expectedMsg, m.message)
	}
}

// Helper function to check if a string contains a substring
func stringContains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && stringContainsAt(s, substr))
}

func stringContainsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}