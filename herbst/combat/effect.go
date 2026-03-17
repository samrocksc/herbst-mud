package combat

import (
	"fmt"
	"sync"
)

// EffectType defines the type of status effect.
type EffectType int

const (
	// DoT Effects (damage over time)
	EffectBleeding EffectType = iota // Take 1 damage per tick
	EffectPoison                     // Take 2 damage per tick
	EffectBurning                    // Take 1 damage + -10% accuracy

	// Control Effects
	EffectStunned // Can't act for duration
	EffectBlinded // -50% accuracy

	// Buff Effects
	EffectBuffStrength // +25% damage
	EffectBuffShield   // -25% incoming damage
)

// EffectCategory groups effects by their behavior.
type EffectCategory int

const (
	CategoryDoT     EffectCategory = iota // Damage over time
	CategoryControl                       // Impair actions
	CategoryBuff                          // Enhance stats
	CategoryDebuff                        // Weaken stats
)

// effectTypeData holds static data for each effect type.
var effectTypeData = map[EffectType]struct {
	name     string
	category EffectCategory
	potency  int // default potency if not specified
	duration int // default duration in ticks
}{
	EffectBleeding:    {"Bleeding", CategoryDoT, 1, 3},
	EffectPoison:       {"Poison", CategoryDoT, 2, 5},
	EffectBurning:      {"Burning", CategoryDoT, 1, 4},
	EffectStunned:      {"Stunned", CategoryControl, 0, 2},
	EffectBlinded:      {"Blinded", CategoryDebuff, 0, 3},
	EffectBuffStrength: {"Strength", CategoryBuff, 25, 5},
	EffectBuffShield:   {"Shield", CategoryBuff, 25, 3},
}

// String returns the string representation of an EffectType.
func (e EffectType) String() string {
	if data, ok := effectTypeData[e]; ok {
		return data.name
	}
	return fmt.Sprintf("EffectType(%d)", int(e))
}

// Category returns the effect category for this type.
func (e EffectType) Category() EffectCategory {
	if data, ok := effectTypeData[e]; ok {
		return data.category
	}
	return CategoryDebuff
}

// DefaultDuration returns the default duration for this effect type.
func (e EffectType) DefaultDuration() int {
	if data, ok := effectTypeData[e]; ok {
		return data.duration
	}
	return 1
}

// DefaultPotency returns the default potency for this effect type.
func (e EffectType) DefaultPotency() int {
	if data, ok := effectTypeData[e]; ok {
		return data.potency
	}
	return 0
}

// IsDoT returns true if this effect is a damage-over-time effect.
func (e EffectType) IsDoT() bool {
	return e.Category() == CategoryDoT
}

// IsControl returns true if this effect impairs actions.
func (e EffectType) IsControl() bool {
	return e.Category() == CategoryControl
}

// IsBuff returns true if this effect is a beneficial buff.
func (e EffectType) IsBuff() bool {
	return e.Category() == CategoryBuff
}

// IsDebuff returns true if this effect is a harmful debuff.
func (e EffectType) IsDebuff() bool {
	return e.Category() == CategoryDebuff || e.Category() == CategoryDoT
}

// ActiveEffect represents an active status effect on a combatant.
type ActiveEffect struct {
	mu sync.RWMutex

	ID         int
	Type       EffectType
	Duration   int // remaining ticks
	Potency    int // effect strength (damage amount, % modifier, etc.)
	SourceID   int // who applied this effect
	TargetID   int // who this effect is on
	TickCount  int // how many ticks have processed (for tracking)
}

// NewActiveEffect creates a new active effect with the given parameters.
func NewActiveEffect(id int, effectType EffectType, duration, potency, sourceID, targetID int) *ActiveEffect {
	// Use defaults if not specified
	if duration <= 0 {
		duration = effectType.DefaultDuration()
	}
	if potency <= 0 {
		potency = effectType.DefaultPotency()
	}

	return &ActiveEffect{
		ID:        id,
		Type:      effectType,
		Duration:  duration,
		Potency:   potency,
		SourceID:  sourceID,
		TargetID:  targetID,
		TickCount: 0,
	}
}

// Name returns the effect type's name.
func (e *ActiveEffect) Name() string {
	return e.Type.String()
}

// ProcessTick processes a single tick for this effect.
// Returns the damage dealt (if any) and whether the effect expired.
func (e *ActiveEffect) ProcessTick() (damageDealt int, expired bool) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.TickCount++

	var damage int

	switch e.Type {
	case EffectBleeding, EffectPoison, EffectBurning:
		// DoT effects deal damage each tick
		damage = e.Potency
	}

	e.Duration--
	expired = e.Duration <= 0

	return damage, expired
}

// IsExpired returns true if this effect has run its course.
func (e *ActiveEffect) IsExpired() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.Duration <= 0
}

// GetDuration returns the remaining duration.
func (e *ActiveEffect) GetDuration() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.Duration
}

// ExtendDuration adds ticks to the effect's duration.
func (e *ActiveEffect) ExtendDuration(ticks int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.Duration += ticks
}

// String returns a string representation of the effect.
func (e *ActiveEffect) String() string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return fmt.Sprintf("ActiveEffect{id: %d, type: %s, duration: %d, potency: %d}",
		e.ID, e.Type, e.Duration, e.Potency)
}

// EffectRegistry manages all active effects in combat.
type EffectRegistry struct {
	mu      sync.RWMutex
	effects map[int]*ActiveEffect // effect ID -> effect
	nextID  int
}

// NewEffectRegistry creates a new effect registry.
func NewEffectRegistry() *EffectRegistry {
	return &EffectRegistry{
		effects: make(map[int]*ActiveEffect),
		nextID:  1,
	}
}

