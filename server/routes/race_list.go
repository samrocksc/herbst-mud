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
	races := r.Group("/api/races")
	races.Use(middleware.AuthMiddleware())
	races.Use(middleware.AdminMiddleware())
	{
		races.GET("", listRaces(client))
		races.GET("/:id", getRace(client))
		races.POST("", createRace(client))
		races.PUT("/:id", updateRace(client))
		races.DELETE("/:id", deleteRace(client))
		races.POST("/:id/apply-tags", applyRaceTags(client))
	}
}

// listRaces returns all races.
func listRaces(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		races, err := client.Race.Query().WithTags().Order(race.ByDisplayName()).All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		views := make([]raceView, len(races))
		for i, r := range races {
			views[i] = raceToView(r)
		}
		c.JSON(http.StatusOK, gin.H{"races": views})
	}
}

// getRace returns a single race by ID.
func getRace(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid race id"})
			return
		}
		r, err := client.Race.Query().Where(race.ID(id)).WithTags().Only(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "race not found"})
			return
		}
		c.JSON(http.StatusOK, raceToView(r))
	}
}