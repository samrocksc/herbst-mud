package routes

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/dblog"
	"herbst-server/worldexport"
)

// RegisterGameExportRoutes registers export/import routes.
func RegisterGameExportRoutes(router *gin.Engine, client *db.Client) {
	router.GET("/admin/export/worlds", listWorlds(client))
	router.GET("/admin/export", exportWorld(client))
	router.POST("/admin/import", importWorld(client))
	router.POST("/admin/import/validate", validateImport())
}

func listWorlds(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		worlds, err := client.World.Query().All(c.Request.Context())
		if err != nil {
			dblog.Error("failed to fetch worlds", err, slog.String("path", "/admin/export/worlds"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		out := make([]gin.H, len(worlds))
		for i, w := range worlds {
			out[i] = gin.H{"id": w.ID, "name": w.Name, "title": w.Title, "description": w.Description, "active": w.Active}
		}
		c.JSON(http.StatusOK, gin.H{"worlds": out})
	}
}

func exportWorld(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		worldID := c.Query("world_id")
		if worldID == "" {
			worldID = "1"
		}
		snap, err := worldexport.ExportWorld(c.Request.Context(), client, worldID)
		if err != nil {
			dblog.Error("failed to export world", err, slog.String("world_id", worldID))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, snap)
	}
}

func importWorld(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Snapshot      *worldexport.WorldSnapshot `json:"snapshot" binding:"required"`
			NewWorldName  string                     `json:"new_world_name" binding:"required"`
			NewWorldTitle string                     `json:"new_world_title"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			slog.Warn("invalid import request", slog.String("error", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result, err := worldexport.ImportWorld(c.Request.Context(), client, req.Snapshot, req.NewWorldName, req.NewWorldTitle)
		if err != nil {
			dblog.Error("failed to import world", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "result": result})
	}
}

func validateImport() gin.HandlerFunc {
	return func(c *gin.Context) {
		var snap worldexport.WorldSnapshot
		if err := c.ShouldBindJSON(&snap); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"is_valid": false, "errors": []string{err.Error()}})
			return
		}
		if snap.Version != worldexport.CurrentVersion {
			c.JSON(http.StatusBadRequest, gin.H{"is_valid": false, "errors": []string{"unsupported version: " + snap.Version}})
			return
		}
		c.JSON(http.StatusOK, gin.H{"is_valid": true, "world": snap.World["name"], "rooms": len(snap.Rooms), "npcs": len(snap.NPCs)})
	}
}
