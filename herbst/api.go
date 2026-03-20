package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// ============================================================
// API LOADING - Load data from REST API
// ============================================================

// loadRoomItems fetches items for the current room from the API
func (m *model) loadRoomItems() {
	if m.currentRoom == 0 {
		return
	}

	resp, err := http.Get(fmt.Sprintf("%s/rooms/%d/equipment", RESTAPIBase, m.currentRoom))
	if err != nil {
		log.Printf("Error fetching room items: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	var items []RoomItem
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		log.Printf("Error decoding room items: %v", err)
		return
	}

	m.roomItems = items
}

// loadRoomCharacters fetches characters (NPCs and players) in the current room from the API
func (m *model) loadRoomCharacters() {
	if m.currentRoom == 0 {
		return
	}

	resp, err := http.Get(fmt.Sprintf("%s/rooms/%d/characters", RESTAPIBase, m.currentRoom))
	if err != nil {
		log.Printf("Error fetching room characters: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	var characters []roomCharacter
	if err := json.NewDecoder(resp.Body).Decode(&characters); err != nil {
		log.Printf("Error decoding room characters: %v", err)
		return
	}

	m.roomCharacters = characters
}