// AddEffect adds a new effect to the registry.
// Returns the assigned effect ID.
func (r *EffectRegistry) AddEffect(effectType EffectType, duration, potency, sourceID, targetID int) int {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := r.nextID
	r.nextID++

	effect := NewActiveEffect(id, effectType, duration, potency, sourceID, targetID)
	r.effects[id] = effect

	return id
}

// AddEffectInstance adds a pre-created ActiveEffect to the registry.
// Returns the effect ID.
func (r *EffectRegistry) AddEffectInstance(effect *ActiveEffect) int {
	r.mu.Lock()
	defer r.mu.Unlock()

	if effect.ID == 0 {
		effect.ID = r.nextID
		r.nextID++
	}

	r.effects[effect.ID] = effect
	return effect.ID
}

// RemoveEffect removes an effect by ID.
func (r *EffectRegistry) RemoveEffect(effectID int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.effects, effectID)
}

// GetEffect retrieves an effect by ID.
func (r *EffectRegistry) GetEffect(effectID int) (*ActiveEffect, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	effect, ok := r.effects[effectID]
	return effect, ok
}

// GetEffectsOnTarget returns all effects on a specific target.
func (r *EffectRegistry) GetEffectsOnTarget(targetID int) []*ActiveEffect {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*ActiveEffect
	for _, effect := range r.effects {
		if effect.TargetID == targetID {
			result = append(result, effect)
		}
	}
	return result
}

// GetEffectsFromSource returns all effects applied by a specific source.
func (r *EffectRegistry) GetEffectsFromSource(sourceID int) []*ActiveEffect {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*ActiveEffect
	for _, effect := range r.effects {
		if effect.SourceID == sourceID {
			result = append(result, effect)
		}
	}
	return result
}

// ProcessTickResult holds the result of processing all effects for a tick.
type ProcessTickResult struct {
	DamageByTarget map[int]int    // target ID -> total damage
	ExpiredEffects []*ActiveEffect // effects that expired this tick
}

// ProcessAllEffects processes all effects for one tick.
// Returns damage totals per target and list of expired effects.
func (r *EffectRegistry) ProcessAllEffects() ProcessTickResult {
	r.mu.Lock()
	defer r.mu.Unlock()

	result := ProcessTickResult{
		DamageByTarget: make(map[int]int),
		ExpiredEffects: make([]*ActiveEffect, 0),
	}

	for id, effect := range r.effects {
		damage, expired := effect.ProcessTick()

		// Accumulate damage for each target
		if damage > 0 {
			result.DamageByTarget[effect.TargetID] += damage
		}

		// Track expired effects
		if expired {
			result.ExpiredEffects = append(result.ExpiredEffects, effect)
			delete(r.effects, id)
		}
	}

	return result
}

// Count returns the total number of active effects.
func (r *EffectRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.effects)
}

// CountByType returns the number of effects of a specific type.
func (r *EffectRegistry) CountByType(effectType EffectType) int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, effect := range r.effects {
		if effect.Type == effectType {
			count++
		}
	}
	return count
}

// HasEffect checks if a target has a specific effect type active.
func (r *EffectRegistry) HasEffect(targetID int, effectType EffectType) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, effect := range r.effects {
		if effect.TargetID == targetID && effect.Type == effectType {
			return true
		}
	}
	return false
}

// CalculateDamageModifier calculates the damage modifier for a target based on active buffs/debuffs.
// Returns a multiplier (1.0 = normal, >1.0 = increased, <1.0 = reduced).
func (r *EffectRegistry) CalculateDamageModifier(targetID int) float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()

	modifier := 1.0

	for _, effect := range r.effects {
		if effect.TargetID != targetID {
			continue
		}

		switch effect.Type {
		case EffectBuffStrength:
			// +25% damage
			modifier += 0.25
		case EffectBuffShield:
			// This affects incoming damage, not outgoing
			// Will be handled in CalculateIncomingDamageModifier
		}
	}

	return modifier
}

// CalculateIncomingDamageModifier calculates the incoming damage modifier for a target.
// Returns a multiplier for incoming damage (1.0 = normal, <1.0 = reduced).
func (r *EffectRegistry) CalculateIncomingDamageModifier(targetID int) float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()

	modifier := 1.0

	for _, effect := range r.effects {
		if effect.TargetID != targetID {
			continue
		}

		switch effect.Type {
		case EffectBuffShield:
			// -25% incoming damage
			modifier -= 0.25
		}
	}

	// Floor at 0 - can't take negative damage
	if modifier < 0 {
		modifier = 0
	}

	return modifier
}

// CalculateAccuracyModifier calculates the accuracy modifier for a target.
// Returns a multiplier (1.0 = normal, <1.0 = reduced accuracy).
func (r *EffectRegistry) CalculateAccuracyModifier(targetID int) float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()

	modifier := 1.0

	for _, effect := range r.effects {
		if effect.TargetID != targetID {
			continue
		}

		switch effect.Type {
		case EffectBlinded:
			// -50% accuracy
			modifier -= 0.5
		case EffectBurning:
			// -10% accuracy
			modifier -= 0.1
		}
	}

	// Floor at 0 - can't have negative accuracy
	if modifier < 0 {
		modifier = 0
	}

	return modifier
}

// CanAct returns true if the target can take actions this tick.
// Stunned targets cannot act.
func (r *EffectRegistry) CanAct(targetID int) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, effect := range r.effects {
		if effect.TargetID == targetID && effect.Type == EffectStunned {
			return false
		}
	}
	return true
}

// String returns a string representation of the registry.
func (r *EffectRegistry) String() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return fmt.Sprintf("EffectRegistry{count: %d, nextID: %d}", len(r.effects), r.nextID)
}