package combat

import (
	"context"
	"testing"
)

// TestSpawnEnemy tests basic enemy spawning.
func TestSpawnEnemy(t *testing.T) {
	enemy := SpawnEnemy(EnemyTypeScrapRat)
	
	if enemy == nil {
		t.Fatal("Expected non-nil enemy")
	}
	if enemy.Name != "Scrap Rat" {
		t.Errorf("Expected name 'Scrap Rat', got %s", enemy.Name)
	}
	if enemy.MaxHP != 15 {
		t.Errorf("Expected MaxHP 15, got %d", enemy.MaxHP)
	}
	if enemy.HP != 15 {
		t.Errorf("Expected HP 15, got %d", enemy.HP)
	}
}

// TestSpawnEnemyWithAI tests spawning enemy with AI.
func TestSpawnEnemyWithAI(t *testing.T) {
	participant, ai := SpawnEnemyWithAI(EnemyTypeJunkDog)
	
	if participant == nil {
		t.Fatal("Expected non-nil participant")
	}
	if ai == nil {
		t.Fatal("Expected non-nil AI")
	}
	
	if participant.Name != "Junk Dog" {
		t.Errorf("Expected name 'Junk Dog', got %s", participant.Name)
	}
	if participant.MaxHP != 25 {
		t.Errorf("Expected MaxHP 25, got %d", participant.MaxHP)
	}
	if ai.Definition == nil {
		t.Error("Expected AI to have definition")
	}
}

// TestGetEnemyDefinition tests retrieving enemy definitions.
func TestGetEnemyDefinition(t *testing.T) {
	def, ok := GetEnemyDefinition(EnemyTypeScrapRat)
	
	if !ok {
		t.Fatal("Expected to find Scrap Rat definition")
	}
	if def.Name != "Scrap Rat" {
		t.Errorf("Expected name 'Scrap Rat', got %s", def.Name)
	}
	if len(def.Abilities) != 0 {
		t.Errorf("Expected 0 abilities for Scrap Rat, got %d", len(def.Abilities))
	}
}

// TestGetEnemyDefinition_NotFound tests unknown enemy type.
func TestGetEnemyDefinition_NotFound(t *testing.T) {
	_, ok := GetEnemyDefinition(EnemyType("unknown"))
	
	if ok {
		t.Error("Expected false for unknown enemy type")
	}
}

// TestListEnemyTypes tests listing all enemy types.
func TestListEnemyTypes(t *testing.T) {
	types := ListEnemyTypes()
	
	if len(types) != 4 {
		t.Errorf("Expected 4 enemy types, got %d", len(types))
	}
	
	// Check that all expected types exist
	typeMap := make(map[EnemyType]bool)
	for _, et := range types {
		typeMap[et] = true
	}
	
	expected := []EnemyType{EnemyTypeScrapRat, EnemyTypeJunkDog, EnemyTypeOozeSpawn, EnemyTypeOldScrap}
	for _, exp := range expected {
		if !typeMap[exp] {
			t.Errorf("Missing enemy type: %s", exp)
		}
	}
}

// TestEnemyCombat tests EnemyCombat wrapper.
func TestEnemyCombat(t *testing.T) {
	ec := NewEnemyCombat(EnemyTypeOldScrap)
	
	if ec == nil {
		t.Fatal("Expected non-nil EnemyCombat")
	}
	if ec.Participant == nil {
		t.Error("Expected non-nil Participant")
	}
	if ec.AI == nil {
		t.Error("Expected non-nil AI")
	}
	if ec.Definition == nil {
		t.Error("Expected non-nil Definition")
	}
}

// TestEnemyCombat_TakeTurn tests taking a turn.
func TestEnemyCombat_TakeTurn(t *testing.T) {
	ec := NewEnemyCombat(EnemyTypeOldScrap)
	
	combat := &Combat{
		ID:           1,
		StartTick:    0,
		participants: make(map[int]*Participant),
		actionQueue:  make(map[int][]*Action),
		effects:      make(map[int][]*Effect),
	}
	
	combat.participants[ec.Participant.ID] = ec.Participant
	
	// Add target
	target := NewParticipant(999, "Player", 100)
	combat.participants[999] = target
	
	action := ec.TakeTurn(context.Background(), combat)
	
	if action == nil {
		t.Error("Expected non-nil action from TakeTurn")
	}
}

// TestEnemyCombat_IsAlive tests alive check.
func TestEnemyCombat_IsAlive(t *testing.T) {
	ec := NewEnemyCombat(EnemyTypeScrapRat)
	
	if !ec.IsAlive() {
		t.Error("Expected enemy to be alive at full HP")
	}
	
	// Damage the enemy
	ec.Participant.HP = 0
	
	if ec.IsAlive() {
		t.Error("Expected enemy to be dead at 0 HP")
	}
}

// TestEnemyCombat_HPPercent tests HP percentage calculation.
func TestEnemyCombat_HPPercent(t *testing.T) {
	ec := NewEnemyCombat(EnemyTypeScrapRat) // 15 HP
	
	percent := ec.HPPercent()
	if percent != 1.0 {
		t.Errorf("Expected HP percent 1.0, got %.2f", percent)
	}
	
	ec.Participant.HP = 7 // ~50%
	percent = ec.HPPercent()
	if percent < 0.45 || percent > 0.55 {
		t.Errorf("Expected HP percent ~0.5, got %.2f", percent)
	}
	
	ec.Participant.HP = 0
	percent = ec.HPPercent()
	if percent != 0 {
		t.Errorf("Expected HP percent 0, got %.2f", percent)
	}
}

