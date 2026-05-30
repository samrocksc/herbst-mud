package routes

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/dblog"
	"herbst-server/repository"
)

// GameConfigResponse is the API shape for a single game config entry.
type GameConfigResponse struct {
	ID    int    `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

// RegisterGameConfigRoutes sets up game config management endpoints.
// GET /game-configs        — list all configs (auth required)
// POST /game-configs       — create a config (auth required)
// GET /game-configs/:key  — get a single config (public, used by game server)
// PUT /game-configs/:key  — update a config (auth required)
// DELETE /game-configs/:key — delete a config (auth required)
func RegisterGameConfigRoutes(router *gin.RouterGroup, repos *repository.Container) {
	// List all configs — auth required
	router.GET("/game-configs", func(c *gin.Context) {
		configs, err := repos.GameConfig.List(c.Request.Context())
		if err != nil {
			dblog.Error("Failed to list game configs", err, slog.String("service", "game"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		out := make([]GameConfigResponse, 0, len(configs))
		for _, cfg := range configs {
			out = append(out, GameConfigResponse{ID: cfg.ID, Key: cfg.Key, Value: cfg.Value})
		}
		c.JSON(http.StatusOK, out)
	})

	// Create a new config — auth required
	router.POST("/game-configs", func(c *gin.Context) {
		var req struct {
			Key   string `json:"key" binding:"required"`
			Value string `json:"value" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			slog.Warn("Invalid create game config request", "error", err, slog.String("service", "game"))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		created, err := repos.GameConfig.GetOrCreate(c.Request.Context(), req.Key, req.Value)
		if err != nil {
			dblog.Error("Failed to create game config", err, slog.String("service", "game"), slog.String("key", req.Key))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		slog.Info("Game config created", slog.String("service", "game"), slog.String("key", created.Key))
		c.JSON(http.StatusCreated, GameConfigResponse{ID: created.ID, Key: created.Key, Value: created.Value})
	})

	// Get a single config by key — public (game server reads without auth)
	router.GET("/game-configs/:key", func(c *gin.Context) {
		key := c.Param("key")
		cfg, err := repos.GameConfig.Get(c.Request.Context(), key)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "config key not found"})
			return
		}
		c.JSON(http.StatusOK, GameConfigResponse{ID: cfg.ID, Key: cfg.Key, Value: cfg.Value})
	})

	// Update a config — auth required
	router.PUT("/game-configs/:key", func(c *gin.Context) {
		key := c.Param("key")
		var req struct {
			Value string `json:"value" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			slog.Warn("Invalid update game config request", "error", err, slog.String("service", "game"))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		updated, err := repos.GameConfig.Set(c.Request.Context(), key, req.Value)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "config key not found"})
			return
		}
		slog.Info("Game config updated", slog.String("service", "game"), slog.String("key", updated.Key))
		c.JSON(http.StatusOK, GameConfigResponse{ID: updated.ID, Key: updated.Key, Value: updated.Value})
	})

	// Delete a config — auth required
	router.DELETE("/game-configs/:key", func(c *gin.Context) {
		key := c.Param("key")
		if err := repos.GameConfig.Delete(c.Request.Context(), key); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "config key not found"})
			return
		}
		slog.Info("Game config deleted", slog.String("service", "game"), slog.String("key", key))
		c.JSON(http.StatusOK, gin.H{"deleted": key})
	})
}
