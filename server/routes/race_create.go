package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/constants"
	"herbst-server/db"
	"herbst-server/db/tag"
	"herbst-server/repository"
)

// createRace creates a new race.
func createRace(repos *repository.Container, client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name            string   `json:"name" binding:"required"`
			DisplayName     string   `json:"display_name"`
			Description     string   `json:"description"`
			StatModifiers   *string  `json:"stat_modifiers"`
			SkillGrants     []string `json:"skill_grants"`
			EquipmentSlots  []string `json:"equipment_slots"`
			RequirementTags []string `json:"requirement_tags"`
			Color           string   `json:"color"`
			Tags            []string `json:"tags"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
			return
		}

		if err := validateSlots(req.EquipmentSlots); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check for duplicate name
		existing, err := repos.Race.GetByName(c.Request.Context(), req.Name)
		if err == nil && existing != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "a race with this name already exists"})
			return
		}

		displayName := req.DisplayName
		if displayName == "" {
			displayName = req.Name
		}

		// Resolve tag names to IDs
		var tagIDs []int
		if len(req.Tags) > 0 {
			tagIDs, err = resolveTagIDs(c, client, req.Tags)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}

		r, err := repos.Race.Create(c.Request.Context(), repository.CreateRaceInput{
			Name:            req.Name,
			DisplayName:     displayName,
			Description:     req.Description,
			StatModifiers:   req.StatModifiers,
			RequirementTags: req.RequirementTags,
			Color:           req.Color,
			EquipmentSlots:  req.EquipmentSlots,
			TagIDs:          tagIDs,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, raceToView(r))
	}
}

// resolveTagIDs resolves tag names to IDs, creating tags that don't exist yet.
// TODO: Move tag resolution to TagRepo
func resolveTagIDs(c *gin.Context, client *db.Client, names []string) ([]int, error) {
	if len(names) == 0 {
		return nil, nil
	}
	existing, err := client.Tag.Query().Where(tag.NameIn(names...)).All(c.Request.Context())
	if err != nil {
		return nil, err
	}
	existingNames := make(map[string]bool)
	for _, t := range existing {
		existingNames[t.Name] = true
	}
	var ids []int
	for _, t := range existing {
		ids = append(ids, t.ID)
	}
	for _, name := range names {
		if !existingNames[name] {
			created, err := client.Tag.Create().SetName(name).Save(c.Request.Context())
			if err != nil {
				return nil, err
			}
			ids = append(ids, created.ID)
		}
	}
	return ids, nil
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
