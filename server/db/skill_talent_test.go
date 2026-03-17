package db_test

import (
	"context"
	"testing"

	"herbst-server/db"
	"herbst-server/db/skill"
	"herbst-server/db/talent"

	_ "github.com/mattn/go-sqlite3"
	"entgo.io/ent/dialect"
)

// TestSkillSchema tests the Skill entity schema
func TestSkillSchema(t *testing.T) {
	client, err := db.Open(dialect.SQLite, "file:ent?mode=memory&_fk=1")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}

	// Create a skill
	s, err := client.Skill.Create().
		SetName("Fireball").
		SetDescription("A fiery projectile").
		SetSkillType("combat").
		SetCost(10).
		SetCooldown(30).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create skill: %v", err)
	}

	if s.Name != "Fireball" {
		t.Errorf("expected name Fireball, got %s", s.Name)
	}

	if s.SkillType != "combat" {
		t.Errorf("expected skill_type combat, got %s", s.SkillType)
	}

	if s.Cost != 10 {
		t.Errorf("expected cost 10, got %d", s.Cost)
	}

	if s.Cooldown != 30 {
		t.Errorf("expected cooldown 30, got %d", s.Cooldown)
	}

	// Query the skill
	skills, err := client.Skill.Query().All(ctx)
	if err != nil {
		t.Fatalf("failed to query skills: %v", err)
	}

	if len(skills) != 1 {
		t.Errorf("expected 1 skill, got %d", len(skills))
	}

	// Test GetByName
	foundSkill, err := client.Skill.Query().Where(skill.Name("Fireball")).Only(ctx)
	if err != nil {
		t.Fatalf("failed to get skill by name: %v", err)
	}

	if foundSkill.ID != s.ID {
		t.Errorf("expected skill ID %d, got %d", s.ID, foundSkill.ID)
	}
}

// TestTalentSchema tests the Talent entity schema
func TestTalentSchema(t *testing.T) {
	client, err := db.Open(dialect.SQLite, "file:ent?mode=memory&_fk=1")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}

	// Create a talent
	tal, err := client.Talent.Create().
		SetName("Berserker").
		SetDescription("Enter a rage state").
		SetRequirements(`{"min_level": 5, "required_skill": "Melee"}`).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create talent: %v", err)
	}

	if tal.Name != "Berserker" {
		t.Errorf("expected name Berserker, got %s", tal.Name)
	}

	if tal.Description != "Enter a rage state" {
		t.Errorf("expected description 'Enter a rage state', got %s", tal.Description)
	}

	// Query the talent
	talents, err := client.Talent.Query().All(ctx)
	if err != nil {
		t.Fatalf("failed to query talents: %v", err)
	}

	if len(talents) != 1 {
		t.Errorf("expected 1 talent, got %d", len(talents))
	}

	// Test GetByName
	foundTalent, err := client.Talent.Query().Where(talent.Name("Berserker")).Only(ctx)
	if err != nil {
		t.Fatalf("failed to get talent by name: %v", err)
	}

	if foundTalent.ID != tal.ID {
		t.Errorf("expected talent ID %d, got %d", tal.ID, foundTalent.ID)
	}
}

