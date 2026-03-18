package combat

import (
	"sync"
	"time"
)

// EffectType defines the type of effect
type EffectType string

const (
	// EffectDamage applies damage over time (DoT)
	EffectDamage EffectType = "DAMAGE"
	// EffectHeal applies healing over time (HoT)
	EffectHeal EffectType = "HEAL"
	// EffectBuff enhances a stat
	EffectBuff EffectType = "BUFF"
	// EffectDebuff reduces a stat
	EffectDebuff EffectType = "DEBUFF"
	// EffectStun prevents actions
	EffectStun EffectType = "STUN"
	// EffectRoot prevents movement
	EffectRoot EffectType = "ROOT"
	// EffectShield absorbs damage
	EffectShield EffectType = "SHIELD"
)

// StatType defines which stat an effect modifies
type StatType string

const (
	StatAttack   StatType = "ATTACK"
	StatDefense  StatType = "DEFENSE"
	StatSpeed    StatType = "SPEED"
	StatDexterity StatType = "DEXTERITY"
	StatManaRegen StatType = "MANA_REGEN"
	StatHPRegen  StatType = "HP_REGEN"
)

// ActiveEffect represents an effect currently active on a participant
type ActiveEffect struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	Type         EffectType  `json:"type"`
	TargetID     int         `json:"targetId"`
	SourceID     int         `json:"sourceId"`

	// Effect parameters
	Value        int         `json:"value"`        // Damage/heal amount or stat modifier
	Stat         StatType    `json:"stat"`         // Which stat to modify (for buffs/debuffs)

	// Timing
	TicksRemaining int        `json:"ticksRemaining"`
	AppliedAt      time.Time  `json:"appliedAt"`

	// Stacking
	Stackable    bool        `json:"stackable"`
	Stacks       int         `json:"stacks"`
	MaxStacks    int         `json:"maxStacks"`

	// Dispelling
	Dispellable  bool        `json:"dispellable"`
	DispelType  string      `json:"dispelType"` // "magic", "curse", "poison", etc.
}

// EffectRegistry manages all active effects in combat
type EffectRegistry struct {
	mu      sync.RWMutex
	effects map[int][]*ActiveEffect // participantID -> effects
}

// NewEffectRegistry creates a new effect registry
func NewEffectRegistry() *EffectRegistry {
	return &EffectRegistry{
		effects: make(map[int][]*ActiveEffect),
	}
}

// ApplyEffect adds a new effect to a participant
func (er *EffectRegistry) ApplyEffect(effect *ActiveEffect) {
	er.mu.Lock()
	defer er.mu.Unlock()

	// Check for existing effect if not stackable
	if !effect.Stackable {
		for _, existing := range er.effects[effect.TargetID] {
			if existing.ID == effect.ID {
				// Refresh duration
				existing.TicksRemaining = effect.TicksRemaining
				return
			}
		}
	}

	er.effects[effect.TargetID] = append(er.effects[effect.TargetID], effect)
}

// RemoveEffect removes an effect from a participant
func (er *EffectRegistry) RemoveEffect(participantID int, effectID string) {
	er.mu.Lock()
	defer er.mu.Unlock()

	effects, exists := er.effects[participantID]
	if !exists {
		return
	}

	newEffects := make([]*ActiveEffect, 0)
	for _, e := range effects {
		if e.ID != effectID {
			newEffects = append(newEffects, e)
		}
	}
	er.effects[participantID] = newEffects
}

// RemoveAllEffects removes all effects from a participant
func (er *EffectRegistry) RemoveAllEffects(participantID int) {
	er.mu.Lock()
	defer er.mu.Unlock()
	delete(er.effects, participantID)
}

// GetEffectsForParticipant returns all effects on a participant
func (er *EffectRegistry) GetEffectsForParticipant(participantID int) []*ActiveEffect {
	er.mu.RLock()
	defer er.mu.RUnlock()

	effects, exists := er.effects[participantID]
	if !exists {
		return []*ActiveEffect{}
	}
	return effects
}

