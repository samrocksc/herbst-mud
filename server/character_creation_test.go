package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

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

// uniqueTestEmail returns an email that won't collide with leftover test data.
func uniqueTestEmail(prefix string) string {
	return prefix + "_" + strconv.FormatInt(time.Now().UnixNano(), 10) + "@example.com"
}

// uniqueNameSuffix returns a letters-only suffix for character names, which must
// only contain a-z or A-Z.
func uniqueNameSuffix(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, n)
	for i := range b {
		buf := make([]byte, 1)
		_, _ = rand.Read(buf)
		b[i] = letters[int(buf[0])%len(letters)]
	}
	return string(b)
}

func TestCharacterCreationFlow(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	router := setupTestRouter(client)

	// Unique suffix for all names/emails created by this test run.
	uniqueSuffix := uniqueNameSuffix(8)

	// Create a test user and capture the returned id.
	userReq := map[string]interface{}{
		"email":    uniqueTestEmail("test_char_creation"),
		"password": "testpass123",
	}
	userJSON, _ := json.Marshal(userReq)

	userHTTPReq, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(userJSON))
	userHTTPReq.Header.Set("Content-Type", "application/json")
	userWriter := httptest.NewRecorder()
	router.ServeHTTP(userWriter, userHTTPReq)

	var createdUser struct {
		ID int `json:"id"`
	}
	json.Unmarshal(userWriter.Body.Bytes(), &createdUser)
	if createdUser.ID == 0 {
		t.Fatalf("Failed to create test user: status=%d body=%s", userWriter.Code, userWriter.Body.String())
	}

	userID := createdUser.ID
	t.Logf("Created user with ID: %d", userID)

	// Check if the user needs to create a character.
	needsReq, _ := http.NewRequest("GET", "/user-characters/"+strconv.Itoa(userID)+"/needed", nil)
	needsWriter := httptest.NewRecorder()
	router.ServeHTTP(needsWriter, needsReq)

	t.Logf("Needs character response: %s", needsWriter.Body.String())

	// Create a character for the user using the actual character creation fields.
	// Use a unique name so the global name uniqueness check does not collide
	// with data left over from previous test runs.
	charName := "TestHero" + uniqueSuffix
	charReq := map[string]interface{}{
		"name":        charName,
		"password":    "heropass123",
		"race":        "Ooze",
		"gender":      "It",
		"class":       "survivor",
		"world":       "2",
		"description": "A test hero",
	}
	charJSON, _ := json.Marshal(charReq)
	charHTTPReq, _ := http.NewRequest("POST", "/user-characters/"+strconv.Itoa(userID), bytes.NewBuffer(charJSON))
	charHTTPReq.Header.Set("Content-Type", "application/json")
	charWriter := httptest.NewRecorder()
	router.ServeHTTP(charWriter, charHTTPReq)

	if charWriter.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d. Body: %s", charWriter.Code, charWriter.Body.String())
	}

	var createdChar struct {
		ID             int    `json:"id"`
		Name           string `json:"name"`
		Hitpoints      int    `json:"hitpoints"`
		MaxHitpoints    int    `json:"max_hitpoints"`
		StartingRoomId int    `json:"startingRoomId"`
	}
	json.Unmarshal(charWriter.Body.Bytes(), &createdChar)

	t.Logf("Created character: ID=%d, Name=%s, HP=%d/%d, StartingRoom=%d",
		createdChar.ID, createdChar.Name, createdChar.Hitpoints, createdChar.MaxHitpoints, createdChar.StartingRoomId)

	// Verify default stats.
	if createdChar.Hitpoints != 100 || createdChar.MaxHitpoints != 100 {
		t.Errorf("Expected default HP 100, got %d/%d", createdChar.Hitpoints, createdChar.MaxHitpoints)
	}

	// Get characters for the user.
	getReq, _ := http.NewRequest("GET", "/user-characters/"+strconv.Itoa(userID), nil)
	getWriter := httptest.NewRecorder()
	router.ServeHTTP(getWriter, getReq)

	if getWriter.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", getWriter.Code)
	}

	t.Logf("User characters response: %s", getWriter.Body.String())

	// After creating a character, the user should no longer need one.
	needsReq2, _ := http.NewRequest("GET", "/user-characters/"+strconv.Itoa(userID)+"/needed", nil)
	needsWriter2 := httptest.NewRecorder()
	router.ServeHTTP(needsWriter2, needsReq2)

	t.Logf("Needs character after creation: %s", needsWriter2.Body.String())

	// Duplicate character name should fail.
	dupCharReq := map[string]interface{}{
		"name":        charName,
		"password":    "heropass123",
		"race":        "Ooze",
		"gender":      "It",
		"class":       "survivor",
		"world":       "2",
		"description": "A duplicate hero",
	}
	charJSON2, _ := json.Marshal(dupCharReq)
	dupeReq, _ := http.NewRequest("POST", "/user-characters/"+strconv.Itoa(userID), bytes.NewBuffer(charJSON2))
	dupeReq.Header.Set("Content-Type", "application/json")
	dupeWriter := httptest.NewRecorder()
	router.ServeHTTP(dupeWriter, dupeReq)

	if dupeWriter.Code != http.StatusConflict {
		t.Errorf("Expected status 409 for duplicate name, got %d. Body: %s", dupeWriter.Code, dupeWriter.Body.String())
	}

	t.Logf("Duplicate character test passed")
}
