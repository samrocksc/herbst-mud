package combat

import (
	"context"
	"testing"
)

// TestBasicEnemyAI_BasicAttack tests that AI falls back to basic attack.
func TestBasicEnemyAI_BasicAttack(t *testing.T) {
	def := &EnemyDefinition{
		ID:            "test_enemy",
		Name:          "Test Enemy",
		HP:            100,
		AttackTick:    1,
		FleeThreshold: 0.20,
		DEX:           5,
		Abilities:     []AbilityDef{}, // No abilities
	}

	ai := NewBasicEnemyAI(def)
	
	// Create combat with enemy and a target
	combat := &Combat{
		ID:           1,
		StartTick:    0,
		participants: make(map[int]*Participant),
		actionQueue:  make(map[int][]*Action),
		effects:      make(map[int][]*Effect),
	}
	
	// Enemy participant (self)
	enemy := NewParticipant(1, "Test Enemy", 100) // Full HP, won't flee
	combat.participants[1] = enemy
	
	// Target participant
	target := NewParticipant(2, "Player", 100)
	combat.participants[2] = target

	// AI should choose basic attack
	action := ai.Decide(context.Background(), combat, enemy)
	
	if action == nil {
		t.Fatal("Expected non-nil action")
	}
	if action.Type != ActionAttack {
		t.Errorf("Expected ActionAttack, got %v", action.Type)
	}
	if action.Participant != 1 {
		t.Errorf("Expected participant 1, got %d", action.Participant)
	}
	if action.Target != 2 {
		t.Errorf("Expected target 2, got %d", action.Target)
	}
}

// TestBasicEnemyAI_FleeBehavior tests that low HP triggers flee.
func TestBasicEnemyAI_FleeBehavior(t *testing.T) {
	def := &EnemyDefinition{
		ID:            "flee_test",
		Name:          "Flee Test Enemy",
		HP:            100,
		AttackTick:    1,
		FleeThreshold: 0.50, // 50% - high threshold for testing
		DEX:           100,   // Very high DEX = almost guaranteed flee
		Abilities:     []AbilityDef{},
	}

	ai := NewBasicEnemyAI(def)
	
	combat := &Combat{
		ID:           1,
		StartTick:    0,
		participants: make(map[int]*Participant),
		actionQueue:  make(map[int][]*Action),
		effects:      make(map[int][]*Effect),
	}
	
	// Enemy at 30% HP (below 50% threshold)
	enemy := NewParticipant(1, "Flee Test", 30)
	enemy.MaxHP = 100
	combat.participants[1] = enemy
	
	// Target
	target := NewParticipant(2, "Player", 100)
	combat.participants[2] = target

	// Run multiple times to check flee behavior
	fleeCount := 0
	iterations := 100
	
	for i := 0; i < iterations; i++ {
		ai.ResetCooldowns()
		action := ai.Decide(context.Background(), combat, enemy)
		if action.Type == ActionFlee {
			fleeCount++
		}
	}
	
	// With DEX 100 and 30% HP (below 50% threshold), should flee most of the time
	fleeRate := float64(fleeCount) / float64(iterations)
	if fleeRate < 0.70 {
		t.Errorf("Expected flee rate >= 0.70, got %.2f", fleeRate)
	}
}

