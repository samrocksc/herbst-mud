package routes

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"herbst-server/dblog"
	"herbst-server/db"
	"herbst-server/repository"
)

func listEffectDefs(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		effects, err := repos.Effect.ListWithHooks(c.Request.Context())
		if err != nil {
			dblog.Error("failed to list effect definitions", err, slog.String("service", "effects"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if search := c.Query("search"); search != "" {
			s := strings.ToLower(search)
			filtered := make([]*db.Effect, 0, len(effects))
			for _, e := range effects {
				if strings.Contains(strings.ToLower(e.Name), s) {
					filtered = append(filtered, e)
				}
			}
			effects = filtered
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
			slog.Warn("bad request", slog.String("service", "effects"), slog.String("reason", "invalid effect id"), slog.String("client_ip", c.ClientIP()))
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
