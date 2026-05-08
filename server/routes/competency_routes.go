package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/db/charactercompetency"
	"herbst-server/db/competencycategory"
	"herbst-server/db/competencylevelthreshold"
	"herbst-server/middleware"
)

func RegisterCompetencyRoutes(r *gin.Engine, client *db.Client) {
	competencies := r.Group("/api")
	competencies.Use(middleware.AuthMiddleware())
	competencies.Use(middleware.AdminMiddleware())
	{
		competencies.GET("/competency-categories", listCompetencyCategories(client))
		competencies.POST("/competency-categories", createCompetencyCategory(client))
		competencies.PUT("/competency-categories/:id", updateCompetencyCategory(client))
		competencies.DELETE("/competency-categories/:id", deleteCompetencyCategory(client))
		competencies.GET("/characters/:id/competencies", listCharacterCompetencies(client))
	}
}

func listCompetencyCategories(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		cats, err := client.CompetencyCategory.Query().
			WithThresholds().
			All(c.Request.Context())
		if err != nil {
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

func createCompetencyCategory(client *db.Client) gin.HandlerFunc {
	var req struct {
		ID           string            `json:"id" binding:"required"`
		Name         string            `json:"name" binding:"required"`
		XpMultiplier float64          `json:"xp_multiplier"`
		Thresholds   []thresholdInput  `json:"thresholds"`
	}
	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		cat, err := client.CompetencyCategory.Create().
			SetID(req.ID).
			SetName(req.Name).
			SetXpMultiplier(req.XpMultiplier).
			Save(c.Request.Context())
		if err != nil {
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
				SetCategory(cat).
				Save(c.Request.Context())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		cat, err = client.CompetencyCategory.Query().
			Where(competencycategory.ID(cat.ID)).
			WithThresholds().
			Only(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, compCatToJSON(cat))
	}
}

func updateCompetencyCategory(client *db.Client) gin.HandlerFunc {
	var req struct {
		Name         string            `json:"name"`
		XpMultiplier *float64          `json:"xp_multiplier"`
		Thresholds   []thresholdInput  `json:"thresholds"`
	}
	return func(c *gin.Context) {
		id := c.Param("id")
		_, err := client.CompetencyCategory.Query().
			Where(competencycategory.ID(id)).
			Only(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
			return
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		mutation := client.CompetencyCategory.UpdateOneID(id).
			SetName(req.Name)
		if req.XpMultiplier != nil {
			mutation = mutation.SetXpMultiplier(*req.XpMultiplier)
		}
		if req.Thresholds != nil {
			_, err := client.CompetencyLevelThreshold.Delete().
				Where(competencylevelthreshold.HasCategoryWith(competencycategory.ID(id))).
				Exec(c.Request.Context())
			if err != nil {
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
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}
		}
		cat, err := mutation.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		cat, err = client.CompetencyCategory.Query().
			Where(competencycategory.ID(cat.ID)).
			WithThresholds().
			Only(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, compCatToJSON(cat))
	}
}

func deleteCompetencyCategory(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		count, err := client.CharacterCompetency.Query().
			Where(charactercompetency.HasCategoryWith(competencycategory.ID(id))).
			Count(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Cannot delete: characters have XP in this category"})
			return
		}
		if err := client.CompetencyCategory.DeleteOneID(id).Exec(c.Request.Context()); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func listCharacterCompetencies(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		charIDStr := c.Param("id")
		charID, err := strconv.Atoi(charIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid character id"})
			return
		}
		if _, err := client.Character.Get(c.Request.Context(), charID); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "character not found"})
			return
		}
		ccs, err := client.CharacterCompetency.Query().
			Where(charactercompetency.HasCharacterWith(character.ID(charID))).
			WithCategory(func(q *db.CompetencyCategoryQuery) {
				q.WithThresholds()
			}).
			All(c.Request.Context())
		if err != nil {
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