package routes

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/dblog"
	"herbst-server/repository"
)

var validEffectTypes = map[string]bool{
	"xp_drain": true, "xp_gain": true, "xp_set": true,
	"xp_multiplier": true,
	"bind_point_set": true, "hp_change": true, "stamina_change": true,
	"mana_change": true, "message": true, "teleport": true,
	"apply_effect": true, "tag_add": true, "tag_remove": true,
	"change_race": true, "change_class": true,
}

func createEffectDef(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input effectDefInput
		if err := c.ShouldBindJSON(&input); err != nil {
			slog.Warn("bad request", slog.String("service", "effects"), slog.String("reason", "invalid json"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if input.EffectType == nil || !validEffectTypes[*input.EffectType] {
			slog.Warn("bad request", slog.String("service", "effects"), slog.String("reason", "invalid effect_type"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid effect_type"})
			return
		}
		if input.Name == nil || *input.Name == "" {
			slog.Warn("bad request", slog.String("service", "effects"), slog.String("reason", "name is required"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
			return
		}

		desc := ""
		if input.Description != nil {
			desc = *input.Description
		}
		stackMode := ""
		if input.StackMode != nil {
			stackMode = *input.StackMode
		}
		stackLimit := 0
		if input.StackLimit != nil {
			stackLimit = *input.StackLimit
		}
		var params map[string]interface{}
		if input.Parameters != nil {
			params = *input.Parameters
		}
		var msgs map[string]string
		if input.Messages != nil {
			msgs = *input.Messages
		}
		e, err := repos.Effect.Create(c.Request.Context(), repository.CreateEffectInput{
			Name:         *input.Name,
			Description:  desc,
			EffectType:   *input.EffectType,
			Parameters:   params,
			StackMode:    stackMode,
			StackLimit:   stackLimit,
			IsPermanent:  input.IsPermanent != nil && *input.IsPermanent,
			DurationSecs: 0,
			Messages:     msgs,
		})
		if err != nil {
			dblog.Error("failed to create effect definition", err, slog.String("service", "effects"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Re-fetch with hooks edge loaded for hook count in response
		e, _ = repos.Effect.GetWithHooks(c.Request.Context(), e.ID)
		slog.Info("effect definition created", slog.Int("effect_id", e.ID), slog.String("service", "effects"), slog.String("user_email", c.GetString("email")))
		c.JSON(http.StatusCreated, effectDefToView(e))
	}
}
