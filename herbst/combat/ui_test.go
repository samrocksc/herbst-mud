package combat

import (
	"strings"
	"testing"
	"time"
)

func TestNewCombatUI(t *testing.T) {
	ui := NewCombatUI(80, 24)
	
	if ui.width != 80 {
		t.Errorf("Expected width 80, got %d", ui.width)
	}
	
	if ui.height != 24 {
		t.Errorf("Expected height 24, got %d", ui.height)
	}
}

func TestCombatUI_RenderTitleBar(t *testing.T) {
	ui := NewCombatUI(80, 24)
	
	combat := &Combat{
		ID:           1,
		State:        StateActive,
		TickNumber:   5,
		Participants: []*Participant{},
	}
	
	titleBar := ui.renderTitleBar(combat)
	
	if !strings.Contains(titleBar, "COMBAT") {
		t.Error("Title bar should contain 'COMBAT'")
	}
}

func TestCombatUI_RenderEnemyHP(t *testing.T) {
	ui := NewCombatUI(80, 24)
	
	combat := &Combat{
		ID: 1,
		Participants: []*Participant{
			{ID: 1, Name: "Hero", IsPlayer: true, Team: 0, HP: 100, MaxHP: 100, IsAlive: true},
			{ID: 2, Name: "Goblin", IsNPC: true, Team: 1, HP: 30, MaxHP: 50, IsAlive: true},
		},
	}
	
	enemyHP := ui.renderEnemyHP(combat)
	
	if !strings.Contains(enemyHP, "Goblin") {
		t.Error("Enemy HP should contain enemy name")
	}
}

func TestCombatUI_RenderHPBar(t *testing.T) {
	ui := NewCombatUI(80, 24)
	
	// Full HP
	bar := ui.renderHPBar(100, 100, 10)
	if !strings.Contains(bar, "█") {
		t.Error("Full HP bar should contain filled blocks")
	}
	
	// Half HP
	bar = ui.renderHPBar(50, 100, 10)
	if !strings.Contains(bar, "█") || !strings.Contains(bar, "░") {
		t.Error("Half HP bar should contain both filled and empty blocks")
	}
	
	// Low HP
	bar = ui.renderHPBar(10, 100, 10)
	if !strings.Contains(bar, "░") {
		t.Error("Low HP bar should contain empty blocks")
	}
}

func TestCombatUI_RenderTurnOrder(t *testing.T) {
	ui := NewCombatUI(80, 24)
	
	participants := []*Participant{
		{ID: 1, Name: "Hero", IsPlayer: true, Team: 0, Initiative: 25, IsAlive: true},
		{ID: 2, Name: "Goblin", IsNPC: true, Team: 1, Initiative: 15, IsAlive: true},
	}
	
	combat := &Combat{
		ID:           1,
		Participants: participants,
	}
	
	// Roll initiative to set turn order
	combat.GetTurnOrder()
	
	turnOrder := ui.renderTurnOrder(combat)
	
	if !strings.Contains(turnOrder, "Hero") {
		t.Error("Turn order should contain 'Hero'")
	}
	
	if !strings.Contains(turnOrder, "Goblin") {
		t.Error("Turn order should contain 'Goblin'")
	}
}

func TestCombatUI_RenderCombatLog(t *testing.T) {
	ui := NewCombatUI(80, 24)
	
	combat := &Combat{
		ID: 1,
		Log: []CombatLogEntry{
			{Tick: 1, Message: "Hero attacks Goblin", Type: "damage"},
			{Tick: 2, Message: "Goblin takes 10 damage", Type: "damage"},
			{Tick: 3, Message: "Hero healed for 5", Type: "heal"},
		},
	}
	
	logOutput := ui.renderCombatLog(combat, 3)
	
	if !strings.Contains(logOutput, "Combat Log") {
		t.Error("Combat log should contain header")
	}
}

