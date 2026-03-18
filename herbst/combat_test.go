package main

import (
	"strings"
	"testing"
)

// TestCombatScreenState tests that combat state is properly managed
func TestCombatScreenState(t *testing.T) {
	m := &model{
		screen:            ScreenCombat,
		inCombat:          true,
		currentCharacterID: 1,
		characterName:    "TestPlayer",
		characterHP:      35,
		characterMaxHP:   45,
		characterStamina: 12,
		characterMaxStamina: 20,
		characterMana:    10,
		characterMaxMana: 20,
		enemyName:        "Old Scrap",
		enemyType:        "Mutant Raccoon",
		enemyHP:          24,
		enemyMaxHP:       45,
		currentTick:      4,
		channeledTalent:  "heavy_strike",
		channelingTicksLeft: 1,
		combatLog:        []string{"Old Scrap snarls and prepares to attack!"},
		equippedTalents:  []string{"slash", "parry", "heavy_strike", "second_wind"},
		talentCosts:      map[string]int{"slash": 0, "parry": 0, "heavy_strike": 2, "second_wind": 2},
	}

	// Test combat screen state
	if m.screen != ScreenCombat {
		t.Error("Expected ScreenCombat")
	}
	if !m.inCombat {
		t.Error("Expected inCombat to be true")
	}
	if m.enemyName != "Old Scrap" {
		t.Error("Expected enemy name Old Scrap")
	}
	if m.currentTick != 4 {
		t.Error("Expected current tick 4")
	}
}

// TestCombatScreenTransition tests transition to/from combat screen
func TestCombatScreenTransition(t *testing.T) {
	m := &model{
		screen:   ScreenPlaying,
		inCombat: false,
	}

	// Enter combat
	m.EnterCombat("Old Scrap", "Mutant Raccoon", 45, 24)
	if m.screen != ScreenCombat {
		t.Error("Expected screen to be ScreenCombat after entering combat")
	}
	if !m.inCombat {
		t.Error("Expected inCombat to be true")
	}
	if m.enemyName != "Old Scrap" {
		t.Error("Expected enemy name Old Scrap")
	}

	// Exit combat
	m.ExitCombat()
	if m.screen != ScreenPlaying {
		t.Error("Expected screen to be ScreenPlaying after exiting combat")
	}
	if m.inCombat {
		t.Error("Expected inCombat to be false after exiting combat")
	}
}

// TestCombatUIRendering tests the combat screen renders correctly
func TestCombatUIRendering(t *testing.T) {
	m := &model{
		screen:            ScreenCombat,
		inCombat:          true,
		currentCharacterID: 1,
		characterName:    "TestPlayer",
		characterHP:      35,
		characterMaxHP:   45,
		characterStamina: 12,
		characterMaxStamina: 20,
		characterMana:    10,
		characterMaxMana: 20,
		enemyName:        "Old Scrap",
		enemyType:        "Mutant Raccoon",
		enemyHP:          24,
		enemyMaxHP:       45,
		currentTick:      4,
		channeledTalent:  "heavy_strike",
		channelingTicksLeft: 1,
		combatLog:        []string{"Old Scrap snarls and prepares to attack!"},
		equippedTalents:  []string{"slash", "parry", "heavy_strike", "second_wind"},
		talentCosts:      map[string]int{"slash": 0, "parry": 0, "heavy_strike": 2, "second_wind": 2},
		width:            60,
		height:           24,
	}

	// Render combat view
	view := m.combatView()

	// Check for key elements
	if !strings.Contains(view, "COMBAT") {
		t.Error("Expected view to contain COMBAT")
	}
	if !strings.Contains(view, "Old Scrap") {
		t.Error("Expected view to contain enemy name")
	}
	if !strings.Contains(view, "Mutant Raccoon") {
		t.Error("Expected view to contain enemy type")
	}
	if !strings.Contains(view, "slash") {
		t.Error("Expected view to contain slash talent")
	}
	if !strings.Contains(view, "parry") {
		t.Error("Expected view to contain parry talent")
	}
	if !strings.Contains(view, "heavy_strike") {
		t.Error("Expected view to contain heavy_strike talent")
	}
}

