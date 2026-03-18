package combat

import (
	"sync"
	"time"
)

// InputMode defines how the input system works
type InputMode string

const (
	// InputModeRealTime requires input within the tick window
	InputModeRealTime InputMode = "REAL_TIME"
	// InputModeTurnBased waits for input before advancing
	InputModeTurnBased InputMode = "TURN_BASED"
	// InputModeHybrid allows real-time input with auto-actions on timeout
	InputModeHybrid InputMode = "HYBRID"
)

// InputConfig holds input system configuration
type InputConfig struct {
	// Input window duration (default 1.5s)
	InputWindow time.Duration
	
	// Mode for input handling
	Mode InputMode
	
	// Auto-attack when no input received
	AutoAttackEnabled bool
	
	// Auto-defend when HP below threshold
	AutoDefendEnabled bool
	AutoDefendHPTreshold float64 // 0.0-1.0 (percentage)
	
	// Allow late input queuing
	LateInputQueueing bool
	
	// Number of talent slots
	TalentSlots int
}

// DefaultInputConfig returns the default input configuration
func DefaultInputConfig() InputConfig {
	return InputConfig{
		InputWindow:           TickDuration, // 1.5s
		Mode:                  InputModeHybrid,
		AutoAttackEnabled:     true,
		AutoDefendEnabled:     true,
		AutoDefendHPTreshold:  0.3, // 30% HP
		LateInputQueueing:     true,
		TalentSlots:           4,
	}
}

// InputState represents the current input state for a participant
type InputState struct {
	mu sync.RWMutex
	
	// Participant info
	ParticipantID int
	
	// Whether waiting for input
	AwaitingInput bool
	
	// When input window closes
	Deadline time.Time
	
	// Selected action (if any)
	SelectedAction *ActionDefinition
	SelectedTarget *Participant
	
	// Input received
	InputReceived bool
	
	// Input time
	InputTime time.Time
	
	// Late input (arrived after deadline)
	IsLate bool
}

// NewInputState creates a new input state for a participant
func NewInputState(participantID int) *InputState {
	return &InputState{
		ParticipantID: participantID,
	}
}

// StartInputWindow begins an input window
func (is *InputState) StartInputWindow(duration time.Duration) {
	is.mu.Lock()
	defer is.mu.Unlock()
	
	is.AwaitingInput = true
	is.Deadline = time.Now().Add(duration)
	is.SelectedAction = nil
	is.SelectedTarget = nil
	is.InputReceived = false
	is.IsLate = false
}

// SubmitInput records an input action
func (is *InputState) SubmitInput(action *ActionDefinition, target *Participant) (accepted bool, isLate bool) {
	is.mu.Lock()
	defer is.mu.Unlock()
	
	now := time.Now()
	isLate = now.After(is.Deadline)
	
	if isLate {
		is.IsLate = true
		// Still record late input for queueing
		is.SelectedAction = action
		is.SelectedTarget = target
		is.InputReceived = true
		is.InputTime = now
		return false, true
	}
	
	is.SelectedAction = action
	is.SelectedTarget = target
	is.InputReceived = true
	is.InputTime = now
	is.AwaitingInput = false
	
	return true, false
}

// HasInput returns true if input was received
func (is *InputState) HasInput() bool {
	is.mu.RLock()
	defer is.mu.RUnlock()
	return is.InputReceived
}

// IsLateInput returns true if the input was late
func (is *InputState) WasLate() bool {
	is.mu.RLock()
	defer is.mu.RUnlock()
	return is.IsLate
}

// TimeRemaining returns time left in the input window
func (is *InputState) TimeRemaining() time.Duration {
	is.mu.RLock()
	defer is.mu.RUnlock()
	
	remaining := time.Until(is.Deadline)
	if remaining < 0 {
		return 0
	}
	return remaining
}

// Clear resets the input state
func (is *InputState) Clear() {
	is.mu.Lock()
	defer is.mu.Unlock()
	
	is.AwaitingInput = false
	is.SelectedAction = nil
	is.SelectedTarget = nil
	is.InputReceived = false
	is.IsLate = false
}

