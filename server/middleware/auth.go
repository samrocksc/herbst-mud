package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"herbst-server/db"
	"herbst-server/db/user"
)

// GetWorldID extracts world_id from a Gin context. It checks the query
// string first, then the stored "world_id" context key set by
// WorldAccessMiddleware. Returns "" if neither is set.
func GetWorldID(c *gin.Context) string {
	if w := c.Query("world_id"); w != "" {
		return w
	}
	if v, ok := c.Get("world_id"); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// Claims represents JWT claims structure
type Claims struct {
	UserID  uint   `json:"user_id"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

// getJWTSecret returns the JWT secret from environment variable
func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return []byte("dev-secret-key-not-for-production-use-only")
	}
	return []byte(secret)
}

// AuthMiddleware creates authentication middleware
// It validates JWT tokens and extracts user information
// The dbClient parameter is used for querying user allowed_worlds
func AuthMiddleware(dbClient *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Check Bearer token format
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := validateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("is_admin", claims.IsAdmin)
		c.Set("db_client", dbClient)

		// Store allowed worlds from user's whitelist
		allowedWorlds, err := getAllowedWorlds(c, claims.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load user permissions"})
			c.Abort()
			return
		}
		c.Set("allowed_worlds", allowedWorlds)

		c.Next()
	}
}

// getAllowedWorlds retrieves the whitelist of world IDs for a user
func getAllowedWorlds(c *gin.Context, userID uint) ([]string, error) {
	// db_client should be in context from AuthMiddleware
	client, ok := c.Get("db_client")
	if !ok {
		// db_client not in context - return empty list (admin with no restrictions)
		return nil, nil
	}

	// Check if client is nil - note that an interface containing a typed nil pointer
	// is not equal to nil, so we need to check after type assertion
	dbClient, ok := client.(*db.Client)
	if !ok || dbClient == nil {
		// No db client available - return empty list (admin with no restrictions)
		return nil, nil
	}

	u, err := dbClient.User.Query().Where(user.ID(int(userID))).Only(c.Request.Context())
	if err != nil {
		return nil, err
	}

	allowedWorldsStr := u.AllowedWorlds
	if allowedWorldsStr == "" {
		// Empty string means admin can access all worlds
		return nil, nil
	}

	// Split comma-separated list and trim whitespace
	worlds := strings.Split(allowedWorldsStr, ",")
	result := make([]string, 0, len(worlds))
	for _, w := range worlds {
		trimmed := strings.TrimSpace(w)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result, nil
}

// AdminMiddleware creates admin-only middleware
// Must be used after AuthMiddleware
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, exists := c.Get("is_admin")
		if !exists || !isAdmin.(bool) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// WorldIDRequiredMiddleware validates that world_id query parameter is present
// Returns 400 error if missing. Only enforced for mutating methods (POST/PUT/PATCH/DELETE).
func WorldIDRequiredMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Query("world_id") == "" {
			// GET requests (list/read) can omit world_id
			if c.Request.Method != "GET" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "world_id query parameter is required"})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

// WorldAccessMiddleware checks if the user has access to a specific world
// Must be used after AuthMiddleware
func WorldAccessMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedWorlds, exists := c.Get("allowed_worlds")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication middleware not run"})
			c.Abort()
			return
		}

		// Get world_id from query parameter, form, or JSON body
		worldID := c.Query("world_id")
		if worldID == "" {
			// Try to get from JSON body for POST/PUT requests.
			// Read raw body so we can restore it — ShouldBindJSON consumes the stream.
			rawData, _ := c.GetRawData()
			if len(rawData) > 0 {
				var body struct {
					WorldID string `json:"world_id"`
				}
				if err := json.Unmarshal(rawData, &body); err == nil {
					worldID = body.WorldID
				}
				// Restore body for downstream handlers
				c.Request.Body = io.NopCloser(bytes.NewReader(rawData))
			}
		}

		// If no world_id specified, allow access (for list operations)
		if worldID == "" {
			c.Next()
			return
		}

		// Check if user has access to this world
		wl := allowedWorlds.([]string)
		if wl == nil {
			// Nil means admin can access all worlds
			c.Next()
			return
		}

		// Check if world_id is in whitelist
		for _, w := range wl {
			if w == worldID {
				c.Next()
				return
			}
		}

		// World not in whitelist
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this world"})
		c.Abort()
	}
}

// OptionalAuthMiddleware creates optional authentication middleware
// It attaches user info if a valid token is provided, but doesn't require it
func OptionalAuthMiddleware(dbClient *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// No token provided - continue without auth
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			// Invalid format - continue without auth
			c.Next()
			return
		}

		tokenString := parts[1]
		claims, err := validateToken(tokenString)
		if err != nil {
			// Invalid token - continue without auth
			c.Next()
			return
		}

		// Set user info if valid
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("is_admin", claims.IsAdmin)
		c.Set("db_client", dbClient)

		// Store allowed worlds from user's whitelist
		allowedWorlds, err := getAllowedWorlds(c, claims.UserID)
		if err == nil {
			c.Set("allowed_worlds", allowedWorlds)
		}

		c.Next()
	}
}

// ValidateToken validates a JWT token and returns (userID, isAdmin, error).
// Exported for use by handlers that need to verify tokens outside of middleware
// (e.g., SSE endpoints where EventSource cannot send Authorization headers).
func ValidateToken(tokenString string) (uint, bool, error) {
	claims, err := validateToken(tokenString)
	if err != nil {
		return 0, false, err
	}
	return claims.UserID, claims.IsAdmin, nil
}

// validateToken validates JWT token and returns claims
func validateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, &ValidationError{Message: "invalid signing method"}
		}
		return getJWTSecret(), nil
	})

	if err != nil {
		return nil, &ValidationError{Message: err.Error()}
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, &ValidationError{Message: "invalid token"}
}

// ValidationError represents a token validation error
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
