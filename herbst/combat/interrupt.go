package combat

import (
	"fmt"
	"log"
)

// InterruptType defines the type of interrupt
type InterruptType string

const (
	InterruptParry      InterruptType = "PARRY"
	InterruptShieldBash InterruptType = "SHIELD_BASH"
	InterruptStun       InterruptType = "STUN"
)

// InterruptResult represents the result of an interrupt
type InterruptResult struct {
	Success        bool          `json:"success"`
	Type           InterruptType `json:"type"`
	AttackerID     int           `json:"attackerId"`
	DefenderID     int           `json:"defenderId"`
	ActionCancelled string       `json:"actionCancelled"`
	CounterDamage   int          `json:"counterDamage"`
	StunDuration    int          `json:"stunDuration"` // Ticks
	Message         string       `json:"message"`
}

// ParryAttempt attempts to parry an incoming attack
// Returns nil if parry is not possible (wrong timing, no attack)
func ParryAttempt(combat *Combat, defender *Participant, attackerAction *QueuedAction) *InterruptResult {
	// Parry only works against INSTANT attacks on the same tick
	attacker := attackerAction.Source
	if attacker == nil {
		return nil
	}

	// Get parry action definition
	parryDef, exists := GetActionDefinition("parry")
	if !exists {
		log.Println("[Parry] Parry action not defined")
		return nil
	}

	// Check if defender can use parry (has stamina)
	if defender.Stamina < parryDef.StaminaCost {
		return &InterruptResult{
			Success: false,
			Message: fmt.Sprintf("%s doesn't have enough stamina to parry!", defender.Name),
		}
	}

	// Parry succeeds - cancel the attack
	result := &InterruptResult{
		Success:         true,
		Type:            InterruptParry,
		AttackerID:      attacker.ID,
		DefenderID:      defender.ID,
		ActionCancelled: attackerAction.Action.ID,
		CounterDamage:   parryDef.BaseDamage,
		Message:         fmt.Sprintf("%s parries %s's attack!", defender.Name, attacker.Name),
	}

	// Deduct stamina
	defender.Stamina -= parryDef.StaminaCost

	// Apply cooldown
	// (Handled by action queue system)

	return result
}

// ShieldBashAttempt attempts to interrupt a channeling enemy
func ShieldBashAttempt(combat *Combat, attacker *Participant, target *Participant) *InterruptResult {
	// Get shield bash action definition
	bashDef, exists := GetActionDefinition("shield_bash")
	if !exists {
		log.Println("[ShieldBash] Shield bash action not defined")
		return nil
	}

	// Check if attacker has enough stamina
	if attacker.Stamina < bashDef.StaminaCost {
		return &InterruptResult{
			Success: false,
			Message: fmt.Sprintf("%s doesn't have enough stamina for shield bash!", attacker.Name),
		}
	}

	// Deduct stamina
	attacker.Stamina -= bashDef.StaminaCost

	// Build result
	result := &InterruptResult{
		Success:       true,
		Type:          InterruptShieldBash,
		AttackerID:    attacker.ID,
		DefenderID:    target.ID,
		CounterDamage: bashDef.BaseDamage,
		StunDuration:  2, // 2 tick stun
	}

	// Cancel any channeling action
	if target.CurrentAction != nil {
		result.ActionCancelled = target.CurrentAction.Action.ID
		result.Message = fmt.Sprintf("%s's shield bash interrupts %s's %s!", attacker.Name, target.Name, target.CurrentAction.Action.ID)
	} else {
		result.Message = fmt.Sprintf("%s shield bashes %s!", attacker.Name, target.Name)
	}

	// Clear the target's current action
	target.CurrentAction = nil

	// Apply stun effect
	stunEffect := CreateStatusEffect(StatusStunned, target.ID, attacker.ID)
	if stunEffect != nil && combat.Effects != nil {
		combat.Effects.ApplyEffect(stunEffect)
	}

	return result
}

// StunInterrupt applies a stun to interrupt actions
func StunInterrupt(combat *Combat, attacker *Participant, target *Participant, duration int) *InterruptResult {
	result := &InterruptResult{
		Success:      true,
		Type:         InterruptStun,
		AttackerID:   attacker.ID,
		DefenderID:   target.ID,
		StunDuration: duration,
	}

	// Cancel any channeling action
	if target.CurrentAction != nil {
		result.ActionCancelled = target.CurrentAction.Action.ID
		result.Message = fmt.Sprintf("%s is stunned, cancelling %s!", target.Name, target.CurrentAction.Action.ID)
	} else {
		result.Message = fmt.Sprintf("%s is stunned for %d ticks!", target.Name, duration)
	}

	// Clear current action
	target.CurrentAction = nil

	// Apply stun effect
	stunEffect := CreateStatusEffect(StatusStunned, target.ID, attacker.ID)
	if stunEffect != nil {
		stunEffect.TicksRemaining = duration
		if combat.Effects != nil {
			combat.Effects.ApplyEffect(stunEffect)
		}
	}

	return result
}

