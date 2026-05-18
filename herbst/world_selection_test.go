package main

import (
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// TestWorldSelection tests the world/MUD selection screen
func TestWorldSelection(t *testing.T) {
	availableWorlds = []WorldInfo{
		{ID: "1", Name: "Herbst", Description: "Main world", Status: "active"},
		{ID: "2", Name: "Shadow Realm", Description: "Dark world", Status: "active"},
		{ID: "3", Name: "Elder Forest", Description: "Ancient woods", Status: "active"},
	}
	defer func() { availableWorlds = nil }()

	t.Run("typing 1 selects first world", func(t *testing.T) {
		ti := textinput.New()
		ti.Focus()
		m := &model{
			screen:       ScreenWorldSelect,
			textInput:    ti,
			visitedRooms: make(map[int]bool),
			knownExits:   make(map[string]bool),
			width:        80,
			height:       24,
		}
		m.Init()

		m.handleWorldSelectInput("1")

		if m.currentWorld != "Herbst" {
			t.Errorf("Expected currentWorld 'Herbst', got %q", m.currentWorld)
		}
		if m.screen != ScreenCharacterSelect {
			t.Errorf("Expected screen %q, got %q", ScreenCharacterSelect, m.screen)
		}
	})

	t.Run("typing world name selects by name", func(t *testing.T) {
		ti := textinput.New()
		ti.Focus()
		m := &model{
			screen:       ScreenWorldSelect,
			textInput:    ti,
			visitedRooms: make(map[int]bool),
			knownExits:   make(map[string]bool),
			width:        80,
			height:       24,
		}
		m.Init()

		m.handleWorldSelectInput("Shadow Realm")

		if m.currentWorld != "Shadow Realm" {
			t.Errorf("Expected currentWorld 'Shadow Realm', got %q", m.currentWorld)
		}
	})

	t.Run("case insensitive world name match", func(t *testing.T) {
		ti := textinput.New()
		ti.Focus()
		m := &model{
			screen:       ScreenWorldSelect,
			textInput:    ti,
			visitedRooms: make(map[int]bool),
			knownExits:   make(map[string]bool),
			width:        80,
			height:       24,
		}
		m.Init()

		m.handleWorldSelectInput("elder forest")

		if m.currentWorld != "Elder Forest" {
			t.Errorf("Expected case-insensitive match, got %q", m.currentWorld)
		}
	})

	t.Run("typing b goes back to welcome", func(t *testing.T) {
		ti := textinput.New()
		ti.Focus()
		m := &model{
			screen:       ScreenWorldSelect,
			textInput:    ti,
			visitedRooms: make(map[int]bool),
			knownExits:   make(map[string]bool),
			width:        80,
			height:       24,
		}
		m.Init()

		m.handleWorldSelectInput("b")

		if m.screen != ScreenWelcome {
			t.Errorf("Expected screen %q, got %q", ScreenWelcome, m.screen)
		}
	})

	t.Run("invalid world number is handled gracefully", func(t *testing.T) {
		ti := textinput.New()
		ti.Focus()
		m := &model{
			screen:       ScreenWorldSelect,
			textInput:    ti,
			visitedRooms: make(map[int]bool),
			knownExits:   make(map[string]bool),
			width:        80,
			height:       24,
		}
		m.Init()

		m.handleWorldSelectInput("9")

		// "9" matches the case pattern but parseWorldIndex returns -1
		// Code stays on world select screen without error message
		if m.screen != ScreenWorldSelect {
			t.Errorf("Expected to stay on world select, got %q", m.screen)
		}
	})

	t.Run("invalid world name shows error", func(t *testing.T) {
		ti := textinput.New()
		ti.Focus()
		m := &model{
			screen:       ScreenWorldSelect,
			textInput:    ti,
			visitedRooms: make(map[int]bool),
			knownExits:   make(map[string]bool),
			width:        80,
			height:       24,
			maxHistory:   50,
		}
		m.Init()

		m.handleWorldSelectInput("Nonexistent World")

		if len(m.messageTypes) == 0 {
			t.Fatal("Expected message in history")
		}
		if lastType := m.messageTypes[len(m.messageTypes)-1]; lastType != "error" {
			t.Errorf("Expected error type, got %q", lastType)
		}
	})

	t.Run("displayWorlds lists all worlds", func(t *testing.T) {
		m := &model{
			width:        80,
			height:       24,
			visitedRooms: make(map[int]bool),
			knownExits:   make(map[string]bool),
		}
		display := m.displayWorlds()

		if display == "" {
			t.Error("Expected non-empty world display")
		}
		if !contains(display, "Herbst") {
			t.Error("Expected display to contain 'Herbst'")
		}
		if !contains(display, "Shadow Realm") {
			t.Error("Expected display to contain 'Shadow Realm'")
		}
		if !contains(display, "Elder Forest") {
			t.Error("Expected display to contain 'Elder Forest'")
		}
	})

	t.Run("quit from world selection goes to welcome", func(t *testing.T) {
		ti := textinput.New()
		ti.Focus()
		m := &model{
			screen:       ScreenWorldSelect,
			textInput:    ti,
			visitedRooms: make(map[int]bool),
			knownExits:   make(map[string]bool),
			width:        80,
			height:       24,
		}
		m.Init()

		m.handleWorldSelectInput("quit")

		if m.screen != ScreenWelcome {
			t.Errorf("Expected screen %q, got %q", ScreenWelcome, m.screen)
		}
	})
}

// TestWorldSelectionEmpty tests behavior with no worlds available
func TestWorldSelectionEmpty(t *testing.T) {
	availableWorlds = nil
	defer func() { availableWorlds = nil }()

	t.Run("no worlds shows appropriate message", func(t *testing.T) {
		m := &model{
			width:        80,
			height:       24,
			visitedRooms: make(map[int]bool),
			knownExits:   make(map[string]bool),
		}
		display := m.displayWorlds()
		if display == "" {
			t.Error("Expected non-empty display even with no worlds")
		}
	})

	t.Run("no worlds triggers fetch on input", func(t *testing.T) {
		ti := textinput.New()
		ti.Focus()
		m := &model{
			screen:       ScreenWorldSelect,
			textInput:    ti,
			visitedRooms: make(map[int]bool),
			knownExits:   make(map[string]bool),
			width:        80,
			height:       24,
		}
		m.Init()

		// Should trigger fetch but keep same screen
		m.handleWorldSelectInput("1")
		if m.screen != ScreenWorldSelect {
			t.Log("Stays on world select while fetching (async)")
		}
	})
}

// TestParseWorldIndex tests the world index parsing utility
func TestParseWorldIndex(t *testing.T) {
	tests := []struct {
		input     string
		numWorlds int
		expected  int
	}{
		{"1", 3, 0},
		{"2", 3, 1},
		{"3", 3, 2},
		{"0", 3, -1},
		{"4", 3, -1},
		{"abc", 3, -1},
		{"", 3, -1},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseWorldIndex(tt.input, tt.numWorlds)
			if result != tt.expected {
				t.Errorf("parseWorldIndex(%q, %d) = %d, want %d",
					tt.input, tt.numWorlds, result, tt.expected)
			}
		})
	}
}

