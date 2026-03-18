package main

import (
	"context"
	"testing"

	"herbst-server/db"
	"herbst-server/db/available_talent"
	"herbst-server/db/character"
	"herbst-server/db/skill"
	"herbst-server/db/talent"
)

// TestAvailableTalentSchema tests the AvailableTalent entity schema
func TestAvailableTalentSchema(t *testing.T) {
	client, err := db.Open("sqlite3", "file:ent?mode=memory&_fk=1")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Create a character first
	char, err := client.Character.Create().
		SetName("TestChar").
		SetCurrentRoomId(1).
		SetStartingRoomId(1).
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

	// Create an available talent (unlocked but not equipped)
	at, err := client.AvailableTalent.Create().
		SetCharacter(char).
		SetTalent(tal).
		SetUnlockReason("level_up").
		SetUnlockedAtLevel(5).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create available talent: %v", err)
	}

	if at.UnlockReason != "level_up" {
		t.Errorf("expected unlock_reason level_up, got %s", at.UnlockReason)
	}

	if at.UnlockedAtLevel != 5 {
		t.Errorf("expected unlocked_at_level 5, got %d", at.UnlockedAtLevel)
	}

	// Query available talents for character
	availTalents, err := char.QueryAvailableTalents().All(ctx)
	if err != nil {
		t.Fatalf("failed to query available talents: %v", err)
	}

	if len(availTalents) != 1 {
		t.Errorf("expected 1 available talent, got %d", len(availTalents))
	}
}

// TestCharacterSkillTalentRelation tests the full chain: Character -> Skills/Talents/AvailableTalents
func TestCharacterSkillTalentRelation(t *testing.T) {
	client, err := db.Open("sqlite3", "file:ent?mode=memory&_fk=1")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Create character
	char, err := client.Character.Create().
		SetName("Hero").
		SetCurrentRoomId(1).
		SetStartingRoomId(1).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create character: %v", err)
	}

	// Create skill
	sk, err := client.Skill.Create().
		SetName("Swordsmanship").
		SetDescription("Mastery of bladed weapons").
		SetSkillType("combat").
		SetCost(10).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create skill: %v", err)
	}

	// Link skill to character
	_, err = client.CharacterSkill.Create().
		SetCharacter(char).
		SetSkill(sk).
		SetLevel(5).
		SetExperience(100).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to link skill to character: %v", err)
	}

	// Create talent
	tal, err := client.Talent.Create().
		SetName("Cleave").
		SetDescription("Hit multiple enemies").
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create talent: %v", err)
	}

	// Link talent to character (equipped)
	_, err = client.CharacterTalent.Create().
		SetCharacter(char).
		SetTalent(tal).
		SetSlot(0).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to link talent to character: %v", err)
	}

	// Create another talent as available (unlocked but not equipped)
	tal2, err := client.Talent.Create().
		SetName("Shield Block").
		SetDescription("Block with shield").
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create talent: %v", err)
	}

	_, err = client.AvailableTalent.Create().
		SetCharacter(char).
		SetTalent(tal2).
		SetUnlockReason("quest").
		SetUnlockedAtLevel(3).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create available talent: %v", err)
	}

	// Verify: Query character's skills
	charSkills, err := char.QuerySkills().All(ctx)
	if err != nil {
		t.Fatalf("failed to query character skills: %v", err)
	}
	if len(charSkills) != 1 {
		t.Errorf("expected 1 character skill, got %d", len(charSkills))
	}

	// Verify: Query character's equipped talents
	charTalents, err := char.QueryTalents().All(ctx)
	if err != nil {
		t.Fatalf("failed to query character talents: %v", err)
	}
	if len(charTalents) != 1 {
		t.Errorf("expected 1 equipped talent, got %d", len(charTalents))
	}

	// Verify: Query character's available talents
	availableTalents, err := char.QueryAvailableTalents().All(ctx)
	if err != nil {
		t.Fatalf("failed to query available talents: %v", err)
	}
	if len(availableTalents) != 1 {
		t.Errorf("expected 1 available talent, got %d", len(availableTalents))
	}

	// Verify talent points to available_talents edge
	availChars, err := tal2.QueryAvailableToCharacters().All(ctx)
	if err != nil {
		t.Fatalf("failed to query available to characters: %v", err)
	}
	if len(availChars) != 1 {
		t.Errorf("expected 1 character for available talent, got %d", len(availChars))
	}
}

