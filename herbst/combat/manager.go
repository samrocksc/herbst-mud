package combat

import (
	"context"
	"log"
	"sync"
)

// CombatID is a unique identifier for a combat encounter
type CombatID int

// CombatState represents the current state of a combat encounter
type CombatState string

const (
	// StateIdle means no combat is active
	StateIdle CombatState = "IDLE"
	// StateInit means combat is initializing
	StateInit CombatState = "COMBAT_INIT"
	// StateActive means combat is ongoing
	StateActive CombatState = "COMBAT_ACTIVE"
	// StateEnded means combat has concluded
	StateEnded CombatState = "COMBAT_END"
)

// CombatManager tracks and manages all active combat encounters
type CombatManager struct {
	mu      sync.RWMutex
	ctx     context.Context
	cancel  context.CancelFunc
	running bool

	// Active combats indexed by ID
	combats map[CombatID]*Combat

	// Next combat ID for auto-increment
	nextID CombatID

	// Tick manager reference
	tickManager *TickManager

	// Callbacks
	onCombatStart   func(combat *Combat)
	onCombatEnd     func(combat *Combat, reason string)
	onCombatStateChange func(combat *Combat, oldState, newState CombatState)
}

// NewCombatManager creates a new combat manager
func NewCombatManager() *CombatManager {
	return &CombatManager{
		combats: make(map[CombatID]*Combat),
		nextID:  1,
	}
}

// SetTickManager sets the tick manager for combat timing
func (cm *CombatManager) SetTickManager(tm *TickManager) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.tickManager = tm
}

// Start begins the combat manager's main loop
func (cm *CombatManager) Start(ctx context.Context) {
	cm.mu.Lock()
	cm.ctx, cm.cancel = context.WithCancel(ctx)
	cm.running = true
	cm.mu.Unlock()

	log.Println("[CombatManager] Started")

	// Start the tick manager if available
	if cm.tickManager != nil {
		cm.tickManager.Start(ctx)
	}
}

// Stop halts the combat manager
func (cm *CombatManager) Stop() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.cancel != nil {
		cm.cancel()
	}
	cm.running = false

	// End all active combats
	for id, combat := range cm.combats {
		combat.State = StateEnded
		delete(cm.combats, id)
	}

	// Stop tick manager
	if cm.tickManager != nil {
		cm.tickManager.Stop()
	}

	log.Println("[CombatManager] Stopped")
}

// CreateCombat creates a new combat encounter
// Returns the combat ID
func (cm *CombatManager) CreateCombat(roomID int, participants []*Participant) CombatID {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	id := cm.nextID
	cm.nextID++

	combat := &Combat{
		ID:           id,
		RoomID:       roomID,
		Participants: participants,
		State:        StateInit,
		TickNumber:   0,
		ActionQueue:  NewActionQueue(),
		Effects:      NewEffectRegistry(),
	}

	cm.combats[id] = combat

	// Subscribe to ticks
	if cm.tickManager != nil {
		combat.tickChan = cm.tickManager.Subscribe(int(id))
	}

	log.Printf("[CombatManager] Created combat %d in room %d with %d participants",
		id, roomID, len(participants))

	// Fire callback
	if cm.onCombatStart != nil {
		cm.onCombatStart(combat)
	}

	return id
}

// GetCombat retrieves a combat by ID
func (cm *CombatManager) GetCombat(id CombatID) (*Combat, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	combat, exists := cm.combats[id]
	return combat, exists
}

// EndCombat ends a combat encounter
func (cm *CombatManager) EndCombat(id CombatID, reason string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	combat, exists := cm.combats[id]
	if !exists {
		return
	}

	combat.State = StateEnded

	// Unsubscribe from ticks
	if cm.tickManager != nil {
		cm.tickManager.Unsubscribe(int(id))
	}

	// Fire callback
	if cm.onCombatEnd != nil {
		cm.onCombatEnd(combat, reason)
	}

	delete(cm.combats, id)
	log.Printf("[CombatManager] Ended combat %d: %s", id, reason)
}

// GetActiveCombats returns all active combats
func (cm *CombatManager) GetActiveCombats() []*Combat {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	combats := make([]*Combat, 0, len(cm.combats))
	for _, combat := range cm.combats {
		combats = append(combats, combat)
	}
	return combats
}

// GetCombatsByRoom returns all combats in a specific room
func (cm *CombatManager) GetCombatsByRoom(roomID int) []*Combat {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	combats := make([]*Combat, 0)
	for _, combat := range cm.combats {
		if combat.RoomID == roomID {
			combats = append(combats, combat)
		}
	}
	return combats
}

// GetCombatCount returns the number of active combats
func (cm *CombatManager) GetCombatCount() int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return len(cm.combats)
}

// SetOnCombatStart sets the callback for combat start events
func (cm *CombatManager) SetOnCombatStart(callback func(combat *Combat)) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.onCombatStart = callback
}

// SetOnCombatEnd sets the callback for combat end events
func (cm *CombatManager) SetOnCombatEnd(callback func(combat *Combat, reason string)) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.onCombatEnd = callback
}

// SetOnCombatStateChange sets the callback for combat state changes
func (cm *CombatManager) SetOnCombatStateChange(callback func(combat *Combat, oldState, newState CombatState)) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.onCombatStateChange = callback
}

// ChangeCombatState changes a combat's state and fires callbacks
func (cm *CombatManager) ChangeCombatState(combat *Combat, newState CombatState) {
	cm.mu.Lock()
	oldState := combat.State
	combat.State = newState
	callback := cm.onCombatStateChange
	cm.mu.Unlock()

	if callback != nil {
		callback(combat, oldState, newState)
	}
}

// IsRunning returns whether the combat manager is running
func (cm *CombatManager) IsRunning() bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.running
}