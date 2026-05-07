package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/db/race"
)

// deleteRace deletes a race by ID.
func deleteRace(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid race id"})
			return
		}

		raceName := client.Race.Query().Where(race.ID(id)).OnlyX(c.Request.Context()).Name
		count, err := client.Character.Query().
			Where(character.RaceEQ(raceName)).
			Count(c.Request.Context())
		if err == nil && count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "cannot delete: race is in use by characters"})
			return
		}

		err = client.Race.DeleteOneID(id).Exec(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "race not found"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}