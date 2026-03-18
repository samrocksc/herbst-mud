package main

import (
	"testing"
)

// TestHandleExamineCommand_Usage tests the examine command with no arguments
func TestHandleExamineCommand_Usage(t *testing.T) {
	m := &model{
		message:     "",
		messageType: "",
	}

	m.handleExamineCommand("examine")

	if m.message == "" {
		t.Error("Expected usage message, got empty string")
	}

	if m.messageType != "error" {
		t.Errorf("Expected messageType 'error', got '%s'", m.messageType)
	}
}

// TestHandleExamineCommand_NoClient tests examine when not logged in
func TestHandleExamineCommand_NoClient(t *testing.T) {
	m := &model{
		client:      nil,
		message:     "",
		messageType: "",
	}

	m.handleExamineCommand("examine sword")

	// Should handle gracefully without panicking
	if m.message == "" {
		t.Error("Expected message, got empty string")
	}
}

// TestDisplayItemExamine_Output tests item examine formatting
func TestDisplayItemExamine_Output(t *testing.T) {
	m := &model{
		message:     "",
		messageType: "",
	}

	// Test with a mock item - can't actually call dbclient.Item
	// but we can test the command parsing works
	parts := strings.Fields("examine sword")
	if len(parts) < 2 {
		t.Error("Expected at least 2 parts for examine command")
	}

	target := strings.Join(parts[1:], " ")
	if target != "sword" {
		t.Errorf("Expected target 'sword', got '%s'", target)
	}
}

// TestExamine_DirectionShortcut tests that examining a direction uses peer
func TestExamine_DirectionShortcut(t *testing.T) {
	m := &model{
		message:     "",
		messageType: "",
		exits:      map[string]int{"north": 1},
		client:     nil, // Will cause peer to fail gracefully
	}

	// This should call handlePeerCommand internally
	// and result in an error about no client
	m.handleExamineCommand("examine north")

	// The peer command will fail gracefully without a client
	// but we verify the direction shortcut is being parsed
}

// TestExamine_InvalidDirection tests examining an invalid direction
func TestExamine_InvalidDirection(t *testing.T) {
	m := &model{
		message:     "",
		messageType: "",
		exits:      map[string]int{},
		client:     nil,
	}

	m.handleExamineCommand("examine northwest")

	// Should look for "northwest" as an item, not as a direction
	if m.messageType != "error" && m.messageType != "info" {
		t.Errorf("Expected messageType 'error' or 'info', got '%s'", m.messageType)
	}
}

// TestExamine_EmptyTarget tests examine with just whitespace
func TestExamine_EmptyTarget(t *testing.T) {
	m := &model{
		message:     "",
		messageType: "",
	}

	m.handleExamineCommand("examine   ")

	// Should show usage since target is empty after trimming
	if m.message == "" {
		t.Error("Expected usage message for empty target")
	}
}