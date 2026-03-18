package combat

import (
	"math"
)

// DamageConfig holds configuration for damage calculation
type DamageConfig struct {
	// Minimum damage (floor)
	MinDamage int
	// Maximum damage before reduction
	MaxDamage int
}

// DefaultDamageConfig is the default damage configuration
var DefaultDamageConfig = DamageConfig{
	MinDamage: 1,
	MaxDamage: 999,
}

// SkillBonusTable defines skill level to damage bonus mapping
// Level 0-25: 0%, 26-50: 10%, 51-75: 25%, 76-90: 50%, 91-99: 75%, 100: 100%
var SkillBonusTable = []struct {
	MinLevel int
	MaxLevel int
	Bonus    float64
}{
	{0, 25, 0.0},
	{26, 50, 0.10},
	{51, 75, 0.25},
	{76, 90, 0.50},
	{91, 99, 0.75},
	{100, 100, 1.0},
}

// GetSkillBonus returns the damage bonus percentage for a skill level
func GetSkillBonus(level int) float64 {
	for _, tier := range SkillBonusTable {
		if level >= tier.MinLevel && level <= tier.MaxLevel {
			return tier.Bonus
		}
	}
	// Level > 100 gets max bonus
	if level > 100 {
		return 1.0
	}
	return 0.0
}

// SkillLevelProvider is a function type for getting skill levels
// This allows different implementations (from DB, from memory, etc.)
type SkillLevelProvider func(participantID int, skillName string) int

// DamageResult represents the result of a damage calculation
type DamageResult struct {
	BaseDamage      int     `json:"baseDamage"`
	SkillBonus      float64 `json:"skillBonus"`
	BuffBonus       float64 `json:"buffBonus"`
	RawDamage       float64 `json:"rawDamage"`
	Defense         int     `json:"defense"`
	ArmorPercent    float64 `json:"armorPercent"`
	FinalDamage     int     `json:"finalDamage"`
	WasBlocked      bool    `json:"wasBlocked"`
	WasCritical     bool    `json:"wasCritical"`
	DamageAbsorbed  int     `json:"damageAbsorbed"`
	ShieldRemaining int     `json:"shieldRemaining"`
}

// CalculateDamage calculates damage using the formula:
// Damage = Base × (1 + SkillBonus) × (1 + BuffBonus) - Defense × (1 - Armor)
// Minimum 1 damage
func CalculateDamage(attacker, defender *Participant, action *ActionDefinition, registry *EffectRegistry) *DamageResult {
	result := &DamageResult{}

	// Base damage from action
	result.BaseDamage = action.BaseDamage
	if result.BaseDamage <= 0 {
		result.BaseDamage = 10 // Default for basic attacks
	}

	// Skill bonus (would need skill level from participant)
	// For now, assume skill level is stored in attacker's stats
	// In real implementation, this would use SkillLevelProvider
	skillLevel := attacker.Attack // Use attack stat as proxy for skill level
	result.SkillBonus = GetSkillBonus(skillLevel)

	// Buff bonus from strength buff
	result.BuffBonus = GetDamageModifier(registry, attacker.ID)

	// Calculate raw damage
	result.RawDamage = float64(result.BaseDamage) * (1.0 + result.SkillBonus) * (1.0 + result.BuffBonus)

	// Defense reduction
	result.Defense = defender.Defense
	result.ArmorPercent = 0.0 // TODO: Get armor from equipment

	// Shield absorbs damage first
	shield := registry.HasShield(defender.ID)
	if shield > 0 {
		if shield >= int(result.RawDamage) {
			result.DamageAbsorbed = int(result.RawDamage)
			result.ShieldRemaining = shield - int(result.RawDamage)
			result.FinalDamage = 0
			result.WasBlocked = true
			return result
		}
		// Shield partially absorbs
		result.DamageAbsorbed = shield
		result.RawDamage -= float64(shield)
	}

	// Defense reduction (after shield)
	defenseReduction := float64(result.Defense) * (1.0 - result.ArmorPercent)
	damageAfterDefense := result.RawDamage - defenseReduction

	// Incoming damage modifier (from buffs/debuffs)
	incomingMod := GetIncomingDamageModifier(registry, defender.ID)
	damageAfterDefense *= (1.0 + incomingMod)

	// Floor to minimum
	result.FinalDamage = int(math.Max(1.0, math.Floor(damageAfterDefense)))

	// Cap at max
	if result.FinalDamage > DefaultDamageConfig.MaxDamage {
		result.FinalDamage = DefaultDamageConfig.MaxDamage
	}

	return result
}

