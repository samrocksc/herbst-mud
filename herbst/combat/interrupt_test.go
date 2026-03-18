package combat

import (
	"testing"
)

func TestParryAttempt_Success(t *testing.T) {
	combat := &Combat{
		ID:    1,
		State: StateActive,
		Participants: []*Participant{
			{ID: 1, Name: "Player", IsPlayer: true, Team: 0, HP: 50, MaxHP: 50, Stamina: 20, IsAlive: true},
			{ID: 2, Name: "Enemy", IsPlayer: false, Team: 1, HP: 30, MaxHP: 30, IsAlive: true},
		},
		Effects: NewEffectRegistry(),
	}

	defender := combat.Participants[0]
	attacker := combat.Participants[1]

	attackDef, _ := GetActionDefinition("attack")
	attackerAction := &QueuedAction{
		Action: attackDef,
		Source: attacker,
		Target: defender,
	}

	result := ParryAttempt(combat, defender, attackerAction)

	if result == nil {
		t.Fatal("Expected parry result, got nil")
	}
	if !result.Success {
		t.Errorf("Expected successful parry, got failure: %s", result.Message)
	}
	if result.Type != InterruptParry {
		t.Errorf("Expected PARRY type, got %s", result.Type)
	}
	if result.DefenderID != defender.ID {
		t.Errorf("Expected defender ID %d, got %d", defender.ID, result.DefenderID)
	}
}

func TestParryAttempt_NoStamina(t *testing.T) {
	combat := &Combat{
		ID:    1,
		State: StateActive,
		Participants: []*Participant{
			{ID: 1, Name: "Player", IsPlayer: true, Team: 0, HP: 50, MaxHP: 50, Stamina: 5, IsAlive: true}, // Not enough stamina
			{ID: 2, Name: "Enemy", IsPlayer: false, Team: 1, HP: 30, MaxHP: 30, IsAlive: true},
		},
		Effects: NewEffectRegistry(),
	}

	defender := combat.Participants[0]
	attacker := combat.Participants[1]

	attackDef, _ := GetActionDefinition("attack")
	attackerAction := &QueuedAction{
		Action: attackDef,
		Source: attacker,
		Target: defender,
	}

	result := ParryAttempt(combat, defender, attackerAction)

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
	if result.Success {
		t.Error("Expected parry to fail due to no stamina")
	}
}

func TestShieldBashAttempt_Success(t *testing.T) {
	combat := &Combat{
		ID:    1,
		State: StateActive,
		Participants: []*Participant{
			{ID: 1, Name: "Player", IsPlayer: true, Team: 0, HP: 50, MaxHP: 50, Stamina: 20, IsAlive: true},
			{ID: 2, Name: "Enemy", IsPlayer: false, Team: 1, HP: 30, MaxHP: 30, IsAlive: true},
		},
		Effects: NewEffectRegistry(),
	}

	attacker := combat.Participants[0]
	target := combat.Participants[1]

	// Set up enemy channeling an action
	fireballDef, _ := GetActionDefinition("fireball")
	target.CurrentAction = &QueuedAction{
		Action: fireballDef,
		Source: target,
	}

	result := ShieldBashAttempt(combat, attacker, target)

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
	if !result.Success {
		t.Errorf("Expected successful shield bash, got failure: %s", result.Message)
	}
	if result.Type != InterruptShieldBash {
		t.Errorf("Expected SHIELD_BASH type, got %s", result.Type)
	}
	if result.ActionCancelled != "fireball" {
		t.Errorf("Expected 'fireball' action cancelled, got '%s'", result.ActionCancelled)
	}
	if result.StunDuration != 2 {
		t.Errorf("Expected 2 tick stun, got %d", result.StunDuration)
	}

	// Check stun was applied
	if !combat.Effects.IsStunned(target.ID) {
		t.Error("Expected target to be stunned")
	}
}

