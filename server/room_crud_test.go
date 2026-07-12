package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
)

func TestRoomCRUD(t *testing.T) {
	gin.SetMode(gin.TestMode)

	client, err := db.Open("postgres", "host=localhost port=5432 user=herbst password=herbst_password dbname=herbst_mud sslmode=disable")
	if err != nil {
		t.Fatalf("failed connecting to postgres: %v", err)
	}
	defer client.Close()

	router := setupTestRouter(client)

	// Use a unique room name to avoid collisions across test runs.
	uniqueSuffix := time.Now().Format("20060102150405")

	t.Run("CreateRoom", func(t *testing.T) {
		roomData := map[string]interface{}{
			"name":           "Test Room " + uniqueSuffix,
			"description":    "A test room for unit testing",
			"world_id":       "2",
			"isStartingRoom": false,
			"exits": map[string]int{
				"north": 1,
				"south": 2,
			},
		}

		jsonData, _ := json.Marshal(roomData)
		req, _ := http.NewRequest("POST", "/api/rooms?world_id=2", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+adminToken())

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusCreated, resp.Code, resp.Body.String())
		}
	})

	t.Run("GetAllRooms", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/rooms?world_id=2", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken())
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, resp.Code)
		}
	})
}
