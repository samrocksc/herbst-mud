package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"herbst/effects"
)

// TestCharacterCreationE2E tests the full character creation flow from
// start to finish, including API interactions via a mock server.
func TestCharacterCreationE2E(t *testing.T) {
	// Track API calls made by the model
	var createCharReq struct {
		body map[string]interface{}
		url  string
	}
	var loadCharCalled bool

	// Start a mock API server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET" && r.URL.Path == "/playable-races":
			json.NewEncoder(w).Encode(map[string]interface{}{
				"races": []map[string]string{
					{"name": "human", "display_name": "Human"},
					{"name": "elf", "display_name": "Elf"},
					{"name": "dwarf", "display_name": "Dwarf"},
				},
			})

		case r.Method == "GET" && r.URL.Path == "/genders":
			json.NewEncoder(w).Encode([]map[string]string{
				{"name": "he_him", "display_name": "He/Him", "subject_pronoun": "he", "object_pronoun": "him", "possessive_pronoun": "his"},
				{"name": "she_her", "display_name": "She/Her", "subject_pronoun": "she", "object_pronoun": "her", "possessive_pronoun": "hers"},
				{"name": "they_them", "display_name": "They/Them", "subject_pronoun": "they", "object_pronoun": "them", "possessive_pronoun": "theirs"},
			})

		case r.Method == "GET" && r.URL.Path == "/api/faction-categories":
			// No factions configured — return empty
			json.NewEncoder(w).Encode([]map[string]interface{}{})

		case r.Method == "POST" && strings.HasPrefix(r.URL.Path, "/user-characters/"):
			// Record the request
			createCharReq.url = r.URL.Path
			json.NewDecoder(r.Body).Decode(&createCharReq.body)

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id": float64(42),
				"name": createCharReq.body["name"],
			})

		case r.Method == "GET" && r.URL.Path == "/characters/42":
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":            float64(42),
				"name":          "HeroName",
				"race":          "elf",
				"class":         "warrior",
				"gender":        "",
				"description":   "",
				"hitpoints":     float64(100),
				"max_hitpoints": float64(100),
				"stamina":       float64(50),
				"max_stamina":   float64(50),
				"mana":          float64(20),
				"max_mana":      float64(20),
				"level":         float64(1),
				"currentRoomId": float64(5),
				"respawnRoomId": float64(5),
				"lastSeenAt":    nil,
			})
			loadCharCalled = true

		case r.Method == "PUT" && strings.HasPrefix(r.URL.Path, "/characters/"):
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{})

		default:
			t.Logf("Unexpected API call: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "unexpected"})
		}
	}))
	defer ts.Close()

	// Override the API base URL
	origBase := RESTAPIBase
	RESTAPIBase = ts.URL
	defer func() { RESTAPIBase = origBase }()

	// Setup model
	ti := textinput.New()
	ti.Focus()

	m := &model{
		screen:              ScreenCharacterSelect,
		textInput:           ti,
		currentUserID:       1,
		currentWorld:        "TestWorld",
		connectedAt:         time.Now(),
		visitedRooms:        make(map[int]bool),
		knownExits:          make(map[string]bool),
		width:               80,
		height:              24,
		maxHistory:          50,
		isCreatingCharacter: false,
		effectsService:      effects.NewService(ts.URL, slog.Default()),
	}
	m.Init()

	// --- Step 1: Start character creation ---
	m.startCharacterCreation()

	if !m.isCreatingCharacter {
		t.Fatal("Step 1: Expected isCreatingCharacter after startCharacterCreation")
	}
	if m.inputField != "char_name" {
		t.Fatalf("Step 1: Expected inputField 'char_name', got %q", m.inputField)
	}
	if len(availableRaces) == 0 {
		t.Fatal("Step 1: Expected races to be loaded")
	}

	// --- Step 2: Enter character name ---
	m.handleCharacterCreationInput("HeroName")

	if m.inputField != "char_race" {
		t.Fatalf("Step 2: Expected transition to char_race, got %q", m.inputField)
	}

	// --- Step 3: Select race via cursor ---
	// Elf is at index 1 (0=Human, 1=Elf, 2=Dwarf)
	m.createCursor = 1
	m.handleCharacterCreationInput("") // Enter = select highlighted

	// Should transition to char_gender
	if m.inputField != "char_gender" {
		t.Fatalf("Step 3: Expected transition to char_gender, got %q", m.inputField)
	}

	// --- Step 4: Select gender via cursor ---
	// She/Her is at index 1
	m.createCursor = 1
	m.handleCharacterCreationInput("") // Enter = select highlighted

	if m.inputField != "char_description" {
		t.Fatalf("Step 4: Expected transition to char_description, got %q", m.inputField)
	}

	// --- Step 5: Enter description ---
	m.handleCharacterCreationInput("A brave elven warrior.")

	// Step 5 should trigger createCharacter (no factions configured)
	if createCharReq.url == "" {
		t.Fatal("Step 5: Expected createCharacter API call after description")
	}
	expectedPath := "/user-characters/1"
	if createCharReq.url != expectedPath {
		t.Errorf("Step 5: Expected API path %q, got %q", expectedPath, createCharReq.url)
	}

	// Verify the create payload
	if name, ok := createCharReq.body["name"].(string); !ok || name != "HeroName" {
		t.Errorf("Step 5: Expected character name 'HeroName', got %v", createCharReq.body["name"])
	}
	if race, ok := createCharReq.body["race"].(string); !ok || race != "elf" {
		t.Errorf("Step 5: Expected race 'elf', got %v", createCharReq.body["race"])
	}
	if gender, ok := createCharReq.body["gender"].(string); !ok || gender != "she_her" {
		t.Errorf("Step 5: Expected gender 'she_her', got %v", createCharReq.body["gender"])
	}
	if desc, ok := createCharReq.body["description"].(string); !ok || desc != "A brave elven warrior." {
		t.Errorf("Step 5: Expected description 'A brave elven warrior.', got %v", createCharReq.body["description"])
	}
	if world, ok := createCharReq.body["world"].(string); !ok || world != "TestWorld" {
		t.Errorf("Step 5: Expected world 'TestWorld', got %v", createCharReq.body["world"])
	}

	// --- Step 6: Verify loadCharacter was called and we're in playing screen ---
	if !loadCharCalled {
		t.Error("Step 6: Expected loadCharacter to call GET /characters/42")
	}
	if m.screen != ScreenPlaying {
		t.Errorf("Step 6: Expected screen 'playing', got %q", m.screen)
	}
	if m.currentCharacterName != "HeroName" {
		t.Errorf("Step 6: Expected character name 'HeroName', got %q", m.currentCharacterName)
	}
	if m.characterRace != "elf" {
		t.Errorf("Step 6: Expected race 'elf', got %q", m.characterRace)
	}
	if m.characterClass != "warrior" {
		t.Errorf("Step 6: Expected class 'warrior', got %q", m.characterClass)
	}
	if m.characterLevel != 1 {
		t.Errorf("Step 6: Expected level 1, got %d", m.characterLevel)
	}
	if m.characterHP != 100 {
		t.Errorf("Step 6: Expected HP 100, got %d", m.characterHP)
	}
	if m.isCreatingCharacter {
		t.Error("Step 6: Expected isCreatingCharacter to be false after load")
	}
}

