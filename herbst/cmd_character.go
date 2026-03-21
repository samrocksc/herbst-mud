package main

import (
	"context"
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
	// Handle case when database client is not available
	if m.client == nil {
		m.currentCharacterName = m.currentUserName
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
		m.currentCharacterName = m.currentUserName
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

	char := chars[0]
	m.currentCharacterID = char.ID
	m.currentCharacterName = char.Name
	m.characterLevel = 1
	m.characterExperience = 0
}

func (m *model) saveProfileToDB() {
	if m.currentCharacterID == 0 {
		return
	}

	// Profile saved via REST API — stub for now, DB update deferred
	_ = m.currentCharacterID
	_ = m.characterGender
	_ = m.characterDescription
}
