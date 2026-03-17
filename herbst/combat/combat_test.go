package combat

import (
	"testing"
	"time"
)

func TestConfig_DefaultValues(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.TickIntervalMs != DefaultTickInterval {
		t.Errorf("expected default tick interval %d, got %d", DefaultTickInterval, cfg.TickIntervalMs)
	}
}

func TestNewTickLoop(t *testing.T) {
	tl := NewTickLoop(100) // 100ms for fast tests
	if tl == nil {
		t.Fatal("NewTickLoop returned nil")
	}
	if tl.interval != 100*time.Millisecond {
		t.Errorf("expected interval 100ms, got %v", tl.interval)
	}
	if tl.tick != 0 {
		t.Errorf("expected initial tick 0, got %d", tl.tick)
	}
	if tl.stopped {
		t.Error("expected stopped to be false initially")
	}
}

func TestTickLoop_GetTick(t *testing.T) {
	tl := NewTickLoop(1000)
	if tl.GetTick() != 0 {
		t.Error("expected initial tick 0")
	}
}

func TestTickLoop_GetInterval(t *testing.T) {
	tl := NewTickLoop(1500)
	expected := 1500 * time.Millisecond
	if tl.GetInterval() != expected {
		t.Errorf("expected interval %v, got %v", expected, tl.GetInterval())
	}
}

func TestTickLoop_TickChan(t *testing.T) {
	tl := NewTickLoop(1000)
	ch := tl.TickChan()
	if ch == nil {
		t.Error("TickChan returned nil")
	}
}

func TestTickLoop_RegisterUnregisterCombat(t *testing.T) {
	tl := NewTickLoop(1000)
	combat := &Combat{
		ID:           1,
		StartTick:    0,
		participants: make(map[int]*Participant),
		actionQueue:  make(map[int][]*Action),
		effects:      make(map[int][]*Effect),
	}

	tl.RegisterCombat(combat)
	if len(tl.combats) != 1 {
		t.Errorf("expected 1 combat registered, got %d", len(tl.combats))
	}

	tl.UnregisterCombat(1)
	if len(tl.combats) != 0 {
		t.Errorf("expected 0 combats after unregister, got %d", len(tl.combats))
	}
}

func TestTickLoop_String(t *testing.T) {
	tl := NewTickLoop(1500)
	s := tl.String()
	if s == "" {
		t.Error("String() returned empty string")
	}
	// Should contain TickLoop
	if len(s) < 8 || s[:8] != "TickLoop" {
		t.Errorf("expected String to start with TickLoop, got %s", s)
	}
}

func TestCombatManager_CreateCombat(t *testing.T) {
	tl := NewTickLoop(1000)
	m := NewCombatManager(tl)

	combatID := m.CreateCombat(0)
	if combatID != 1 {
		t.Errorf("expected first combat ID 1, got %d", combatID)
	}

	// Get the combat
	combat, ok := m.GetCombat(combatID)
	if !ok {
		t.Fatal("failed to get created combat")
	}
	if combat.ID != combatID {
		t.Errorf("expected combat ID %d, got %d", combatID, combat.ID)
	}
}

func TestCombatManager_GetCombat_NotFound(t *testing.T) {
	tl := NewTickLoop(1000)
	m := NewCombatManager(tl)

	_, ok := m.GetCombat(999)
	if ok {
		t.Error("expected GetCombat to return false for non-existent combat")
	}
}

func TestCombatManager_EndCombat(t *testing.T) {
	tl := NewTickLoop(1000)
	m := NewCombatManager(tl)

	combatID := m.CreateCombat(0)
	err := m.EndCombat(combatID)
	if err != nil {
		t.Errorf("EndCombat returned error: %v", err)
	}

	_, ok := m.GetCombat(combatID)
	if ok {
		t.Error("expected combat to be removed after EndCombat")
	}
}

func TestCombatManager_AddParticipant(t *testing.T) {
	tl := NewTickLoop(1000)
	m := NewCombatManager(tl)

	combatID := m.CreateCombat(0)
	participant := NewParticipant(1, "TestPlayer", 100)

	err := m.AddParticipant(combatID, participant)
	if err != nil {
		t.Errorf("AddParticipant returned error: %v", err)
	}

	combat, _ := m.GetCombat(combatID)
	p, ok := combat.GetParticipant(1)
	if !ok {
		t.Fatal("failed to get participant")
	}
	if p.Name != "TestPlayer" {
		t.Errorf("expected participant name TestPlayer, got %s", p.Name)
	}
	if p.HP != 100 {
		t.Errorf("expected HP 100, got %d", p.HP)
	}
}

func TestCombatManager_RemoveParticipant(t *testing.T) {
	tl := NewTickLoop(1000)
	m := NewCombatManager(tl)

	combatID := m.CreateCombat(0)
	participant := NewParticipant(1, "TestPlayer", 100)
	m.AddParticipant(combatID, participant)

	err := m.RemoveParticipant(combatID, 1)
	if err != nil {
		t.Errorf("RemoveParticipant returned error: %v", err)
	}

	combat, _ := m.GetCombat(combatID)
	_, ok := combat.GetParticipant(1)
	if ok {
		t.Error("expected participant to be removed")
	}
}