// TestCharacterCreationE2EView tests that the character creation screens
// render without panic and show the expected content.
func TestCharacterCreationE2EView(t *testing.T) {
	ti := textinput.New()
	ti.Focus()

	availableRaces = []RaceInfo{
		{Name: "human", DisplayName: "Human"},
		{Name: "elf", DisplayName: "Elf"},
		{Name: "dwarf", DisplayName: "Dwarf"},
	}
	defer func() { availableRaces = nil }()

	t.Run("character select list shows cursor", func(t *testing.T) {
		m := &model{
			screen: ScreenCharacterSelect,
			selectedWorldCharacters: []CharacterInfo{
				{ID: 1, Name: "Aragorn", Level: 5, Race: "Human", Class: "Warrior", Hitpoints: 80, MaxHitpoints: 100},
				{ID: 2, Name: "Legolas", Level: 3, Race: "Elf", Class: "Archer", Hitpoints: 60, MaxHitpoints: 70},
			},
			characterCursor: 1, // Highlight Legolas
			textInput:       ti,
			currentWorld:    "TestWorld",
			connectedAt:     time.Now(),
			visitedRooms:    make(map[int]bool),
			knownExits:      make(map[string]bool),
			width:           80,
			height:          24,
			maxHistory:      50,
		}
		m.Init()

		view := m.View()

		if !contains(view, "▸") {
			t.Error("Expected cursor (▸) in character list view")
		}
		if !contains(view, "Legolas") {
			t.Error("Expected 'Legolas' in view (should be highlighted)")
		}
		if !contains(view, "Aragorn") {
			t.Error("Expected 'Aragorn' in view")
		}
		if !contains(view, "Create new character") {
			t.Error("Expected 'Create new character' option in view")
		}
		if !contains(view, "j/k navigate") {
			t.Error("Expected navigation hint in view")
		}
		if contains(view, "> ") {
			t.Error("Did not expect text input prompt in character select")
		}
	})

	t.Run("race selection shows cursor", func(t *testing.T) {
		m := &model{
			screen:              ScreenCharacterSelect,
			isCreatingCharacter: true,
			inputField:          "char_race",
			createCursor:        2, // Highlight Dwarf
			textInput:           ti,
			connectedAt:         time.Now(),
			visitedRooms:        make(map[int]bool),
			knownExits:          make(map[string]bool),
			width:               80,
			height:              24,
			maxHistory:          50,
		}
		m.Init()

		view := m.View()

		if !contains(view, "▸") {
			t.Error("Expected cursor (▸) in race selection view")
		}
		if !contains(view, "Dwarf") {
			t.Error("Expected 'Dwarf' in view (should be highlighted)")
		}
		if !contains(view, "Human") || !contains(view, "Elf") {
			t.Error("Expected all races in view")
		}
	})
}
