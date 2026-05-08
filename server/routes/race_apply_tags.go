package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/db/charactertag"
	"herbst-server/db/race"
)

// applyRaceTags syncs a race's tags to all characters of that race.
func applyRaceTags(client *db.Client) gin.HandlerFunc {
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

		tagNames := make([]string, len(r.Edges.Tags))
		for i, t := range r.Edges.Tags {
			tagNames[i] = t.Name
		}

		characters, err := client.Character.Query().
			Where(character.RaceEQ(r.Name)).
			All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		updated := 0
		for _, ch := range characters {
			_, err := client.CharacterTag.Delete().
				Where(
					charactertag.HasCharacterWith(character.ID(ch.ID)),
					charactertag.SourceEQ("race"),
				).
				Exec(c.Request.Context())
			if err != nil {
				continue
			}

			for _, tagName := range tagNames {
				_, err := client.CharacterTag.Create().
					SetTag(tagName).
					SetSource("race").
					SetCharacter(ch).
					Save(c.Request.Context())
				if err != nil {
					continue
				}
			}
			updated++
		}

		c.JSON(http.StatusOK, gin.H{
			"race":              r.Name,
			"characters_updated": updated,
			"tags_applied":      tagNames,
		})
	}
}