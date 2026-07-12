package routes

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/dblog"
	"herbst-server/middleware"
	"herbst-server/repository"
	"log/slog"
)

type TargetSearchResult struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
}

func RegisterTargetRoutes(r *gin.Engine, repos *repository.Container, client *db.Client) {
	targets := r.Group("/api")
	targets.Use(middleware.AuthMiddleware(nil))
	targets.Use(middleware.AdminMiddleware())
	targets.Use(middleware.WorldAccessMiddleware())

	targets.GET("/targets", func(c *gin.Context) {
		search := c.Query("search")
		worldID := c.Query("world_id")

		s := strings.ToLower(search)

		// Collect all results first, then filter or limit based on query
		allResults := make([]TargetSearchResult, 0)

		// Search effects (AbilityEffect)
		effects, err := repos.Effect.List(c.Request.Context())
		if err == nil {
			for _, e := range effects {
				allResults = append(allResults, TargetSearchResult{
					ID:   e.ID,
					Name: e.Name,
					Type: "effect",
				})
			}
		} else {
			dblog.Error("list effects for targets failed", err, slog.String("service", "targets"))
		}

		// Search dialog nodes - id is string, display text is npc_text
		dn, err := repos.DialogNode.List(c.Request.Context(), worldID)
		if err == nil {
			for _, d := range dn {
				allResults = append(allResults, TargetSearchResult{
					ID:   hashCodeToInt(d.ID),
					Name: d.NpcText,
					Type: "dialog_node",
				})
			}
		} else {
			dblog.Error("list dialog nodes for targets failed", err, slog.String("service", "targets"))
		}

		// Search crafting recipes
		recipes, err := repos.CraftingRecipe.List(c.Request.Context(), worldID, "")
		if err == nil {
			for _, r := range recipes {
				allResults = append(allResults, TargetSearchResult{
					ID:   r.ID,
					Name: r.Name,
					Type: "recipe",
				})
			}
		} else {
			dblog.Error("list recipes for targets failed", err, slog.String("service", "targets"))
		}

		// Filter by search query if provided
		var results []TargetSearchResult
		if search == "" {
			// Return up to 4 initial suggestions
			results = allResults
			if len(results) > 4 {
				results = results[:4]
			}
		} else {
			// Filter by search query
			results = make([]TargetSearchResult, 0)
			for _, r := range allResults {
				if strings.Contains(strings.ToLower(r.Name), s) {
					results = append(results, r)
				}
			}
		}

		c.JSON(http.StatusOK, results)
	})

	targets.GET("/targets/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid target id"})
			return
		}
		worldID := c.Query("world_id")

		// Search across all target types to find the one with this ID

		// Check effects first
		effects, _ := repos.Effect.List(c.Request.Context())
		for _, e := range effects {
			if e.ID == id {
				c.JSON(http.StatusOK, gin.H{"id": e.ID, "name": e.Name})
				return
			}
		}

		// Check recipes
		recipes, _ := repos.CraftingRecipe.List(c.Request.Context(), worldID, "")
		for _, r := range recipes {
			if r.ID == id {
				c.JSON(http.StatusOK, gin.H{"id": r.ID, "name": r.Name})
				return
			}
		}

		// Check dialog nodes (hash match)
		dn, _ := repos.DialogNode.List(c.Request.Context(), worldID)
		for _, d := range dn {
			if hashCodeToInt(d.ID) == id {
				c.JSON(http.StatusOK, gin.H{"id": id, "name": d.NpcText})
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "target not found"})
	})
}

// hashCodeToInt converts a string to an int for use as a pseudo-ID
func hashCodeToInt(s string) int {
	hash := 0
	for i := 0; i < len(s); i++ {
		hash = 31*hash + int(s[i])
	}
	// Make sure it's positive and not zero (0 is reserved for "no value")
	if hash <= 0 {
		hash = 1
	}
	return hash
}
