package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"herbst/db/character"
	"herbst/db/user"
)

// ============================================================
// PROFILE + WHOAMI
// ============================================================

func (m *model) handleProfileInput(input string) {
	input = strings.ToLower(input)
	switch input {
	case "1":
		m.editField = "gender"
		m.screen = ScreenEditField
		m.textInput.SetValue("")
		m.inputBuffer = ""
		m.message = ""
		m.messageType = ""
	case "2":
		m.editField = "description"
		m.screen = ScreenEditField
		m.textInput.SetValue("")
		m.inputBuffer = ""
		m.message = ""
		m.messageType = ""
	case "3", "back", "b", "esc":
		m.screen = ScreenPlaying
		m.message = ""
		m.messageType = ""
		m.menuItems = []string{"Login", "Register", "Quit"}
		m.menuCursor = 0
	default:
		m.message = "Invalid choice. Enter 1, 2, or 3"
		m.messageType = "error"
	}
}

func (m *model) handleEditFieldInput(input string) {
	if m.editField == "gender" {
		m.characterGender = input
		m.saveProfileToDB()
		m.message = "Gender updated!"
		m.messageType = "success"
	} else if m.editField == "description" {
		m.characterDescription = input
		m.saveProfileToDB()
		m.message = "Description updated!"
		m.messageType = "success"
	}
	m.screen = ScreenProfile
	m.textInput.SetValue("")
	m.inputBuffer = ""
}

// ============================================================
// CHARACTER
// ============================================================

func (m *model) loadOrCreateCharacter() {
	if m.client == nil {
		m.currentCharacterName = m.currentUserName
		m.characterRace = "human"
		m.characterGender = "unspecified"
		m.characterDescription = "A mysterious figure."
		m.characterHP = 100
		m.characterMaxHP = 100
		m.characterStamina = 50
		m.characterMaxStamina = 50
		m.characterMana = 25
		m.characterMaxMana = 25
		m.characterLevel = 1
		m.characterExperience = 0
		return
	}

	ctx := context.Background()
	chars, err := m.client.Character.Query().Where(character.HasUserWith(user.IDEQ(m.currentUserID))).All(ctx)
	if err != nil || len(chars) == 0 {
		// Check if world has races and genders before auto-creating
		if m.client != nil {
			m.fetchRaces()
			m.fetchGenders()
			if len(availableRaces) == 0 || len(availableGenders) == 0 {
				m.message = "This world is not ready for character creation. Please contact an admin."
				m.messageType = "error"
				m.screen = ScreenCharacterSelect
				return
			}
		}
		m.createDefaultCharacter()
		return
	}

	char := chars[0]
	m.currentCharacterID = char.ID
	m.currentCharacterName = char.Name
	m.characterRace = char.Race
	m.characterLevel = char.Level
	m.characterExperience = 0
	m.characterHP = char.Hitpoints
	m.characterMaxHP = char.MaxHitpoints
	m.characterStamina = char.Stamina
	m.characterMaxStamina = char.MaxStamina
	m.characterMana = char.Mana
	m.characterMaxMana = char.MaxMana
	m.isTest = char.IsTest
}

func (m *model) createDefaultCharacter() {
	jsonData, _ := json.Marshal(map[string]interface{}{
		"name":           m.currentUserName,
		"userId":         m.currentUserID,
		"isNPC":          false,
		"currentRoomId":  m.getRootRoomID(),
		"startingRoomId": m.getRootRoomID(),
		"race":           "human",
		"class":          "adventurer",
		"level":          1,
		"hitpoints":      100,
		"max_hitpoints":  100,
		"stamina":        50,
		"max_stamina":    50,
		"mana":           25,
		"max_mana":       25,
	})

	resp, err := http.Post(RESTAPIBase+"/characters", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		m.currentCharacterName = m.currentUserName
		m.setCharacterDefaults()
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		m.currentCharacterName = m.currentUserName
		m.setCharacterDefaults()
		return
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		m.currentCharacterName = m.currentUserName
		m.setCharacterDefaults()
		return
	}

	if id, ok := result["id"].(float64); ok {
		m.currentCharacterID = int(id)
	}
	if race, ok := result["race"].(string); ok {
		m.characterRace = race
	}
	m.currentCharacterName = m.currentUserName
	m.setCharacterDefaults()
}

func (m *model) setCharacterDefaults() {
	if m.characterRace == "" {
		m.characterRace = "human"
	}
	m.characterGender = "unspecified"
	m.characterDescription = "A mysterious figure."
	m.characterHP = 100
	m.characterMaxHP = 100
	m.characterStamina = 50
	m.characterMaxStamina = 50
	m.characterMana = 25
	m.characterMaxMana = 25
	m.characterLevel = 1
	m.characterExperience = 0
}

func (m *model) saveProfileToDB() {
	if m.currentCharacterID == 0 {
		return
	}
	_ = m.currentCharacterID
	_ = m.characterGender
	_ = m.characterDescription
}