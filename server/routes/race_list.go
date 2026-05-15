package routes

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/middleware"
	"herbst-server/repository"
)

// RegisterRaceRoutes registers REST endpoints for races.
func RegisterRaceRoutes(r *gin.Engine, repos *repository.Container, client *db.Client) {
	races := r.Group("/api/races")
	races.Use(middleware.AuthMiddleware(nil))
	races.Use(middleware.AdminMiddleware())
	{
		races.GET("", listRaces(repos))
		races.GET("/:id", getRace(repos))
		races.POST("", createRace(repos, client))
		races.PUT("/:id", updateRace(repos, client))
		races.DELETE("/:id", deleteRace(repos))
		races.POST("/:id/apply-tags", applyRaceTags(repos, client))
	}
}

// listRaces returns all races.
func listRaces(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		races, err := repos.Race.ListWithTags(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if search := c.Query("search"); search != "" {
			s := strings.ToLower(search)
			filtered := make([]*db.Race, 0, len(races))
			for _, r := range races {
				if strings.Contains(strings.ToLower(r.Name), s) {
					filtered = append(filtered, r)
				}
			}
			races = filtered
		}
		views := make([]raceView, len(races))
		for i, r := range races {
			views[i] = raceToView(r)
		}
		c.JSON(http.StatusOK, gin.H{"races": views})
	}
}

// getRace returns a single race by ID.
func getRace(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid race id"})
			return
		}
		r, err := repos.Race.GetWithTags(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "race not found"})
			return
		}
		c.JSON(http.StatusOK, raceToView(r))
	}
}