package combat

import (
	"fmt"
	"time"
)

// Combat represents a single combat encounter
type Combat struct {
	ID           CombatID
	RoomID       int
	Participants []*Participant
	State        CombatState
	TickNumber   int
	TickCountdown float64 // Seconds until next tick (for UI display)
	StartedAt    time.Time
	EndedAt      time.Time

	// Action queue for pending actions
	ActionQueue *ActionQueue

	// Effect registry for ongoing effects (DoT, HoT, buffs, debuffs)
	Effects *EffectRegistry

	// Tick channel for receiving tick events
	tickChan <-chan Tick

	// Combat log entries
	Log []CombatLogEntry
}

// CombatLogEntry represents a single entry in the combat log
type CombatLogEntry struct {
	Tick      int       `json:"tick"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	Type      string    `json:"type"` // "damage", "heal", "effect", "system", "miss", "info"
	SourceID  int       `json:"sourceId,omitempty"`
	TargetID  int       `json:"targetId,omitempty"`
	Value     int       `json:"value,omitempty"`
}

// Participant represents a character in combat
type Participant struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	IsPlayer     bool      `json:"isPlayer"`
	IsNPC        bool      `json:"isNPC"`
	UserID       int       `json:"userId,omitempty"` // For players
	NPCID        int       `json:"npcId,omitempty"`  // For NPCs

	// Stats
	HP           int       `json:"hp"`
	MaxHP        int       `json:"maxHP"`
	Mana         int       `json:"mana"`
	MaxMana      int       `json:"maxMana"`
	Stamina      int       `json:"stamina"`
	MaxStamina   int       `json:"maxStamina"`
	Level        int       `json:"level"`
	Class        string    `json:"class,omitempty"`

	// Combat stats
	Attack       int       `json:"attack"`
	Defense      int       `json:"defense"`
	Dexterity    int       `json:"dexterity"` // For initiative
	Initiative   int       `json:"initiative"` // Calculated at combat start

	// Status
	IsAlive      bool      `json:"isAlive"`
	IsActive     bool      `json:"isActive"` // In current turn order
	Team         int       `json:"team"` // 0 = player team, 1 = enemy team

	// Position in turn order (set during initiative roll)
	TurnPosition int       `json:"turnPosition"`

	// Current action being prepared (for UI display)
	CurrentAction *QueuedAction `json:"currentAction,omitempty"`

	// Active effects on this participant
	ActiveEffects []ActiveEffect `json:"activeEffects"`
}

// GetTeam returns the team number for this participant
func (p *Participant) GetTeam() int {
	return p.Team
}

// IsEnemy returns true if this participant is an enemy of the other
func (p *Participant) IsEnemy(other *Participant) bool {
	return p.Team != other.Team
}

// TakeDamage applies damage to the participant
func (p *Participant) TakeDamage(damage int) int {
	if damage < 0 {
		damage = 0
	}
	p.HP -= damage
	if p.HP < 0 {
		p.HP = 0
	}
	if p.HP == 0 {
		p.IsAlive = false
	}
	return damage
}

// Heal restores HP to the participant
func (p *Participant) Heal(amount int) int {
	if amount < 0 {
		amount = 0
	}
	oldHP := p.HP
	p.HP += amount
	if p.HP > p.MaxHP {
		p.HP = p.MaxHP
	}
	return p.HP - oldHP
}

// CanAct returns true if the participant can take actions
func (p *Participant) CanAct() bool {
	return p.IsAlive && p.HP > 0
}

// AddLogEntry adds an entry to the combat log
func (c *Combat) AddLogEntry(entry CombatLogEntry) {
	c.Log = append(c.Log, entry)
}

// FormatCombatLog is a helper for consistent log formatting
func FormatCombatLog(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

// GetParticipantByID finds a participant by their ID
func (c *Combat) GetParticipantByID(id int) *Participant {
	for _, p := range c.Participants {
		if p.ID == id {
			return p
		}
	}
	return nil
}

// GetAliveParticipants returns all alive participants
func (c *Combat) GetAliveParticipants() []*Participant {
	alive := make([]*Participant, 0)
	for _, p := range c.Participants {
		if p.IsAlive {
			alive = append(alive, p)
		}
	}
	return alive
}

// GetAliveByTeam returns alive participants on a specific team
func (c *Combat) GetAliveByTeam(team int) []*Participant {
	alive := make([]*Participant, 0)
	for _, p := range c.Participants {
		if p.IsAlive && p.Team == team {
			alive = append(alive, p)
		}
	}
	return alive
}

// AllEnemiesDefeated returns true if all enemies are defeated
func (c *Combat) AllEnemiesDefeated() bool {
	enemies := c.GetAliveByTeam(1)
	return len(enemies) == 0
}

// AllPlayersDefeated returns true if all players are defeated
func (c *Combat) AllPlayersDefeated() bool {
	players := c.GetAliveByTeam(0)
	return len(players) == 0
}