package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/character"
)

// WipeRequest configures what to wipe and reload
type WipeRequest struct {
	WipeNPCs      bool `json:"wipe_npcs"`
	WipeRooms     bool `json:"wipe_rooms"`
	WipeItems     bool `json:"wipe_items"`
	WipeSkills    bool `json:"wipe_skills"`
	PreserveUsers bool `json:"preserve_users"`
}

// WipeResult reports what was wiped and reloaded
type WipeResult struct {
	NPCsWiped    int      `json:"npcs_wiped"`
	RoomsWiped   int      `json:"rooms_wiped"`
	ItemsWiped   int      `json:"items_wiped"`
	SkillsWiped  int      `json:"skills_wiped"`
	Reinitialized []string `json:"reinitialized"`
	Errors       []string `json:"errors,omitempty"`
}

// RegisterAdminWipeRoutes registers the wipe/reload routes
func RegisterAdminWipeRoutes(router *gin.Engine, client *db.Client) {
	// Wipe and reload game data
	router.POST("/admin/wipe", func(c *gin.Context) {
		var req WipeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result := WipeResult{
			Reinitialized: make([]string, 0),
			Errors:        make([]string, 0),
		}

		ctx := c.Request.Context()

		// Wipe NPCs (characters with isNPC=true)
		if req.WipeNPCs {
			npcs, err := client.Character.Query().
				Where(character.IsNPCEQ(true)).
				All(ctx)
			if err == nil {
				for _, npc := range npcs {
					client.Character.DeleteOne(npc).Exec(ctx)
					result.NPCsWiped++
				}
			}

		}

		// Wipe rooms
		if req.WipeRooms {
			// Get count before deletion
			count, _ := client.Room.Query().Count(ctx)
			client.Room.Delete().Exec(ctx)
			result.RoomsWiped = count
		}

		// Wipe items/equipment
		if req.WipeItems {
			count, _ := client.Equipment.Query().Count(ctx)
			client.Equipment.Delete().Exec(ctx)
			result.ItemsWiped = count
		}

		// Wipe abilities (formerly skills)
		if req.WipeSkills {
			count, _ := client.Ability.Query().Count(ctx)
			client.Ability.Delete().Exec(ctx)
			result.SkillsWiped = count
		}

		c.JSON(http.StatusOK, result)
	})

	// Full wipe - resets everything (except users by default)
	router.POST("/admin/wipe/full", func(c *gin.Context) {
		result := WipeResult{
			Reinitialized: make([]string, 0),
			Errors:        make([]string, 0),
		}

		ctx := c.Request.Context()

		// Wipe NPCs
		npcs, _ := client.Character.Query().
			Where(character.IsNPCEQ(true)).
			All(ctx)
		for _, npc := range npcs {
			client.Character.DeleteOne(npc).Exec(ctx)
			result.NPCsWiped++
		}

		// Wipe rooms
		count, _ := client.Room.Query().Count(ctx)
		client.Room.Delete().Exec(ctx)
		result.RoomsWiped = count

		// Wipe equipment
		itemCount, _ := client.Equipment.Query().Count(ctx)
		client.Equipment.Delete().Exec(ctx)
		result.ItemsWiped = itemCount

		// Wipe abilities (formerly skills)
		skillCount, _ := client.Ability.Query().Count(ctx)
		client.Ability.Delete().Exec(ctx)
		result.SkillsWiped = skillCount

		c.JSON(http.StatusOK, result)
	})
}