// TestEnemyCombat_ShouldFlee tests flee threshold check.
func TestEnemyCombat_ShouldFlee(t *testing.T) {
	// OldScrap has FleeThreshold 0.25
	ec := NewEnemyCombat(EnemyTypeOldScrap)
	
	// Full HP - should not flee
	if ec.ShouldFlee() {
		t.Error("Should not flee at full HP")
	}
	
	// 50% HP - above threshold
	ec.Participant.HP = 20 // 50% of 40
	if ec.ShouldFlee() {
		t.Error("Should not flee above threshold")
	}
	
	// 20% HP - below threshold
	ec.Participant.HP = 8 // 20% of 40
	if !ec.ShouldFlee() {
		t.Error("Should flee below threshold")
	}
}

// TestScrapRatDefinition tests Scrap Rat enemy definition.
func TestScrapRatDefinition(t *testing.T) {
	def, ok := GetEnemyDefinition(EnemyTypeScrapRat)
	if !ok {
		t.Fatal("Scrap Rat not found")
	}
	
	if def.HP != 15 {
		t.Errorf("Expected HP 15, got %d", def.HP)
	}
	if def.AttackTick != 1 {
		t.Errorf("Expected AttackTick 1, got %d", def.AttackTick)
	}
	if def.FleeThreshold != 0.20 {
		t.Errorf("Expected FleeThreshold 0.20, got %.2f", def.FleeThreshold)
	}
	if def.DEX != 8 {
		t.Errorf("Expected DEX 8, got %d", def.DEX)
	}
	if len(def.Abilities) != 0 {
		t.Errorf("Expected 0 abilities, got %d", len(def.Abilities))
	}
}

// TestJunkDogDefinition tests Junk Dog enemy definition.
func TestJunkDogDefinition(t *testing.T) {
	def, ok := GetEnemyDefinition(EnemyTypeJunkDog)
	if !ok {
		t.Fatal("Junk Dog not found")
	}
	
	if def.HP != 25 {
		t.Errorf("Expected HP 25, got %d", def.HP)
	}
	if len(def.Abilities) != 1 {
		t.Errorf("Expected 1 ability, got %d", len(def.Abilities))
	}
	
	// Check Bite ability
	bite := def.Abilities[0]
	if bite.Name != "Bite" {
		t.Errorf("Expected ability 'Bite', got %s", bite.Name)
	}
	if bite.Damage != 8 {
		t.Errorf("Expected Bite damage 8, got %d", bite.Damage)
	}
	if bite.Cooldown != 3 {
		t.Errorf("Expected Bite cooldown 3, got %d", bite.Cooldown)
	}
}

// TestOozeSpawnDefinition tests Ooze Spawn enemy definition.
func TestOozeSpawnDefinition(t *testing.T) {
	def, ok := GetEnemyDefinition(EnemyTypeOozeSpawn)
	if !ok {
		t.Fatal("Ooze Spawn not found")
	}
	
	if def.HP != 20 {
		t.Errorf("Expected HP 20, got %d", def.HP)
	}
	if def.DEX != 3 {
		t.Errorf("Expected DEX 3, got %d", def.DEX)
	}
	
	// Check Explode ability (for death trigger)
	if len(def.Abilities) != 1 {
		t.Errorf("Expected 1 ability, got %d", len(def.Abilities))
	}
	explode := def.Abilities[0]
	if explode.Name != "Explode" {
		t.Errorf("Expected ability 'Explode', got %s", explode.Name)
	}
}

// TestOldScrapDefinition tests Old Scrap enemy definition.
func TestOldScrapDefinition(t *testing.T) {
	def, ok := GetEnemyDefinition(EnemyTypeOldScrap)
	if !ok {
		t.Fatal("Old Scrap not found")
	}
	
	if def.HP != 40 {
		t.Errorf("Expected HP 40, got %d", def.HP)
	}
	if def.DEX != 10 {
		t.Errorf("Expected DEX 10, got %d", def.DEX)
	}
	
	// Should have 2 abilities: Crush and Scavenge
	if len(def.Abilities) != 2 {
		t.Errorf("Expected 2 abilities, got %d", len(def.Abilities))
	}
}

// TestSpawnEnemy_UnknownType tests fallback for unknown enemy type.
func TestSpawnEnemy_UnknownType(t *testing.T) {
	enemy := SpawnEnemy(EnemyType("unknown_enemy"))
	
	// Should fallback to Scrap Rat
	if enemy.Name != "Scrap Rat" {
		t.Errorf("Expected fallback to Scrap Rat, got %s", enemy.Name)
	}
}

// TestSpawnEnemyUniqueIDs tests that each spawn gets unique ID.
func TestSpawnEnemyUniqueIDs(t *testing.T) {
	enemies := make(map[int]bool)
	
	for i := 0; i < 100; i++ {
		enemy := SpawnEnemy(EnemyTypeScrapRat)
		if enemies[enemy.ID] {
			t.Errorf("Duplicate enemy ID: %d", enemy.ID)
		}
		enemies[enemy.ID] = true
	}
}