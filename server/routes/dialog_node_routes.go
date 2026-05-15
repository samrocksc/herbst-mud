package routes

import (
	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/dialognode"
	"herbst-server/middleware"
	"herbst-server/repository"
)

// RegisterDialogNodeRoutes registers CRUD endpoints for DialogNode definitions.
func RegisterDialogNodeRoutes(r *gin.Engine, repos *repository.Container, client *db.Client) {
	nodes := r.Group("/api/dialog-nodes")
	nodes.Use(middleware.AuthMiddleware(nil))
	nodes.Use(middleware.AdminMiddleware())
	nodes.Use(middleware.WorldAccessMiddleware())
	{
		nodes.GET("", listDialogNodes(repos, client))
		nodes.POST("", createDialogNode(repos, client))
		nodes.GET("/:id", getDialogNode(repos, client))
		nodes.PUT("/:id", updateDialogNode(repos, client))
		nodes.DELETE("/:id", deleteDialogNode(repos))
	}
	// Public: game client fetches dialog tree for a specific NPC template.
	r.GET("/api/npc-templates/:id/dialog-nodes", middleware.AuthMiddleware(nil), middleware.AdminMiddleware(), middleware.WorldAccessMiddleware(), getDialogNodesForTemplate(repos))
}

// TODO: Use repos.DialogNode.Get once repo supports WithNpcTemplate edge loading
func getDialogNode(repos *repository.Container, client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(400, gin.H{"error": "invalid dialog node id"})
			return
		}
		dn, err := client.DialogNode.Query().
			Where(dialognode.IDEQ(id)).
			WithNpcTemplate().
			Only(c.Request.Context())
		if err != nil {
			c.JSON(404, gin.H{"error": "dialog node not found"})
			return
		}
		c.JSON(200, dialogNodeToView(dn))
	}
}

// getDialogNodesForTemplate returns all dialog nodes for an NPC template (public).
func getDialogNodesForTemplate(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		templateID := c.Param("id")
		if templateID == "" {
			c.JSON(400, gin.H{"error": "invalid template id"})
			return
		}
		// Verify the NPC template exists
		_, err := repos.NPCTemplate.Get(c.Request.Context(), templateID)
		if err != nil {
			c.JSON(404, gin.H{"error": "npc template not found"})
			return
		}
		nodes, err := repos.DialogNode.ListByTemplate(c.Request.Context(), templateID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		result := make([]dialogNodeView, len(nodes))
		for i, n := range nodes {
			result[i] = dialogNodeToViewSimple(n)
		}
		c.JSON(200, gin.H{"dialog_nodes": result})
	}
}