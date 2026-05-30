package routes

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/dblog"
	"herbst-server/db/character"
	"herbst-server/db/charactertag"
	"herbst-server/repository"
)

// applyRaceTags syncs a race's tags to all characters of that race.
// TODO: Migrate character tag operations to CharacterTagRepo
func applyRaceTags(repos *repository.Container, client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			slog.Warn("bad request", slog.String("service", "races"), slog.String("reason", "invalid race id"))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid race id"})
			return
		}

		r, err := repos.Race.GetWithTags(c.Request.Context(), id)
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
			dblog.Error("failed to list characters for race tag application", err, slog.String("service", "races"), slog.Int("race_id", id))
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

		slog.Info("race tags applied", slog.String("service", "races"), slog.String("race", r.Name), slog.Int("characters_updated", updated), slog.Any("tags", tagNames))
		c.JSON(http.StatusOK, gin.H{
			"race":               r.Name,
			"characters_updated": updated,
			"tags_applied":      tagNames,
		})
	}
}