package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// yellCooldownDuration is the cooldown between yell commands (5 seconds)
const yellCooldownDuration = 5 * time.Second

// handleYellCommand sends a message to everyone in the zone
func (m *model) handleYellCommand(args []string) {
	if m.currentCharacterID == 0 {
		m.AppendMessage("You must be playing to yell.", "error")
		return
	}

	if len(args) < 1 {
		m.AppendMessage("Usage: yell <message>", "error")
		return
	}

	message := strings.Join(args, " ")
	if message == "" {
		m.AppendMessage("Usage: yell <message>", "error")
		return
	}

	// Check cooldown
	if m.yellCooldown != nil && time.Since(*m.yellCooldown) < yellCooldownDuration {
		remaining := yellCooldownDuration - time.Since(*m.yellCooldown)
		m.AppendMessage(fmt.Sprintf("You must wait %v before yelling again.", remaining.Round(time.Second)), "error")
		return
	}

	// Build request body with character data
	bodyData := map[string]interface{}{
		"character_id": m.currentCharacterID,
		"room_id":      m.currentRoom,
		"message":      message,
	}
	body, _ := json.Marshal(bodyData)

	// Call the chat API with auth header
	req, err := http.NewRequest("POST", RESTAPIBase+"/api/chat/yell", bytes.NewBuffer(body))
	if err != nil {
		m.AppendMessage("Failed to yell.", "error")
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+m.characterToken)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		m.AppendMessage("Failed to yell.", "error")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		if msg, ok := result["error"].(string); ok {
			m.AppendMessage(msg, "error")
		} else {
			m.AppendMessage("Failed to yell.", "error")
		}
		return
	}

	// Set cooldown
	now := time.Now()
	m.yellCooldown = &now

	m.AppendMessage(fmt.Sprintf("You yell, \"%s\"", message), "info")
}
