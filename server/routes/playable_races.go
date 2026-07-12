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

// RegisterPlayableRaceRoutes registers public endpoints for playable races.
func RegisterPlayableRaceRoutes(r *gin.Engine, repos *repository.Container) {
	// Public endpoint for playable races - no auth required
	r.GET("/playable-races", middleware.WorldIDRequiredMiddleware(), listPlayableRacesHandler(repos))
}

// listPlayableRacesHandler returns all playable races for the specified world.
func listPlayableRacesHandler(repos *repository.Container) gin.HandlerFunc {
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
		races, err := repos.Race.ListPlayable(c.Request.Context(), worldID)
		if err != nil {
			dblog.Error("failed to list playable races", err, slog.String("service", "races"), slog.String("world_id", worldID))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		views := make([]playableRaceView, len(races))
		for i, r := range races {
			views[i] = playableRaceToView(r)
		}

		c.JSON(http.StatusOK, gin.H{"races": views})
	}
}

// playableRaceView is the minimal JSON shape for playable races.
type playableRaceView struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	WorldID     string `json:"world_id,omitempty"`
}

// playableRaceToView converts a Race ent model to a playableRaceView.
func playableRaceToView(r *repository.Race) playableRaceView {
	return playableRaceView{
		Name:        r.Name,
		DisplayName: r.DisplayName,
		WorldID:     r.WorldID,
	}
}
