package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/channelconfig"
	"herbst-server/middleware"
)

// RegisterChannelRoutes registers REST endpoints for global channel configurations.
func RegisterChannelRoutes(r *gin.Engine, client *db.Client) {
	group := r.Group("/api")
	group.Use(middleware.AuthMiddleware())
	group.Use(middleware.AdminMiddleware())
	{
		group.GET("/channels", listChannels(client))
		group.GET("/channels/:name", getChannel(client))
		group.POST("/channels", createChannel(client))
		group.PUT("/channels/:name", updateChannel(client))
		group.DELETE("/channels/:name", deleteChannel(client))
	}
}

// ─── Channel Config CRUD ─────────────────────────────────────────────────────

func listChannels(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := client.ChannelConfig.Query()
		if search := c.Query("search"); search != "" {
			query = query.Where(channelconfig.NameContains(search))
		}
		configs, err := query.All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result := make([]gin.H, len(configs))
		for i, cfg := range configs {
			result[i] = channelConfigToJSON(cfg)
		}
		c.JSON(http.StatusOK, result)
	}
}

func getChannel(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		cfg, err := client.ChannelConfig.Query().
			Where(channelconfig.NameEQ(name)).
			Only(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "channel configuration not found"})
			return
		}
		c.JSON(http.StatusOK, channelConfigToJSON(cfg))
	}
}

func createChannel(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name            string `json:"name" binding:"required"`
			Description     string `json:"description"`
			Color           string `json:"color"`
			DefaultEnabled  bool   `json:"default_enabled"`
			CooldownSeconds int    `json:"cooldown_seconds"`
			AdminOnly       bool   `json:"admin_only"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		created, err := client.ChannelConfig.Create().
			SetName(req.Name).
			SetDescription(req.Description).
			SetColor(req.Color).
			SetDefaultEnabled(req.DefaultEnabled).
			SetCooldownSeconds(req.CooldownSeconds).
			SetAdminOnly(req.AdminOnly).
			Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, channelConfigToJSON(created))
	}
}

func updateChannel(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		var req struct {
			Description     *string `json:"description"`
			Color           *string `json:"color"`
			DefaultEnabled  *bool   `json:"default_enabled"`
			CooldownSeconds *int    `json:"cooldown_seconds"`
			AdminOnly       *bool   `json:"admin_only"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		builder := client.ChannelConfig.Update().
			Where(channelconfig.NameEQ(name))
		if req.Description != nil {
			builder.SetDescription(*req.Description)
		}
		if req.Color != nil {
			builder.SetColor(*req.Color)
		}
		if req.DefaultEnabled != nil {
			builder.SetDefaultEnabled(*req.DefaultEnabled)
		}
		if req.CooldownSeconds != nil {
			builder.SetCooldownSeconds(*req.CooldownSeconds)
		}
		if req.AdminOnly != nil {
			builder.SetAdminOnly(*req.AdminOnly)
		}
		count, err := builder.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if count == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "channel configuration not found"})
			return
		}

		// Fetch the updated object to return it as requested by the frontend
		updated, err := client.ChannelConfig.Query().
			Where(channelconfig.NameEQ(name)).
			Only(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, channelConfigToJSON(updated))
	}
}

func deleteChannel(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		_, err := client.ChannelConfig.Delete().
			Where(channelconfig.NameEQ(name)).
			Exec(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNoContent, nil)
	}
}

// ─── JSON Helper ──────────────────────────────────────────────────────────────

func channelConfigToJSON(cfg *db.ChannelConfig) gin.H {
	return gin.H{
		"name":             cfg.Name,
		"description":      cfg.Description,
		"color":            cfg.Color,
		"default_enabled":  cfg.DefaultEnabled,
		"cooldown_seconds": cfg.CooldownSeconds,
		"admin_only":       cfg.AdminOnly,
	}
}
