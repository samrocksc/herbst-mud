package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"herbst-server/db"
	"herbst-server/middleware"
)

// generateTokenWithSecret creates a JWT token signed with the given secret.
func generateTokenWithSecret(userID uint, email string, isAdmin bool, secret string) (string, error) {
	claims := &middleware.Claims{
		UserID:  userID,
		Email:   email,
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// uniqueEmail returns a test email that won't collide across runs.
func uniqueEmail(prefix string) string {
	return prefix + "_" + strconv.FormatInt(time.Now().UnixNano(), 10) + "@example.com"
}

// TestAuthMiddlewareIntegration verifies that protected routes reject unauthenticated
// requests and accept authenticated ones using the real route setup.
func TestAuthMiddlewareIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("JWT_SECRET", "test-secret-key-for-integration-tests")

	client, err := db.Open("postgres", "host=localhost port=5432 user=herbst password=herbst_password dbname=herbst_mud sslmode=disable")
	if err != nil {
		t.Skipf("Skipping test - no database available: %v", err)
	}
	defer client.Close()

	router := setupTestRouter(client)

	t.Run("PublicRoutesWorkWithoutAuth", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/races?world_id=1", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code == http.StatusUnauthorized {
			t.Errorf("Public route /races returned 401 - should be accessible without auth")
		}
	})

	t.Run("ProtectedRoutesRequireAuth", func(t *testing.T) {
		roomData := map[string]interface{}{
			"name":        "Test Room",
			"description": "A test room",
			"world_id":    "1",
		}
		jsonData, _ := json.Marshal(roomData)
		req, _ := http.NewRequest("POST", "/api/rooms?world_id=1", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401 for protected POST /api/rooms, got %d", resp.Code)
		}
	})

	t.Run("ProtectedRoutesWorkWithValidToken", func(t *testing.T) {
		token, err := generateTokenWithSecret(1, "test@example.com", true, "test-secret-key-for-integration-tests")
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}

		// /api/me/characters only requires a valid JWT.
		req, _ := http.NewRequest("GET", "/api/me/characters", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code == http.StatusUnauthorized {
			t.Errorf("Protected route returned 401 even with valid token")
		}
	})

	t.Run("UserRegistrationReturnsUserData", func(t *testing.T) {
		testEmail := uniqueEmail("auth_test")
		userData := map[string]interface{}{
			"email":    testEmail,
			"password": "testpassword123",
			"isAdmin":  false,
		}
		jsonData, _ := json.Marshal(userData)
		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusCreated {
			t.Errorf("Expected status 201 for user registration, got %d. Body: %s", resp.Code, resp.Body.String())
			return
		}

		var result map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &result)
		if result["id"] == nil {
			t.Error("User registration should return user data with id")
		}
		if result["email"] == nil {
			t.Error("User registration should return user data with email")
		}
	})

	t.Run("UserLoginReturnsToken", func(t *testing.T) {
		testEmail := uniqueEmail("auth_login_test")
		userData := map[string]interface{}{
			"email":    testEmail,
			"password": "testpassword123",
			"isAdmin":  false,
		}
		jsonData, _ := json.Marshal(userData)
		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		if resp.Code != http.StatusCreated {
			t.Fatalf("Failed to create user: %s", resp.Body.String())
		}

		loginData := map[string]interface{}{
			"email":    testEmail,
			"password": "testpassword123",
		}
		loginJSON, _ := json.Marshal(loginData)
		loginReq, _ := http.NewRequest("POST", "/users/auth", bytes.NewBuffer(loginJSON))
		loginReq.Header.Set("Content-Type", "application/json")
		loginResp := httptest.NewRecorder()
		router.ServeHTTP(loginResp, loginReq)

		if loginResp.Code != http.StatusOK {
			t.Errorf("Login failed with status %d: %s", loginResp.Code, loginResp.Body.String())
			return
		}

		var result map[string]interface{}
		json.Unmarshal(loginResp.Body.Bytes(), &result)
		if result["token"] == nil {
			t.Error("User login should return a JWT token")
		}
	})

	t.Run("ProtectedMeRoutesRequireAuth", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/me/characters", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401 for protected GET /api/me/characters, got %d", resp.Code)
		}
	})

	t.Run("ProtectedAdminRoutesRequireAuth", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/item-instances?world_id=1", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401 for protected GET /api/item-instances, got %d", resp.Code)
		}
	})
}

// TestAdminMiddleware verifies that admin-only routes reject non-admin tokens.
func TestAdminMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testSecret := "test-admin-secret"
	os.Setenv("JWT_SECRET", testSecret)

	client, err := db.Open("postgres", "host=localhost port=5432 user=herbst password=herbst_password dbname=herbst_mud sslmode=disable")
	if err != nil {
		t.Skipf("Skipping test - no database available: %v", err)
	}
	defer client.Close()

	router := setupTestRouter(client)

	t.Run("AdminRoutesRequireAdminToken", func(t *testing.T) {
		token, _ := generateTokenWithSecret(1, "user@example.com", false, testSecret)

		req, _ := http.NewRequest("GET", "/api/rooms?world_id=1", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusForbidden {
			t.Errorf("Expected 403 for non-admin user accessing admin route, got %d", resp.Code)
		}
	})

	t.Run("AdminRoutesWorkWithAdminToken", func(t *testing.T) {
		token, _ := generateTokenWithSecret(1, "admin@example.com", true, testSecret)

		req, _ := http.NewRequest("GET", "/api/rooms?world_id=1", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code == http.StatusForbidden {
			t.Error("Admin routes should be accessible with admin token")
		}
	})
}
