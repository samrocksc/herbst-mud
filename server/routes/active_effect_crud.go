package routes

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"herbst-server/dblog"
	"herbst-server/repository"
)

func listActiveEffects(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		charID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			slog.Warn("bad request", slog.String("service", "active_effects"), slog.String("reason", "invalid character id"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid character id"})
			return
		}
		effects, err := repos.ActiveEffect.ListActiveByCharacter(c.Request.Context(), charID)
		if err != nil {
			dblog.Error("failed to list active effects", err, slog.String("service", "active_effects"), slog.Int("character_id", charID))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result := make([]activeEffectView, len(effects))
		for i, ae := range effects {
			result[i] = activeEffectToView(ae)
		}
		c.JSON(http.StatusOK, gin.H{"active_effects": result})
	}
}

func removeActiveEffect(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		effectID, err := strconv.Atoi(c.Param("effect_id"))
		if err != nil {
			slog.Warn("bad request", slog.String("service", "active_effects"), slog.String("reason", "invalid effect id"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid effect id"})
			return
		}
		_, err = repos.ActiveEffect.Update(c.Request.Context(), effectID, repository.ActiveEffectUpdates{
			IsActive: ptrBool(false),
		})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "active effect not found"})
			return
		}
		slog.Info("active effect removed", slog.Int("effect_id", effectID), slog.String("user_email", c.GetString("email")), slog.String("service", "active_effects"))
		c.Status(http.StatusNoContent)
	}
}

func applyEffect(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		charID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			slog.Warn("bad request", slog.String("service", "active_effects"), slog.String("reason", "invalid character id"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid character id"})
			return
		}
		var input applyEffectInput
		if err := c.ShouldBindJSON(&input); err != nil {
			slog.Warn("bad request", slog.String("service", "active_effects"), slog.String("reason", "invalid json"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if _, err := repos.Character.Get(c.Request.Context(), charID); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "character not found"})
			return
		}
		eff, err := repos.Effect.Get(c.Request.Context(), input.EffectID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "effect not found"})
			return
		}
		// Handle stack mode — find existing active effect for same character+effect
		existing, _ := repos.ActiveEffect.GetActiveByCharacterAndEffect(c.Request.Context(), charID, input.EffectID)

		if existing != nil {
			switch eff.StackMode {
			case "replace":
				now := time.Now()
				repos.ActiveEffect.Update(c.Request.Context(), existing.ID, repository.ActiveEffectUpdates{
					StackCount: ptrInt(1),
					StartedAt:   &now,
				})
			case "refresh":
				if !eff.IsPermanent && eff.DurationSecs > 0 {
					expiresAt := time.Now().Add(time.Duration(eff.DurationSecs) * time.Second)
					repos.ActiveEffect.Update(c.Request.Context(), existing.ID, repository.ActiveEffectUpdates{
						ExpiresAt: &expiresAt,
					})
				}
			case "stack":
				if existing.StackCount < eff.StackLimit {
					repos.ActiveEffect.Update(c.Request.Context(), existing.ID, repository.ActiveEffectUpdates{
						StackCount: ptrInt(existing.StackCount + 1),
					})
				}
			}
			updated, _ := repos.ActiveEffect.GetWithEffect(c.Request.Context(), existing.ID)
			slog.Info("active effect updated", slog.Int("active_effect_id", updated.ID), slog.String("user_email", c.GetString("email")), slog.String("service", "active_effects"))
			c.JSON(http.StatusOK, activeEffectToView(updated))
			return
		}
		// Create new active effect
		var expiresAt *time.Time
		if input.DurationSecs != nil && *input.DurationSecs > 0 {
			t := time.Now().Add(time.Duration(*input.DurationSecs) * time.Second)
			expiresAt = &t
		} else if !eff.IsPermanent && eff.DurationSecs > 0 {
			t := time.Now().Add(time.Duration(eff.DurationSecs) * time.Second)
			expiresAt = &t
		}
		ae, err := repos.ActiveEffect.Create(c.Request.Context(), repository.CreateActiveEffectInput{
			CharacterID: charID,
			EffectID:    input.EffectID,
			AppliedByID: input.AppliedByID,
			StackCount:  1,
			ExpiresAt:   expiresAt,
		})
		if err != nil {
			dblog.Error("failed to create active effect", err, slog.String("service", "active_effects"), slog.Int("character_id", charID))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ae, _ = repos.ActiveEffect.GetWithEffect(c.Request.Context(), ae.ID)
		slog.Info("active effect applied", slog.Int("active_effect_id", ae.ID), slog.Int("character_id", charID), slog.String("user_email", c.GetString("email")), slog.String("service", "active_effects"))
		c.JSON(http.StatusCreated, activeEffectToView(ae))
	}
}

func ptrInt(v int) *int    { return &v }
func ptrBool(v bool) *bool { return &v }