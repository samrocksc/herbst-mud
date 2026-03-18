package combat

import (
	"strings"
	"testing"
	"time"
)

func TestRenderVictoryScreen_Basic(t *testing.T) {
	ui := NewCombatUI(80, 20)

	result := &VictoryResult{
		Loot: []LootItem{
			{ID: 1, Name: "Rusty Pipe", Description: "A bent pipe", Quantity: 1, Rarity: "common"},
		},
		Coins:       23,
		XP:          150,
		VictoryTime: time.Now(),
	}

	output := ui.RenderVictoryScreen(result, 80)

	// Check for victory elements
	if !strings.Contains(output, "VICTORY") {
		t.Error("Expected victory title in output")
	}

	if !strings.Contains(output, "Rusty Pipe") {
		t.Error("Expected loot item in output")
	}

	if !strings.Contains(output, "23") {
		t.Error("Expected coin amount in output")
	}

	if !strings.Contains(output, "150") {
		t.Error("Expected XP amount in output")
	}
}

func TestRenderVictoryScreen_WithWeapons(t *testing.T) {
	ui := NewCombatUI(80, 20)

	result := &VictoryResult{
		Loot:        []LootItem{},
		XP:          200,
		VictoryTime: time.Now(),
		WeaponDrops: []WeaponDrop{
			{Name: "Scrap Machete", WeaponType: "sword", ClassRestriction: "warrior", MinDamage: 4, MaxDamage: 8, Guaranteed: true},
		},
	}

	output := ui.RenderVictoryScreen(result, 80)

	// Check for weapon section
	if !strings.Contains(output, "WEAPONS") {
		t.Error("Expected weapons section in output")
	}

	if !strings.Contains(output, "Scrap Machete") {
		t.Error("Expected weapon name in output")
	}

	if !strings.Contains(output, "warrior") {
		t.Error("Expected class restriction in output")
	}
}

func TestRenderVictoryScreen_WithSkillUps(t *testing.T) {
	ui := NewCombatUI(80, 20)

	result := &VictoryResult{
		XP:          100,
		VictoryTime: time.Now(),
		SkillUps: []SkillUp{
			{SkillName: "blades", Amount: 2, NewLevel: 47},
			{SkillName: "brawling", Amount: 1, NewLevel: 15},
		},
	}

	output := ui.RenderVictoryScreen(result, 80)

	// Check for skill up section
	if !strings.Contains(output, "SKILL UP") {
		t.Error("Expected skill up section in output")
	}

	if !strings.Contains(output, "blades") {
		t.Error("Expected blades skill in output")
	}

	if !strings.Contains(output, "brawling") {
		t.Error("Expected brawling skill in output")
	}
}

func TestRenderVictoryScreen_Rarity(t *testing.T) {
	ui := NewCombatUI(80, 20)

	result := &VictoryResult{
		Loot: []LootItem{
			{ID: 1, Name: "Common Item", Rarity: "common"},
			{ID: 2, Name: "Uncommon Item", Rarity: "uncommon"},
			{ID: 3, Name: "Rare Item", Rarity: "rare"},
			{ID: 4, Name: "Legendary Item", Rarity: "legendary"},
		},
		XP:          100,
		VictoryTime: time.Now(),
	}

	output := ui.RenderVictoryScreen(result, 80)

	// All items should appear
	if !strings.Contains(output, "Common Item") {
		t.Error("Expected common item in output")
	}

	if !strings.Contains(output, "Uncommon Item") {
		t.Error("Expected uncommon item in output")
	}

	if !strings.Contains(output, "Rare Item") {
		t.Error("Expected rare item in output")
	}

	if !strings.Contains(output, "Legendary Item") {
		t.Error("Expected legendary item in output")
	}
}

func TestRenderDefeatScreen_Basic(t *testing.T) {
	ui := NewCombatUI(80, 20)

	result := &DefeatResult{
		XPLost:         10,
		RespawnPoint:   1,
		CorpseLocation:  5,
		CorpseExpiry:   time.Now().Add(5 * time.Minute),
		Consequences: []string{
			"Lose 10 XP (10% of current level progress)",
			"Your equipment has been left in the combat area",
		},
		DefeatTime: time.Now(),
	}

	output := ui.RenderDefeatScreen(result, "Old Scrap", 80)

	// Check for defeat elements
	if !strings.Contains(output, "DEFEAT") {
		t.Error("Expected defeat title in output")
	}

	if !strings.Contains(output, "Old Scrap") {
		t.Error("Expected enemy name in output")
	}

	if !strings.Contains(output, "CONSEQUENCES") {
		t.Error("Expected consequences section in output")
	}

	if !strings.Contains(output, "equipment") {
		t.Error("Expected equipment consequence in output")
	}
}

