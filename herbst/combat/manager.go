package combat

import (
	"fmt"
	"sync"
)

// Combat represents an active combat encounter.
type Combat struct {
	ID        int
	StartTick int
	mu        sync.RWMutex
	// participants are the combatants in this combat
	participants map[int]*Participant
	// actionQueue holds pending actions for this combat
	actionQueue map[int][]*Action
	// effects holds active effects on combatants
	effects map[int][]*Effect
}

// Participant represents a combatant in a combat.
type Participant struct {
	ID       int
	Name     string
	HP       int
	MaxHP    int
	Ready    bool // ready for this tick's input phase
	Input    *Action
	Cooldowns map[string]int // skill cooldowns in ticks
	mu       sync.RWMutex
}

// Action represents a combat action taken by a participant.
type Action struct {
	ID          int
	Participant int
	Type        ActionType
	Target      int
	Priority    int // lower = higher priority (1 is highest)
	Payload     interface{}
}

// ActionType defines the type of combat action.
type ActionType string

const (
	ActionAttack    ActionType = "attack"
	ActionDefend    ActionType = "defend"
	ActionSkill     ActionType = "skill"
	ActionItem      ActionType = "item"
	ActionFlee      ActionType = "flee"
	ActionWait      ActionType = "wait"
)

// Effect represents a status effect on a participant.
type Effect struct {
	ID        int
	Name      string
	Duration  int // remaining ticks
	OnApply   func(*Participant)
	OnTick    func(*Participant)
	OnExpire  func(*Participant)
}

// CombatManager manages all active combats.
type CombatManager struct {
	mu           sync.RWMutex
	combats      map[int]*Combat
	nextCombatID int
	tickLoop     *TickLoop
}

// NewCombatManager creates a new CombatManager.
func NewCombatManager(tickLoop *TickLoop) *CombatManager {
	return &CombatManager{
		combats:      make(map[int]*Combat),
		nextCombatID: 1,
		tickLoop:     tickLoop,
	}
}

// CreateCombat creates a new combat and returns its ID.
func (m *CombatManager) CreateCombat(startTick int) int {
	m.mu.Lock()
	defer m.mu.Unlock()

	combat := &Combat{
		ID:           m.nextCombatID,
		StartTick:    startTick,
		participants: make(map[int]*Participant),
		actionQueue:  make(map[int][]*Action),
		effects:      make(map[int][]*Effect),
	}
	m.combats[combat.ID] = combat
	m.nextCombatID++

	// Register with tick loop if available
	if m.tickLoop != nil {
		m.tickLoop.RegisterCombat(combat)
	}

	return combat.ID
}

// GetCombat retrieves a combat by ID.
func (m *CombatManager) GetCombat(id int) (*Combat, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, ok := m.combats[id]
	return c, ok
}

// EndCombat removes a combat from the manager.
func (m *CombatManager) EndCombat(id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	c, ok := m.combats[id]
	if !ok {
		return fmt.Errorf("combat %d not found", id)
	}

	// Unregister from tick loop
	if m.tickLoop != nil {
		m.tickLoop.UnregisterCombat(id)
	}

	delete(m.combats, id)
	_ = c // silence unused variable warning
	return nil
}

// AddParticipant adds a participant to a combat.
func (m *CombatManager) AddParticipant(combatID int, p *Participant) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	c, ok := m.combats[combatID]
	if !ok {
		return fmt.Errorf("combat %d not found", combatID)
	}

	c.participants[p.ID] = p
	c.actionQueue[p.ID] = nil // initialize empty action queue
	c.effects[p.ID] = nil     // initialize empty effects list
	return nil
}

// RemoveParticipant removes a participant from a combat.
func (m *CombatManager) RemoveParticipant(combatID, participantID int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	c, ok := m.combats[combatID]
	if !ok {
		return fmt.Errorf("combat %d not found", combatID)
	}

	delete(c.participants, participantID)
	delete(c.actionQueue, participantID)
	delete(c.effects, participantID)
	return nil
}

// QueueAction adds an action to a participant's queue.
func (c *Combat) QueueAction(action *Action) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.actionQueue[action.Participant] = append(c.actionQueue[action.Participant], action)
}

// GetActions returns the queued actions for a participant.
func (c *Combat) GetActions(participantID int) []*Action {
	c.mu.RLock()
	defer c.mu.RUnlock()
	actions := c.actionQueue[participantID]
	// Return a copy to prevent mutation
	result := make([]*Action, len(actions))
	copy(result, actions)
	return result
}

// ClearActions clears all queued actions for a participant.
func (c *Combat) ClearActions(participantID int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.actionQueue[participantID] = nil
}

// ProcessTick processes a single tick for this combat.
func (c *Combat) ProcessTick(tick int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Phase 1: Input phase - participants submit their actions
	// (handled by SetInput on Participant)

	// Phase 2: Resolve phase - execute actions in priority order
	c.resolveActions()

	// Phase 3: Apply effects
	c.processEffects()

	// Phase 4: Check for combat end
	c.checkCombatEnd()
}

