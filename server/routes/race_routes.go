package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/race"
	"herbst-server/middleware"
)

// RegisterRaceRoutes registers REST endpoints for races.
func RegisterRaceRoutes(r *gin.Engine, client *db.Client) {
	// Protected /api routes — all require JWT auth + admin check
	races := r.Group("/api")
	races.Use(middleware.AuthMiddleware())
	races.Use(middleware.AdminMiddleware())
	{
		races.GET("/races", listRaces(client))
		races.POST("/races", createRace(client))
		races.GET("/races/:id", getRace(client))
		races.PUT("/races/:id", updateRace(client))
		races.DELETE("/races/:id", deleteRace(client))
	}
}

// ─── Race view struct ─────────────────────────────────────────────────────────

type raceView struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	DisplayName   string `json:"display_name"`
	Description   string `json:"description"`
	StatModifiers string `json:"stat_modifiers"`
	SkillGrants   string `json:"skill_grants"`
	IsPlayable    bool   `json:"is_playable"`
}

func raceToView(r *db.Race) raceView {
	return raceView{
		ID:            r.ID,
		Name:          r.Name,
		DisplayName:   r.DisplayName,
		Description:   r.Description,
		StatModifiers: r.StatModifiers,
		SkillGrants:   r.SkillGrants,
		IsPlayable:    r.IsPlayable,
	}
}

// ─── Race CRUD handlers ───────────────────────────────────────────────────────

func listRaces(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		races, err := client.Race.Query().
			Order(race.ByDisplayName()).
			All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result := make([]raceView, len(races))
		for i, r := range races {
			result[i] = raceToView(r)
		}
		c.JSON(http.StatusOK, result)
	}
}

func getRace(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid race id"})
			return
		}
		r, err := client.Race.Query().
			Where(race.ID(id)).
			Only(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "race not found"})
			return
		}
		c.JSON(http.StatusOK, raceToView(r))
	}
}

func createRace(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name          string `json:"name" binding:"required"`
			DisplayName   string `json:"display_name" binding:"required"`
			Description   string `json:"description"`
			StatModifiers string `json:"stat_modifiers"`
			SkillGrants   string `json:"skill_grants"`
			IsPlayable    bool   `json:"is_playable"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		r, err := client.Race.Create().
			SetName(req.Name).
			SetDisplayName(req.DisplayName).
			SetDescription(req.Description).
			SetStatModifiers(req.StatModifiers).
			SetSkillGrants(req.SkillGrants).
			SetIsPlayable(req.IsPlayable).
			Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, raceToView(r))
	}
}

func updateRace(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid race id"})
			return
		}
		var req struct {
			DisplayName   *string `json:"display_name"`
			Description   *string `json:"description"`
			StatModifiers *string `json:"stat_modifiers"`
			SkillGrants   *string `json:"skill_grants"`
			IsPlayable    *bool   `json:"is_playable"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		r, err := client.Race.Query().
			Where(race.ID(id)).
			Only(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "race not found"})
			return
		}
		if req.DisplayName != nil {
			r.DisplayName = *req.DisplayName
		}
		if req.Description != nil {
			r.Description = *req.Description
		}
		if req.StatModifiers != nil {
			r.StatModifiers = *req.StatModifiers
		}
		if req.SkillGrants != nil {
			r.SkillGrants = *req.SkillGrants
		}
		if req.IsPlayable != nil {
			r.IsPlayable = *req.IsPlayable
		}
		if err := client.Race.UpdateOne(r).Exec(c.Request.Context()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, raceToView(r))
	}
}

func deleteRace(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid race id"})
			return
		}
		deleted, err := client.Race.Delete().
			Where(race.ID(id)).
			Exec(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if deleted == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "race not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "race deleted"})
	}
}
