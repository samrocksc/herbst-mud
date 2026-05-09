package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/effect"
)

var validEffectTypes = map[string]bool{
	"xp_drain": true, "xp_gain": true, "xp_set": true,
	"bind_point_set": true, "hp_change": true, "stamina_change": true,
	"mana_change": true, "message": true, "teleport": true,
	"apply_effect": true, "tag_add": true, "tag_remove": true,
}

func createEffectDef(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input effectDefInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if input.EffectType == nil || !validEffectTypes[*input.EffectType] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid effect_type"})
			return
		}
		if input.Name == nil || *input.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
			return
		}
		mut := client.Effect.Create().
		SetName(*input.Name).
			SetEffectType(*input.EffectType)
		if input.Description != nil {
			mut.SetDescription(*input.Description)
		}
		if input.Parameters != nil {
			mut.SetParameters(*input.Parameters)
		}
		if input.StackMode != nil {
			mut.SetStackMode(*input.StackMode)
		}
		if input.StackLimit != nil {
			mut.SetStackLimit(*input.StackLimit)
		}
		if input.IsPermanent != nil {
			mut.SetIsPermanent(*input.IsPermanent)
		}
		if input.DurationSecs != nil {
			mut.SetDurationSecs(*input.DurationSecs)
		}
		if input.Messages != nil {
			mut.SetMessages(*input.Messages)
		}
		e, err := mut.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		e, _ = client.Effect.Query().
			Where(effect.IDEQ(e.ID)).
			WithHooks().
			Only(c.Request.Context())
		c.JSON(http.StatusCreated, effectDefToView(e))
	}
}