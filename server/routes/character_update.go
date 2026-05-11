package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"herbst-server/repository"
)

// updateCharacter handles PUT /characters/:id.
func updateCharacter(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := getIDParam(c)
		if !ok {
			return
		}
		var req struct {
			Name         string  `json:"name"`
			IsNPC        *bool   `json:"isNPC"`
			CurrentRoom  *int    `json:"currentRoomId"`
			StartingRoom *int    `json:"startingRoomId"`
			RespawnRoom  *int    `json:"respawnRoomId"`
			IsAdmin      *bool   `json:"isAdmin"`
			IsTest       *bool   `json:"isTest"`
			Gender       string  `json:"gender"`
			Description  string  `json:"description"`
			LastSeenAt   *string `json:"lastSeenAt"`
			Level        *int    `json:"level"`
			XP           *int    `json:"xp"`
			HP           *int    `json:"hitpoints"`
			MaxHP        *int    `json:"maxHitpoints"`
			Stamina      *int    `json:"stamina"`
			MaxStamina   *int    `json:"maxStamina"`
			Mana         *int    `json:"mana"`
			MaxMana      *int    `json:"maxMana"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		updates := repository.CharacterUpdates{
			Name: strPtr(req.Name), Gender: strPtr(req.Gender), Description: strPtr(req.Description),
			IsNPC: req.IsNPC, CurrentRoomID: req.CurrentRoom, StartingRoomID: req.StartingRoom,
			RespawnRoomID: req.RespawnRoom, IsAdmin: req.IsAdmin, IsTest: req.IsTest,
			Level: req.Level, Xp: req.XP, Hitpoints: req.HP, MaxHitpoints: req.MaxHP,
			Stamina: req.Stamina, MaxStamina: req.MaxStamina, Mana: req.Mana, MaxMana: req.MaxMana,
		}
		if req.LastSeenAt != nil {
			if t, err := time.Parse(time.RFC3339, *req.LastSeenAt); err == nil {
				updates.LastSeenAt = &t
			}
		}
		if req.Name == "" {
			updates.Name = nil
		}
		if req.Gender == "" {
			updates.Gender = nil
		}
		if req.Description == "" {
			updates.Description = nil
		}
		char, err := repos.Character.Update(c.Request.Context(), id, updates)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}
		c.JSON(http.StatusOK, char)
	}
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}