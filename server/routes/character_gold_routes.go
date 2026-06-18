package routes

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/dblog"
	"herbst-server/repository"
)

// getCharacterGold returns the current gold balance for a character.
func getCharacterGold(repos *repository.Container) gin.HandlerFunc {
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

		c.JSON(http.StatusOK, gin.H{
			"character_id":  id,
			"gold_credits": char.GoldCredits,
		})
	}
}

// addCharacterGold adds gold to a character's balance.
func addCharacterGold(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := getIDParam(c)
		if !ok {
			return
		}

		var req struct {
			Amount int    `json:"amount" binding:"required"`
			Source string `json:"source"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			slog.Warn("bad request", slog.String("service", "characters"), slog.String("reason", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		char, err := repos.Character.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}

		newBalance := char.GoldCredits + req.Amount
		updated, err := repos.Character.Update(c.Request.Context(), id, repository.CharacterUpdates{
			GoldCredits: &newBalance,
		})
		if err != nil {
			dblog.Error("failed to add gold", err, slog.String("service", "characters"), slog.Int("character_id", id))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		slog.Info("gold added", slog.String("service", "characters"), slog.Int("character_id", id), slog.Int("amount", req.Amount), slog.String("source", req.Source))
		c.JSON(http.StatusOK, gin.H{
			"character_id":  id,
			"gold_credits": updated.GoldCredits,
			"amount_added":  req.Amount,
		})
	}
}

// spendCharacterGold removes gold from a character's balance.
func spendCharacterGold(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := getIDParam(c)
		if !ok {
			return
		}

		var req struct {
			Amount int    `json:"amount" binding:"required"`
			Reason string `json:"reason"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			slog.Warn("bad request", slog.String("service", "characters"), slog.String("reason", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		char, err := repos.Character.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}

		if char.GoldCredits < req.Amount {
			c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient gold"})
			return
		}

		newBalance := char.GoldCredits - req.Amount
		updated, err := repos.Character.Update(c.Request.Context(), id, repository.CharacterUpdates{
			GoldCredits: &newBalance,
		})
		if err != nil {
			dblog.Error("failed to spend gold", err, slog.String("service", "characters"), slog.Int("character_id", id))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		slog.Info("gold spent", slog.String("service", "characters"), slog.Int("character_id", id), slog.Int("amount", req.Amount), slog.String("reason", req.Reason))
		c.JSON(http.StatusOK, gin.H{
			"character_id":  id,
			"gold_credits": updated.GoldCredits,
			"amount_spent": req.Amount,
		})
	}
}
