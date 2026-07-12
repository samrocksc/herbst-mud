package routes

import (
	"log/slog"
	"herbst-server/dblog"

	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/repository"
)

// deleteDialogNode removes a dialog node by ID.
func deleteDialogNode(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid dialog node id"})
			return
		}
		err := repos.DialogNode.Delete(c.Request.Context(), id)
		if err != nil {
			dblog.Error("delete dialog node failed", err, slog.String("service", "dialog_nodes"))
			c.JSON(http.StatusNotFound, gin.H{"error": "dialog node not found"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}