package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/repository"
	"herbst-server/service"
)

// RegisterMeRoutes registers /api/me/* endpoints (requires AuthMiddleware upstream).
func RegisterMeRoutes(group *gin.RouterGroup, svc *service.Container, repos *repository.Container) {
	// GET /api/me/characters — list characters for the authenticated user
	group.GET("/me/characters", func(c *gin.Context) {
		userIDVal, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
			return
		}
		userID := int(userIDVal.(uint))

		characters, err := repos.Character.ListByUser(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		result := make([]gin.H, len(characters))
		for i, char := range characters {
			result[i] = gin.H{
				"id":            char.ID,
				"name":          char.Name,
				"isNPC":         char.IsNPC,
				"is_admin":      char.IsAdmin,
				"currentRoomId": char.CurrentRoomId,
				"currentWorld":  char.CurrentWorld,
				"hitpoints":     char.Hitpoints,
				"max_hitpoints": char.MaxHitpoints,
				"stamina":       char.Stamina,
				"max_stamina":   char.MaxStamina,
				"mana":          char.Mana,
				"max_mana":      char.MaxMana,
				"race":          char.Race,
				"gender":        char.Gender,
				"level":         char.Level,
				"class":         char.Class,
			}
		}
		c.JSON(http.StatusOK, result)
	})

	// POST /api/me/characters — create a character for the authenticated user
	group.POST("/me/characters", func(c *gin.Context) {
		userIDVal, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
			return
		}
		userID := int(userIDVal.(uint))

		var req struct {
			Name   string `json:"name" binding:"required"`
			Race   string `json:"race"`
			Gender string `json:"gender"`
			World  string `json:"world" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		char, err := svc.Character.CreateCharacter(c.Request.Context(), service.CreateCharacterInput{
			UserID: userID,
			Name:   req.Name,
			Race:   req.Race,
			Gender: req.Gender,
			WorldID: req.World,
		})
		if err != nil {
			switch {
			case err == service.ErrCharacterNameTaken:
				c.JSON(http.StatusConflict, gin.H{"error": "Name already taken"})
			case err == service.ErrTooManyCharacters:
				c.JSON(http.StatusBadRequest, gin.H{"error": "Maximum 3 characters per user"})
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"id":            char.ID,
			"name":          char.Name,
			"currentWorld":  char.CurrentWorld,
			"hitpoints":     char.Hitpoints,
			"max_hitpoints": char.MaxHitpoints,
			"stamina":       char.Stamina,
			"max_stamina":   char.MaxStamina,
			"mana":          char.Mana,
			"max_mana":      char.MaxMana,
			"race":          char.Race,
			"gender":        char.Gender,
			"level":         char.Level,
			"class":         char.Class,
		})
	})
}
