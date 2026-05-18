package routes

import (
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/middleware"
	"herbst-server/repository"
)

// RegisterNPCTemplateRoutes registers REST endpoints for NPC templates.
func RegisterNPCTemplateRoutes(r *gin.Engine, repos *repository.Container) {
	// Protected /api routes — all require JWT auth + admin check + world access
	templates := r.Group("/api")
	templates.Use(middleware.AuthMiddleware(nil))
	templates.Use(middleware.AdminMiddleware())
	templates.Use(middleware.WorldAccessMiddleware())
	{
		templates.GET("/npc-templates", listNPCTemplates(repos))
		templates.GET("/npc-templates/:id", getNPCTemplate(repos))
		templates.POST("/npc-templates", createNPCTemplate(repos))
		templates.PUT("/npc-templates/:id", updateNPCTemplate(repos))
		templates.DELETE("/npc-templates/:id", deleteNPCTemplate(repos))
	}
}

// ─── NPCTemplate CRUD ─────────────────────────────────────────────────────────

// npcTemplateView is the JSON shape returned by the API.
type npcTemplateView struct {
	ID              string         `json:"id"`
	Slug            string         `json:"slug"`
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
	WorldID         string         `json:"world_id"`
}

func listNPCTemplates(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get world_id from query parameter
		worldID := c.Query("world_id")
		templates, err := repos.NPCTemplate.List(c.Request.Context(), worldID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if search := c.Query("search"); search != "" {
			s := strings.ToLower(search)
			filtered := make([]*db.NPCTemplate, 0, len(templates))
			for _, t := range templates {
				if strings.Contains(strings.ToLower(t.Name), s) {
					filtered = append(filtered, t)
				}
			}
			templates = filtered
		}
		sort.Slice(templates, func(i, j int) bool {
			return templates[i].Name < templates[j].Name
		})
		result := make([]npcTemplateView, len(templates))
		for i, t := range templates {
			result[i] = npcTemplateView{
				ID:              t.ID,
				Slug:            t.Slug,
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
				WorldID:         t.WorldID,
			}
		}
		c.JSON(http.StatusOK, result)
	}
}

func getNPCTemplate(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid npc template id"})
			return
		}

		tmpl, err := repos.NPCTemplate.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "npc template not found"})
			return
		}

		c.JSON(http.StatusOK, npcTemplateView{
			ID:              tmpl.ID,
			Slug:            tmpl.Slug,
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
			WorldID:         tmpl.WorldID,
		})
	}
}

func createNPCTemplate(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ID              string         `json:"id"`
			Slug            string         `json:"slug"`
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
		if req.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
			return
		}

		disposition := "neutral"
		if req.Disposition != "" {
			switch req.Disposition {
			case "hostile", "friendly", "neutral":
				disposition = req.Disposition
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid disposition: " + req.Disposition})
				return
			}
		}

		cooldown := req.RespawnCooldown
		created, err := repos.NPCTemplate.Create(c.Request.Context(), repository.CreateNPCTemplateInput{
			ID:              req.ID,
			Slug:            req.Slug,
			Name:            req.Name,
			Description:     req.Description,
			Race:            req.Race,
			Disposition:     disposition,
			Level:           req.Level,
			XPValue:         req.XpValue,
			Skills:          req.Skills,
			TradesWith:      req.TradesWith,
			Greeting:        req.Greeting,
			RespawnRooms:    req.RespawnRooms,
			RespawnCooldown: &cooldown,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, npcTemplateView{
			ID:              created.ID,
			Slug:            created.Slug,
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
		Slug            *string         `json:"slug"`
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
	WorldID         *string         `json:"world_id"`
}

func updateNPCTemplate(repos *repository.Container) gin.HandlerFunc {
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

		updates := repository.NPCTemplateUpdates{
			Name:            req.Name,
			Slug:            req.Slug,
			Description:     req.Description,
			Race:            req.Race,
			Level:           req.Level,
			XPValue:         req.XpValue,
			Skills:          req.Skills,
			TradesWith:      req.TradesWith,
			Greeting:        req.Greeting,
			RespawnRooms:    req.RespawnRooms,
			RespawnCooldown: req.RespawnCooldown,
		}
		if req.Disposition != nil {
			switch *req.Disposition {
			case "hostile", "friendly", "neutral":
				updates.Disposition = req.Disposition
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid disposition: " + *req.Disposition})
				return
			}
		}

		updated, err := repos.NPCTemplate.Update(c.Request.Context(), id, updates)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, npcTemplateView{
			ID:              updated.ID,
			Slug:            updated.Slug,
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
func deleteNPCTemplate(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid npc template id"})
			return
		}

		err := repos.NPCTemplate.Delete(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "npc template not found"})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}