package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// shoutCooldownDuration is the cooldown between shout commands (30 seconds)
const shoutCooldownDuration = 30 * time.Second

// shoutStaminaCost is the stamina cost for shout
const shoutStaminaCost = 5

// handleShoutCommand sends a message to everyone in the world
func (m *model) handleShoutCommand(args []string) {
	if m.currentCharacterID == 0 {
		m.AppendMessage("You must be playing to shout.", "error")
		return
	}

	if len(args) < 1 {
		m.AppendMessage("Usage: shout <message>", "error")
		return
	}

	message := strings.Join(args, " ")
	if message == "" {
		m.AppendMessage("Usage: shout <message>", "error")
		return
	}

	// Check stamina
	if m.characterStamina < shoutStaminaCost {
		m.AppendMessage(fmt.Sprintf("You need %d stamina to shout. You only have %d.", shoutStaminaCost, m.characterStamina), "error")
		return
	}

	// Check cooldown
	if m.shoutCooldown != nil && time.Since(*m.shoutCooldown) < shoutCooldownDuration {
		remaining := shoutCooldownDuration - time.Since(*m.shoutCooldown)
		m.AppendMessage(fmt.Sprintf("You must wait %v before shouting again.", remaining.Round(time.Second)), "error")
		return
	}

	// Build request body with character data
	bodyData := map[string]interface{}{
		"character_id": m.currentCharacterID,
		"message":      message,
	}
	body, _ := json.Marshal(bodyData)

	// Call the chat API with auth header
	req, err := http.NewRequest("POST", RESTAPIBase+"/api/chat/shout", bytes.NewBuffer(body))
	if err != nil {
		m.AppendMessage("Failed to shout.", "error")
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+m.characterToken)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		m.AppendMessage("Failed to shout.", "error")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		if msg, ok := result["error"].(string); ok {
			m.AppendMessage(msg, "error")
		} else {
			m.AppendMessage("Failed to shout.", "error")
		}
		return
	}

	// Deduct stamina
	m.characterStamina -= shoutStaminaCost

	// Set cooldown
	now := time.Now()
	m.shoutCooldown = &now

	m.AppendMessage(fmt.Sprintf("You shout to the world, \"%s\"", message), "info")
}
