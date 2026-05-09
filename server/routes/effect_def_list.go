package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/effect"
)

func listEffectDefs(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		effects, err := client.Effect.Query().
			Order(db.Asc(effect.FieldName)).
			WithHooks().
			All(c.Request.Context())
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

func getEffectDef(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid effect id"})
			return
		}
		e, err := client.Effect.Query().
			Where(effect.IDEQ(id)).
			WithHooks().
			Only(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "effect not found"})
			return
		}
		c.JSON(http.StatusOK, effectDefToView(e))
	}
}