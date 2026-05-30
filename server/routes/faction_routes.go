package routes

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/factioncategory"
	"herbst-server/dblog"
	"herbst-server/middleware"
	"herbst-server/repository"
)

// RegisterFactionRoutes registers REST endpoints for factions and faction categories.
func RegisterFactionRoutes(r *gin.Engine, repos *repository.Container, client *db.Client) {
	// Protected /api routes — all require JWT auth + admin check
	factions := r.Group("/api")
	factions.Use(middleware.AuthMiddleware(nil))
	factions.Use(middleware.AdminMiddleware())
	factions.Use(middleware.WorldAccessMiddleware())
	{
		factions.GET("/factions", listFactions(repos))
		factions.GET("/factions/:id", getFaction(repos))
		factions.POST("/factions", createFaction(repos))
		factions.PUT("/factions/:id", updateFaction(repos))
		factions.DELETE("/factions/:id", deleteFaction(repos))
		factions.GET("/factions/:id/members", getFactionMembers(repos))

		// TODO: Migrate to FactionCategoryRepo once created
		factions.GET("/faction-categories", listFactionCategories(client))
		factions.POST("/faction-categories", createFactionCategory(client))
		factions.DELETE("/faction-categories/:id", deleteFactionCategory(client))
	}
}

// ─── Faction CRUD ─────────────────────────────────────────────────────────────

func listFactions(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		worldID := c.Query("world_id")
		factions, err := repos.Faction.List(c.Request.Context(), worldID)
		if err != nil {
			dblog.Error("failed to list factions", err, slog.String("service", "factions"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if search := c.Query("search"); search != "" {
			s := strings.ToLower(search)
			filtered := make([]*db.Faction, 0, len(factions))
			for _, f := range factions {
				if strings.Contains(strings.ToLower(f.Name), s) {
					filtered = append(filtered, f)
				}
			}
			factions = filtered
		}
		result := make([]gin.H, len(factions))
		for i, f := range factions {
			result[i] = factionToJSON(f)
		}
		c.JSON(http.StatusOK, result)
	}
}

func getFaction(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			slog.Warn("bad request", slog.String("service", "factions"), slog.String("reason", "invalid faction id"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid faction id"})
			return
		}
		f, err := repos.Faction.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "faction not found"})
			return
		}
		c.JSON(http.StatusOK, factionToJSON(f))
	}
}

