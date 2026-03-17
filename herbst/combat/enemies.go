package combat

import (
	"context"
	"sync/atomic"
)

// enemyIDCounter generates unique IDs for enemy instances.
var enemyIDCounter int64

// EnemyType defines the type identifier for enemies.
type EnemyType string

const (
	EnemyTypeScrapRat  EnemyType = "scrap_rat"
	EnemyTypeJunkDog   EnemyType = "junk_dog"
	EnemyTypeOozeSpawn EnemyType = "ooze_spawn"
	EnemyTypeOldScrap  EnemyType = "old_scrap"
)

// EnemyRegistry contains all enemy definitions.
var EnemyRegistry = map[EnemyType]*EnemyDefinition{
	EnemyTypeScrapRat: {
		ID:            string(EnemyTypeScrapRat),
		Name:          "Scrap Rat",
		HP:            15,
		AttackTick:    1,
		FleeThreshold: 0.20, // Flee at 20% HP
		DEX:           8,
		Abilities:     []AbilityDef{}, // No special abilities
	},
	EnemyTypeJunkDog: {
		ID:            string(EnemyTypeJunkDog),
		Name:          "Junk Dog",
		HP:            25,
		AttackTick:    1,
		FleeThreshold: 0.15, // Flee at 15% HP
		DEX:           6,
		Abilities: []AbilityDef{
			{
				Name:     "Bite",
				TickCost: 1,
				Damage:   8,
				Cooldown: 3,
				Priority: 5,
			},
		},
	},
	EnemyTypeOozeSpawn: {
		ID:            string(EnemyTypeOozeSpawn),
		Name:          "Ooze Spawn",
		HP:            20,
		AttackTick:    1,
		FleeThreshold: 0.10, // Flee at 10% HP (rarely flees)
		DEX:           3,
		Abilities: []AbilityDef{
			{
				Name:     "Explode",
				TickCost: 3,
				Damage:   15,
				Cooldown: 0,    // Only used on death trigger
				Priority: 1,    // High priority when triggered
			},
		},
	},
	EnemyTypeOldScrap: {
		ID:            string(EnemyTypeOldScrap),
		Name:          "Old Scrap",
		HP:            40,
		AttackTick:    1,
		FleeThreshold: 0.25, // Flee at 25% HP
		DEX:           10,
		Abilities: []AbilityDef{
			{
				Name:     "Crush",
				TickCost: 2,
				Damage:   12,
				Cooldown: 4,
				Priority: 3,
			},
			{
				Name:       "Scavenge",
				TickCost:   2,
				HealAmount: 10,
				Cooldown:   5,
				Priority:   2, // Higher priority than Crush
				TargetSelf: true,
			},
		},
	},
}

// SpawnEnemy creates a new enemy combat participant from an EnemyType.
// Returns a Participant ready to be added to combat.
func SpawnEnemy(enemyType EnemyType) *Participant {
	def, ok := EnemyRegistry[enemyType]
	if !ok {
		// Fallback to basic Scrap Rat
		def = EnemyRegistry[EnemyTypeScrapRat]
	}

	// Generate unique ID
	id := int(atomic.AddInt64(&enemyIDCounter, 1))

	// Create participant with enemy stats
	p := NewParticipant(id, def.Name, def.HP)
	
	// Initialize cooldowns map for abilities
	if p.Cooldowns == nil {
		p.Cooldowns = make(map[string]int)
	}

	return p
}

// SpawnEnemyWithAI creates an enemy participant with its AI controller.
// Returns both the participant and the AI for decision making.
func SpawnEnemyWithAI(enemyType EnemyType) (*Participant, *BasicEnemyAI) {
	def, ok := EnemyRegistry[enemyType]
	if !ok {
		def = EnemyRegistry[EnemyTypeScrapRat]
	}

	// Create participant
	p := SpawnEnemy(enemyType)
	
	// Create AI instance
	ai := NewBasicEnemyAI(def)

	return p, ai
}

// GetEnemyDefinition returns the definition for an enemy type.
func GetEnemyDefinition(enemyType EnemyType) (*EnemyDefinition, bool) {
	def, ok := EnemyRegistry[enemyType]
	return def, ok
}

// ListEnemyTypes returns all available enemy types.
func ListEnemyTypes() []EnemyType {
	types := make([]EnemyType, 0, len(EnemyRegistry))
	for t := range EnemyRegistry {
		types = append(types, t)
	}
	return types
}

// EnemyCombat wraps an enemy participant with its AI controller.
// This is a convenience type for managing enemy state in combat.
type EnemyCombat struct {
	Participant *Participant
	AI          EnemyAI
	Definition  *EnemyDefinition
}

// NewEnemyCombat creates a new EnemyCombat for the given enemy type.
func NewEnemyCombat(enemyType EnemyType) *EnemyCombat {
	p, ai := SpawnEnemyWithAI(enemyType)
	def, _ := GetEnemyDefinition(enemyType)
	
	return &EnemyCombat{
		Participant: p,
		AI:          ai,
		Definition:  def,
	}
}

// TakeTurn executes one turn of AI decision making and returns the action.
func (ec *EnemyCombat) TakeTurn(ctx context.Context, combat *Combat) *Action {
	return ec.AI.Decide(ctx, combat, ec.Participant)
}

// IsAlive returns true if the enemy is still alive (HP > 0).
func (ec *EnemyCombat) IsAlive() bool {
	return ec.Participant.HP > 0
}

// HPPercent returns the current HP as a percentage of max HP.
func (ec *EnemyCombat) HPPercent() float64 {
	if ec.Participant.MaxHP == 0 {
		return 0
	}
	return float64(ec.Participant.HP) / float64(ec.Participant.MaxHP)
}

// ShouldFlee returns true if the enemy is below flee threshold.
func (ec *EnemyCombat) ShouldFlee() bool {
	return ec.HPPercent() < ec.Definition.FleeThreshold
}