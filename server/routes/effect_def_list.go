package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/repository"
)

func listEffectDefs(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		effects, err := repos.Effect.ListWithHooks(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result := make([]effectDefView, len(effects))
		for i, e := range effects {
			result[i] = effectDefToView(e)
		}
		c.JSON(http.StatusOK, gin.H{"effects": result})
	}
}

func getEffectDef(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid effect id"})
			return
		}
		e, err := repos.Effect.GetWithHooks(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "effect not found"})
			return
		}
		c.JSON(http.StatusOK, effectDefToView(e))
	}
}