package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/dbinit"
)

// WipeRequest configures what to wipe and reload
type WipeRequest struct {
	WipeNPCs      bool `json:"wipe_npcs"`
	WipeRooms     bool `json:"wipe_rooms"`
	WipeItems     bool `json:"wipe_items"`
	WipeSkills    bool `json:"wipe_skills"`
	WipeTalents   bool `json:"wipe_talents"`
	PreserveUsers bool `json:"preserve_users"`
}

// WipeResult reports what was wiped and reloaded
type WipeResult struct {
	NPCsWiped    int      `json:"npcs_wiped"`
	RoomsWiped   int      `json:"rooms_wiped"`
	ItemsWiped   int      `json:"items_wiped"`
	SkillsWiped  int      `json:"skills_wiped"`
	TalentsWiped int      `json:"talents_wiped"`
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

			// Re-initialize characters
			if err := dbinit.InitCharacters(client); err != nil {
				result.Errors = append(result.Errors, "Failed to reinitialize characters: "+err.Error())
			} else {
				result.Reinitialized = append(result.Reinitialized, "characters")
			}
		}

		// Wipe rooms
		if req.WipeRooms {
			// Get count before deletion
			count, _ := client.Room.Query().Count(ctx)
			
			// Delete all rooms
			client.Room.Delete().Exec(ctx)
			result.RoomsWiped = count

			// Re-initialize rooms
			if err := dbinit.InitCrossWay(client); err != nil {
				result.Errors = append(result.Errors, "Failed to reinitialize cross-way rooms: "+err.Error())
			} else {
				result.Reinitialized = append(result.Reinitialized, "cross-way rooms")
			}

			if err := dbinit.InitFountain(client); err != nil {
				result.Errors = append(result.Errors, "Failed to reinitialize fountain: "+err.Error())
			} else {
				result.Reinitialized = append(result.Reinitialized, "fountain")
			}

			if err := dbinit.InitJunkyard(client); err != nil {
				result.Errors = append(result.Errors, "Failed to reinitialize junkyard: "+err.Error())
			} else {
				result.Reinitialized = append(result.Reinitialized, "junkyard")
			}
		}

		// Wipe items/equipment
		if req.WipeItems {
			count, _ := client.Equipment.Query().Count(ctx)
			client.Equipment.Delete().Exec(ctx)
			result.ItemsWiped = count

			// Re-initialize items
			if err := dbinit.InitConsumables(client); err != nil {
				result.Errors = append(result.Errors, "Failed to reinitialize consumables: "+err.Error())
			} else {
				result.Reinitialized = append(result.Reinitialized, "consumables")
			}
		}

		// Wipe skills
		if req.WipeSkills {
			count, _ := client.Skill.Query().Count(ctx)
			client.Skill.Delete().Exec(ctx)
			result.SkillsWiped = count

			if err := dbinit.InitSkills(client); err != nil {
				result.Errors = append(result.Errors, "Failed to reinitialize skills: "+err.Error())
			} else {
				result.Reinitialized = append(result.Reinitialized, "skills")
			}
		}

		// Wipe talents
		if req.WipeTalents {
			count, _ := client.Talent.Query().Count(ctx)
			client.Talent.Delete().Exec(ctx)
			result.TalentsWiped = count

			if err := dbinit.InitTalents(client); err != nil {
				result.Errors = append(result.Errors, "Failed to reinitialize talents: "+err.Error())
			} else {
				result.Reinitialized = append(result.Reinitialized, "talents")
			}
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

		// Wipe skills
		skillCount, _ := client.Skill.Query().Count(ctx)
		client.Skill.Delete().Exec(ctx)
		result.SkillsWiped = skillCount

		// Wipe talents
		talentCount, _ := client.Talent.Query().Count(ctx)
		client.Talent.Delete().Exec(ctx)
		result.TalentsWiped = talentCount

		// Re-initialize everything
		initializers := []struct {
			name string
			fn   func(*db.Client) error
		}{
			{"cross-way rooms", dbinit.InitCrossWay},
			{"fountain", dbinit.InitFountain},
			{"junkyard", dbinit.InitJunkyard},
			{"characters", dbinit.InitCharacters},
			{"consumables", dbinit.InitConsumables},
			{"skills", dbinit.InitSkills},
			{"talents", dbinit.InitTalents},
			{"fountain", dbinit.InitFountain},
			{"gizmo", dbinit.InitGizmoNPC},
		}

		for _, init := range initializers {
			if err := init.fn(client); err != nil {
				result.Errors = append(result.Errors, "Failed to reinitialize "+init.name+": "+err.Error())
			} else {
				result.Reinitialized = append(result.Reinitialized, init.name)
			}
		}

		c.JSON(http.StatusOK, result)
	})
}
