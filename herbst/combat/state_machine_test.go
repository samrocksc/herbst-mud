package combat

import (
	"context"
	"testing"
	"time"
)

func TestCombatStateMachine_TransitionTo(t *testing.T) {
	combat := &Combat{
		ID:           1,
		Participants: []*Participant{},
		State:        StateIdle,
	}
	
	csm := NewCombatStateMachine(combat, nil, nil)
	
	// Test state transition
	enteredInit := false
	csm.SetOnStateEnter(StateInit, func(c *Combat) {
		enteredInit = true
	})
	
	csm.TransitionTo(StateInit)
	
	if combat.State != StateInit {
		t.Errorf("Expected state %s, got %s", StateInit, combat.State)
	}
	
	if !enteredInit {
		t.Error("State enter callback should have fired")
	}
}

func TestCombatStateMachine_Start(t *testing.T) {
	participants := []*Participant{
		{ID: 1, Name: "Hero", IsPlayer: true, Team: 0, IsAlive: true, HP: 100, MaxHP: 100, Dexterity: 15},
		{ID: 2, Name: "Goblin", IsNPC: true, Team: 1, IsAlive: true, HP: 30, MaxHP: 30, Dexterity: 10},
	}
	
	combat := &Combat{
		ID:           1,
		Participants: participants,
		State:        StateIdle,
		ActionQueue:  NewActionQueue(),
		Effects:      NewEffectRegistry(),
	}
	
	tm := NewTickManager()
	cm := NewCombatManager()
	csm := NewCombatStateMachine(combat, cm, tm)
	
	ctx := context.Background()
	csm.Start(ctx)
	
	// Should have transitioned to Active
	if combat.State != StateActive {
		t.Errorf("Expected state %s, got %s", StateActive, combat.State)
	}
	
	// Should have rolled initiative
	if participants[0].Initiative == 0 {
		t.Error("Participant 0 should have rolled initiative")
	}
	
	if participants[1].Initiative == 0 {
		t.Error("Participant 1 should have rolled initiative")
	}
	
	// Stop to clean up
	csm.Stop()
}

func TestCombatStateMachine_SubmitAction(t *testing.T) {
	participants := []*Participant{
		{ID: 1, Name: "Hero", IsPlayer: true, Team: 0, IsAlive: true, HP: 100, MaxHP: 100, Dexterity: 15},
		{ID: 2, Name: "Goblin", IsNPC: true, Team: 1, IsAlive: true, HP: 30, MaxHP: 30, Dexterity: 10},
	}
	
	combat := &Combat{
		ID:           1,
		Participants: participants,
		State:        StateActive,
		ActionQueue:  NewActionQueue(),
		Effects:      NewEffectRegistry(),
	}
	
	csm := NewCombatStateMachine(combat, nil, nil)
	csm.currentActorIdx = 0
	csm.awaitingInput[1] = true
	
	attack, _ := GetActionDefinition("attack")
	target := participants[1] // Goblin
	
	success := csm.SubmitAction(1, attack, target)
	
	if !success {
		t.Error("SubmitAction should succeed")
	}
	
	// Should have cleared awaiting input
	if csm.awaitingInput[1] {
		t.Error("Should have cleared awaiting input after action submission")
	}
	
	// Action should be in selected actions
	if csm.selectedActions[1] == nil {
		t.Error("Action should be stored in selectedActions")
	}
}

func TestCombatStateMachine_GetCurrentActor(t *testing.T) {
	participants := []*Participant{
		{ID: 1, Name: "Hero", IsPlayer: true, Team: 0, Initiative: 25, IsAlive: true, HP: 100},
		{ID: 2, Name: "Goblin", IsNPC: true, Team: 1, Initiative: 15, IsAlive: true, HP: 30},
	}
	
	combat := &Combat{
		ID:           1,
		Participants: participants,
		ActionQueue:  NewActionQueue(),
		Effects:      NewEffectRegistry(),
	}
	
	csm := NewCombatStateMachine(combat, nil, nil)
	
	// First actor should be Hero (higher initiative)
	csm.currentActorIdx = 0
	actor := csm.GetCurrentActor()
	
	if actor == nil || actor.ID != 1 {
		t.Errorf("Expected Hero (ID 1), got %v", actor)
	}
	
	// Second actor should be Goblin
	csm.currentActorIdx = 1
	actor = csm.GetCurrentActor()
	
	if actor == nil || actor.ID != 2 {
		t.Errorf("Expected Goblin (ID 2), got %v", actor)
	}
}