func TestRenderDefeatScreen_Consequences(t *testing.T) {
	ui := NewCombatUI(80, 20)

	result := &DefeatResult{
		XPLost:         50,
		RespawnPoint:    1,
		CorpseLocation:  10,
		CorpseExpiry:    time.Now().Add(5 * time.Minute),
		Consequences:    []string{"Lose 50 XP", "Drop 1 item"},
		DefeatTime:      time.Now(),
	}

	output := ui.RenderDefeatScreen(result, "Boss", 80)

	// All consequences should appear
	for _, c := range result.Consequences {
		if !strings.Contains(output, c) {
			t.Errorf("Expected consequence '%s' in output", c)
		}
	}
}

func TestRenderCombatEnd_Victory(t *testing.T) {
	combat := &Combat{
		ID:    1,
		State: StateEnded,
		Participants: []*Participant{
			{ID: 1, Name: "Player", IsPlayer: true, Team: 0, HP: 50, MaxHP: 50, IsAlive: true},
			{ID: 2, Name: "Enemy", IsNPC: true, Team: 1, HP: 0, MaxHP: 20, Level: 1, IsAlive: false},
		},
		Log: []CombatLogEntry{},
	}

	output := RenderCombatEnd(combat, 1, 1, 80)

	// Should show victory
	if !strings.Contains(output, "VICTORY") {
		t.Error("Expected victory screen when all enemies defeated")
	}
}

func TestRenderCombatEnd_Defeat(t *testing.T) {
	combat := &Combat{
		ID:    1,
		State: StateEnded,
		Participants: []*Participant{
			{ID: 1, Name: "Player", IsPlayer: true, Team: 0, HP: 0, MaxHP: 50, IsAlive: false},
			{ID: 2, Name: "Boss", IsNPC: true, Team: 1, HP: 30, MaxHP: 100, Level: 10, IsAlive: true},
		},
		Log: []CombatLogEntry{},
	}

	output := RenderCombatEnd(combat, 1, 1, 80)

	// Should show defeat
	if !strings.Contains(output, "DEFEAT") {
		t.Error("Expected defeat screen when all players defeated")
	}

	// Should show the enemy that delivered the killing blow
	if !strings.Contains(output, "Boss") {
		t.Error("Expected enemy name in defeat screen")
	}
}

func TestRenderCombatEnd_Ongoing(t *testing.T) {
	combat := &Combat{
		ID:    1,
		State: StateActive,
		Participants: []*Participant{
			{ID: 1, Name: "Player", IsPlayer: true, Team: 0, HP: 50, MaxHP: 50, IsAlive: true},
			{ID: 2, Name: "Enemy", IsNPC: true, Team: 1, HP: 20, MaxHP: 20, Level: 1, IsAlive: true},
		},
		Log: []CombatLogEntry{},
	}

	output := RenderCombatEnd(combat, 1, 1, 80)

	// Should return empty string for ongoing combat
	if output != "" {
		t.Errorf("Expected empty string for ongoing combat, got: %s", output)
	}
}

func TestVictoryResult_EmptyLoot(t *testing.T) {
	combat := &Combat{
		ID:    1,
		State: StateActive,
		Participants: []*Participant{
			{ID: 1, Name: "Player", IsPlayer: true, Team: 0, HP: 50, MaxHP: 50, IsAlive: true},
			{ID: 2, Name: "Weak Enemy", IsNPC: true, Team: 1, HP: 0, MaxHP: 5, Level: 1, IsAlive: false},
		},
		Log: []CombatLogEntry{},
	}

	result := combat.HandleVictory()

	// Should still give XP even with minimal enemies
	if result.XP < 1 {
		t.Error("Expected at least some XP for victory")
	}

	// Victory time should be set
	if result.VictoryTime.IsZero() {
		t.Error("Expected victory time to be set")
	}
}

func TestDefeatResult_PlayerRespawn(t *testing.T) {
	combat := &Combat{
		ID:     1,
		RoomID: 5,
		State:  StateActive,
		Participants: []*Participant{
			{ID: 1, Name: "Player", IsPlayer: true, Team: 0, HP: 0, MaxHP: 50, IsAlive: false},
		},
	}

	startingRoom := 1 // Foggy Gate
	result := combat.HandleDefeat(startingRoom)

	// Should use provided respawn room
	if result.RespawnPoint != startingRoom {
		t.Errorf("Expected respawn point %d, got %d", startingRoom, result.RespawnPoint)
	}

	// Corpse location should be combat room
	if result.CorpseLocation != combat.RoomID {
		t.Errorf("Expected corpse location %d, got %d", combat.RoomID, result.CorpseLocation)
	}
}