// InputManager handles input collection for all combat participants
type InputManager struct {
	mu sync.RWMutex
	
	// Configuration
	config InputConfig
	
	// Input states by participant ID
	states map[int]*InputState
	
	// Talent bindings by participant ID (slot -> action ID)
	talentBindings map[int]map[int]string // participantID -> {slot: actionID}
	
	// Callback when input received
	onInput func(participantID int, action *ActionDefinition, target *Participant)
	
	// Callback when input times out
	onTimeout func(participantID int)
	
	// Callback for auto-action
	onAutoAction func(participantID int, action *ActionDefinition)
}

// NewInputManager creates a new input manager
func NewInputManager(config InputConfig) *InputManager {
	return &InputManager{
		config:         config,
		states:         make(map[int]*InputState),
		talentBindings: make(map[int]map[int]string),
	}
}

// SetOnInput sets the input received callback
func (im *InputManager) SetOnInput(callback func(int, *ActionDefinition, *Participant)) {
	im.mu.Lock()
	defer im.mu.Unlock()
	im.onInput = callback
}

// SetOnTimeout sets the timeout callback
func (im *InputManager) SetOnTimeout(callback func(int)) {
	im.mu.Lock()
	defer im.mu.Unlock()
	im.onTimeout = callback
}

// SetOnAutoAction sets the auto-action callback
func (im *InputManager) SetOnAutoAction(callback func(int, *ActionDefinition)) {
	im.mu.Lock()
	defer im.mu.Unlock()
	im.onAutoAction = callback
}

// RegisterParticipant adds a participant to the input manager
func (im *InputManager) RegisterParticipant(participant *Participant) {
	im.mu.Lock()
	defer im.mu.Unlock()
	
	im.states[participant.ID] = NewInputState(participant.ID)
	
	// Set default talent bindings (slots 1-4)
	im.talentBindings[participant.ID] = map[int]string{
		1: "attack",
		2: "defend",
		3: "item",
		4: "wait",
	}
}

// UnregisterParticipant removes a participant from the input manager
func (im *InputManager) UnregisterParticipant(participantID int) {
	im.mu.Lock()
	defer im.mu.Unlock()
	
	delete(im.states, participantID)
	delete(im.talentBindings, participantID)
}

// StartInputWindow starts an input window for a participant
func (im *InputManager) StartInputWindow(participantID int) {
	im.mu.RLock()
	state, exists := im.states[participantID]
	im.mu.RUnlock()
	
	if !exists {
		return
	}
	
	state.StartInputWindow(im.config.InputWindow)
}

// StartAllInputWindows starts input windows for all participants
func (im *InputManager) StartAllInputWindows(participants []*Participant) {
	for _, p := range participants {
		if p.IsAlive && p.CanAct() {
			im.StartInputWindow(p.ID)
		}
	}
}

// HandleKeyInput processes a key input (1-4 for talents)
func (im *InputManager) HandleKeyInput(participantID int, key int, target *Participant) bool {
	im.mu.RLock()
	state, exists := im.states[participantID]
	bindings := im.talentBindings[participantID]
	im.mu.RUnlock()
	
	if !exists {
		return false
	}
	
	// Get action for key
	actionID, hasBinding := bindings[key]
	if !hasBinding {
		return false
	}
	
	action, found := GetActionDefinition(actionID)
	if !found {
		return false
	}
	
	// Submit input
	accepted, isLate := state.SubmitInput(action, target)
	
	if isLate {
		// Late input - queue for next tick if enabled
		if im.config.LateInputQueueing {
			// The state already stores the late input
			return true
		}
		return false
	}
	
	// Fire callback
	if accepted && im.onInput != nil {
		im.onInput(participantID, action, target)
	}
	
	return accepted
}

// HandleActionInput processes a direct action input
func (im *InputManager) HandleActionInput(participantID int, actionID string, target *Participant) bool {
	im.mu.RLock()
	state, exists := im.states[participantID]
	im.mu.RUnlock()
	
	if !exists {
		return false
	}
	
	action, found := GetActionDefinition(actionID)
	if !found {
		return false
	}
	
	accepted, isLate := state.SubmitInput(action, target)
	
	if isLate {
		return im.config.LateInputQueueing
	}
	
	if accepted && im.onInput != nil {
		im.onInput(participantID, action, target)
	}
	
	return accepted
}

