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
			All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		type categoryView struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			XpMultiplier float64 `json:"xp_multiplier"`
		}
		result := make([]categoryView, len(cats))
		for i, cat := range cats {
			result[i] = categoryView{ID: cat.ID, Name: cat.Name, XpMultiplier: cat.XpMultiplier}
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

		// Verify character exists
		if _, err := client.Character.Get(c.Request.Context(), charID); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "character not found"})
			return
		}

		ccs, err := client.CharacterCompetency.Query().
			Where(charactercompetency.HasCharacterWith(character.ID(charID))).
			All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		type competencyView struct {
			CategoryID    string  `json:"category_id"`
			CategoryName  string  `json:"category_name"`
			Xp            int     `json:"xp"`
			Level         int     `json:"level"`
			XpMultiplier  float64 `json:"xp_multiplier"`
		}
		result := []competencyView{}
		for _, cc := range ccs {
			cat, err := cc.QueryCategory().Only(c.Request.Context())
			if err != nil {
				continue
			}
			result = append(result, competencyView{
				CategoryID:   cat.ID,
				CategoryName: cat.Name,
				Xp:           cc.Xp,
				Level:        cc.Level,
				XpMultiplier: cat.XpMultiplier,
			})
		}
		c.JSON(http.StatusOK, result)
	}
}