func TestShieldBashAttempt_InterruptsChannel(t *testing.T) {
	combat := &Combat{
		ID:    1,
		State: StateActive,
		Participants: []*Participant{
			{ID: 1, Name: "Player", IsPlayer: true, Team: 0, HP: 50, MaxHP: 50, Stamina: 20, IsAlive: true},
			{ID: 2, Name: "Enemy", IsPlayer: false, Team: 1, HP: 30, MaxHP: 30, IsAlive: true},
		},
		Effects: NewEffectRegistry(),
	}

	attacker := combat.Participants[0]
	target := combat.Participants[1]

	// Enemy is channeling a powerful attack
	fireballDef, _ := GetActionDefinition("fireball")
	target.CurrentAction = &QueuedAction{
		Action: fireballDef,
		Source: target,
	}

	result := ShieldBashAttempt(combat, attacker, target)

	// Check channel was cancelled
	if target.CurrentAction != nil {
		t.Error("Expected target's current action to be cancelled")
	}
	if result.ActionCancelled != "fireball" {
		t.Errorf("Expected 'fireball' to be cancelled, got '%s'", result.ActionCancelled)
	}
}

func TestStunInterrupt(t *testing.T) {
	combat := &Combat{
		ID:    1,
		State: StateActive,
		Participants: []*Participant{
			{ID: 1, Name: "Player", IsPlayer: true, Team: 0, HP: 50, MaxHP: 50, IsAlive: true},
			{ID: 2, Name: "Enemy", IsPlayer: false, Team: 1, HP: 30, MaxHP: 30, IsAlive: true},
		},
		Effects: NewEffectRegistry(),
	}

	attacker := combat.Participants[0]
	target := combat.Participants[1]

	result := StunInterrupt(combat, attacker, target, 3)

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
	if !result.Success {
		t.Error("Expected successful stun")
	}
	if result.StunDuration != 3 {
		t.Errorf("Expected 3 tick stun, got %d", result.StunDuration)
	}

	// Check stun was applied
	if !combat.Effects.IsStunned(target.ID) {
		t.Error("Expected target to be stunned")
	}
}

func TestCanInterrupt_Channel(t *testing.T) {
	// Create a channeling action that can be interrupted
	// Using channel_mana which is defined as ActionChannel
	channelDef, _ := GetActionDefinition("channel_mana")
	if channelDef == nil {
		t.Skip("channel_mana action not defined")
	}
	
	channelAction := &QueuedAction{
		Action: channelDef,
	}

	if !CanInterrupt(channelAction) {
		t.Error("Channeling action should be interruptible")
	}
}

func TestCanInterrupt_Charge(t *testing.T) {
	// Charging action (still charging) can be interrupted
	heavyStrikeDef, _ := GetActionDefinition("heavy_strike")
	chargeAction := &QueuedAction{
		Action:     heavyStrikeDef,
		ChargeTick: 1, // Still charging
	}

	if !CanInterrupt(chargeAction) {
		t.Error("Charging action should be interruptible")
	}

	// Fully charged action is not interruptible
	readyAction := &QueuedAction{
		Action:     heavyStrikeDef,
		ChargeTick: 0,
	}

	if CanInterrupt(readyAction) {
		t.Error("Ready charged action should not be interruptible")
	}
}

func TestCanInterrupt_Instant(t *testing.T) {
	// Instant action cannot be interrupted
	attackDef, _ := GetActionDefinition("attack")
	instantAction := &QueuedAction{
		Action: attackDef,
	}

	if CanInterrupt(instantAction) {
		t.Error("Instant action should not be interruptible")
	}
}

func TestCheckParryTiming_Valid(t *testing.T) {
	// Parry action
	parryDef, _ := GetActionDefinition("parry")
	parryAction := &QueuedAction{
		Action: parryDef,
	}

	// Attack action targeting the parrying player
	attackDef, _ := GetActionDefinition("attack")
	attackAction := &QueuedAction{
		Action: attackDef,
	}

	// Same tick check
	valid := CheckParryTiming(parryAction, attackAction, 1)
	if !valid {
		t.Error("Parry should be valid against attack on same tick")
	}
}

func TestCheckParryTiming_WrongAction(t *testing.T) {
	// Not a parry
	defendDef, _ := GetActionDefinition("defend")
	defendAction := &QueuedAction{
		Action: defendDef,
	}

	attackDef, _ := GetActionDefinition("attack")
	attackAction := &QueuedAction{
		Action: attackDef,
	}

	valid := CheckParryTiming(defendAction, attackAction, 1)
	if valid {
		t.Error("Defend should not count as parry")
	}
}

