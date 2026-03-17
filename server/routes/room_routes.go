package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/character"
)

// RegisterRoomRoutes registers all room-related routes
func RegisterRoomRoutes(router *gin.Engine, client *db.Client) {
	// Create a new room
	router.POST("/rooms", func(c *gin.Context) {
		var req struct {
			Name        string            `json:"name" binding:"required"`
			Description string            `json:"description" binding:"required"`
			IsStartingRoom bool           `json:"isStartingRoom"`
			Exits       map[string]int    `json:"exits"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		room, err := client.Room.
			Create().
			SetName(req.Name).
			SetDescription(req.Description).
			SetIsStartingRoom(req.IsStartingRoom).
			SetExits(req.Exits).
			Save(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, room)
	})

	// Get all rooms
	router.GET("/rooms", func(c *gin.Context) {
		rooms, err := client.Room.Query().All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, rooms)
	})

	// Get a single room by ID
	router.GET("/rooms/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
			return
		}

		room, err := client.Room.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
			return
		}

		c.JSON(http.StatusOK, room)
	})

	// Update a room by ID
	router.PUT("/rooms/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
			return
		}

		var req struct {
			Name        string            `json:"name"`
			Description string            `json:"description"`
			IsStartingRoom bool           `json:"isStartingRoom"`
			Exits       map[string]int    `json:"exits"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updater := client.Room.UpdateOneID(id)

		// Only update fields that are provided
		if req.Name != "" {
			updater.SetName(req.Name)
		}
		if req.Description != "" {
			updater.SetDescription(req.Description)
		}
		// For boolean and map fields, we'll always update them if provided
		updater.SetIsStartingRoom(req.IsStartingRoom)
		if req.Exits != nil {
			updater.SetExits(req.Exits)
		}

		room, err := updater.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
			return
		}

		c.JSON(http.StatusOK, room)
	})

	// Delete a room by ID
	router.DELETE("/rooms/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
			return
		}

		err = client.Room.DeleteOneID(id).Exec(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	})

	// Get characters in a room (for displaying NPCs vs players)
	router.GET("/rooms/:id/characters", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
			return
		}

		characters, err := client.Character.Query().
			Where(character.CurrentRoomId(id)).
			All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return characters without sensitive data
		result := make([]gin.H, len(characters))
		for i, char := range characters {
			result[i] = gin.H{
				"id":       char.ID,
				"name":     char.Name,
				"isNPC":    char.IsNPC,
				"level":    char.Level,
				"class":    char.Class,
				"race":     char.Race,
			}
		}

		c.JSON(http.StatusOK, result)
	})
}
