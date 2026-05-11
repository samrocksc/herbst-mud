package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/repository"
)

func updateEffectDef(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid effect id"})
			return
		}
		_, err = repos.Effect.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "effect not found"})
			return
		}
		var input effectDefInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if input.EffectType != nil && !validEffectTypes[*input.EffectType] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid effect_type"})
			return
		}

		updated, err := repos.Effect.Update(c.Request.Context(), id, repository.EffectUpdates{
			Name:         input.Name,
			Description:  input.Description,
			EffectType:   input.EffectType,
			Parameters:   input.Parameters,
			StackMode:    input.StackMode,
			StackLimit:   input.StackLimit,
			IsPermanent:  input.IsPermanent,
			DurationSecs: input.DurationSecs,
			Messages:     input.Messages,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Re-fetch with hooks edge loaded for hook count in response
		updated, _ = repos.Effect.GetWithHooks(c.Request.Context(), updated.ID)
		c.JSON(http.StatusOK, effectDefToView(updated))
	}
}