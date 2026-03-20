package main

import (
	"testing"
)

// Helper to check if string contains substring
func strContains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > 0 && 
		(s[:len(substr)] == substr || strContains(s[1:], substr))))
}

// TestReadableCommandBasic tests basic read functionality
func TestReadableCommandBasic(t *testing.T) {
	m := &model{
		roomItems: []RoomItem{
			{
				ID:          1,
				Name:        "survival_manual",
				Description: "A worn survival guide",
				IsReadable:  true,
				Content:     "SURVIVAL IN THE NEW WORLD\nA Guide for the Post-Ooze Era\n\nContents:\n  1. Finding Water\n  2. Identifying Safe Food\n  3. Avoiding Ooze Pools...",
			},
		},
	}

	// Test reading a readable item
	m.handleReadCommand("read survival_manual")
	if m.messageType != "info" {
		t.Errorf("Expected info message type, got %s", m.messageType)
	}
	if m.message == "" {
		t.Error("Expected message content, got empty")
	}
}

// TestReadCommandNonReadable tests reading a non-readable item
func TestReadCommandNonReadable(t *testing.T) {
	m := &model{
		roomItems: []RoomItem{
			{
				ID:          1,
				Name:        "rusty_sword",
				Description: "A rusty iron sword",
				IsReadable:  false,
			},
		},
	}

	m.handleReadCommand("read rusty_sword")
	if m.messageType != "error" {
		t.Errorf("Expected error message type, got %s", m.messageType)
	}
	if !strContains(m.message, "can't read") {
		t.Errorf("Expected 'can't read' message, got: %s", m.message)
	}
}

// TestReadCommandNotFound tests reading a non-existent item
func TestReadCommandNotFound(t *testing.T) {
	m := &model{
		roomItems: []RoomItem{},
	}

	m.handleReadCommand("read nonexistent")
	if m.messageType != "error" {
		t.Errorf("Expected error message type, got %s", m.messageType)
	}
	if !strContains(m.message, "don't see") {
		t.Errorf("Expected 'don't see' message, got: %s", m.message)
	}
}

// TestReadCommandNoArgs tests read command with no arguments
func TestReadCommandNoArgs(t *testing.T) {
	m := &model{}

	m.handleReadCommand("read")
	if m.messageType != "error" {
		t.Errorf("Expected error message type, got %s", m.messageType)
	}
	if !strContains(m.message, "Read what") {
		t.Errorf("Expected 'Read what' message, got: %s", m.message)
	}
}

// TestReadCommandSkillGated tests reading an item with skill requirement (insufficient skill)
func TestReadCommandSkillGated(t *testing.T) {
	m := &model{
		roomItems: []RoomItem{
			{
				ID:             1,
				Name:           "encrypted_terminal",
				Description:    "A high-tech terminal",
				IsReadable:     true,
				Content:        "[Encrypted - requires tech skill level 5 to decode]",
				ReadSkill:      "tech",
				ReadSkillLevel: 5,
			},
		},
		currentCharacterID: 1,
		// Note: getCharacterSkillLevel would return 0 (no skill) in this mock
	}

	m.handleReadCommand("read encrypted_terminal")
	if m.messageType != "info" {
		t.Errorf("Expected info message type, got %s", m.messageType)
	}
	if !strContains(m.message, "Requires tech skill level 5") {
		t.Errorf("Expected skill requirement message, got: %s", m.message)
	}
}

// TestReadCommandPartialMatch tests that partial name matching works
func TestReadCommandPartialMatch(t *testing.T) {
	m := &model{
		roomItems: []RoomItem{
			{
				ID:       1,
				Name:     "ancient_tome",
				IsReadable: true,
				Content:  "The secrets of the universe...",
			},
		},
	}

	m.handleReadCommand("read ancient")
	if m.messageType != "info" {
		t.Errorf("Expected info message type, got %s", m.messageType)
	}
}

// TestRoomItemReadableFields tests that RoomItem struct has readable fields
func TestRoomItemReadableFields(t *testing.T) {
	item := RoomItem{
		ID:             1,
		Name:           "test_book",
		IsReadable:     true,
		Content:        "Test content",
		ReadSkill:      "tech",
		ReadSkillLevel: 3,
	}

	if !item.IsReadable {
		t.Error("Expected IsReadable to be true")
	}
	if item.Content != "Test content" {
		t.Errorf("Expected Content 'Test content', got '%s'", item.Content)
	}
	if item.ReadSkill != "tech" {
		t.Errorf("Expected ReadSkill 'tech', got '%s'", item.ReadSkill)
	}
	if item.ReadSkillLevel != 3 {
		t.Errorf("Expected ReadSkillLevel 3, got %d", item.ReadSkillLevel)
	}
}