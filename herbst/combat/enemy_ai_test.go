package combat

import (
	"context"
	"testing"
)

func TestBasicEnemyAI_BasicAttack(t *testing.T) {
	ai := NewBasicEnemyAI()

	combat := &Combat{
		ID:    1,
		State: StateActive,
		Participants: []*Participant{
			{ID: 1, Name: "Player", IsPlayer: true, Team: 0, HP: 50, MaxHP: 50, IsAlive: true},
			{ID: 2, Name: "Enemy", IsPlayer: false, Team: 1, HP: 30, MaxHP: 30, IsAlive: true},
		},
	}

	self := combat.Participants[1]

	action := ai.Decide(context.Background(), combat, self)

	if action == nil {
		t.Fatal("Expected action, got nil")
	}
	if action.ActionID != "attack" {
		t.Errorf("Expected action 'attack', got '%s'", action.ActionID)
	}
	if action.TargetID != 1 {
		t.Errorf("Expected target ID 1 (player), got %d", action.TargetID)
	}
}

func TestBasicEnemyAI_FleeWhenLowHP(t *testing.T) {
	ai := &BasicEnemyAI{
		FleeThreshold: 0.25,
		HealThreshold: 0.30,
		Aggression:    0.5,
		Abilities:     []string{},
	}

	combat := &Combat{
		ID:    1,
		State: StateActive,
		Participants: []*Participant{
			{ID: 1, Name: "Player", IsPlayer: true, Team: 0, HP: 50, MaxHP: 50, IsAlive: true},
			{ID: 2, Name: "Enemy", IsPlayer: false, Team: 1, HP: 5, MaxHP: 30, Dexterity: 80, IsAlive: true}, // 16% HP
		},
	}

	self := combat.Participants[1]

	// With high DEX, flee chance should be 80%
	// Run multiple times to check flee behavior
	fleeCount := 0
	runs := 100
	for i := 0; i < runs; i++ {
		action := ai.Decide(context.Background(), combat, self)
		if action.ActionID == "flee" {
			fleeCount++
		}
	}

	// Should flee approximately 80% of the time given high DEX
	if fleeCount < 50 { // Allow some variance
		t.Errorf("Expected flee behavior at least 50%% of time (got %d%%)", fleeCount)
	}
}

func TestBasicEnemyAI_HealAbility(t *testing.T) {
	// First, register a heal ability for testing
	SkillActions["enemy_heal"] = &ActionDefinition{
		ID:         "enemy_heal",
		Name:       "Enemy Heal",
		Type:       ActionInstant,
		TickCost:   1,
		BaseHeal:   10,
		Cooldown:   0,
	}
	defer delete(SkillActions, "enemy_heal")

	ai := &BasicEnemyAI{
		FleeThreshold: 0.25,
		HealThreshold: 0.30,
		Aggression:    0.0, // No aggression - should still try to heal
		Abilities:     []string{"enemy_heal"},
	}

	combat := &Combat{
		ID:    1,
		State: StateActive,
		Participants: []*Participant{
			{ID: 1, Name: "Player", IsPlayer: true, Team: 0, HP: 50, MaxHP: 50, IsAlive: true},
			{ID: 2, Name: "Enemy", IsPlayer: false, Team: 1, HP: 5, MaxHP: 30, Dexterity: 0, IsAlive: true}, // Low HP, no DEX for flee
		},
	}

	self := combat.Participants[1]

	action := ai.Decide(context.Background(), combat, self)

	if action == nil {
		t.Fatal("Expected action, got nil")
	}
	if action.ActionID != "enemy_heal" {
		t.Errorf("Expected 'enemy_heal' action, got '%s'", action.ActionID)
	}
	if action.Reason != "low_hp_heal" {
		t.Errorf("Expected reason 'low_hp_heal', got '%s'", action.Reason)
	}
}

func TestBasicEnemyAI_OffensiveAbility(t *testing.T) {
	// Register test ability
	SkillActions["crush"] = &ActionDefinition{
		ID:         "crush",
		Name:       "Crush",
		Type:       ActionInstant,
		TickCost:   2,
		BaseDamage: 12,
		Cooldown:   2,
	}
	defer delete(SkillActions, "crush")

	ai := &BasicEnemyAI{
		FleeThreshold: 0.25,
		Aggression:    1.0, // Always use abilities
		Abilities:     []string{"crush"},
	}

	combat := &Combat{
		ID:    1,
		State: StateActive,
		Participants: []*Participant{
			{ID: 1, Name: "Player", IsPlayer: true, Team: 0, HP: 50, MaxHP: 50, IsAlive: true},
			{ID: 2, Name: "Enemy", IsPlayer: false, Team: 1, HP: 30, MaxHP: 30, IsAlive: true},
		},
	}

	self := combat.Participants[1]

	action := ai.Decide(context.Background(), combat, self)

	if action == nil {
		t.Fatal("Expected action, got nil")
	}
	if action.ActionID != "crush" {
		t.Errorf("Expected 'crush' action, got '%s'", action.ActionID)
	}
	if action.Reason != "ability_offensive" {
		t.Errorf("Expected reason 'ability_offensive', got '%s'", action.Reason)
	}
}

