package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/dialognode"
	"herbst-server/repository"
)

// TODO: Use repos.DialogNode.List once repo supports ordering and WithNpcTemplate edge loading
func listDialogNodes(repos *repository.Container, client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		nodes, err := client.DialogNode.Query().
			WithNpcTemplate().
			Order(dialognode.ByID()).
			All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result := make([]dialogNodeView, len(nodes))
		for i, n := range nodes {
			result[i] = dialogNodeToView(n)
		}
		c.JSON(http.StatusOK, gin.H{"dialog_nodes": result})
	}
}