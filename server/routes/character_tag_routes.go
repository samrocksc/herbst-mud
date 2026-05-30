package routes

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/dblog"
	"herbst-server/repository"
)

// RegisterCharacterTagRoutes registers character tag endpoints.
func RegisterCharacterTagRoutes(r *gin.Engine, repos *repository.Container) {
	r.GET("/characters/:id/tags", getCharacterTags(repos))
	r.POST("/characters/:id/tags", addCharacterTag(repos))
	r.DELETE("/characters/:id/tags/:tagId", removeCharacterTag(repos))
}

// getCharacterTags returns all tags for a character.
func getCharacterTags(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			slog.Warn("bad request: invalid character id", slog.String("service", "characters"), slog.String("character_id", c.Param("id")))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid character id"})
			return
		}

		tags, err := repos.CharacterTag.ListByCharacter(c.Request.Context(), id)
		if err != nil {
			dblog.Error("failed to list character tags", err, slog.String("service", "characters"), slog.Int("character_id", id))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		result := make([]gin.H, len(tags))
		for i, t := range tags {
			result[i] = gin.H{
				"id":       t.ID,
				"tag":      t.Tag,
				"source":   t.Source,
				"earned_at": t.EarnedAt,
			}
		}
		c.JSON(http.StatusOK, gin.H{"tags": result})
	}
}

// addCharacterTag adds a tag to a character.
func addCharacterTag(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			slog.Warn("bad request: invalid character id", slog.String("service", "characters"), slog.String("character_id", c.Param("id")))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid character id"})
			return
		}

		var req struct {
			Tag    string `json:"tag" binding:"required"`
			Source string `json:"source"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			slog.Warn("bad request: invalid add tag request", slog.String("service", "characters"), slog.Int("character_id", id), slog.String("error", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if req.Source == "" {
			req.Source = "admin"
		}

		ct, err := repos.CharacterTag.Create(c.Request.Context(), id, req.Tag, req.Source)
		if err != nil {
			dblog.Error("failed to add character tag", err, slog.String("service", "characters"), slog.Int("character_id", id), slog.String("tag", req.Tag))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		slog.Info("character tag added", slog.String("service", "characters"), slog.Int("character_id", id), slog.String("tag", ct.Tag), slog.String("source", ct.Source))
		c.JSON(http.StatusCreated, gin.H{
			"id":     ct.ID,
			"tag":    ct.Tag,
			"source": ct.Source,
		})
	}
}

// removeCharacterTag removes a tag from a character.
func removeCharacterTag(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		tagID, err := strconv.Atoi(c.Param("tagId"))
		if err != nil {
			slog.Warn("bad request: invalid tag id", slog.String("service", "characters"), slog.String("tag_id", c.Param("tagId")))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tag id"})
			return
		}

		if err := repos.CharacterTag.Delete(c.Request.Context(), tagID); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "tag not found"})
			return
		}
		slog.Info("character tag removed", slog.String("service", "characters"), slog.Int("tag_id", tagID))
		c.Status(http.StatusNoContent)
	}
}
