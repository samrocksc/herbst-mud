package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
)

// RegisterCharacterRoutes registers all character-related routes
func RegisterCharacterRoutes(router *gin.Engine, client *db.Client) {
	// Create a new character
	router.POST("/characters", func(c *gin.Context) {
		var req struct {
			Name        string `json:"name" binding:"required"`
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

		builder := client.Character.
			Create().
			SetName(req.Name).
			SetIsNPC(req.IsNPC).
			SetIsAdmin(req.IsAdmin)

		if req.CurrentRoom > 0 {
			builder.SetCurrentRoomID(req.CurrentRoom)
		}
		if req.StartingRoom > 0 {
			builder.SetStartingRoomID(req.StartingRoom)
		}
		if req.UserID > 0 {
			builder.SetUserID(req.UserID)
		}

		character, err := builder.Save(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, character)
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
			updater.SetCurrentRoomID(*req.CurrentRoom)
		}
		if req.StartingRoom != nil {
			updater.SetStartingRoomID(*req.StartingRoom)
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
}