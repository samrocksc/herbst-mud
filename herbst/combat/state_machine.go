package combat

import (
	"context"
	"log"
	"time"
)

// CombatStateMachine manages the lifecycle of a combat encounter
type CombatStateMachine struct {
	combat        *Combat
	combatManager *CombatManager
	tickManager   *TickManager

	// State transition callbacks
	onStateEnter map[CombatState]func(*Combat)
	onStateExit  map[CombatState]func(*Combat)
	onTick       func(*Combat, Tick)

	// Input handling
	inputWindow     time.Duration
	inputDeadline   time.Time
	awaitingInput   map[int]bool // participantID -> awaiting
	selectedActions map[int]*QueuedAction // participantID -> action

	// Turn tracking
	currentTurn     int
	currentActorIdx int

	// State
	ctx    context.Context
	cancel context.CancelFunc
}

// NewCombatStateMachine creates a new state machine for a combat
func NewCombatStateMachine(combat *Combat, cm *CombatManager, tm *TickManager) *CombatStateMachine {
	return &CombatStateMachine{
		combat:          combat,
		combatManager:   cm,
		tickManager:     tm,
		onStateEnter:    make(map[CombatState]func(*Combat)),
		onStateExit:     make(map[CombatState]func(*Combat)),
		awaitingInput:   make(map[int]bool),
		selectedActions: make(map[int]*QueuedAction),
		inputWindow:     TickDuration,
	}
}

// SetOnStateEnter sets a callback for entering a state
func (csm *CombatStateMachine) SetOnStateEnter(state CombatState, callback func(*Combat)) {
	csm.onStateEnter[state] = callback
}

// SetOnStateExit sets a callback for exiting a state
func (csm *CombatStateMachine) SetOnStateExit(state CombatState, callback func(*Combat)) {
	csm.onStateExit[state] = callback
}

// SetOnTick sets a callback for tick events
func (csm *CombatStateMachine) SetOnTick(callback func(*Combat, Tick)) {
	csm.onTick = callback
}

// TransitionTo changes the combat state and fires callbacks
func (csm *CombatStateMachine) TransitionTo(newState CombatState) {
	oldState := csm.combat.State

	// Fire exit callback for old state
	if callback, exists := csm.onStateExit[oldState]; exists {
		callback(csm.combat)
	}

	// Update state
	csm.combat.State = newState

	// Log state transition
	log.Printf("[Combat %d] State transition: %s -> %s", csm.combat.ID, oldState, newState)
	csm.combat.AddLogEntry(CombatLogEntry{
		Tick:      csm.combat.TickNumber,
		Timestamp: time.Now(),
		Message:   string(newState),
		Type:      "system",
	})

	// Fire enter callback for new state
	if callback, exists := csm.onStateEnter[newState]; exists {
		callback(csm.combat)
	}

	// Notify combat manager
	if csm.combatManager != nil {
		csm.combatManager.ChangeCombatState(csm.combat, newState)
	}
}

// Start begins the combat state machine
func (csm *CombatStateMachine) Start(ctx context.Context) {
	csm.ctx, csm.cancel = context.WithCancel(ctx)

	// Transition to INIT state
	csm.TransitionTo(StateInit)

	// Roll initiative for all participants
	RollAllInitiative(csm.combat.Participants)
	csm.combat.AddLogEntry(CombatLogEntry{
		Tick:      0,
		Timestamp: time.Now(),
		Message:   "Initiative rolled",
		Type:      "system",
	})

	// Sort by initiative
	order := csm.combat.GetTurnOrder()
	for i, p := range order {
		csm.combat.AddLogEntry(CombatLogEntry{
			Tick:      0,
			Timestamp: time.Now(),
			Message:   FormatCombatLog("%s rolls %d initiative", p.Name, p.Initiative),
			Type:      "info",
		})
		_ = i // position in turn order
	}

	// Start input window for first turn
	csm.TransitionTo(StateActive)
	csm.startInputWindow()

	// Start processing ticks
	go csm.processTicks()
}

// Stop halts the state machine
func (csm *CombatStateMachine) Stop() {
	if csm.cancel != nil {
		csm.cancel()
	}
}

// startInputWindow begins the input window for the current actor
func (csm *CombatStateMachine) startInputWindow() {
	// Get current actor
	order := csm.combat.GetTurnOrder()
	if len(order) == 0 {
		return
	}

	// Find next alive actor
	for i := 0; i < len(order); i++ {
		idx := (csm.currentActorIdx + i) % len(order)
		actor := order[idx]
		if actor.IsAlive && actor.CanAct() {
			csm.currentActorIdx = idx
			csm.awaitingInput[actor.ID] = true
			csm.inputDeadline = time.Now().Add(csm.inputWindow)
			csm.combat.TickCountdown = float64(csm.inputWindow) / float64(time.Second)
			csm.combat.AddLogEntry(CombatLogEntry{
				Tick:      csm.combat.TickNumber,
				Timestamp: time.Now(),
				Message:   FormatCombatLog("Waiting for %s's action", actor.Name),
				Type:      "system",
			})
			return
		}
	}

	// No alive actors who can act
	csm.checkEndCondition()
}