func TestProcessInterrupt_Damage(t *testing.T) {
	combat := &Combat{
		ID:    1,
		State: StateActive,
		Participants: []*Participant{
			{ID: 1, Name: "Player", IsPlayer: true, Team: 0, HP: 50, MaxHP: 50, IsAlive: true},
			{ID: 2, Name: "Enemy", IsPlayer: false, Team: 1, HP: 30, MaxHP: 30, IsAlive: true},
		},
		Effects: NewEffectRegistry(),
	}

	attacker := combat.Participants[0]
	target := combat.Participants[1]

	result := &InterruptResult{
		Success:       true,
		Type:          InterruptShieldBash,
		AttackerID:    attacker.ID,
		DefenderID:    target.ID,
		CounterDamage: 8,
		Message:       "Player shield bashes Enemy!",
	}

	ProcessInterrupt(combat, result)

	// Check damage was applied
	if target.HP != 22 {
		t.Errorf("Expected target HP 22, got %d", target.HP)
	}

	// Check log entry
	if len(combat.Log) != 1 {
		t.Errorf("Expected 1 log entry, got %d", len(combat.Log))
	}
}

func TestProcessInterrupt_NoDamage(t *testing.T) {
	combat := &Combat{
		ID:    1,
		State: StateActive,
		Participants: []*Participant{
			{ID: 1, Name: "Player", IsPlayer: true, Team: 0, HP: 50, MaxHP: 50, IsAlive: true},
			{ID: 2, Name: "Enemy", IsPlayer: false, Team: 1, HP: 30, MaxHP: 30, IsAlive: true},
		},
		Effects: NewEffectRegistry(),
	}

	result := &InterruptResult{
		Success:      true,
		Type:         InterruptStun,
		AttackerID:   1,
		DefenderID:   2,
		StunDuration: 2,
		Message:      "Enemy is stunned for 2 ticks!",
	}

	ProcessInterrupt(combat, result)

	// Check log entry exists
	if len(combat.Log) != 1 {
		t.Errorf("Expected 1 log entry, got %d", len(combat.Log))
	}
	if combat.Log[0].Type != "interrupt" {
		t.Errorf("Expected log type 'interrupt', got '%s'", combat.Log[0].Type)
	}
}

func TestApplyStunFromAction(t *testing.T) {
	combat := &Combat{
		ID:    1,
		State: StateActive,
		Participants: []*Participant{
			{ID: 1, Name: "Player", IsPlayer: true, Team: 0, HP: 50, MaxHP: 50, IsAlive: true},
			{ID: 2, Name: "Enemy", IsPlayer: false, Team: 1, HP: 30, MaxHP: 30, IsAlive: true},
		},
		Effects: NewEffectRegistry(),
	}

	attackDef, _ := GetActionDefinition("attack")
	action := &QueuedAction{
		Action: attackDef,
		Source: combat.Participants[0],
	}

	target := combat.Participants[1]

	success := ApplyStunFromAction(combat, action, target, 2)

	if !success {
		t.Error("Expected successful stun application")
	}

	// Check stun was applied
	if !combat.Effects.IsStunned(target.ID) {
		t.Error("Expected target to be stunned")
	}
}

func TestResolveParryOnTick(t *testing.T) {
	combat := &Combat{
		ID:    1,
		State: StateActive,
		Participants: []*Participant{
			{ID: 1, Name: "Player", IsPlayer: true, Team: 0, HP: 50, MaxHP: 50, Stamina: 20, IsAlive: true},
			{ID: 2, Name: "Enemy", IsPlayer: false, Team: 1, HP: 30, MaxHP: 30, IsAlive: true},
		},
		Effects: NewEffectRegistry(),
	}

	player := combat.Participants[0]
	enemy := combat.Participants[1]

	parryDef, _ := GetActionDefinition("parry")
	attackDef, _ := GetActionDefinition("attack")

	actions := []*QueuedAction{
		{Action: parryDef, Source: player, Target: enemy},
		{Action: attackDef, Source: enemy, Target: player},
	}

	results := ResolveParryOnTick(combat, actions)

	// Should have one parry result
	if len(results) != 1 {
		t.Errorf("Expected 1 parry result, got %d", len(results))
	}

	parryResult, exists := results[1] // Player's parry
	if !exists {
		t.Fatal("Expected parry result for player (ID 1)")
	}

	if !parryResult.Success {
		t.Error("Expected successful parry")
	}
}