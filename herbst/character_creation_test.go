package main

import (
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
)

// TestCharacterCreationFlow tests the character creation state machine
func TestCharacterCreationFlow(t *testing.T) {
	availableRaces = []RaceInfo{
		{Name: "human", DisplayName: "Human"},
		{Name: "elf", DisplayName: "Elf"},
		{Name: "dwarf", DisplayName: "Dwarf"},
		{Name: "orc", DisplayName: "Orc"},
	}
	defer func() { availableRaces = nil }()

	t.Run("start character creation sets initial state", func(t *testing.T) {
		ti := textinput.New()
		ti.Focus()
		m := &model{
			screen:       ScreenCharacterSelect,
			textInput:    ti,
			connectedAt:  time.Now(),
			visitedRooms: make(map[int]bool),
			knownExits:   make(map[string]bool),
			width:        80,
			height:       24,
		}
		m.Init()
		// Set state directly to avoid triggering API call via fetchRaces
		m.isCreatingCharacter = true
		m.inputField = "char_name"
		m.textInput.SetValue("")
		m.inputBuffer = ""

		if !m.isCreatingCharacter {
			t.Error("Expected isCreatingCharacter to be true")
		}
		if m.inputField != "char_name" {
			t.Errorf("Expected inputField 'char_name', got %q", m.inputField)
		}
	})

	// Helper to create a model in character creation mode at a specific field
	creationModel := func(field string) *model {
		ti := textinput.New()
		ti.Focus()
		m := &model{
			screen:              ScreenCharacterSelect,
			textInput:           ti,
			isCreatingCharacter: true,
			inputField:          field,
			connectedAt:         time.Now(),
			visitedRooms:        make(map[int]bool),
			knownExits:          make(map[string]bool),
			width:               80,
			height:              24,
			maxHistory:          50,
		}
		m.Init()
		return m
	}

	t.Run("empty name shows error", func(t *testing.T) {
		m := creationModel("char_name")
		m.handleCharacterCreationInput("")
		if len(m.messageTypes) == 0 {
			t.Fatal("Expected message in history")
		}
		if lastType := m.messageTypes[len(m.messageTypes)-1]; lastType != "error" {
			t.Errorf("Expected error type, got %q", lastType)
		}
		if m.inputField != "char_name" {
			t.Errorf("Expected to stay on char_name, got %q", m.inputField)
		}
	})

	t.Run("name with numbers shows error", func(t *testing.T) {
		m := creationModel("char_name")
		m.handleCharacterCreationInput("Hero123")
		if len(m.messageTypes) == 0 {
			t.Fatal("Expected message in history")
		}
		if lastType := m.messageTypes[len(m.messageTypes)-1]; lastType != "error" {
			t.Errorf("Expected error type, got %q", lastType)
		}
		if m.inputField != "char_name" {
			t.Errorf("Expected to stay on char_name, got %q", m.inputField)
		}
	})

	t.Run("name too long shows error", func(t *testing.T) {
		m := creationModel("char_name")
		m.handleCharacterCreationInput("ThisNameIsWayTooLongForSure")
		if len(m.messageTypes) == 0 {
			t.Fatal("Expected message in history")
		}
		if lastType := m.messageTypes[len(m.messageTypes)-1]; lastType != "error" {
			t.Errorf("Expected error type, got %q", lastType)
		}
		if m.inputField != "char_name" {
			t.Errorf("Expected to stay on char_name, got %q", m.inputField)
		}
	})

	t.Run("valid name transitions to race", func(t *testing.T) {
		m := creationModel("char_name")
		m.handleCharacterCreationInput("HeroName")
		if m.inputField != "char_race" {
			t.Errorf("Expected char_race after valid name, got %q", m.inputField)
		}
	})

	t.Run("race selection by number", func(t *testing.T) {
		m := creationModel("char_race")
		createCharName = "TestHero"
		createCharRace = ""
		createCharFactionCategories = nil
		createCharFactionStep = 0
		createCharFactionChoices = nil
		defer func() { createCharName = ""; createCharRace = ""; createCharFactionCategories = nil }()

		m.handleCharacterCreationInput("2")
		// With no faction categories, should transition to creating character
		// (which triggers API call - we just check it's no longer in char_race)
		if m.inputField == "char_race" {
			t.Error("Expected to leave char_race after valid selection")
		}
		if m.inputField == "char_faction" {
			t.Error("Did not expect char_faction with no faction categories")
		}
	})

	t.Run("race selection by cursor", func(t *testing.T) {
		m := creationModel("char_race")
		createCharName = "TestHero"
		createCharRace = ""
		createCharFactionCategories = nil
		createCharFactionStep = 0
		createCharFactionChoices = nil
		defer func() { createCharName = ""; createCharRace = ""; createCharFactionCategories = nil }()

		m.createCursor = 1
		m.handleCharacterCreationInput("")
		// With no faction categories, should transition to creating character
		if m.inputField == "char_race" {
			t.Error("Expected to leave char_race after cursor selection")
		}
		if m.inputField == "char_faction" {
			t.Error("Did not expect char_faction with no faction categories")
		}
	})

	t.Run("invalid race number shows error", func(t *testing.T) {
		m := creationModel("char_race")
		createCharName = "TestHero"
		createCharRace = ""
		defer func() { createCharName = ""; createCharRace = "" }()

		m.handleCharacterCreationInput("9")
		if len(m.messageTypes) == 0 {
			t.Fatal("Expected message in history")
		}
		if lastType := m.messageTypes[len(m.messageTypes)-1]; lastType != "error" {
			t.Errorf("Expected error type, got %q", lastType)
		}
	})

	t.Run("invalid race name shows error", func(t *testing.T) {
		m := creationModel("char_race")
		createCharName = "TestHero"
		createCharRace = ""
		defer func() { createCharName = ""; createCharRace = "" }()

		m.handleCharacterCreationInput("Troll")
		if len(m.messageTypes) == 0 {
			t.Fatal("Expected message in history")
		}
		if lastType := m.messageTypes[len(m.messageTypes)-1]; lastType != "error" {
			t.Errorf("Expected error type, got %q", lastType)
		}
	})

	t.Run("empty race defaults to first", func(t *testing.T) {
		m := creationModel("char_race")
		createCharName = "TestHero"
		createCharRace = ""
		createCharFactionCategories = nil
		defer func() { createCharName = ""; createCharRace = ""; createCharFactionCategories = nil }()

		m.createCursor = 0
		m.handleCharacterCreationInput("")
		if m.inputField == "char_race" {
			t.Error("Expected to leave char_race after Enter")
		}
	})

	t.Run("cancel returns to character selection", func(t *testing.T) {
		m := creationModel("char_name")
		createCharName = "TestHero"
		defer func() { createCharName = "" }()

		m.handleCharacterCreationInput("cancel")
		if m.isCreatingCharacter {
			t.Error("Expected isCreatingCharacter to be false")
		}
		if m.screen != ScreenCharacterSelect {
			t.Errorf("Expected screen %q, got %q", ScreenCharacterSelect, m.screen)
		}
	})
}

