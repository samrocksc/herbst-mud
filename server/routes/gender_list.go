package routes

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/dblog"
	"herbst-server/middleware"
	"herbst-server/repository"
)

// RegisterGenderRoutes registers REST endpoints for genders.
func RegisterGenderRoutes(r *gin.Engine, repos *repository.Container, client *db.Client) {
	genders := r.Group("/api/genders")
	genders.Use(middleware.AuthMiddleware(nil))
	genders.Use(middleware.AdminMiddleware())
	{
		genders.GET("", listGenders(repos))
		genders.GET("/:id", getGender(repos))
		genders.POST("", createGender(repos, client))
		genders.PUT("/:id", updateGender(repos, client))
		genders.DELETE("/:id", deleteGender(repos))
	}
}

// listGenders returns all genders for the specified world_id (default: "1").
func listGenders(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		worldID := c.Query("world_id")
		if worldID == "" {
			worldID = "1"
		}
		genders, err := repos.Gender.List(c.Request.Context(), worldID)
		if err != nil {
			dblog.Error("failed to list genders", err, slog.String("service", "genders"), slog.String("world_id", worldID))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		views := make([]genderView, 0, len(genders))
		for _, g := range genders {
			views = append(views, genderToView(g))
		}
		c.JSON(http.StatusOK, gin.H{"genders": views})
	}
}

// getGender returns a single gender by ID.
func getGender(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			slog.Warn("bad request", slog.String("service", "genders"), slog.String("reason", "invalid gender id"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid gender id"})
			return
		}
		g, err := repos.Gender.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "gender not found"})
			return
		}
		c.JSON(http.StatusOK, genderToView(g))
	}
}

// createGender creates a new gender.
func createGender(repos *repository.Container, client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name              string `json:"name" binding:"required"`
			DisplayName       string `json:"display_name"`
			SubjectPronoun    string `json:"subject_pronoun"`
			ObjectPronoun     string `json:"object_pronoun"`
			PossessivePronoun string `json:"possessive_pronoun"`
			WorldID           string `json:"world_id" default:"1"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			slog.Warn("bad request", slog.String("service", "genders"), slog.String("reason", "invalid request body"), slog.String("error", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
			return
		}

		// Check for duplicate name in this world
		existing, err := repos.Gender.GetByWorld(c.Request.Context(), req.Name, req.WorldID)
		if err == nil && existing != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "a gender with this name already exists in this world"})
			return
		}

		displayName := req.DisplayName
		if displayName == "" {
			displayName = req.Name
		}

		g, err := repos.Gender.Create(c.Request.Context(), repository.CreateGenderInput{
			Name:              req.Name,
			DisplayName:       displayName,
			SubjectPronoun:    req.SubjectPronoun,
			ObjectPronoun:     req.ObjectPronoun,
			PossessivePronoun: req.PossessivePronoun,
			WorldID:           req.WorldID,
		})
		if err != nil {
			dblog.Error("failed to create gender", err, slog.String("service", "genders"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		slog.Info("gender created", slog.String("service", "genders"), slog.String("gender_name", g.Name))
		c.JSON(http.StatusCreated, genderToView(g))
	}
}

// updateGender updates an existing gender.
func updateGender(repos *repository.Container, client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			slog.Warn("bad request", slog.String("service", "genders"), slog.String("reason", "invalid gender id"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid gender id"})
			return
		}

		var req struct {
			Name              *string `json:"name"`
			DisplayName       *string `json:"display_name"`
			SubjectPronoun    *string `json:"subject_pronoun"`
			ObjectPronoun     *string `json:"object_pronoun"`
			PossessivePronoun *string `json:"possessive_pronoun"`
			WorldID           *string `json:"world_id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			slog.Warn("bad request", slog.String("service", "genders"), slog.String("reason", "invalid request body"), slog.String("error", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		updates := repository.GenderUpdates{}
		if req.Name != nil {
			updates.Name = req.Name
		}
		if req.DisplayName != nil {
			updates.DisplayName = req.DisplayName
		}
		if req.SubjectPronoun != nil {
			updates.SubjectPronoun = req.SubjectPronoun
		}
		if req.ObjectPronoun != nil {
			updates.ObjectPronoun = req.ObjectPronoun
		}
		if req.PossessivePronoun != nil {
			updates.PossessivePronoun = req.PossessivePronoun
		}
		if req.WorldID != nil {
			updates.WorldID = req.WorldID
		}

		g, err := repos.Gender.Update(c.Request.Context(), id, updates)
		if err != nil {
			dblog.Error("failed to update gender", err, slog.String("service", "genders"), slog.Int("gender_id", id))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		slog.Info("gender updated", slog.String("service", "genders"), slog.Int("gender_id", id), slog.String("new_name", g.Name))
		c.JSON(http.StatusOK, genderToView(g))
	}
}

// deleteGender deletes a gender.
func deleteGender(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			slog.Warn("bad request", slog.String("service", "genders"), slog.String("reason", "invalid gender id"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid gender id"})
			return
		}
		if err := repos.Gender.Delete(c.Request.Context(), id); err != nil {
			dblog.Error("failed to delete gender", err, slog.String("service", "genders"), slog.Int("gender_id", id))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		slog.Info("gender deleted", slog.String("service", "genders"), slog.Int("gender_id", id))
		c.JSON(http.StatusOK, gin.H{"message": "gender deleted"})
	}
}

// --- View conversions ---

type genderView struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	DisplayName       string `json:"display_name"`
	SubjectPronoun    string `json:"subject_pronoun"`
	ObjectPronoun     string `json:"object_pronoun"`
	PossessivePronoun string `json:"possessive_pronoun"`
	WorldID           string `json:"world_id"`
}

func genderToView(g *db.Gender) genderView {
	return genderView{
		ID:                g.ID,
		Name:              g.Name,
		DisplayName:       g.DisplayName,
		SubjectPronoun:    g.SubjectPronoun,
		ObjectPronoun:     g.ObjectPronoun,
		PossessivePronoun: g.PossessivePronoun,
		WorldID:           g.WorldID,
	}
}
