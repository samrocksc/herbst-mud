package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/tag"
	"herbst-server/middleware"
)

// RegisterTagRoutes registers REST endpoints for tags.
func RegisterTagRoutes(r *gin.Engine, client *db.Client) {
	tags := r.Group("/api/tags")
	tags.Use(middleware.AuthMiddleware())
	tags.Use(middleware.AdminMiddleware())
	{
		tags.GET("", listTags(client))
		tags.POST("", createTag(client))
		tags.PUT("/:id", updateTag(client))
		tags.DELETE("/:id", deleteTag(client))
	}
}

// tagView is the JSON shape returned by the API.
type tagView struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color,omitempty"`
}

// listTags returns all tags.
func listTags(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		tags, err := client.Tag.Query().Order(tag.ByName()).All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query tags"})
			return
		}
		views := make([]tagView, len(tags))
		for i, t := range tags {
			views[i] = tagView{ID: t.ID, Name: t.Name, Color: t.Color}
		}
		c.JSON(http.StatusOK, gin.H{"tags": views})
	}
}

// createTag creates a new tag.
func createTag(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Name  string `json:"name" binding:"required"`
			Color string `json:"color"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
			return
		}

		t, err := client.Tag.Create().
			SetName(input.Name).
			SetColor(input.Color).
			Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create tag"})
			return
		}
		c.JSON(http.StatusCreated, tagView{ID: t.ID, Name: t.Name, Color: t.Color})
	}
}

// updateTag updates an existing tag.
func updateTag(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		var input struct {
			Name  *string `json:"name"`
			Color *string `json:"color"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		// Build update mutation
		mut := client.Tag.UpdateOneID(id)
		if input.Name != nil {
			mut = mut.SetName(*input.Name)
		}
		if input.Color != nil {
			mut = mut.SetColor(*input.Color)
		}

		t, err := mut.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "tag not found"})
			return
		}

		c.JSON(http.StatusOK, tagView{ID: t.ID, Name: t.Name, Color: t.Color})
	}
}

// deleteTag deletes a tag by ID.
func deleteTag(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		err = client.Tag.DeleteOneID(id).Exec(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "tag not found"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
