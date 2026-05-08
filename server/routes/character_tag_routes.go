package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/db/charactertag"
)

// RegisterCharacterTagRoutes registers character tag endpoints.
func RegisterCharacterTagRoutes(r *gin.Engine, client *db.Client) {
	r.GET("/characters/:id/tags", getCharacterTags(client))
	r.POST("/characters/:id/tags", addCharacterTag(client))
	r.DELETE("/characters/:id/tags/:tagId", removeCharacterTag(client))
}

// getCharacterTags returns all tags for a character.
func getCharacterTags(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid character id"})
			return
		}

		tags, err := client.CharacterTag.Query().
			Where(charactertag.HasCharacterWith(character.ID(id))).
			All(c.Request.Context())
		if err != nil {
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
func addCharacterTag(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid character id"})
			return
		}

		var req struct {
			Tag    string `json:"tag" binding:"required"`
			Source string `json:"source"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if req.Source == "" {
			req.Source = "admin"
		}

		ct, err := client.CharacterTag.Create().
			SetTag(req.Tag).
			SetSource(req.Source).
			SetCharacterID(id).
			Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"id":       ct.ID,
			"tag":      ct.Tag,
			"source":   ct.Source,
		})
	}
}

// removeCharacterTag removes a tag from a character.
func removeCharacterTag(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		tagID, err := strconv.Atoi(c.Param("tagId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tag id"})
			return
		}

		err = client.CharacterTag.DeleteOneID(tagID).Exec(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "tag not found"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}