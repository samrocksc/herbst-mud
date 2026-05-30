package routes

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"herbst-server/dblog"
	"herbst-server/db"
	"herbst-server/repository"
	"herbst-server/service"
)

// createCharacter handles POST /characters (admin/debug endpoint).
func createCharacter(svc *service.Container, repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name         string `json:"name" binding:"required"`
			IsNPC        bool   `json:"isNPC"`
			CurrentRoom  int    `json:"currentRoomId"`
			StartingRoom int    `json:"startingRoomId"`
			UserID       int    `json:"userId"`
			IsAdmin      bool   `json:"isAdmin"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			slog.Warn("bad request: invalid create character request", slog.String("service", "characters"), slog.String("error", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if len(req.Name) < 1 || len(req.Name) > 23 {
			slog.Warn("bad request: character name length invalid", slog.String("service", "characters"), slog.String("name", req.Name))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Character name must be 1-23 characters"})
			return
		}
		for _, ch := range req.Name {
			if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')) {
				slog.Warn("bad request: character name contains invalid characters", slog.String("service", "characters"), slog.String("name", req.Name))
				c.JSON(http.StatusBadRequest, gin.H{"error": "Character name can only contain letters (a-z, A-Z)"})
				return
			}
		}
		if req.UserID > 0 {
			count, err := repos.Character.CountByUser(c.Request.Context(), req.UserID)
			if err == nil && count >= 3 {
				slog.Warn("bad request: max characters per user reached", slog.String("service", "characters"), slog.Int("user_id", req.UserID))
				c.JSON(http.StatusBadRequest, gin.H{"error": "Maximum of 3 characters per user reached"})
				return
			}
		}
		char, err := repos.Character.Create(c.Request.Context(), repository.CreateCharacterInput{
			Name:      req.Name,
			UserID:    req.UserID,
			RoomID:    req.CurrentRoom,
			IsAdmin:   req.IsAdmin,
			IsNPC:     req.IsNPC,
			HP:        100,
			MaxHP:     100,
			Stamina:   50,
			MaxStamina: 50,
			Mana:      50,
			MaxMana:   50,
			Level:     1,
		})
		if err != nil {
			dblog.Error("failed to create character", err, slog.String("service", "characters"), slog.String("name", req.Name), slog.Int("user_id", req.UserID))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		slog.Info("character created", slog.String("service", "characters"), slog.Int("character_id", char.ID), slog.String("name", char.Name))
		c.JSON(http.StatusCreated, gin.H{
			"id":             char.ID,
			"name":           char.Name,
			"race":           char.Race,
			"gender":         char.Gender,
			"class":          char.Class,
			"isNPC":          char.IsNPC,
			"is_admin":       char.IsAdmin,
			"currentRoomId":  char.CurrentRoomId,
			"startingRoomId": char.StartingRoomId,
		})
	}
}

// listCharacters handles GET /characters, optionally filtered by name or ID.
func listCharacters(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		characters, err := repos.Character.ListAll(c.Request.Context())
		if err != nil {
			dblog.Error("failed to list characters", err, slog.String("service", "characters"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if search := c.Query("search"); search != "" {
			s := strings.ToLower(search)
			filtered := make([]*db.Character, 0, len(characters))
			for _, ch := range characters {
				// Match by name or by ID
				nameMatch := strings.Contains(strings.ToLower(ch.Name), s)
				idMatch := strings.Contains(strings.ToLower(strconv.Itoa(ch.ID)), s)
				if nameMatch || idMatch {
					filtered = append(filtered, ch)
				}
			}
			characters = filtered
		}
		c.JSON(http.StatusOK, characters)
	}
}

// getCharacter handles GET /characters/:id.
func getCharacter(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := getIDParam(c)
		if !ok {
			return
		}
		char, err := repos.Character.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}
		c.JSON(http.StatusOK, char)
	}
}

// deleteCharacter handles DELETE /characters/:id.
func deleteCharacter(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := getIDParam(c)
		if !ok {
			return
		}
		if err := repos.Character.Delete(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}
		slog.Info("character deleted", slog.String("service", "characters"), slog.Int("character_id", id))
		c.JSON(http.StatusNoContent, nil)
	}
}
