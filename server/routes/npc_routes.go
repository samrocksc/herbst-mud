package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
)

// RegisterNPCRoutes registers all NPC-related routes
func RegisterNPCRoutes(router *gin.Engine, client *db.Client) {
	group := router.Group("/npcs")
	{
		// Get all NPCs
		group.GET("", func(c *gin.Context) {
			npchars, err := client.Character.Query().
				Where(db.Character.IsNPC(true)).
				All(c.Request.Context())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, npchars)
		})

		// Get NPC by ID
		group.GET("/:id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid NPC ID"})
				return
			}

			npc, err := client.Character.Get(c.Request.Context(), id)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "NPC not found"})
				return
			}

			if !npc.IsNPC {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Character is not an NPC"})
				return
			}

			c.JSON(http.StatusOK, npc)
		})

		// Get NPC examine details
		group.GET("/:id/examine", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid NPC ID"})
				return
			}

			npc, err := client.Character.Get(c.Request.Context(), id)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "NPC not found"})
				return
			}

			if !npc.IsNPC {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Character is not an NPC"})
				return
			}

			// Get character's examine skill level
			examineLevel := 1
			if charIDStr := c.Query("character_id"); charIDStr != "" {
				charID, err := strconv.Atoi(charIDStr)
				if err == nil {
					char, err := client.Character.Get(c.Request.Context(), charID)
					if err == nil {
						// Derive examine level from INT stat
						examineLevel = 1 + (char.IntStat / 5)
					}
				}
			}

			// Build NPC examine response
			// In a full implementation, NPCs would have examine_desc and hidden_details
			c.JSON(http.StatusOK, gin.H{
				"id":            npc.ID,
				"name":          npc.Name,
				"description":   "A " + npc.Name + " stands here.",
				"examineDesc":   "You examine " + npc.Name + " closely.",
				"isNPC":         npc.IsNPC,
				"level":         npc.Level,
				"hitpoints":     npc.Hitpoints,
				"maxHitpoints":  npc.MaxHitpoints,
				"examineLevel":  examineLevel,
				"disposition":   getNPCDisposition(npc),
				"trades":        getNPCCredits(npc),
			})
		})
	}
}

// getNPCDisposition returns the NPC's general disposition based on stats
func getNPCDisposition(npc *db.Character) string {
	// Simple disposition logic based on NPC type/name
	// In full implementation, this would be stored in the database
	return "neutral"
}

// getNPCCredits returns the NPC's credit count for trading
func getNPCCredits(npc *db.Character) int {
	// In full implementation, this would query a separate NPC inventory table
	return 0
}