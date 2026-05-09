package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/equipment"
	"herbst-server/middleware"
)

// RegisterEquipmentTemplateRoutes registers REST endpoints for equipment templates.
func RegisterEquipmentTemplateRoutes(r *gin.Engine, client *db.Client) {
	g := r.Group("/api")
	g.Use(middleware.AuthMiddleware())
	g.Use(middleware.AdminMiddleware())
	{
		g.GET("/equipment-templates", listEquipmentTemplates(client))
		g.GET("/equipment-templates/:id", getEquipmentTemplate(client))
		g.POST("/equipment-templates", createEquipmentTemplate(client))
		g.PUT("/equipment-templates/:id", updateEquipmentTemplate(client))
		g.DELETE("/equipment-templates/:id", deleteEquipmentTemplate(client))
	}
}

func listEquipmentTemplates(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		templates, err := client.EquipmentTemplate.Query().All(c.Request.Context())
		if err != nil {
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

func getEquipmentTemplate(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template id"})
			return
		}
		t, err := client.EquipmentTemplate.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "template not found: " + id})
			return
		}
		c.JSON(http.StatusOK, templateToMap(t))
	}
}


func createEquipmentTemplate(client *db.Client) gin.HandlerFunc {
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
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		builder := client.EquipmentTemplate.Create().
			SetName(req.Name).
			SetDescription(req.Description).
			SetSlot(req.Slot).
			SetLevel(req.Level).
			SetWeight(req.Weight).
			SetItemType(req.ItemType).
			SetColor(req.Color).
			SetEffectType(req.EffectType).
			SetEffectValue(req.EffectValue).
			SetEffectDuration(req.EffectDuration).
			SetContainerCapacity(req.ContainerCapacity).
			SetKeyItemID(req.KeyItemID).
			SetRevealCondition(req.RevealCondition).
			SetArmorRating(req.ArmorRating).
			SetArmorType(req.ArmorType).
			SetRarity(req.Rarity).
			SetSkillRequirement(req.SkillRequirement).
			SetSkillRequirementLevel(req.SkillRequirementLvl).
			SetDamageDiceCount(req.DamageDiceCount).
			SetDamageDiceSides(req.DamageDiceSides).
			SetDamageBonus(req.DamageBonus).
			SetDamageType(req.DamageType).
			SetWeaponType(req.WeaponType)

		if req.Stats != nil {
			builder.SetStats(req.Stats)
		}
		if req.IsVisible != nil {
			builder.SetIsVisible(*req.IsVisible)
		}
		if req.IsImmovable != nil {
			builder.SetIsImmovable(*req.IsImmovable)
		}
		if req.IsContainer != nil {
			builder.SetIsContainer(*req.IsContainer)
		}
		if req.IsLocked != nil {
			builder.SetIsLocked(*req.IsLocked)
		}
		if req.IsTwoHanded != nil {
			builder.SetIsTwoHanded(*req.IsTwoHanded)
		}

		t, err := builder.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, templateToMap(t))
	}
}

func updateEquipmentTemplate(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
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
			Stats               map[string]int `json:"stats"`
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
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		builder := client.EquipmentTemplate.UpdateOneID(id)
		if req.Name != nil {
			builder.SetName(*req.Name)
		}
		if req.Description != nil {
			builder.SetDescription(*req.Description)
		}
		if req.Slot != nil {
			builder.SetSlot(*req.Slot)
		}
		if req.Level != nil {
			builder.SetLevel(*req.Level)
		}
		if req.Weight != nil {
			builder.SetWeight(*req.Weight)
		}
		if req.ItemType != nil {
			builder.SetItemType(*req.ItemType)
		}
		if req.Stats != nil {
			builder.SetStats(req.Stats)
		}
		if req.Color != nil {
			builder.SetColor(*req.Color)
		}
		if req.IsVisible != nil {
			builder.SetIsVisible(*req.IsVisible)
		}
		if req.IsImmovable != nil {
			builder.SetIsImmovable(*req.IsImmovable)
		}
		if req.EffectType != nil {
			builder.SetEffectType(*req.EffectType)
		}
		if req.EffectValue != nil {
			builder.SetEffectValue(*req.EffectValue)
		}
		if req.EffectDuration != nil {
			builder.SetEffectDuration(*req.EffectDuration)
		}
		if req.IsContainer != nil {
			builder.SetIsContainer(*req.IsContainer)
		}
		if req.ContainerCapacity != nil {
			builder.SetContainerCapacity(*req.ContainerCapacity)
		}
		if req.IsLocked != nil {
			builder.SetIsLocked(*req.IsLocked)
		}
		if req.KeyItemID != nil {
			builder.SetKeyItemID(*req.KeyItemID)
		}
		if req.RevealCondition != nil {
			builder.SetRevealCondition(*req.RevealCondition)
		}
		if req.ArmorRating != nil {
			builder.SetArmorRating(*req.ArmorRating)
		}
		if req.ArmorType != nil {
			builder.SetArmorType(*req.ArmorType)
		}
		if req.Rarity != nil {
			builder.SetRarity(*req.Rarity)
		}
		if req.SkillRequirement != nil {
			builder.SetSkillRequirement(*req.SkillRequirement)
		}
		if req.SkillRequirementLvl != nil {
			builder.SetSkillRequirementLevel(*req.SkillRequirementLvl)
		}
		if req.DamageDiceCount != nil {
			builder.SetDamageDiceCount(*req.DamageDiceCount)
		}
		if req.DamageDiceSides != nil {
			builder.SetDamageDiceSides(*req.DamageDiceSides)
		}
		if req.DamageBonus != nil {
			builder.SetDamageBonus(*req.DamageBonus)
		}
		if req.DamageType != nil {
			builder.SetDamageType(*req.DamageType)
		}
		if req.WeaponType != nil {
			builder.SetWeaponType(*req.WeaponType)
		}
		if req.IsTwoHanded != nil {
			builder.SetIsTwoHanded(*req.IsTwoHanded)
		}

		t, err := builder.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
			return
		}

		c.JSON(http.StatusOK, templateToMap(t))
	}
}

func deleteEquipmentTemplate(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template id"})
			return
		}

		// Check if any instances reference this template
		count, err := client.Equipment.Query().
			Where(equipment.EquipmentTemplateIDEQ(id)).
			Count(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "cannot delete template with existing instances"})
			return
		}

		err = client.EquipmentTemplate.DeleteOneID(id).Exec(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"deleted": id})
	}
}

func templateToMap(t *db.EquipmentTemplate) gin.H {
	return gin.H{
		"id":                      t.ID,
		"name":                    t.Name,
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