// GetEffectsByType returns all effects of a specific type on a participant
func (er *EffectRegistry) GetEffectsByType(participantID int, effectType EffectType) []*ActiveEffect {
	er.mu.RLock()
	defer er.mu.RUnlock()

	var result []*ActiveEffect
	for _, e := range er.effects[participantID] {
		if e.Type == effectType {
			result = append(result, e)
		}
	}
	return result
}

// ProcessEffects processes all effects for a combat tick
// Returns: (damageEffects, healEffects) for logging
func (er *EffectRegistry) ProcessEffects(participants []*Participant) (damageEffects, healEffects []EffectResult) {
	er.mu.Lock()
	defer er.mu.Unlock()

	for _, participant := range participants {
		effects, exists := er.effects[participant.ID]
		if !exists {
			continue
		}

		var newEffects []*ActiveEffect
		for _, effect := range effects {
			// Apply effect
			switch effect.Type {
			case EffectDamage:
				damageEffects = append(damageEffects, EffectResult{
					Effect:  effect,
					Target:  participant,
					Value:   effect.Value * effect.Stacks,
				})
				participant.TakeDamage(effect.Value * effect.Stacks)
			case EffectHeal:
				healEffects = append(healEffects, EffectResult{
					Effect:  effect,
					Target:  participant,
					Value:   effect.Value * effect.Stacks,
				})
				participant.Heal(effect.Value * effect.Stacks)
			}

			// Decrement duration
			effect.TicksRemaining--
			if effect.TicksRemaining > 0 {
				newEffects = append(newEffects, effect)
			}
		}
		er.effects[participant.ID] = newEffects
	}

	return damageEffects, healEffects
}

// GetStatModifier calculates the total stat modifier from buffs/debuffs
func (er *EffectRegistry) GetStatModifier(participantID int, stat StatType) int {
	er.mu.RLock()
	defer er.mu.RUnlock()

	total := 0
	for _, effect := range er.effects[participantID] {
		if effect.Stat == stat {
			if effect.Type == EffectBuff {
				total += effect.Value * effect.Stacks
			} else if effect.Type == EffectDebuff {
				total -= effect.Value * effect.Stacks
			}
		}
	}
	return total
}

// IsStunned returns true if the participant is stunned
func (er *EffectRegistry) IsStunned(participantID int) bool {
	er.mu.RLock()
	defer er.mu.RUnlock()

	for _, effect := range er.effects[participantID] {
		// Check for STUNNED status effect
		if effect.ID == string(StatusStunned) && effect.TicksRemaining > 0 {
			return true
		}
		// Also check for generic EffectStun type
		if effect.Type == EffectStun && effect.TicksRemaining > 0 {
			return true
		}
	}
	return false
}

// IsRooted returns true if the participant is rooted
func (er *EffectRegistry) IsRooted(participantID int) bool {
	er.mu.RLock()
	defer er.mu.RUnlock()

	for _, effect := range er.effects[participantID] {
		if effect.Type == EffectRoot && effect.TicksRemaining > 0 {
			return true
		}
	}
	return false
}

// HasShield returns the total shield value
func (er *EffectRegistry) HasShield(participantID int) int {
	er.mu.RLock()
	defer er.mu.RUnlock()

	total := 0
	for _, effect := range er.effects[participantID] {
		if effect.Type == EffectShield && effect.TicksRemaining > 0 {
			total += effect.Value * effect.Stacks
		}
	}
	return total
}

// EffectResult represents the result of processing an effect
type EffectResult struct {
	Effect *ActiveEffect
	Target *Participant
	Value  int
}

// CountEffects returns the number of effects on a participant
func (er *EffectRegistry) CountEffects(participantID int) int {
	er.mu.RLock()
	defer er.mu.RUnlock()
	return len(er.effects[participantID])
}