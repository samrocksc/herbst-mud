package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// determineReconnectRoom decides where the player should appear on login.
// Offline < 1 hour: return to last known room (currentRoomId from server).
// Offline >= 1 hour: return to bind point (respawnRoomId).
// Fallback: root room.
func (m *model) determineReconnectRoom() int {
	if m.currentCharacterID == 0 {
		return m.getRootRoomID()
	}

	// Fetch character data from server to get currentRoomId, respawnRoomId, lastSeenAt
	resp, err := httpGet(fmt.Sprintf("%s/characters/%d", RESTAPIBase, m.currentCharacterID))
	if err != nil {
		return m.getRootRoomID()
	}
	defer resp.Body.Close()

	var char struct {
		CurrentRoomId int        `json:"currentRoomId"`
		RespawnRoomId int        `json:"respawnRoomId"`
		LastSeenAt    *time.Time `json:"lastSeenAt"`
	}
	if json.NewDecoder(resp.Body).Decode(&char) != nil {
		return m.getRootRoomID()
	}

	// Store respawn room for later use
	m.respawnRoom = char.RespawnRoomId

	// If lastSeenAt is nil (first login or no data), use current room or root
	if char.LastSeenAt == nil {
		if char.CurrentRoomId > 0 {
			return char.CurrentRoomId
		}
		return m.getRootRoomID()
	}

	// Calculate time since last seen
	timeSinceLastSeen := time.Since(*char.LastSeenAt)

	if timeSinceLastSeen >= time.Hour {
		// Offline >= 1 hour: return to bind point (respawn room)
		if char.RespawnRoomId > 0 {
			return char.RespawnRoomId
		}
		return m.getRootRoomID()
	}

	// Offline < 1 hour: return to last known room
	if char.CurrentRoomId > 0 {
		return char.CurrentRoomId
	}
	return m.getRootRoomID()
}

// getRootRoomID fetches the root room ID from the server.
func (m *model) getRootRoomID() int {
	if m.client == nil {
		return StartingRoomID
	}

	rooms, err := m.client.Room.Query().All(context.Background())
	if err != nil || len(rooms) == 0 {
		return StartingRoomID
	}

	for _, r := range rooms {
		if r.IsRootRoom {
			return r.ID
		}
	}

	// Fallback: use starting room flag
	for _, r := range rooms {
		if r.IsStartingRoom {
			return r.ID
		}
	}

	return rooms[0].ID
}

// loadRoom loads a room by ID and sets the model's room state.
func (m *model) loadRoom(roomID int) {
	if m.client == nil {
		m.currentRoom = roomID
		return
	}

	room, err := m.client.Room.Get(context.Background(), roomID)
	if err != nil {
		m.AppendMessage(fmt.Sprintf("Failed to load room %d, using fallback.", roomID), "error")
		m.currentRoom = StartingRoomID
		return
	}

	m.currentRoom = room.ID
	m.roomName = room.Name
	m.roomDesc = room.Description
	m.exits = room.Exits
	for dir := range m.exits {
		m.knownExits[dir] = true
	}
}

// updateLastSeenAt updates the character's lastSeenAt timestamp on the server.
func (m *model) updateLastSeenAt() {
	if m.currentCharacterID == 0 {
		return
	}

	url := fmt.Sprintf("%s/characters/%d", RESTAPIBase, m.currentCharacterID)
	now := time.Now().UTC().Format(time.RFC3339)
	payload := fmt.Sprintf(`{"lastSeenAt": "%s"}`, now)

	resp, err := httpPut(url, payload)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}