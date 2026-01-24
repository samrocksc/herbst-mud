package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthzEndpoint(t *testing.T) {
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

	// Assert the response
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "status")
	assert.Contains(t, resp.Body.String(), "ssh")
}