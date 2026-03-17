package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	dbcharacter "herbst-server/db/character"
	"herbst-server/routes"
)

func TestCharacterAuthentication(t *testing.T) {
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

	// Test creating a character with password
	t.Run("CreateCharacterWithPassword", func(t *testing.T) {
		charData := map[string]interface{}{
			"name":     "Warrior1",
			"password": "securePass123",
			"isNPC":    false,
		}

		jsonData, _ := json.Marshal(charData)
		req, _ := http.NewRequest("POST", "/characters", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusCreated, resp.Code, resp.Body.String())
		}

		// Verify password is NOT returned in response
		var createdChar map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &createdChar)
		if _, exists := createdChar["password"]; exists {
			t.Error("Password should not be returned in response")
		}
	})

	// Test character authentication success
	t.Run("AuthenticateCharacterSuccess", func(t *testing.T) {
		authData := map[string]interface{}{
			"name":     "Warrior1",
			"password": "securePass123",
		}

		jsonData, _ := json.Marshal(authData)
		req, _ := http.NewRequest("POST", "/characters/authenticate", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, resp.Code, resp.Body.String())
		}

		var authResult map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &authResult)

		if !authResult["authenticated"].(bool) {
			t.Error("Authentication should be successful")
		}

		if authResult["access_token"] == nil || authResult["access_token"] == "" {
			t.Error("Should receive an access token")
		}
	})

	// Test failed authentication with wrong password
	t.Run("AuthenticateCharacterWrongPassword", func(t *testing.T) {
		authData := map[string]interface{}{
			"name":     "Warrior1",
			"password": "wrongPassword",
		}

		jsonData, _ := json.Marshal(authData)
		req, _ := http.NewRequest("POST", "/characters/authenticate", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusUnauthorized, resp.Code, resp.Body.String())
		}

		var authResult map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &authResult)

		if authResult["error"] != "Invalid password" {
			t.Errorf("Expected 'Invalid password' error, got: %v", authResult["error"])
		}
	})

	// Test failed authentication for non-existent character
	t.Run("AuthenticateCharacterNotFound", func(t *testing.T) {
		authData := map[string]interface{}{
			"name":     "UnknownHero",
			"password": "anyPassword",
		}

		jsonData, _ := json.Marshal(authData)
		req, _ := http.NewRequest("POST", "/characters/authenticate", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusUnauthorized, resp.Code, resp.Body.String())
		}

		var authResult map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &authResult)

		if authResult["error"] != "Character not found" {
			t.Errorf("Expected 'Character not found' error, got: %v", authResult["error"])
		}
	})

	// Cleanup - delete test character
	t.Cleanup(func() {
		char, err := client.Character.Query().Where(
			dbcharacter.NameEQ("Warrior1"),
		).Only(context.Background())
		if err == nil {
			client.Character.DeleteOne(char).Exec(context.Background())
		}
	})
}