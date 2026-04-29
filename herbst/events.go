package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// adminEventsToken is cached after first fetch.
var (
	adminEventsToken string
	tokenMu          sync.Mutex // protects adminEventsToken
)

// fetchAdminEventsToken logs in as the admin user and caches the JWT.
// It reads ADMIN_EMAIL and ADMIN_PASSWORD from env vars, falling back to
// the server's hardcoded default credentials.
func fetchAdminEventsToken() (string, error) {
	tokenMu.Lock()
	defer tokenMu.Unlock()
	if adminEventsToken != "" {
		return adminEventsToken, nil
	}

	email := os.Getenv("ADMIN_EMAIL")
	password := os.Getenv("ADMIN_PASSWORD")
	if email == "" {
		email = "admin@herbstmud.local"
	}
	if password == "" {
		password = "herb5t2026!"
	}

	jsonData, _ := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
	})
	resp, err := http.Post(RESTAPIBase+"/users/auth", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to auth for events token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("auth failed with status %d", resp.StatusCode)
	}

	var result struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode auth response: %w", err)
	}
	adminEventsToken = result.Token
	return adminEventsToken, nil
}

// FireDefeatEvent sends an npc_defeated event to the API server's event bus.
// It is called from the game combat loop when a player defeats an NPC.
// Network errors are logged but do not crash the game.
func FireDefeatEvent(characterID, xpValue int) {
	if xpValue <= 0 {
		return // No XP to award — skip
	}

	token, err := fetchAdminEventsToken()
	if err != nil {
		log.Printf("[events] failed to get admin token: %v", err)
		return
	}

	payload := map[string]interface{}{
		"type": "npc_defeated",
		"payload": map[string]int{
			"character_id": characterID,
			"xp_value":     xpValue,
		},
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", RESTAPIBase+"/api/events", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("[events] failed to build request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[events] failed to fire npc_defeated event: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		log.Printf("[events] npc_defeated event rejected: status %d", resp.StatusCode)
	}
}
