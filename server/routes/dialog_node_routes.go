package routes

import (
	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/dialognode"
	"herbst-server/db/npctemplate"
	"herbst-server/middleware"
)

// RegisterDialogNodeRoutes registers CRUD endpoints for DialogNode definitions.
func RegisterDialogNodeRoutes(r *gin.Engine, client *db.Client) {
	nodes := r.Group("/api/dialog-nodes")
	nodes.Use(middleware.AuthMiddleware())
	nodes.Use(middleware.AdminMiddleware())
	{
		nodes.GET("", listDialogNodes(client))
		nodes.POST("", createDialogNode(client))
		nodes.GET("/:id", getDialogNode(client))
		nodes.PUT("/:id", updateDialogNode(client))
		nodes.DELETE("/:id", deleteDialogNode(client))
	}
	// Public: game client fetches dialog tree for a specific NPC template.
	r.GET("/api/npc-templates/:templateId/dialog-nodes", getDialogNodesForTemplate(client))
}

// getDialogNode returns a single dialog node by ID (with npc_template edge loaded).
func getDialogNode(client *db.Client) gin.HandlerFunc {
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
func getDialogNodesForTemplate(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		templateID := c.Param("templateId")
		if templateID == "" {
			c.JSON(400, gin.H{"error": "invalid template id"})
			return
		}
		// Verify the NPC template exists.
		_, err := client.NPCTemplate.Get(c.Request.Context(), templateID)
		if err != nil {
			c.JSON(404, gin.H{"error": "npc template not found"})
			return
		}
		nodes, err := client.DialogNode.Query().
			Where(dialognode.HasNpcTemplateWith(npctemplate.IDEQ(templateID))).
			Order(dialognode.ByID()).
			All(c.Request.Context())
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