package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// Helper to generate a valid JWT token for testing
func generateTestToken(userID uint, email string, isAdmin bool) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userID,
		"email":    email,
		"is_admin": isAdmin,
		"exp":      time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString(jwtSecret)
	return tokenString
}

// Helper to generate an expired JWT token for testing
func generateExpiredToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  1,
		"email":    "test@example.com",
		"is_admin": false,
		"exp":      time.Now().Add(-time.Hour).Unix(), // Expired 1 hour ago
	})
	tokenString, _ := token.SignedString(jwtSecret)
	return tokenString
}

func TestAuthMiddleware_NoHeader(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	AuthMiddleware()(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Authorization header required")
}

func TestAuthMiddleware_InvalidFormat(t *testing.T) {
	tests := []struct {
		name   string
		header string
	}{
		{"No Bearer prefix", "some-token"},
		{"Empty after Bearer", "Bearer "},
		{"Wrong prefix", "Basic some-token"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/", nil)
			c.Request.Header.Set("Authorization", tt.header)

			AuthMiddleware()(c)

			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})
	}
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer "+generateTestToken(1, "test@example.com", true))

	AuthMiddleware()(c)

	// Should pass through
	assert.Equal(t, 0, w.Code)
	
	// Check context values
	userID, exists := c.Get("user_id")
	assert.True(t, exists)
	assert.Equal(t, uint(1), userID)
	
	email, exists := c.Get("email")
	assert.True(t, exists)
	assert.Equal(t, "test@example.com", email)
	
	isAdmin, exists := c.Get("is_admin")
	assert.True(t, exists)
	assert.Equal(t, true, isAdmin)
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer "+generateExpiredToken())

	AuthMiddleware()(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer invalid.token.here")

	AuthMiddleware()(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAdminMiddleware_NoUser(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	AdminMiddleware()(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Admin access required")
}

func TestAdminMiddleware_NonAdminUser(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Set("is_admin", false)

	AdminMiddleware()(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestAdminMiddleware_AdminUser(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Set("is_admin", true)

	AdminMiddleware()(c)

	// Should pass through without writing response
	assert.Equal(t, 0, w.Code)
}

func TestOptionalAuthMiddleware_NoHeader(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	OptionalAuthMiddleware()(c)

	// Should pass through without setting user info
	assert.Equal(t, 0, w.Code)
}

func TestOptionalAuthMiddleware_InvalidToken(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer invalid")

	OptionalAuthMiddleware()(c)

	// Should pass through (optional auth doesn't fail on invalid token)
	assert.Equal(t, 0, w.Code)
}

func TestOptionalAuthMiddleware_ValidToken(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer "+generateTestToken(1, "test@example.com", true))

	OptionalAuthMiddleware()(c)

	// Should pass through and set user info
	userID, exists := c.Get("user_id")
	assert.True(t, exists)
	assert.Equal(t, uint(1), userID)
}

func TestValidationError(t *testing.T) {
	err := &ValidationError{Message: "test error"}
	assert.Equal(t, "test error", err.Error())
}