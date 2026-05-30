package routes

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/dblog"
	"herbst-server/repository"
)

func deleteEffectDef(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			slog.Warn("bad request", slog.String("service", "effect_defs"), slog.String("reason", "invalid effect id"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid effect id"})
			return
		}
		hookCount, err := repos.EffectHook.CountByEffect(c.Request.Context(), id)
		if err != nil {
			dblog.Error("failed to count effect hooks", err, slog.String("service", "effect_defs"), slog.Int("effect_id", id))
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
		slog.Info("effect definition deleted", slog.Int("effect_id", id), slog.String("user_email", c.GetString("email")), slog.String("service", "effect_defs"))
		c.Status(http.StatusNoContent)
	}
}