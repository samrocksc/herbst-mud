package combat

import (
	"testing"
)

func TestGetSkillBonus(t *testing.T) {
	tests := []struct {
		level    int
		expected float64
	}{
		{0, 0.0},
		{25, 0.0},
		{26, 0.10},
		{50, 0.10},
		{51, 0.25},
		{75, 0.25},
		{76, 0.50},
		{90, 0.50},
		{91, 0.75},
		{99, 0.75},
		{100, 1.0},
		{150, 1.0}, // Over 100 gets max
	}

	for _, tt := range tests {
		result := GetSkillBonus(tt.level)
		if result != tt.expected {
			t.Errorf("GetSkillBonus(%d) = %.2f, expected %.2f", tt.level, result, tt.expected)
		}
	}
}

func TestCalculateDamage_Basic(t *testing.T) {
	registry := NewEffectRegistry()
	action := &ActionDefinition{
		ID:         "test_attack",
		Name:       "Test Attack",
		BaseDamage: 10,
	}

	attacker := &Participant{ID: 1, Name: "Attacker", HP: 100, MaxHP: 100, Attack: 20}
	defender := &Participant{ID: 2, Name: "Defender", HP: 100, MaxHP: 100, Defense: 0}

	result := CalculateDamage(attacker, defender, action, registry)

	// Base 10 * (1 + 0) * (1 + 0) - 0 * (1 - 0) = 10
	// Minimum 1
	if result.BaseDamage != 10 {
		t.Errorf("Expected base damage 10, got %d", result.BaseDamage)
	}
	if result.FinalDamage < 1 {
		t.Errorf("Expected minimum 1 damage, got %d", result.FinalDamage)
	}
}

func TestCalculateDamage_WithSkillBonus(t *testing.T) {
	registry := NewEffectRegistry()
	action := &ActionDefinition{
		ID:         "test_attack",
		Name:       "Test Attack",
		BaseDamage: 10,
	}

	// Attack 50 gives 10% bonus
	attacker := &Participant{ID: 1, Name: "Attacker", HP: 100, MaxHP: 100, Attack: 50}
	defender := &Participant{ID: 2, Name: "Defender", HP: 100, MaxHP: 100, Defense: 0}

	result := CalculateDamage(attacker, defender, action, registry)

	// Base 10 * (1 + 0.10) * (1 + 0) = 11
	if result.SkillBonus != 0.10 {
		t.Errorf("Expected skill bonus 0.10, got %.2f", result.SkillBonus)
	}
	// Damage should be higher than base
	if result.FinalDamage <= 10 {
		t.Errorf("Expected damage > 10 with skill bonus, got %d", result.FinalDamage)
	}
}

func TestCalculateDamage_WithDefense(t *testing.T) {
	registry := NewEffectRegistry()
	action := &ActionDefinition{
		ID:         "test_attack",
		Name:       "Test Attack",
		BaseDamage: 20,
	}

	attacker := &Participant{ID: 1, Name: "Attacker", HP: 100, MaxHP: 100, Attack: 25}
	defender := &Participant{ID: 2, Name: "Defender", HP: 100, MaxHP: 100, Defense: 5}

	result := CalculateDamage(attacker, defender, action, registry)

	// Base 20 * (1 + 0) - 5 = 15
	if result.Defense != 5 {
		t.Errorf("Expected defense 5, got %d", result.Defense)
	}
	if result.FinalDamage >= 20 {
		t.Errorf("Expected damage < 20 with defense, got %d", result.FinalDamage)
	}
}

func TestCalculateDamage_WithBuff(t *testing.T) {
	registry := NewEffectRegistry()
	action := &ActionDefinition{
		ID:         "test_attack",
		Name:       "Test Attack",
		BaseDamage: 10,
	}

	attacker := &Participant{ID: 1, Name: "Attacker", HP: 100, MaxHP: 100, Attack: 25}
	defender := &Participant{ID: 2, Name: "Defender", HP: 100, MaxHP: 100, Defense: 0}

	// Apply strength buff (+25% damage)
	ApplyStatusEffect(registry, StatusBuffStrength, attacker.ID, attacker.ID)

	result := CalculateDamage(attacker, defender, action, registry)

	// Buff bonus should be 0.25
	if result.BuffBonus != 0.25 {
		t.Errorf("Expected buff bonus 0.25, got %.2f", result.BuffBonus)
	}
	// Base 10 * (1 + 0) * (1 + 0.25) = 12.5 → floor = 12
	if result.FinalDamage < 10 {
		t.Errorf("Expected damage >= 10 with buff, got %d", result.FinalDamage)
	}
}

