package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"entgo.io/ent/dialect/sql"
	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/systemlog"
)

// BiographyEntry is a single event in a character's biography timeline.
type BiographyEntry struct {
	ID          int64       `json:"id"`
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Timestamp   string      `json:"timestamp"`
	Details     interface{} `json:"details"`
}

// BiographySummary holds aggregated stats for a character's biography.
type BiographySummary struct {
	TotalEvents     int `json:"total_events"`
	BossKills       int `json:"boss_kills"`
	NPCKills        int `json:"npc_kills"`
	Deaths          int `json:"deaths"`
	LevelUps        int `json:"level_ups"`
	ClassChanges    int `json:"class_changes"`
	QuestCompletions int `json:"quest_completions"`
	SkillsLearned   int `json:"skills_learned"`
}

// BiographyResponse is the full biography API response.
type BiographyResponse struct {
	CharacterID   int64             `json:"character_id"`
	CharacterName string            `json:"character_name"`
	Entries       []BiographyEntry  `json:"entries"`
	Summary       BiographySummary  `json:"summary"`
}

// formatBiographyEntry creates a human-readable description from an event log entry.
func formatBiographyEntry(action string, details map[string]interface{}) string {
	switch action {
	case "boss_killed":
		bossName, _ := details["npc_name"].(string)
		bossLevel, _ := details["npc_level"].(float64)
		charName, _ := details["character_name"].(string)
		return fmt.Sprintf("%s defeated %s (Level %d boss)", charName, bossName, int(bossLevel))
	case "npc_defeated":
		npcName, _ := details["npc_name"].(string)
		charName, _ := details["character_name"].(string)
		if npcName == "" {
			npcID, _ := details["npc_id"].(float64)
			return fmt.Sprintf("%s defeated NPC #%d", charName, int(npcID))
		}
		return fmt.Sprintf("%s defeated %s", charName, npcName)
	case "character_died":
		charName, _ := details["character_name"].(string)
		attackerName, _ := details["attacker_name"].(string)
		if attackerName != "" {
			return fmt.Sprintf("%s was slain by %s", charName, attackerName)
		}
		return fmt.Sprintf("%s died", charName)
	case "character_respawned":
		charName, _ := details["character_name"].(string)
		return fmt.Sprintf("%s respawned at rally point", charName)
	case "level_up":
		newLevel, _ := details["new_level"].(float64)
		return fmt.Sprintf("Reached level %d", int(newLevel))
	case "reclass":
		return "Changed class"
	case "rerace":
		return "Changed race"
	case "quest_complete":
		return "Completed a quest"
	case "skill_learned":
		return "Learned a new skill"
	case "skill.leveled_up":
		return "Skill leveled up"
	case "skill.ability_unlocked":
		return "New ability unlocked"
	case "xp.gained":
		return "Gained XP"
	case "skill_xp.gained":
		return "Gained skill XP"
	default:
		return action
	}
}

// handleGetBiography returns a character's biography timeline from system_logs.
// GET /characters/:id/biography?type=boss_killed&limit=50
func handleGetBiography(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := getIDParam(c)
		if !ok {
			return
		}

		// Get character name
		char, err := client.Character.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "character not found"})
			return
		}

		// Build query
		query := client.SystemLog.Query().
			Where(systemlog.CharacterIDEQ(id)).
			Order(systemlog.ByTimestamp(sql.OrderDesc()))

		// Filter by event type if provided
		if eventType := c.Query("type"); eventType != "" {
			query = query.Where(systemlog.ActionEQ(eventType))
		}

		// Default limit 100, max 500
		limit := 100
		if l := c.Query("limit"); l != "" {
			if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 500 {
				limit = parsed
			}
		}
		query = query.Limit(limit)

		logs, err := query.All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("query biography: %v", err)})
			return
		}

		// Build entries
		entries := make([]BiographyEntry, 0, len(logs))
		summary := BiographySummary{}
		summary.TotalEvents = len(logs)

		for _, log := range logs {
			var details map[string]interface{}
			if log.Details != "" {
				_ = json.Unmarshal([]byte(log.Details), &details)
			}

			entries = append(entries, BiographyEntry{
				ID:          int64(log.ID),
				Type:        log.Action,
				Description: formatBiographyEntry(log.Action, details),
				Timestamp:   log.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
				Details:     details,
			})

			// Update summary
			switch log.Action {
			case "boss_killed":
				summary.BossKills++
			case "npc_defeated":
				summary.NPCKills++
			case "character_died":
				summary.Deaths++
			case "level_up":
				summary.LevelUps++
			case "reclass", "rerace":
				summary.ClassChanges++
			case "quest_complete":
				summary.QuestCompletions++
			case "skill_learned":
				summary.SkillsLearned++
			}
		}

		c.JSON(http.StatusOK, BiographyResponse{
			CharacterID:   int64(id),
			CharacterName: char.Name,
			Entries:       entries,
			Summary:       summary,
		})
	}
}