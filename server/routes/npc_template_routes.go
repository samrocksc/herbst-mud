package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/npctemplate"
	"herbst-server/middleware"
)

// RegisterNPCTemplateRoutes registers REST endpoints for NPC templates.
func RegisterNPCTemplateRoutes(r *gin.Engine, client *db.Client) {
	// Protected /api routes — all require JWT auth + admin check
	templates := r.Group("/api")
	templates.Use(middleware.AuthMiddleware())
	templates.Use(middleware.AdminMiddleware())
	{
		templates.GET("/npc-templates", listNPCTemplates(client))
		templates.PUT("/npc-templates/:id", updateNPCTemplate(client))
	}
}

// ─── NPCTemplate CRUD ─────────────────────────────────────────────────────────

// npcTemplateView is the JSON shape returned by the API.
type npcTemplateView struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Level   int    `json:"level"`
	XpValue int    `json:"xp_value"`
}

func listNPCTemplates(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		templates, err := client.NPCTemplate.Query().
			Order(npctemplate.ByName()).
			All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result := make([]npcTemplateView, len(templates))
		for i, t := range templates {
			result[i] = npcTemplateView{
				ID:      t.ID,
				Name:    t.Name,
				Level:   t.Level,
				XpValue: t.XpValue,
			}
		}
		c.JSON(http.StatusOK, result)
	}
}

func updateNPCTemplate(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid npc template id"})
			return
		}

		var req struct {
			XpValue *int `json:"xp_value"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if req.XpValue == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "xp_value is required"})
			return
		}

		updated, err := client.NPCTemplate.UpdateOneID(id).
			SetXpValue(*req.XpValue).
			Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, npcTemplateView{
			ID:      updated.ID,
			Name:    updated.Name,
			Level:   updated.Level,
			XpValue: updated.XpValue,
		})
	}
}