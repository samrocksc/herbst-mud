package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/content"
)

// RegisterContentRoutes registers content API endpoints
func RegisterContentRoutes(router *gin.Engine, mgr *content.Manager) {
	// Get all skills
	router.GET("/content/skills", func(c *gin.Context) {
		skills := mgr.Skills.GetAll()
		c.JSON(http.StatusOK, gin.H{
			"skills": skills,
			"count":  len(skills),
		})
	})

	// Get skill by ID
	router.GET("/content/skills/:id", func(c *gin.Context) {
		id := c.Param("id")
		skill, exists := mgr.Skills.Get(id)
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Skill not found"})
			return
		}
		c.JSON(http.StatusOK, skill)
	})

	// Get skills by tag
	router.GET("/content/skills/tag/:tag", func(c *gin.Context) {
		tag := c.Param("tag")
		skills := mgr.Skills.GetByTag(tag)
		c.JSON(http.StatusOK, gin.H{
			"skills": skills,
			"count":  len(skills),
		})
	})

	// Get skills by class
	router.GET("/content/skills/class/:class", func(c *gin.Context) {
		class := c.Param("class")
		skills := mgr.Skills.GetByClass(class)
		c.JSON(http.StatusOK, gin.H{
			"skills": skills,
			"count":  len(skills),
		})
	})

	// Get all items
	router.GET("/content/items", func(c *gin.Context) {
		items := mgr.Items.GetAll()
		c.JSON(http.StatusOK, gin.H{
			"items": items,
			"count": len(items),
		})
	})

	// Get item by ID
	router.GET("/content/items/:id", func(c *gin.Context) {
		id := c.Param("id")
		item, exists := mgr.Items.Get(id)
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
			return
		}
		c.JSON(http.StatusOK, item)
	})

	// Get all NPCs
	router.GET("/content/npcs", func(c *gin.Context) {
		npcs := mgr.NPCs.GetAll()
		c.JSON(http.StatusOK, gin.H{
			"npcs":  npcs,
			"count": len(npcs),
		})
	})

	// Get NPC by ID
	router.GET("/content/npcs/:id", func(c *gin.Context) {
		id := c.Param("id")
		npc, exists := mgr.NPCs.Get(id)
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "NPC not found"})
			return
		}
		c.JSON(http.StatusOK, npc)
	})

	// Content statistics
	router.GET("/content/stats", func(c *gin.Context) {
		stats := mgr.GetStats()
		c.JSON(http.StatusOK, stats)
	})

	// Get all rooms
	router.GET("/content/rooms", func(c *gin.Context) {
		rooms := mgr.Rooms.GetAll()
		c.JSON(http.StatusOK, gin.H{
			"rooms": rooms,
			"count": len(rooms),
		})
	})

	// Get room by ID
	router.GET("/content/rooms/:id", func(c *gin.Context) {
		id := c.Param("id")
		room, exists := mgr.Rooms.Get(id)
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
			return
		}
		c.JSON(http.StatusOK, room)
	})

	// Get connected rooms (exits)
	router.GET("/content/rooms/:id/exits", func(c *gin.Context) {
		id := c.Param("id")
		room, exists := mgr.Rooms.Get(id)
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
			return
		}
		
		exits := make(map[string]interface{})
		for direction, targetID := range room.Exits {
			if targetRoom, exists := mgr.Rooms.Get(targetID); exists {
				exits[direction] = gin.H{
					"room_id": targetID,
					"name":    targetRoom.Name,
				}
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"room_id": id,
			"exits":   exits,
		})
	})

	// Get all quests
	router.GET("/content/quests", func(c *gin.Context) {
		quests := mgr.Quests.GetAll()
		c.JSON(http.StatusOK, gin.H{
			"quests": quests,
			"count":  len(quests),
		})
	})

	// Get quest by ID
	router.GET("/content/quests/:id", func(c *gin.Context) {
		id := c.Param("id")
		quest, exists := mgr.Quests.Get(id)
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Quest not found"})
			return
		}
		c.JSON(http.StatusOK, quest)
	})

	// Get quests by difficulty
	router.GET("/content/quests/difficulty/:difficulty", func(c *gin.Context) {
		difficulty := c.Param("difficulty")
		quests := mgr.Quests.GetByDifficulty(difficulty)
		c.JSON(http.StatusOK, gin.H{
			"quests": quests,
			"count":  len(quests),
		})
	})

	// Get quests by type
	router.GET("/content/quests/type/:type", func(c *gin.Context) {
		questType := c.Param("type")
		quests := mgr.Quests.GetByType(questType)
		c.JSON(http.StatusOK, gin.H{
			"quests": quests,
			"count":  len(quests),
		})
	})

	// Admin: Validate content changes before saving
	router.POST("/admin/content/validate", func(c *gin.Context) {
		var request struct {
			Type string          `json:"type"`
			Data json.RawMessage `json:"data"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate based on content type
		valid, errors := mgr.ValidateContentChange(request.Type, request.Data)
		
		c.JSON(http.StatusOK, gin.H{
			"valid":  valid,
			"errors": errors,
		})
	})

	// Admin: Preview content before saving
	router.POST("/admin/content/preview", func(c *gin.Context) {
		var request struct {
			Type string          `json:"type"`
			Data json.RawMessage `json:"data"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		preview, err := mgr.PreviewContent(request.Type, request.Data)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"preview": preview,
			"valid":   true,
		})
	})

	// Validate content (admin endpoint)
	router.GET("/content/validate", func(c *gin.Context) {
		errors := mgr.Validate()
		if len(errors) > 0 {
			c.JSON(http.StatusOK, gin.H{
				"valid":  false,
				"errors": errors,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"valid":   true,
			"message": "All content validated successfully",
		})
	})
}
