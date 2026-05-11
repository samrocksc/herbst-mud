package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/repository"
)

func deleteEffectDef(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid effect id"})
			return
		}
		hookCount, err := repos.EffectHook.CountByEffect(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if hookCount > 0 {
			c.JSON(http.StatusConflict, gin.H{
				"error":      "cannot delete effect: referenced by hooks",
				"hook_count": hookCount,
			})
			return
		}
		err = repos.Effect.Delete(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "effect not found"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}