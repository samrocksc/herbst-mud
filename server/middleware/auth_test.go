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

// TestAuthMiddleware_ValidToken tests that a valid JWT token allows access
func TestAuthMiddleware_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	secret := "test-secret"
	token, err := jwt.GenerateToken(1, "test@example.com", false, secret, time.Hour)
	assert.NoError(t, err)

	router := gin.New()
	router.Use(AuthMiddleware(secret))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestAuthMiddleware_MissingToken tests that missing token returns 401
func TestAuthMiddleware_MissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(AuthMiddleware("test-secret"))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestAuthMiddleware_InvalidToken tests that invalid token returns 401
func TestAuthMiddleware_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(AuthMiddleware("test-secret"))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestAuthMiddleware_WrongSecret tests that token with wrong secret returns 401
func TestAuthMiddleware_WrongSecret(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Generate token with secret1
	token, err := jwt.GenerateToken(1, "test@example.com", false, "secret1", time.Hour)
	assert.NoError(t, err)

	router := gin.New()
	router.Use(AuthMiddleware("secret2"))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestAuthMiddleware_ExpiredToken tests that expired token returns 401
func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	secret := "test-secret"
	// Generate token that expired 1 hour ago
	token, err := jwt.GenerateToken(1, "test@example.com", false, secret, -time.Hour)
	assert.NoError(t, err)

	router := gin.New()
	router.Use(AuthMiddleware(secret))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestAuthMiddleware_ExtractsUserInfo tests that user info is extracted from token
func TestAuthMiddleware_ExtractsUserInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	secret := "test-secret"
	userID := 42
	email := "testuser@example.com"
	isAdmin := true
	token, err := jwt.GenerateToken(userID, email, isAdmin, secret, time.Hour)
	assert.NoError(t, err)

	router := gin.New()
	router.Use(AuthMiddleware(secret))
	router.GET("/test", func(c *gin.Context) {
		// Verify user info was extracted
		assert.Equal(t, float64(userID), c.Get("userID"))
		assert.Equal(t, email, c.Get("email"))
		assert.Equal(t, isAdmin, c.Get("isAdmin"))
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestAdminMiddleware_AdminAllowed tests that admin user can access admin routes
func TestAdminMiddleware_AdminAllowed(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("isAdmin", true)
		c.Next()
	})
	router.Use(AdminMiddleware())
	router.GET("/admin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	req, _ := http.NewRequest("GET", "/admin", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestAdminMiddleware_NonAdminDenied tests that non-admin user is denied
func TestAdminMiddleware_NonAdminDenied(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("isAdmin", false)
		c.Next()
	})
	router.Use(AdminMiddleware())
	router.GET("/admin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	req, _ := http.NewRequest("GET", "/admin", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// TestAdminMiddleware_NoAuthHeader tests that missing isAdmin claim returns 403
func TestAdminMiddleware_NoAuthHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(AdminMiddleware())
	router.GET("/admin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	req, _ := http.NewRequest("GET", "/admin", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// TestOptionalAuthMiddleware_WithToken tests that optional auth works with valid token
func TestOptionalAuthMiddleware_WithToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	secret := "test-secret"
	token, err := jwt.GenerateToken(1, "test@example.com", false, secret, time.Hour)
	assert.NoError(t, err)

	router := gin.New()
	router.Use(OptionalAuthMiddleware(secret))
	router.GET("/test", func(c *gin.Context) {
		userID, exists := c.Get("userID")
		assert.True(t, exists)
		assert.Equal(t, float64(1), userID)
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestOptionalAuthMiddleware_NoToken tests that optional auth allows no token
func TestOptionalAuthMiddleware_NoToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(OptionalAuthMiddleware("test-secret"))
	router.GET("/test", func(c *gin.Context) {
		_, exists := c.Get("userID")
		assert.False(t, exists) // No userID should be set
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestOptionalAuthMiddleware_InvalidToken tests that optional auth allows invalid token
func TestOptionalAuthMiddleware_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(OptionalAuthMiddleware("test-secret"))
	router.GET("/test", func(c *gin.Context) {
		_, exists := c.Get("userID")
		assert.False(t, exists) // No userID should be set
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}