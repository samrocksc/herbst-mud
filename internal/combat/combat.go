package combat

import (
	"github.com/sam/makeathing/internal/characters"
	"time"
)

// CombatRound represents a combat round
type CombatRound struct {
	Duration  time.Duration // 3 seconds
	StartTime time.Time
	Actions   []CombatAction
}

// CombatAction represents an action in combat
type CombatAction struct {
	Character *characters.Character
	Action    ActionType
	Target    *characters.Character
	Timestamp time.Time
}

// ActionType represents the type of combat action
type ActionType string

const (
	Attack ActionType = "attack"
	Spell  ActionType = "spell"
	Defend ActionType = "defend"
)

// CalculateHitsPerRound calculates how many hits a character can make per round
// based on their dexterity
func CalculateHitsPerRound(dexterity int) int {
	// High dexterity = 3 hits per 5 seconds
	// Low dexterity = 1 hit per 5 seconds
	// Linear scaling between 1 and 3
	if dexterity <= 8 {
		return 1
	} else if dexterity >= 17 {
		return 3
	} else {
		// Linear interpolation between 1 and 3 for values 9-16
		return 1 + (dexterity-8)/4
	}
}

// SpellCost represents the cost of casting a spell in terms of combat speed
const SpellCost = 1.0 / 3.0 // Spells take up 1/3 of the round's speed
