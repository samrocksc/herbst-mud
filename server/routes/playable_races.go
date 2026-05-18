package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/repository"
)

// RegisterPlayableRaceRoutes registers public endpoints for playable races.
func RegisterPlayableRaceRoutes(r *gin.Engine, repos *repository.Container) {
	// Public endpoint for playable races - no auth required
	r.GET("/playable-races", listPlayableRacesHandler(repos))
}

// listPlayableRacesHandler returns all playable races.
func listPlayableRacesHandler(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		races, err := repos.Race.ListPlayable(c.Request.Context())
		if err != nil {
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
}

// playableRaceToView converts a Race ent model to a playableRaceView.
func playableRaceToView(r *repository.Race) playableRaceView {
	return playableRaceView{
		Name:        r.Name,
		DisplayName: r.DisplayName,
	}
}