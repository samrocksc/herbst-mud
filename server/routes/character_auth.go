package routes

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/repository"
	"golang.org/x/crypto/bcrypt"
)

// authenticateCharacter handles POST /characters/authenticate.
func authenticateCharacter(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name     string `json:"name" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		char, err := repos.Character.GetByName(c.Request.Context(), req.Name)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Character not found"})
			return
		}
		if char.Password == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Character has no password set"})
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(char.Password), []byte(req.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
			return
		}
		tokenBytes := make([]byte, 32)
		rand.Read(tokenBytes)
		accessToken := hex.EncodeToString(tokenBytes)
		c.JSON(http.StatusOK, gin.H{
			"authenticated": true,
			"access_token":  accessToken,
			"character": gin.H{
				"id":       char.ID,
				"name":     char.Name,
				"is_admin": char.IsAdmin,
			},
		})
	}
}

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