// TestBasicEnemyAI_HealAbility tests that low HP triggers heal if available.
func TestBasicEnemyAI_HealAbility(t *testing.T) {
	def := &EnemyDefinition{
		ID:            "heal_test",
		Name:          "Heal Test Enemy",
		HP:            100,
		AttackTick:    1,
		FleeThreshold: 0.30,
		DEX:           5,
		Abilities: []AbilityDef{
			{
				Name:       "Heal",
				TickCost:   1,
				HealAmount: 20,
				Cooldown:   5,
				Priority:   1,
				TargetSelf: true,
			},
		},
	}

	ai := NewBasicEnemyAI(def)
	
	combat := &Combat{
		ID:           1,
		StartTick:    0,
		participants: make(map[int]*Participant),
		actionQueue:  make(map[int][]*Action),
		effects:      make(map[int][]*Effect),
	}
	
	// Enemy at 20% HP (below 30% threshold)
	enemy := NewParticipant(1, "Heal Test", 20)
	enemy.MaxHP = 100
	combat.participants[1] = enemy
	
	// Target
	target := NewParticipant(2, "Player", 100)
	combat.participants[2] = target

	// AI should choose to heal
	action := ai.Decide(context.Background(), combat, enemy)
	
	if action == nil {
		t.Fatal("Expected non-nil action")
	}
	if action.Type != ActionSkill {
		t.Errorf("Expected ActionSkill (heal), got %v", action.Type)
	}
	
	// Check payload
	payload, ok := action.Payload.(AbilityPayload)
	if !ok {
		t.Fatal("Expected AbilityPayload in action")
	}
	if payload.Name != "Heal" {
		t.Errorf("Expected ability 'Heal', got %s", payload.Name)
	}
	if payload.HealAmount != 20 {
		t.Errorf("Expected heal amount 20, got %d", payload.HealAmount)
	}
}

// TestBasicEnemyAI_UseOffensiveAbility tests that abilities are used when ready.
func TestBasicEnemyAI_UseOffensiveAbility(t *testing.T) {
	def := &EnemyDefinition{
		ID:            "ability_test",
		Name:          "Ability Test Enemy",
		HP:            100,
		AttackTick:    1,
		FleeThreshold: 0.10, // Low threshold
		DEX:           5,
		Abilities: []AbilityDef{
			{
				Name:     "PowerStrike",
				TickCost: 2,
				Damage:   15,
				Cooldown: 3,
				Priority: 2,
			},
		},
	}

	ai := NewBasicEnemyAI(def)
	
	combat := &Combat{
		ID:           1,
		StartTick:    0,
		participants: make(map[int]*Participant),
		actionQueue:  make(map[int][]*Action),
		effects:      make(map[int][]*Effect),
	}
	
	// Enemy at high HP (no flee)
	enemy := NewParticipant(1, "Ability Test", 100)
	combat.participants[1] = enemy
	
	// Target
	target := NewParticipant(2, "Player", 100)
	combat.participants[2] = target

	// First decision should use the ability
	action := ai.Decide(context.Background(), combat, enemy)
	
	if action == nil {
		t.Fatal("Expected non-nil action")
	}
	if action.Type != ActionSkill {
		t.Errorf("Expected ActionSkill, got %v", action.Type)
	}
	
	payload, ok := action.Payload.(AbilityPayload)
	if !ok {
		t.Fatal("Expected AbilityPayload in action")
	}
	if payload.Name != "PowerStrike" {
		t.Errorf("Expected ability 'PowerStrike', got %s", payload.Name)
	}
	if payload.Damage != 15 {
		t.Errorf("Expected damage 15, got %d", payload.Damage)
	}
	
	// Cooldown should now be set
	if ai.GetCooldown("PowerStrike") != 3 {
		t.Errorf("Expected cooldown 3, got %d", ai.GetCooldown("PowerStrike"))
	}
}

// TestBasicEnemyAI_CooldownRespected tests that abilities on cooldown are skipped.
func TestBasicEnemyAI_CooldownRespected(t *testing.T) {
	def := &EnemyDefinition{
		ID:            "cooldown_test",
		Name:          "Cooldown Test Enemy",
		HP:            100,
		AttackTick:    1,
		FleeThreshold: 0.10,
		DEX:           5,
		Abilities: []AbilityDef{
			{
				Name:     "BigHit",
				TickCost: 1,
				Damage:   20,
				Cooldown: 5,
				Priority: 1,
			},
		},
	}

	ai := NewBasicEnemyAI(def)
	ai.SetCooldown("BigHit", 3) // Ability on cooldown for 3 more ticks
	
	combat := &Combat{
		ID:           1,
		StartTick:    0,
		participants: make(map[int]*Participant),
		actionQueue:  make(map[int][]*Action),
		effects:      make(map[int][]*Effect),
	}
	
	enemy := NewParticipant(1, "Cooldown Test", 100)
	combat.participants[1] = enemy
	
	target := NewParticipant(2, "Player", 100)
	combat.participants[2] = target

	// Should fall back to basic attack since ability on cooldown
	action := ai.Decide(context.Background(), combat, enemy)
	
	if action.Type != ActionAttack {
		t.Errorf("Expected ActionAttack (ability on cooldown), got %v", action.Type)
	}
}