// processTicks handles the tick loop
func (csm *CombatStateMachine) processTicks() {
	if csm.tickManager == nil {
		return
	}

	tickChan := csm.tickManager.Subscribe(int(csm.combat.ID))
	defer csm.tickManager.Unsubscribe(int(csm.combat.ID))

	for {
		select {
		case <-csm.ctx.Done():
			return
		case tick := <-tickChan:
			csm.processTick(tick)
		}
	}
}

// processTick handles a single combat tick
func (csm *CombatStateMachine) processTick(tick Tick) {
	// Skip if combat ended
	if csm.combat.State == StateEnded {
		return
	}

	// Fire tick callback
	if csm.onTick != nil {
		csm.onTick(csm.combat, tick)
	}

	// Increment tick counter
	csm.combat.TickNumber = tick.ID

	// Process pending actions
	csm.processPendingActions()

	// Process effects (DoT, HoT, buffs/debuffs)
	csm.processEffects()

	// Check end condition
	if csm.checkEndCondition() {
		return
	}

	// Advance to next turn
	csm.advanceTurn()
}

// processPendingActions executes queued actions
func (csm *CombatStateMachine) processPendingActions() {
	// Get actions for this tick
	actions := csm.combat.ActionQueue.GetActionsForTick(csm.combat.TickNumber)

	for _, qa := range actions {
		// Skip if source is dead
		if !qa.Source.IsAlive {
			continue
		}

		// Execute the action
		csm.executeAction(qa)

		// Remove from queue
		csm.combat.ActionQueue.Remove(qa)
	}
}

// executeAction executes a queued action
func (csm *CombatStateMachine) executeAction(qa *QueuedAction) {
	action := qa.Action
	source := qa.Source
	target := qa.Target

	// Log action
	csm.combat.AddLogEntry(CombatLogEntry{
		Tick:      csm.combat.TickNumber,
		Timestamp: time.Now(),
		SourceID:  source.ID,
		TargetID:  target.ID,
		Type:      "action",
		Message:   FormatCombatLog("%s uses %s", source.Name, action.Name),
	})

	// Apply damage/healing
	if action.BaseDamage > 0 && target != nil {
		damage := action.BaseDamage
		actualDamage := target.TakeDamage(damage)
		csm.combat.AddLogEntry(CombatLogEntry{
			Tick:      csm.combat.TickNumber,
			Timestamp: time.Now(),
			SourceID:  source.ID,
			TargetID:  target.ID,
			Type:      "damage",
			Value:     actualDamage,
			Message:   FormatCombatLog("%s takes %d damage", target.Name, actualDamage),
		})
	}

	if action.BaseHeal > 0 && target != nil {
		heal := action.BaseHeal
		actualHeal := target.Heal(heal)
		csm.combat.AddLogEntry(CombatLogEntry{
			Tick:      csm.combat.TickNumber,
			Timestamp: time.Now(),
			SourceID:  source.ID,
			TargetID:  target.ID,
			Type:      "heal",
			Value:     actualHeal,
			Message:   FormatCombatLog("%s heals for %d", target.Name, actualHeal),
		})
	}
}

// processEffects handles ongoing effects (DoT, HoT, etc.)
func (csm *CombatStateMachine) processEffects() {
	// Process all effects
	damageEffects, healEffects := csm.combat.Effects.ProcessEffects(csm.combat.Participants)

	// Log DoT effects
	for _, e := range damageEffects {
		csm.combat.AddLogEntry(CombatLogEntry{
			Tick:      csm.combat.TickNumber,
			Timestamp: time.Now(),
			TargetID:  e.Target.ID,
			Type:      "damage",
			Value:     e.Value,
			Message:   FormatCombatLog("%s takes %d damage from %s", e.Target.Name, e.Value, e.Effect.Name),
		})
	}

	// Log HoT effects
	for _, e := range healEffects {
		csm.combat.AddLogEntry(CombatLogEntry{
			Tick:      csm.combat.TickNumber,
			Timestamp: time.Now(),
			TargetID:  e.Target.ID,
			Type:      "heal",
			Value:     e.Value,
			Message:   FormatCombatLog("%s heals for %d from %s", e.Target.Name, e.Value, e.Effect.Name),
		})
	}
}

// checkEndCondition checks if combat should end
func (csm *CombatStateMachine) checkEndCondition() bool {
	if csm.combat.AllEnemiesDefeated() {
		csm.endCombat("victory", "All enemies defeated!")
		return true
	}

	if csm.combat.AllPlayersDefeated() {
		csm.endCombat("defeat", "All players defeated!")
		return true
	}

	return false
}

