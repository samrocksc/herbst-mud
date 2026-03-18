package combat

import (
	"fmt"
	"time"
)

// StatusEffectType defines specific status effect types
type StatusEffectType string

const (
	// DoT effects
	StatusBleeding StatusEffectType = "BLEEDING" // 3 ticks, 1 dmg/tick
	StatusPoison   StatusEffectType = "POISON"   // 5 ticks, 2 dmg/tick
	StatusBurning  StatusEffectType = "BURNING"  // 4 ticks, 1 dmg/tick + -10% accuracy

	// Control effects
	StatusStunned  StatusEffectType = "STUNNED"  // 2 ticks, can't act
	StatusBlinded  StatusEffectType = "BLINDED"  // 3 ticks, -50% accuracy

	// Buff effects
	StatusBuffStrength StatusEffectType = "BUFF_STRENGTH" // 5 ticks, +25% damage
	StatusBuffShield   StatusEffectType = "BUFF_SHIELD"   // 3 ticks, -25% incoming damage
)

// StatusEffectDefinition defines a status effect's properties
type StatusEffectDefinition struct {
	ID          StatusEffectType `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Duration    int              `json:"defaultDuration"` // Default duration in ticks
	DamagePerTick int            `json:"damagePerTick"`
	HealPerTick   int            `json:"healPerTick"`
	
	// Stat modifiers
	AccuracyMod   float64 `json:"accuracyMod"`   // -0.50 for -50%
	DamageMod     float64 `json:"damageMod"`     // +0.25 for +25%
	IncomingDamageMod float64 `json:"incomingDamageMod"` // -0.25 for -25%
	
	// Flags
	PreventsAction bool `json:"preventsAction"` // Stunned, etc.
	IsDoT          bool `json:"isDoT"`
	IsDebuff      bool `json:"isDebuff"`
	IsBuff        bool `json:"isBuff"`
	Dispellable   bool `json:"dispellable"`
}

// StatusEffectDefinitions holds all status effect definitions
var StatusEffectDefinitions = map[StatusEffectType]*StatusEffectDefinition{
	StatusBleeding: {
		ID:            StatusBleeding,
		Name:          "Bleeding",
		Description:   "Taking 1 damage per tick from blood loss.",
		Duration:      3,
		DamagePerTick: 1,
		IsDoT:         true,
		IsDebuff:      true,
		Dispellable:   true,
	},
	StatusPoison: {
		ID:            StatusPoison,
		Name:          "Poisoned",
		Description:   "Taking 2 damage per tick from poison.",
		Duration:      5,
		DamagePerTick: 2,
		IsDoT:         true,
		IsDebuff:      true,
		Dispellable:   true,
	},
	StatusBurning: {
		ID:            StatusBurning,
		Name:          "Burning",
		Description:   "Taking 1 damage per tick and -10% accuracy from flames.",
		Duration:      4,
		DamagePerTick: 1,
		AccuracyMod:   -0.10,
		IsDoT:         true,
		IsDebuff:      true,
		Dispellable:   true,
	},
	StatusStunned: {
		ID:             StatusStunned,
		Name:           "Stunned",
		Description:    "Cannot take any actions.",
		Duration:       2,
		PreventsAction: true,
		IsDebuff:       true,
		Dispellable:    true,
	},
	StatusBlinded: {
		ID:           StatusBlinded,
		Name:         "Blinded",
		Description:  "-50% accuracy due to impaired vision.",
		Duration:     3,
		AccuracyMod: -0.50,
		IsDebuff:     true,
		Dispellable:  true,
	},
	StatusBuffStrength: {
		ID:        StatusBuffStrength,
		Name:      "Strength Buff",
		Description: "+25% damage from battle rage.",
		Duration:  5,
		DamageMod: 0.25,
		IsBuff:    true,
	},
	StatusBuffShield: {
		ID:                StatusBuffShield,
		Name:              "Shielded",
		Description:       "-25% incoming damage from protective barrier.",
		Duration:          3,
		IncomingDamageMod: -0.25,
		IsBuff:            true,
	},
}

// CreateStatusEffect creates a new active status effect
func CreateStatusEffect(effectType StatusEffectType, targetID, sourceID int) *ActiveEffect {
	def, exists := StatusEffectDefinitions[effectType]
	if !exists {
		return nil
	}

	return &ActiveEffect{
		ID:             string(effectType),
		Name:           def.Name,
		Type:           mapEffectType(def),
		TargetID:       targetID,
		SourceID:       sourceID,
		Value:          def.DamagePerTick,
		TicksRemaining: def.Duration,
		AppliedAt:      time.Now(),
		Stackable:      false,
		Stacks:         1,
		Dispellable:    def.Dispellable,
	}
}

// mapEffectType converts StatusEffectType to EffectType for the registry
func mapEffectType(def *StatusEffectDefinition) EffectType {
	if def.IsDoT {
		return EffectDamage
	}
	if def.IsBuff {
		return EffectBuff
	}
	if def.IsDebuff {
		return EffectDebuff
	}
	if def.PreventsAction {
		return EffectStun
	}
	return EffectDebuff
}

// ApplyStatusEffect applies a status effect to a participant
func ApplyStatusEffect(registry *EffectRegistry, effectType StatusEffectType, targetID, sourceID int) *ActiveEffect {
	effect := CreateStatusEffect(effectType, targetID, sourceID)
	if effect == nil {
		return nil
	}
	registry.ApplyEffect(effect)
	return effect
}

// ProcessStatusEffectTick processes a single status effect for one tick
func ProcessStatusEffectTick(effect *ActiveEffect, target *Participant) (damageDealt int, action string) {
	def, exists := StatusEffectDefinitions[StatusEffectType(effect.ID)]
	if !exists {
		return 0, ""
	}

	// Apply DoT damage
	if def.IsDoT && def.DamagePerTick > 0 {
		damageDealt = def.DamagePerTick * effect.Stacks
		target.TakeDamage(damageDealt)
		action = fmt.Sprintf("%s takes %d %s damage", target.Name, damageDealt, def.Name)
	}

	// Stun prevents actions
	if def.PreventsAction {
		action = fmt.Sprintf("%s is stunned and cannot act", target.Name)
	}

	return damageDealt, action
}

// GetAccuracyModifier returns the total accuracy modifier from all effects
func GetAccuracyModifier(registry *EffectRegistry, participantID int) float64 {
	effects := registry.GetEffectsForParticipant(participantID)
	totalMod := 0.0

	for _, effect := range effects {
		def, exists := StatusEffectDefinitions[StatusEffectType(effect.ID)]
		if exists {
			totalMod += def.AccuracyMod
		}
	}

	return totalMod
}

// GetDamageModifier returns the total damage modifier from all effects (buffs)
func GetDamageModifier(registry *EffectRegistry, participantID int) float64 {
	effects := registry.GetEffectsForParticipant(participantID)
	totalMod := 0.0

	for _, effect := range effects {
		def, exists := StatusEffectDefinitions[StatusEffectType(effect.ID)]
		if exists && effect.TicksRemaining > 0 {
			totalMod += def.DamageMod
		}
	}

	return totalMod
}

// GetIncomingDamageModifier returns the modifier for incoming damage (shields)
func GetIncomingDamageModifier(registry *EffectRegistry, participantID int) float64 {
	effects := registry.GetEffectsForParticipant(participantID)
	totalMod := 0.0

	for _, effect := range effects {
		def, exists := StatusEffectDefinitions[StatusEffectType(effect.ID)]
		if exists && effect.TicksRemaining > 0 {
			totalMod += def.IncomingDamageMod
		}
	}

	return totalMod
}

// CanAct returns false if the participant is stunned or otherwise prevented from acting
func CanAct(registry *EffectRegistry, participantID int) bool {
	return !registry.IsStunned(participantID)
}

// GetActiveStatusEffects returns all active status effects with their definitions
func GetActiveStatusEffects(registry *EffectRegistry, participantID int) []StatusEffectInstance {
	effects := registry.GetEffectsForParticipant(participantID)
	instances := make([]StatusEffectInstance, 0, len(effects))

	for _, effect := range effects {
		def, exists := StatusEffectDefinitions[StatusEffectType(effect.ID)]
		if exists {
			instances = append(instances, StatusEffectInstance{
				Definition:     def,
				ActiveEffect:   effect,
				TicksRemaining: effect.TicksRemaining,
			})
		}
	}

	return instances
}

// StatusEffectInstance combines a definition with its active state
type StatusEffectInstance struct {
	Definition     *StatusEffectDefinition
	ActiveEffect   *ActiveEffect
	TicksRemaining int
}

// ProcessAllStatusEffects processes all status effects for a tick
// This applies DoT damage and decrements durations
func ProcessAllStatusEffects(registry *EffectRegistry, participants []*Participant) []StatusEffectLog {
	var logs []StatusEffectLog

	registry.mu.Lock()
	defer registry.mu.Unlock()

	for _, participant := range participants {
		effects, exists := registry.effects[participant.ID]
		if !exists {
			continue
		}

		var remainingEffects []*ActiveEffect
		for _, effect := range effects {
			// Get definition for this effect
			def, exists := StatusEffectDefinitions[StatusEffectType(effect.ID)]
			if !exists {
				// Not a status effect, keep it
				remainingEffects = append(remainingEffects, effect)
				continue
			}

			// Apply DoT damage
			var damageDealt int
			var action string
			if def.IsDoT && def.DamagePerTick > 0 {
				damageDealt = def.DamagePerTick * effect.Stacks
				participant.TakeDamage(damageDealt)
				action = fmt.Sprintf("%s takes %d %s damage", participant.Name, damageDealt, def.Name)
			}

			// Log stun
			if def.PreventsAction {
				action = fmt.Sprintf("%s is stunned and cannot act", participant.Name)
			}

			// Log effect
			if action != "" {
				logs = append(logs, StatusEffectLog{
					TargetID:    participant.ID,
					TargetName:  participant.Name,
					EffectID:    effect.ID,
					EffectName:  effect.Name,
					DamageDealt: damageDealt,
					Action:      action,
					Tick:        effect.TicksRemaining - 1, // After this tick
				})
			}

			// Decrement duration
			effect.TicksRemaining--
			if effect.TicksRemaining > 0 {
				remainingEffects = append(remainingEffects, effect)
			}
		}
		registry.effects[participant.ID] = remainingEffects
	}

	return logs
}

// StatusEffectLog represents a log entry for status effect processing
type StatusEffectLog struct {
	TargetID    int    `json:"targetId"`
	TargetName  string `json:"targetName"`
	EffectID    string `json:"effectId"`
	EffectName  string `json:"effectName"`
	DamageDealt int    `json:"damageDealt"`
	Action      string `json:"action"`
	Tick        int    `json:"tick"` // Ticks remaining after processing
}

// GetStatusEffectDefinition retrieves a status effect definition
func GetStatusEffectDefinition(effectType StatusEffectType) (*StatusEffectDefinition, bool) {
	def, exists := StatusEffectDefinitions[effectType]
	return def, exists
}

// ListAllStatusEffects returns all available status effect definitions
func ListAllStatusEffects() []*StatusEffectDefinition {
	defs := make([]*StatusEffectDefinition, 0, len(StatusEffectDefinitions))
	for _, def := range StatusEffectDefinitions {
		defs = append(defs, def)
	}
	return defs
}

// IsDebuff checks if a status effect is a debuff
func IsDebuff(effectType StatusEffectType) bool {
	def, exists := StatusEffectDefinitions[effectType]
	if !exists {
		return false
	}
	return def.IsDebuff
}

// IsBuff checks if a status effect is a buff
func IsBuff(effectType StatusEffectType) bool {
	def, exists := StatusEffectDefinitions[effectType]
	if !exists {
		return false
	}
	return def.IsBuff
}

// IsDoT checks if a status effect is damage over time
func IsDoT(effectType StatusEffectType) bool {
	def, exists := StatusEffectDefinitions[effectType]
	if !exists {
		return false
	}
	return def.IsDoT
}

// PreventsAction checks if a status effect prevents actions
func PreventsAction(effectType StatusEffectType) bool {
	def, exists := StatusEffectDefinitions[effectType]
	if !exists {
		return false
	}
	return def.PreventsAction
}