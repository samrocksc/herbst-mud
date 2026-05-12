package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// handleSayCommand sends a message to everyone in the same room
func (m *model) handleSayCommand(args []string) {
	if m.currentCharacterID == 0 {
		m.AppendMessage("You must be playing to speak.", "error")
		return
	}

	if len(args) < 1 {
		m.AppendMessage("Usage: say <message>", "error")
		return
	}

	message := strings.Join(args, " ")
	if message == "" {
		m.AppendMessage("Usage: say <message>", "error")
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
	req, err := http.NewRequest("POST", RESTAPIBase+"/api/chat/say", bytes.NewBuffer(body))
	if err != nil {
		m.AppendMessage("Failed to send message.", "error")
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+m.characterToken)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		m.AppendMessage("Failed to send message.", "error")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		if msg, ok := result["error"].(string); ok {
			m.AppendMessage(msg, "error")
		} else {
			m.AppendMessage("Failed to send message.", "error")
		}
		return
	}

	// Echo the message to self
	m.AppendMessage(fmt.Sprintf("You say, \"%s\"", message), "info")
}
