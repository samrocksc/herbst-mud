package main

import (
	"testing"
)

// TestIsContainer tests container identification
func TestIsContainer(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"rusty chest", true},
		{"wooden crate", true},
		{"old bag", true},
		{"dusty barrel", true},
		{"magic box", true},
		{"sack of gold", true},
		{"rusty sword", false},
		{"gold coin", false},
		{"potion", false},
		{"statue", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsContainer(tt.name)
			if result != tt.expected {
				t.Errorf("IsContainer(%q) = %v, want %v", tt.name, result, tt.expected)
			}
		})
	}
}

// TestContainerAddRemove tests adding and removing items
func TestContainerAddRemove(t *testing.T) {
	c := NewContainer("test-container", "Test Chest", "A test container")

	// Add items
	err := c.AddItem("rusty_sword")
	if err != nil {
		t.Errorf("Failed to add item: %v", err)
	}

	err = c.AddItem("gold_coin")
	if err != nil {
		t.Errorf("Failed to add item: %v", err)
	}

	// Check contents
	contents := c.GetContents()
	if len(contents) != 2 {
		t.Errorf("Expected 2 items, got %d", len(contents))
	}

	// Remove item
	removed := c.RemoveItem("rusty_sword")
	if !removed {
		t.Error("Failed to remove item")
	}

	contents = c.GetContents()
	if len(contents) != 1 {
		t.Errorf("Expected 1 item after removal, got %d", len(contents))
	}
}

// TestContainerCapacity tests capacity limits
func TestContainerCapacity(t *testing.T) {
	c := NewContainer("small-chest", "Small Chest", "A small chest")
	c.Capacity = 2

	// Should succeed
	err := c.AddItem("item1")
	if err != nil {
		t.Errorf("Should be able to add first item: %v", err)
	}

	err = c.AddItem("item2")
	if err != nil {
		t.Errorf("Should be able to add second item: %v", err)
	}

	// Should fail - container full
	err = c.AddItem("item3")
	if err == nil {
		t.Error("Expected error when container is full")
	}
}

// TestContainerOpenClose tests open/close functionality
func TestContainerOpenClose(t *testing.T) {
	c := NewContainer("locked-chest", "Locked Chest", "A locked chest")

	// Test close
	c.Close()
	if c.IsOpen {
		t.Error("Container should be closed after Close()")
	}

	// Test open
	c.IsLocked = false // Unlock first
	err := c.Open()
	if err != nil {
		t.Errorf("Open failed: %v", err)
	}
	if !c.IsOpen {
		t.Error("Container should be open after Open()")
	}

	// Test open when locked
	c.IsLocked = true
	err = c.Open()
	if err == nil {
		t.Error("Expected error when opening locked container")
	}
}

// TestContainerLockUnlock tests lock/unlock functionality
func TestContainerLockUnlock(t *testing.T) {
	c := NewContainer("chest", "Chest", "A chest")
	keyID := "golden_key"

	// Lock
	c.Lock(keyID)
	if !c.IsLocked {
		t.Error("Container should be locked")
	}
	if c.KeyID != keyID {
		t.Errorf("Expected keyID %s, got %s", keyID, c.KeyID)
	}

	// Unlock with wrong key
	err := c.Unlock("wrong_key")
	if err == nil {
		t.Error("Expected error with wrong key")
	}

	// Unlock with correct key
	err = c.Unlock(keyID)
	if err != nil {
		t.Errorf("Unlock failed: %v", err)
	}
	if c.IsLocked {
		t.Error("Container should be unlocked")
	}
}

// TestContainerFindItem tests finding items in container
func TestContainerFindItem(t *testing.T) {
	c := NewContainer("chest", "Chest", "A chest")
	c.Items = []string{"rusty_sword", "golden_amulet", "broken_dagger"}

	// Test exact match
	found := c.FindItemInContainer("rusty_sword")
	if found != "rusty_sword" {
		t.Errorf("Expected rusty_sword, got %s", found)
	}

	// Test partial match
	found = c.FindItemInContainer("sword")
	if found != "rusty_sword" {
		t.Errorf("Expected rusty_sword for 'sword', got %s", found)
	}

	// Test case insensitive
	found = c.FindItemInContainer("GOLDEN")
	if found != "golden_amulet" {
		t.Errorf("Expected golden_amulet for 'GOLDEN', got %s", found)
	}

	// Test not found
	found = c.FindItemInContainer("magic_wand")
	if found != "" {
		t.Error("Expected empty for not found")
	}
}

// TestContainerFormatContents tests formatting container contents
func TestContainerFormatContents(t *testing.T) {
	c := NewContainer("chest", "Dusty Crate", "A dusty wooden crate")
	c.Items = []string{"rusty_sword", "gold_coin"}

	itemDetails := map[string]string{
		"rusty_sword": "rusty sword (weapon)",
		"gold_coin":   "gold coin (currency)",
	}

	result := c.FormatContents(itemDetails)

	if result == "" {
		t.Error("Expected non-empty result")
	}

	// Check content
	if !containerContainsString(result, "Dusty Crate") {
		t.Error("Expected 'Dusty Crate' in output")
	}
	if !containerContainsString(result, "rusty sword") {
		t.Error("Expected 'rusty sword' in output")
	}
}

// TestContainerFormatClosed tests formatting when container is closed
func TestContainerFormatClosed(t *testing.T) {
	c := NewContainer("chest", "Iron Chest", "An iron chest")
	c.IsOpen = false

	result := c.FormatContents(map[string]string{})

	if !containerContainsString(result, "closed") {
		t.Error("Expected 'closed' in output for closed container")
	}
}

// TestContainerFormatEmpty tests formatting empty container
func TestContainerFormatEmpty(t *testing.T) {
	c := NewContainer("chest", "Empty Box", "An empty box")
	c.Items = []string{}

	result := c.FormatContents(map[string]string{})

	if !containerContainsString(result, "empty") {
		t.Error("Expected 'empty' in output for empty container")
	}
}

// TestContainerJSON tests JSON serialization
func TestContainerJSON(t *testing.T) {
	c := NewContainer("test", "Test", "Test desc")
	c.AddItem("item1")

	jsonStr, err := ContainerToJSON(c)
	if err != nil {
		t.Errorf("Failed to serialize: %v", err)
	}

	c2, err := ContainerFromJSON(jsonStr)
	if err != nil {
		t.Errorf("Failed to deserialize: %v", err)
	}

	if c2.Name != c.Name {
		t.Errorf("Expected name %s, got %s", c.Name, c2.Name)
	}
	if len(c2.Items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(c2.Items))
	}
}

// TestContainerFromEmptyJSON tests deserializing empty string
func TestContainerFromEmptyJSON(t *testing.T) {
	c, err := ContainerFromJSON("")
	if err != nil {
		t.Errorf("Failed to deserialize empty string: %v", err)
	}
	if c != nil {
		t.Error("Expected nil container for empty JSON")
	}
}

// containerContainsString is a test helper
func containerContainsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containerContainsSubstring(s, substr))
}

func containerContainsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}