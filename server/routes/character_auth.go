package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/repository"
)


// getUserCharacters handles GET /user-characters/:id.
func getUserCharacters(repos *repository.Container) gin.HandlerFunc {
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
				"startingRoomId": char.StartingRoomId,
				"hitpoints":     char.Hitpoints,
				"max_hitpoints": char.MaxHitpoints,
				"stamina":       char.Stamina,
				"max_stamina":   char.MaxStamina,
				"mana":          char.Mana,
				"max_mana":      char.MaxMana,
			}
		}
		c.JSON(http.StatusOK, result)
	}
}

// needsCharacter handles GET /user-characters/:id/needed.
func needsCharacter(repos *repository.Container) gin.HandlerFunc {
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
		count, err := repos.Character.CountByUser(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"needs_character": count == 0, "character_count": count})
	}
}

func parseIntParam(c *gin.Context, param string) (int, error) {
	return strconv.Atoi(c.Param(param))
}