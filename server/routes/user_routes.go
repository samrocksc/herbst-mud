package routes

import (
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"herbst-server/dblog"
	"herbst-server/repository"
)

// getJWTSecret returns the JWT secret from environment variable
// Must match the function in middleware/auth.go
func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return []byte("dev-secret-key-not-for-production-use-only")
	}
	return []byte(secret)
}

// RegisterUserRoutes registers all user-related routes
func RegisterUserRoutes(router *gin.Engine, repos *repository.Container) {
	// Create a new user
	router.POST("/users", func(c *gin.Context) {
		var req struct {
			Email         string `json:"email" binding:"required"`
			Password      string `json:"password" binding:"required"`
			IsAdmin       bool   `json:"isAdmin"`
			AllowedWorlds string `json:"allowed_worlds"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			slog.Warn("invalid create user request", slog.String("error", err.Error()), slog.String("service", "users"))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Hash the password with bcrypt
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			dblog.Error("failed to hash password", err, slog.String("service", "users"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		user, err := repos.User.Create(c.Request.Context(), repository.CreateUserInput{
			Email:         req.Email,
			Password:      string(hashedPassword),
			IsAdmin:       req.IsAdmin,
			AllowedWorlds: req.AllowedWorlds,
		})

		if err != nil {
			dblog.Error("failed to create user", err, slog.String("service", "users"), slog.String("email", req.Email))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return user without password
		slog.Info("user created", slog.Int("user_id", user.ID), slog.String("service", "users"))
		c.JSON(http.StatusCreated, gin.H{
			"id":             user.ID,
			"email":          user.Email,
			"is_admin":       user.IsAdmin,
			"allowed_worlds": user.AllowedWorlds,
			"created_at":     user.CreatedAt.Format(time.RFC3339),
		})
	})

	// Authenticate a user (login)
	router.POST("/users/auth", func(c *gin.Context) {
		var req struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			slog.Warn("invalid login request", slog.String("error", err.Error()), slog.String("service", "users"))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Find user by email
		user, err := repos.User.GetByEmail(c.Request.Context(), req.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		// Compare password with bcrypt hash
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		// Login successful - generate JWT token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id":  user.ID,
			"email":    user.Email,
			"is_admin": user.IsAdmin,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

		tokenString, err := token.SignedString(getJWTSecret())
		if err != nil {
			dblog.Error("failed to sign token", err, slog.String("service", "users"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		slog.Info("user logged in", slog.Int("user_id", user.ID), slog.String("service", "users"))
		c.JSON(http.StatusOK, gin.H{
			"id":             user.ID,
			"email":          user.Email,
			"is_admin":       user.IsAdmin,
			"allowed_worlds": user.AllowedWorlds,
			"token":          tokenString,
			"expires_in":     86400,
		})
	})

	// Get current authenticated user
	router.GET("/users/me", func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}

		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return getJWTSecret(), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			return
		}

		user, err := repos.User.Get(c.Request.Context(), int(userID))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":             user.ID,
			"email":          user.Email,
			"is_admin":       user.IsAdmin,
			"allowed_worlds": user.AllowedWorlds,
		})
	})

	// Get all users
	router.GET("/users", func(c *gin.Context) {
		users, err := repos.User.List(c.Request.Context())
		if err != nil {
			dblog.Error("failed to list users", err, slog.String("service", "users"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return users without passwords
		result := make([]gin.H, len(users))
		for i, user := range users {
			result[i] = gin.H{
				"id":         user.ID,
				"email":      user.Email,
				"is_admin":   user.IsAdmin,
				"created_at": user.CreatedAt.Format(time.RFC3339),
			}
		}

		c.JSON(http.StatusOK, result)
	})

	// Get a single user by ID
	router.GET("/users/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			slog.Warn("invalid user id", slog.String("error", err.Error()), slog.String("service", "users"))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		user, err := repos.User.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":             user.ID,
			"email":          user.Email,
			"is_admin":       user.IsAdmin,
			"allowed_worlds": user.AllowedWorlds,
			"created_at":     user.CreatedAt.Format(time.RFC3339),
		})
	})

	// Update a user by ID
	router.PUT("/users/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			slog.Warn("invalid user id", slog.String("error", err.Error()), slog.String("service", "users"))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var req struct {
			Email         string `json:"email"`
			Password      string `json:"password"`
			IsAdmin       *bool  `json:"isAdmin"`
			AllowedWorlds string `json:"allowed_worlds"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			slog.Warn("invalid update user request", slog.String("error", err.Error()), slog.String("service", "users"))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updates := repository.UserUpdates{
			Email:         &req.Email,
			IsAdmin:       req.IsAdmin,
			AllowedWorlds: &req.AllowedWorlds,
		}

		// Only update password if provided
		if req.Password != "" {
			// Hash the new password with bcrypt
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
			if err != nil {
				dblog.Error("failed to hash password", err, slog.String("service", "users"), slog.Int("user_id", id))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
				return
			}
			updates.Password = new(string)
			*updates.Password = string(hashedPassword)
		}

		// Only set email if provided
		if req.Email == "" {
			updates.Email = nil
		}

		// Only set allowed_worlds if provided
		if req.AllowedWorlds == "" {
			updates.AllowedWorlds = nil
		}

		user, err := repos.User.Update(c.Request.Context(), id, updates)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		slog.Info("user updated", slog.Int("user_id", user.ID), slog.String("service", "users"))
		c.JSON(http.StatusOK, gin.H{
			"id":             user.ID,
			"email":          user.Email,
			"is_admin":       user.IsAdmin,
			"allowed_worlds": user.AllowedWorlds,
			"created_at":     user.CreatedAt.Format(time.RFC3339),
		})
	})

	// Reset password for a user (sets to "password")
	router.POST("/users/:id/reset-password", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			slog.Warn("invalid user id", slog.String("error", err.Error()), slog.String("service", "users"))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		// Hash the default password with bcrypt
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		if err != nil {
			dblog.Error("failed to hash password", err, slog.String("service", "users"), slog.Int("user_id", id))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		user, err := repos.User.Update(c.Request.Context(), id, repository.UserUpdates{
			Password: newString(string(hashedPassword)),
		})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		slog.Info("user password reset", slog.Int("user_id", user.ID), slog.String("service", "users"))
		c.JSON(http.StatusOK, gin.H{
			"id":             user.ID,
			"email":          user.Email,
			"is_admin":       user.IsAdmin,
			"allowed_worlds": user.AllowedWorlds,
			"created_at":     user.CreatedAt.Format(time.RFC3339),
		})
	})

	// Delete a user by ID
	router.DELETE("/users/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			slog.Warn("invalid user id", slog.String("error", err.Error()), slog.String("service", "users"))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		if err := repos.User.Delete(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		slog.Info("user deleted", slog.Int("user_id", id), slog.String("service", "users"))
		c.JSON(http.StatusNoContent, nil)
	})
}

// newString is a helper to create a *string from a string value.
func newString(s string) *string {
	return &s
}
