package routes

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/constants"
	"herbst-server/db"
	"herbst-server/dblog"
	"herbst-server/db/tag"
	"herbst-server/repository"
)

// createRace creates a new race in the specified world.
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
			WorldID         string   `json:"world_id" default:"1"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			slog.Warn("bad request", slog.String("service", "races"), slog.String("reason", "invalid request body"), slog.String("error", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
			return
		}

		if err := validateSlots(req.EquipmentSlots); err != nil {
			slog.Warn("bad request", slog.String("service", "races"), slog.String("reason", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check for duplicate name in this world
		existing, err := repos.Race.GetByName(c.Request.Context(), req.Name, req.WorldID)
		if err == nil && existing != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "a race with this name already exists in this world"})
			return
		}

		displayName := req.DisplayName
		if displayName == "" {
			displayName = req.Name
		}

		// Resolve tag names to IDs
		var tagIDs []int
		if len(req.Tags) > 0 {
			tagIDs, err = resolveTagIDs(c, client, req.Tags, req.WorldID)
			if err != nil {
				slog.Warn("bad request", slog.String("service", "races"), slog.String("reason", err.Error()))
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
			dblog.Error("failed to create race", err, slog.String("service", "races"), slog.String("world_id", req.WorldID))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		slog.Info("race created", slog.String("service", "races"), slog.String("race_name", req.Name), slog.String("world_id", req.WorldID))
		c.JSON(http.StatusCreated, raceToView(r))
	}
}

// resolveTagIDs resolves tag names to IDs, creating tags that don't exist yet.
// Tags created here are scoped to the given worldID so they appear in the
// admin Tags page for that world (instead of leaking into world 1).
// TODO: Move tag resolution to TagRepo
func resolveTagIDs(c *gin.Context, client *db.Client, names []string, worldID string) ([]int, error) {
	if len(names) == 0 {
		return nil, nil
	}
	if worldID == "" {
		worldID = "1"
	}
	existing, err := client.Tag.Query().Where(tag.NameIn(names...), tag.WorldID(worldID)).All(c.Request.Context())
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
			created, err := client.Tag.Create().SetName(name).SetWorldID(worldID).Save(c.Request.Context())
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
