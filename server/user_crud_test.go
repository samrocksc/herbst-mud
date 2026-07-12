package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/repository"
	"herbst-server/routes"
)

func uniqueUserEmail(prefix string) string {
	return prefix + "_" + strconv.FormatInt(time.Now().UnixNano(), 10) + "@example.com"
}

func TestUserCRUD(t *testing.T) {
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
	routes.RegisterUserRoutes(router, repos)

	// Test creating a user
	t.Run("CreateUser", func(t *testing.T) {
		userData := map[string]interface{}{
			"email":    uniqueUserEmail("testuser"),
			"password": "testpassword123",
			"isAdmin":  false,
		}

		jsonData, _ := json.Marshal(userData)
		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusCreated, resp.Code, resp.Body.String())
		}
	})

	// Test getting all users
	t.Run("GetAllUsers", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, resp.Code)
		}
	})

	// Test getting a single user by ID
	t.Run("GetUserByID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users/1", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		// Could be 200 (found) or 404 (not found) depending on data
		if resp.Code != http.StatusOK && resp.Code != http.StatusNotFound {
			t.Errorf("Expected status 200 or 404, got %d", resp.Code)
		}
	})

	// Test updating a user
	t.Run("UpdateUser", func(t *testing.T) {
		// First create a user to update
		userData := map[string]interface{}{
			"email":    uniqueUserEmail("updatetest"),
			"password": "originalpassword",
			"isAdmin":  false,
		}
		jsonData, _ := json.Marshal(userData)
		createReq, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
		createReq.Header.Set("Content-Type", "application/json")
		createResp := httptest.NewRecorder()
		router.ServeHTTP(createResp, createReq)

		if createResp.Code == http.StatusCreated {
			var createdUser map[string]interface{}
			json.Unmarshal(createResp.Body.Bytes(), &createdUser)
			userID := int(createdUser["id"].(float64))

			// Now update it
			updateData := map[string]interface{}{
				"email":   uniqueUserEmail("updated"),
				"isAdmin": true,
			}
			updateJSON, _ := json.Marshal(updateData)
			updateReq, _ := http.NewRequest("PUT", "/users/"+strconv.Itoa(userID), bytes.NewBuffer(updateJSON))
			updateReq.Header.Set("Content-Type", "application/json")
			updateResp := httptest.NewRecorder()
			router.ServeHTTP(updateResp, updateReq)

			if updateResp.Code != http.StatusOK && updateResp.Code != http.StatusNotFound {
				t.Errorf("Expected status 200 or 404, got %d", updateResp.Code)
			}
		}
	})

	// Test deleting a user
	t.Run("DeleteUser", func(t *testing.T) {
		// First create a user to delete
		userData := map[string]interface{}{
			"email":    uniqueUserEmail("deletetest"),
			"password": "deletepassword",
			"isAdmin":  false,
		}
		jsonData, _ := json.Marshal(userData)
		createReq, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
		createReq.Header.Set("Content-Type", "application/json")
		createResp := httptest.NewRecorder()
		router.ServeHTTP(createResp, createReq)

		if createResp.Code == http.StatusCreated {
			var createdUser map[string]interface{}
			json.Unmarshal(createResp.Body.Bytes(), &createdUser)
			userID := int(createdUser["id"].(float64))

			// Now delete it
			deleteReq, _ := http.NewRequest("DELETE", "/users/"+strconv.Itoa(userID), nil)
			deleteResp := httptest.NewRecorder()
			router.ServeHTTP(deleteResp, deleteReq)

			if deleteResp.Code != http.StatusNoContent && deleteResp.Code != http.StatusNotFound {
				t.Errorf("Expected status 204 or 404, got %d", deleteResp.Code)
			}
		}
	})
}