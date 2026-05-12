package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// handleWhisperCommand sends a private message to a player in the same room
func (m *model) handleWhisperCommand(args []string) {
	if m.currentCharacterID == 0 {
		m.AppendMessage("You must be playing to whisper.", "error")
		return
	}

	// Skip the first arg which is the command name "whisper"
	if len(args) < 2 {
		m.AppendMessage("Usage: whisper <player> <message>", "error")
		return
	}

	// args[0] is "whisper", args[1] is target, args[2:] is message
	targetPlayer := args[1]
	message := strings.Join(args[2:], " ")

	if targetPlayer == "" || message == "" {
		m.AppendMessage("Usage: whisper <player> <message>", "error")
		return
	}

	// Check if ignored
	for _, ignored := range m.ignoredPlayers {
		if strings.EqualFold(ignored, targetPlayer) {
			m.AppendMessage(fmt.Sprintf("You are ignoring %s.", targetPlayer), "error")
			return
		}
	}

	// Build request body with character data
	bodyData := map[string]interface{}{
		"from_id": m.currentCharacterID,
		"to_id":   0, // Target will be resolved by name in handler
		"message": message,
	}
	// Add target name for resolver
	bodyData["target_name"] = targetPlayer
	body, _ := json.Marshal(bodyData)

	// Call the chat API with auth header
	req, err := http.NewRequest("POST", RESTAPIBase+"/api/chat/whisper", bytes.NewBuffer(body))
	if err != nil {
		m.AppendMessage("Failed to whisper.", "error")
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+m.characterToken)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		m.AppendMessage("Failed to whisper.", "error")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		if msg, ok := result["error"].(string); ok {
			m.AppendMessage(msg, "error")
		} else {
			m.AppendMessage("Failed to whisper.", "error")
		}
		return
	}

	m.AppendMessage(fmt.Sprintf("You whisper to %s, \"%s\"", targetPlayer, message), "info")
}
