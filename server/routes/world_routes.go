package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/content"
)

// RegisterWorldRoutes registers world-scoped content endpoints
func RegisterWorldRoutes(router *gin.Engine, wm *content.WorldManager) {
	// Get all worlds
	router.GET("/worlds", func(c *gin.Context) {
		worlds := wm.GetAllWorlds()
		c.JSON(http.StatusOK, gin.H{
			"worlds": worlds,
			"count":  len(worlds),
		})
	})

	// Get active worlds
	router.GET("/worlds/active", func(c *gin.Context) {
		worlds := wm.GetActiveWorlds()
		c.JSON(http.StatusOK, gin.H{
			"worlds": worlds,
			"count":  len(worlds),
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
