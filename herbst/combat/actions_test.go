package combat

import (
	"testing"
)

func TestGetActionDefinition_BasicActions(t *testing.T) {
	// Test all basic actions exist
	basicActions := []string{"attack", "defend", "flee", "item", "wait"}
	
	for _, actionID := range basicActions {
		action, exists := GetActionDefinition(actionID)
		if !exists {
			t.Errorf("Basic action '%s' should exist", actionID)
			continue
		}
		
		if action.Category != CategoryBasic {
			t.Errorf("Action '%s' should be CategoryBasic, got %v", actionID, action.Category)
		}
	}
}

func TestGetActionDefinition_SkillActions(t *testing.T) {
	// Test some skill actions exist
	skillActions := []string{"slash", "heavy_strike", "parry", "fireball", "heal", "backstab"}
	
	for _, actionID := range skillActions {
		action, exists := GetActionDefinition(actionID)
		if !exists {
			t.Errorf("Skill action '%s' should exist", actionID)
			continue
		}
		
		if action.Category != CategorySkill {
			t.Errorf("Action '%s' should be CategorySkill, got %v", actionID, action.Category)
		}
	}
}

func TestGetActionDefinition_NonExistent(t *testing.T) {
	_, exists := GetActionDefinition("nonexistent_action")
	if exists {
		t.Error("Non-existent action should not return true")
	}
}

func TestActionDefinition_TickCosts(t *testing.T) {
	tests := []struct {
		actionID  string
		tickCost  int
		actionType ActionType
	}{
		{"attack", 1, ActionInstant},
		{"defend", 1, ActionInstant},
		{"wait", 0, ActionInstant},
		{"heavy_strike", 1, ActionCharge}, // Has charge ticks
		{"fireball", 1, ActionCharge},
		{"channel_mana", 0, ActionChannel}, // Channel type
	}
	
	for _, tt := range tests {
		action, exists := GetActionDefinition(tt.actionID)
		if !exists {
			t.Errorf("Action '%s' should exist", tt.actionID)
			continue
		}
		
		if action.TickCost != tt.tickCost {
			t.Errorf("Action '%s' should have tick cost %d, got %d", 
				tt.actionID, tt.tickCost, action.TickCost)
		}
		
		if action.Type != tt.actionType {
			t.Errorf("Action '%s' should have type %s, got %s",
				tt.actionID, tt.actionType, action.Type)
		}
	}
}

func TestActionDefinition_ResourceCosts(t *testing.T) {
	// Attack should be free
	attack, _ := GetActionDefinition("attack")
	if attack.ManaCost != 0 || attack.StaminaCost != 0 {
		t.Errorf("Attack should be free, got mana=%d stamina=%d", 
			attack.ManaCost, attack.StaminaCost)
	}
	
	// Defend should cost stamina
	defend, _ := GetActionDefinition("defend")
	if defend.StaminaCost != 5 {
		t.Errorf("Defend should cost 5 stamina, got %d", defend.StaminaCost)
	}
	
	// Fireball should cost mana
	fireball, _ := GetActionDefinition("fireball")
	if fireball.ManaCost != 20 {
		t.Errorf("Fireball should cost 20 mana, got %d", fireball.ManaCost)
	}
}

func TestGetActionsByCategory(t *testing.T) {
	basicActions := GetActionsByCategory(CategoryBasic)
	if len(basicActions) != 5 {
		t.Errorf("Expected 5 basic actions, got %d", len(basicActions))
	}
	
	skillActions := GetActionsByCategory(CategorySkill)
	if len(skillActions) < 12 {
		t.Errorf("Expected at least 12 skill actions, got %d", len(skillActions))
	}
}

func TestGetActionsByType(t *testing.T) {
	instantActions := GetActionsByType(ActionInstant)
	if len(instantActions) == 0 {
		t.Error("Should have instant actions")
	}
	
	chargeActions := GetActionsByType(ActionCharge)
	if len(chargeActions) == 0 {
		t.Error("Should have charge actions")
	}
	
	channelActions := GetActionsByType(ActionChannel)
	if len(channelActions) == 0 {
		t.Error("Should have channel actions")
	}
}

func TestGetAllActionDefinitions(t *testing.T) {
	all := GetAllActionDefinitions()
	
	// Should include both basic and skill actions
	if _, exists := all["attack"]; !exists {
		t.Error("Should include basic 'attack'")
	}
	
	if _, exists := all["fireball"]; !exists {
		t.Error("Should include skill 'fireball'")
	}
}