package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/db/charactertag"
	"herbst-server/db/faction"
	"herbst-server/db/factionrequiredtag"
	"herbst-server/db/skill"
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
		tags.GET("/:id/usages", tagUsages(client))
	}
}

// tagView is the JSON shape returned by the API.
type tagView struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color,omitempty"`
}

// tagUsageView is a single reference to an entity that uses a tag.
type tagUsageView struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// tagUsageReport is the JSON shape for /api/tags/:id/usages.
type tagUsageReport struct {
	TagName     string         `json:"tag_name"`
	TotalUsages int            `json:"total_usages"`
	Skills      []tagUsageView `json:"skills"`
	Factions    []tagUsageView `json:"factions"`
	Characters  []tagUsageView `json:"characters"`
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

// tagUsages returns every entity that references the given tag name.
func tagUsages(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		tagEntity, err := client.Tag.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "tag not found"})
			return
		}

		report := tagUsageReport{
			TagName: tagEntity.Name,
		}

		// --- Skills (required_tag field) ---
		skills, err := client.Skill.Query().
			Where(skill.RequiredTag(tagEntity.Name)).
			All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query skills"})
			return
		}
		report.Skills = make([]tagUsageView, len(skills))
		for i, s := range skills {
			report.Skills[i] = tagUsageView{ID: s.ID, Name: s.Name, Type: "skill"}
			report.TotalUsages++
		}

		// --- Factions (required_tags via faction_required_tag join) ---
		factionIDs, err := client.FactionRequiredTag.Query().
			Where(factionrequiredtag.RequiredTag(tagEntity.Name)).
			WithFaction(func(q *db.FactionQuery) {
				q.Select(faction.FieldID, faction.FieldName)
			}).
			All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query faction tags"})
			return
		}
		report.Factions = make([]tagUsageView, 0, len(factionIDs))
		for _, frt := range factionIDs {
			if frt.Edges.Faction != nil {
				report.Factions = append(report.Factions, tagUsageView{
					ID:   frt.Edges.Faction.ID,
					Name: frt.Edges.Faction.Name,
					Type: "faction",
				})
				report.TotalUsages++
			}
		}

		// --- Characters (character_tag join) ---
		charTags, err := client.CharacterTag.Query().
			Where(charactertag.Tag(tagEntity.Name)).
			WithCharacter(func(q *db.CharacterQuery) {
				q.Select(character.FieldID, character.FieldName)
			}).
			All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query character tags"})
			return
		}
		report.Characters = make([]tagUsageView, 0, len(charTags))
		for _, ct := range charTags {
			if ct.Edges.Character != nil {
				report.Characters = append(report.Characters, tagUsageView{
					ID:   ct.Edges.Character.ID,
					Name: ct.Edges.Character.Name,
					Type: "character",
				})
				report.TotalUsages++
			}
		}

		c.JSON(http.StatusOK, report)
	}
}
