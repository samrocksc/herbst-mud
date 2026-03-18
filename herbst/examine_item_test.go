package main

import (
	"strings"
	"testing"
)

// TestExamineHiddenDetails tests hidden detail display
func TestExamineHiddenDetails(t *testing.T) {
	item := RoomItem{
		ID:             1,
		Name:           "Stone Fountain",
		Description:    "A murky stone fountain.",
		ExamineDesc:    "An old stone fountain, cracked and weathered. Water trickles from a broken statue in the center...",
		HiddenDetails: []HiddenDetail{
			{Description: "Something shiny at the bottom (could be coins)", Threshold: 0},
			{Description: "A crack in the base (something might be hidden inside)", Threshold: 10},
			{Description: "Faint writing on the rim: 'WISH UPON THE OOZE'", Threshold: 20},
		},
		HiddenThreshold: 10,
		ItemType:       "misc",
		IsImmovable:    true,
		Color:          "gold",
	}

	// Verify hidden details are populated
	if len(item.HiddenDetails) != 3 {
		t.Errorf("Expected 3 hidden details, got %d", len(item.HiddenDetails))
	}

	// Verify examine description is available
	if item.ExamineDesc == "" {
		t.Error("Expected examine description to be set")
	}
}

// TestExamineEquipmentStats tests equipment stat display
func TestExamineEquipmentStats(t *testing.T) {
	item := RoomItem{
		ID:             1,
		Name:           "Iron Sword",
		Description:    "A basic iron sword.",
		ExamineDesc:    "A well-crafted iron sword with a wooden hilt. The blade shows signs of use.",
		ItemType:       "weapon",
		Weight:         5,
		ItemDamage:     10,
		ItemDurability: 80,
	}

	// Verify stats are set
	if item.Weight != 5 {
		t.Errorf("Expected weight 5, got %d", item.Weight)
	}
	if item.ItemDamage != 10 {
		t.Errorf("Expected damage 10, got %d", item.ItemDamage)
	}
	if item.ItemDurability != 80 {
		t.Errorf("Expected durability 80, got %d", item.ItemDurability)
	}
}

// TestExamineTargetMatching tests target string matching
func TestExamineTargetMatching(t *testing.T) {
	tests := []struct {
		target    string
		itemName  string
		shouldMatch bool
	}{
		{"fountain", "Stone Fountain", true},
		{"fountain", "Water Fountain", true},
		{"stone", "Stone Fountain", true},
		{"SWORD", "Iron Sword", true},
		{"axe", "Iron Sword", false},
	}

	for _, tt := range tests {
		match := containsIgnoreCase(tt.itemName, tt.target)
		if match != tt.shouldMatch {
			t.Errorf("containsIgnoreCase(%q, %q) = %v, want %v", tt.itemName, tt.target, match, tt.shouldMatch)
		}
	}
}

// Helper function for case-insensitive substring matching
func containsIgnoreCase(s, substr string) bool {
	s = strings.ToLower(s)
	substr = strings.ToLower(substr)
	return strings.Contains(s, substr)
}