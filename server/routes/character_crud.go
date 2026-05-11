package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/repository"
	"herbst-server/service"
)

// createCharacter handles POST /characters (admin/debug endpoint).
func createCharacter(svc *service.Container, repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name         string `json:"name" binding:"required"`
			Password     string `json:"password"`
			IsNPC        bool   `json:"isNPC"`
			CurrentRoom  int    `json:"currentRoomId"`
			StartingRoom int    `json:"startingRoomId"`
			UserID       int    `json:"userId"`
			IsAdmin      bool   `json:"isAdmin"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if len(req.Name) < 1 || len(req.Name) > 23 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Character name must be 1-23 characters"})
			return
		}
		for _, ch := range req.Name {
			if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Character name can only contain letters (a-z, A-Z)"})
				return
			}
		}
		if req.UserID > 0 {
			count, err := repos.Character.CountByUser(c.Request.Context(), req.UserID)
			if err == nil && count >= 3 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Maximum of 3 characters per user reached"})
				return
			}
		}
		char, err := repos.Character.Create(c.Request.Context(), repository.CreateCharacterInput{
			Name:      req.Name,
			Password:  req.Password,
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
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

// listCharacters handles GET /characters.
func listCharacters(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		characters, err := repos.Character.ListAll(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
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
		c.JSON(http.StatusNoContent, nil)
	}
}