// TestSkillTalentEndpoints tests API-style access patterns
func TestSkillTalentEndpoints(t *testing.T) {
	client, err := db.Open("sqlite3", "file:ent?mode=memory&_fk=1")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Seed skills
	skills := []struct {
		name      string
		skillType string
		cost      int
	}{
		{"Blades", "combat", 10},
		{"Staves", "magic", 10},
		{"Tech", "utility", 8},
	}

	for _, s := range skills {
		_, err := client.Skill.Create().
			SetName(s.name).
			SetDescription(s.name + " skill").
			SetSkillType(s.skillType).
			SetCost(s.cost).
			Save(ctx)
		if err != nil {
			t.Fatalf("failed to create skill %s: %v", s.name, err)
		}
	}

	// Seed talents
	talents := []struct {
		name  string
		desc  string
	}{
		{"Double Strike", "Attack twice"},
		{"Healing Light", "Restore health"},
		{"Quick Step", "Move faster"},
	}

	for _, t := range talents {
		_, err := client.Talent.Create().
			SetName(t.name).
			SetDescription(t.desc).
			Save(ctx)
		if err != nil {
			t.Fatalf("failed to create talent %s: %v", t.name, err)
		}
	}

	// Count skills
	skillCount, err := client.Skill.Query().Count(ctx)
	if err != nil {
		t.Fatalf("failed to count skills: %v", err)
	}
	if skillCount != 3 {
		t.Errorf("expected 3 skills, got %d", skillCount)
	}

	// Count talents
	talentCount, err := client.Talent.Query().Count(ctx)
	if err != nil {
		t.Fatalf("failed to count talents: %v", err)
	}
	if talentCount != 3 {
		t.Errorf("expected 3 talents, got %d", talentCount)
	}

	// Query by type
	combatSkills, err := client.Skill.Query().Where(skill.SkillType("combat")).All(ctx)
	if err != nil {
		t.Fatalf("failed to query combat skills: %v", err)
	}
	if len(combatSkills) != 1 {
		t.Errorf("expected 1 combat skill, got %d", len(combatSkills))
	}

	// Query by name
	bladeSkill, err := client.Skill.Query().Where(skill.Name("Blades")).Only(ctx)
	if err != nil {
		t.Fatalf("failed to get blades skill: %v", err)
	}
	if bladeSkill.Cost != 10 {
		t.Errorf("expected cost 10, got %d", bladeSkill.Cost)
	}

	// Talent by name
	doubleStrike, err := client.Talent.Query().Where(talent.Name("Double Strike")).Only(ctx)
	if err != nil {
		t.Fatalf("failed to get double strike talent: %v", err)
	}
	if doubleStrike.Description != "Attack twice" {
		t.Errorf("expected description 'Attack twice', got %s", doubleStrike.Description)
	}
}

// TestCharacterSkillUpdate tests updating character skill level
func TestCharacterSkillUpdate(t *testing.T) {
	client, err := db.Open("sqlite3", "file:ent?mode=memory&_fk=1")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Create character and skill
	char, err := client.Character.Create().
		SetName("UpdateTest").
		SetCurrentRoomId(1).
		SetStartingRoomId(1).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create character: %v", err)
	}

	sk, err := client.Skill.Create().
		SetName("UpdateSkill").
		SetCost(5).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create skill: %v", err)
	}

	charSkill, err := client.CharacterSkill.Create().
		SetCharacter(char).
		SetSkill(sk).
		SetLevel(1).
		SetExperience(0).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create character skill: %v", err)
	}

	// Update skill level and experience
	updated, err := charSkill.Update().
		SetLevel(10).
		SetExperience(500).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to update character skill: %v", err)
	}

	if updated.Level != 10 {
		t.Errorf("expected level 10, got %d", updated.Level)
	}

	if updated.Experience != 500 {
		t.Errorf("expected experience 500, got %d", updated.Experience)
	}
}

// TestCharacterTalentSlot tests talent slot assignment
func TestCharacterTalentSlot(t *testing.T) {
	client, err := db.Open("sqlite3", "file:ent?mode=memory&_fk=1")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	char, err := client.Character.Create().
		SetName("SlotTest").
		SetCurrentRoomId(1).
		SetStartingRoomId(1).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create character: %v", err)
	}

	// Create 4 talents for slots 0-3
	for i := 0; i < 4; i++ {
		tal, err := client.Talent.Create().
			SetName(string(rune('A' + i))).
			SetDescription("Talent " + string(rune('A' + i))).
			Save(ctx)
		if err != nil {
			t.Fatalf("failed to create talent: %v", err)
		}

		_, err = client.CharacterTalent.Create().
			SetCharacter(char).
			SetTalent(tal).
			SetSlot(i).
			Save(ctx)
		if err != nil {
			t.Fatalf("failed to create character talent: %v", err)
		}
	}

	// Verify all 4 slots filled
	talents, err := char.QueryTalents().All(ctx)
	if err != nil {
		t.Fatalf("failed to query talents: %v", err)
	}

	slotMap := make(map[int]bool)
	for _, ct := range talents {
		slotMap[ct.Slot] = true
	}

	for i := 0; i < 4; i++ {
		if !slotMap[i] {
			t.Errorf("expected slot %d to be filled", i)
		}
	}
}