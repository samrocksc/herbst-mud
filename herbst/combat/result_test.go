package combat

import (
	"testing"
	"time"
)

func TestHandleVictory_Basic(t *testing.T) {
	combat := &Combat{
		ID:    1,
		State: StateActive,
		Participants: []*Participant{
			{ID: 1, Name: "Player", IsPlayer: true, Team: 0, HP: 50, MaxHP: 50, IsAlive: true},
			{ID: 2, Name: "Scrap Rat", IsNPC: true, Team: 1, HP: 0, MaxHP: 20, Level: 1, IsAlive: false},
		},
		Log: []CombatLogEntry{
			{Tick: 1, Message: "Player attacks Scrap Rat", Type: "action"},
			{Tick: 2, Message: "Scrap Rat takes 15 damage", Type: "damage"},
		},
	}

	result := combat.HandleVictory()

	if result == nil {
		t.Fatal("Expected victory result, got nil")
	}

	if result.XP < 1 {
		t.Errorf("Expected XP reward, got %d", result.XP)
	}

	if len(result.Loot) == 0 {
		t.Error("Expected loot items from defeated enemy")
	}

	if result.VictoryTime.IsZero() {
		t.Error("Expected victory time to be set")
	}
}

func TestHandleVictory_MultipleEnemies(t *testing.T) {
	combat := &Combat{
		ID:    1,
		State: StateActive,
		Participants: []*Participant{
			{ID: 1, Name: "Player", IsPlayer: true, Team: 0, HP: 50, MaxHP: 50, IsAlive: true},
			{ID: 2, Name: "Rat A", IsNPC: true, Team: 1, HP: 0, MaxHP: 15, Level: 1, IsAlive: false},
			{ID: 3, Name: "Rat B", IsNPC: true, Team: 1, HP: 0, MaxHP: 15, Level: 1, IsAlive: false},
			{ID: 4, Name: "Rat C", IsNPC: true, Team: 1, HP: 0, MaxHP: 15, Level: 2, IsAlive: false},
		},
		Log: []CombatLogEntry{},
	}

	result := combat.HandleVictory()

	// Should get XP from all defeated enemies
	if result.XP < 3 {
		t.Errorf("Expected XP from 3 enemies, got %d", result.XP)
	}

	// Should have loot from all enemies
	if len(result.Loot) < 1 {
		t.Error("Expected loot from multiple enemies")
	}
}

func TestHandleVictory_WeaponDrop(t *testing.T) {
	combat := &Combat{
		ID:    1,
		State: StateActive,
		Participants: []*Participant{
			{ID: 1, Name: "Player", IsPlayer: true, Team: 0, HP: 50, MaxHP: 50, IsAlive: true},
			{ID: 2, Name: "Old Scrap", IsNPC: true, Team: 1, HP: 0, MaxHP: 100, Level: 5, IsAlive: false,
				HasGuaranteedDrop: true, DropWeapon: "Scrap Machete"},
		},
		Log: []CombatLogEntry{},
	}

	result := combat.HandleVictory()

	// Should have a weapon drop
	if len(result.WeaponDrops) == 0 {
		t.Fatal("Expected weapon drop from guaranteed enemy")
	}

	drop := result.WeaponDrops[0]
	if drop.Name != "Scrap Machete" {
		t.Errorf("Expected 'Scrap Machete', got '%s'", drop.Name)
	}

	if !drop.Guaranteed {
		t.Error("Expected guaranteed drop")
	}
}

func TestHandleDefeat_Basic(t *testing.T) {
	combat := &Combat{
		ID:     1,
		RoomID: 5, // Junkyard entrance
		State:  StateActive,
		Participants: []*Participant{
			{ID: 1, Name: "Player", IsPlayer: true, Team: 0, HP: 0, MaxHP: 50, IsAlive: false},
			{ID: 2, Name: "Boss", IsNPC: true, Team: 1, HP: 30, MaxHP: 100, Level: 10, IsAlive: true},
		},
		Log: []CombatLogEntry{
			{Tick: 5, Message: "Boss attacks Player", Type: "action"},
			{Tick: 5, Message: "Player takes 50 damage", Type: "damage"},
		},
	}

	startingRoom := 1 // Foggy Gate
	result := combat.HandleDefeat(startingRoom)

	if result == nil {
		t.Fatal("Expected defeat result, got nil")
	}

	if result.RespawnPoint != startingRoom {
		t.Errorf("Expected respawn at room %d, got %d", startingRoom, result.RespawnPoint)
	}

	if result.CorpseLocation != combat.RoomID {
		t.Errorf("Expected corpse at room %d, got %d", combat.RoomID, result.CorpseLocation)
	}

	if len(result.Consequences) == 0 {
		t.Error("Expected defeat consequences")
	}

	if result.DefeatTime.IsZero() {
		t.Error("Expected defeat time to be set")
	}
}

func TestHandleDefeat_CorpseExpiry(t *testing.T) {
	combat := &Combat{
		ID:     1,
		RoomID: 10,
		State:  StateActive,
		Participants: []*Participant{
			{ID: 1, Name: "Player", IsPlayer: true, Team: 0, HP: 0, MaxHP: 50, IsAlive: false},
		},
	}

	result := combat.HandleDefeat(1)

	// Corpse should expire in 5 minutes
	if result.CorpseExpiry.IsZero() {
		t.Fatal("Expected corpse expiry to be set")
	}

	// Should be approximately 5 minutes in the future
	expectedExpiry := result.DefeatTime.Add(5 * time.Minute)
	diff := result.CorpseExpiry.Sub(expectedExpiry)
	if diff > time.Second {
		t.Errorf("Expected corpse expiry ~5 minutes from defeat, got %v", diff)
	}
}

