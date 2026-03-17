package combat

import (
	"testing"
)

func TestCreateStatusEffect_Bleeding(t *testing.T) {
	effect := CreateStatusEffect(StatusBleeding, 1, 2)

	if effect == nil {
		t.Fatal("Expected effect, got nil")
	}
	if effect.Name != "Bleeding" {
		t.Errorf("Expected name 'Bleeding', got '%s'", effect.Name)
	}
	if effect.TicksRemaining != 3 {
		t.Errorf("Expected 3 ticks, got %d", effect.TicksRemaining)
	}
	if effect.Value != 1 {
		t.Errorf("Expected 1 damage per tick, got %d", effect.Value)
	}
	if effect.TargetID != 1 {
		t.Errorf("Expected target ID 1, got %d", effect.TargetID)
	}
}

func TestCreateStatusEffect_Poison(t *testing.T) {
	effect := CreateStatusEffect(StatusPoison, 1, 2)

	if effect == nil {
		t.Fatal("Expected effect, got nil")
	}
	if effect.Name != "Poisoned" {
		t.Errorf("Expected name 'Poisoned', got '%s'", effect.Name)
	}
	if effect.TicksRemaining != 5 {
		t.Errorf("Expected 5 ticks, got %d", effect.TicksRemaining)
	}
	if effect.Value != 2 {
		t.Errorf("Expected 2 damage per tick, got %d", effect.Value)
	}
}

func TestCreateStatusEffect_Stunned(t *testing.T) {
	effect := CreateStatusEffect(StatusStunned, 1, 2)

	if effect == nil {
		t.Fatal("Expected effect, got nil")
	}
	if effect.Name != "Stunned" {
		t.Errorf("Expected name 'Stunned', got '%s'", effect.Name)
	}
	if effect.TicksRemaining != 2 {
		t.Errorf("Expected 2 ticks, got %d", effect.TicksRemaining)
	}
}

func TestCreateStatusEffect_BuffStrength(t *testing.T) {
	effect := CreateStatusEffect(StatusBuffStrength, 1, 2)

	if effect == nil {
		t.Fatal("Expected effect, got nil")
	}
	if effect.Name != "Strength Buff" {
		t.Errorf("Expected name 'Strength Buff', got '%s'", effect.Name)
	}
	if effect.TicksRemaining != 5 {
		t.Errorf("Expected 5 ticks, got %d", effect.TicksRemaining)
	}
}

func TestApplyStatusEffect(t *testing.T) {
	registry := NewEffectRegistry()

	effect := ApplyStatusEffect(registry, StatusBleeding, 1, 2)

	if effect == nil {
		t.Fatal("Expected effect, got nil")
	}

	effects := registry.GetEffectsForParticipant(1)
	if len(effects) != 1 {
		t.Errorf("Expected 1 effect, got %d", len(effects))
	}
	if effects[0].Name != "Bleeding" {
		t.Errorf("Expected 'Bleeding', got '%s'", effects[0].Name)
	}
}

func TestProcessStatusEffectTick_Damage(t *testing.T) {
	effect := CreateStatusEffect(StatusBleeding, 1, 2)
	target := &Participant{
		ID:     1,
		Name:   "Test Target",
		HP:     30,
		MaxHP:  30,
		IsAlive: true,
	}

	damage, action := ProcessStatusEffectTick(effect, target)

	if damage != 1 {
		t.Errorf("Expected 1 damage, got %d", damage)
	}
	if action == "" {
		t.Error("Expected action description, got empty string")
	}
	if target.HP != 29 {
		t.Errorf("Expected HP 29 after bleed, got %d", target.HP)
	}
}

func TestProcessStatusEffectTick_Stun(t *testing.T) {
	effect := CreateStatusEffect(StatusStunned, 1, 2)
	target := &Participant{
		ID:      1,
		Name:    "Test Target",
		HP:      30,
		MaxHP:   30,
		IsAlive: true,
	}

	damage, action := ProcessStatusEffectTick(effect, target)

	if damage != 0 {
		t.Errorf("Expected 0 damage for stun, got %d", damage)
	}
	if action == "" {
		t.Error("Expected stun action description")
	}
	if target.HP != 30 {
		t.Errorf("Stun should not deal damage, HP should be 30, got %d", target.HP)
	}
}

