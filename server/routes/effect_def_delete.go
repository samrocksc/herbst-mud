package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/effect"
	"herbst-server/db/effecthook"
)

func deleteEffectDef(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid effect id"})
			return
		}
		// Check if any hooks reference this effect
		hookCount, err := client.EffectHook.Query().
			Where(effecthook.HasEffectWith(effect.IDEQ(id))).
			Count(c.Request.Context())
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
		err = client.Effect.DeleteOneID(id).Exec(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "effect not found"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}