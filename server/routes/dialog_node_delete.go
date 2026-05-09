package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
)

// deleteDialogNode removes a dialog node by ID.
func deleteDialogNode(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid dialog node id"})
			return
		}
		err := client.DialogNode.DeleteOneID(id).Exec(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "dialog node not found"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}