// endCombat transitions to ended state
func (csm *CombatStateMachine) endCombat(result, reason string) {
	csm.combat.EndedAt = time.Now()
	csm.combat.AddLogEntry(CombatLogEntry{
		Tick:      csm.combat.TickNumber,
		Timestamp: time.Now(),
		Type:      "system",
		Message:   FormatCombatLog("Combat ended: %s - %s", result, reason),
	})

	// Handle weapon drops on NPC defeat (victory)
	if result == "victory" {
		csm.handleWeaponDrops()
	}

	csm.TransitionTo(StateEnded)

	if csm.combatManager != nil {
		csm.combatManager.EndCombat(csm.combat.ID, reason)
	}
}

// handleWeaponDrops processes weapon drops when enemies are defeated
func (csm *CombatStateMachine) handleWeaponDrops() {
	// Get all defeated NPCs in this combat
	for _, participant := range csm.combat.Participants {
		if participant.IsNPC && !participant.IsAlive {
			// NPC was defeated - check if it has a guaranteed drop
			csm.combat.AddLogEntry(CombatLogEntry{
				Tick:      csm.combat.TickNumber,
				Timestamp: time.Now(),
				Type:      "loot",
				TargetID:  participant.ID,
				Message:   FormatCombatLog("%s has been defeated!", participant.Name),
			})

			// Log that a weapon can be picked up (actual pickup happens via game command)
			csm.combat.AddLogEntry(CombatLogEntry{
				Tick:      csm.combat.TickNumber,
				Timestamp: time.Now(),
				Type:      "loot",
				Message:   FormatCombatLog("You see weapons among the remains..."),
			})
		}
	}

	// Log weapon availability to players
	for _, participant := range csm.combat.Participants {
		if participant.IsPlayer && participant.IsAlive {
			csm.combat.AddLogEntry(CombatLogEntry{
				Tick:      csm.combat.TickNumber,
				Timestamp: time.Now(),
				Type:      "loot",
				TargetID:  participant.ID,
				Message:   FormatCombatLog("Victory! Use 'pickup' to claim your weapon from the fallen!"),
			})
		}
	}
}

// advanceTurn moves to the next actor's turn
func (csm *CombatStateMachine) advanceTurn() {
	order := csm.combat.GetTurnOrder()
	if len(order) == 0 {
		return
	}

	// Move to next alive actor
	startIdx := csm.currentActorIdx
	for {
		csm.currentActorIdx = (csm.currentActorIdx + 1) % len(order)
		if csm.currentActorIdx == 0 {
			// Completed a full round
			csm.currentTurn++
		}

		actor := order[csm.currentActorIdx]
		if actor.IsAlive && actor.CanAct() {
			csm.awaitingInput[actor.ID] = true
			csm.inputDeadline = time.Now().Add(csm.inputWindow)
			csm.combat.TickCountdown = float64(csm.inputWindow) / float64(time.Second)
			return
		}

		// Wrapped around, no one can act
		if csm.currentActorIdx == startIdx {
			csm.checkEndCondition()
			return
		}
	}
}

// SubmitAction allows a participant to submit their action
func (csm *CombatStateMachine) SubmitAction(participantID int, action *ActionDefinition, target *Participant) bool {
	// Check if this participant can act
	if !csm.awaitingInput[participantID] {
		return false
	}

	// Find the participant
	source := csm.combat.GetParticipantByID(participantID)
	if source == nil {
		return false
	}

	// Create queued action
	qa := csm.combat.ActionQueue.Queue(action, source, target, csm.combat.TickNumber+action.TickCost)
	csm.selectedActions[participantID] = qa

	// Clear input awaiting
	delete(csm.awaitingInput, participantID)

	// Log the selection
	csm.combat.AddLogEntry(CombatLogEntry{
		Tick:      csm.combat.TickNumber,
		Timestamp: time.Now(),
		SourceID:  source.ID,
		TargetID:  target.ID,
		Type:      "info",
		Message:   FormatCombatLog("%s prepares %s", source.Name, action.Name),
	})

	return true
}

// GetCurrentActor returns the participant whose turn it is
func (csm *CombatStateMachine) GetCurrentActor() *Participant {
	order := csm.combat.GetTurnOrder()
	if len(order) == 0 || csm.currentActorIdx >= len(order) {
		return nil
	}
	return order[csm.currentActorIdx]
}

// GetInputDeadline returns when the input window closes
func (csm *CombatStateMachine) GetInputDeadline() time.Time {
	return csm.inputDeadline
}

// IsAwaitingInput returns true if a participant needs to input an action
func (csm *CombatStateMachine) IsAwaitingInput(participantID int) bool {
	return csm.awaitingInput[participantID]
}

// GetTickCountdown returns seconds until next tick
func (csm *CombatStateMachine) GetTickCountdown() float64 {
	remaining := time.Until(csm.inputDeadline).Seconds()
	if remaining < 0 {
		return 0
	}
	return remaining
}