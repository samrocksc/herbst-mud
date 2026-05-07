package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/race"
)

// updateRace updates an existing race.
func updateRace(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid race id"})
			return
		}

		var req struct {
			Name           *string  `json:"name"`
			DisplayName    *string  `json:"display_name"`
			Description    *string  `json:"description"`
			StatModifiers  *string  `json:"stat_modifiers"`
			SkillGrants    []string `json:"skill_grants"`
			EquipmentSlots []string `json:"equipment_slots"`
			IsPlayable     *bool    `json:"is_playable"`
			Color          *string  `json:"color"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		if err := validateSlots(req.EquipmentSlots); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Name != nil {
			existing, err := client.Race.Query().Where(race.NameEQ(*req.Name)).Only(c.Request.Context())
			if err == nil && existing != nil && existing.ID != id {
				c.JSON(http.StatusConflict, gin.H{"error": "a race with this name already exists"})
				return
			}
		}

		mut := client.Race.UpdateOneID(id)
		if req.Name != nil {
			mut = mut.SetName(*req.Name)
		}
		if req.DisplayName != nil {
			mut = mut.SetDisplayName(*req.DisplayName)
		}
		if req.Description != nil {
			mut = mut.SetDescription(*req.Description)
		}
		if req.StatModifiers != nil {
			mut = mut.SetStatModifiers(*req.StatModifiers)
		}
		if req.IsPlayable != nil {
			mut = mut.SetIsPlayable(*req.IsPlayable)
		}
		if req.Color != nil {
			mut = mut.SetColor(*req.Color)
		}
		if req.EquipmentSlots != nil {
			mut = mut.SetEquipmentSlots(req.EquipmentSlots)
		}

		r, err := mut.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "race not found or update failed"})
			return
		}
		c.JSON(http.StatusOK, raceToView(r))
	}
}