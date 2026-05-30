package routes

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/dblog"
	"herbst-server/repository"
)

// RegisterPlayableRaceRoutes registers public endpoints for playable races.
func RegisterPlayableRaceRoutes(r *gin.Engine, repos *repository.Container) {
	// Public endpoint for playable races - no auth required
	r.GET("/playable-races", listPlayableRacesHandler(repos))
}

// listPlayableRacesHandler returns all playable races for the specified world.
func listPlayableRacesHandler(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Default to world "1" for public API
		worldID := c.Query("world_id")
		if worldID == "" {
			worldID = "1"
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