func (c *Combat) resolveActions() {
	// Collect all actions and sort by priority
	type queuedAction struct {
		participant int
		action      *Action
	}

	var allActions []queuedAction
	for participantID, actions := range c.actionQueue {
		for _, action := range actions {
			allActions = append(allActions, queuedAction{
				participant: participantID,
				action:      action,
			})
		}
	}

	// Sort by priority (lower = higher priority)
	// Note: This is a simple sort; could be enhanced for ties
	// Using insertion sort for small lists is efficient
	for i := 1; i < len(allActions); i++ {
		j := i
		for j > 0 && allActions[j].action.Priority < allActions[j-1].action.Priority {
			allActions[j], allActions[j-1] = allActions[j-1], allActions[j]
			j--
		}
	}

	// Execute actions in order
	for _, qa := range allActions {
		c.executeAction(qa.participant, qa.action)
	}

	// Clear action queues after resolution
	for participantID := range c.actionQueue {
		c.actionQueue[participantID] = nil
	}
}

func (c *Combat) executeAction(participantID int, action *Action) {
	participant, ok := c.participants[participantID]
	if !ok {
		return
	}

	switch action.Type {
	case ActionAttack:
		// Attack logic - target takes damage
		_ = participant // Attack execution would modify target
		_ = action      // Could contain damage amount, weapon, etc.
	case ActionDefend:
		// Defend logic - grant defense bonus this tick
		participant.Ready = true
	case ActionSkill:
		// Skill execution
		_ = action
	case ActionItem:
		// Item usage
		_ = action
	case ActionFlee:
		// Flee attempt
		_ = action
	case ActionWait:
		// Do nothing this tick
		participant.Ready = true
	}
}

func (c *Combat) processEffects() {
	for participantID, effects := range c.effects {
		participant, ok := c.participants[participantID]
		if !ok {
			continue
		}

		for _, effect := range effects {
			if effect.OnTick != nil {
				effect.OnTick(participant)
			}
		}
	}
}

func (c *Combat) checkCombatEnd() {
	alive := 0
	for _, p := range c.participants {
		if p.HP > 0 {
			alive++
		}
	}
	// Combat ends when only one participant remains or all flee
	_ = alive // Would trigger combat end logic
}

// SetInput sets a participant's action for the current tick.
func (p *Participant) SetInput(action *Action) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Input = action
	p.Ready = true
}

// AddEffect adds a status effect to a participant.
func (c *Combat) AddEffect(participantID int, effect *Effect) {
	c.mu.Lock()
	defer c.mu.Unlock()

	effects := c.effects[participantID]
	effects = append(effects, effect)
	c.effects[participantID] = effects

	if effect.OnApply != nil {
		if participant, ok := c.participants[participantID]; ok {
			effect.OnApply(participant)
		}
	}
}

// RemoveEffect removes a status effect from a participant.
func (c *Combat) RemoveEffect(participantID, effectID int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	effects := c.effects[participantID]
	for i, e := range effects {
		if e.ID == effectID {
			// Remove from slice
			effects = append(effects[:i], effects[i+1:]...)
			c.effects[participantID] = effects

			if e.OnExpire != nil {
				if participant, ok := c.participants[participantID]; ok {
					e.OnExpire(participant)
				}
			}
			return
		}
	}
}

// GetEffects returns all effects on a participant.
func (c *Combat) GetEffects(participantID int) []*Effect {
	c.mu.RLock()
	defer c.mu.RUnlock()
	effects := c.effects[participantID]
	result := make([]*Effect, len(effects))
	copy(result, effects)
	return result
}

// GetParticipants returns all participants in the combat.
func (c *Combat) GetParticipants() []*Participant {
	c.mu.RLock()
	defer c.mu.RUnlock()

	participants := make([]*Participant, 0, len(c.participants))
	for _, p := range c.participants {
		participants = append(participants, p)
	}
	return participants
}

// GetParticipant returns a specific participant.
func (c *Combat) GetParticipant(id int) (*Participant, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	p, ok := c.participants[id]
	return p, ok
}

// NewParticipant creates a new participant with the given parameters.
func NewParticipant(id int, name string, hp int) *Participant {
	return &Participant{
		ID:          id,
		Name:        name,
		HP:          hp,
		MaxHP:       hp,
		Ready:       false,
		Input:       nil,
		Cooldowns:   make(map[string]int),
	}
}

// NewAction creates a new action.
func NewAction(id, participant, target int, actionType ActionType, priority int) *Action {
	return &Action{
		ID:          id,
		Participant: participant,
		Type:        actionType,
		Target:      target,
		Priority:    priority,
		Payload:     nil,
	}
}

// NewEffect creates a new status effect.
func NewEffect(id int, name string, duration int) *Effect {
	return &Effect{
		ID:       id,
		Name:     name,
		Duration: duration,
	}
}