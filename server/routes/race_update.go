package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/repository"
)

// updateRace updates an existing race.
func updateRace(repos *repository.Container, client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid race id"})
			return
		}

		var req struct {
			Name            *string  `json:"name"`
			DisplayName     *string  `json:"display_name"`
			Description     *string  `json:"description"`
			StatModifiers   *string  `json:"stat_modifiers"`
			SkillGrants     []string `json:"skill_grants"`
			EquipmentSlots  []string `json:"equipment_slots"`
			RequirementTags []string `json:"requirement_tags"`
			Color           *string  `json:"color"`
			Tags            []string `json:"tags"`
			WorldID         *string  `json:"world_id"`
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
			// Get the existing race to get its world_id for duplicate check
			existingRace, err := repos.Race.Get(c.Request.Context(), id)
			if err == nil {
				existing, err := repos.Race.GetByName(c.Request.Context(), *req.Name, existingRace.WorldID)
				if err == nil && existing != nil && existing.ID != id {
					c.JSON(http.StatusConflict, gin.H{"error": "a race with this name already exists in this world"})
					return
				}
			}
		}

		updates := repository.RaceUpdates{
			Name:            req.Name,
			DisplayName:     req.DisplayName,
			Description:     req.Description,
			StatModifiers:   req.StatModifiers,
			RequirementTags: req.RequirementTags,
			Color:           req.Color,
			EquipmentSlots:  req.EquipmentSlots,
		}

		if req.Tags != nil {
			tagIDs, err := resolveTagIDs(c, client, req.Tags)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			updates.ClearTags = true
			updates.AddTagIDs = tagIDs
		}

		r, err := repos.Race.Update(c.Request.Context(), id, updates)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "race not found or update failed"})
			return
		}

		c.JSON(http.StatusOK, raceToView(r))
	}
}