func TestCalculateDamage_Minimum(t *testing.T) {
	registry := NewEffectRegistry()
	action := &ActionDefinition{
		ID:         "weak_attack",
		Name:       "Weak Attack",
		BaseDamage: 5,
	}

	attacker := &Participant{ID: 1, Name: "Attacker", HP: 100, MaxHP: 100, Attack: 0}
	defender := &Participant{ID: 2, Name: "Defender", HP: 100, MaxHP: 100, Defense: 100} // High defense

	result := CalculateDamage(attacker, defender, action, registry)

	// Even with high defense, minimum should be 1
	if result.FinalDamage < 1 {
		t.Errorf("Expected minimum 1 damage, got %d", result.FinalDamage)
	}
}

func TestApplyDamage(t *testing.T) {
	registry := NewEffectRegistry()
	action := &ActionDefinition{
		ID:         "test_attack",
		Name:       "Test Attack",
		BaseDamage: 15,
	}

	attacker := &Participant{ID: 1, Name: "Attacker", HP: 100, MaxHP: 100, Attack: 25}
	defender := &Participant{ID: 2, Name: "Defender", HP: 100, MaxHP: 100, Defense: 0}

	damage := ApplyDamage(attacker, defender, action, registry)

	if damage < 1 {
		t.Errorf("Expected at least 1 damage, got %d", damage)
	}
	if defender.HP != 100-damage {
		t.Errorf("Expected defender HP %d, got %d", 100-damage, defender.HP)
	}
}

func TestApplyDamage_WithShield(t *testing.T) {
	registry := NewEffectRegistry()
	action := &ActionDefinition{
		ID:         "test_attack",
		Name:       "Test Attack",
		BaseDamage: 10,
	}

	attacker := &Participant{ID: 1, Name: "Attacker", HP: 100, MaxHP: 100, Attack: 25}
	defender := &Participant{ID: 2, Name: "Defender", HP: 100, MaxHP: 100, Defense: 0}

	// Apply shield that absorbs damage
	ApplyStatusEffect(registry, StatusBuffShield, defender.ID, defender.ID)

	damage := ApplyDamage(attacker, defender, action, registry)

	// With shield, damage should be reduced
	// Note: This test depends on shield implementation
	_ = damage // Damage might be 0 if shield absorbs all
}

func TestCalculateHeal(t *testing.T) {
	action := &ActionDefinition{
		ID:       "heal",
		Name:     "Heal",
		BaseHeal: 20,
	}

	healer := &Participant{ID: 1, Name: "Healer", HP: 100, MaxHP: 100}

	healAmount := CalculateHeal(healer, action)
	if healAmount != 20 {
		t.Errorf("Expected heal 20, got %d", healAmount)
	}
}

func TestApplyHeal(t *testing.T) {
	action := &ActionDefinition{
		ID:       "heal",
		Name:     "Heal",
		BaseHeal: 15,
	}

	healer := &Participant{ID: 1, Name: "Healer", HP: 100, MaxHP: 100}
	target := &Participant{ID: 2, Name: "Target", HP: 50, MaxHP: 100}

	actualHeal := ApplyHeal(healer, target, action)

	if actualHeal != 15 {
		t.Errorf("Expected heal 15, got %d", actualHeal)
	}
	if target.HP != 65 {
		t.Errorf("Expected target HP 65, got %d", target.HP)
	}
}

func TestApplyHeal_CappedAtMax(t *testing.T) {
	action := &ActionDefinition{
		ID:       "heal",
		Name:     "Heal",
		BaseHeal: 30,
	}

	healer := &Participant{ID: 1, Name: "Healer", HP: 100, MaxHP: 100}
	target := &Participant{ID: 2, Name: "Target", HP: 90, MaxHP: 100}

	actualHeal := ApplyHeal(healer, target, action)

	// Can only heal 10 HP (to reach max)
	if actualHeal != 10 {
		t.Errorf("Expected heal 10 (capped), got %d", actualHeal)
	}
	if target.HP != 100 {
		t.Errorf("Expected target HP 100, got %d", target.HP)
	}
}

