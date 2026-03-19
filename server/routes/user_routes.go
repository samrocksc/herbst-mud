package routes

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/user"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// JWT secret key - in production use environment variable
var jwtSecret = []byte("your-secret-key-change-in-production")

// RegisterUserRoutes registers all user-related routes
func RegisterUserRoutes(router *gin.Engine, client *db.Client) {
	// Create a new user
	router.POST("/users", func(c *gin.Context) {
		var req struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
			IsAdmin  bool   `json:"isAdmin"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Hash the password with bcrypt
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		user, err := client.User.
			Create().
			SetEmail(req.Email).
			SetPassword(string(hashedPassword)).
			SetIsAdmin(req.IsAdmin).
			Save(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return user without password
		c.JSON(http.StatusCreated, gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"is_admin": user.IsAdmin,
		})
	})

	// Authenticate a user (login)
	router.POST("/users/auth", func(c *gin.Context) {
		var req struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Find user by email
		users, err := client.User.Query().Where(user.Email(req.Email)).All(c.Request.Context())
		if err != nil || len(users) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		user := users[0]

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

		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":        user.ID,
			"email":     user.Email,
			"is_admin":  user.IsAdmin,
			"token":      tokenString,
			"expires_in": 86400,
		})
	})

	// Get all users
	router.GET("/users", func(c *gin.Context) {
		users, err := client.User.Query().All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return users without passwords
		result := make([]gin.H, len(users))
		for i, user := range users {
			result[i] = gin.H{
				"id":       user.ID,
				"email":    user.Email,
				"is_admin": user.IsAdmin,
			}
		}

		c.JSON(http.StatusOK, result)
	})

	// Get a single user by ID
	router.GET("/users/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		user, err := client.User.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"is_admin": user.IsAdmin,
		})
	})

	// Update a user by ID
	router.PUT("/users/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
			IsAdmin  *bool  `json:"isAdmin"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updater := client.User.UpdateOneID(id)

		// Only update fields that are provided
		if req.Email != "" {
			updater.SetEmail(req.Email)
		}
		if req.Password != "" {
			// Hash the new password with bcrypt
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
				return
			}
			updater.SetPassword(string(hashedPassword))
		}
		if req.IsAdmin != nil {
			updater.SetIsAdmin(*req.IsAdmin)
		}

		user, err := updater.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"is_admin": user.IsAdmin,
		})
	})

	// Delete a user by ID
	router.DELETE("/users/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		err = client.User.DeleteOneID(id).Exec(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	})
}