package combat

import (
	"testing"
)

func TestCombatManager_CreateCombat(t *testing.T) {
	cm := NewCombatManager()
	
	participants := []*Participant{
		{ID: 1, Name: "Hero", IsPlayer: true, Team: 0, HP: 100, MaxHP: 100, Dexterity: 15},
		{ID: 2, Name: "Goblin", IsNPC: true, Team: 1, HP: 30, MaxHP: 30, Dexterity: 10},
	}
	
	id := cm.CreateCombat(5, participants)
	
	if id != 1 {
		t.Errorf("Expected combat ID 1, got %d", id)
	}
	
	combat, exists := cm.GetCombat(id)
	if !exists {
		t.Error("Combat should exist after creation")
	}
	
	if len(combat.Participants) != 2 {
		t.Errorf("Expected 2 participants, got %d", len(combat.Participants))
	}
}

func TestCombatManager_EndCombat(t *testing.T) {
	cm := NewCombatManager()
	
	participants := []*Participant{
		{ID: 1, Name: "Hero", IsPlayer: true, Team: 0},
		{ID: 2, Name: "Goblin", IsNPC: true, Team: 1},
	}
	
	id := cm.CreateCombat(5, participants)
	cm.EndCombat(id, "test end")
	
	_, exists := cm.GetCombat(id)
	if exists {
		t.Error("Combat should not exist after ending")
	}
	
	if count := cm.GetCombatCount(); count != 0 {
		t.Errorf("Expected 0 combats, got %d", count)
	}
}

func TestCombatManager_GetCombatsByRoom(t *testing.T) {
	cm := NewCombatManager()
	
	participants := []*Participant{
		{ID: 1, Name: "Hero", IsPlayer: true, Team: 0},
	}
	
	_ = cm.CreateCombat(5, participants)
	_ = cm.CreateCombat(10, participants)
	
	room5Combats := cm.GetCombatsByRoom(5)
	if len(room5Combats) != 1 {
		t.Errorf("Expected 1 combat in room 5, got %d", len(room5Combats))
	}
	
	room10Combats := cm.GetCombatsByRoom(10)
	if len(room10Combats) != 1 {
		t.Errorf("Expected 1 combat in room 10, got %d", len(room10Combats))
	}
	
	room99Combats := cm.GetCombatsByRoom(99)
	if len(room99Combats) != 0 {
		t.Errorf("Expected 0 combats in room 99, got %d", len(room99Combats))
	}
}

func TestCombat_GetAliveParticipants(t *testing.T) {
	participants := []*Participant{
		{ID: 1, Name: "Hero", IsPlayer: true, Team: 0, IsAlive: true, HP: 100},
		{ID: 2, Name: "Goblin", IsNPC: true, Team: 1, IsAlive: false, HP: 0},
		{ID: 3, Name: "Orc", IsNPC: true, Team: 1, IsAlive: true, HP: 50},
	}
	
	combat := &Combat{
		ID:           1,
		Participants: participants,
	}
	
	alive := combat.GetAliveParticipants()
	if len(alive) != 2 {
		t.Errorf("Expected 2 alive participants, got %d", len(alive))
	}
}

func TestCombat_AllEnemiesDefeated(t *testing.T) {
	participants := []*Participant{
		{ID: 1, Name: "Hero", IsPlayer: true, Team: 0, IsAlive: true, HP: 100},
		{ID: 2, Name: "Goblin", IsNPC: true, Team: 1, IsAlive: false, HP: 0},
	}
	
	combat := &Combat{
		ID:           1,
		Participants: participants,
	}
	
	if !combat.AllEnemiesDefeated() {
		t.Error("All enemies should be defeated")
	}
}

func TestCombat_AllPlayersDefeated(t *testing.T) {
	participants := []*Participant{
		{ID: 1, Name: "Hero", IsPlayer: true, Team: 0, IsAlive: false, HP: 0},
		{ID: 2, Name: "Goblin", IsNPC: true, Team: 1, IsAlive: true, HP: 30},
	}
	
	combat := &Combat{
		ID:           1,
		Participants: participants,
	}
	
	if !combat.AllPlayersDefeated() {
		t.Error("All players should be defeated")
	}
}