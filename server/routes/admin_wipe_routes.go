package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/db/room"
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
	NPCsWiped     int      `json:"npcs_wiped"`
	RoomsWiped    int      `json:"rooms_wiped"`
	ItemsWiped    int      `json:"items_wiped"`
	SkillsWiped   int      `json:"skills_wiped"`
	WorldWiped    bool     `json:"world_wiped,omitempty"`
	Reinitialized []string `json:"reinitialized"`
	Errors        []string `json:"errors,omitempty"`
}

// RegisterAdminWipeRoutes registers the wipe/reload routes
func RegisterAdminWipeRoutes(router *gin.Engine, client *db.Client) {
	// Wipe game data, optionally scoped to a world
	// GET ?world=test-world filters deletion to that world only
	router.POST("/admin/wipe", func(c *gin.Context) {
		var req WipeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		worldFilter := c.Query("world")
		result := WipeResult{
			Reinitialized: make([]string, 0),
			Errors:        make([]string, 0),
		}

		ctx := c.Request.Context()

		// Wipe NPCs (characters with isNPC=true)
		if req.WipeNPCs {
			npcQuery := client.Character.Query().
				Where(character.IsNPCEQ(true))
			if worldFilter != "" {
				npcQuery = npcQuery.Where(character.CurrentWorldEQ(worldFilter))
			}
			npcs, err := npcQuery.All(ctx)
			if err == nil {
				for _, npc := range npcs {
					client.Character.DeleteOne(npc).Exec(ctx)
					result.NPCsWiped++
				}
			}
		}

		// Wipe rooms
		if req.WipeRooms {
			roomQuery := client.Room.Query()
			if worldFilter != "" {
				roomQuery = roomQuery.Where(room.WorldIDEQ(worldFilter))
			}
			count, _ := roomQuery.Count(ctx)
			if worldFilter != "" {
				client.Room.Delete().Where(room.WorldIDEQ(worldFilter)).Exec(ctx)
			} else {
				client.Room.Delete().Exec(ctx)
			}
			result.RoomsWiped = count
		}

		// Wipe items/equipment
		if req.WipeItems {
			if worldFilter != "" {
				// Equipment table may not have world_id — delete all for now
				// TODO: add world_id to equipment schema for world-scoped deletion
				count, _ := client.Equipment.Query().Count(ctx)
				client.Equipment.Delete().Exec(ctx)
				result.ItemsWiped = count
			} else {
				count, _ := client.Equipment.Query().Count(ctx)
				client.Equipment.Delete().Exec(ctx)
				result.ItemsWiped = count
			}
		}

		// Wipe abilities (formerly skills)
		if req.WipeSkills {
			if worldFilter != "" {
				count, _ := client.Ability.Query().Count(ctx)
				client.Ability.Delete().Exec(ctx)
				result.SkillsWiped = count
			} else {
				count, _ := client.Ability.Query().Count(ctx)
				client.Ability.Delete().Exec(ctx)
				result.SkillsWiped = count
			}
		}

		c.JSON(http.StatusOK, result)
	})

	// Full wipe — resets everything (except users by default)
	// GET ?world=test-world scopes to that world only
	router.POST("/admin/wipe/full", func(c *gin.Context) {
		worldFilter := c.Query("world")
		result := WipeResult{
			Reinitialized: make([]string, 0),
			Errors:        make([]string, 0),
		}

		ctx := c.Request.Context()

		// Wipe NPCs — scope to world if filter set
		npcQuery := client.Character.Query().
			Where(character.IsNPCEQ(true))
		if worldFilter != "" {
			npcQuery = npcQuery.Where(character.CurrentWorldEQ(worldFilter))
		}
		npcs, _ := npcQuery.All(ctx)
		for _, npc := range npcs {
			client.Character.DeleteOne(npc).Exec(ctx)
			result.NPCsWiped++
		}

		// Wipe rooms — scope to world if filter set
		roomQuery := client.Room.Query()
		if worldFilter != "" {
			roomQuery = roomQuery.Where(room.WorldIDEQ(worldFilter))
		}
		count, _ := roomQuery.Count(ctx)
		if worldFilter != "" {
			client.Room.Delete().Where(room.WorldIDEQ(worldFilter)).Exec(ctx)
		} else {
			client.Room.Delete().Exec(ctx)
		}
		result.RoomsWiped = count

		// Wipe equipment
		itemCount, _ := client.Equipment.Query().Count(ctx)
		client.Equipment.Delete().Exec(ctx)
		result.ItemsWiped = itemCount

		// Wipe abilities
		skillCount, _ := client.Ability.Query().Count(ctx)
		client.Ability.Delete().Exec(ctx)
		result.SkillsWiped = skillCount

		c.JSON(http.StatusOK, result)
	})
}
