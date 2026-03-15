package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupTestClient(t *testing.T) *db.Client {
	client, err := db.Open("postgres", "host=localhost port=5432 user=herbst password=herbst_password dbname=herbst_mud sslmode=disable")
	if err != nil {
		t.Skipf("Skipping test: cannot connect to database: %v", err)
	}
	return client
}

func TestCharacterCreationFlow(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	router := gin.New()
	RegisterCharacterRoutes(router, client)

	// First, create a test user
	userReq := map[string]interface{}{
		"email":    "test_char_creation@example.com",
		"password": "testpass123",
	}
	userJSON, _ := json.Marshal(userReq)
	
	userReq2, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(userJSON))
	userReq2.Header.Set("Content-Type", "application/json")
	userWriter := httptest.NewRecorder()
	router.ServeHTTP(userWriter, userReq2)

	var createdUser struct {
		ID int `json:"id"`
	}
	json.Unmarshal(userWriter.Body.Bytes(), &createdUser)

	t.Logf("Created user with ID: %d", createdUser.ID)

	// Test: Check if user needs to create a character (should be true)
	needsReq, _ := http.NewRequest("GET", "/users/1/characters/needed", nil)
	needsWriter := httptest.NewRecorder()
	router.ServeHTTP(needsWriter, needsReq)

	t.Logf("Needs character response: %s", needsWriter.Body.String())

	// Test: Create a character for the user
	charReq := map[string]interface{}{
		"name":     "TestHero",
		"password": "heropass123",
	}
	charJSON, _ := json.Marshal(charReq)
	charReq2, _ := http.NewRequest("POST", "/users/1/characters", bytes.NewBuffer(charJSON))
	charReq2.Header.Set("Content-Type", "application/json")
	charWriter := httptest.NewRecorder()
	router.ServeHTTP(charWriter, charReq2)

	if charWriter.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d. Body: %s", charWriter.Code, charWriter.Body.String())
	}

	var createdChar struct {
		ID              int    `json:"id"`
		Name            string `json:"name"`
		Hitpoints       int    `json:"hitpoints"`
		MaxHitpoints    int    `json:"max_hitpoints"`
		StartingRoomId  int    `json:"startingRoomId"`
	}
	json.Unmarshal(charWriter.Body.Bytes(), &createdChar)

	t.Logf("Created character: ID=%d, Name=%s, HP=%d/%d, StartingRoom=%d", 
		createdChar.ID, createdChar.Name, createdChar.Hitpoints, createdChar.MaxHitpoints, createdChar.StartingRoomId)

	// Verify default stats
	if createdChar.Hitpoints != 100 || createdChar.MaxHitpoints != 100 {
		t.Errorf("Expected default HP 100, got %d/%d", createdChar.Hitpoints, createdChar.MaxHitpoints)
	}

	// Test: Get characters for the user
	getReq, _ := http.NewRequest("GET", "/users/1/characters", nil)
	getWriter := httptest.NewRecorder()
	router.ServeHTTP(getWriter, getReq)

	if getWriter.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", getWriter.Code)
	}

	t.Logf("User characters response: %s", getWriter.Body.String())

	// Test: Check if user needs to create a character now (should be false)
	needsReq2, _ := http.NewRequest("GET", "/users/1/characters/needed", nil)
	needsWriter2 := httptest.NewRecorder()
	router.ServeHTTP(needsWriter2, needsReq2)

	t.Logf("Needs character after creation: %s", needsWriter2.Body.String())

	// Test: Duplicate character name should fail
	charReq2 = map[string]interface{}{
		"name":     "TestHero",
		"password": "heropass123",
	}
	charJSON2, _ := json.Marshal(charReq2)
	dupeReq, _ := http.NewRequest("POST", "/users/1/characters", bytes.NewBuffer(charJSON2))
	dupeReq.Header.Set("Content-Type", "application/json")
	dupeWriter := httptest.NewRecorder()
	router.ServeHTTP(dupeWriter, dupeReq)

	if dupeWriter.Code != http.StatusConflict {
		t.Errorf("Expected status 409 for duplicate name, got %d. Body: %s", dupeWriter.Code, dupeWriter.Body.String())
	}

	t.Logf("Duplicate character test passed")
}