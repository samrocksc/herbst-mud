package combat

import (
	"math/rand"
	"sort"
)

// InitiativeRoll represents an initiative roll result
type InitiativeRoll struct {
	Participant *Participant
	Roll        int // Random 1-20
	Dexterity   int // DEX stat
	Total       int // Roll + DEX (used for tie-breaking)
}

// RollInitiative calculates initiative for a participant
// Initiative = DEX + random(1,20)
func RollInitiative(p *Participant) int {
	roll := rand.Intn(20) + 1 // 1-20
	p.Initiative = p.Dexterity + roll
	return p.Initiative
}

// RollAllInitiative calculates initiative for all participants in a combat
func RollAllInitiative(participants []*Participant) {
	for _, p := range participants {
		RollInitiative(p)
	}
}

// SortByInitiative sorts participants by initiative (descending)
// Ties are broken by DEX, then by random roll
func SortByInitiative(participants []*Participant) {
	sort.Slice(participants, func(i, j int) bool {
		// Higher initiative goes first
		if participants[i].Initiative != participants[j].Initiative {
			return participants[i].Initiative > participants[j].Initiative
		}
		// Tie-breaker: higher DEX goes first
		if participants[i].Dexterity != participants[j].Dexterity {
			return participants[i].Dexterity > participants[j].Dexterity
		}
		// Final tie-breaker: random
		return rand.Intn(2) == 0
	})
}

// GetTurnOrder returns the initiative order as a slice of participants
func (c *Combat) GetTurnOrder() []*Participant {
	// Get all alive participants
	alive := c.GetAliveParticipants()
	
	// Copy and sort by initiative
	order := make([]*Participant, len(alive))
	copy(order, alive)
	SortByInitiative(order)
	
	// Set turn positions
	for i, p := range order {
		p.TurnPosition = i + 1
	}
	
	return order
}

// GetTurnOrderString returns a formatted string of turn order for display
func (c *Combat) GetTurnOrderString() string {
	order := c.GetTurnOrder()
	if len(order) == 0 {
		return "No combatants"
	}
	
	result := "Turn Order:\n"
	for i, p := range order {
		teamLabel := "Ally"
		if p.Team == 1 {
			teamLabel = "Enemy"
		}
		status := ""
		if !p.IsAlive {
			status = " [DEAD]"
		}
		result += FormatCombatLog("  %d. %s (%s) - Initiative: %d%s\n",
			i+1, p.Name, teamLabel, p.Initiative, status)
	}
	return result
}

// GetCurrentActor returns the participant whose turn it is
// This assumes round-robin turn order based on initiative
func (c *Combat) GetCurrentActor(turnIndex int) *Participant {
	order := c.GetTurnOrder()
	if len(order) == 0 {
		return nil
	}
	// Wrap around for round-robin
	return order[turnIndex % len(order)]
}

// GetNextActor returns the next participant in turn order
func (c *Combat) GetNextActor(currentIndex int) (nextIndex int, actor *Participant) {
	order := c.GetTurnOrder()
	if len(order) == 0 {
		return 0, nil
	}
	nextIndex = (currentIndex + 1) % len(order)
	return nextIndex, order[nextIndex]
}

// InitiativeRollDetailed returns detailed initiative info for display
type InitiativeRollDetailed struct {
	Name        string `json:"name"`
	IsPlayer    bool   `json:"isPlayer"`
	Dexterity   int    `json:"dexterity"`
	Roll        int    `json:"roll"`
	Total       int    `json:"total"`
	TurnOrder   int    `json:"turnOrder"`
}

// GetInitiativeRolls returns detailed initiative information for all participants
func (c *Combat) GetInitiativeRolls() []InitiativeRollDetailed {
	order := c.GetTurnOrder()
	rolls := make([]InitiativeRollDetailed, len(order))
	
	for i, p := range order {
		rolls[i] = InitiativeRollDetailed{
			Name:      p.Name,
			IsPlayer:  p.IsPlayer,
			Dexterity: p.Dexterity,
			Roll:      p.Initiative - p.Dexterity, // The dice roll portion
			Total:     p.Initiative,
			TurnOrder: i + 1,
		}
	}
	
	return rolls
}