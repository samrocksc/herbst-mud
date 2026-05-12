package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// handleReplyCommand responds to the last player who sent a tell
func (m *model) handleReplyCommand(args []string) {
	if m.currentCharacterID == 0 {
		m.AppendMessage("You must be playing to reply.", "error")
		return
	}

	if len(args) < 1 {
		m.AppendMessage("Usage: reply <message>", "error")
		return
	}

	if m.lastTeller == "" {
		m.AppendMessage("No one to reply to. Use tell <player> <message> instead.", "error")
		return
	}

	message := strings.Join(args, " ")
	if message == "" {
		m.AppendMessage("Usage: reply <message>", "error")
		return
	}

	// Check if ignored
	for _, ignored := range m.ignoredPlayers {
		if strings.EqualFold(ignored, m.lastTeller) {
			m.AppendMessage(fmt.Sprintf("You are ignoring %s.", m.lastTeller), "error")
			return
		}
	}

	// Call the chat API
	jsonData, _ := json.Marshal(map[string]interface{}{
		"target":  m.lastTeller,
		"message": message,
	})
	resp, err := http.Post(RESTAPIBase+"/api/chat/tell", "application/json", bytes.NewReader(jsonData))
	if err != nil {
		m.AppendMessage("Failed to send reply.", "error")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		if msg, ok := result["error"].(string); ok {
			m.AppendMessage(msg, "error")
		} else {
			m.AppendMessage("Failed to send reply.", "error")
		}
		return
	}

	m.AppendMessage(fmt.Sprintf("You reply to %s, \"%s\"", m.lastTeller, message), "info")
}
