package routes

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/dblog"
	"herbst-server/middleware"
	"herbst-server/repository"
)

// RegisterPlayableGenderRoutes registers public endpoints for playable genders.
func RegisterPlayableGenderRoutes(r *gin.Engine, repos *repository.Container) {
	// Public endpoint for playable genders - no auth required
	r.GET("/playable-genders", middleware.WorldIDRequiredMiddleware(), listPlayableGendersHandler(repos))
}

// listPlayableGendersHandler returns all playable genders for the specified world.
func listPlayableGendersHandler(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		queryID := c.Query("world_id")
		// Empty / "default" / non-numeric values are treated as world 1 (dev default).
		// "default" is the UI sentinel for an unconfigured world context.
		if queryID == "" || queryID == "default" {
			queryID = "1"
		}
		// Check if queryID is a numeric ID or a world name
		var worldID string
		if _, err := strconv.Atoi(queryID); err == nil {
			worldID = queryID
		} else {
			// Look up world by name to get the numeric ID
			world, err := repos.World.GetByName(c.Request.Context(), queryID)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "world not found"})
				return
			}
			worldID = strconv.Itoa(world.ID)
		}
		genders, err := repos.Gender.ListPlayable(c.Request.Context(), worldID)
		if err != nil {
			dblog.Error("failed to list playable genders", err, slog.String("service", "genders"), slog.String("world_id", worldID))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"genders": genders})
	}
}
