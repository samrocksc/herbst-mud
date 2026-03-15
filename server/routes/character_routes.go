package routes

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/character"
	"golang.org/x/crypto/bcrypt"
)

// RegisterCharacterRoutes registers all character-related routes
func RegisterCharacterRoutes(router *gin.Engine, client *db.Client) {
	// Create a new character
	router.POST("/characters", func(c *gin.Context) {
		var req struct {
			Name        string `json:"name" binding:"required"`
			Password    string `json:"password"`
			IsNPC       bool   `json:"isNPC"`
			CurrentRoom int    `json:"currentRoomId"`
			StartingRoom int   `json:"startingRoomId"`
			UserID      int    `json:"userId"`
			IsAdmin     bool   `json:"isAdmin"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Hash the password if provided
		var hashedPassword string
		if req.Password != "" {
			hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
				return
			}
			hashedPassword = string(hash)
		}

		builder := client.Character.
			Create().
			SetName(req.Name).
			SetIsNPC(req.IsNPC).
			SetIsAdmin(req.IsAdmin)

		if hashedPassword != "" {
			builder.SetPassword(hashedPassword)
		}

		if req.CurrentRoom > 0 {
			builder.SetCurrentRoomId(req.CurrentRoom)
		}
		if req.StartingRoom > 0 {
			builder.SetStartingRoomId(req.StartingRoom)
		}
		if req.UserID > 0 {
			builder.SetUserID(req.UserID)
		}

		character, err := builder.Save(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"id":           character.ID,
			"name":         character.Name,
			"isNPC":        character.IsNPC,
			"is_admin":     character.IsAdmin,
			"currentRoomId": character.CurrentRoomId,
			"startingRoomId": character.StartingRoomId,
		})
	})

	// Get all characters
	router.GET("/characters", func(c *gin.Context) {
		characters, err := client.Character.Query().All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, characters)
	})

	// Get a single character by ID
	router.GET("/characters/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		character, err := client.Character.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}

		c.JSON(http.StatusOK, character)
	})

	// Update a character by ID
	router.PUT("/characters/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		var req struct {
			Name        string `json:"name"`
			IsNPC       *bool  `json:"isNPC"`
			CurrentRoom *int   `json:"currentRoomId"`
			StartingRoom *int  `json:"startingRoomId"`
			IsAdmin     *bool  `json:"isAdmin"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updater := client.Character.UpdateOneID(id)

		if req.Name != "" {
			updater.SetName(req.Name)
		}
		if req.IsNPC != nil {
			updater.SetIsNPC(*req.IsNPC)
		}
		if req.CurrentRoom != nil {
			updater.SetCurrentRoomId(*req.CurrentRoom)
		}
		if req.StartingRoom != nil {
			updater.SetStartingRoomId(*req.StartingRoom)
		}
		if req.IsAdmin != nil {
			updater.SetIsAdmin(*req.IsAdmin)
		}

		character, err := updater.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}

		c.JSON(http.StatusOK, character)
	})

	// Delete a character by ID
	router.DELETE("/characters/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		err = client.Character.DeleteOneID(id).Exec(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	})

	// Authenticate a character
	router.POST("/characters/authenticate", func(c *gin.Context) {
		var req struct {
			Name     string `json:"name" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Find character by name
		char, err := client.Character.Query().Where(character.NameEQ(req.Name)).Only(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Character not found"})
			return
		}

		// Check if character has a password set
		if char.Password == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Character has no password set"})
			return
		}

		// Verify password
		err = bcrypt.CompareHashAndPassword([]byte(char.Password), []byte(req.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
			return
		}

		// Generate access token
		tokenBytes := make([]byte, 32)
		rand.Read(tokenBytes)
		accessToken := hex.EncodeToString(tokenBytes)

		c.JSON(http.StatusOK, gin.H{
			"authenticated": true,
			"access_token": accessToken,
			"character": gin.H{
				"id":       char.ID,
				"name":     char.Name,
				"is_admin": char.IsAdmin,
			},
		})
	})
}