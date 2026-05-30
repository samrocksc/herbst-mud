package routes

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/dblog"
	"herbst-server/middleware"
	"herbst-server/repository"
)

// RegisterEquipmentTemplateRoutes registers REST endpoints for equipment templates.
func RegisterEquipmentTemplateRoutes(r *gin.Engine, repos *repository.Container) {
	g := r.Group("/api")
	g.Use(middleware.AuthMiddleware(nil))
	g.Use(middleware.AdminMiddleware())
	g.Use(middleware.WorldAccessMiddleware())
	{
		g.GET("/equipment-templates", listEquipmentTemplates(repos))
		g.GET("/equipment-templates/:id", getEquipmentTemplate(repos))
		g.POST("/equipment-templates", createEquipmentTemplate(repos))
		g.PUT("/equipment-templates/:id", updateEquipmentTemplate(repos))
		g.DELETE("/equipment-templates/:id", deleteEquipmentTemplate(repos))
	}
}

func listEquipmentTemplates(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		worldID := c.Query("world_id")
		templates, err := repos.EquipmentTemplate.List(c.Request.Context(), worldID)
		if err != nil {
			dblog.Error("failed to list equipment templates", err, slog.String("service", "equipment"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result := make([]gin.H, len(templates))
		for i, t := range templates {
			result[i] = templateToMap(t)
		}
		c.JSON(http.StatusOK, result)
	}
}

func getEquipmentTemplate(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			slog.Warn("bad request", slog.String("service", "equipment"), slog.String("reason", "invalid template id"))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template id"})
			return
		}
		t, err := repos.EquipmentTemplate.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "template not found: " + idStr})
			return
		}
		c.JSON(http.StatusOK, templateToMap(t))
	}
}

