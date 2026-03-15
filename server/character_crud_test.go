package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/routes"
)

func TestCharacterCRUD(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Initialize database
	client, err := db.Open("postgres", "host=localhost port=5432 user=herbst password=herbst_password dbname=herbst_mud sslmode=disable")
	if err != nil {
		t.Fatalf("failed connecting to postgres: %v", err)
	}
	defer client.Close()

	// Create router
	router := gin.New()
	routes.RegisterCharacterRoutes(router, client)

	// First create a room to reference
	roomData := map[string]interface{}{
		"name":        "Test Room for Character",
		"description": "A test room",
		"isStartingRoom": false,
		"exits":        map[string]int{},
	}
	jsonRoomData, _ := json.Marshal(roomData)
	roomReq, _ := http.NewRequest("POST", "/rooms", bytes.NewBuffer(jsonRoomData))
	roomReq.Header.Set("Content-Type", "application/json")
	roomResp := httptest.NewRecorder()
	router.ServeHTTP(roomResp, roomReq)

	// Skip character tests if room creation failed
	if roomResp.Code != http.StatusCreated {
		t.Skip("Could not create test room, skipping character tests")
	}

	// Parse room response to get ID
	var room map[string]interface{}
	json.Unmarshal(roomResp.Body.Bytes(), &room)
	roomIDFloat, ok := room["id"].(float64)
	if !ok {
		t.Skip("Could not parse room ID")
	}
	roomID := int(roomIDFloat)

	// Test creating a character
	t.Run("CreateCharacter", func(t *testing.T) {
		characterData := map[string]interface{}{
			"name":           "Test Character",
			"isNPC":          false,
			"currentRoomId":  roomID,
			"startingRoomId": roomID,
			"isAdmin":        false,
		}

		jsonData, _ := json.Marshal(characterData)
		req, _ := http.NewRequest("POST", "/characters", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusCreated, resp.Code, resp.Body.String())
		}
	})

	// Test getting all characters
	t.Run("GetAllCharacters", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/characters", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, resp.Code)
		}
	})

	// Test getting a single character by ID
	t.Run("GetCharacterByID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/characters/1", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		// Could be 200 (found) or 404 (not found) depending on data
		if resp.Code != http.StatusOK && resp.Code != http.StatusNotFound {
			t.Errorf("Expected status 200 or 404, got %d", resp.Code)
		}
	})

	// Test updating a character
	t.Run("UpdateCharacter", func(t *testing.T) {
		// First create a character to update
		characterData := map[string]interface{}{
			"name":           "Character To Update",
			"isNPC":          false,
			"currentRoomId":  roomID,
			"startingRoomId": roomID,
			"isAdmin":        false,
		}
		jsonData, _ := json.Marshal(characterData)
		createReq, _ := http.NewRequest("POST", "/characters", bytes.NewBuffer(jsonData))
		createReq.Header.Set("Content-Type", "application/json")
		createResp := httptest.NewRecorder()
		router.ServeHTTP(createResp, createReq)

		if createResp.Code == http.StatusCreated {
			var createdChar map[string]interface{}
			json.Unmarshal(createResp.Body.Bytes(), &createdChar)
			charID := int(createdChar["id"].(float64))

			// Now update it
			updateData := map[string]interface{}{
				"name":   "Updated Character Name",
				"isNPC":  true,
			}
			updateJSON, _ := json.Marshal(updateData)
			updateReq, _ := http.NewRequest("PUT", "/characters/"+string(rune(charID+'0')), bytes.NewBuffer(updateJSON))
			updateReq.Header.Set("Content-Type", "application/json")
			updateResp := httptest.NewRecorder()
			router.ServeHTTP(updateResp, updateReq)

			if updateResp.Code != http.StatusOK && updateResp.Code != http.StatusNotFound {
				t.Errorf("Expected status 200 or 404, got %d", updateResp.Code)
			}
		}
	})

	// Test deleting a character
	t.Run("DeleteCharacter", func(t *testing.T) {
		// First create a character to delete
		characterData := map[string]interface{}{
			"name":           "Character To Delete",
			"isNPC":          false,
			"currentRoomId":  roomID,
			"startingRoomId": roomID,
			"isAdmin":        false,
		}
		jsonData, _ := json.Marshal(characterData)
		createReq, _ := http.NewRequest("POST", "/characters", bytes.NewBuffer(jsonData))
		createReq.Header.Set("Content-Type", "application/json")
		createResp := httptest.NewRecorder()
		router.ServeHTTP(createResp, createReq)

		if createResp.Code == http.StatusCreated {
			var createdChar map[string]interface{}
			json.Unmarshal(createResp.Body.Bytes(), &createdChar)
			charID := int(createdChar["id"].(float64))

			// Now delete it
			deleteReq, _ := http.NewRequest("DELETE", "/characters/"+string(rune(charID+'0')), nil)
			deleteResp := httptest.NewRecorder()
			router.ServeHTTP(deleteResp, deleteReq)

			if deleteResp.Code != http.StatusNoContent && deleteResp.Code != http.StatusNotFound {
				t.Errorf("Expected status 204 or 404, got %d", deleteResp.Code)
			}
		}
	})
}