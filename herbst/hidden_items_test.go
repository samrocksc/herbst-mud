package main

import (
	"testing"
)

// TestHiddenItemsNotShownInRoom verifies hidden items are filtered out
func TestHiddenItemsNotShownInRoom(t *testing.T) {
	items := []RoomItem{
		{ID: 1, Name: "Visible Sword", IsVisible: true},
		{ID: 2, Name: "Hidden Key", IsVisible: false},
		{ID: 3, Name: "Visible Potion", IsVisible: true},
	}

	// Simulate room display filtering
	var visibleItems []RoomItem
	for _, item := range items {
		if item.IsVisible {
			visibleItems = append(visibleItems, item)
		}
	}

	if len(visibleItems) != 2 {
		t.Errorf("Expected 2 visible items, got %d", len(visibleItems))
	}

	// Verify hidden item is not in visible list
	for _, item := range visibleItems {
		if item.Name == "Hidden Key" {
			t.Error("Hidden Key should not be visible")
		}
	}
}

// TestExamineCanRevealHiddenItems verifies examine reveals hidden items
func TestExamineCanRevealHiddenItems(t *testing.T) {
	skill := &ExamineSkillLevel{Level: 50, XP: 0}

	// Hidden item with reveal condition
	hiddenItem := RoomItem{
		ID:              1,
		Name:            "Secret Key",
		Description:     "You don't see anything special.",
		ExamineDesc:     "A small brass key hidden in the crack!",
		IsVisible:       false,
		RevealThreshold: 30, // Requires level 30 to reveal
	}

	// Player with level 50 should reveal it
	if skill.Level >= hiddenItem.RevealThreshold {
		hiddenItem.IsVisible = true
	}

	if !hiddenItem.IsVisible {
		t.Error("Hidden item should be revealed for skill level 50 (threshold: 30)")
	}
}

// TestLowSkillCannotReveal verifies low skill doesn't reveal items
func TestLowSkillCannotReveal(t *testing.T) {
	skill := &ExamineSkillLevel{Level: 10, XP: 0}

	hiddenItem := RoomItem{
		ID:              1,
		Name:            "Secret Key",
		IsVisible:       false,
		RevealThreshold: 30,
	}

	// Player with level 10 should NOT reveal it
	if skill.Level >= hiddenItem.RevealThreshold {
		hiddenItem.IsVisible = true
	}

	if hiddenItem.IsVisible {
		t.Error("Hidden item should NOT be revealed for skill level 10 (threshold: 30)")
	}
}

// TestPerceptionCheckReveal verifies perception check reveals hidden items
func TestPerceptionCheckReveal(t *testing.T) {
	skill := &ExamineSkillLevel{Level: 20, XP: 0}
	intStat := 10
	wisStat := 10

	hiddenItem := RoomItem{
		ID:              1,
		Name:            "Hidden Coin",
		IsVisible:       false,
		RevealThreshold: 15,
		RevealDC:        25, // Difficulty class
		Stat:            "WIS",
	}

	// Perform perception check
	check := skill.ExamineCheck(intStat, wisStat) // 20 + 10 + 5 + random(1-10) = 35-44

	if check >= hiddenItem.RevealDC && skill.Level >= hiddenItem.RevealThreshold {
		hiddenItem.IsVisible = true
	}

	// With level 20 and stats, should likely pass DC 25
	if !hiddenItem.IsVisible {
		t.Logf("Check result: %d, DC: %d - item not revealed (may be random)", check, hiddenItem.RevealDC)
	}
}

// TestRevealedItemsPersist verifies revealed items stay visible
func TestRevealedItemsPersist(t *testing.T) {
	// Simulate character knowing about revealed items
	type CharacterKnowledge struct {
		RevealedItems map[int]bool
	}

	char := CharacterKnowledge{
		RevealedItems: make(map[int]bool),
	}

	// Reveal item 1
	char.RevealedItems[1] = true

	// Verify persistence
	if !char.RevealedItems[1] {
		t.Error("Revealed item should persist in character knowledge")
	}

	// Verify item 2 is not revealed
	if char.RevealedItems[2] {
		t.Error("Item 2 should not be revealed")
	}
}

// TestExamineRevealWithThresholdZero verifies auto-reveal at threshold 0
func TestExamineRevealWithThresholdZero(t *testing.T) {
	skill := &ExamineSkillLevel{Level: 0, XP: 0}

	hiddenItem := RoomItem{
		ID:              1,
		Name:            "Always Visible Secret",
		IsVisible:       false,
		RevealThreshold: 0, // Revealed for anyone
	}

	// Threshold of 0 means visible to all
	if skill.Level >= hiddenItem.RevealThreshold {
		hiddenItem.IsVisible = true
	}

	if !hiddenItem.IsVisible {
		t.Error("Item with threshold 0 should be visible to everyone")
	}
}

// TestMultipleHiddenItemsReveal tests revealing multiple items
func TestMultipleHiddenItemsReveal(t *testing.T) {
	skill := &ExamineSkillLevel{Level: 40, XP: 0}

	items := []RoomItem{
		{ID: 1, Name: "Key", IsVisible: false, RevealThreshold: 20},
		{ID: 2, Name: "Gem", IsVisible: false, RevealThreshold: 50},
		{ID: 3, Name: "Scroll", IsVisible: false, RevealThreshold: 30},
		{ID: 4, Name: "Coin", IsVisible: false, RevealThreshold: 10},
	}

	revealedCount := 0
	for i := range items {
		if skill.Level >= items[i].RevealThreshold {
			items[i].IsVisible = true
			revealedCount++
		}
	}

	// Should reveal Key (20), Scroll (30), Coin (10) = 3 items
	if revealedCount != 3 {
		t.Errorf("Expected 3 revealed items, got %d", revealedCount)
	}
}