// CheckTimeout checks for timed-out input windows and applies auto-actions
func (im *InputManager) CheckTimeout(participantID int, combat *Combat) *ActionDefinition {
	im.mu.RLock()
	state, exists := im.states[participantID]
	im.mu.RUnlock()
	
	if !exists {
		return nil
	}
	
	// Check if still awaiting input
	if !state.AwaitingInput {
		return nil
	}
	
	// Check if time has expired
	if time.Now().Before(state.Deadline) {
		return nil
	}
	
	// Timeout - apply auto-action
	var autoAction *ActionDefinition
	
	// Get participant
	participant := combat.GetParticipantByID(participantID)
	if participant == nil {
		return nil
	}
	
	// Decide auto-action
	if im.config.AutoDefendEnabled && participant.HP < int(float64(participant.MaxHP)*im.config.AutoDefendHPTreshold) {
		// Low HP - auto-defend
		autoAction, _ = GetActionDefinition("defend")
	} else if im.config.AutoAttackEnabled {
		// Default - auto-attack
		autoAction, _ = GetActionDefinition("attack")
	}
	
	// Fire timeout callback
	if im.onTimeout != nil {
		im.onTimeout(participantID)
	}
	
	// Fire auto-action callback
	if autoAction != nil && im.onAutoAction != nil {
		im.onAutoAction(participantID, autoAction)
	}
	
	// Clear input state
	state.Clear()
	
	return autoAction
}

// GetInputState returns the input state for a participant
func (im *InputManager) GetInputState(participantID int) *InputState {
	im.mu.RLock()
	defer im.mu.RUnlock()
	return im.states[participantID]
}

// SetTalentBinding binds an action to a talent slot
func (im *InputManager) SetTalentBinding(participantID int, slot int, actionID string) bool {
	// Validate slot
	if slot < 1 || slot > im.config.TalentSlots {
		return false
	}
	
	// Validate action exists
	if _, found := GetActionDefinition(actionID); !found {
		return false
	}
	
	im.mu.Lock()
	defer im.mu.Unlock()
	
	if im.talentBindings[participantID] == nil {
		im.talentBindings[participantID] = make(map[int]string)
	}
	
	im.talentBindings[participantID][slot] = actionID
	return true
}

// GetTalentBindings returns all talent bindings for a participant
func (im *InputManager) GetTalentBindings(participantID int) map[int]string {
	im.mu.RLock()
	defer im.mu.RUnlock()
	
	bindings := im.talentBindings[participantID]
	if bindings == nil {
		return map[int]string{}
	}
	
	// Return a copy
	result := make(map[int]string)
	for k, v := range bindings {
		result[k] = v
	}
	return result
}

// GetLateInputs returns all late inputs that should be queued for next tick
func (im *InputManager) GetLateInputs() map[int]*InputState {
	im.mu.RLock()
	defer im.mu.RUnlock()
	
	result := make(map[int]*InputState)
	for id, state := range im.states {
		if state.IsLate {
			result[id] = state
		}
	}
	return result
}

// ClearLateInputs clears all late inputs after queuing
func (im *InputManager) ClearLateInputs() {
	im.mu.Lock()
	defer im.mu.Unlock()
	
	for _, state := range im.states {
		if state.IsLate {
			state.Clear()
		}
	}
}

// GetTimeRemaining returns formatted time remaining for input
func (im *InputManager) GetTimeRemaining(participantID int) string {
	state := im.GetInputState(participantID)
	if state == nil {
		return "0.0"
	}
	
	remaining := state.TimeRemaining().Seconds()
	if remaining < 0 {
		return "0.0"
	}
	return FormatCombatLog("%.1f", remaining)
}

// IsAwaitingInput returns true if waiting for participant input
func (im *InputManager) IsAwaitingInput(participantID int) bool {
	state := im.GetInputState(participantID)
	if state == nil {
		return false
	}
	return state.AwaitingInput
}