func createEquipmentTemplate(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name                string         `json:"name" binding:"required"`
			Description         string         `json:"description"`
			Slot                string         `json:"slot"`
			Level               int            `json:"level"`
			Weight              int            `json:"weight"`
			ItemType            string         `json:"item_type"`
			Stats               map[string]int `json:"stats"`
			Color               string         `json:"color"`
			IsVisible           *bool          `json:"is_visible"`
			IsImmovable         *bool          `json:"is_immovable"`
			EffectType          string         `json:"effect_type"`
			EffectValue         int            `json:"effect_value"`
			EffectDuration      int            `json:"effect_duration"`
			IsContainer         *bool          `json:"is_container"`
			ContainerCapacity   int            `json:"container_capacity"`
			IsLocked            *bool          `json:"is_locked"`
			KeyItemID           string         `json:"key_item_id"`
			RevealCondition     string         `json:"reveal_condition"`
			ArmorRating         int            `json:"armor_rating"`
			ArmorType           string         `json:"armor_type"`
			Rarity              string         `json:"rarity"`
			SkillRequirement    string         `json:"skill_requirement"`
			SkillRequirementLvl int            `json:"skill_requirement_level"`
			DamageDiceCount     int            `json:"damage_dice_count"`
			DamageDiceSides     int            `json:"damage_dice_sides"`
			DamageBonus         int            `json:"damage_bonus"`
			DamageType          string         `json:"damage_type"`
			WeaponType          string         `json:"weapon_type"`
			IsTwoHanded         *bool          `json:"is_two_handed"`
			WorldID             string         `json:"world_id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			slog.Warn("bad request", slog.String("service", "equipment"), slog.String("error", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Auto-derive slug from name
		slug := strings.ToLower(req.Name)
		slug = strings.NewReplacer(" ", "_", "-", "_", "'", "", "\"", "").Replace(slug)

		// Inherit world from query param if not in body, default to first world
		worldID := req.WorldID
		if worldID == "" {
			worldID = c.Query("world_id")
		}
		if worldID == "" {
			worldID = "default"
		}

		// Check if a template with this slug already exists
		existing, err := repos.EquipmentTemplate.GetBySlug(c.Request.Context(), slug, worldID)
		if err == nil && existing != nil {
			c.JSON(http.StatusConflict, gin.H{
				"error":       "An item template with this name already exists",
				"existing":    templateToMap(existing),
				"suggestions": []string{
					"Choose a different name",
					"Edit the existing template instead",
					"Add a suffix like '_2' or ' (new)'",
				},
			})
			return
		}

		isVisible := false
		if req.IsVisible != nil {
			isVisible = *req.IsVisible
		}
		isImmovable := false
		if req.IsImmovable != nil {
			isImmovable = *req.IsImmovable
		}
		isContainer := false
		if req.IsContainer != nil {
			isContainer = *req.IsContainer
		}
		isLocked := false
		if req.IsLocked != nil {
			isLocked = *req.IsLocked
		}
		isTwoHanded := false
		if req.IsTwoHanded != nil {
			isTwoHanded = *req.IsTwoHanded
		}

		t, err := repos.EquipmentTemplate.Create(c.Request.Context(), repository.CreateEquipmentTemplateInput{
			Slug:                  slug,
			Name:                  req.Name,
			Description:           req.Description,
			Slot:                  req.Slot,
			Level:                 req.Level,
			Weight:                req.Weight,
			ItemType:              req.ItemType,
			Stats:                 req.Stats,
			Color:                 req.Color,
			IsVisible:             isVisible,
			IsImmovable:           isImmovable,
			EffectType:            req.EffectType,
			EffectValue:           req.EffectValue,
			EffectDuration:        req.EffectDuration,
			IsContainer:           isContainer,
			ContainerCapacity:     req.ContainerCapacity,
			IsLocked:              isLocked,
			KeyItemID:             req.KeyItemID,
			RevealCondition:       req.RevealCondition,
			ArmorRating:           req.ArmorRating,
			ArmorType:             req.ArmorType,
			Rarity:                req.Rarity,
			SkillRequirement:      req.SkillRequirement,
			SkillRequirementLevel: req.SkillRequirementLvl,
			DamageDiceCount:       req.DamageDiceCount,
			DamageDiceSides:       req.DamageDiceSides,
			DamageBonus:           req.DamageBonus,
			DamageType:            req.DamageType,
			WeaponType:            req.WeaponType,
			IsTwoHanded:           isTwoHanded,
			WorldID:               worldID,
		})
		if err != nil {
			dblog.Error("failed to create equipment template", err, slog.String("service", "equipment"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		slog.Info("equipment template created", slog.String("service", "equipment"), slog.Int("template_id", t.ID))
		c.JSON(http.StatusCreated, templateToMap(t))
	}
}

func updateEquipmentTemplate(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			slog.Warn("bad request", slog.String("service", "equipment"), slog.String("reason", "invalid template id"))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template id"})
			return
		}

		var req struct {
			Name                *string         `json:"name"`
			Description         *string         `json:"description"`
			Slot                *string         `json:"slot"`
			Level               *int            `json:"level"`
			Weight              *int            `json:"weight"`
			ItemType            *string         `json:"item_type"`
			Stats               map[string]int  `json:"stats"`
			Color               *string         `json:"color"`
			IsVisible           *bool           `json:"is_visible"`
			IsImmovable         *bool           `json:"is_immovable"`
			EffectType          *string         `json:"effect_type"`
			EffectValue         *int            `json:"effect_value"`
			EffectDuration      *int            `json:"effect_duration"`
			IsContainer         *bool           `json:"is_container"`
			ContainerCapacity   *int            `json:"container_capacity"`
			IsLocked            *bool           `json:"is_locked"`
			KeyItemID           *string         `json:"key_item_id"`
			RevealCondition     *string         `json:"reveal_condition"`
			ArmorRating         *int            `json:"armor_rating"`
			ArmorType           *string         `json:"armor_type"`
			Rarity              *string         `json:"rarity"`
			SkillRequirement    *string         `json:"skill_requirement"`
			SkillRequirementLvl *int            `json:"skill_requirement_level"`
			DamageDiceCount     *int            `json:"damage_dice_count"`
			DamageDiceSides     *int            `json:"damage_dice_sides"`
			DamageBonus         *int            `json:"damage_bonus"`
			DamageType          *string         `json:"damage_type"`
			WeaponType          *string         `json:"weapon_type"`
			IsTwoHanded         *bool           `json:"is_two_handed"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			slog.Warn("bad request", slog.String("service", "equipment"), slog.String("error", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		t, err := repos.EquipmentTemplate.Update(c.Request.Context(), id, repository.EquipmentTemplateUpdates{
			Name:                  req.Name,
			Description:           req.Description,
			Slot:                  req.Slot,
			Level:                 req.Level,
			Weight:                req.Weight,
			ItemType:              req.ItemType,
			Stats:                 req.Stats,
			Color:                 req.Color,
			IsVisible:             req.IsVisible,
			IsImmovable:           req.IsImmovable,
			EffectType:            req.EffectType,
			EffectValue:           req.EffectValue,
			EffectDuration:        req.EffectDuration,
			IsContainer:           req.IsContainer,
			ContainerCapacity:      req.ContainerCapacity,
			IsLocked:              req.IsLocked,
			KeyItemID:             req.KeyItemID,
			RevealCondition:       req.RevealCondition,
			ArmorRating:           req.ArmorRating,
			ArmorType:             req.ArmorType,
			Rarity:                req.Rarity,
			SkillRequirement:      req.SkillRequirement,
			SkillRequirementLevel: req.SkillRequirementLvl,
			DamageDiceCount:       req.DamageDiceCount,
			DamageDiceSides:       req.DamageDiceSides,
			DamageBonus:           req.DamageBonus,
			DamageType:            req.DamageType,
			WeaponType:            req.WeaponType,
			IsTwoHanded:           req.IsTwoHanded,
		})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
			return
		}

		slog.Info("equipment template updated", slog.String("service", "equipment"), slog.Int("template_id", t.ID))
		c.JSON(http.StatusOK, templateToMap(t))
	}
}

