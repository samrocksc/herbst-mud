package db_test

import (
	"context"
	"testing"

	"herbst-server/db"

	_ "github.com/mattn/go-sqlite3"
	"entgo.io/ent/dialect"
)

// TestAvailableTalentSchema tests the AvailableTalent entity schema
func TestAvailableTalentSchema(t *testing.T) {
	client, err := db.Open(dialect.SQLite, "file:ent?mode=memory&_fk=1")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}

	// Create a room first (required for character)
	room, err := client.Room.Create().
		SetName("Test Room").
		SetDescription("A test room").
		SetIsStartingRoom(true).
		SetExits(map[string]int{}).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create room: %v", err)
	}

	// Create a character
	char, err := client.Character.Create().
		SetName("TestChar").
		SetCurrentRoomId(room.ID).
		SetStartingRoomId(room.ID).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create character: %v", err)
	}

	// Create a talent
	tal, err := client.Talent.Create().
		SetName("Power Strike").
		SetDescription("A powerful strike").
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create talent: %v", err)
	}

	// Create an available talent (unlocked for character)
	availTal, err := client.AvailableTalent.Create().
		SetCharacterID(char.ID).
		SetTalentID(tal.ID).
		SetUnlockReason("level_up").
		SetUnlockedAtLevel(5).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create available talent: %v", err)
	}

	if availTal.UnlockReason != "level_up" {
		t.Errorf("expected unlock_reason 'level_up', got %s", availTal.UnlockReason)
	}

	if availTal.UnlockedAtLevel != 5 {
		t.Errorf("expected unlocked_at_level 5, got %d", availTal.UnlockedAtLevel)
	}

	// Query available talents for character
	charAvailTalents, err := char.QueryAvailableTalents().All(ctx)
	if err != nil {
		t.Fatalf("failed to query available talents: %v", err)
	}
	if len(charAvailTalents) != 1 {
		t.Errorf("expected 1 available talent, got %d", len(charAvailTalents))
	}
}

// TestAvailableTalentCRUD tests basic CRUD for AvailableTalent
func TestAvailableTalentCRUD(t *testing.T) {
	client, err := db.Open(dialect.SQLite, "file:ent?mode=memory&_fk=1")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}

	// Create a room first (required for character)
	room, err := client.Room.Create().
		SetName("CRUD Test Room").
		SetDescription("A test room").
		SetIsStartingRoom(true).
		SetExits(map[string]int{}).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create room: %v", err)
	}

	// Create a character
	char, err := client.Character.Create().
		SetName("CRUDChar").
		SetCurrentRoomId(room.ID).
		SetStartingRoomId(room.ID).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create character: %v", err)
	}

	// Create multiple talents and make them available
	talentNames := []string{"Double Strike", "Fast Healing", "Sprint"}
	for _, name := range talentNames {
		tal, err := client.Talent.Create().
			SetName(name).
			SetDescription(name + " ability").
			Save(ctx)
		if err != nil {
			t.Fatalf("failed to create talent %s: %v", name, err)
		}

		_, err = client.AvailableTalent.Create().
			SetCharacterID(char.ID).
			SetTalentID(tal.ID).
			SetUnlockReason("quest").
			SetUnlockedAtLevel(3).
			Save(ctx)
		if err != nil {
			t.Fatalf("failed to create available talent %s: %v", name, err)
		}
	}

	// Count available talents
	count, err := client.AvailableTalent.Query().Count(ctx)
	if err != nil {
		t.Fatalf("failed to count available talents: %v", err)
	}
	if count != 3 {
		t.Errorf("expected 3 available talents, got %d", count)
	}

	// Update one available talent - get the first one from the character
	charAvailTalents, err := char.QueryAvailableTalents().All(ctx)
	if err != nil {
		t.Fatalf("failed to query character available talents: %v", err)
	}
	if len(charAvailTalents) == 0 {
		t.Fatal("expected available talents")
	}

	availTal := charAvailTalents[0]
	updated, err := availTal.Update().
		SetUnlockReason("skill_trainer").
		SetUnlockedAtLevel(10).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to update available talent: %v", err)
	}

	if updated.UnlockReason != "skill_trainer" {
		t.Errorf("expected unlock_reason 'skill_trainer', got %s", updated.UnlockReason)
	}
	if updated.UnlockedAtLevel != 10 {
		t.Errorf("expected unlocked_at_level 10, got %d", updated.UnlockedAtLevel)
	}

	// Verify query works - should have 3 available talents
	count, err = client.AvailableTalent.Query().Count(ctx)
	if err != nil {
		t.Fatalf("failed to count available talents: %v", err)
	}
	if count != 3 {
		t.Errorf("expected 3 available talents, got %d", count)
	}
}