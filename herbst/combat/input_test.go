package combat

import (
	"testing"
	"time"
)

func TestNewInputState(t *testing.T) {
	is := NewInputState(1)
	
	if is.ParticipantID != 1 {
		t.Errorf("Expected participant ID 1, got %d", is.ParticipantID)
	}
	
	if is.AwaitingInput {
		t.Error("Should not be awaiting input initially")
	}
}

func TestInputState_StartInputWindow(t *testing.T) {
	is := NewInputState(1)
	
	is.StartInputWindow(TickDuration)
	
	if !is.AwaitingInput {
		t.Error("Should be awaiting input after start")
	}
	
	if is.Deadline.IsZero() {
		t.Error("Deadline should be set")
	}
	
	// Deadline should be about 1.5s in the future
	expected := time.Now().Add(TickDuration)
	diff := is.Deadline.Sub(expected)
	if diff < -100*time.Millisecond || diff > 100*time.Millisecond {
		t.Errorf("Deadline should be ~1.5s from now, diff: %v", diff)
	}
}

func TestInputState_SubmitInput(t *testing.T) {
	is := NewInputState(1)
	is.StartInputWindow(TickDuration)
	
	action := BasicActions["attack"]
	target := &Participant{ID: 2, Name: "Goblin"}
	
	accepted, isLate := is.SubmitInput(action, target)
	
	if !accepted {
		t.Error("Input should be accepted")
	}
	
	if isLate {
		t.Error("Input should not be late")
	}
	
	if !is.InputReceived {
		t.Error("Input should be marked as received")
	}
	
	if is.SelectedAction != action {
		t.Error("Selected action should be stored")
	}
	
	if is.SelectedTarget != target {
		t.Error("Selected target should be stored")
	}
}

func TestInputState_SubmitInput_Late(t *testing.T) {
	is := NewInputState(1)
	
	// Set deadline in the past
	is.StartInputWindow(time.Microsecond)
	time.Sleep(2 * time.Millisecond) // Ensure timeout
	
	action := BasicActions["attack"]
	target := &Participant{ID: 2, Name: "Goblin"}
	
	accepted, isLate := is.SubmitInput(action, target)
	
	if accepted {
		t.Error("Late input should not be accepted")
	}
	
	if !isLate {
		t.Error("Input should be marked as late")
	}
	
	if !is.InputReceived {
		t.Error("Late input should still be recorded")
	}
}

func TestInputState_TimeRemaining(t *testing.T) {
	is := NewInputState(1)
	is.StartInputWindow(TickDuration)
	
	remaining := is.TimeRemaining()
	
	if remaining <= 0 {
		t.Error("Should have time remaining")
	}
	
	if remaining > TickDuration {
		t.Errorf("Remaining time should not exceed input window")
	}
}

func TestInputState_Clear(t *testing.T) {
	is := NewInputState(1)
	is.StartInputWindow(TickDuration)
	is.SubmitInput(BasicActions["attack"], &Participant{ID: 2})
	
	is.Clear()
	
	if is.AwaitingInput {
		t.Error("Should not be awaiting after clear")
	}
	
	if is.InputReceived {
		t.Error("Input received should be cleared")
	}
	
	if is.SelectedAction != nil {
		t.Error("Selected action should be cleared")
	}
}

func TestNewInputManager(t *testing.T) {
	config := DefaultInputConfig()
	im := NewInputManager(config)
	
	if im.config.InputWindow != TickDuration {
		t.Errorf("Expected input window %v, got %v", TickDuration, im.config.InputWindow)
	}
	
	if im.config.AutoAttackEnabled != true {
		t.Error("Auto-attack should be enabled by default")
	}
}

func TestInputManager_RegisterParticipant(t *testing.T) {
	im := NewInputManager(DefaultInputConfig())
	
	p := &Participant{ID: 1, Name: "Hero"}
	im.RegisterParticipant(p)
	
	state := im.GetInputState(1)
	if state == nil {
		t.Error("State should exist for registered participant")
	}
	
	bindings := im.GetTalentBindings(1)
	if len(bindings) != 4 {
		t.Errorf("Expected 4 default bindings, got %d", len(bindings))
	}
	
	// Check default bindings
	if bindings[1] != "attack" {
		t.Errorf("Slot 1 should be attack, got %s", bindings[1])
	}
}

func TestInputManager_HandleKeyInput(t *testing.T) {
	im := NewInputManager(DefaultInputConfig())
	
	p := &Participant{ID: 1, Name: "Hero"}
	target := &Participant{ID: 2, Name: "Goblin"}
	
	im.RegisterParticipant(p)
	im.StartInputWindow(1)
	
	inputReceived := false
	im.SetOnInput(func(participantID int, action *ActionDefinition, target *Participant) {
		inputReceived = true
	})
	
	accepted := im.HandleKeyInput(1, 1, target) // Slot 1 = attack
	
	if !accepted {
		t.Error("Key input should be accepted")
	}
	
	if !inputReceived {
		t.Error("Callback should have fired")
	}
}

