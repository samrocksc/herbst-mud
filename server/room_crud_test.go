package main

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/repository"
	"herbst-server/routes"
	"herbst-server/service"
)

func TestRoomCRUD(t *testing.T) {
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
	repos := repository.NewContainer(client)
	services := service.NewContainer(client, repos, slog.Default())
	routes.RegisterRoomRoutes(router, client, services)

	// Test creating a room
	t.Run("CreateRoom", func(t *testing.T) {
		roomData := map[string]interface{}{
			"name":        "Test Room",
			"description": "A test room for unit testing",
			"isStartingRoom": false,
			"exits": map[string]int{
				"north": 1,
				"south": 2,
			},
		}

		jsonData, _ := json.Marshal(roomData)
		req, _ := http.NewRequest("POST", "/rooms", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, resp.Code)
		}
	})

	// Test getting all rooms
	t.Run("GetAllRooms", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/rooms", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, resp.Code)
		}
	})
}