func TestGetAccuracyModifier(t *testing.T) {
	registry := NewEffectRegistry()

	// Apply Blinded effect (-50% accuracy)
	ApplyStatusEffect(registry, StatusBlinded, 1, 2)

	mod := GetAccuracyModifier(registry, 1)
	if mod != -0.50 {
		t.Errorf("Expected -0.50 accuracy mod, got %.2f", mod)
	}

	// Apply Burning (-10% accuracy)
	ApplyStatusEffect(registry, StatusBurning, 1, 2)

	mod = GetAccuracyModifier(registry, 1)
	if mod != -0.60 {
		t.Errorf("Expected -0.60 combined accuracy mod, got %.2f", mod)
	}
}

func TestGetDamageModifier(t *testing.T) {
	registry := NewEffectRegistry()

	// Apply Strength buff (+25% damage)
	ApplyStatusEffect(registry, StatusBuffStrength, 1, 2)

	mod := GetDamageModifier(registry, 1)
	if mod != 0.25 {
		t.Errorf("Expected 0.25 damage mod, got %.2f", mod)
	}
}

func TestGetIncomingDamageModifier(t *testing.T) {
	registry := NewEffectRegistry()

	// Apply Shield buff (-25% incoming damage)
	ApplyStatusEffect(registry, StatusBuffShield, 1, 2)

	mod := GetIncomingDamageModifier(registry, 1)
	if mod != -0.25 {
		t.Errorf("Expected -0.25 incoming damage mod, got %.2f", mod)
	}
}

func TestCanAct_Stunned(t *testing.T) {
	registry := NewEffectRegistry()

	// Without stun, can act
	if !CanAct(registry, 1) {
		t.Error("Expected to be able to act without stun")
	}

	// Apply stun
	ApplyStatusEffect(registry, StatusStunned, 1, 2)

	// With stun, cannot act
	if CanAct(registry, 1) {
		t.Error("Expected to NOT be able to act while stunned")
	}
}

func TestProcessAllStatusEffects(t *testing.T) {
	registry := NewEffectRegistry()
	participants := []*Participant{
		{ID: 1, Name: "Player", HP: 30, MaxHP: 30, IsAlive: true},
		{ID: 2, Name: "Enemy", HP: 20, MaxHP: 20, IsAlive: true},
	}

	// Apply bleed to player (1 damage/tick)
	ApplyStatusEffect(registry, StatusBleeding, 1, 2)
	// Apply poison to enemy (2 damage/tick)
	ApplyStatusEffect(registry, StatusPoison, 2, 1)

	logs := ProcessAllStatusEffects(registry, participants)

	if len(logs) != 2 {
		t.Errorf("Expected 2 log entries, got %d", len(logs))
	}

	// Check player took bleed damage (1 dmg)
	if participants[0].HP != 29 {
		t.Errorf("Expected player HP 29 after bleed, got %d", participants[0].HP)
	}

	// Check enemy took poison damage (2 dmg)
	if participants[1].HP != 18 {
		t.Errorf("Expected enemy HP 18 after poison, got %d", participants[1].HP)
	}
}

func TestStatusEffectDuration(t *testing.T) {
	registry := NewEffectRegistry()

	// Apply 3-tick bleeding
	ApplyStatusEffect(registry, StatusBleeding, 1, 2)

	participant := &Participant{ID: 1, Name: "Test", HP: 30, MaxHP: 30, IsAlive: true}
	participants := []*Participant{participant}

	// Tick 1
	ProcessAllStatusEffects(registry, participants)
	effects := registry.GetEffectsForParticipant(1)
	if len(effects) != 1 {
		t.Errorf("Expected effect to remain after tick 1, got %d effects", len(effects))
	}
	if effects[0].TicksRemaining != 2 {
		t.Errorf("Expected 2 ticks remaining, got %d", effects[0].TicksRemaining)
	}

	// Tick 2
	ProcessAllStatusEffects(registry, participants)
	effects = registry.GetEffectsForParticipant(1)
	if len(effects) != 1 {
		t.Errorf("Expected effect to remain after tick 2, got %d effects", len(effects))
	}

	// Tick 3
	ProcessAllStatusEffects(registry, participants)
	effects = registry.GetEffectsForParticipant(1)
	if len(effects) != 0 {
		t.Errorf("Expected effect to be removed after tick 3, got %d effects", len(effects))
	}
}

