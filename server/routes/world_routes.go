package routes

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/content"
	"herbst-server/dblog"
	"herbst-server/middleware"
	"herbst-server/repository"
)

// RegisterWorldCRUDRoutes registers DB-backed world CRUD endpoints (independent of worldManager)
func RegisterWorldCRUDRoutes(router *gin.Engine, repos *repository.Container) {
	// Protected /api routes — all require JWT auth + admin check
	worlds := router.Group("/api")
	worlds.Use(middleware.AuthMiddleware(nil))
	worlds.Use(middleware.AdminMiddleware())
	{
		// POST /api/worlds - CreateWorldHandler
		worlds.POST("/worlds", func(c *gin.Context) {
			var input repository.CreateWorldInput
			if err := c.ShouldBindJSON(&input); err != nil {
				slog.Warn("invalid create world request", slog.String("error", err.Error()), slog.String("service", "worlds"))
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			world, err := repos.World.Create(c.Request.Context(), input)
			if err != nil {
				dblog.Error("failed to create world", err, slog.String("service", "worlds"))
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			slog.Info("world created", slog.Int("world_id", world.ID), slog.String("service", "worlds"))
			c.JSON(http.StatusCreated, world)
		})

		// GET /api/worlds/db - ListWorldsHandler (DB-backed)
		worlds.GET("/worlds/db", func(c *gin.Context) {
			worlds, err := repos.World.List(c.Request.Context())
			if err != nil {
				dblog.Error("failed to list worlds", err, slog.String("service", "worlds"))
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"worlds": worlds, "count": len(worlds)})
		})

		// GET /api/worlds/:id - GetWorldHandler (DB-backed)
		worlds.GET("/worlds/:id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				slog.Warn("invalid world id", slog.String("error", err.Error()), slog.String("service", "worlds"))
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid world ID"})
				return
			}
			world, err := repos.World.Get(c.Request.Context(), id)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "World not found"})
				return
			}
			c.JSON(http.StatusOK, world)
		})

		// PUT /api/worlds/:id - UpdateWorldHandler
		worlds.PUT("/worlds/:id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				slog.Warn("invalid world id", slog.String("error", err.Error()), slog.String("service", "worlds"))
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid world ID"})
				return
			}
			var updates repository.WorldUpdates
			if err := c.ShouldBindJSON(&updates); err != nil {
				slog.Warn("invalid update world request", slog.String("error", err.Error()), slog.String("service", "worlds"))
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			world, err := repos.World.Update(c.Request.Context(), id, updates)
			if err != nil {
				dblog.Error("failed to update world", err, slog.String("service", "worlds"), slog.Int("world_id", id))
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			slog.Info("world updated", slog.Int("world_id", world.ID), slog.String("service", "worlds"))
			c.JSON(http.StatusOK, world)
		})

		// DELETE /api/worlds/:id - DeleteWorldHandler
		worlds.DELETE("/worlds/:id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				slog.Warn("invalid world id", slog.String("error", err.Error()), slog.String("service", "worlds"))
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid world ID"})
				return
			}
			if err := repos.World.Delete(c.Request.Context(), id); err != nil {
				dblog.Error("failed to delete world", err, slog.String("service", "worlds"), slog.Int("world_id", id))
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			slog.Info("world deleted", slog.Int("world_id", id), slog.String("service", "worlds"))
			c.JSON(http.StatusOK, gin.H{"message": "World deleted"})
		})

		// GET /api/worlds/active - GetActiveWorldHandler
		worlds.GET("/worlds/active", func(c *gin.Context) {
			worlds, err := repos.World.GetActive(c.Request.Context())
			if err != nil {
				dblog.Error("failed to get active worlds", err, slog.String("service", "worlds"))
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if len(worlds) == 0 {
				c.JSON(http.StatusOK, gin.H{"error": "No active worlds"})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"world": worlds[0],
				"count": len(worlds),
			})
		})
	}

	fmt.Println("World CRUD routes registered")
}