// TestSkillTalentCRUD tests basic CRUD operations for Skills and Talents
func TestSkillTalentCRUD(t *testing.T) {
	client, err := db.Open(dialect.SQLite, "file:ent?mode=memory&_fk=1")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}

	// Create multiple skills
	skills := []struct {
		name      string
		skillType string
		cost      int
	}{
		{"Slash", "combat", 5},
		{"Heal", "magic", 8},
		{"Sneak", "utility", 3},
	}

	for _, s := range skills {
		_, err := client.Skill.Create().
			SetName(s.name).
			SetDescription(s.name + " ability").
			SetSkillType(s.skillType).
			SetCost(s.cost).
			Save(ctx)
		if err != nil {
			t.Fatalf("failed to create skill %s: %v", s.name, err)
		}
	}

	// Create multiple talents
	talentList := []struct {
		name         string
		requirements string
	}{
		{"Double Strike", `{"min_level": 3}`},
		{"Fast Healing", `{"min_level": 5}`},
	}

	for _, ta := range talentList {
		_, err := client.Talent.Create().
			SetName(ta.name).
			SetDescription(ta.name + " talent").
			SetRequirements(ta.requirements).
			Save(ctx)
		if err != nil {
			t.Fatalf("failed to create talent %s: %v", ta.name, err)
		}
	}

	// Count skills
	count, err := client.Skill.Query().Count(ctx)
	if err != nil {
		t.Fatalf("failed to count skills: %v", err)
	}
	if count != 3 {
		t.Errorf("expected 3 skills, got %d", count)
	}

	// Count talents
	count, err = client.Talent.Query().Count(ctx)
	if err != nil {
		t.Fatalf("failed to count talents: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 talents, got %d", count)
	}

	// Query by type
	combatSkills, err := client.Skill.Query().Where(skill.SkillType("combat")).All(ctx)
	if err != nil {
		t.Fatalf("failed to query combat skills: %v", err)
	}
	if len(combatSkills) != 1 {
		t.Errorf("expected 1 combat skill, got %d", len(combatSkills))
	}
}

// TestSkillUpdate tests updating a skill
func TestSkillUpdate(t *testing.T) {
	client, err := db.Open(dialect.SQLite, "file:ent?mode=memory&_fk=1")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}

	// Create a skill
	s, err := client.Skill.Create().
		SetName("Original").
		SetDescription("Original description").
		SetSkillType("combat").
		SetCost(5).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create skill: %v", err)
	}

	// Update the skill
	updated, err := s.Update().
		SetName("Updated").
		SetCost(10).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to update skill: %v", err)
	}

	if updated.Name != "Updated" {
		t.Errorf("expected name Updated, got %s", updated.Name)
	}

	if updated.Cost != 10 {
		t.Errorf("expected cost 10, got %d", updated.Cost)
	}
}

// TestTalentUpdate tests updating a talent
func TestTalentUpdate(t *testing.T) {
	client, err := db.Open(dialect.SQLite, "file:ent?mode=memory&_fk=1")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}

	// Create a talent
	tal, err := client.Talent.Create().
		SetName("Original").
		SetDescription("Original description").
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create talent: %v", err)
	}

	// Update the talent
	updated, err := tal.Update().
		SetDescription("Updated description").
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to update talent: %v", err)
	}

	if updated.Description != "Updated description" {
		t.Errorf("expected description 'Updated description', got %s", updated.Description)
	}
}

// TestSkillDelete tests deleting a skill via client
func TestSkillDelete(t *testing.T) {
	client, err := db.Open(dialect.SQLite, "file:ent?mode=memory&_fk=1")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}

	// Create a skill
	s, err := client.Skill.Create().
		SetName("ToDelete").
		SetDescription("To be deleted").
		SetSkillType("utility").
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create skill: %v", err)
	}

	// Delete the skill via client
	_, err = client.Skill.Delete().Where(skill.ID(s.ID)).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to delete skill: %v", err)
	}

	// Verify deletion
	exists, err := client.Skill.Query().Where(skill.ID(s.ID)).Exist(ctx)
	if err != nil {
		t.Fatalf("failed to check skill existence: %v", err)
	}
	if exists {
		t.Error("expected skill to be deleted, but it still exists")
	}
}

// TestTalentDelete tests deleting a talent via client
func TestTalentDelete(t *testing.T) {
	client, err := db.Open(dialect.SQLite, "file:ent?mode=memory&_fk=1")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}

	// Create a talent
	tal, err := client.Talent.Create().
		SetName("ToDelete").
		SetDescription("To be deleted").
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create talent: %v", err)
	}

	// Delete the talent via client
	_, err = client.Talent.Delete().Where(talent.ID(tal.ID)).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to delete talent: %v", err)
	}

	// Verify deletion
	exists, err := client.Talent.Query().Where(talent.ID(tal.ID)).Exist(ctx)
	if err != nil {
		t.Fatalf("failed to check talent existence: %v", err)
	}
	if exists {
		t.Error("expected talent to be deleted, but it still exists")
	}
}