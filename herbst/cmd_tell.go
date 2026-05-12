package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// handleTellCommand sends a direct message to a specific player
func (m *model) handleTellCommand(args []string) {
	if m.currentCharacterID == 0 {
		m.AppendMessage("You must be playing to tell.", "error")
		return
	}

	// Skip the first arg which is the command name "tell"
	if len(args) < 2 {
		m.AppendMessage("Usage: tell <player> <message>", "error")
		return
	}

	// args[0] is "tell", args[1] is target, args[2:] is message
	targetPlayer := args[1]
	message := strings.Join(args[2:], " ")

	if targetPlayer == "" || message == "" {
		m.AppendMessage("Usage: tell <player> <message>", "error")
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

	// Debug: log what we're about to send
	if m.debugMode {
		m.debugLogf("tell: from_id=%d, target_name=%s, message=%s", m.currentCharacterID, targetPlayer, message)
		m.debugLogf("tell: body JSON: %s", string(body))
	}

	// Call the chat API with auth header
	req, err := http.NewRequest("POST", RESTAPIBase+"/api/chat/tell", bytes.NewBuffer(body))
	if err != nil {
		m.AppendMessage("Failed to send tell.", "error")
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+m.characterToken)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		m.AppendMessage("Failed to send tell.", "error")
		if m.debugMode {
			m.debugLogf("tell: HTTP request error: %v", err)
		}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		if msg, ok := result["error"].(string); ok {
			m.AppendMessage(msg, "error")
			if m.debugMode {
				m.debugLogf("tell: API error: %s", msg)
			}
		} else {
			m.AppendMessage("Failed to send tell.", "error")
			if m.debugMode {
				m.debugLogf("tell: unexpected response: %v", result)
			}
		}
		return
	}

	// Store as lastTeller for reply
	m.lastTeller = targetPlayer

	// Use the display message from the API result if available
	var response struct {
		DisplayMessage string `json:"display_message,omitempty"`
	}
	json.NewDecoder(resp.Body).Decode(&response)

	if response.DisplayMessage != "" {
		m.AppendMessage(response.DisplayMessage, "info")
	} else {
		// Fallback
		m.AppendMessage(fmt.Sprintf("You tell %s, \"%s\"", targetPlayer, message), "info")
	}
}
