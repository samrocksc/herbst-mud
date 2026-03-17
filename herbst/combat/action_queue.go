package combat

import (
	"sort"
	"sync"
)

// ActionPriority defines when an action executes in the tick
type ActionPriority int

const (
	// PriorityImmediate executes first (instant reactions)
	PriorityImmediate ActionPriority = iota
	// PriorityFast executes early (quick actions like parry)
	PriorityFast
	// PriorityNormal executes in normal order (most actions)
	PriorityNormal
	// PrioritySlow executes late (heavy attacks, spells)
	PrioritySlow
	// PriorityLast executes last (delayed effects)
	PriorityLast
)

// ActionType defines how an action executes
type ActionType string

const (
	// ActionInstant executes immediately in the current tick
	ActionInstant ActionType = "INSTANT"
	// ActionChannel requires multiple ticks to complete
	ActionChannel ActionType = "CHANNEL"
	// ActionCharge requires buildup before execution
	ActionCharge ActionType = "CHARGE"
)

// Action represents a combat action that can be performed
type Action struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Type         ActionType    `json:"type"`
	Priority     ActionPriority `json:"priority"`
	TickCost     int           `json:"tickCost"`      // Ticks required to execute
	ChannelTicks int           `json:"channelTicks"`  // For channeling: total ticks
	ChargeTicks  int           `json:"chargeTicks"`   // For charging: ticks to charge

	// Base damage/healing
	BaseDamage   int           `json:"baseDamage"`
	BaseHeal     int           `json:"baseHeal"`

	// Requirements
	ManaCost     int           `json:"manaCost"`
	StaminaCost  int           `json:"staminaCost"`
	Cooldown     int           `json:"cooldown"`      // Ticks before can use again

	// Flags
	IsAttack     bool          `json:"isAttack"`
	IsDefend     bool           `json:"isDefend"`
	IsHeal       bool           `json:"isHeal"`
	IsFlee       bool           `json:"isFlee"`
	IsWait       bool           `json:"isWait"`
}

// QueuedAction represents an action queued for a specific tick
type QueuedAction struct {
	Action      *Action     `json:"action"`
	Source      *Participant `json:"source"`
	Target      *Participant `json:"target"`
	TickNumber  int         `json:"tickNumber"`  // Which tick this executes on
	Priority    ActionPriority `json:"priority"`
	QueuedAt    int         `json:"queuedAt"`    // Tick when this was queued
	ChannelTick int         `json:"channelTick"` // Current channel progress
	ChargeTick   int         `json:"chargeTick"`  // Current charge progress
}

// ActionQueue manages pending combat actions
type ActionQueue struct {
	mu     sync.RWMutex
	actions []*QueuedAction
}

// NewActionQueue creates a new action queue
func NewActionQueue() *ActionQueue {
	return &ActionQueue{
		actions: make([]*QueuedAction, 0),
	}
}

// Queue adds an action to the queue
func (aq *ActionQueue) Queue(action *Action, source, target *Participant, tickNumber int) *QueuedAction {
	aq.mu.Lock()
	defer aq.mu.Unlock()

	qa := &QueuedAction{
		Action:     action,
		Source:     source,
		Target:     target,
		TickNumber: tickNumber,
		Priority:   action.Priority,
		QueuedAt:   tickNumber,
	}

	// For channel actions, set the current channel tick
	if action.Type == ActionChannel {
		qa.ChannelTick = 1
	}
	// For charge actions, set the current charge tick
	if action.Type == ActionCharge {
		qa.ChargeTick = 1
	}

	aq.actions = append(aq.actions, qa)
	return qa
}

// QueueImmediate queues an action for immediate execution (current tick)
func (aq *ActionQueue) QueueImmediate(action *Action, source, target *Participant, currentTick int) *QueuedAction {
	return aq.Queue(action, source, target, currentTick)
}

// QueueForNextTick queues an action for the next tick
func (aq *ActionQueue) QueueForNextTick(action *Action, source, target *Participant, currentTick int) *QueuedAction {
	return aq.Queue(action, source, target, currentTick+1)
}

// GetActionsForTick returns all actions that should execute on a given tick, sorted by priority
func (aq *ActionQueue) GetActionsForTick(tick int) []*QueuedAction {
	aq.mu.RLock()
	defer aq.mu.RUnlock()

	var result []*QueuedAction
	for _, qa := range aq.actions {
		if qa.TickNumber == tick {
			result = append(result, qa)
		}
	}

	// Sort by priority (lower = earlier execution)
	sort.Slice(result, func(i, j int) bool {
		if result[i].Priority != result[j].Priority {
			return result[i].Priority < result[j].Priority
		}
		// Same priority: higher initiative acts first
		return result[i].Source.Initiative > result[j].Source.Initiative
	})

	return result
}

// Remove removes a queued action
func (aq *ActionQueue) Remove(qa *QueuedAction) {
	aq.mu.Lock()
	defer aq.mu.Unlock()

	for i, action := range aq.actions {
		if action == qa {
			aq.actions = append(aq.actions[:i], aq.actions[i+1:]...)
			break
		}
	}
}

// RemoveByID removes a queued action by source ID
func (aq *ActionQueue) RemoveByID(sourceID int) {
	aq.mu.Lock()
	defer aq.mu.Unlock()

	newActions := make([]*QueuedAction, 0)
	for _, qa := range aq.actions {
		if qa.Source.ID != sourceID {
			newActions = append(newActions, qa)
		}
	}
	aq.actions = newActions
}

// Clear removes all actions from the queue
func (aq *ActionQueue) Clear() {
	aq.mu.Lock()
	defer aq.mu.Unlock()
	aq.actions = make([]*QueuedAction, 0)
}

// GetPendingActionsForParticipant returns all pending actions for a specific participant
func (aq *ActionQueue) GetPendingActionsForParticipant(participantID int) []*QueuedAction {
	aq.mu.RLock()
	defer aq.mu.RUnlock()

	var result []*QueuedAction
	for _, qa := range aq.actions {
		if qa.Source.ID == participantID {
			result = append(result, qa)
		}
	}
	return result
}

// Count returns the number of actions in the queue
func (aq *ActionQueue) Count() int {
	aq.mu.RLock()
	defer aq.mu.RUnlock()
	return len(aq.actions)
}