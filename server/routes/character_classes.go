package routes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/faction"
	"herbst-server/db/factioncategory"
)

// ClassInfo holds public class data for the client.
type ClassInfo struct {
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	Description string   `json:"description"`
	Specialties []string `json:"specialties"`
	WorldID     string   `json:"world_id,omitempty"`
	Source      string   `json:"source"` // always "db" now
}

// listClasses handles GET /classes — returns world-scoped classes from DB.
func listClasses(c *gin.Context) {
	worldID := c.Query("world_id")
	ctx := context.Background()
	client := clientFromFactionRoutes(c)

	// World-scoped query: factions whose category.name="class" for the given world.
	if worldID != "" {
		cat, err := client.FactionCategory.Query().
			Where(
				factioncategory.Name("class"),
				factioncategory.WorldID(worldID),
			).
			WithFactions().
			Only(ctx)
		if err == nil && cat != nil && cat.Edges.Factions != nil && len(cat.Edges.Factions) > 0 {
			classes := make([]ClassInfo, 0, len(cat.Edges.Factions))
			for _, f := range cat.Edges.Factions {
				classes = append(classes, ClassInfo{
					Name:        f.Name,
					DisplayName: f.DisplayName,
					Description: f.Description,
					Specialties: specialtyIDsFromFaction(f),
					WorldID:     worldID,
					Source:      "db",
				})
			}
			c.JSON(http.StatusOK, gin.H{
				"classes": classes,
				"count":   len(classes),
			})
			return
		}
		// No class category for this world — return empty array.
		c.JSON(http.StatusOK, gin.H{
			"classes": []ClassInfo{},
			"count":   0,
		})
		return
	}

	// No world_id: return all class factions across all worlds.
	factions, err := client.Faction.Query().
		Where(faction.HasCategoryWith(factioncategory.Name("class"))).
		WithCategory().
		All(ctx)
	if err == nil && len(factions) > 0 {
		classes := make([]ClassInfo, 0, len(factions))
		for _, f := range factions {
			worldIDVal := ""
			if f.Edges.Category != nil {
				worldIDVal = f.Edges.Category.WorldID
			}
			classes = append(classes, ClassInfo{
				Name:        f.Name,
				DisplayName: f.DisplayName,
				Description: f.Description,
				Specialties: specialtyIDsFromFaction(f),
				WorldID:     worldIDVal,
				Source:      "db",
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"classes": classes,
			"count":   len(classes),
		})
		return
	}

	// No DB classes found at all — return empty array.
	c.JSON(http.StatusOK, gin.H{
		"classes": []ClassInfo{},
		"count":   0,
	})
}

// specialtyIDsFromFaction extracts specialty IDs from a faction's Specialties JSON field.
func specialtyIDsFromFaction(f *db.Faction) []string {
	var ids []string
	for _, s := range f.Specialties {
		ids = append(ids, s.ID)
	}
	return ids
}

// clientFromFactionRoutes extracts the *db.Client from the Gin context
// (stored there by RegisterCharacterRoutes).
func clientFromFactionRoutes(c *gin.Context) *db.Client {
	return c.MustGet("db_client").(*db.Client)
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]-32) + s[1:]
}

// isValidClassDB checks if a class name exists as a faction in a "class"
// faction category for the given world.
func isValidClassDB(ctx context.Context, client *db.Client, className, worldID string) bool {
	f, err := getClassFactionByName(ctx, client, className, worldID)
	return err == nil && f != nil
}

// getClassFactionByName queries the DB for a faction with the given name in a
// "class" faction category for the given world. Returns nil if not found.
func getClassFactionByName(ctx context.Context, client *db.Client, name, worldID string) (*db.Faction, error) {
	cat, err := client.FactionCategory.Query().
		Where(
			factioncategory.Name("class"),
			factioncategory.WorldID(worldID),
		).
		WithFactions().
		Only(ctx)
	if err != nil {
		return nil, err
	}
	for _, f := range cat.Edges.Factions {
		if f.Name == name {
			return f, nil
		}
	}
	return nil, fmt.Errorf("class faction %q not found in world %q", name, worldID)
}