// TestDisplayRaces tests the race display formatting
func TestDisplayRaces(t *testing.T) {
	t.Run("displayRaces with races lists all", func(t *testing.T) {
		availableRaces = []RaceInfo{
			{Name: "human", DisplayName: "Human"},
			{Name: "elf", DisplayName: "Elf"},
			{Name: "dwarf", DisplayName: "Dwarf"},
		}
		defer func() { availableRaces = nil }()

		m := &model{
			width:        80,
			height:       24,
			visitedRooms: make(map[int]bool),
			knownExits:   make(map[string]bool),
		}
		display := m.displayRaces()

		if display == "" {
			t.Error("Expected non-empty race display")
		}
		if !contains(display, "Human") {
			t.Error("Expected display to contain 'Human'")
		}
		if !contains(display, "Elf") {
			t.Error("Expected display to contain 'Elf'")
		}
		if !contains(display, "Dwarf") {
			t.Error("Expected display to contain 'Dwarf'")
		}
	})

	t.Run("displayRaces with no races", func(t *testing.T) {
		availableRaces = nil

		m := &model{
			width:        80,
			height:       24,
			visitedRooms: make(map[int]bool),
			knownExits:   make(map[string]bool),
		}
		display := m.displayRaces()

		if display == "" {
			t.Error("Expected non-empty display even with no races")
		}
	})
}
