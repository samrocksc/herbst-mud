package combat

import (
	"testing"
)

// === EffectType Tests ===

func TestEffectType_String(t *testing.T) {
	tests := []struct {
		et     EffectType
		expect string
	}{
		{EffectBleeding, "Bleeding"},
		{EffectPoison, "Poison"},
		{EffectBurning, "Burning"},
		{EffectStunned, "Stunned"},
		{EffectBlinded, "Blinded"},
		{EffectBuffStrength, "Strength"},
		{EffectBuffShield, "Shield"},
	}

	for _, test := range tests {
		if got := test.et.String(); got != test.expect {
			t.Errorf("EffectType(%d).String() = %s, want %s", int(test.et), got, test.expect)
		}
	}
}

func TestEffectType_Category(t *testing.T) {
	tests := []struct {
		et       EffectType
		category EffectCategory
	}{
		{EffectBleeding, CategoryDoT},
		{EffectPoison, CategoryDoT},
		{EffectBurning, CategoryDoT},
		{EffectStunned, CategoryControl},
		{EffectBlinded, CategoryDebuff},
		{EffectBuffStrength, CategoryBuff},
		{EffectBuffShield, CategoryBuff},
	}

	for _, test := range tests {
		if got := test.et.Category(); got != test.category {
			t.Errorf("%s.Category() = %v, want %v", test.et, got, test.category)
		}
	}
}

func TestEffectType_IsDoT(t *testing.T) {
	if !EffectBleeding.IsDoT() {
		t.Error("Bleeding should be DoT")
	}
	if !EffectPoison.IsDoT() {
		t.Error("Poison should be DoT")
	}
	if !EffectBurning.IsDoT() {
		t.Error("Burning should be DoT")
	}
	if EffectStunned.IsDoT() {
		t.Error("Stunned should not be DoT")
	}
	if EffectBuffStrength.IsDoT() {
		t.Error("BuffStrength should not be DoT")
	}
}

func TestEffectType_IsControl(t *testing.T) {
	if !EffectStunned.IsControl() {
		t.Error("Stunned should be control")
	}
	if EffectBleeding.IsControl() {
		t.Error("Bleeding should not be control")
	}
}

func TestEffectType_IsBuff(t *testing.T) {
	if !EffectBuffStrength.IsBuff() {
		t.Error("BuffStrength should be buff")
	}
	if !EffectBuffShield.IsBuff() {
		t.Error("BuffShield should be buff")
	}
	if EffectPoison.IsBuff() {
		t.Error("Poison should not be buff")
	}
}

func TestEffectType_IsDebuff(t *testing.T) {
	if !EffectBleeding.IsDebuff() {
		t.Error("Bleeding should be debuff")
	}
	if !EffectPoison.IsDebuff() {
		t.Error("Poison should be debuff")
	}
	if !EffectBlinded.IsDebuff() {
		t.Error("Blinded should be debuff")
	}
	if EffectBuffStrength.IsDebuff() {
		t.Error("BuffStrength should not be debuff")
	}
}

func TestEffectType_DefaultDuration(t *testing.T) {
	tests := []struct {
		et       EffectType
		duration int
	}{
		{EffectBleeding, 3},
		{EffectPoison, 5},
		{EffectBurning, 4},
		{EffectStunned, 2},
		{EffectBlinded, 3},
		{EffectBuffStrength, 5},
		{EffectBuffShield, 3},
	}

	for _, test := range tests {
		if got := test.et.DefaultDuration(); got != test.duration {
			t.Errorf("%s.DefaultDuration() = %d, want %d", test.et, got, test.duration)
		}
	}
}

func TestEffectType_DefaultPotency(t *testing.T) {
	tests := []struct {
		et       EffectType
		potency  int
	}{
		{EffectBleeding, 1},
		{EffectPoison, 2},
		{EffectBurning, 1},
		{EffectStunned, 0},
		{EffectBlinded, 0},
		{EffectBuffStrength, 25},
		{EffectBuffShield, 25},
	}

	for _, test := range tests {
		if got := test.et.DefaultPotency(); got != test.potency {
			t.Errorf("%s.DefaultPotency() = %d, want %d", test.et, got, test.potency)
		}
	}
}

// === ActiveEffect Tests ===

