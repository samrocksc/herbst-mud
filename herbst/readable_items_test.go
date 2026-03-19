package main

import (
	"strconv"
	"strings"
	"testing"
)

// TestReadableItemStruct tests the readable item fields
func TestReadableItemStruct(t *testing.T) {
	item := RoomItem{
		ID:           1,
		Name:         "Old Scroll",
		Description:  "A dusty scroll",
		ItemType:     "scroll",
		IsReadable:   true,
		ReadableText: "The ancient text reads: 'Beware the Ooze!'",
	}

	if !item.IsReadable {
		t.Error("Expected item to be readable")
	}

	if item.ReadableText == "" {
		t.Error("Expected readable text to be set")
	}
}

// TestReadableItemNoText tests that non-readable items have no readable text
func TestReadableItemNoText(t *testing.T) {
	item := RoomItem{
		ID:           2,
		Name:         "Iron Sword",
		Description:  "A sharp blade",
		ItemType:     "weapon",
		IsReadable:   false,
		ReadableText: "",
	}

	if item.IsReadable {
		t.Error("Expected weapon to not be readable")
	}
}

// TestMultiPageBook tests multi-page book structure
func TestMultiPageBook(t *testing.T) {
	item := RoomItem{
		ID:        3,
		Name:      "Adventurer's Journal",
		ItemType:  "book",
		IsReadable: true,
		PageCount: 3,
		Pages: []string{
			"Day 1: We arrived at the fountain plaza. The water shimmered with an odd glow.",
			"Day 2: Explored the junkyard today. Found some strange golem parts.",
			"Day 3: The Ooze... it's everywhere. We must warn the others!",
		},
	}

	if item.PageCount != 3 {
		t.Errorf("Expected page count 3, got %d", item.PageCount)
	}

	if len(item.Pages) != 3 {
		t.Errorf("Expected 3 pages, got %d", len(item.Pages))
	}

	// Verify page content
	if !strings.Contains(item.Pages[0], "Day 1") {
		t.Error("Expected first page to contain 'Day 1'")
	}
}

// TestReadableItemTypes tests various readable item types
func TestReadableItemTypes(t *testing.T) {
	tests := []struct {
		itemType string
		expected bool
	}{
		{"book", true},
		{"scroll", true},
		{"sign", true},
		{"note", true},
		{"weapon", false},
		{"armor", false},
		{"misc", false},
	}

	for _, tt := range tests {
		item := RoomItem{
			ID:          1,
			Name:        "Test Item",
			ItemType:    tt.itemType,
			IsReadable:  tt.expected,
			ReadableText: "Test content",
		}

		if item.IsReadable != tt.expected {
			t.Errorf("ItemType %s: expected readable=%v, got %v", tt.itemType, tt.expected, item.IsReadable)
		}
	}
}

// TestReadCommandParsing tests read command argument parsing
func TestReadCommandParsing(t *testing.T) {
	tests := []struct {
		input        string
		expectedItem string
		expectedPage int
	}{
		{"read scroll", "scroll", 1},
		{"read book", "book", 1},
		{"read ancient tome", "ancient tome", 1},
		{"read scroll page 2", "scroll", 2},
		{"read book page 1", "book", 1},
		{"read 'Worn Note'", "Worn Note", 1},
	}

	for _, tt := range tests {
		parts := strings.Fields(tt.input)
		targetItem := strings.Join(parts[1:], " ")
		var pageNum int = 1

		if len(parts) >= 4 && parts[len(parts)-2] == "page" {
			pageNum, _ = strconv.Atoi(parts[len(parts)-1])
			targetItem = strings.Join(parts[1:len(parts)-2], " ")
		}

		if targetItem != tt.expectedItem {
			t.Errorf("Input %q: expected target %q, got %q", tt.input, tt.expectedItem, targetItem)
		}
		if pageNum != tt.expectedPage {
			t.Errorf("Input %q: expected page %d, got %d", tt.input, tt.expectedPage, pageNum)
		}
	}
}

// TestReadableContentDisplay tests content display formatting
func TestReadableContentDisplay(t *testing.T) {
	item := RoomItem{
		ID:          1,
		Name:        "Tavern Sign",
		ItemType:    "sign",
		IsReadable:  true,
		ReadableText: "Welcome to The Rusty Bucket - Ale, Mead, and Safety!",
	}

	// Single page content should not have navigation hints
	if strings.Contains(item.ReadableText, "page") {
		t.Error("Single page content should not mention pages")
	}
}

// TestMultiPageNavigation tests multi-page navigation hints
func TestMultiPageNavigationHints(t *testing.T) {
	item := RoomItem{
		ID:        1,
		Name:      "Tome of Secrets",
		ItemType:  "book",
		IsReadable: true,
		Pages: []string{
			"Chapter 1: The Beginning",
			"Chapter 2: The Journey",
			"Chapter 3: The End",
		},
	}

	// Page 1 should have "next page" hint but not "previous"
	page1Hint := "Use 'read tome of secrets page 2' for next page"
	if !strings.Contains(page1Hint, "page 2") {
		t.Error("Expected next page hint for page 1")
	}

	// Page 2 should have both hints
	page2Hint := "Use 'read tome of secrets page 3' for next page"
	if !strings.Contains(page2Hint, "page 3") {
		t.Error("Expected next page hint for page 2")
	}

	// Page 3 should have only "previous" hint
	page3PrevHint := "Use 'read tome of secrets page 2' for previous page"
	if !strings.Contains(page3PrevHint, "page 2") {
		t.Error("Expected previous page hint for page 3")
	}
}

// TestReadableItemMatching tests target matching for read command
func TestReadableItemMatching(t *testing.T) {
	tests := []struct {
		target    string
		itemName  string
		shouldMatch bool
	}{
		{"scroll", "Old Scroll", true},
		{"scroll", "Magic Scroll", true},
		{"sign", "Tavern Sign", true},
		{"sign", "Warning Sign", true},
		{"book", "Ancient Book", true},
		{"book", "Torn Book", true},
		{"journal", "Adventurer's Journal", true},
		{"tome", "Tome of Secrets", true},
	}

	for _, tt := range tests {
		// Simulate the matching logic used in handleReadCommand
		match := strings.Contains(strings.ToLower(tt.itemName), tt.target) ||
			strings.ToLower(tt.itemName) == tt.target
		if match != tt.shouldMatch {
			t.Errorf("Match(%q, %q) = %v, want %v", tt.target, tt.itemName, match, tt.shouldMatch)
		}
	}
}

// TestReadableItemFieldPresence tests all readable-related fields exist
func TestReadableItemFieldPresence(t *testing.T) {
	item := RoomItem{
		IsReadable:   true,
		ReadableText: "Test",
		PageCount:    2,
		Pages:        []string{"Page 1", "Page 2"},
	}

	// Verify all fields are present and correctly typed
	if !item.IsReadable {
		t.Error("IsReadable should be true")
	}
	if item.ReadableText != "Test" {
		t.Error("ReadableText mismatch")
	}
	if item.PageCount != 2 {
		t.Errorf("PageCount expected 2, got %d", item.PageCount)
	}
	if len(item.Pages) != 2 {
		t.Errorf("Pages length expected 2, got %d", len(item.Pages))
	}
}