// TestCombatStatusBar tests combat-specific status bar
func TestCombatStatusBar(t *testing.T) {
	m := &model{
		characterHP:       35,
		characterMaxHP:    45,
		characterStamina:  12,
		characterMaxStamina: 20,
		characterMana:     10,
		characterMaxMana:  20,
		channeledTalent:   "heavy_strike",
		channelingTicksLeft: 1,
	}

	bar := m.combatStatusBar()

	if !strings.Contains(bar, "35") || !strings.Contains(bar, "45") {
		t.Error("Expected bar to contain HP values")
	}
	if !strings.Contains(bar, "heavy_strike") {
		t.Error("Expected bar to show channeling talent")
	}
}

// TestEnemyHPBar tests enemy HP bar rendering
func TestEnemyHPBar(t *testing.T) {
	m := &model{}

	// 53% HP (24/45)
	bar := m.enemyHPBar(24, 45)
	if !strings.Contains(bar, "53%") {
		t.Error("Expected HP bar to show 53%")
	}

	// 100% HP
	bar = m.enemyHPBar(45, 45)
	if !strings.Contains(bar, "100%") {
		t.Error("Expected HP bar to show 100%")
	}

	// 10% HP
	bar = m.enemyHPBar(5, 50)
	if !strings.Contains(bar, "10%") {
		t.Error("Expected HP bar to show 10%")
	}
}

// TestCombatTalentSelection tests selecting talents in combat
func TestCombatTalentSelection(t *testing.T) {
	m := &model{
		inCombat:       true,
		equippedTalents: []string{"slash", "parry", "heavy_strike", "second_wind"},
		talentCosts:    map[string]int{"slash": 0, "parry": 0, "heavy_strike": 2, "second_wind": 2},
		characterStamina: 5,
		characterMana:   10,
	}

	// Test selecting slash (should work - 0 cost)
	selected := m.selectTalent(0) // First talent (index 0 = slash)
	if selected != "slash" {
		t.Errorf("Expected slash, got %s", selected)
	}

	// Test selecting heavy_strike with insufficient stamina
	m.characterStamina = 1
	selected = m.selectTalent(2) // heavy_strike
	if selected != "" {
		t.Error("Expected empty selection for insufficient stamina")
	}

	// Test selecting with sufficient resources
	m.characterStamina = 5
	selected = m.selectTalent(2)
	if selected != "heavy_strike" {
		t.Errorf("Expected heavy_strike, got %s", selected)
	}
}

// TestCombatTick tests combat tick advancement
func TestCombatTick(t *testing.T) {
	m := &model{
		inCombat:            true,
		currentTick:         4,
		channeledTalent:     "heavy_strike",
		channelingTicksLeft: 2,
		enemyName:           "Old Scrap",
		enemyHP:             45,  // Must be > 0 to avoid instant defeat
		enemyMaxHP:          45,
		characterLevel:      1,
		characterHP:         50,
		combatLog:           []string{},
	}

	// Advance tick
	m.AdvanceCombatTick()

	if m.currentTick != 5 {
		t.Errorf("Expected tick to advance to 5, got %d", m.currentTick)
	}
	if m.channelingTicksLeft != 1 {
		t.Errorf("Expected channeling ticks left to be 1, got %d", m.channelingTicksLeft)
	}

	// Another tick - should resolve
	m.AdvanceCombatTick()

	if m.channeledTalent != "" {
		t.Error("Expected channeled talent to clear after resolving")
	}
	// After second tick, we should have at least one log entry (enemy attack)
	if len(m.combatLog) < 1 {
		t.Errorf("Expected combat log to have entries after tick, got %d entries: %v", len(m.combatLog), m.combatLog)
	}
}

// TestCombatLog tests combat log management
func TestCombatLog(t *testing.T) {
	m := &model{
		combatLog: []string{},
	}

	// Add combat message
	m.AddCombatLog("Test combat message")
	if len(m.combatLog) != 1 {
		t.Error("Expected 1 combat log entry")
	}
	if !strings.Contains(m.combatLog[0], "Test combat message") {
		t.Error("Expected combat log to contain message")
	}

	// Add more messages (should cap at max)
	for i := 0; i < 20; i++ {
		m.AddCombatLog("Message " + string(rune(i+'0')))
	}

	// Should not exceed maxCombatLogSize
	if len(m.combatLog) > 10 {
		t.Errorf("Expected combat log to be capped at 10, got %d", len(m.combatLog))
	}
}