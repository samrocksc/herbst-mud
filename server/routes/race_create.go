package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/constants"
	"herbst-server/db"
	"herbst-server/db/race"
)

// createRace creates a new race.
func createRace(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name           string   `json:"name" binding:"required"`
			DisplayName    string   `json:"display_name"`
			Description    string   `json:"description"`
			StatModifiers  *string  `json:"stat_modifiers"`
			SkillGrants    []string `json:"skill_grants"`
			EquipmentSlots []string `json:"equipment_slots"`
			IsPlayable     *bool    `json:"is_playable"`
			Color          string   `json:"color"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
			return
		}

		if err := validateSlots(req.EquipmentSlots); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		existing, err := client.Race.Query().Where(race.NameEQ(req.Name)).Only(c.Request.Context())
		if err == nil && existing != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "a race with this name already exists"})
			return
		}

		isPlayable := true
		if req.IsPlayable != nil {
			isPlayable = *req.IsPlayable
		}
		displayName := req.DisplayName
		if displayName == "" {
			displayName = req.Name
		}

		mut := client.Race.Create().
			SetName(req.Name).
			SetDisplayName(displayName).
			SetDescription(req.Description).
			SetIsPlayable(isPlayable).
			SetColor(req.Color)

		if req.StatModifiers != nil {
			mut = mut.SetStatModifiers(*req.StatModifiers)
		}
		if len(req.EquipmentSlots) > 0 {
			mut = mut.SetEquipmentSlots(req.EquipmentSlots)
		}

		r, err := mut.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, raceToView(r))
	}
}

// validateSlots checks that every slot name is in SlotCatalog.
func validateSlots(slots []string) error {
	for _, slot := range slots {
		if !constants.IsValidSlot(slot) {
			return fmt.Errorf("invalid equipment slot: %s", slot)
		}
	}
	return nil
}