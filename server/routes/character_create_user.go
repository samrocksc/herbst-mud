package routes

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/repository"
	"herbst-server/service"
)

// createCharacterForUser handles POST /user-characters/:id (create character for existing user).
func createCharacterForUser(svc *service.Container, repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := parseIntParam(c, "id")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
		if _, err := repos.User.Get(c.Request.Context(), userID); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		var req struct {
			Name        string   `json:"name" binding:"required"`
			Race        string   `json:"race"`
			Gender      string   `json:"gender"`
			Description string   `json:"description"`
			World       string   `json:"world"`
			Factions    []string `json:"factions"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
			char, err := svc.Character.CreateCharacter(c.Request.Context(), service.CreateCharacterInput{
			UserID:   userID,
			Name:     req.Name,
			Race:     req.Race,
			Gender:   req.Gender,
			Description: req.Description,
			WorldID:  req.World,
			Factions: req.Factions,
		})
		if err != nil {
			switch {
			case errors.Is(err, service.ErrCharacterNameTaken):
				c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			case errors.Is(err, service.ErrTooManyCharacters):
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			case errors.Is(err, service.ErrInvalidRace):
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid race"})
			case errors.Is(err, service.ErrInvalidGender):
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid gender"})
			case errors.Is(err, service.ErrWorldNotReady):
				c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"id":             char.ID,
			"name":           char.Name,
			"isNPC":          char.IsNPC,
			"is_admin":       char.IsAdmin,
			"currentRoomId":  char.CurrentRoomId,
			"startingRoomId":  char.StartingRoomId,
			"hitpoints":      char.Hitpoints,
			"max_hitpoints":  char.MaxHitpoints,
			"stamina":        char.Stamina,
			"max_stamina":    char.MaxStamina,
			"mana":           char.Mana,
			"max_mana":       char.MaxMana,
			"race":           char.Race,
			"class":          char.Class,
			"specialty":      char.Specialty,
			"currentWorld":   char.CurrentWorld,
			"strength":       char.Strength,
			"dexterity":      char.Dexterity,
			"constitution":   char.Constitution,
			"intelligence":   char.Intelligence,
			"wisdom":         char.Wisdom,
			"level":          char.Level,
			"xp":            char.Xp,
		})
	}
}