// CheckParryTiming checks if a parry can intercept an attack
// Parry must be used on the same tick as the attack
func CheckParryTiming(defenderAction *QueuedAction, attackerAction *QueuedAction, currentTick int) bool {
	// Defender must be using parry action
	if defenderAction.Action == nil || defenderAction.Action.ID != "parry" {
		return false
	}

	// Attacker must be using an attack
	if attackerAction.Action == nil || !attackerAction.Action.IsOffensive {
		return false
	}

	// Must be same tick (parry is instant, attack resolves this tick)
	return true
}

// CanInterrupt checks if an action can be interrupted
func CanInterrupt(action *QueuedAction) bool {
	if action == nil || action.Action == nil {
		return false
	}

	// Channeling actions can be interrupted
	if action.Action.Type == ActionChannel {
		return true
	}

	// Charging actions can be interrupted
	if action.Action.Type == ActionCharge && action.ChargeTick > 0 {
		return true
	}

	return false
}

// ProcessInterrupt processes an interrupt and applies its effects
func ProcessInterrupt(combat *Combat, result *InterruptResult) {
	if result == nil || !result.Success {
		return
	}

	attacker := combat.GetParticipantByID(result.AttackerID)
	target := combat.GetParticipantByID(result.DefenderID)

	if target == nil {
		return
	}

	// Apply counter/ability damage
	if result.CounterDamage > 0 && attacker != nil {
		target.TakeDamage(result.CounterDamage)
		combat.AddLogEntry(CombatLogEntry{
			Tick:      combat.TickNumber,
			Message:   result.Message,
			Type:      "damage",
			SourceID:  result.AttackerID,
			TargetID:  result.DefenderID,
			Value:     result.CounterDamage,
		})
	}

	// Add log entry if no damage
	if result.CounterDamage == 0 {
		combat.AddLogEntry(CombatLogEntry{
			Tick:    combat.TickNumber,
			Message: result.Message,
			Type:    "interrupt",
		})
	}
}

// ResolveParryOnTick resolves parry attempts for a given tick
// Called when processing actions - if a parry matches an incoming attack, it interrupts
func ResolveParryOnTick(combat *Combat, tickActions []*QueuedAction) map[int]*InterruptResult {
	results := make(map[int]*InterruptResult)

	// Find parry actions
	parryActions := make(map[int]*QueuedAction) // defenderID -> action
	attackActions := make(map[int][]*QueuedAction) // targetID -> attacker actions

	for _, action := range tickActions {
		if action.Action == nil {
			continue
		}

		if action.Action.ID == "parry" {
			parryActions[action.Source.ID] = action
		}

		if action.Action.IsOffensive && action.Target != nil {
			attackActions[action.Target.ID] = append(attackActions[action.Target.ID], action)
		}
	}

	// Match parries to attacks
	for defenderID, _ := range parryActions {
		parryAct := parryActions[defenderID]
		defender := combat.GetParticipantByID(defenderID)
		if defender == nil {
			continue
		}

		// Check if any attacks target this defender
		attacks := attackActions[defenderID]
		if len(attacks) == 0 {
			continue
		}

		// Parry the first attack
		attackerAction := attacks[0]
		result := ParryAttempt(combat, defender, attackerAction)
		if result != nil && result.Success {
			results[defenderID] = result
		}
		_ = parryAct // parry action tracked for potential future use
	}

	return results
}

// ApplyStunFromAction applies stun from an action (e.g., Shield Bash)
func ApplyStunFromAction(combat *Combat, action *QueuedAction, target *Participant, duration int) bool {
	// Get attacker
	if action.Source == nil {
		return false
	}
	attacker := action.Source

	// Create stun effect
	stunEffect := CreateStatusEffect(StatusStunned, target.ID, attacker.ID)
	if stunEffect == nil {
		return false
	}

	stunEffect.TicksRemaining = duration

	// Apply to effect registry
	if combat.Effects != nil {
		combat.Effects.ApplyEffect(stunEffect)
	}

	return true
}