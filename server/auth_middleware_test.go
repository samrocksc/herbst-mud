package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/middleware"
	"herbst-server/routes"
)

// TestAuthMiddlewareIntegration tests the auth middleware integration with routes
func TestAuthMiddlewareIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Set a test JWT secret
	os.Setenv("JWT_SECRET", "test-secret-key-for-integration-tests")

	// Initialize database
	client, err := db.Open("postgres", "host=localhost port=5432 user=herbst password=herbst_password dbname=herbst_mud sslmode=disable")
	if err != nil {
		t.Skipf("Skipping test - no database available: %v", err)
	}
	defer client.Close()

	// Create router with routes
	router := gin.New()
	routes.RegisterRoomRoutes(router, client)
	routes.RegisterUserRoutes(router, client)
	routes.RegisterCharacterRoutes(router, client)
	routes.RegisterEquipmentRoutes(router, client)

	t.Run("PublicRoutesWorkWithoutAuth", func(t *testing.T) {
		// GET /rooms should work without auth (public)
		req, _ := http.NewRequest("GET", "/rooms", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		// Should succeed (200) or fail for DB reasons (500), not 401
		if resp.Code == http.StatusUnauthorized {
			t.Errorf("Public route /rooms returned 401 - should be accessible without auth")
		}
	})

	t.Run("ProtectedRoutesRequireAuth", func(t *testing.T) {
		// POST /rooms should require auth (protected)
		roomData := map[string]interface{}{
			"name":        "Test Room",
			"description": "A test room",
		}
		jsonData, _ := json.Marshal(roomData)
		req, _ := http.NewRequest("POST", "/rooms", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401 for protected route /rooms POST, got %d", resp.Code)
		}
	})

	t.Run("ProtectedRoutesWorkWithValidToken", func(t *testing.T) {
		// Generate a valid token
		token, err := middleware.GenerateTokenWithSecret(1, "test@example.com", false, "user", "test-secret-key-for-integration-tests")
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}

		// POST /rooms with valid token
		roomData := map[string]interface{}{
			"name":        "Test Room",
			"description": "A test room",
		}
		jsonData, _ := json.Marshal(roomData)
		req, _ := http.NewRequest("POST", "/rooms", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		// Should not return 401 (might fail for DB reasons, but not auth)
		if resp.Code == http.StatusUnauthorized {
			t.Errorf("Protected route returned 401 even with valid token")
		}
	})

	t.Run("UserRegistrationReturnsToken", func(t *testing.T) {
		// Create a unique test email
		testEmail := "auth_test_" + t.Name() + "@example.com"
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
			t.Logf("Response: %s", resp.Body.String())
		}

		var result map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &result)

		// Should return a token
		if result["token"] == nil {
			t.Error("User registration should return a JWT token")
		}
	})

	t.Run("UserLoginReturnsToken", func(t *testing.T) {
		// First create a user
		testEmail := "auth_login_test_" + t.Name() + "@example.com"
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

		// Now try to login
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

		// Should return a token
		if result["token"] == nil {
			t.Error("User login should return a JWT token")
		}
	})

	t.Run("ProtectedUserRoutesRequireAuth", func(t *testing.T) {
		// GET /users (list all) should require auth
		req, _ := http.NewRequest("GET", "/users", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401 for protected GET /users, got %d", resp.Code)
		}
	})

	t.Run("ProtectedCharacterRoutesRequireAuth", func(t *testing.T) {
		// GET /characters (list all) should require auth
		req, _ := http.NewRequest("GET", "/characters", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401 for protected GET /characters, got %d", resp.Code)
		}
	})

	t.Run("ProtectedEquipmentRoutesRequireAuth", func(t *testing.T) {
		// POST /equipment should require auth
		equipData := map[string]interface{}{
			"name": "Test Sword",
		}
		jsonData, _ := json.Marshal(equipData)
		req, _ := http.NewRequest("POST", "/equipment", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401 for protected POST /equipment, got %d", resp.Code)
		}
	})
}

// TestAdminMiddleware tests admin-only route protection
func TestAdminMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testSecret := "test-admin-secret"
	os.Setenv("JWT_SECRET", testSecret)

	// Create router with admin routes
	router := gin.New()
	client, err := db.Open("postgres", "host=localhost port=5432 user=herbst password=herbst_password dbname=herbst_mud sslmode=disable")
	if err != nil {
		t.Skipf("Skipping test - no database available: %v", err)
	}
	defer client.Close()

	routes.RegisterAdminRoutes(router, client)

	t.Run("AdminRoutesRequireAdminToken", func(t *testing.T) {
		// Generate a non-admin token
		token, _ := middleware.GenerateTokenWithSecret(1, "user@example.com", false, "user", testSecret)

		req, _ := http.NewRequest("GET", "/admin/stats", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusForbidden {
			t.Errorf("Expected 403 for non-admin user, got %d", resp.Code)
		}
	})

	t.Run("AdminRoutesWorkWithAdminToken", func(t *testing.T) {
		// Generate an admin token
		token, _ := middleware.GenerateTokenWithSecret(1, "admin@example.com", true, "user", testSecret)

		req, _ := http.NewRequest("GET", "/admin/stats", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		// Should not be forbidden (might be 404 if route doesn't exist, but not 403)
		if resp.Code == http.StatusForbidden {
			t.Error("Admin routes should be accessible with admin token")
		}
	})
}