func TestBasicEnemyAI_FocusWeakestTarget(t *testing.T) {
	ai := NewBasicEnemyAI()

	combat := &Combat{
		ID:    1,
		State: StateActive,
		Participants: []*Participant{
			{ID: 1, Name: "Player1", IsPlayer: true, Team: 0, HP: 50, MaxHP: 50, IsAlive: true},
			{ID: 2, Name: "Player2", IsPlayer: true, Team: 0, HP: 10, MaxHP: 50, IsAlive: true}, // Weak target
			{ID: 3, Name: "Enemy", IsPlayer: false, Team: 1, HP: 30, MaxHP: 30, IsAlive: true},
		},
	}

	self := combat.Participants[2]

	// Run multiple times - should always target Player2 (lowest HP %)
	for i := 0; i < 10; i++ {
		action := ai.Decide(context.Background(), combat, self)
		if action.TargetID != 2 {
			t.Errorf("Expected target ID 2 (weakest), got %d", action.TargetID)
		}
	}
}

func TestEnemyRegistry(t *testing.T) {
	// Test that all enemies are defined
	expectedEnemies := []string{"scrap_rat", "junk_dog", "ooze_spawn", "old_scrap"}

	for _, id := range expectedEnemies {
		def, exists := GetEnemyDefinition(id)
		if !exists {
			t.Errorf("Expected enemy '%s' in registry", id)
			continue
		}
		if def.Name == "" {
			t.Errorf("Enemy '%s' has no name", id)
		}
		if def.HP <= 0 {
			t.Errorf("Enemy '%s' has invalid HP: %d", id, def.HP)
		}
		if def.AI == nil {
			t.Errorf("Enemy '%s' has no AI", id)
		}
	}
}

func TestCreateParticipantFromEnemy(t *testing.T) {
	participant := CreateParticipantFromEnemy("scrap_rat", 1)

	if participant == nil {
		t.Fatal("Expected participant, got nil")
	}

	def, _ := GetEnemyDefinition("scrap_rat")

	if participant.Name != def.Name {
		t.Errorf("Expected name '%s', got '%s'", def.Name, participant.Name)
	}
	if participant.HP != def.HP {
		t.Errorf("Expected HP %d, got %d", def.HP, participant.HP)
	}
	if participant.MaxHP != def.HP {
		t.Errorf("Expected MaxHP %d, got %d", def.HP, participant.MaxHP)
	}
	if participant.Team != 1 {
		t.Errorf("Expected team 1 (enemy), got %d", participant.Team)
	}
	if participant.IsPlayer {
		t.Error("Enemy should not be marked as player")
	}
	if !participant.IsNPC {
		t.Error("Enemy should be marked as NPC")
	}
}

func TestCreateParticipantFromEnemy_InvalidID(t *testing.T) {
	participant := CreateParticipantFromEnemy("nonexistent_enemy", 1)

	if participant != nil {
		t.Error("Expected nil for invalid enemy ID, got participant")
	}
}

func TestGetAIAction(t *testing.T) {
	combat := &Combat{
		ID:    1,
		State: StateActive,
		Participants: []*Participant{
			{ID: 1, Name: "Player", IsPlayer: true, Team: 0, HP: 50, MaxHP: 50, IsAlive: true},
			{ID: 2, Name: "Scrap Rat", IsPlayer: false, Team: 1, HP: 15, MaxHP: 15, IsAlive: true},
		},
	}

	self := combat.Participants[1]

	action := GetAIAction("scrap_rat", combat, self)

	if action == nil {
		t.Fatal("Expected action, got nil")
	}
	if action.ActionID != "attack" {
		t.Errorf("Expected 'attack' for scrap_rat, got '%s'", action.ActionID)
	}
}

func TestGetAIAction_InvalidEnemy(t *testing.T) {
	combat := &Combat{
		ID:    1,
		State: StateActive,
		Participants: []*Participant{
			{ID: 1, Name: "Player", IsPlayer: true, Team: 0, HP: 50, MaxHP: 50, IsAlive: true},
			{ID: 2, Name: "Unknown", IsPlayer: false, Team: 1, HP: 10, MaxHP: 10, IsAlive: true},
		},
	}

	self := combat.Participants[1]

	// Should fall back to basic AI
	action := GetAIAction("invalid_enemy", combat, self)

	if action == nil {
		t.Fatal("Expected fallback action, got nil")
	}
	if action.ActionID != "attack" {
		t.Errorf("Expected fallback 'attack', got '%s'", action.ActionID)
	}
}

func TestAIAction_Priority(t *testing.T) {
	// Flee should have highest priority (100)
	// Heal should have high priority (90)
	// Ability should have medium priority (70)
	// Basic attack should have lowest priority (50)

	ai := NewBasicEnemyAI()

	combat := &Combat{
		ID:    1,
		State: StateActive,
		Participants: []*Participant{
			{ID: 1, Name: "Player", IsPlayer: true, Team: 0, HP: 50, MaxHP: 50, IsAlive: true},
			{ID: 2, Name: "Enemy", IsPlayer: false, Team: 1, HP: 30, MaxHP: 30, IsAlive: true},
		},
	}

	self := combat.Participants[1]
	action := ai.Decide(context.Background(), combat, self)

	if action.Priority < 50 {
		t.Errorf("Basic attack should have priority >= 50, got %d", action.Priority)
	}
}

func TestEnemyDefinition_SpecialAbilities(t *testing.T) {
	// Junk Dog should have bite ability
	junkDog, exists := GetEnemyDefinition("junk_dog")
	if !exists {
		t.Fatal("junk_dog not found in registry")
	}
	if len(junkDog.Abilities) != 1 || junkDog.Abilities[0] != "bite" {
		t.Errorf("junk_dog should have 'bite' ability, got %v", junkDog.Abilities)
	}

	// Old Scrap should have crush and scavenge
	oldScrap, exists := GetEnemyDefinition("old_scrap")
	if !exists {
		t.Fatal("old_scrap not found in registry")
	}
	if len(oldScrap.Abilities) != 2 {
		t.Errorf("old_scrap should have 2 abilities, got %d", len(oldScrap.Abilities))
	}
}