func TestNewActiveEffect_Defaults(t *testing.T) {
	// Test with zero values - should use defaults
	e := NewActiveEffect(1, EffectPoison, 0, 0, 100, 200)

	if e.ID != 1 {
		t.Errorf("ID = %d, want 1", e.ID)
	}
	if e.Type != EffectPoison {
		t.Errorf("Type = %v, want Poison", e.Type)
	}
	if e.Duration != 5 { // default poison duration
		t.Errorf("Duration = %d, want 5", e.Duration)
	}
	if e.Potency != 2 { // default poison potency
		t.Errorf("Potency = %d, want 2", e.Potency)
	}
	if e.SourceID != 100 {
		t.Errorf("SourceID = %d, want 100", e.SourceID)
	}
	if e.TargetID != 200 {
		t.Errorf("TargetID = %d, want 200", e.TargetID)
	}
}

func TestNewActiveEffect_CustomValues(t *testing.T) {
	// Test with explicit values
	e := NewActiveEffect(42, EffectBleeding, 10, 5, 1, 2)

	if e.ID != 42 {
		t.Errorf("ID = %d, want 42", e.ID)
	}
	if e.Duration != 10 {
		t.Errorf("Duration = %d, want 10", e.Duration)
	}
	if e.Potency != 5 {
		t.Errorf("Potency = %d, want 5", e.Potency)
	}
}

func TestActiveEffect_ProcessTick_Damage(t *testing.T) {
	// Bleeding: 1 damage per tick, 3 duration
	e := NewActiveEffect(1, EffectBleeding, 3, 1, 100, 200)

	// First tick
	damage, expired := e.ProcessTick()
	if damage != 1 {
		t.Errorf("First tick damage = %d, want 1", damage)
	}
	if expired {
		t.Error("First tick should not expire")
	}
	if e.Duration != 2 {
		t.Errorf("Duration after first tick = %d, want 2", e.Duration)
	}
	if e.TickCount != 1 {
		t.Errorf("TickCount = %d, want 1", e.TickCount)
	}

	// Second tick
	damage, expired = e.ProcessTick()
	if damage != 1 {
		t.Errorf("Second tick damage = %d, want 1", damage)
	}
	if expired {
		t.Error("Second tick should not expire")
	}

	// Third tick - should expire
	damage, expired = e.ProcessTick()
	if damage != 1 {
		t.Errorf("Third tick damage = %d, want 1", damage)
	}
	if !expired {
		t.Error("Third tick should expire")
	}
}

func TestActiveEffect_ProcessTick_Poison(t *testing.T) {
	// Poison: 2 damage per tick
	e := NewActiveEffect(1, EffectPoison, 2, 2, 100, 200)

	damage, expired := e.ProcessTick()
	if damage != 2 {
		t.Errorf("Poison damage = %d, want 2", damage)
	}
	if expired {
		t.Error("Poison should not expire on first tick")
	}
}

func TestActiveEffect_ProcessTick_NoDamage(t *testing.T) {
	// Stunned: no damage, just control
	e := NewActiveEffect(1, EffectStunned, 2, 0, 100, 200)

	damage, expired := e.ProcessTick()
	if damage != 0 {
		t.Errorf("Stunned should deal 0 damage, got %d", damage)
	}
	if expired {
		t.Error("Stunned should not expire on first tick")
	}
}

func TestActiveEffect_IsExpired(t *testing.T) {
	e := NewActiveEffect(1, EffectBleeding, 1, 1, 100, 200)

	if e.IsExpired() {
		t.Error("Effect should not be expired initially")
	}

	e.ProcessTick()

	if !e.IsExpired() {
		t.Error("Effect should be expired after last tick")
	}
}

func TestActiveEffect_ExtendDuration(t *testing.T) {
	e := NewActiveEffect(1, EffectPoison, 3, 2, 100, 200)

	e.ExtendDuration(2)

	if e.Duration != 5 {
		t.Errorf("Duration after extend = %d, want 5", e.Duration)
	}
}