// RegisterWorldRoutes registers world-scoped content endpoints
func RegisterWorldRoutes(router *gin.Engine, wm *content.WorldManager, repos *repository.Container) {
	router.GET("/worlds", func(c *gin.Context) {
		worlds := wm.GetAllWorlds()
		c.JSON(http.StatusOK, gin.H{
			"worlds": worlds,
			"count":  len(worlds),
		})
	})

	// Get active worlds (returns first active world for backward compatibility)
	router.GET("/worlds/active", func(c *gin.Context) {
		worlds := wm.GetActiveWorlds()
		if len(worlds) == 0 {
			c.JSON(http.StatusOK, gin.H{"error": "No active worlds"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"world": worlds[0],
			"count": len(worlds),
		})
	})

	// Get specific world
	router.GET("/worlds/:world_id", func(c *gin.Context) {
		worldID := c.Param("world_id")
		world, exists := wm.GetWorld(worldID)
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "World not found"})
			return
		}
		c.JSON(http.StatusOK, world)
	})

	// Get world stats
	router.GET("/worlds/:world_id/stats", func(c *gin.Context) {
		worldID := c.Param("world_id")
		stats, exists := wm.GetWorldStats(worldID)
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "World not found"})
			return
		}
		c.JSON(http.StatusOK, stats)
	})

	// World-scoped content routes
	// GET /worlds/:world_id/content/skills
	router.GET("/worlds/:world_id/content/skills", func(c *gin.Context) {
		worldID := c.Param("world_id")
		mgr, exists := wm.GetWorldManager(worldID)
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "World not found"})
			return
		}
		skills := mgr.Skills.GetAll()
		c.JSON(http.StatusOK, gin.H{
			"world_id": worldID,
			"skills":   skills,
			"count":    len(skills),
		})
	})

	// GET /worlds/:world_id/content/skills/:id
	router.GET("/worlds/:world_id/content/skills/:skill_id", func(c *gin.Context) {
		worldID := c.Param("world_id")
		skillID := c.Param("skill_id")
		mgr, exists := wm.GetWorldManager(worldID)
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "World not found"})
			return
		}
		skill, exists := mgr.Skills.Get(skillID)
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Skill not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"world_id": worldID,
			"skill":    skill,
		})
	})

	// GET /worlds/:world_id/content/npcs
	router.GET("/worlds/:world_id/content/npcs", func(c *gin.Context) {
		worldID := c.Param("world_id")
		mgr, exists := wm.GetWorldManager(worldID)
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "World not found"})
			return
		}
		npcs := mgr.NPCs.GetAll()
		c.JSON(http.StatusOK, gin.H{
			"world_id": worldID,
			"npcs":     npcs,
			"count":    len(npcs),
		})
	})

	// GET /worlds/:world_id/content/npcs/:id
	router.GET("/worlds/:world_id/content/npcs/:npc_id", func(c *gin.Context) {
		worldID := c.Param("world_id")
		npcID := c.Param("npc_id")
		mgr, exists := wm.GetWorldManager(worldID)
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "World not found"})
			return
		}
		npc, exists := mgr.NPCs.Get(npcID)
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "NPC not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"world_id": worldID,
			"npc":      npc,
		})
	})

	// GET /worlds/:world_id/content/items
	router.GET("/worlds/:world_id/content/items", func(c *gin.Context) {
		worldID := c.Param("world_id")
		mgr, exists := wm.GetWorldManager(worldID)
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "World not found"})
			return
		}
		items := mgr.Items.GetAll()
		c.JSON(http.StatusOK, gin.H{
			"world_id": worldID,
			"items":    items,
			"count":    len(items),
		})
	})

	// GET /worlds/:world_id/content/rooms
	router.GET("/worlds/:world_id/content/rooms", func(c *gin.Context) {
		worldID := c.Param("world_id")
		mgr, exists := wm.GetWorldManager(worldID)
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "World not found"})
			return
		}
		rooms := mgr.Rooms.GetAll()
		c.JSON(http.StatusOK, gin.H{
			"world_id": worldID,
			"rooms":    rooms,
			"count":    len(rooms),
		})
	})

	// GET /worlds/:world_id/content/quests
	router.GET("/worlds/:world_id/content/quests", func(c *gin.Context) {
		worldID := c.Param("world_id")
		mgr, exists := wm.GetWorldManager(worldID)
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "World not found"})
			return
		}
		quests := mgr.Quests.GetAll()
		c.JSON(http.StatusOK, gin.H{
			"world_id": worldID,
			"quests":   quests,
			"count":    len(quests),
		})
	})

	// Validate content for a world
	router.POST("/worlds/:world_id/content/validate", func(c *gin.Context) {
		worldID := c.Param("world_id")
		mgr, exists := wm.GetWorldManager(worldID)
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "World not found"})
			return
		}

		errors := mgr.Validate()
		if len(errors) > 0 {
			c.JSON(http.StatusOK, gin.H{
				"world_id": worldID,
				"valid":    false,
				"errors":   errors,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"world_id": worldID,
			"valid":    true,
			"message":  "All content validated successfully",
		})
	})

	fmt.Println("World routes registered")
}
