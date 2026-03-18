package main

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

// TestRoomItemDisplay tests the roomItemDisplay struct
func TestRoomItemDisplay(t *testing.T) {
	// Test creating an immovable item (gold color)
	immovableItem := roomItemDisplay{
		name:        "Ancient Statue",
		description: "A weathered statue of a forgotten hero",
		color:       lipgloss.Color("220"), // Gold
		isImmovable: true,
		isVisible:   true,
	}

	if !immovableItem.isImmovable {
		t.Error("Expected immovable item to be marked as immovable")
	}

	if immovableItem.color != "220" {
		t.Errorf("Expected gold color 220, got %v", immovableItem.color)
	}

	// Test creating a visible item with custom color
	customColorItem := roomItemDisplay{
		name:        "Magic Potion",
		description: "A glowing red potion",
		color:       lipgloss.Color("196"), // Red
		isImmovable: false,
		isVisible:   true,
	}

	if customColorItem.isImmovable {
		t.Error("Expected movable item to not be immovable")
	}

	if customColorItem.color != "196" {
		t.Errorf("Expected red color 196, got %v", customColorItem.color)
	}

	// Test creating an invisible item
	invisibleItem := roomItemDisplay{
		name:        "Hidden Treasure",
		description: "A secret stash",
		color:       lipgloss.Color("15"),
		isImmovable: false,
		isVisible:   false,
	}

	if invisibleItem.isVisible {
		t.Error("Expected invisible item to not be visible")
	}
}

// TestFormatRoomItemsWithColor tests the formatRoomItemsWithColor function
func TestFormatRoomItemsWithColor(t *testing.T) {
	m := &model{
		roomItems: []roomItemDisplay{
			{
				name:        "Golden Idol",
				description: "An ancient golden idol",
				color:       lipgloss.Color("220"),
				isImmovable: true,
				isVisible:   true,
			},
			{
				name:        "Magic Sword",
				description: "A blade humming with energy",
				color:       lipgloss.Color("75"), // Blue
				isImmovable: false,
				isVisible:   true,
			},
		},
	}

	result := m.formatRoomItemsWithColor()

	if result == "" {
		t.Error("Expected non-empty result when items exist")
	}

	// Check that items are included
	if !contains(result, "Golden Idol") {
		t.Error("Expected result to contain 'Golden Idol'")
	}

	if !contains(result, "Magic Sword") {
		t.Error("Expected result to contain 'Magic Sword'")
	}

	// Check immovable indicator
	if !contains(result, "(fixed)") {
		t.Error("Expected result to contain '(fixed)' for immovable items")
	}
}

// TestFormatRoomItemsWithColorEmpty tests empty items
func TestFormatRoomItemsWithColorEmpty(t *testing.T) {
	m := &model{
		roomItems: []roomItemDisplay{},
	}

	result := m.formatRoomItemsWithColor()

	if result != "" {
		t.Errorf("Expected empty result for empty items, got: %s", result)
	}
}

// TestFormatRoomItemsWithColorInvisible tests invisible items are not shown
func TestFormatRoomItemsWithColorInvisible(t *testing.T) {
	m := &model{
		roomItems: []roomItemDisplay{
			{
				name:        "Hidden Item",
				description: "Should not be shown",
				color:       lipgloss.Color("15"),
				isImmovable: false,
				isVisible:   false,
			},
		},
	}

	result := m.formatRoomItemsWithColor()

	if result != "" {
		t.Errorf("Expected empty result for invisible items, got: %s", result)
	}
}

// TestFormatRoomItemsWithColorDefaults tests default color when no custom color
func TestFormatRoomItemsWithColorDefaults(t *testing.T) {
	m := &model{
		roomItems: []roomItemDisplay{
			{
				name:        "Plain Item",
				description: "A plain item with default color",
				color:       lipgloss.Color("15"), // White default
				isImmovable: false,
				isVisible:   true,
			},
		},
	}

	result := m.formatRoomItemsWithColor()

	if !contains(result, "Plain Item") {
		t.Error("Expected result to contain 'Plain Item'")
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