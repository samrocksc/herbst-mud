package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/ability"
	"herbst-server/db/character"
	"herbst-server/db/charactertag"
	"herbst-server/db/faction"
	"herbst-server/db/factionrequiredtag"
	"herbst-server/middleware"
	"herbst-server/repository"
)

// RegisterTagRoutes registers REST endpoints for tags.
func RegisterTagRoutes(r *gin.Engine, repos *repository.Container, client *db.Client) {
	tags := r.Group("/api/tags")
	tags.Use(middleware.AuthMiddleware())
	tags.Use(middleware.AdminMiddleware())
	{
		tags.GET("", listTags(repos))
		tags.POST("", createTag(repos))
		tags.PUT("/:id", updateTag(repos))
		tags.DELETE("/:id", deleteTag(repos))
		// TODO: migrate tagUsages to repo methods for cross-entity queries
		tags.GET("/:id/usages", tagUsages(repos, client))
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
	Abilities   []tagUsageView `json:"abilities"`
	Factions    []tagUsageView `json:"factions"`
	Characters  []tagUsageView `json:"characters"`
}

// listTags returns all tags.
func listTags(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		tags, err := repos.Tag.List(c.Request.Context())
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
func createTag(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Name  string `json:"name" binding:"required"`
			Color string `json:"color"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
			return
		}

		t, err := repos.Tag.Create(c.Request.Context(), repository.CreateTagInput{
			Name:  input.Name,
			Color: input.Color,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create tag"})
			return
		}
		c.JSON(http.StatusCreated, tagView{ID: t.ID, Name: t.Name, Color: t.Color})
	}
}

// updateTag updates an existing tag.
func updateTag(repos *repository.Container) gin.HandlerFunc {
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

		t, err := repos.Tag.Update(c.Request.Context(), id, repository.TagUpdates{
			Name:  input.Name,
			Color: input.Color,
		})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "tag not found"})
			return
		}

		c.JSON(http.StatusOK, tagView{ID: t.ID, Name: t.Name, Color: t.Color})
	}
}

// deleteTag deletes a tag by ID.
func deleteTag(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		if err := repos.Tag.Delete(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "tag not found"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

// tagUsages returns every entity that references the given tag name.
// TODO: Add usage-tracking methods to repos — currently uses client directly for cross-entity queries.
func tagUsages(repos *repository.Container, client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		tagEntity, err := repos.Tag.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "tag not found"})
			return
		}

		report := tagUsageReport{
			TagName: tagEntity.Name,
		}

		// --- Abilities (required_tag field) ---
		abilities, err := client.Ability.Query().
			Where(ability.RequiredTag(tagEntity.Name)).
			All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query abilities"})
			return
		}
		report.Abilities = make([]tagUsageView, len(abilities))
		for i, s := range abilities {
			report.Abilities[i] = tagUsageView{ID: s.ID, Name: s.Name, Type: "ability"}
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