func TestNewParticipant(t *testing.T) {
	p := NewParticipant(1, "Warrior", 150)
	if p.ID != 1 {
		t.Errorf("expected ID 1, got %d", p.ID)
	}
	if p.Name != "Warrior" {
		t.Errorf("expected Name Warrior, got %s", p.Name)
	}
	if p.HP != 150 {
		t.Errorf("expected HP 150, got %d", p.HP)
	}
	if p.MaxHP != 150 {
		t.Errorf("expected MaxHP 150, got %d", p.MaxHP)
	}
	if p.Ready {
		t.Error("expected Ready to be false initially")
	}
}

func TestNewAction(t *testing.T) {
	a := NewAction(1, 2, 3, ActionAttack, 5)
	if a.ID != 1 {
		t.Errorf("expected ID 1, got %d", a.ID)
	}
	if a.Participant != 2 {
		t.Errorf("expected Participant 2, got %d", a.Participant)
	}
	if a.Target != 3 {
		t.Errorf("expected Target 3, got %d", a.Target)
	}
	if a.Type != ActionAttack {
		t.Errorf("expected Type ActionAttack, got %v", a.Type)
	}
	if a.Priority != 5 {
		t.Errorf("expected Priority 5, got %d", a.Priority)
	}
}

func TestNewEffect(t *testing.T) {
	e := NewEffect(1, "Poison", 3)
	if e.ID != 1 {
		t.Errorf("expected ID 1, got %d", e.ID)
	}
	if e.Name != "Poison" {
		t.Errorf("expected Name Poison, got %s", e.Name)
	}
	if e.Duration != 3 {
		t.Errorf("expected Duration 3, got %d", e.Duration)
	}
}

func TestCombat_QueueAction(t *testing.T) {
	combat := &Combat{
		ID:           1,
		StartTick:    0,
		participants: make(map[int]*Participant),
		actionQueue:  make(map[int][]*Action),
		effects:      make(map[int][]*Effect),
	}

	action := NewAction(1, 1, 2, ActionAttack, 1)
	combat.QueueAction(action)

	actions := combat.GetActions(1)
	if len(actions) != 1 {
		t.Errorf("expected 1 action in queue, got %d", len(actions))
	}
}

func TestCombat_ClearActions(t *testing.T) {
	combat := &Combat{
		ID:           1,
		StartTick:    0,
		participants: make(map[int]*Participant),
		actionQueue:  make(map[int][]*Action),
		effects:      make(map[int][]*Effect),
	}

	combat.QueueAction(NewAction(1, 1, 2, ActionAttack, 1))
	combat.ClearActions(1)

	actions := combat.GetActions(1)
	if len(actions) != 0 {
		t.Errorf("expected 0 actions after clear, got %d", len(actions))
	}
}

func TestCombat_GetParticipants(t *testing.T) {
	combat := &Combat{
		ID:           1,
		StartTick:    0,
		participants: make(map[int]*Participant),
		actionQueue:  make(map[int][]*Action),
		effects:      make(map[int][]*Effect),
	}

	p1 := NewParticipant(1, "Warrior", 100)
	p2 := NewParticipant(2, "Mage", 80)
	combat.participants[1] = p1
	combat.participants[2] = p2

	participants := combat.GetParticipants()
	if len(participants) != 2 {
		t.Errorf("expected 2 participants, got %d", len(participants))
	}
}

func TestCombat_AddEffect(t *testing.T) {
	combat := &Combat{
		ID:           1,
		StartTick:    0,
		participants: make(map[int]*Participant),
		actionQueue:  make(map[int][]*Action),
		effects:      make(map[int][]*Effect),
	}

	p := NewParticipant(1, "Warrior", 100)
	combat.participants[1] = p

	effect := NewEffect(1, "Poison", 3)
	applied := false
	effect.OnApply = func(p *Participant) {
		applied = true
	}

	combat.AddEffect(1, effect)

	effects := combat.GetEffects(1)
	if len(effects) != 1 {
		t.Errorf("expected 1 effect, got %d", len(effects))
	}
	if !applied {
		t.Error("expected OnApply to be called")
	}
}

func TestCombat_RemoveEffect(t *testing.T) {
	combat := &Combat{
		ID:           1,
		StartTick:    0,
		participants: make(map[int]*Participant),
		actionQueue:  make(map[int][]*Action),
		effects:      make(map[int][]*Effect),
	}

	p := NewParticipant(1, "Warrior", 100)
	combat.participants[1] = p

	effect := NewEffect(1, "Poison", 3)
	expired := false
	effect.OnExpire = func(p *Participant) {
		expired = true
	}

	combat.AddEffect(1, effect)
	combat.RemoveEffect(1, 1)

	effects := combat.GetEffects(1)
	if len(effects) != 0 {
		t.Errorf("expected 0 effects after removal, got %d", len(effects))
	}
	if !expired {
		t.Error("expected OnExpire to be called")
	}
}

func TestActionTypes(t *testing.T) {
	tests := []struct {
		at    ActionType
		expect string
	}{
		{ActionAttack, "attack"},
		{ActionDefend, "defend"},
		{ActionSkill, "skill"},
		{ActionItem, "item"},
		{ActionFlee, "flee"},
		{ActionWait, "wait"},
	}

	for _, test := range tests {
		if string(test.at) != test.expect {
			t.Errorf("expected %s, got %s", test.expect, string(test.at))
		}
	}
}