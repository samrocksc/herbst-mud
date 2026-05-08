package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/db/charactercompetency"
	"herbst-server/middleware"
)

func RegisterCompetencyRoutes(r *gin.Engine, client *db.Client) {
	competencies := r.Group("/api")
	competencies.Use(middleware.AuthMiddleware())
	competencies.Use(middleware.AdminMiddleware())
	{
		competencies.GET("/competency-categories", listCompetencyCategories(client))
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

		type thresholdView struct {
			Level              int     `json:"level"`
			XpRequired         int     `json:"xp_required"`
			DamageMultiplier   float64 `json:"damage_multiplier"`
			DefenseMultiplier  float64 `json:"defense_multiplier"`
		}
		type categoryView struct {
			ID           string          `json:"id"`
			Name         string          `json:"name"`
			XpMultiplier float64         `json:"xp_multiplier"`
			Thresholds   []thresholdView `json:"thresholds"`
		}
		result := make([]categoryView, len(cats))
		for i, cat := range cats {
			thresholds := make([]thresholdView, len(cat.Edges.Thresholds))
			for j, t := range cat.Edges.Thresholds {
				thresholds[j] = thresholdView{
					Level:             t.Level,
					XpRequired:        t.XpRequired,
					DamageMultiplier:  t.DamageMultiplier,
					DefenseMultiplier: t.DefenseMultiplier,
				}
			}
			result[i] = categoryView{
				ID:           cat.ID,
				Name:         cat.Name,
				XpMultiplier: cat.XpMultiplier,
				Thresholds:   thresholds,
			}
		}
		c.JSON(http.StatusOK, result)
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
			cat := cc.Edges.Category
			cv := competencyView{
				Xp:           cc.Xp,
				Level:        cc.Level,
			}
			if cat != nil {
				cv.CategoryID = cat.ID
				cv.CategoryName = cat.Name
				cv.XpMultiplier = cat.XpMultiplier
				// Find threshold for current level
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