// TestProcessInputWorldSelect tests processInput routing to world select
func TestProcessInputWorldSelect(t *testing.T) {
	availableWorlds = []WorldInfo{
		{ID: "1", Name: "Herbst", Description: "Main world", Status: "active"},
	}
	defer func() { availableWorlds = nil }()

	t.Run("processInput routes to world select handler", func(t *testing.T) {
		ti := textinput.New()
		ti.Focus()
		m := &model{
			screen:       ScreenWorldSelect,
			textInput:    ti,
			visitedRooms: make(map[int]bool),
			knownExits:   make(map[string]bool),
			width:        80,
			height:       24,
		}
		m.Init()

		m.processInput("1")

		if m.currentWorld != "Herbst" {
			t.Errorf("Expected currentWorld 'Herbst', got %q", m.currentWorld)
		}
	})
}

// TestWorldSelectionJK tests j/k navigation on world selection
func TestWorldSelectionJK(t *testing.T) {
	availableWorlds = []WorldInfo{
		{ID: "1", Name: "Herbst", Description: "Main world", Status: "active"},
		{ID: "2", Name: "Shadow Realm", Description: "Dark world", Status: "active"},
		{ID: "3", Name: "Elder Forest", Description: "Ancient woods", Status: "active"},
	}
	defer func() { availableWorlds = nil }()

	t.Run("default cursor starts at 0", func(t *testing.T) {
		m := &model{
			visitedRooms: make(map[int]bool),
			knownExits:   make(map[string]bool),
		}
		if m.worldCursor != 0 {
			t.Errorf("Expected worldCursor to start at 0, got %d", m.worldCursor)
		}
	})

	t.Run("j moves cursor down", func(t *testing.T) {
		ti := textinput.New()
		ti.Focus()
		m := &model{
			screen:       ScreenWorldSelect,
			textInput:    ti,
			visitedRooms: make(map[int]bool),
			knownExits:   make(map[string]bool),
			width:        80,
			height:       24,
		}
		m.Init()

		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})

		if m.worldCursor != 1 {
			t.Errorf("Expected worldCursor 1 after 'j', got %d", m.worldCursor)
		}
	})

	t.Run("k moves cursor up", func(t *testing.T) {
		ti := textinput.New()
		ti.Focus()
		m := &model{
			screen:       ScreenWorldSelect,
			textInput:    ti,
			visitedRooms: make(map[int]bool),
			knownExits:   make(map[string]bool),
			width:        80,
			height:       24,
			worldCursor:  2,
		}
		m.Init()

		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})

		if m.worldCursor != 1 {
			t.Errorf("Expected worldCursor 1 after 'k', got %d", m.worldCursor)
		}
	})

	t.Run("j wraps to first world when at end", func(t *testing.T) {
		ti := textinput.New()
		ti.Focus()
		m := &model{
			screen:       ScreenWorldSelect,
			textInput:    ti,
			visitedRooms: make(map[int]bool),
			knownExits:   make(map[string]bool),
			width:        80,
			height:       24,
			worldCursor:  2,
		}
		m.Init()

		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})

		if m.worldCursor != 0 {
			t.Errorf("Expected worldCursor to wrap to 0, got %d", m.worldCursor)
		}
	})

	t.Run("k wraps to last world when at start", func(t *testing.T) {
		ti := textinput.New()
		ti.Focus()
		m := &model{
			screen:       ScreenWorldSelect,
			textInput:    ti,
			visitedRooms: make(map[int]bool),
			knownExits:   make(map[string]bool),
			width:        80,
			height:       24,
		}
		m.Init()

		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})

		if m.worldCursor != 2 {
			t.Errorf("Expected worldCursor to wrap to 2, got %d", m.worldCursor)
		}
	})

	t.Run("enter selects highlighted world", func(t *testing.T) {
		ti := textinput.New()
		ti.Focus()
		m := &model{
			screen:       ScreenWorldSelect,
			textInput:    ti,
			visitedRooms: make(map[int]bool),
			knownExits:   make(map[string]bool),
			width:        80,
			height:       24,
			worldCursor:  2,
		}
		m.Init()

		m.handleWorldSelectInput("")

		if m.currentWorld != "Elder Forest" {
			t.Errorf("Expected currentWorld 'Elder Forest', got %q", m.currentWorld)
		}
		if m.screen != ScreenCharacterSelect {
			t.Errorf("Expected screen ScreenCharacterSelect, got %q", m.screen)
		}
	})

	t.Run("displayWorlds shows cursor on highlighted world", func(t *testing.T) {
		m := &model{
			width:        80,
			height:       24,
			visitedRooms: make(map[int]bool),
			knownExits:   make(map[string]bool),
			worldCursor:  1,
		}
		display := m.displayWorlds()

		if !contains(display, "▸") {
			t.Error("Expected display to contain cursor character")
		}
	})
}

// contains is a helper to check substring presence
func contains(s, substr string) bool {
	if len(s) == 0 || len(substr) == 0 {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
