package combat

import (
	"testing"
)

func TestRollInitiative(t *testing.T) {
	p := &Participant{
		ID:         1,
		Name:       "Hero",
		Dexterity:  15,
	}
	
	result := RollInitiative(p)
	
	// Initiative should be DEX + roll (1-20)
	// So result should be between 16 (15+1) and 35 (15+20)
	if result < 16 || result > 35 {
		t.Errorf("Initiative %d out of expected range [16, 35]", result)
	}
	
	if p.Initiative != result {
		t.Errorf("Participant initiative not set correctly, expected %d, got %d", result, p.Initiative)
	}
}

func TestSortByInitiative(t *testing.T) {
	participants := []*Participant{
		{ID: 1, Name: "Hero", Dexterity: 10, Initiative: 25},
		{ID: 2, Name: "Goblin", Dexterity: 8, Initiative: 15},
		{ID: 3, Name: "Mage", Dexterity: 16, Initiative: 30},
	}
	
	SortByInitiative(participants)
	
	if participants[0].ID != 3 {
		t.Errorf("Expected Mage (ID 3) first, got %s (ID %d)", participants[0].Name, participants[0].ID)
	}
	if participants[1].ID != 1 {
		t.Errorf("Expected Hero (ID 1) second, got %s (ID %d)", participants[1].Name, participants[1].ID)
	}
	if participants[2].ID != 2 {
		t.Errorf("Expected Goblin (ID 2) third, got %s (ID %d)", participants[2].Name, participants[2].ID)
	}
}

func TestGetTurnOrder(t *testing.T) {
	participants := []*Participant{
		{ID: 1, Name: "Hero", Team: 0, IsAlive: true, Initiative: 20},
		{ID: 2, Name: "Goblin", Team: 1, IsAlive: true, Initiative: 15},
		{ID: 3, Name: "Orc", Team: 1, IsAlive: true, Initiative: 18},
	}
	
	combat := &Combat{
		ID:           1,
		Participants: participants,
	}
	
	order := combat.GetTurnOrder()
	
	if len(order) != 3 {
		t.Errorf("Expected 3 participants in turn order, got %d", len(order))
	}
	
	// Check turn positions are set
	if order[0].TurnPosition != 1 {
		t.Errorf("Expected first participant to have position 1, got %d", order[0].TurnPosition)
	}
}

func TestGetCurrentActor(t *testing.T) {
	participants := []*Participant{
		{ID: 1, Name: "Hero", Team: 0, IsAlive: true, Initiative: 20},
		{ID: 2, Name: "Goblin", Team: 1, IsAlive: true, Initiative: 15},
	}
	
	combat := &Combat{
		ID:           1,
		Participants: participants,
	}
	
	// Get turn order first
	combat.GetTurnOrder()
	
	// Turn index 0 should return first actor
	actor := combat.GetCurrentActor(0)
	if actor == nil || actor.ID != 1 {
		t.Errorf("Expected Hero (ID 1) at turn 0, got %v", actor)
	}
	
	// Turn index 1 should return second actor
	actor = combat.GetCurrentActor(1)
	if actor == nil || actor.ID != 2 {
		t.Errorf("Expected Goblin (ID 2) at turn 1, got %v", actor)
	}
	
	// Turn index 2 should wrap back to first (modular)
	actor = combat.GetCurrentActor(2)
	if actor == nil || actor.ID != 1 {
		t.Errorf("Expected Hero (ID 1) at turn 2 (wrap), got %v", actor)
	}
}

func TestGetNextActor(t *testing.T) {
	participants := []*Participant{
		{ID: 1, Name: "Hero", Team: 0, IsAlive: true, Initiative: 20},
		{ID: 2, Name: "Goblin", Team: 1, IsAlive: true, Initiative: 15},
	}
	
	combat := &Combat{
		ID:           1,
		Participants: participants,
	}
	
	combat.GetTurnOrder()
	
	nextIdx, nextActor := combat.GetNextActor(0)
	
	if nextIdx != 1 {
		t.Errorf("Expected next index 1, got %d", nextIdx)
	}
	
	if nextActor == nil || nextActor.ID != 2 {
		t.Errorf("Expected Goblin (ID 2) as next actor, got %v", nextActor)
	}
	
	// Test wrap-around
	nextIdx, nextActor = combat.GetNextActor(1)
	if nextIdx != 0 {
		t.Errorf("Expected next index 0 (wrap), got %d", nextIdx)
	}
	if nextActor == nil || nextActor.ID != 1 {
		t.Errorf("Expected Hero (ID 1) as next actor (wrap), got %v", nextActor)
	}
}