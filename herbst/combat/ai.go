package combat

import (
	"context"
	"math/rand"
)

// EnemyAI defines the interface for enemy AI decision making.
// Each tick, the AI decides what action the enemy should take.
type EnemyAI interface {
	// Decide returns the combat action the enemy should take this tick.
	// ctx provides context for cancellation, combat provides combat state,
	// and self is the participant being controlled.
	Decide(ctx context.Context, combat *Combat, self *Participant) *Action
}

// EnemyDefinition defines the stats and behavior of an enemy type.
type EnemyDefinition struct {
	ID           string        // Unique identifier (e.g., "scrap_rat")
	Name         string        // Display name
	HP           int           // Base HP
	AttackTick   int           // Ticks to perform basic attack (1 = instant)
	Abilities    []AbilityDef  // Available abilities
	FleeThreshold float64      // HP % threshold to consider fleeing (0.0-1.0)
	DEX          int           // Dexterity affects flee chance
}

// AbilityDef defines an enemy ability.
type AbilityDef struct {
	Name       string // Ability name
	TickCost   int    // Cast time in ticks
	Damage     int    // Base damage (0 for non-damage abilities)
	HealAmount int    // Heal amount (0 for non-heal abilities)
	Cooldown   int    // Cooldown in ticks after use
	Priority   int    // Action priority (lower = higher)
	TargetSelf bool   // Whether ability targets self
}

// BasicEnemyAI implements a simple decision tree AI for enemies.
type BasicEnemyAI struct {
	Definition *EnemyDefinition
	AbilityMap map[string]int // Tracks cooldowns: ability name -> ticks remaining
}

// NewBasicEnemyAI creates a new BasicEnemyAI for the given definition.
func NewBasicEnemyAI(def *EnemyDefinition) *BasicEnemyAI {
	return &BasicEnemyAI{
		Definition: def,
		AbilityMap: make(map[string]int),
	}
}

// Decide implements EnemyAI.Decide - makes a combat decision each tick.
func (ai *BasicEnemyAI) Decide(ctx context.Context, combat *Combat, self *Participant) *Action {
	// Check for context cancellation
	select {
	case <-ctx.Done():
		return NewAction(0, self.ID, -1, ActionWait, 100)
	default:
	}

	// Decrement cooldowns from previous tick
	ai.tickCooldowns()

	// Get current HP percentage
	hpPercent := float64(self.HP) / float64(self.MaxHP)

	// PHASE 1: Health Check - Flee or Heal
	if hpPercent < ai.Definition.FleeThreshold {
		// Try to heal first if we have a heal ability ready
		if action := ai.tryHealAbility(self); action != nil {
			return action
		}

		// Try to flee - chance based on DEX
		if ai.shouldFlee(self) {
			return NewAction(0, self.ID, -1, ActionFlee, 1)
		}
	}

	// PHASE 2: Ability Check - Use ready abilities
	if action := ai.tryUseAbility(self, combat); action != nil {
		return action
	}

	// PHASE 3: Basic Attack - Fallback
	return ai.basicAttack(self, combat)
}

// tickCooldowns decrements all ability cooldowns by 1.
func (ai *BasicEnemyAI) tickCooldowns() {
	for name, ticks := range ai.AbilityMap {
		if ticks > 0 {
			ai.AbilityMap[name] = ticks - 1
		}
	}
}

// shouldFlee determines if the enemy should attempt to flee.
func (ai *BasicEnemyAI) shouldFlee(self *Participant) bool {
	// Base flee chance is 25%
	// DEX adds 1% per point (higher DEX = smarter about fleeing)
	baseChance := 0.25
	dexBonus := float64(ai.Definition.DEX) * 0.01
	fleeChance := baseChance + dexBonus

	// Cap at 75% flee chance
	if fleeChance > 0.75 {
		fleeChance = 0.75
	}

	return rand.Float64() < fleeChance
}

// tryHealAbility attempts to use a healing ability if available and ready.
func (ai *BasicEnemyAI) tryHealAbility(self *Participant) *Action {
	for _, ability := range ai.Definition.Abilities {
		if ability.HealAmount > 0 && ai.isAbilityReady(ability.Name) {
			// Mark cooldown
			ai.AbilityMap[ability.Name] = ability.Cooldown

			// Create heal action
			action := NewAction(0, self.ID, self.ID, ActionSkill, ability.Priority)
			action.Payload = AbilityPayload{
				Name:       ability.Name,
				Damage:     0,
				HealAmount: ability.HealAmount,
				TickCost:   ability.TickCost,
			}
			return action
		}
	}
	return nil
}

// tryUseAbility attempts to use the best available offensive ability.
func (ai *BasicEnemyAI) tryUseAbility(self *Participant, combat *Combat) *Action {
	// Find target (first enemy participant)
	targetID := ai.findTarget(self.ID, combat)
	if targetID == -1 {
		return nil // No valid target
	}

	// Find best ready ability (highest priority = lowest number)
	var bestAbility *AbilityDef
	for i := range ai.Definition.Abilities {
		ability := &ai.Definition.Abilities[i]
		// Skip healing abilities here (handled in tryHealAbility)
		if ability.HealAmount > 0 {
			continue
		}
		// Check if ability is ready (cooldown == 0)
		if !ai.isAbilityReady(ability.Name) {
			continue
		}
		// First ability found, or better priority
		if bestAbility == nil || ability.Priority < bestAbility.Priority {
			bestAbility = ability
		}
	}

	if bestAbility != nil {
		// Mark cooldown
		ai.AbilityMap[bestAbility.Name] = bestAbility.Cooldown

		// Create ability action
		action := NewAction(0, self.ID, targetID, ActionSkill, bestAbility.Priority)
		action.Payload = AbilityPayload{
			Name:       bestAbility.Name,
			Damage:     bestAbility.Damage,
			HealAmount: 0,
			TickCost:   bestAbility.TickCost,
		}
		return action
	}

	return nil
}

// basicAttack performs a basic attack as fallback.
func (ai *BasicEnemyAI) basicAttack(self *Participant, combat *Combat) *Action {
	targetID := ai.findTarget(self.ID, combat)
	if targetID == -1 {
		// No valid target, wait
		return NewAction(0, self.ID, -1, ActionWait, 100)
	}

	// Basic attack has priority 10
	return NewAction(0, self.ID, targetID, ActionAttack, 10)
}

// findTarget returns the ID of the first enemy participant.
func (ai *BasicEnemyAI) findTarget(selfID int, combat *Combat) int {
	participants := combat.GetParticipants()
	for _, p := range participants {
		if p.ID != selfID && p.HP > 0 {
			return p.ID
		}
	}
	return -1
}

// isAbilityReady checks if an ability is off cooldown.
func (ai *BasicEnemyAI) isAbilityReady(name string) bool {
	ticks, exists := ai.AbilityMap[name]
	return !exists || ticks == 0
}

// AbilityPayload is the payload for skill actions.
type AbilityPayload struct {
	Name       string
	Damage     int
	HealAmount int
	TickCost   int
}

// SetCooldown sets the cooldown for an ability (useful for testing).
func (ai *BasicEnemyAI) SetCooldown(name string, ticks int) {
	ai.AbilityMap[name] = ticks
}

// GetCooldown returns the remaining cooldown for an ability (useful for testing).
func (ai *BasicEnemyAI) GetCooldown(name string) int {
	return ai.AbilityMap[name]
}

// ResetCooldowns resets all ability cooldowns to 0 (useful for testing).
func (ai *BasicEnemyAI) ResetCooldowns() {
	ai.AbilityMap = make(map[string]int)
}