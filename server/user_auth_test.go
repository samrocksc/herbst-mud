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

func TestUserAuthBcrypt(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Initialize database
	client, err := db.Open("postgres", "host=localhost port=5432 user=herbst password=herbst_password dbname=herbst_mud sslmode=disable")
	if err != nil {
		t.Skipf("Skipping test - no database available: %v", err)
	}
	defer client.Close()

	// Create router
	router := gin.New()
	routes.RegisterUserRoutes(router, client)

	// Test user authentication with bcrypt
	t.Run("AuthenticateUserWithBcrypt", func(t *testing.T) {
		// Create a user with a known password
		testEmail := "bcrypt_test_" + t.Name() + "@example.com"
		userData := map[string]interface{}{
			"email":    testEmail,
			"password": "mysecretpassword123",
			"isAdmin":  false,
		}

		jsonData, _ := json.Marshal(userData)
		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusCreated, resp.Code, resp.Body.String())
			return
		}

		// Now try to authenticate with the correct password
		authData := map[string]interface{}{
			"email":    testEmail,
			"password": "mysecretpassword123",
		}
		authJSON, _ := json.Marshal(authData)
		authReq, _ := http.NewRequest("POST", "/users/auth", bytes.NewBuffer(authJSON))
		authReq.Header.Set("Content-Type", "application/json")

		authResp := httptest.NewRecorder()
		router.ServeHTTP(authResp, authReq)

		if authResp.Code != http.StatusOK {
			t.Errorf("Expected status %d for correct password, got %d. Body: %s", http.StatusOK, authResp.Code, authResp.Body.String())
		}

		// Verify response contains user info
		var authResult map[string]interface{}
		json.Unmarshal(authResp.Body.Bytes(), &authResult)
		if authResult["email"] != testEmail {
			t.Errorf("Expected email %s, got %v", testEmail, authResult["email"])
		}

		// Try to authenticate with wrong password
		wrongAuthData := map[string]interface{}{
			"email":    testEmail,
			"password": "wrongpassword",
		}
		wrongAuthJSON, _ := json.Marshal(wrongAuthData)
		wrongAuthReq, _ := http.NewRequest("POST", "/users/auth", bytes.NewBuffer(wrongAuthJSON))
		wrongAuthReq.Header.Set("Content-Type", "application/json")

		wrongAuthResp := httptest.NewRecorder()
		router.ServeHTTP(wrongAuthResp, wrongAuthReq)

		if wrongAuthResp.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d for wrong password, got %d. Body: %s", http.StatusUnauthorized, wrongAuthResp.Code, wrongAuthResp.Body.String())
		}
	})

	// Test authentication with non-existent user
	t.Run("AuthenticateNonExistentUser", func(t *testing.T) {
		authData := map[string]interface{}{
			"email":    "nonexistent@example.com",
			"password": "anypassword",
		}
		authJSON, _ := json.Marshal(authData)
		authReq, _ := http.NewRequest("POST", "/users/auth", bytes.NewBuffer(authJSON))
		authReq.Header.Set("Content-Type", "application/json")

		authResp := httptest.NewRecorder()
		router.ServeHTTP(authResp, authReq)

		if authResp.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d for non-existent user, got %d. Body: %s", http.StatusUnauthorized, authResp.Code, authResp.Body.String())
		}
	})
}