func TestCombatStateMachine_CheckEndCondition(t *testing.T) {
	t.Run("enemies defeated", func(t *testing.T) {
		participants := []*Participant{
			{ID: 1, Name: "Hero", IsPlayer: true, Team: 0, IsAlive: true, HP: 100},
			{ID: 2, Name: "Goblin", IsNPC: true, Team: 1, IsAlive: false, HP: 0},
		}
		
		combat := &Combat{
			ID:           1,
			Participants: participants,
			State:        StateActive,
			ActionQueue:  NewActionQueue(),
			Effects:      NewEffectRegistry(),
		}
		
		csm := NewCombatStateMachine(combat, nil, nil)
		result := csm.checkEndCondition()
		
		if !result {
			t.Error("Should detect enemies defeated")
		}
		
		if combat.State != StateEnded {
			t.Errorf("Combat should be ended, got %s", combat.State)
		}
	})
	
	t.Run("players defeated", func(t *testing.T) {
		participants := []*Participant{
			{ID: 1, Name: "Hero", IsPlayer: true, Team: 0, IsAlive: false, HP: 0},
			{ID: 2, Name: "Goblin", IsNPC: true, Team: 1, IsAlive: true, HP: 30},
		}
		
		combat := &Combat{
			ID:           1,
			Participants: participants,
			State:        StateActive,
			ActionQueue:  NewActionQueue(),
			Effects:      NewEffectRegistry(),
		}
		
		csm := NewCombatStateMachine(combat, nil, nil)
		result := csm.checkEndCondition()
		
		if !result {
			t.Error("Should detect players defeated")
		}
		
		if combat.State != StateEnded {
			t.Errorf("Combat should be ended, got %s", combat.State)
		}
	})
	
	t.Run("combat continues", func(t *testing.T) {
		participants := []*Participant{
			{ID: 1, Name: "Hero", IsPlayer: true, Team: 0, IsAlive: true, HP: 100},
			{ID: 2, Name: "Goblin", IsNPC: true, Team: 1, IsAlive: true, HP: 30},
		}
		
		combat := &Combat{
			ID:           1,
			Participants: participants,
			State:        StateActive,
			ActionQueue:  NewActionQueue(),
			Effects:      NewEffectRegistry(),
		}
		
		csm := NewCombatStateMachine(combat, nil, nil)
		result := csm.checkEndCondition()
		
		if result {
			t.Error("Combat should continue")
		}
	})
}

func TestCombatStateMachine_ExecuteAction(t *testing.T) {
	participants := []*Participant{
		{ID: 1, Name: "Hero", IsPlayer: true, Team: 0, IsAlive: true, HP: 100, MaxHP: 100},
		{ID: 2, Name: "Goblin", IsNPC: true, Team: 1, IsAlive: true, HP: 30, MaxHP: 30},
	}
	
	combat := &Combat{
		ID:           1,
		Participants: participants,
		ActionQueue:  NewActionQueue(),
		Effects:      NewEffectRegistry(),
		Log:          []CombatLogEntry{},
	}
	
	csm := NewCombatStateMachine(combat, nil, nil)
	
	attack, _ := GetActionDefinition("attack")
	qa := &QueuedAction{
		Action: attack,
		Source: participants[0],
		Target: participants[1],
	}
	
	csm.executeAction(qa)
	
	// Goblin should have taken damage
	if participants[1].HP >= 30 {
		t.Errorf("Goblin should have taken damage, HP: %d", participants[1].HP)
	}
	
	// Should have logged the action
	if len(combat.Log) == 0 {
		t.Error("Should have logged the action")
	}
}

func TestCombatStateMachine_GetTickCountdown(t *testing.T) {
	combat := &Combat{
		ID:          1,
		ActionQueue: NewActionQueue(),
		Effects:     NewEffectRegistry(),
	}
	
	csm := NewCombatStateMachine(combat, nil, nil)
	csm.inputWindow = TickDuration
	csm.inputDeadline = time.Now().Add(TickDuration)
	
	countdown := csm.GetTickCountdown()
	
	// Countdown should be roughly 1.5 seconds
	if countdown < 1.0 || countdown > 1.6 {
		t.Errorf("Expected countdown ~1.5s, got %f", countdown)
	}
}

func TestCombatStateMachine_AwaitInput(t *testing.T) {
	participants := []*Participant{
		{ID: 1, Name: "Hero", IsPlayer: true, Team: 0, IsAlive: true, HP: 100, Initiative: 25},
	}
	
	combat := &Combat{
		ID:           1,
		Participants: participants,
		ActionQueue:  NewActionQueue(),
		Effects:      NewEffectRegistry(),
	}
	
	csm := NewCombatStateMachine(combat, nil, nil)
	
	// Set up input waiting
	csm.awaitingInput[1] = true
	
	if !csm.IsAwaitingInput(1) {
		t.Error("Should be awaiting input for participant 1")
	}
	
	if csm.IsAwaitingInput(2) {
		t.Error("Should not be awaiting input for participant 2")
	}
}