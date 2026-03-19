package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Test secret
	testSecret := "test-secret-key"
	os.Setenv("JWT_SECRET", testSecret)

	t.Run("Missing Authorization Header", func(t *testing.T) {
		router := gin.New()
		router.Use(AuthMiddleware())
		router.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req, _ := http.NewRequest("GET", "/protected", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("Invalid Authorization Format", func(t *testing.T) {
		router := gin.New()
		router.Use(AuthMiddleware())
		router.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "InvalidFormat")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("Valid Token", func(t *testing.T) {
		router := gin.New()
		router.Use(AuthMiddleware())
		router.GET("/protected", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			email, _ := c.Get("email")
			isAdmin, _ := c.Get("is_admin")
			c.JSON(http.StatusOK, gin.H{
				"user_id":  userID,
				"email":    email,
				"is_admin": isAdmin,
			})
		})

		// Generate a valid token
		token, err := GenerateTokenWithSecret(1, "test@example.com", false, "user", testSecret)
		assert.NoError(t, err)

		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		router := gin.New()
		router.Use(AuthMiddleware())
		router.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("Token Signed With Wrong Secret", func(t *testing.T) {
		router := gin.New()
		router.Use(AuthMiddleware())
		router.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// Generate token with different secret
		token, err := GenerateTokenWithSecret(1, "test@example.com", false, "user", "wrong-secret")
		assert.NoError(t, err)

		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})
}

func TestAdminMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Admin User Can Access", func(t *testing.T) {
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("is_admin", true)
			c.Set("user_id", 1)
			c.Next()
		})
		router.Use(AdminMiddleware())
		router.GET("/admin", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "admin access"})
		})

		req, _ := http.NewRequest("GET", "/admin", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("Non-Admin User Cannot Access", func(t *testing.T) {
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("is_admin", false)
			c.Set("user_id", 1)
			c.Next()
		})
		router.Use(AdminMiddleware())
		router.GET("/admin", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "admin access"})
		})

		req, _ := http.NewRequest("GET", "/admin", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusForbidden, resp.Code)
	})

	t.Run("Missing is_admin Context", func(t *testing.T) {
		router := gin.New()
		router.Use(AdminMiddleware())
		router.GET("/admin", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "admin access"})
		})

		req, _ := http.NewRequest("GET", "/admin", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusForbidden, resp.Code)
	})
}

func TestGenerateToken(t *testing.T) {
	testSecret := "test-secret-key"
	os.Setenv("JWT_SECRET", testSecret)

	t.Run("Generate Valid Token", func(t *testing.T) {
		token, err := GenerateToken(1, "test@example.com", false, "user")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("Generate Admin Token", func(t *testing.T) {
		token, err := GenerateToken(1, "admin@example.com", true, "user")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("Token Contains Correct Claims", func(t *testing.T) {
		token, err := GenerateTokenWithSecret(42, "user@example.com", true, "user", testSecret)
		assert.NoError(t, err)

		// Parse token to verify claims
		parsedToken, err := ValidateToken(token, testSecret)
		assert.NoError(t, err)
		assert.NotNil(t, parsedToken)

		claims := parsedToken.Claims.(*JWTClaims)
		assert.Equal(t, 42, claims.UserID)
		assert.Equal(t, "user@example.com", claims.Email)
		assert.True(t, claims.IsAdmin)
		assert.Equal(t, "user", claims.TokenType)
	})
}

// ValidateToken parses and validates a JWT token
func ValidateToken(tokenString, secret string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
}