func TestCombatUI_RenderPlayerStatus(t *testing.T) {
	ui := NewCombatUI(80, 24)
	
	participants := []*Participant{
		{ID: 1, Name: "Hero", IsPlayer: true, Team: 0, HP: 80, MaxHP: 100, Mana: 50, MaxMana: 100, Stamina: 60, MaxStamina: 100},
	}
	
	combat := &Combat{
		ID:           1,
		Participants: participants,
	}
	
	status := ui.renderPlayerStatus(combat, 1)
	
	if !strings.Contains(status, "HP") {
		t.Error("Player status should contain HP")
	}
	
	if !strings.Contains(status, "Mana") {
		t.Error("Player status should contain Mana")
	}
	
	if !strings.Contains(status, "STA") {
		t.Error("Player status should contain Stamina")
	}
}

func TestCombatUI_RenderActionBar(t *testing.T) {
	ui := NewCombatUI(80, 24)
	
	combat := &Combat{
		ID: 1,
		Participants: []*Participant{
			{ID: 1, Name: "Hero", IsPlayer: true, Team: 0},
		},
	}
	
	im := NewInputManager(DefaultInputConfig())
	im.RegisterParticipant(combat.Participants[0])
	
	actionBar := ui.renderActionBar(combat, 1, im)
	
	if !strings.Contains(actionBar, "attack") && !strings.Contains(actionBar, "Attack") {
		t.Error("Action bar should contain attack action")
	}
}

func TestRenderCompactHPBar(t *testing.T) {
	// Full HP
	bar := RenderCompactHPBar(100, 100)
	if !strings.Contains(bar, "100/100") {
		t.Error("Compact HP bar should show HP values")
	}
	
	// Half HP
	bar = RenderCompactHPBar(50, 100)
	if !strings.Contains(bar, "50/100") {
		t.Error("Compact HP bar should show HP values")
	}
}

func TestFormatTimeRemaining(t *testing.T) {
	// 1 second
	formatted := FormatTimeRemaining(time.Second)
	if formatted != "1.0" {
		t.Errorf("Expected '1.0', got '%s'", formatted)
	}
	
	// 1.5 seconds
	formatted = FormatTimeRemaining(1500 * time.Millisecond)
	if formatted != "1.5" {
		t.Errorf("Expected '1.5', got '%s'", formatted)
	}
	
	// Negative (should return 0.0)
	formatted = FormatTimeRemaining(-time.Second)
	if formatted != "0.0" {
		t.Errorf("Expected '0.0' for negative, got '%s'", formatted)
	}
}

func TestCombatUI_Render_Full(t *testing.T) {
	ui := NewCombatUI(80, 24)
	
	participants := []*Participant{
		{ID: 1, Name: "Hero", IsPlayer: true, Team: 0, HP: 80, MaxHP: 100, Mana: 50, MaxMana: 100, Stamina: 60, MaxStamina: 100, Dexterity: 15, IsAlive: true},
		{ID: 2, Name: "Goblin", IsNPC: true, Team: 1, HP: 30, MaxHP: 50, Dexterity: 10, IsAlive: true},
	}
	
	combat := &Combat{
		ID:           1,
		RoomID:       5,
		Participants: participants,
		State:        StateActive,
		TickNumber:   1,
		ActionQueue:  NewActionQueue(),
		Effects:      NewEffectRegistry(),
	}
	
	tm := NewTickManager()
	cm := NewCombatManager()
	sm := NewCombatStateMachine(combat, cm, tm)
	im := NewInputManager(DefaultInputConfig())
	
	for _, p := range participants {
		im.RegisterParticipant(p)
	}
	
	// Roll initiative
	combat.GetTurnOrder()
	
	// Start input window
	sm.currentActorIdx = 0
	im.StartInputWindow(1)
	
	rendered := ui.Render(combat, sm, im, 1)
	
	// Should contain all major sections
	if !strings.Contains(rendered, "COMBAT") {
		t.Error("Full render should contain combat title")
	}
	
	if !strings.Contains(rendered, "Enemy") || !strings.Contains(rendered, "Goblin") {
		t.Error("Full render should contain enemy section")
	}
	
	if !strings.Contains(rendered, "Turn Order") {
		t.Error("Full render should contain turn order")
	}
	
	if !strings.Contains(rendered, "Combat Log") {
		t.Error("Full render should contain combat log")
	}
	
	if !strings.Contains(rendered, "Status") {
		t.Error("Full render should contain player status")
	}
	
	if !strings.Contains(rendered, "Action") {
		t.Error("Full render should contain action bar")
	}
}