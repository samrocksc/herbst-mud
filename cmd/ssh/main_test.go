package main

import (
	"testing"
	tea "github.com/charmbracelet/bubbletea"
)

func TestModelInit(t *testing.T) {
	model := newModel()
	if model == nil {
		t.Error("New model should not be nil")
	}
}

func TestModelUpdate(t *testing.T) {
	model := newModel()
	
	// Test down arrow
	msg := tea.KeyMsg{Type: tea.KeyDown}
	newModel, _ := model.Update(msg)
	
	if newModel == nil {
		t.Error("Updated model should not be nil")
	}
}

func TestQuitting(t *testing.T) {
	model := newModel()
	
	// Test quit key
	msg := tea.KeyMsg{Type: tea.KeyCtrlC}
	newModel, cmd := model.Update(msg)
	
	if cmd == nil {
		t.Error("Expected quit command")
	}
	
	if newModel == nil {
		t.Error("Model should not be nil after quit")
	}
}

func TestView(t *testing.T) {
	model := newModel()
	view := model.View()
	
	if len(view) == 0 {
		t.Error("View should not be empty")
	}
	
	if !contains(view, "MUD Server") {
		t.Error("View should contain MUD Server")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr || 
		(len(s) > len(substr) && s[len(s)-len(substr):] == substr) ||
		(len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
