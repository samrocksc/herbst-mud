package routes

import (
	"herbst-server/dblog"
	"log/slog"

	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/repository"
)

type ExamineDisplayConfig struct {
	ShowName        bool `json:"showName"`
	ShowDescription bool `json:"showDescription"`
	ShowRace        bool `json:"showRace"`
	ShowLevel       bool `json:"showLevel"`
	ShowEquipped    bool `json:"showEquipped"`
	ShowUnequipped  bool `json:"showUnequipped"`
}

type EquipmentItem struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Slot        string `json:"slot,omitempty"`
	ItemType    string `json:"item_type,omitempty"`
}

type CharacterExamineResponse struct {
	Name            string           `json:"name,omitempty"`
	Description     string           `json:"description,omitempty"`
	Race            string           `json:"race,omitempty"`
	Level           int              `json:"level,omitempty"`
	EquippedItems   []EquipmentItem  `json:"equipped_items,omitempty"`
	UnequippedItems []EquipmentItem  `json:"unequipped_items,omitempty"`
}

func parseBoolQuery(c *gin.Context, key string, defaultVal bool) bool {
	if val := c.Query(key); val != "" {
		if b, err := strconv.ParseBool(val); err == nil {
			return b
		}
	}
	return defaultVal
}

func RegisterCharacterExamineRoutes(router *gin.RouterGroup, repos *repository.Container, client *db.Client) {
	// Character examine endpoint
	router.GET("/characters/:id/examine", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		char, err := repos.Character.Get(c.Request.Context(), id)
		if err != nil {
			dblog.Error("examine character failed", err, slog.String("service", "characters"))
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}

		config := getExamineConfig(c, repos)
		response := buildCharacterExamineResponse(c, char, config, repos)

		c.JSON(http.StatusOK, response)
	})

	// NPC instance examine endpoint - NPC instances are Characters with is_instance=true
	router.GET("/npc-instances/:id/examine", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid NPC instance ID"})
			return
		}

		// NPC instance is a Character row with is_instance=true
		npc, err := client.Character.Get(c.Request.Context(), id)
		if err != nil {
			dblog.Error("examine NPC instance failed", err, slog.String("service", "characters"))
			c.JSON(http.StatusNotFound, gin.H{"error": "NPC instance not found"})
			return
		}

		config := getExamineConfig(c, repos)
		response := buildNPCExamineResponse(c, npc, config, repos)

		c.JSON(http.StatusOK, response)
	})
}

func getExamineConfig(c *gin.Context, repos *repository.Container) ExamineDisplayConfig {
	var config ExamineDisplayConfig
	cfg, err := repos.GameConfig.Get(c.Request.Context(), "examine_display_config")
	if err != nil {
		config = ExamineDisplayConfig{
			ShowName: true, ShowDescription: true, ShowRace: true,
			ShowLevel: true, ShowEquipped: true, ShowUnequipped: true,
		}
	} else if err := json.Unmarshal([]byte(cfg.Value), &config); err != nil {
		config = ExamineDisplayConfig{
			ShowName: true, ShowDescription: true, ShowRace: true,
			ShowLevel: true, ShowEquipped: true, ShowUnequipped: true,
		}
	}

	// Query param overrides
	config.ShowName = parseBoolQuery(c, "showName", config.ShowName)
	config.ShowDescription = parseBoolQuery(c, "showDescription", config.ShowDescription)
	config.ShowRace = parseBoolQuery(c, "showRace", config.ShowRace)
	config.ShowLevel = parseBoolQuery(c, "showLevel", config.ShowLevel)
	config.ShowEquipped = parseBoolQuery(c, "showEquipped", config.ShowEquipped)
	config.ShowUnequipped = parseBoolQuery(c, "showUnequipped", config.ShowUnequipped)

	return config
}

func buildCharacterExamineResponse(c *gin.Context, char *db.Character, config ExamineDisplayConfig, repos *repository.Container) CharacterExamineResponse {
	response := CharacterExamineResponse{}

	if config.ShowName {
		response.Name = char.Name
	}
	if config.ShowDescription && char.Description != "" {
		response.Description = char.Description
	}
	if config.ShowRace {
		response.Race = char.Race
	}
	if config.ShowLevel {
		response.Level = char.Level
	}

	// Fetch equipment
	equipment, err := repos.Equipment.ListByOwner(c.Request.Context(), char.ID)
	if err == nil && equipment != nil {
		if config.ShowEquipped {
			response.EquippedItems = []EquipmentItem{}
			for _, eq := range equipment {
				if eq.IsEquipped {
					response.EquippedItems = append(response.EquippedItems, EquipmentItem{
						ID: eq.ID, Name: eq.Name, Description: eq.Description,
						Slot: eq.Slot, ItemType: eq.ItemType,
					})
				}
			}
		}
		if config.ShowUnequipped {
			response.UnequippedItems = []EquipmentItem{}
			for _, eq := range equipment {
				if !eq.IsEquipped {
					response.UnequippedItems = append(response.UnequippedItems, EquipmentItem{
						ID: eq.ID, Name: eq.Name, Description: eq.Description,
						Slot: eq.Slot, ItemType: eq.ItemType,
					})
				}
			}
		}
	}

	return response
}

func buildNPCExamineResponse(c *gin.Context, npc *db.Character, config ExamineDisplayConfig, repos *repository.Container) CharacterExamineResponse {
	response := CharacterExamineResponse{}

	if config.ShowName {
		response.Name = npc.Name
	}
	if config.ShowDescription && npc.Description != "" {
		response.Description = npc.Description
	}
	if config.ShowRace {
		response.Race = npc.Race
	}
	if config.ShowLevel {
		response.Level = npc.Level
	}

	// NPC instances can have equipment too - they're stored in the Equipment table
	// with owner_id pointing to the Character ID
	equipment, err := repos.Equipment.ListByOwner(c.Request.Context(), npc.ID)
	if err == nil && equipment != nil {
		if config.ShowEquipped {
			response.EquippedItems = []EquipmentItem{}
			for _, eq := range equipment {
				if eq.IsEquipped {
					response.EquippedItems = append(response.EquippedItems, EquipmentItem{
						ID: eq.ID, Name: eq.Name, Description: eq.Description,
						Slot: eq.Slot, ItemType: eq.ItemType,
					})
				}
			}
		}
		if config.ShowUnequipped {
			response.UnequippedItems = []EquipmentItem{}
			for _, eq := range equipment {
				if !eq.IsEquipped {
					response.UnequippedItems = append(response.UnequippedItems, EquipmentItem{
						ID: eq.ID, Name: eq.Name, Description: eq.Description,
						Slot: eq.Slot, ItemType: eq.ItemType,
					})
				}
			}
		}
	}

	return response
}
