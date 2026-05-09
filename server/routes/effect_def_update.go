package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/effect"
)

func updateEffectDef(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid effect id"})
			return
		}
		e, err := client.Effect.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "effect not found"})
			return
		}
		var input effectDefInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		mut := client.Effect.UpdateOne(e)
		if input.Name != nil {
			mut.SetName(*input.Name)
		}
		if input.Description != nil {
			mut.SetDescription(*input.Description)
		}
		if input.EffectType != nil {
			if !validEffectTypes[*input.EffectType] {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid effect_type"})
				return
			}
			mut.SetEffectType(*input.EffectType)
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
		updated, err := mut.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		updated, _ = client.Effect.Query().
			Where(effect.IDEQ(updated.ID)).
			WithHooks().
			Only(c.Request.Context())
		c.JSON(http.StatusOK, effectDefToView(updated))
	}
}