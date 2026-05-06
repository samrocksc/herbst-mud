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
		templates.GET("/npc-templates/:id", getNPCTemplate(client))
		templates.POST("/npc-templates", createNPCTemplate(client))
		templates.PUT("/npc-templates/:id", updateNPCTemplate(client))
		templates.DELETE("/npc-templates/:id", deleteNPCTemplate(client))
	}
}

// ─── NPCTemplate CRUD ─────────────────────────────────────────────────────────

// npcTemplateView is the JSON shape returned by the API.
type npcTemplateView struct {
	ID              string         `json:"id"`
	Name            string         `json:"name"`
	Description     string         `json:"description"`
	Race            string         `json:"race"`
	Disposition     string         `json:"disposition"`
	Level           int            `json:"level"`
	XpValue         int            `json:"xp_value"`
	Skills          map[string]int `json:"skills"`
	TradesWith      []string       `json:"trades_with"`
	Greeting        string         `json:"greeting"`
	RespawnRooms    []string       `json:"respawn_rooms"`
	RespawnCooldown int            `json:"respawn_cooldown"`
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
				ID:              t.ID,
				Name:            t.Name,
				Description:     t.Description,
				Race:            t.Race,
				Disposition:     string(t.Disposition),
				Level:           t.Level,
				XpValue:         t.XpValue,
				Skills:          t.Skills,
				TradesWith:      t.TradesWith,
				Greeting:        t.Greeting,
				RespawnRooms:    t.RespawnRooms,
				RespawnCooldown: t.RespawnCooldown,
			}
		}
		c.JSON(http.StatusOK, result)
	}
}

func getNPCTemplate(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid npc template id"})
			return
		}

		tmpl, err := client.NPCTemplate.Query().
			Where(npctemplate.IDEQ(id)).
			Only(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "npc template not found"})
			return
		}

		c.JSON(http.StatusOK, npcTemplateView{
			ID:              tmpl.ID,
			Name:            tmpl.Name,
			Description:     tmpl.Description,
			Race:            tmpl.Race,
			Disposition:     string(tmpl.Disposition),
			Level:           tmpl.Level,
			XpValue:         tmpl.XpValue,
			Skills:          tmpl.Skills,
			TradesWith:      tmpl.TradesWith,
			Greeting:        tmpl.Greeting,
			RespawnRooms:    tmpl.RespawnRooms,
			RespawnCooldown: tmpl.RespawnCooldown,
		})
	}
}

func createNPCTemplate(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ID              string         `json:"id"`
			Name            string         `json:"name"`
			Description     string         `json:"description"`
			Race            string         `json:"race"`
			Disposition     string         `json:"disposition"`
			Level           int            `json:"level"`
			XpValue         int            `json:"xp_value"`
			Skills          map[string]int `json:"skills"`
			TradesWith      []string       `json:"trades_with"`
			Greeting        string         `json:"greeting"`
			RespawnRooms    []string       `json:"respawn_rooms"`
			RespawnCooldown int            `json:"respawn_cooldown"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if req.ID == "" || req.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id and name are required"})
			return
		}

		builder := client.NPCTemplate.Create().
			SetID(req.ID).
			SetName(req.Name).
			SetDescription(req.Description).
			SetRace(req.Race).
			SetLevel(req.Level).
			SetXpValue(req.XpValue).
			SetSkills(req.Skills).
			SetTradesWith(req.TradesWith).
			SetGreeting(req.Greeting).
			SetRespawnRooms(req.RespawnRooms).
			SetRespawnCooldown(req.RespawnCooldown)

		if req.Disposition != "" {
			switch req.Disposition {
			case "hostile", "friendly", "neutral":
				builder.SetDisposition(npctemplate.Disposition(req.Disposition))
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid disposition: " + req.Disposition})
				return
			}
		}

		created, err := builder.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, npcTemplateView{
			ID:              created.ID,
			Name:            created.Name,
			Description:     created.Description,
			Race:            created.Race,
			Disposition:     string(created.Disposition),
			Level:           created.Level,
			XpValue:         created.XpValue,
			Skills:          created.Skills,
			TradesWith:      created.TradesWith,
			Greeting:        created.Greeting,
			RespawnRooms:    created.RespawnRooms,
			RespawnCooldown: created.RespawnCooldown,
		})
	}
}

// updateNPCTemplateRequest accepts all template fields as optional pointers.
// Only non-nil fields are applied.
type updateNPCTemplateRequest struct {
	Name            *string         `json:"name"`
	Description     *string         `json:"description"`
	Race            *string         `json:"race"`
	Disposition     *string         `json:"disposition"`
	Level           *int            `json:"level"`
	XpValue         *int            `json:"xp_value"`
	Skills          *map[string]int `json:"skills"`
	TradesWith      *[]string       `json:"trades_with"`
	Greeting        *string         `json:"greeting"`
	RespawnRooms    *[]string       `json:"respawn_rooms"`
	RespawnCooldown *int            `json:"respawn_cooldown"`
}

func updateNPCTemplate(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid npc template id"})
			return
		}

		var req updateNPCTemplateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updater := client.NPCTemplate.UpdateOneID(id)

		if req.Name != nil {
			updater.SetName(*req.Name)
		}
		if req.Description != nil {
			updater.SetDescription(*req.Description)
		}
		if req.Race != nil {
			updater.SetRace(*req.Race)
		}
		if req.Disposition != nil {
			switch *req.Disposition {
			case "hostile", "friendly", "neutral":
				updater.SetDisposition(npctemplate.Disposition(*req.Disposition))
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid disposition: " + *req.Disposition})
				return
			}
		}
		if req.Level != nil {
			updater.SetLevel(*req.Level)
		}
		if req.XpValue != nil {
			updater.SetXpValue(*req.XpValue)
		}
		if req.Skills != nil {
			updater.SetSkills(*req.Skills)
		}
		if req.TradesWith != nil {
			updater.SetTradesWith(*req.TradesWith)
		}
		if req.Greeting != nil {
			updater.SetGreeting(*req.Greeting)
		}
		if req.RespawnRooms != nil {
			updater.SetRespawnRooms(*req.RespawnRooms)
		}
		if req.RespawnCooldown != nil {
			updater.SetRespawnCooldown(*req.RespawnCooldown)
		}

		updated, err := updater.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, npcTemplateView{
			ID:              updated.ID,
			Name:            updated.Name,
			Description:     updated.Description,
			Race:            updated.Race,
			Disposition:     string(updated.Disposition),
			Level:           updated.Level,
			XpValue:         updated.XpValue,
			Skills:          updated.Skills,
			TradesWith:      updated.TradesWith,
			Greeting:        updated.Greeting,
			RespawnRooms:    updated.RespawnRooms,
			RespawnCooldown: updated.RespawnCooldown,
		})
	}
}

// deleteNPCTemplate removes an NPC template by ID.
func deleteNPCTemplate(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid npc template id"})
			return
		}

		err := client.NPCTemplate.DeleteOneID(id).Exec(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "npc template not found"})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}
