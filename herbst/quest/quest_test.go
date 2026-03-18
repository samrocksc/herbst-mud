package quest

import (
	"testing"
)

// TestGetExamineQuests tests retrieving examine-triggered quests
func TestGetExamineQuests(t *testing.T) {
	quests := GetExamineQuests()
	if len(quests) == 0 {
		t.Error("Expected at least one examine-triggered quest")
	}

	// Verify fountain quest exists
	found := false
	for _, q := range quests {
		if q.ID == "quest_fountain_secret" {
			found = true
			if q.ExamineTrigger == nil {
				t.Error("Fountain quest should have examine trigger")
			}
			if q.ExamineTrigger.MinExamineLevel != 75 {
				t.Errorf("Expected min examine level 75, got %d", q.ExamineTrigger.MinExamineLevel)
			}
		}
	}
	if !found {
		t.Error("Expected to find fountain secret quest")
	}
}

// TestGetQuestByTarget tests finding quests by target name
func TestGetQuestByTarget(t *testing.T) {
	quests := GetQuestByTarget("Stone Fountain")
	if len(quests) == 0 {
		t.Error("Expected to find quest for Stone Fountain")
	}

	// Test case-insensitive matching doesn't work yet (exact match required)
	quests = GetQuestByTarget("stone fountain")
	if len(quests) != 0 {
		t.Error("Expected no quest for lowercase stone fountain (case-sensitive)")
	}

	// Test non-existent target
	quests = GetQuestByTarget("NonExistentItem")
	if len(quests) != 0 {
		t.Error("Expected no quest for non-existent item")
	}
}

// TestGetQuestByID tests retrieving a quest by ID
func TestGetQuestByID(t *testing.T) {
	quest := GetQuestByID("quest_fountain_secret")
	if quest == nil {
		t.Fatal("Expected to find fountain secret quest")
	}

	if quest.Name != "The Fountain's Secret" {
		t.Errorf("Expected name 'The Fountain's Secret', got '%s'", quest.Name)
	}

	if quest.Type != QuestTypeSecret {
		t.Errorf("Expected type 'secret', got '%s'", quest.Type)
	}

	if !quest.Hidden {
		t.Error("Expected quest to be hidden")
	}
}

// TestCharacterQuestStore_UnlockQuest tests unlocking quests
func TestCharacterQuestStore_UnlockQuest(t *testing.T) {
	store := NewCharacterQuestStore()

	// Unlock a quest for character 1
	cq := store.UnlockQuest(1, "quest_fountain_secret")
	if cq == nil {
		t.Fatal("Expected quest to be unlocked")
	}

	if cq.Status != "unlocked" {
		t.Errorf("Expected status 'unlocked', got '%s'", cq.Status)
	}

	// Check if it's unlocked
	if !store.HasQuestUnlocked(1, "quest_fountain_secret") {
		t.Error("Expected quest to be unlocked")
	}

	// Try unlocking again (should return existing)
	cq2 := store.UnlockQuest(1, "quest_fountain_secret")
	if cq2.Status != "unlocked" {
		t.Error("Expected same quest to be returned on re-unlock")
	}
}

// TestCharacterQuestStore_HasQuestUnlocked tests checking if quest is unlocked
func TestCharacterQuestStore_HasQuestUnlocked(t *testing.T) {
	store := NewCharacterQuestStore()

	// Should return false for non-existent quest
	if store.HasQuestUnlocked(1, "quest_fountain_secret") {
		t.Error("Expected false for non-unlocked quest")
	}

	// Unlock the quest
	store.UnlockQuest(1, "quest_fountain_secret")

	// Should return true now
	if !store.HasQuestUnlocked(1, "quest_fountain_secret") {
		t.Error("Expected true for unlocked quest")
	}

	// Different character should not have it unlocked
	if store.HasQuestUnlocked(2, "quest_fountain_secret") {
		t.Error("Expected false for different character")
	}
}