func createFaction(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name        string   `json:"name" binding:"required"`
			DisplayName string   `json:"display_name" binding:"required"`
			Description string   `json:"description"`
			CategoryID  int      `json:"category_id"`
			MemberTags  []string `json:"member_tags"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			slog.Warn("bad request", slog.String("service", "factions"), slog.String("reason", err.Error()), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// TODO: Add CategoryID to CreateFactionInput once repo supports it
		created, err := repos.Faction.Create(c.Request.Context(), repository.CreateFactionInput{
			Name:        req.Name,
			DisplayName: req.DisplayName,
			Description: req.Description,
			MemberTags:  req.MemberTags,
		})
		if err != nil {
			dblog.Error("failed to create faction", err, slog.String("service", "factions"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		slog.Info("faction created", slog.Int("faction_id", created.ID), slog.String("user_email", c.GetString("email")), slog.String("service", "factions"))
		c.JSON(http.StatusCreated, factionToJSON(created))
	}
}

func updateFaction(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			slog.Warn("bad request", slog.String("service", "factions"), slog.String("reason", "invalid faction id"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid faction id"})
			return
		}
		var req struct {
			DisplayName string   `json:"display_name"`
			Description string   `json:"description"`
			CategoryID  int      `json:"category_id"`
			MemberTags  []string `json:"member_tags"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			slog.Warn("bad request", slog.String("service", "factions"), slog.String("reason", err.Error()), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		updates := repository.FactionUpdates{
			DisplayName: &req.DisplayName,
			Description: &req.Description,
		}
		if req.MemberTags != nil {
			updates.MemberTags = req.MemberTags
		}
		_, err = repos.Faction.Update(c.Request.Context(), id, updates)
		if err != nil {
			dblog.Error("failed to update faction", err, slog.String("service", "factions"), slog.Int("faction_id", id))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Reload with edges for response
		f, err := repos.Faction.GetWithEdges(c.Request.Context(), id)
		if err != nil {
			dblog.Error("failed to reload faction after update", err, slog.String("service", "factions"), slog.Int("faction_id", id))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		slog.Info("faction updated", slog.Int("faction_id", f.ID), slog.String("user_email", c.GetString("email")), slog.String("service", "factions"))
		c.JSON(http.StatusOK, factionToJSON(f))
	}
}

func deleteFaction(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid faction id"})
			return
		}
		if err := repos.Faction.Delete(c.Request.Context(), id); err != nil {
			dblog.Error("failed to delete faction", err, slog.String("service", "factions"), slog.Int("faction_id", id))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		slog.Info("faction deleted", slog.Int("faction_id", id), slog.String("user_email", c.GetString("email")), slog.String("service", "factions"))
		c.JSON(http.StatusNoContent, nil)
	}
}

func getFactionMembers(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			slog.Warn("bad request", slog.String("service", "factions"), slog.String("reason", "invalid faction id"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid faction id"})
			return
		}
		memberships, err := repos.CharacterFaction.ListByFactionWithDetails(c.Request.Context(), id)
		if err != nil {
			dblog.Error("failed to list faction members", err, slog.String("service", "factions"), slog.Int("faction_id", id))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result := make([]gin.H, 0, len(memberships))
		for _, m := range memberships {
			char := m.Edges.Character
			if char == nil {
				continue
			}
			result = append(result, gin.H{
				"character_id": char.ID,
				"name":         char.Name,
				"reputation":   m.Reputation,
				"status":       m.Status,
				"joined_at":    m.JoinedAt.Format("2006-01-02T15:04:05Z"),
			})
		}
		c.JSON(http.StatusOK, result)
	}
}

// ─── Faction Category CRUD ────────────────────────────────────────────────────
// TODO: Migrate to FactionCategoryRepo once created

func listFactionCategories(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := client.FactionCategory.Query()
		if search := c.Query("search"); search != "" {
			query = query.Where(factioncategory.NameContains(search))
		}
		if c.Query("initial_config") == "true" {
			query = query.Where(factioncategory.InitialConfig(true))
		}
		cats, err := query.WithFactions().All(c.Request.Context())
		if err != nil {
			dblog.Error("failed to list faction categories", err, slog.String("service", "factions"))
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
			InitialConfig  bool   `json:"initial_config"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			slog.Warn("bad request", slog.String("service", "factions"), slog.String("reason", err.Error()), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		created, err := client.FactionCategory.Create().
			SetName(req.Name).
			SetDisplayName(req.DisplayName).
			SetDescription(req.Description).
			SetMaxMemberships(req.MaxMemberships).
			SetAutoJoin(req.AutoJoin).
			SetInitialConfig(req.InitialConfig).
			Save(c.Request.Context())
		if err != nil {
			dblog.Error("failed to create faction category", err, slog.String("service", "factions"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		slog.Info("faction category created", slog.Int("category_id", created.ID), slog.String("user_email", c.GetString("email")), slog.String("service", "factions"))
		c.JSON(http.StatusCreated, categoryToJSON(created))
	}
}

func deleteFactionCategory(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := parseIntParam(c, "id")
		if err != nil {
			slog.Warn("bad request", slog.String("service", "factions"), slog.String("reason", "invalid category id"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
			return
		}
		if err := client.FactionCategory.DeleteOneID(id).Exec(c.Request.Context()); err != nil {
			dblog.Error("failed to delete faction category", err, slog.String("service", "factions"), slog.Int("category_id", id))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		slog.Info("faction category deleted", slog.Int("category_id", id), slog.String("user_email", c.GetString("email")), slog.String("service", "factions"))
		c.JSON(http.StatusOK, gin.H{"deleted": id})
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
	if f.MemberTags != nil {
		result["member_tags"] = f.MemberTags
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
	result := gin.H{
		"id":              cat.ID,
		"name":            cat.Name,
		"display_name":    cat.DisplayName,
		"description":     cat.Description,
		"max_memberships": cat.MaxMemberships,
		"auto_join":       cat.AutoJoin,
		"initial_config":  cat.InitialConfig,
	}
	if cat.Edges.Factions != nil {
		factions := make([]gin.H, len(cat.Edges.Factions))
		for i, f := range cat.Edges.Factions {
			factions[i] = gin.H{
				"id":           f.ID,
				"name":         f.Name,
				"display_name": f.DisplayName,
			}
		}
		result["factions"] = factions
	}
	return result
}