func deleteEquipmentTemplate(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			slog.Warn("bad request", slog.String("service", "equipment"), slog.String("reason", "invalid template id"))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template id"})
			return
		}

		// Check if any instances reference this template
		count, err := repos.Equipment.CountByTemplateID(c.Request.Context(), id)
		if err != nil {
			dblog.Error("failed to count equipment by template", err, slog.String("service", "equipment"), slog.Int("template_id", id))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "cannot delete template with existing instances"})
			return
		}

		if err := repos.EquipmentTemplate.Delete(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
			return
		}

		slog.Info("equipment template deleted", slog.String("service", "equipment"), slog.Int("template_id", id))
		c.JSON(http.StatusOK, gin.H{"deleted": id})
	}
}

func templateToMap(t *db.EquipmentTemplate) gin.H {
	return gin.H{
		"id":                      t.ID,
		"slug":                    t.Slug,
		"name":                    t.Name,
		"world_id":                t.WorldID,
		"description":             t.Description,
		"slot":                    t.Slot,
		"level":                   t.Level,
		"weight":                  t.Weight,
		"item_type":               t.ItemType,
		"stats":                   t.Stats,
		"color":                   t.Color,
		"is_visible":              t.IsVisible,
		"is_immovable":            t.IsImmovable,
		"effect_type":             t.EffectType,
		"effect_value":            t.EffectValue,
		"effect_duration":         t.EffectDuration,
		"is_container":            t.IsContainer,
		"container_capacity":      t.ContainerCapacity,
		"is_locked":               t.IsLocked,
		"key_item_id":             t.KeyItemID,
		"reveal_condition":        t.RevealCondition,
		"expires_at":              t.ExpiresAt,
		"armor_rating":            t.ArmorRating,
		"armor_type":              t.ArmorType,
		"rarity":                  t.Rarity,
		"skill_requirement":       t.SkillRequirement,
		"skill_requirement_level": t.SkillRequirementLevel,
		"damage_dice_count":       t.DamageDiceCount,
		"damage_dice_sides":       t.DamageDiceSides,
		"damage_bonus":            t.DamageBonus,
		"damage_type":             t.DamageType,
		"weapon_type":             t.WeaponType,
		"is_two_handed":           t.IsTwoHanded,
	}
}
