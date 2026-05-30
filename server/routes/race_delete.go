package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/repository"
)

// deleteRace deletes a race by ID.
func deleteRace(repos *repository.Container) gin.HandlerFunc {
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

		// Use the race's world ID for the character count check
		count, err := repos.Race.CountCharactersByRaceName(c.Request.Context(), r.Name, r.WorldID)
		if err == nil && count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "cannot delete: race is in use by characters"})
			return
		}

		if err := repos.Race.Delete(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "race not found"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