func TestGenerateNPCLoot_CommonItems(t *testing.T) {
	npc := &Participant{
		ID:     2,
		Name:   "Rat",
		Level:  1,
		IsNPC:  true,
	}

	// Run multiple times to check randomness
	for i := 0; i < 10; i++ {
		loot := generateNPCLoot(npc)
		if len(loot) == 0 {
			t.Error("Expected at least one loot item")
		}

		for _, item := range loot {
			if item.Quantity < 1 {
				t.Error("Expected item quantity >= 1")
			}
		}
	}
}

func TestCalculateEnemyXP_BasicEnemy(t *testing.T) {
	npc := &Participant{
		ID:     2,
		Name:   "Scrap Rat",
		Level:  1,
		MaxHP:  15,
		IsNPC:  true,
	}

	xp := calculateEnemyXP(npc)

	// Level 1 enemy = 15 base XP
	if xp < 15 {
		t.Errorf("Expected at least 15 XP for level 1 enemy, got %d", xp)
	}
}

func TestCalculateEnemyXP_ToughEnemy(t *testing.T) {
	npc := &Participant{
		ID:     2,
		Name:   "Old Scrap",
		Level:  5,
		MaxHP:  120, // Tough enemy
		IsNPC:  true,
	}

	xp := calculateEnemyXP(npc)

	// Level 5 = 75 base + 20 (HP>50) + 30 (HP>100) = 125
	if xp < 75 {
		t.Errorf("Expected at least 75 XP for level 5 tough enemy, got %d", xp)
	}
}

func TestGetWeaponType(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"Rusty Sword", "sword"},
		{"Scrap Machete", "sword"},
		{"Twisted Pipe", "pipe"},
		{"Chef's Knife", "dagger"},
		{"Broken Bottle", "dagger"},
		{"Unknown Weapon", "weapon"},
	}

	for _, tt := range tests {
		result := getWeaponType(tt.name)
		if result != tt.expected {
			t.Errorf("getWeaponType(%q) = %q, want %q", tt.name, result, tt.expected)
		}
	}
}

func TestGetClassRestriction(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"Rusty Sword", "warrior"},
		{"Scrap Machete", "warrior"},
		{"Twisted Pipe", "chef"},
		{"Chef's Knife", "chef"},
		{"Unknown Weapon", ""}, // No restriction
	}

	for _, tt := range tests {
		result := getClassRestriction(tt.name)
		if result != tt.expected {
			t.Errorf("getClassRestriction(%q) = %q, want %q", tt.name, result, tt.expected)
		}
	}
}

func TestSkillUps(t *testing.T) {
	combat := &Combat{
		ID:    1,
		State: StateActive,
		Participants: []*Participant{
			{ID: 1, Name: "Player", IsPlayer: true, Team: 0, HP: 50, MaxHP: 50, IsAlive: true},
		},
		Log: []CombatLogEntry{
			{Tick: 1, Message: "Player attacks", Type: "action"},
		},
	}

	skillUps := calculateSkillUps(combat)

	if len(skillUps) == 0 {
		t.Error("Expected skill ups from combat")
	}

	// Should include brawling for basic attacks
	foundBrawling := false
	for _, su := range skillUps {
		if su.SkillName == "brawling" {
			foundBrawling = true
			if su.Amount < 1 {
				t.Errorf("Expected positive skill amount, got %d", su.Amount)
			}
		}
	}

	if !foundBrawling {
		t.Error("Expected brawling skill up for basic attacks")
	}
}

func TestVictoryResult_LootRarity(t *testing.T) {
	combat := &Combat{
		ID:    1,
		State: StateActive,
		Participants: []*Participant{
			{ID: 1, Name: "Player", IsPlayer: true, Team: 0, HP: 50, MaxHP: 50, IsAlive: true},
			{ID: 2, Name: "Enemy", IsNPC: true, Team: 1, HP: 0, MaxHP: 20, Level: 1, IsAlive: false},
		},
	}

	// Run multiple times to check for different rarities
	rarities := make(map[string]int)
	for i := 0; i < 100; i++ {
		result := combat.HandleVictory()
		for _, item := range result.Loot {
			rarities[item.Rarity]++
		}
	}

	// Should see common items most often
	if rarities["common"] == 0 {
		t.Error("Expected common loot items")
	}

	// Rare items should be rare (if we see them)
	// Note: Due to randomness, we might not see rare items in every test run
}

func TestDefeatResult_XPLoss(t *testing.T) {
	combat := &Combat{
		ID:     1,
		RoomID: 5,
		State:  StateActive,
		Participants: []*Participant{
			{ID: 1, Name: "Player", IsPlayer: true, Team: 0, HP: 0, MaxHP: 50, IsAlive: false},
		},
	}

	result := combat.HandleDefeat(1)

	// XP loss should be 0 since we can't calculate it without player data
	// In a real implementation, this would be 10% of current level XP
	if result.XPLost < 0 {
		t.Error("XP loss should not be negative")
	}
}