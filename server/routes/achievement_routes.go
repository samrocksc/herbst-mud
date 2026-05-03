package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/achievement"
	"herbst-server/middleware"
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
// Protected /api routes — all require JWT auth + admin check
func RegisterAchievementRoutes(r *gin.Engine, client *db.Client) {
	achievements := r.Group("/api")
	achievements.Use(middleware.AuthMiddleware())
	achievements.Use(middleware.AdminMiddleware())
	{
		achievements.GET("/achievements", listAchievements(client))
		achievements.POST("/achievements", createAchievement(client))
		achievements.GET("/achievements/:id", getAchievement(client))
		achievements.PUT("/achievements/:id", updateAchievement(client))
		achievements.DELETE("/achievements/:id", deleteAchievement(client))
	}
}

func listAchievements(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		achievements, err := client.Achievement.Query().
			Order(achievement.ByName()).
			All(c.Request.Context())
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

func createAchievement(client *db.Client) gin.HandlerFunc {
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

		mut := client.Achievement.Create().
			SetName(input.Name).
			SetDescription(input.Desc).
			SetIcon(input.Icon).
			SetCriteria(input.Criteria)

		if input.XPReward != nil {
			mut.SetXpReward(*input.XPReward)
		}

		a, err := mut.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, achievementToView(a))
	}
}

func getAchievement(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid achievement id"})
			return
		}

		a, err := client.Achievement.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "achievement not found"})
			return
		}

		c.JSON(http.StatusOK, achievementToView(a))
	}
}

func updateAchievement(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid achievement id"})
			return
		}

		a, err := client.Achievement.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "achievement not found"})
			return
		}

		var input achievementInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		mut := client.Achievement.UpdateOne(a)

		if input.Name != "" {
			mut.SetName(input.Name)
		}
		mut.SetDescription(input.Desc)
		mut.SetIcon(input.Icon)
		mut.SetCriteria(input.Criteria)

		if input.XPReward != nil {
			mut.SetXpReward(*input.XPReward)
		}

		updated, err := mut.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, achievementToView(updated))
	}
}

func deleteAchievement(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid achievement id"})
			return
		}

		err = client.Achievement.DeleteOneID(id).Exec(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "achievement not found"})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}