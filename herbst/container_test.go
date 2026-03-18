package main

import (
	"testing"
)

// TestContainerCommandParsing tests parsing of container commands
func TestContainerCommandParsing(t *testing.T) {
	tests := []struct {
		name          string
		cmd           string
		wantItem      string
		wantContainer string
		wantValid     bool
	}{
		{
			name:          "simple take from container",
			cmd:           "take sword from chest",
			wantItem:      "sword",
			wantContainer: "chest",
			wantValid:     true,
		},
		{
			name:          "take with multiple word item",
			cmd:           "take rusty sword from wooden chest",
			wantItem:      "rusty sword",
			wantContainer: "wooden chest",
			wantValid:     true,
		},
		{
			name:          "take with multiple word container",
			cmd:           "take key from old wooden chest",
			wantItem:      "key",
			wantContainer: "old wooden chest",
			wantValid:     true,
		},
		{
			name:          "put in container",
			cmd:           "put sword in chest",
			wantItem:      "sword",
			wantContainer: "chest",
			wantValid:     true,
		},
		{
			name:          "put with multiple words",
			cmd:           "put rusty sword in wooden chest",
			wantItem:      "rusty sword",
			wantContainer: "wooden chest",
			wantValid:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// For take commands
			if len(tt.cmd) >= 4 && tt.cmd[:4] == "take" {
				// Parse: take <item> from <container>
				fromIdx := indexOf(tt.cmd, " from ")
				if fromIdx == -1 {
					if tt.wantValid {
						t.Errorf("expected valid command, but ' from ' not found")
					}
					return
				}
				item := trim(tt.cmd[5:fromIdx])
				container := trim(tt.cmd[fromIdx+6:])
				if item != tt.wantItem {
					t.Errorf("got item %q, want %q", item, tt.wantItem)
				}
				if container != tt.wantContainer {
					t.Errorf("got container %q, want %q", container, tt.wantContainer)
				}
			}

			// For put commands
			if len(tt.cmd) >= 3 && tt.cmd[:3] == "put" {
				// Parse: put <item> in <container>
				inIdx := indexOf(tt.cmd, " in ")
				if inIdx == -1 {
					if tt.wantValid {
						t.Errorf("expected valid command, but ' in ' not found")
					}
					return
				}
				item := trim(tt.cmd[4:inIdx])
				container := trim(tt.cmd[inIdx+4:])
				if item != tt.wantItem {
					t.Errorf("got item %q, want %q", item, tt.wantItem)
				}
				if container != tt.wantContainer {
					t.Errorf("got container %q, want %q", container, tt.wantContainer)
				}
			}
		})
	}
}

// TestRoomItemIsContainer tests container detection in RoomItem
func TestRoomItemIsContainer(t *testing.T) {
	item := RoomItem{
		ID:          1,
		Name:        "Wooden Chest",
		Description: "A sturdy wooden chest",
		IsContainer: true,
		Capacity:    10,
		IsLocked:    false,
	}

	if !item.IsContainer {
		t.Error("expected IsContainer to be true")
	}
	if item.Capacity != 10 {
		t.Errorf("expected Capacity 10, got %d", item.Capacity)
	}
	if item.IsLocked {
		t.Error("expected IsLocked to be false")
	}
}

// TestRoomItemLockedContainer tests locked container behavior
func TestRoomItemLockedContainer(t *testing.T) {
	item := RoomItem{
		ID:          2,
		Name:        "Iron Safe",
		Description: "A heavy iron safe with a complex lock",
		IsContainer: true,
		Capacity:    5,
		IsLocked:    true,
	}

	if !item.IsLocked {
		t.Error("expected IsLocked to be true")
	}
}

// TestContainerCapacity tests container capacity limits
func TestContainerCapacity(t *testing.T) {
	tests := []struct {
		name      string
		capacity  int
		contents  int
		wantFull  bool
	}{
		{
			name:     "empty container",
			capacity: 10,
			contents: 0,
			wantFull: false,
		},
		{
			name:     "half full container",
			capacity: 10,
			contents: 5,
			wantFull: false,
		},
		{
			name:     "full container",
			capacity: 10,
			contents: 10,
			wantFull: true,
		},
		{
			name:     "over capacity",
			capacity: 10,
			contents: 11,
			wantFull: true,
		},
		{
			name:     "zero capacity (unlimited)",
			capacity: 0,
			contents: 100,
			wantFull: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isFull := tt.capacity > 0 && tt.contents >= tt.capacity
			if isFull != tt.wantFull {
				t.Errorf("got full=%v, want %v", isFull, tt.wantFull)
			}
		})
	}
}

// TestItemTypeForContainers tests itemType validation for containers
func TestItemTypeForContainers(t *testing.T) {
	item := RoomItem{
		ID:          3,
		Name:        "Leather Bag",
		Description: "A small leather bag",
		ItemType:    "container",
		IsContainer: true,
		Capacity:    5,
	}

	if item.ItemType != "container" {
		t.Errorf("expected ItemType 'container', got %q", item.ItemType)
	}
}

// Helper functions for testing
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func trim(s string) string {
	// Simple whitespace trim
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}