func TestStatusEffectNoRefresh(t *testing.T) {
	registry := NewEffectRegistry()

	// Apply bleeding
	ApplyStatusEffect(registry, StatusBleeding, 1, 2)
	effects := registry.GetEffectsForParticipant(1)
	if len(effects) != 1 {
		t.Errorf("Expected 1 effect, got %d", len(effects))
	}
	firstEffect := effects[0]

	// Apply same effect again - should refresh duration, not stack
	ApplyStatusEffect(registry, StatusBleeding, 1, 2)
	effects = registry.GetEffectsForParticipant(1)
	if len(effects) != 1 {
		t.Errorf("Expected 1 effect (refreshed), got %d", len(effects))
	}
	if effects[0] != firstEffect {
		t.Error("Effect should be same instance (refreshed)")
	}
	if effects[0].TicksRemaining != 3 {
		t.Errorf("Expected refreshed duration of 3, got %d", effects[0].TicksRemaining)
	}
}

func TestStatusEffectDefinitions(t *testing.T) {
	// Test all status effects are defined
	expectedEffects := []StatusEffectType{
		StatusBleeding, StatusPoison, StatusBurning,
		StatusStunned, StatusBlinded,
		StatusBuffStrength, StatusBuffShield,
	}

	for _, effectType := range expectedEffects {
		def, exists := GetStatusEffectDefinition(effectType)
		if !exists {
			t.Errorf("Status effect '%s' not defined", effectType)
			continue
		}
		if def.Name == "" {
			t.Errorf("Status effect '%s' has no name", effectType)
		}
		if def.Duration <= 0 {
			t.Errorf("Status effect '%s' has invalid duration: %d", effectType, def.Duration)
		}
	}
}

func TestIsDebuff(t *testing.T) {
	if !IsDebuff(StatusBleeding) {
		t.Error("Bleeding should be a debuff")
	}
	if !IsDebuff(StatusPoison) {
		t.Error("Poison should be a debuff")
	}
	if !IsDebuff(StatusStunned) {
		t.Error("Stunned should be a debuff")
	}
	if IsDebuff(StatusBuffStrength) {
		t.Error("BuffStrength should NOT be a debuff")
	}
}

func TestIsBuff(t *testing.T) {
	if !IsBuff(StatusBuffStrength) {
		t.Error("BuffStrength should be a buff")
	}
	if !IsBuff(StatusBuffShield) {
		t.Error("BuffShield should be a buff")
	}
	if IsBuff(StatusPoison) {
		t.Error("Poison should NOT be a buff")
	}
}

func TestIsDoT(t *testing.T) {
	if !IsDoT(StatusBleeding) {
		t.Error("Bleeding should be DoT")
	}
	if !IsDoT(StatusPoison) {
		t.Error("Poison should be DoT")
	}
	if !IsDoT(StatusBurning) {
		t.Error("Burning should be DoT")
	}
	if IsDoT(StatusStunned) {
		t.Error("Stunned should NOT be DoT")
	}
}

func TestPreventsAction(t *testing.T) {
	if !PreventsAction(StatusStunned) {
		t.Error("Stunned should prevent actions")
	}
	if PreventsAction(StatusPoison) {
		t.Error("Poison should NOT prevent actions")
	}
}

func TestListAllStatusEffects(t *testing.T) {
	effects := ListAllStatusEffects()
	if len(effects) != 7 {
		t.Errorf("Expected 7 status effects, got %d", len(effects))
	}
}

func TestGetActiveStatusEffects(t *testing.T) {
	registry := NewEffectRegistry()

	// No effects initially
	instances := GetActiveStatusEffects(registry, 1)
	if len(instances) != 0 {
		t.Errorf("Expected 0 active effects, got %d", len(instances))
	}

	// Apply multiple effects
	ApplyStatusEffect(registry, StatusBleeding, 1, 2)
	ApplyStatusEffect(registry, StatusBuffStrength, 1, 2)

	instances = GetActiveStatusEffects(registry, 1)
	if len(instances) != 2 {
		t.Errorf("Expected 2 active effects, got %d", len(instances))
	}

	// Check each has definition
	for _, inst := range instances {
		if inst.Definition == nil {
			t.Error("Expected definition, got nil")
		}
		if inst.ActiveEffect == nil {
			t.Error("Expected active effect, got nil")
		}
	}
}