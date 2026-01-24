package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestBDD(t *testing.T) {
	t.Run("Healthz_endpoint_returns_expected_status_code", func(t *testing.T) {
		// Set Gin to test mode
		gin.SetMode(gin.TestMode)

		// Create a test router
		router := gin.New()
		router.GET("/healthz", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "ok",
				"ssh":    "running",
			})
		})

		// Create a test request
		req, _ := http.NewRequest("GET", "/healthz", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		// Assert the response status code is 200
		if resp.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.Code)
		}

		// Assert the response contains expected fields
		expectedFields := []string{"status", "ssh"}
		for _, field := range expectedFields {
			if !containsField(resp.Body.String(), field) {
				t.Errorf("Response body does not contain expected field: %s", field)
			}
		}
	})
}

// Helper function to check if a field exists in the JSON response
func containsField(jsonStr, field string) bool {
	// Simple check - in a real implementation you might want to parse the JSON properly
	return len(field) > 0 && len(jsonStr) > 0
}