// ApplyDamage applies damage to a participant and returns the actual damage dealt
func ApplyDamage(attacker, defender *Participant, action *ActionDefinition, registry *EffectRegistry) int {
	result := CalculateDamage(attacker, defender, action, registry)
	damage := result.FinalDamage

	// If shield absorbed all, no damage
	if result.WasBlocked && damage == 0 {
		return 0
	}

	// Apply damage to defender
	defender.TakeDamage(damage)

	// Update shield if it absorbed damage
	if result.DamageAbsorbed > 0 && registry != nil {
		// Shield damage is handled by the effect system
		// The shield effect should have its value decremented
		// For now, this is a simplified implementation
	}

	return damage
}

// CalculateHeal calculates healing amount
func CalculateHeal(healer *Participant, action *ActionDefinition) int {
	baseHeal := action.BaseHeal
	if baseHeal <= 0 {
		return 0
	}

	// TODO: Add skill bonus for healing
	// For now, just return base heal
	return baseHeal
}

// ApplyHeal applies healing to a participant and returns actual healing
func ApplyHeal(healer, target *Participant, action *ActionDefinition) int {
	healAmount := CalculateHeal(healer, action)
	actualHeal := target.Heal(healAmount)
	return actualHeal
}

// IsHit determines if an attack hits based on accuracy
func IsHit(attacker, defender *Participant, registry *EffectRegistry) bool {
	// Base hit chance
	hitChance := 0.95 // 95% base hit chance

	// Apply accuracy modifier from debuffs
	accuracyMod := GetAccuracyModifier(registry, attacker.ID)
	hitChance += accuracyMod

	// Clamp to valid range
	if hitChance < 0.0 {
		hitChance = 0.0
	}
	if hitChance > 1.0 {
		hitChance = 1.0
	}

	// Roll for hit
	roll := float64(attacker.Dexterity%100) / 100.0 // Simplified random
	return roll < hitChance
}

// DamageModifiers returns all active damage modifiers for a participant
func DamageModifiers(participantID int, registry *EffectRegistry) (damageBonus, defenseBonus, accuracyMod float64) {
	damageBonus = GetDamageModifier(registry, participantID)
	defenseBonus = 0.0 // TODO: Get defense modifier from effects
	accuracyMod = GetAccuracyModifier(registry, participantID)
	return
}

// GetArmorPercent calculates armor percentage from equipment
// TODO: Implement when equipment system is ready
func GetArmorPercent(participant *Participant) float64 {
	// Placeholder - would calculate from equipped items
	return 0.0
}

// CalculateEffectiveHP calculates HP considering buffs/debuffs
func CalculateEffectiveHP(participant *Participant, registry *EffectRegistry) int {
	baseHP := participant.MaxHP

	// TODO: Add HP buffs/debuffs
	// For now, just return base
	return baseHP
}

// CalculateEffectiveDefense calculates defense considering buffs/debuffs
func CalculateEffectiveDefense(participant *Participant, registry *EffectRegistry) int {
	baseDefense := participant.Defense

	// Add defense from buffs
	defenseMod := registry.GetStatModifier(participant.ID, StatDefense)
	return baseDefense + defenseMod
}

// CalculateEffectiveAttack calculates attack considering buffs/debuffs
func CalculateEffectiveAttack(participant *Participant, registry *EffectRegistry) int {
	baseAttack := participant.Attack

	// Add attack from buffs
	attackMod := registry.GetStatModifier(participant.ID, StatAttack)
	return baseAttack + attackMod
}