// TestCheckExamineQuestUnlock tests the examine quest unlock logic
func TestCheckExamineQuestUnlock(t *testing.T) {
	// Use a fresh store for this test
	store := NewCharacterQuestStore()
	originalStore := GlobalQuestStore
	GlobalQuestStore = store
	defer func() { GlobalQuestStore = originalStore }()

	// Test with insufficient examine level
	result := CheckExamineQuestUnlock(1, "Stone Fountain", 50)
	if result != nil {
		t.Error("Expected nil result for insufficient examine level")
	}

	// Test with sufficient examine level
	result = CheckExamineQuestUnlock(1, "Stone Fountain", 75)
	if result == nil {
		t.Fatal("Expected unlock result")
	}

	if !result.Unlocked {
		t.Error("Expected quest to be unlocked")
	}

	if result.QuestID != "quest_fountain_secret" {
		t.Errorf("Expected quest ID 'quest_fountain_secret', got '%s'", result.QuestID)
	}

	if result.XPGained != 5 {
		t.Errorf("Expected XP gain 5, got %d", result.XPGained)
	}

	// Test already unlocked quest
	result = CheckExamineQuestUnlock(1, "Stone Fountain", 75)
	if result == nil {
		t.Fatal("Expected result for already unlocked quest")
	}

	if result.Unlocked {
		t.Error("Expected not unlocked for already unlocked quest")
	}

	if !result.AlreadyUnlocked {
		t.Error("Expected already_unlocked to be true")
	}

	// Test non-existent target
	result = CheckExamineQuestUnlock(1, "NonExistentItem", 100)
	if result != nil {
		t.Error("Expected nil for non-existent target")
	}
}

// TestGetUnlockedQuests tests retrieving unlocked quests
func TestGetUnlockedQuests(t *testing.T) {
	store := NewCharacterQuestStore()
	originalStore := GlobalQuestStore
	GlobalQuestStore = store
	defer func() { GlobalQuestStore = originalStore }()

	// No quests unlocked
	quests := GetUnlockedQuests(1)
	if len(quests) != 0 {
		t.Error("Expected no unlocked quests initially")
	}

	// Unlock a quest
	store.UnlockQuest(1, "quest_fountain_secret")

	quests = GetUnlockedQuests(1)
	if len(quests) != 1 {
		t.Errorf("Expected 1 unlocked quest, got %d", len(quests))
	}

	if quests[0].ID != "quest_fountain_secret" {
		t.Errorf("Expected quest ID 'quest_fountain_secret', got '%s'", quests[0].ID)
	}
}

// TestExamineTrigger_MinExamineLevel tests examine level threshold
func TestExamineTrigger_MinExamineLevel(t *testing.T) {
	store := NewCharacterQuestStore()
	originalStore := GlobalQuestStore
	GlobalQuestStore = store
	defer func() { GlobalQuestStore = originalStore }()

	quest := GetQuestByID("quest_fountain_secret")
	if quest == nil {
		t.Fatal("Quest not found")
	}

	// Test at exactly the threshold
	if quest.ExamineTrigger.MinExamineLevel != 75 {
		t.Errorf("Expected threshold 75, got %d", quest.ExamineTrigger.MinExamineLevel)
	}

	// Test below threshold
	result := CheckExamineQuestUnlock(1, "Stone Fountain", 74)
	if result != nil {
		t.Error("Expected no unlock below threshold")
	}

	// Test at threshold
	result = CheckExamineQuestUnlock(1, "Stone Fountain", 75)
	if result == nil || !result.Unlocked {
		t.Error("Expected unlock at threshold")
	}

	// Test above threshold
	store.UnlockQuest(2, "quest_fountain_secret") // Reset
	store.UnlockQuest(2, "quest_fountain_secret") // Already unlocked

	// Fresh character
	result = CheckExamineQuestUnlock(3, "Stone Fountain", 100)
	if result == nil || !result.Unlocked {
		t.Error("Expected unlock above threshold")
	}
}

// TestQuestRewards tests quest reward configuration
func TestQuestRewards(t *testing.T) {
	quest := GetQuestByID("quest_fountain_secret")
	if quest == nil {
		t.Fatal("Quest not found")
	}

	if quest.Rewards.XP != 50 {
		t.Errorf("Expected XP reward 50, got %d", quest.Rewards.XP)
	}
}