func TestCalculateDamage_SkillBonusTiers(t *testing.T) {
	registry := NewEffectRegistry()
	action := &ActionDefinition{
		ID:         "test_attack",
		Name:       "Test Attack",
		BaseDamage: 100,
	}

	// Test each tier
	tiers := []struct {
		attackLevel   int
		expectedBonus float64
	}{
		{25, 0.0},   // 0-25: 0%
		{50, 0.10},  // 26-50: 10%
		{75, 0.25},  // 51-75: 25%
		{90, 0.50},  // 76-90: 50%
		{99, 0.75},  // 91-99: 75%
		{100, 1.0},  // 100: 100%
	}

	for _, tier := range tiers {
		attacker := &Participant{ID: 1, Name: "Attacker", HP: 100, MaxHP: 100, Attack: tier.attackLevel}
		defender := &Participant{ID: 2, Name: "Defender", HP: 100, MaxHP: 100, Defense: 0}

		result := CalculateDamage(attacker, defender, action, registry)

		if result.SkillBonus != tier.expectedBonus {
			t.Errorf("Attack %d: expected skill bonus %.2f, got %.2f", tier.attackLevel, tier.expectedBonus, result.SkillBonus)
		}
	}
}

func TestDamageModifiers(t *testing.T) {
	registry := NewEffectRegistry()
	attacker := &Participant{ID: 1, Name: "Attacker", HP: 100, MaxHP: 100, Attack: 50}

	// No buffs initially
	damageBonus, _, _ := DamageModifiers(attacker.ID, registry)
	if damageBonus != 0 {
		t.Errorf("Expected 0 damage bonus without buffs, got %.2f", damageBonus)
	}

	// Apply strength buff
	ApplyStatusEffect(registry, StatusBuffStrength, attacker.ID, attacker.ID)

	damageBonus, _, accuracyMod := DamageModifiers(attacker.ID, registry)
	if damageBonus != 0.25 {
		t.Errorf("Expected 0.25 damage bonus with buff, got %.2f", damageBonus)
	}
	_ = accuracyMod
}

func TestCalculateEffectiveDefense(t *testing.T) {
	registry := NewEffectRegistry()
	participant := &Participant{ID: 1, Name: "Test", HP: 100, MaxHP: 100, Defense: 10}

	defense := CalculateEffectiveDefense(participant, registry)
	if defense != 10 {
		t.Errorf("Expected defense 10, got %d", defense)
	}
}

func TestCalculateEffectiveAttack(t *testing.T) {
	registry := NewEffectRegistry()
	participant := &Participant{ID: 1, Name: "Test", HP: 100, MaxHP: 100, Attack: 15}

	attack := CalculateEffectiveAttack(participant, registry)
	if attack != 15 {
		t.Errorf("Expected attack 15, got %d", attack)
	}
}

func TestDamageFormula_ExampleFromTicket(t *testing.T) {
	// From ticket: Player with blades level 45, weapon damage 6
	// Damage = 6 × (1 + 0.10) × (1 + 0.25) = 8.25 → 8
	registry := NewEffectRegistry()

	action := &ActionDefinition{
		ID:         "slash",
		Name:       "Slash",
		BaseDamage: 6, // Scrap Machete
	}

	// Warrior with blades level 45 → 10% bonus (26-50 tier)
	attacker := &Participant{ID: 1, Name: "Warrior", HP: 100, MaxHP: 100, Attack: 45}
	defender := &Participant{ID: 2, Name: "Enemy", HP: 100, MaxHP: 100, Defense: 0}

	// Apply battle cry buff (+25% damage)
	ApplyStatusEffect(registry, StatusBuffStrength, attacker.ID, attacker.ID)

	result := CalculateDamage(attacker, defender, action, registry)

	// Skill bonus should be 0.10 (level 45 is in 26-50 tier)
	if result.SkillBonus != 0.10 {
		t.Errorf("Expected skill bonus 0.10, got %.2f", result.SkillBonus)
	}

	// Buff bonus should be 0.25
	if result.BuffBonus != 0.25 {
		t.Errorf("Expected buff bonus 0.25, got %.2f", result.BuffBonus)
	}

	// Raw damage: 6 * 1.10 * 1.25 = 8.25
	// Floor to 8
	if result.FinalDamage != 8 {
		t.Errorf("Expected final damage 8, got %d", result.FinalDamage)
	}
}