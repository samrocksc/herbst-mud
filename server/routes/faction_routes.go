package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/characterfaction"
	"herbst-server/db/faction"
	"herbst-server/middleware"
)

// RegisterFactionRoutes registers REST endpoints for factions and faction categories.
func RegisterFactionRoutes(r *gin.Engine, client *db.Client) {
	// Protected /api routes — all require JWT auth + admin check
	factions := r.Group("/api")
	factions.Use(middleware.AuthMiddleware())
	factions.Use(middleware.AdminMiddleware())
	{
		factions.GET("/factions", listFactions(client))
		factions.GET("/factions/:id", getFaction(client))
		factions.POST("/factions", createFaction(client))
		factions.PUT("/factions/:id", updateFaction(client))
		factions.DELETE("/factions/:id", deleteFaction(client))
		factions.GET("/factions/:id/members", getFactionMembers(client))

		factions.GET("/faction-categories", listFactionCategories(client))
		factions.POST("/faction-categories", createFactionCategory(client))
	}
}

// ─── Faction CRUD ─────────────────────────────────────────────────────────────

func listFactions(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		factions, err := client.Faction.Query().
			WithCategory().
			WithRequiredTags().
			All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result := make([]gin.H, len(factions))
		for i, f := range factions {
			result[i] = factionToJSON(f)
		}
		c.JSON(http.StatusOK, result)
	}
}

func getFaction(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid faction id"})
			return
		}
		f, err := client.Faction.Query().
			Where(faction.ID(id)).
			WithCategory().
			WithRequiredTags().
			WithSkills().
			Only(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "faction not found"})
			return
		}
		c.JSON(http.StatusOK, factionToJSON(f))
	}
}

func createFaction(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name        string `json:"name" binding:"required"`
			DisplayName string `json:"display_name" binding:"required"`
			Description string `json:"description"`
			CategoryID  int    `json:"category_id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		created, err := client.Faction.Create().
			SetName(req.Name).
			SetDisplayName(req.DisplayName).
			SetDescription(req.Description).
			SetCategoryID(req.CategoryID).
			Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Reload with edges
		f, err := client.Faction.Query().
			Where(faction.ID(created.ID)).
			WithCategory().
			WithRequiredTags().
			Only(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, factionToJSON(f))
	}
}

func updateFaction(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid faction id"})
			return
		}
		var req struct {
			DisplayName string `json:"display_name"`
			Description string `json:"description"`
			CategoryID  int    `json:"category_id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		updated, err := client.Faction.UpdateOneID(id).
			SetDisplayName(req.DisplayName).
			SetDescription(req.Description).
			Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Reload with edges
		f, err := client.Faction.Query().
			Where(faction.ID(updated.ID)).
			WithCategory().
			WithRequiredTags().
			Only(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, factionToJSON(f))
	}
}

func deleteFaction(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid faction id"})
			return
		}
		_, err = client.Faction.Delete().Where(faction.ID(id)).Exec(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNoContent, nil)
	}
}

func getFactionMembers(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid faction id"})
			return
		}
		memberships, err := client.CharacterFaction.Query().
			Where(characterfaction.HasFactionWith(faction.ID(id))).
			WithCharacter().
			All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result := make([]gin.H, len(memberships))
		for i, m := range memberships {
			char := m.Edges.Character
			if char == nil {
				continue
			}
			result[i] = gin.H{
				"character_id": char.ID,
				"name":         char.Name,
				"reputation":    m.Reputation,
				"status":       m.Status,
				"joined_at":    m.JoinedAt.Format("2006-01-02T15:04:05Z"),
			}
		}
		c.JSON(http.StatusOK, result)
	}
}

// ─── Faction Category CRUD ────────────────────────────────────────────────────

func listFactionCategories(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		cats, err := client.FactionCategory.Query().All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result := make([]gin.H, len(cats))
		for i, cat := range cats {
			result[i] = categoryToJSON(cat)
		}
		c.JSON(http.StatusOK, result)
	}
}

func createFactionCategory(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name           string `json:"name" binding:"required"`
			DisplayName    string `json:"display_name" binding:"required"`
			Description    string `json:"description"`
			MaxMemberships int    `json:"max_memberships"`
			AutoJoin       bool   `json:"auto_join"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		created, err := client.FactionCategory.Create().
			SetName(req.Name).
			SetDisplayName(req.DisplayName).
			SetDescription(req.Description).
			SetMaxMemberships(req.MaxMemberships).
			SetAutoJoin(req.AutoJoin).
			Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, categoryToJSON(created))
	}
}

// ─── JSON Helpers ────────────────────────────────────────────────────────────

func factionToJSON(f *db.Faction) gin.H {
	result := gin.H{
		"id":           f.ID,
		"name":         f.Name,
		"display_name": f.DisplayName,
		"description":  f.Description,
	}
	if f.Edges.Category != nil {
		result["category"] = gin.H{
			"id":           f.Edges.Category.ID,
			"name":         f.Edges.Category.Name,
			"display_name": f.Edges.Category.DisplayName,
		}
	}
	if f.Edges.RequiredTags != nil {
		tags := make([]string, len(f.Edges.RequiredTags))
		for i, t := range f.Edges.RequiredTags {
			tags[i] = t.RequiredTag
		}
		result["required_tags"] = tags
	}
	if f.Edges.CharacterFactions != nil {
		result["member_count"] = len(f.Edges.CharacterFactions)
	}
	return result
}

func categoryToJSON(cat *db.FactionCategory) gin.H {
	return gin.H{
		"id":              cat.ID,
		"name":            cat.Name,
		"display_name":    cat.DisplayName,
		"description":     cat.Description,
		"max_memberships": cat.MaxMemberships,
		"auto_join":       cat.AutoJoin,
	}
}
