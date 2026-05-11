package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/middleware"
	"herbst-server/repository"
)

// achievementView is the JSON shape returned by the API.
type achievementView struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Desc      string `json:"description"`
	Icon      string `json:"icon"`
	XPReward  int    `json:"xp_reward"`
	Criteria  string `json:"criteria"`
}

// achievementInput is the request body for create and update.
type achievementInput struct {
	Name     string `json:"name"`
	Desc     string `json:"description"`
	Icon     string `json:"icon"`
	XPReward *int  `json:"xp_reward"`
	Criteria string `json:"criteria"`
}

func achievementToView(a *db.Achievement) achievementView {
	return achievementView{
		ID:       a.ID,
		Name:     a.Name,
		Desc:     a.Description,
		Icon:     a.Icon,
		XPReward: a.XpReward,
		Criteria: a.Criteria,
	}
}

// RegisterAchievementRoutes registers REST endpoints for achievements.
func RegisterAchievementRoutes(r *gin.Engine, repos *repository.Container) {
	achievements := r.Group("/api")
	achievements.Use(middleware.AuthMiddleware())
	achievements.Use(middleware.AdminMiddleware())
	{
		achievements.GET("/achievements", listAchievements(repos))
		achievements.POST("/achievements", createAchievement(repos))
		achievements.GET("/achievements/:id", getAchievement(repos))
		achievements.PUT("/achievements/:id", updateAchievement(repos))
		achievements.DELETE("/achievements/:id", deleteAchievement(repos))
	}
}

func listAchievements(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		achievements, err := repos.Achievement.List(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result := make([]achievementView, len(achievements))
		for i, a := range achievements {
			result[i] = achievementToView(a)
		}
		c.JSON(http.StatusOK, gin.H{"achievements": result})
	}
}

func createAchievement(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input achievementInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if input.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
			return
		}

		xpReward := 0
		if input.XPReward != nil {
			xpReward = *input.XPReward
		}

		a, err := repos.Achievement.Create(c.Request.Context(), repository.CreateAchievementInput{
			Name:        input.Name,
			Description: input.Desc,
			Icon:        input.Icon,
			XPReward:    xpReward,
			Criteria:    input.Criteria,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, achievementToView(a))
	}
}

func getAchievement(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid achievement id"})
			return
		}

		a, err := repos.Achievement.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "achievement not found"})
			return
		}

		c.JSON(http.StatusOK, achievementToView(a))
	}
}

func updateAchievement(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid achievement id"})
			return
		}

		_, err = repos.Achievement.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "achievement not found"})
			return
		}

		var input achievementInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updates := repository.AchievementUpdates{
			Description: &input.Desc,
			Icon:        &input.Icon,
			Criteria:    &input.Criteria,
		}
		if input.Name != "" {
			updates.Name = &input.Name
		}
		if input.XPReward != nil {
			updates.XPReward = input.XPReward
		}

		updated, err := repos.Achievement.Update(c.Request.Context(), id, updates)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, achievementToView(updated))
	}
}

func deleteAchievement(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid achievement id"})
			return
		}

		if err := repos.Achievement.Delete(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "achievement not found"})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}