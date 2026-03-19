package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/user"
	"herbst-server/middleware"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUserRoutes registers all user-related routes
// Public routes (no auth required):
//   - POST /users - Create new user
//   - POST /users/auth - Login (returns JWT token)
//
// Protected routes (auth required):
//   - GET /users - Get all users
//   - GET /users/:id - Get user by ID
//   - PUT /users/:id - Update user
//   - DELETE /users/:id - Delete user
func RegisterUserRoutes(router *gin.Engine, client *db.Client) {
	// === PUBLIC ROUTES (no authentication required) ===

	// Create a new user (public - for registration)
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

		// Generate JWT token for the newly registered user
		token, err := middleware.GenerateToken(user.ID, user.Email, user.IsAdmin, "user")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		// Return user with token
		c.JSON(http.StatusCreated, gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"is_admin": user.IsAdmin,
			"token":    token,
		})
	})

	// Authenticate a user (login) - returns JWT token
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

		// Generate JWT token
		token, err := middleware.GenerateToken(user.ID, user.Email, user.IsAdmin, "user")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		// Login successful with token
		c.JSON(http.StatusOK, gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"is_admin": user.IsAdmin,
			"token":    token,
		})
	})

	// === PROTECTED ROUTES (authentication required) ===

	// Protected user routes group
	protected := router.Group("/users")
	protected.Use(middleware.AuthMiddleware())
	{
		// Get all users (protected)
		protected.GET("", func(c *gin.Context) {
			users, err := client.User.Query().All(c.Request.Context())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Return users without passwords
			result := make([]gin.H, len(users))
			for i, u := range users {
				result[i] = gin.H{
					"id":       u.ID,
					"email":    u.Email,
					"is_admin": u.IsAdmin,
				}
			}

			c.JSON(http.StatusOK, result)
		})

		// Get a single user by ID (protected)
		protected.GET("/:id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
				return
			}

			u, err := client.User.Get(c.Request.Context(), id)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"id":       u.ID,
				"email":    u.Email,
				"is_admin": u.IsAdmin,
			})
		})

		// Update a user by ID (protected)
		protected.PUT("/:id", func(c *gin.Context) {
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

			u, err := updater.Save(c.Request.Context())
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"id":       u.ID,
				"email":    u.Email,
				"is_admin": u.IsAdmin,
			})
		})

		// Delete a user by ID (protected)
		protected.DELETE("/:id", func(c *gin.Context) {
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
}

// RegisterAdminRoutes registers admin-only routes
// These routes require both authentication and admin privileges
func RegisterAdminRoutes(router *gin.Engine, client *db.Client) {
	admin := router.Group("/admin")
	admin.Use(middleware.AuthMiddleware())
	admin.Use(middleware.AdminMiddleware())
	{
		// Admin-only endpoints can be added here
		// Example:
		// admin.GET("/stats", func(c *gin.Context) { ... })
	}
}