package combat

import (
	"context"
	"math/rand"
)

// EnemyAI defines the interface for enemy AI decision making
type EnemyAI interface {
	// Decide chooses an action for the enemy to take
	Decide(ctx context.Context, combat *Combat, self *Participant) *AIAction
}

// AIAction represents an AI's chosen action
type AIAction struct {
	ActionID   string      `json:"actionId"`   // Action to perform (attack, defend, flee, etc.)
	TargetID   int         `json:"targetId"`   // Target participant ID (if single target)
	TargetIDs  []int       `json:"targetIds"`  // Target IDs (if multi-target)
	Priority   int         `json:"priority"`   // Decision priority (higher = more urgent)
	Reason     string      `json:"reason"`     // Why this action was chosen
}

// BasicEnemyAI implements basic enemy AI behavior
type BasicEnemyAI struct {
	// FleeThreshold is the HP percentage below which flee is considered (0.0-1.0)
	FleeThreshold float64
	// HealThreshold is the HP percentage below which heal abilities are prioritized
	HealThreshold float64
	// Aggression controls tendency to use offensive abilities (0.0-1.0)
	Aggression float64
	// Abilities available to this enemy
	Abilities []string
}

// NewBasicEnemyAI creates a basic enemy AI with default settings
func NewBasicEnemyAI() *BasicEnemyAI {
	return &BasicEnemyAI{
		FleeThreshold: 0.25,  // 25% HP
		HealThreshold: 0.30,  // 30% HP
		Aggression:     0.7,  // 70% chance to use offensive abilities
		Abilities:      []string{},
	}
}

// Decide implements EnemyAI.Decide
func (ai *BasicEnemyAI) Decide(ctx context.Context, combat *Combat, self *Participant) *AIAction {
	// Get current HP percentage
	hpPercent := float64(self.HP) / float64(self.MaxHP)

	// 1. HEALTH CHECK - Flee or heal if critically low
	if hpPercent < ai.FleeThreshold {
		// Try to flee based on DEX
		fleeChance := float64(self.Dexterity) / 100.0
		if rand.Float64() < fleeChance {
			return &AIAction{
				ActionID: "flee",
				Priority: 100,
				Reason:   "low_hp_flee",
			}
		}

		// Check for heal abilities
		healAction := ai.findHealAbility()
		if healAction != nil {
			return healAction
		}
	}

	// 2. ABILITY CHECK - Use special abilities when ready
	if len(ai.Abilities) > 0 && rand.Float64() < ai.Aggression {
		abilityAction := ai.chooseAbility(combat, self)
		if abilityAction != nil {
			return abilityAction
		}
	}

	// 3. BASIC ATTACK - Fallback to basic attack
	return ai.basicAttack(combat, self)
}

// findHealAbility searches for a heal ability the enemy can use
func (ai *BasicEnemyAI) findHealAbility() *AIAction {
	for _, abilityID := range ai.Abilities {
		action, exists := GetActionDefinition(abilityID)
		if !exists {
			continue
		}
		if action.BaseHeal > 0 && action.Cooldown == 0 {
			return &AIAction{
				ActionID: abilityID,
				Priority: 90,
				Reason:   "low_hp_heal",
			}
		}
	}
	return nil
}

// chooseAbility selects the best ability to use
func (ai *BasicEnemyAI) chooseAbility(combat *Combat, self *Participant) *AIAction {
	// Get list of alive enemies (players)
	enemies := combat.GetAliveByTeam(0) // Team 0 = players
	if len(enemies) == 0 {
		return nil
	}

	// Find the highest damage ability ready to use
	var bestAbility string
	var bestDamage int
	for _, abilityID := range ai.Abilities {
		action, exists := GetActionDefinition(abilityID)
		if !exists {
			continue
		}
		// Skip heal abilities here (handled separately)
		if action.BaseHeal > 0 {
			continue
		}
		if action.BaseDamage > bestDamage {
			bestAbility = abilityID
			bestDamage = action.BaseDamage
		}
	}

	if bestAbility != "" {
		// Target the player with lowest HP (focus fire)
		target := ai.findWeakestTarget(enemies)
		return &AIAction{
			ActionID: bestAbility,
			TargetID: target.ID,
			Priority: 70,
			Reason:   "ability_offensive",
		}
	}

	return nil
}

// basicAttack returns a basic attack action
func (ai *BasicEnemyAI) basicAttack(combat *Combat, self *Participant) *AIAction {
	enemies := combat.GetAliveByTeam(0)
	if len(enemies) == 0 {
		return nil
	}

	// Target the player with lowest HP
	target := ai.findWeakestTarget(enemies)

	return &AIAction{
		ActionID: "attack",
		TargetID: target.ID,
		Priority: 50,
		Reason:   "basic_attack",
	}
}

