package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// handleEmoteCommand performs a freeform action visible to the room
func (m *model) handleEmoteCommand(args []string) {
	if m.currentCharacterID == 0 {
		m.AppendMessage("You must be playing to emote.", "error")
		return
	}

	// Skip the first arg which is the command name "emote"
	if len(args) < 2 {
		m.AppendMessage("Usage: emote <action>", "error")
		return
	}

	action := strings.Join(args[1:], " ")
	if action == "" {
		m.AppendMessage("Usage: emote <action>", "error")
		return
	}

	// Build request body with character data
	bodyData := map[string]interface{}{
		"character_id": m.currentCharacterID,
		"action":       action,
	}
	body, _ := json.Marshal(bodyData)

	// Call the chat API with auth header
	req, err := http.NewRequest("POST", RESTAPIBase+"/api/chat/emote", bytes.NewBuffer(body))
	if err != nil {
		m.AppendMessage("Failed to emote.", "error")
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+m.characterToken)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		m.AppendMessage("Failed to emote.", "error")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		if msg, ok := result["error"].(string); ok {
			m.AppendMessage(msg, "error")
		} else {
			m.AppendMessage("Failed to emote.", "error")
		}
		return
	}

	m.AppendMessage(fmt.Sprintf("You %s", action), "info")
}