func TestInputManager_HandleActionInput(t *testing.T) {
	im := NewInputManager(DefaultInputConfig())
	
	p := &Participant{ID: 1, Name: "Hero"}
	target := &Participant{ID: 2, Name: "Goblin"}
	
	im.RegisterParticipant(p)
	im.StartInputWindow(1)
	
	accepted := im.HandleActionInput(1, "defend", target)
	
	if !accepted {
		t.Error("Action input should be accepted")
	}
	
	state := im.GetInputState(1)
	if state.SelectedAction.ID != "defend" {
		t.Errorf("Expected defend action, got %s", state.SelectedAction.ID)
	}
}

func TestInputManager_SetTalentBinding(t *testing.T) {
	im := NewInputManager(DefaultInputConfig())
	
	p := &Participant{ID: 1, Name: "Hero"}
	im.RegisterParticipant(p)
	
	// Set custom binding
	success := im.SetTalentBinding(1, 1, "slash")
	if !success {
		t.Error("Setting binding should succeed")
	}
	
	bindings := im.GetTalentBindings(1)
	if bindings[1] != "slash" {
		t.Errorf("Slot 1 should be slash, got %s", bindings[1])
	}
	
	// Invalid slot
	success = im.SetTalentBinding(1, 5, "slash")
	if success {
		t.Error("Setting invalid slot should fail")
	}
	
	// Invalid action
	success = im.SetTalentBinding(1, 1, "nonexistent")
	if success {
		t.Error("Setting invalid action should fail")
	}
}

func TestInputManager_CheckTimeout(t *testing.T) {
	im := NewInputManager(DefaultInputConfig())
	
	p := &Participant{ID: 1, Name: "Hero", HP: 100, MaxHP: 100}
	combat := &Combat{
		Participants: []*Participant{p},
	}
	
	im.RegisterParticipant(p)
	
	// Manually set up input state with short deadline
	state := im.GetInputState(1)
	state.StartInputWindow(time.Millisecond)
	
	time.Sleep(5 * time.Millisecond) // Wait for timeout
	
	autoAction := im.CheckTimeout(1, combat)
	
	if autoAction == nil {
		t.Error("Should return auto-action on timeout")
	}
	
	// Should be attack (healthy player)
	if autoAction.ID != "attack" {
		t.Errorf("Expected attack, got %s", autoAction.ID)
	}
}

func TestInputManager_CheckTimeout_LowHP(t *testing.T) {
	im := NewInputManager(DefaultInputConfig())
	
	// Low HP player (20% of max)
	p := &Participant{ID: 1, Name: "Hero", HP: 20, MaxHP: 100}
	combat := &Combat{
		Participants: []*Participant{p},
	}
	
	im.RegisterParticipant(p)
	
	// Manually set up input state with short deadline
	state := im.GetInputState(1)
	state.StartInputWindow(time.Millisecond)
	
	time.Sleep(5 * time.Millisecond) // Wait for timeout
	
	autoAction := im.CheckTimeout(1, combat)
	
	if autoAction == nil {
		t.Error("Should return auto-action on timeout")
	}
	
	// Should be defend (low HP)
	if autoAction.ID != "defend" {
		t.Errorf("Expected defend for low HP, got %s", autoAction.ID)
	}
}

func TestInputManager_GetLateInputs(t *testing.T) {
	im := NewInputManager(DefaultInputConfig())
	
	p := &Participant{ID: 1, Name: "Hero"}
	im.RegisterParticipant(p)
	
	// Set up late input - use the InputState directly
	state := im.GetInputState(1)
	state.StartInputWindow(time.Millisecond)
	time.Sleep(5 * time.Millisecond)
	
	action, _ := GetActionDefinition("attack")
	state.SubmitInput(action, nil)
	
	lateInputs := im.GetLateInputs()
	
	if len(lateInputs) != 1 {
		t.Errorf("Expected 1 late input, got %d", len(lateInputs))
	}
	
	if !lateInputs[1].IsLate {
		t.Error("Input should be marked as late")
	}
	
	// Clear and check
	im.ClearLateInputs()
	lateInputs = im.GetLateInputs()
	
	if len(lateInputs) != 0 {
		t.Errorf("Expected 0 late inputs after clear, got %d", len(lateInputs))
	}
}

func TestInputManager_GetTimeRemaining(t *testing.T) {
	im := NewInputManager(DefaultInputConfig())
	
	p := &Participant{ID: 1, Name: "Hero"}
	im.RegisterParticipant(p)
	im.StartInputWindow(1)
	
	timeStr := im.GetTimeRemaining(1)
	
	// Should be a number string
	if timeStr == "" || timeStr == "0.0" {
		t.Errorf("Time remaining should not be zero, got %s", timeStr)
	}
}

func TestDefaultInputConfig(t *testing.T) {
	config := DefaultInputConfig()
	
	if config.InputWindow != TickDuration {
		t.Errorf("Expected input window %v, got %v", TickDuration, config.InputWindow)
	}
	
	if !config.AutoAttackEnabled {
		t.Error("Auto-attack should be enabled")
	}
	
	if !config.AutoDefendEnabled {
		t.Error("Auto-defend should be enabled")
	}
	
	if config.AutoDefendHPTreshold != 0.3 {
		t.Errorf("Expected HP threshold 0.3, got %f", config.AutoDefendHPTreshold)
	}
	
	if !config.LateInputQueueing {
		t.Error("Late input queueing should be enabled")
	}
	
	if config.TalentSlots != 4 {
		t.Errorf("Expected 4 talent slots, got %d", config.TalentSlots)
	}
}