// findWeakestTarget finds the target with the lowest HP percentage
func (ai *BasicEnemyAI) findWeakestTarget(targets []*Participant) *Participant {
	if len(targets) == 0 {
		return nil
	}

	weakest := targets[0]
	weakestPercent := float64(weakest.HP) / float64(weakest.MaxHP)

	for _, t := range targets[1:] {
		percent := float64(t.HP) / float64(t.MaxHP)
		if percent < weakestPercent {
			weakest = t
			weakestPercent = percent
		}
	}

	return weakest
}

// EnemyDefinition defines an enemy type's stats and behavior
type EnemyDefinition struct {
	ID             string      `json:"id"`
	Name           string      `json:"name"`
	HP             int         `json:"hp"`
	Attack         int         `json:"attack"`
	Defense        int         `json:"defense"`
	Dexterity      int         `json:"dexterity"`
	AttackTick     int         `json:"attackTick"`     // Ticks for basic attack
	SpecialTick    int         `json:"specialTick"`    // Ticks for special ability (0 if none)
	SpecialName    string      `json:"specialName"`    // Name of special ability
	SpecialDamage  int         `json:"specialDamage"`  // Damage of special ability
	SpecialCooldown int        `json:"specialCooldown"` // Cooldown for special
	FleeChance     float64     `json:"fleeChance"`     // Base flee chance when low HP
	AI             EnemyAI     `json:"-"`              // AI implementation
	Abilities      []string    `json:"abilities"`      // Ability IDs available
}

// EnemyRegistry holds all enemy definitions
var EnemyRegistry = map[string]*EnemyDefinition{
	"scrap_rat": {
		ID:         "scrap_rat",
		Name:       "Scrap Rat",
		HP:         15,
		Attack:     5,
		Defense:    2,
		Dexterity:  8,
		AttackTick: 1,
		FleeChance: 0.20,
		AI: &BasicEnemyAI{
			FleeThreshold: 0.25,
			Aggression:    0.5,
			Abilities:     []string{},
		},
	},
	"junk_dog": {
		ID:              "junk_dog",
		Name:            "Junk Dog",
		HP:              25,
		Attack:          8,
		Defense:         4,
		Dexterity:       6,
		AttackTick:      1,
		SpecialTick:     1,
		SpecialName:     "Bite",
		SpecialDamage:   8,
		SpecialCooldown: 3,
		FleeChance:      0.15,
		Abilities:       []string{"bite"},
		AI: &BasicEnemyAI{
			FleeThreshold: 0.25,
			Aggression:    0.7,
			Abilities:     []string{"bite"},
		},
	},
	"ooze_spawn": {
		ID:              "ooze_spawn",
		Name:            "Ooze Spawn",
		HP:              20,
		Attack:          6,
		Defense:         3,
		Dexterity:       4,
		AttackTick:      1,
		SpecialTick:     3,
		SpecialName:     "Explode",
		SpecialDamage:   15,
		SpecialCooldown: 0, // One-time use
		FleeChance:      0.10,
		Abilities:       []string{"explode"},
		AI: &BasicEnemyAI{
			FleeThreshold: 0.20,
			Aggression:    0.9,
			Abilities:     []string{"explode"},
		},
	},
	"old_scrap": {
		ID:              "old_scrap",
		Name:            "Old Scrap",
		HP:              40,
		Attack:          12,
		Defense:         8,
		Dexterity:       5,
		AttackTick:      1,
		SpecialTick:     2,
		SpecialName:     "Crush",
		SpecialDamage:   12,
		SpecialCooldown: 2,
		FleeChance:      0.25,
		Abilities:       []string{"crush", "scavenge"},
		AI: &BasicEnemyAI{
			FleeThreshold: 0.30,
			HealThreshold: 0.35,
			Aggression:    0.8,
			Abilities:     []string{"crush", "scavenge"},
		},
	},
}

// GetEnemyDefinition retrieves an enemy definition by ID
func GetEnemyDefinition(id string) (*EnemyDefinition, bool) {
	def, exists := EnemyRegistry[id]
	return def, exists
}

// CreateParticipantFromEnemy creates a combat participant from an enemy definition
func CreateParticipantFromEnemy(enemyID string, id int) *Participant {
	def, exists := GetEnemyDefinition(enemyID)
	if !exists {
		return nil
	}

	return &Participant{
		ID:         id,
		Name:       def.Name,
		IsPlayer:   false,
		IsNPC:      true,
		HP:         def.HP,
		MaxHP:      def.HP,
		Attack:     def.Attack,
		Defense:    def.Defense,
		Dexterity:  def.Dexterity,
		IsAlive:    true,
		IsActive:   true,
		Team:       1, // Enemy team
	}
}

// GetAIAction returns the AI's chosen action for this enemy
func GetAIAction(enemyID string, combat *Combat, self *Participant) *AIAction {
	def, exists := GetEnemyDefinition(enemyID)
	if !exists {
		// Fallback to basic AI
		ai := NewBasicEnemyAI()
		return ai.Decide(nil, combat, self)
	}

	if def.AI != nil {
		return def.AI.Decide(nil, combat, self)
	}

	// Default AI
	ai := NewBasicEnemyAI()
	return ai.Decide(nil, combat, self)
}