package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/gameconfig"
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
func RegisterGameConfigRoutes(router *gin.RouterGroup, client *db.Client) {
	// List all configs — auth required
	router.GET("/game-configs", func(c *gin.Context) {
		configs, err := client.GameConfig.Query().All(c.Request.Context())
		if err != nil {
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
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		created, err := client.GameConfig.Create().
			SetKey(req.Key).
			SetValue(req.Value).
			Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, GameConfigResponse{ID: created.ID, Key: created.Key, Value: created.Value})
	})

	// Get a single config by key — public (game server reads without auth)
	router.GET("/game-configs/:key", func(c *gin.Context) {
		key := c.Param("key")
		cfg, err := client.GameConfig.Query().Where(gameconfig.Key(key)).Only(c.Request.Context())
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
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		cfg, err := client.GameConfig.Query().Where(gameconfig.Key(key)).Only(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "config key not found"})
			return
		}
		updated, err := client.GameConfig.UpdateOne(cfg).
			SetValue(req.Value).
			Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, GameConfigResponse{ID: updated.ID, Key: updated.Key, Value: updated.Value})
	})

	// Delete a config — auth required
	router.DELETE("/game-configs/:key", func(c *gin.Context) {
		key := c.Param("key")
		cfg, err := client.GameConfig.Query().Where(gameconfig.Key(key)).Only(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "config key not found"})
			return
		}
		if err := client.GameConfig.DeleteOne(cfg).Exec(c.Request.Context()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"deleted": key})
	})
}