func TestActiveEffect_String(t *testing.T) {
	e := NewActiveEffect(42, EffectPoison, 5, 2, 100, 200)
	s := e.String()

	if s == "" {
		t.Error("String() returned empty")
	}
	// Should contain type name
	if !containsSubstring(s, "Poison") {
		t.Errorf("String() = %s, should contain 'Poison'", s)
	}
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// === EffectRegistry Tests ===

func TestNewEffectRegistry(t *testing.T) {
	r := NewEffectRegistry()
	if r == nil {
		t.Fatal("NewEffectRegistry returned nil")
	}
	if r.Count() != 0 {
		t.Errorf("new registry should have 0 effects, got %d", r.Count())
	}
}

func TestEffectRegistry_AddEffect(t *testing.T) {
	r := NewEffectRegistry()

	id := r.AddEffect(EffectPoison, 5, 2, 100, 200)
	if id != 1 {
		t.Errorf("first effect ID = %d, want 1", id)
	}

	id2 := r.AddEffect(EffectBleeding, 3, 1, 101, 201)
	if id2 != 2 {
		t.Errorf("second effect ID = %d, want 2", id2)
	}

	if r.Count() != 2 {
		t.Errorf("Count() = %d, want 2", r.Count())
	}
}

func TestEffectRegistry_AddEffectInstance(t *testing.T) {
	r := NewEffectRegistry()

	e := NewActiveEffect(0, EffectPoison, 5, 2, 100, 200)
	id := r.AddEffectInstance(e)

	if id != 1 {
		t.Errorf("effect ID = %d, want 1", id)
	}
	if e.ID != 1 {
		t.Errorf("effect.ID = %d, want 1", e.ID)
	}
}

func TestEffectRegistry_RemoveEffect(t *testing.T) {
	r := NewEffectRegistry()

	id := r.AddEffect(EffectPoison, 5, 2, 100, 200)
	r.RemoveEffect(id)

	if r.Count() != 0 {
		t.Errorf("Count() after remove = %d, want 0", r.Count())
	}
}

func TestEffectRegistry_GetEffect(t *testing.T) {
	r := NewEffectRegistry()

	id := r.AddEffect(EffectPoison, 5, 2, 100, 200)

	e, ok := r.GetEffect(id)
	if !ok {
		t.Fatal("GetEffect returned false")
	}
	if e.Type != EffectPoison {
		t.Errorf("effect type = %v, want Poison", e.Type)
	}

	_, ok = r.GetEffect(999)
	if ok {
		t.Error("GetEffect should return false for non-existent effect")
	}
}

func TestEffectRegistry_GetEffectsOnTarget(t *testing.T) {
	r := NewEffectRegistry()

	// Add effects for different targets
	r.AddEffect(EffectPoison, 5, 2, 100, 200)    // target 200
	r.AddEffect(EffectBleeding, 3, 1, 100, 200)  // target 200
	r.AddEffect(EffectStunned, 2, 0, 101, 201)   // target 201

	effects := r.GetEffectsOnTarget(200)
	if len(effects) != 2 {
		t.Errorf("target 200 has %d effects, want 2", len(effects))
	}

	effects = r.GetEffectsOnTarget(201)
	if len(effects) != 1 {
		t.Errorf("target 201 has %d effects, want 1", len(effects))
	}

	effects = r.GetEffectsOnTarget(999)
	if len(effects) != 0 {
		t.Errorf("non-existent target has %d effects, want 0", len(effects))
	}
}

func TestEffectRegistry_GetEffectsFromSource(t *testing.T) {
	r := NewEffectRegistry()

	// Add effects from different sources
	r.AddEffect(EffectPoison, 5, 2, 100, 200)    // source 100
	r.AddEffect(EffectBleeding, 3, 1, 100, 201) // source 100
	r.AddEffect(EffectStunned, 2, 0, 101, 202)  // source 101

	effects := r.GetEffectsFromSource(100)
	if len(effects) != 2 {
		t.Errorf("source 100 has %d effects, want 2", len(effects))
	}

	effects = r.GetEffectsFromSource(101)
	if len(effects) != 1 {
		t.Errorf("source 101 has %d effects, want 1", len(effects))
	}
}

func TestEffectRegistry_ProcessAllEffects(t *testing.T) {
	r := NewEffectRegistry()

	// Add DoT effects on two targets
	r.AddEffect(EffectPoison, 2, 3, 100, 200)    // target 200: 3 damage/tick
	r.AddEffect(EffectBleeding, 2, 2, 100, 200)  // target 200: 2 damage/tick
	r.AddEffect(EffectPoison, 2, 2, 101, 201)    // target 201: 2 damage/tick

	result := r.ProcessAllEffects()

	// Check damage by target
	if result.DamageByTarget[200] != 5 {
		t.Errorf("target 200 damage = %d, want 5", result.DamageByTarget[200])
	}
	if result.DamageByTarget[201] != 2 {
		t.Errorf("target 201 damage = %d, want 2", result.DamageByTarget[201])
	}

	// Effects should have 1 tick remaining
	if r.Count() != 3 {
		t.Errorf("Count after first tick = %d, want 3", r.Count())
	}

	// Process second tick - effects should expire
	result = r.ProcessAllEffects()
	if len(result.ExpiredEffects) != 3 {
		t.Errorf("expired effects = %d, want 3", len(result.ExpiredEffects))
	}
	if r.Count() != 0 {
		t.Errorf("Count after all expired = %d, want 0", r.Count())
	}
}

func TestEffectRegistry_CountByType(t *testing.T) {
	r := NewEffectRegistry()

	r.AddEffect(EffectPoison, 5, 2, 100, 200)
	r.AddEffect(EffectPoison, 5, 2, 101, 201)
	r.AddEffect(EffectBleeding, 3, 1, 100, 202)

	if r.CountByType(EffectPoison) != 2 {
		t.Errorf("Poison count = %d, want 2", r.CountByType(EffectPoison))
	}
	if r.CountByType(EffectBleeding) != 1 {
		t.Errorf("Bleeding count = %d, want 1", r.CountByType(EffectBleeding))
	}
	if r.CountByType(EffectStunned) != 0 {
		t.Errorf("Stunned count = %d, want 0", r.CountByType(EffectStunned))
	}
}

func TestEffectRegistry_HasEffect(t *testing.T) {
	r := NewEffectRegistry()

	r.AddEffect(EffectPoison, 5, 2, 100, 200)
	r.AddEffect(EffectStunned, 2, 0, 101, 201)

	if !r.HasEffect(200, EffectPoison) {
		t.Error("target 200 should have Poison")
	}
	if r.HasEffect(200, EffectStunned) {
		t.Error("target 200 should not have Stunned")
	}
	if r.HasEffect(999, EffectPoison) {
		t.Error("non-existent target should not have any effect")
	}
}

func TestEffectRegistry_CalculateDamageModifier(t *testing.T) {
	r := NewEffectRegistry()

	// No effects
	mod := r.CalculateDamageModifier(100)
	if mod != 1.0 {
		t.Errorf("no effects: modifier = %f, want 1.0", mod)
	}

	// With Strength buff (+25%)
	r.AddEffect(EffectBuffStrength, 5, 25, 0, 100)
	mod = r.CalculateDamageModifier(100)
	if mod != 1.25 {
		t.Errorf("with Strength: modifier = %f, want 1.25", mod)
	}

	// Shield doesn't affect outgoing damage
	r.AddEffect(EffectBuffShield, 3, 25, 0, 100)
	mod = r.CalculateDamageModifier(100)
	if mod != 1.25 {
		t.Errorf("with Shield: modifier = %f, want 1.25 (Shield doesn't affect outgoing)", mod)
	}
}

func TestEffectRegistry_CalculateIncomingDamageModifier(t *testing.T) {
	r := NewEffectRegistry()

	// No effects
	mod := r.CalculateIncomingDamageModifier(100)
	if mod != 1.0 {
		t.Errorf("no effects: modifier = %f, want 1.0", mod)
	}

	// With Shield buff (-25% incoming)
	r.AddEffect(EffectBuffShield, 3, 25, 0, 100)
	mod = r.CalculateIncomingDamageModifier(100)
	if mod != 0.75 {
		t.Errorf("with Shield: modifier = %f, want 0.75", mod)
	}
}

func TestEffectRegistry_CalculateAccuracyModifier(t *testing.T) {
	r := NewEffectRegistry()

	// No effects
	mod := r.CalculateAccuracyModifier(100)
	if mod != 1.0 {
		t.Errorf("no effects: modifier = %f, want 1.0", mod)
	}

	// Blinded (-50%)
	r.AddEffect(EffectBlinded, 3, 0, 0, 100)
	mod = r.CalculateAccuracyModifier(100)
	if mod != 0.5 {
		t.Errorf("with Blinded: modifier = %f, want 0.5", mod)
	}

	// Burning (-10%)
	r.AddEffect(EffectBurning, 4, 1, 0, 100)
	mod = r.CalculateAccuracyModifier(100)
	if mod != 0.4 { // 0.5 + 0.1 = 0.6 reduction -> 0.4 remaining
		t.Errorf("with Blinded+Burning: modifier = %f, want 0.4", mod)
	}
}

func TestEffectRegistry_CanAct(t *testing.T) {
	r := NewEffectRegistry()

	// No effects
	if !r.CanAct(100) {
		t.Error("should be able to act with no effects")
	}

	// Stunned
	r.AddEffect(EffectStunned, 2, 0, 0, 100)
	if r.CanAct(100) {
		t.Error("stunned target should not be able to act")
	}

	// Other target should still be able to act
	if !r.CanAct(101) {
		t.Error("non-stunned target should be able to act")
	}
}

func TestEffectRegistry_String(t *testing.T) {
	r := NewEffectRegistry()
	s := r.String()
	if s == "" {
		t.Error("String() returned empty")
	}

	r.AddEffect(EffectPoison, 5, 2, 100, 200)
	s = r.String()
	if s == "" {
		t.Error("String() returned empty")
	}
}

// === Integration Tests ===

func TestEffectFlow_FullCycle(t *testing.T) {
	r := NewEffectRegistry()

	// Simulate combat: Player applies poison to enemy
	poisonID := r.AddEffect(EffectPoison, 3, 2, 100, 200) // 2 damage for 3 ticks

	// Tick 1
	result := r.ProcessAllEffects()
	if result.DamageByTarget[200] != 2 {
		t.Errorf("tick 1: damage = %d, want 2", result.DamageByTarget[200])
	}
	if len(result.ExpiredEffects) != 0 {
		t.Error("tick 1: effect should not expire")
	}
	if !r.HasEffect(200, EffectPoison) {
		t.Error("tick 1: poison should still be active")
	}

	// Tick 2
	result = r.ProcessAllEffects()
	if result.DamageByTarget[200] != 2 {
		t.Errorf("tick 2: damage = %d, want 2", result.DamageByTarget[200])
	}

	// Tick 3 - effect expires
	result = r.ProcessAllEffects()
	if result.DamageByTarget[200] != 2 {
		t.Errorf("tick 3: damage = %d, want 2", result.DamageByTarget[200])
	}
	if len(result.ExpiredEffects) != 1 {
		t.Errorf("tick 3: expired = %d, want 1", len(result.ExpiredEffects))
	}
	if result.ExpiredEffects[0].ID != poisonID {
		t.Errorf("tick 3: expired effect ID = %d, want %d", result.ExpiredEffects[0].ID, poisonID)
	}

	// Verify effect is removed
	if r.HasEffect(200, EffectPoison) {
		t.Error("after expiry: poison should be removed")
	}
}

func TestEffectFlow_MultipleEffectsOnSameTarget(t *testing.T) {
	r := NewEffectRegistry()

	// Apply poison and bleeding to same target
	r.AddEffect(EffectPoison, 3, 2, 100, 200)    // 2 damage
	r.AddEffect(EffectBleeding, 3, 1, 100, 200)  // 1 damage

	result := r.ProcessAllEffects()

	// Total damage should be 3 (2 + 1)
	if result.DamageByTarget[200] != 3 {
		t.Errorf("combined damage = %d, want 3", result.DamageByTarget[200])
	}
}

func TestEffectFlow_MultiTickAccumulation(t *testing.T) {
	r := NewEffectRegistry()

	// Poison: 2 damage per tick for 3 ticks = 6 total damage
	r.AddEffect(EffectPoison, 3, 2, 100, 200)

	totalDamage := 0
	for r.Count() > 0 {
		result := r.ProcessAllEffects()
		totalDamage += result.DamageByTarget[200]
	}

	if totalDamage != 6 {
		t.Errorf("total damage = %d, want 6", totalDamage)
	}
}