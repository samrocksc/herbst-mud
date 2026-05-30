package routes

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/db/charactercompetency"
	"herbst-server/db/competencycategory"
	"herbst-server/db/competencylevelthreshold"
	"herbst-server/dblog"
	"herbst-server/middleware"
	"herbst-server/repository"
)

func RegisterCompetencyRoutes(r *gin.Engine, repos *repository.Container, client *db.Client) {
	competencies := r.Group("/api")
	competencies.Use(middleware.AuthMiddleware(nil))
	competencies.Use(middleware.AdminMiddleware())
	{
		competencies.GET("/competency-categories", listCompetencyCategories(repos))
		competencies.POST("/competency-categories", createCompetencyCategory(repos, client))
		competencies.PUT("/competency-categories/:id", updateCompetencyCategory(repos, client))
		competencies.DELETE("/competency-categories/:id", deleteCompetencyCategory(repos, client))
		competencies.GET("/characters/:id/competencies", listCharacterCompetencies(repos, client))
	}
}

func listCompetencyCategories(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		cats, err := repos.Competency.ListCategories(c.Request.Context())
		if err != nil {
			dblog.Error("failed to list competency categories", err, slog.String("service", "competencies"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result := make([]gin.H, len(cats))
		for i, cat := range cats {
			result[i] = compCatToJSON(cat)
		}
		c.JSON(http.StatusOK, result)
	}
}

type thresholdInput struct {
	Level             int     `json:"level"`
	XpRequired        int     `json:"xp_required"`
	DamageMultiplier  float64 `json:"damage_multiplier"`
	DefenseMultiplier float64 `json:"defense_multiplier"`
}

func createCompetencyCategory(repos *repository.Container, client *db.Client) gin.HandlerFunc {
	var req struct {
		ID           string           `json:"id" binding:"required"`
		Name         string           `json:"name" binding:"required"`
		XpMultiplier float64         `json:"xp_multiplier"`
		Thresholds   []thresholdInput `json:"thresholds"`
	}
	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(&req); err != nil {
			slog.Warn("bad request", slog.String("service", "competencies"), slog.String("reason", "invalid json"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		cat, err := repos.Competency.CreateCategory(c.Request.Context(), repository.CreateCompetencyInput{
			ID:           req.ID,
			Name:         req.Name,
			XPMultiplier: req.XpMultiplier,
		})
		if err != nil {
			dblog.Error("failed to create competency category", err, slog.String("service", "competencies"), slog.String("category_id", req.ID))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		for _, t := range req.Thresholds {
			_, err := client.CompetencyLevelThreshold.Create().
				SetID(req.ID+"-"+strconv.Itoa(t.Level)).
				SetLevel(t.Level).
				SetXpRequired(t.XpRequired).
				SetDamageMultiplier(t.DamageMultiplier).
				SetDefenseMultiplier(t.DefenseMultiplier).
				SetCategoryID(cat.ID).
				Save(c.Request.Context())
			if err != nil {
				dblog.Error("failed to create competency threshold", err, slog.String("service", "competencies"), slog.String("threshold_id", req.ID+"-"+strconv.Itoa(t.Level)))
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		// Reload with thresholds
		cat, err = client.CompetencyCategory.Query().
			Where(competencycategory.ID(cat.ID)).
			WithThresholds().
			Only(c.Request.Context())
		if err != nil {
			dblog.Error("failed to reload competency category", err, slog.String("service", "competencies"), slog.String("category_id", cat.ID))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		slog.Info("competency category created", slog.String("category_id", cat.ID), slog.String("user_email", c.GetString("email")), slog.String("service", "competencies"))
		c.JSON(http.StatusCreated, compCatToJSON(cat))
	}
}

func updateCompetencyCategory(repos *repository.Container, client *db.Client) gin.HandlerFunc {
	var req struct {
		Name         string           `json:"name"`
		XpMultiplier *float64         `json:"xp_multiplier"`
		Thresholds   []thresholdInput `json:"thresholds"`
	}
	return func(c *gin.Context) {
		id := c.Param("id")
		_, err := repos.Competency.GetCategory(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
			return
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			slog.Warn("bad request", slog.String("service", "competencies"), slog.String("reason", "invalid json"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Update category via repo
		cat, err := repos.Competency.UpdateCategory(c.Request.Context(), id, repository.CompetencyCategoryUpdates{
			Name:         &req.Name,
			XPMultiplier: req.XpMultiplier,
		})
		if err != nil {
			dblog.Error("failed to update competency category", err, slog.String("service", "competencies"), slog.String("category_id", id))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// TODO: Migrate threshold CRUD to CompetencyRepo
		if req.Thresholds != nil {
			_, err := client.CompetencyLevelThreshold.Delete().
				Where(competencylevelthreshold.HasCategoryWith(competencycategory.ID(id))).
				Exec(c.Request.Context())
			if err != nil {
				dblog.Error("failed to delete competency thresholds", err, slog.String("service", "competencies"), slog.String("category_id", id))
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			for _, t := range req.Thresholds {
				_, err := client.CompetencyLevelThreshold.Create().
					SetID(id+"-"+strconv.Itoa(t.Level)).
					SetLevel(t.Level).
					SetXpRequired(t.XpRequired).
					SetDamageMultiplier(t.DamageMultiplier).
					SetDefenseMultiplier(t.DefenseMultiplier).
					SetCategoryID(id).
					Save(c.Request.Context())
				if err != nil {
					dblog.Error("failed to create competency threshold", err, slog.String("service", "competencies"), slog.String("threshold_id", id+"-"+strconv.Itoa(t.Level)))
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}
		}
		// Reload with thresholds
		cat, err = repos.Competency.GetCategoryWithThresholds(c.Request.Context(), id)
		if err != nil {
			dblog.Error("failed to reload competency category", err, slog.String("service", "competencies"), slog.String("category_id", id))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		slog.Info("competency category updated", slog.String("category_id", id), slog.String("user_email", c.GetString("email")), slog.String("service", "competencies"))
		c.JSON(http.StatusOK, compCatToJSON(cat))
	}
}

func deleteCompetencyCategory(repos *repository.Container, client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		count, err := repos.Competency.CountCompetenciesByCategory(c.Request.Context(), id)
		if err != nil {
			dblog.Error("failed to count competencies by category", err, slog.String("service", "competencies"), slog.String("category_id", id))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Cannot delete: characters have XP in this category"})
			return
		}
		if err := repos.Competency.DeleteCategory(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
			return
		}
		slog.Info("competency category deleted", slog.String("category_id", id), slog.String("user_email", c.GetString("email")), slog.String("service", "competencies"))
		c.Status(http.StatusNoContent)
	}
}

func listCharacterCompetencies(repos *repository.Container, client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		charIDStr := c.Param("id")
		charID, err := strconv.Atoi(charIDStr)
		if err != nil {
			slog.Warn("bad request", slog.String("service", "competencies"), slog.String("reason", "invalid character id"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid character id"})
			return
		}
		if _, err := repos.Character.Get(c.Request.Context(), charID); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "character not found"})
			return
		}
		// TODO: Add ListByCharacterWithDetails to CompetencyRepo for edge-loading
		ccs, err := client.CharacterCompetency.Query().
			Where(charactercompetency.HasCharacterWith(character.ID(charID))).
			WithCategory(func(q *db.CompetencyCategoryQuery) {
				q.WithThresholds()
			}).
			All(c.Request.Context())
		if err != nil {
			dblog.Error("failed to list character competencies", err, slog.String("service", "competencies"), slog.String("character_id", strconv.Itoa(charID)))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		type competencyView struct {
			CategoryID        string  `json:"category_id"`
			CategoryName      string  `json:"category_name"`
			Xp                int     `json:"xp"`
			Level             int     `json:"level"`
			XpMultiplier      float64 `json:"xp_multiplier"`
			DamageMultiplier  float64 `json:"damage_multiplier"`
			DefenseMultiplier float64 `json:"defense_multiplier"`
		}
		result := make([]competencyView, len(ccs))
		for i, cc := range ccs {
			cv := competencyView{Xp: cc.Xp, Level: cc.Level}
			cat := cc.Edges.Category
			if cat != nil {
				cv.CategoryID = cat.ID
				cv.CategoryName = cat.Name
				cv.XpMultiplier = cat.XpMultiplier
				if cc.Level > 0 {
					for _, t := range cat.Edges.Thresholds {
						if t.Level == cc.Level {
							cv.DamageMultiplier = t.DamageMultiplier
							cv.DefenseMultiplier = t.DefenseMultiplier
							break
						}
					}
				}
			}
			result[i] = cv
		}
		c.JSON(http.StatusOK, result)
	}
}

func compCatToJSON(cat *db.CompetencyCategory) gin.H {
	thresholds := make([]gin.H, len(cat.Edges.Thresholds))
	for j, t := range cat.Edges.Thresholds {
		thresholds[j] = gin.H{
			"level":               t.Level,
			"xp_required":        t.XpRequired,
			"damage_multiplier":  t.DamageMultiplier,
			"defense_multiplier": t.DefenseMultiplier,
		}
	}
	return gin.H{
		"id":             cat.ID,
		"name":           cat.Name,
		"xp_multiplier":  cat.XpMultiplier,
		"thresholds":     thresholds,
	}
}