// TestBasicEnemyAI_TickCooldowns tests that cooldowns decrement each tick.
func TestBasicEnemyAI_TickCooldowns(t *testing.T) {
	ai := &BasicEnemyAI{
		Definition: &EnemyDefinition{Name: "test"},
		AbilityMap: map[string]int{
			"Ability1": 3,
			"Ability2": 1,
			"Ability3": 0,
		},
	}
	
	// Tick once
	ai.tickCooldowns()
	
	if ai.AbilityMap["Ability1"] != 2 {
		t.Errorf("Expected Ability1 cooldown 2, got %d", ai.AbilityMap["Ability1"])
	}
	if ai.AbilityMap["Ability2"] != 0 {
		t.Errorf("Expected Ability2 cooldown 0, got %d", ai.AbilityMap["Ability2"])
	}
	if ai.AbilityMap["Ability3"] != 0 {
		t.Errorf("Expected Ability3 cooldown 0, got %d", ai.AbilityMap["Ability3"])
	}
}

// TestBasicEnemyAI_MultipleAbilities tests priority-based ability selection.
func TestBasicEnemyAI_MultipleAbilities(t *testing.T) {
	def := &EnemyDefinition{
		ID:            "multi_ability",
		Name:          "Multi Ability Enemy",
		HP:            100,
		AttackTick:    1,
		FleeThreshold: 0.10,
		DEX:           5,
		Abilities: []AbilityDef{
			{
				Name:     "WeakAttack",
				TickCost: 1,
				Damage:   5,
				Cooldown: 1,
				Priority: 10, // Low priority
			},
			{
				Name:     "StrongAttack",
				TickCost: 2,
				Damage:   15,
				Cooldown: 4,
				Priority: 2, // High priority
			},
			{
				Name:     "MediumAttack",
				TickCost: 1,
				Damage:   10,
				Cooldown: 2,
				Priority: 5, // Medium priority
			},
		},
	}

	ai := NewBasicEnemyAI(def)
	
	combat := &Combat{
		ID:           1,
		StartTick:    0,
		participants: make(map[int]*Participant),
		actionQueue:  make(map[int][]*Action),
		effects:      make(map[int][]*Effect),
	}
	
	enemy := NewParticipant(1, "Multi Ability", 100)
	combat.participants[1] = enemy
	
	target := NewParticipant(2, "Player", 100)
	combat.participants[2] = target

	// Should choose StrongAttack (priority 2)
	action := ai.Decide(context.Background(), combat, enemy)
	
	payload, ok := action.Payload.(AbilityPayload)
	if !ok {
		t.Fatal("Expected AbilityPayload")
	}
	if payload.Name != "StrongAttack" {
		t.Errorf("Expected 'StrongAttack' (priority 2), got %s", payload.Name)
	}
}

// TestBasicEnemyAI_ContextCancellation tests context cancellation handling.
func TestBasicEnemyAI_ContextCancellation(t *testing.T) {
	def := &EnemyDefinition{
		ID:            "cancel_test",
		Name:          "Cancel Test",
		HP:            100,
		AttackTick:    1,
		FleeThreshold: 0.10,
		DEX:           5,
		Abilities:     []AbilityDef{},
	}

	ai := NewBasicEnemyAI(def)
	
	combat := &Combat{
		ID:           1,
		StartTick:    0,
		participants: make(map[int]*Participant),
		actionQueue:  make(map[int][]*Action),
		effects:      make(map[int][]*Effect),
	}
	
	enemy := NewParticipant(1, "Cancel Test", 100)
	combat.participants[1] = enemy
	
	target := NewParticipant(2, "Player", 100)
	combat.participants[2] = target

	// Create cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	
	// Should return wait action when cancelled
	action := ai.Decide(ctx, combat, enemy)
	
	if action == nil {
		t.Fatal("Expected non-nil action")
	}
	if action.Type != ActionWait {
		t.Errorf("Expected ActionWait on context cancellation, got %v", action.Type)
	}
}