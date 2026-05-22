package routes

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/equipment"
	"herbst-server/db/room"
	"herbst-server/middleware"
	"herbst-server/repository"
)

// RegisterItemInstanceRoutes registers REST endpoints for item instances.
// Item instances are Equipment rows with equipment_template_id set.
func RegisterItemInstanceRoutes(r *gin.Engine, repos *repository.Container, client *db.Client) {
	g := r.Group("/api")
	g.Use(middleware.AuthMiddleware(nil))
	g.Use(middleware.AdminMiddleware())
	g.Use(middleware.WorldAccessMiddleware())
	{
		g.GET("/item-instances", listItemInstances(repos, client))
		g.POST("/item-instances", createItemInstance(repos, client))
		g.GET("/item-instances/:id", getItemInstance(repos, client))
		g.PUT("/item-instances/:id", updateItemInstance(repos, client))
		g.DELETE("/item-instances/:id", deleteItemInstance(repos, client))
	}
}

// ─── JSON views ─────────────────────────────────────────────────────────────

// itemInstanceView is the JSON shape returned by the API.
type itemInstanceView struct {
	ID                  int    `json:"id"`
	WorldID             string `json:"world_id"`
	Name                string `json:"name"`
	Description         string `json:"description"`
	Slot                string `json:"slot"`
	Level               int    `json:"level"`
	Weight              int    `json:"weight"`
	IsEquipped          bool   `json:"isEquipped"`
	IsImmovable         bool   `json:"isImmovable"`
	Color               string `json:"color"`
	IsVisible           bool   `json:"isVisible"`
	ItemType            string `json:"itemType"`
	EquipmentTemplateID int    `json:"equipment_template_id"`
	OwnerID             *int   `json:"ownerId,omitempty"`
	RoomID              int    `json:"roomId"`
	EffectType          string `json:"effect_type"`
	EffectValue         int    `json:"effect_value"`
	EffectDuration      int    `json:"effect_duration"`
	Healing             int    `json:"healing"`
	Effect              string `json:"effect"`
	IsContainer         bool   `json:"isContainer"`
	ContainerCapacity   int    `json:"containerCapacity"`
	IsLocked            bool   `json:"isLocked"`
	KeyItemID           string `json:"keyItemID"`
	ContainedItems      string `json:"containedItems"`
	RevealCondition     string `json:"revealCondition"`
	// Combat fields (EQUIP-002)
	ArmorRating           int            `json:"armor_rating"`
	ArmorType             string         `json:"armor_type"`
	Stats                 map[string]int `json:"stats,omitempty"`
	Rarity                string         `json:"rarity"`
	SkillRequirement      string         `json:"skill_requirement"`
	SkillRequirementLevel int            `json:"skill_requirement_level"`
	DamageDiceCount       int            `json:"damage_dice_count"`
	DamageDiceSides       int            `json:"damage_dice_sides"`
	DamageBonus           int            `json:"damage_bonus"`
	DamageType            string         `json:"damage_type"`
	WeaponType            string         `json:"weapon_type"`
	IsTwoHanded           bool           `json:"is_two_handed"`
}

func toItemInstanceView(e *db.Equipment) itemInstanceView {
	v := itemInstanceView{
		ID:                  e.ID,
		Name:                e.Name,
		Description:         e.Description,
		Slot:                e.Slot,
		Level:               e.Level,
		Weight:              e.Weight,
		IsEquipped:          e.IsEquipped,
		IsImmovable:         e.IsImmovable,
		Color:               e.Color,
		IsVisible:           e.IsVisible,
		ItemType:            e.ItemType,
		EquipmentTemplateID: e.EquipmentTemplateID,
		OwnerID:             e.OwnerId,
		EffectType:          e.EffectType,
		EffectValue:         e.EffectValue,
		EffectDuration:      e.EffectDuration,
		Healing:             e.Healing,
		Effect:              e.Effect,
		IsContainer:         e.IsContainer,
		ContainerCapacity:   e.ContainerCapacity,
		IsLocked:            e.IsLocked,
		KeyItemID:           e.KeyItemID,
		ContainedItems:      e.ContainedItems,
		RevealCondition:     e.RevealCondition,
		// Combat fields (EQUIP-002)
		ArmorRating:           e.ArmorRating,
		ArmorType:             e.ArmorType,
		Stats:                 e.Stats,
		Rarity:                e.Rarity,
		SkillRequirement:      e.SkillRequirement,
		SkillRequirementLevel: e.SkillRequirementLevel,
		DamageDiceCount:       e.DamageDiceCount,
		DamageDiceSides:       e.DamageDiceSides,
		DamageBonus:           e.DamageBonus,
		DamageType:            e.DamageType,
		WeaponType:            e.WeaponType,
		IsTwoHanded:           e.IsTwoHanded,
	}
	if r, err := e.QueryRoom().Only(context.TODO()); err == nil {
		v.RoomID = r.ID
		v.WorldID = r.WorldID
	}
	return v
}

// ─── Handlers ───────────────────────────────────────────────────────────────

// GET /api/item-instances?ownerId=X&templateId=X&type=X&world_id=X
// TODO: Add filtered list method to EquipmentRepo for complex queries
func listItemInstances(repos *repository.Container, client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get world_id for filtering via room
		worldID := c.Query("world_id")

		query := client.Equipment.Query().
			Where(equipment.EquipmentTemplateIDNotNil())

		// Filter by world_id via room
		if worldID != "" {
			query = query.Where(equipment.HasRoomWith(room.WorldIDEQ(worldID)))
		}

		// Optional filters
		if ownerIDStr := c.Query("ownerId"); ownerIDStr != "" {
			ownerID, err := strconv.Atoi(ownerIDStr)
			if err == nil {
				query = query.Where(equipment.OwnerIdEQ(ownerID))
			}
		}
		if templateIDStr := c.Query("templateId"); templateIDStr != "" {
			if tid, err := strconv.Atoi(templateIDStr); err == nil {
				query = query.Where(equipment.EquipmentTemplateIDEQ(tid))
			}
		}
		if itemType := c.Query("type"); itemType != "" {
			query = query.Where(equipment.ItemTypeEQ(itemType))
		}
		if roomIDStr := c.Query("roomId"); roomIDStr != "" {
			roomID, err := strconv.Atoi(roomIDStr)
			if err == nil {
				query = query.Where(equipment.HasRoomWith(room.IDEQ(roomID)))
			}
		}

		items, err := query.All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		result := make([]itemInstanceView, len(items))
		for i, item := range items {
			result[i] = toItemInstanceView(item)
		}

		c.JSON(http.StatusOK, result)
	}
}

// POST /api/item-instances — create instance from template or bare item
// TODO: Move creation logic to EquipmentRepo.Create with template expansion
func createItemInstance(repos *repository.Container, client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			WorldID            string `json:"world_id"`
			EquipmentTemplateID int    `json:"equipment_template_id"`
			Name                string `json:"name"`
			Description         string `json:"description"`
			Slot                string `json:"slot"`
			Level               int    `json:"level"`
			Weight              int    `json:"weight"`
			IsEquipped          bool   `json:"isEquipped"`
			IsImmovable         bool   `json:"isImmovable"`
			Color               string `json:"color"`
			IsVisible           bool   `json:"isVisible"`
			ItemType            string `json:"itemType"`
			OwnerID             *int   `json:"ownerId"`
			EffectType          string `json:"effect_type"`
			EffectValue         int    `json:"effect_value"`
			EffectDuration      int    `json:"effect_duration"`
			Healing             int    `json:"healing"`
			Effect              string `json:"effect"`
			IsContainer         bool   `json:"isContainer"`
			ContainerCapacity   int    `json:"containerCapacity"`
			IsLocked            bool   `json:"isLocked"`
			KeyItemID           string `json:"keyItemID"`
			ContainedItems      string `json:"containedItems"`
			RevealCondition     string         `json:"revealCondition"`
			RoomID              int            `json:"room_id"`
			// Combat fields (EQUIP-002)
			ArmorRating           int            `json:"armor_rating"`
			ArmorType             string         `json:"armor_type"`
			Stats                 map[string]int `json:"stats"`
			Rarity                string         `json:"rarity"`
			SkillRequirement      string         `json:"skill_requirement"`
			SkillRequirementLevel int            `json:"skill_requirement_level"`
			DamageDiceCount       int            `json:"damage_dice_count"`
			DamageDiceSides       int            `json:"damage_dice_sides"`
			DamageBonus           int            `json:"damage_bonus"`
			DamageType            string         `json:"damage_type"`
			WeaponType            string         `json:"weapon_type"`
			IsTwoHanded           bool           `json:"is_two_handed"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		builder := client.Equipment.Create()

		// If template provided by slug, look it up
		if req.EquipmentTemplateID > 0 {
			tmpl, err := repos.EquipmentTemplate.Get(c.Request.Context(), req.EquipmentTemplateID)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "equipment template not found"})
				return
			}
			builder.SetEquipmentTemplateID(tmpl.ID)
			if req.Name == "" {
				builder.SetName(tmpl.Name)
			} else {
				builder.SetName(req.Name)
			}
			if req.Description == "" {
				builder.SetDescription(tmpl.Description)
			} else {
				builder.SetDescription(req.Description)
			}
			if req.Slot == "" {
				builder.SetSlot(tmpl.Slot)
			} else {
				builder.SetSlot(req.Slot)
			}
			if req.Level == 0 {
				builder.SetLevel(tmpl.Level)
			} else {
				builder.SetLevel(req.Level)
			}
			if req.Weight == 0 {
				builder.SetWeight(tmpl.Weight)
			} else {
				builder.SetWeight(req.Weight)
			}
			if req.ItemType == "" {
				builder.SetItemType(tmpl.ItemType)
			} else {
				builder.SetItemType(req.ItemType)
			}
			if req.Color == "" {
				builder.SetColor(tmpl.Color)
			} else {
				builder.SetColor(req.Color)
			}
			builder.SetIsImmovable(tmpl.IsImmovable)
			builder.SetIsVisible(tmpl.IsVisible)
			builder.SetEffectType(tmpl.EffectType)
			builder.SetEffectValue(tmpl.EffectValue)
			builder.SetEffectDuration(tmpl.EffectDuration)
			builder.SetIsContainer(tmpl.IsContainer)
			builder.SetContainerCapacity(tmpl.ContainerCapacity)
			builder.SetIsLocked(tmpl.IsLocked)
			// Copy combat fields from template (EQUIP-002)
			builder.SetArmorRating(tmpl.ArmorRating)
			builder.SetArmorType(tmpl.ArmorType)
			builder.SetStats(tmpl.Stats)
			builder.SetRarity(tmpl.Rarity)
			builder.SetSkillRequirement(tmpl.SkillRequirement)
			builder.SetSkillRequirementLevel(tmpl.SkillRequirementLevel)
			builder.SetDamageDiceCount(tmpl.DamageDiceCount)
			builder.SetDamageDiceSides(tmpl.DamageDiceSides)
			builder.SetDamageBonus(tmpl.DamageBonus)
			builder.SetDamageType(tmpl.DamageType)
			builder.SetWeaponType(tmpl.WeaponType)
			builder.SetIsTwoHanded(tmpl.IsTwoHanded)
		} else {
			// Bare item — name is required
			if req.Name == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "name is required when no template is provided"})
				return
			}
			builder.SetName(req.Name)
			builder.SetDescription(req.Description)
			builder.SetSlot(req.Slot)
			builder.SetLevel(req.Level)
			builder.SetWeight(req.Weight)
			builder.SetItemType(req.ItemType)
			builder.SetColor(req.Color)
			builder.SetIsImmovable(req.IsImmovable)
			builder.SetIsVisible(req.IsVisible)
			builder.SetEffectType(req.EffectType)
			builder.SetEffectValue(req.EffectValue)
			builder.SetEffectDuration(req.EffectDuration)
			builder.SetIsContainer(req.IsContainer)
			builder.SetContainerCapacity(req.ContainerCapacity)
			builder.SetIsLocked(req.IsLocked)
			// Set combat fields for bare items (EQUIP-002)
			builder.SetArmorRating(req.ArmorRating)
			if req.ArmorType != "" {
				builder.SetArmorType(req.ArmorType)
			}
			if req.Stats != nil {
				builder.SetStats(req.Stats)
			}
			if req.Rarity != "" {
				builder.SetRarity(req.Rarity)
			}
			if req.SkillRequirement != "" {
				builder.SetSkillRequirement(req.SkillRequirement)
			}
			builder.SetSkillRequirementLevel(req.SkillRequirementLevel)
			builder.SetDamageDiceCount(req.DamageDiceCount)
			builder.SetDamageDiceSides(req.DamageDiceSides)
			builder.SetDamageBonus(req.DamageBonus)
			if req.DamageType != "" {
				builder.SetDamageType(req.DamageType)
			}
			if req.WeaponType != "" {
				builder.SetWeaponType(req.WeaponType)
			}
			builder.SetIsTwoHanded(req.IsTwoHanded)
		}

		// Apply explicit overrides common to both paths
		builder.SetIsEquipped(req.IsEquipped)
		if req.OwnerID != nil {
			builder.SetOwnerId(*req.OwnerID)
		}
		if req.Healing != 0 {
			builder.SetHealing(req.Healing)
		}
		if req.Effect != "" {
			builder.SetEffect(req.Effect)
		}
		if req.KeyItemID != "" {
			builder.SetKeyItemID(req.KeyItemID)
		}
		if req.ContainedItems != "" {
			builder.SetContainedItems(req.ContainedItems)
		}
		if req.RevealCondition != "" {
			builder.SetRevealCondition(req.RevealCondition)
		}
		if req.RoomID > 0 {
			builder.SetRoomID(req.RoomID)
		}

		eq, err := builder.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, toItemInstanceView(eq))
	}
}

// GET /api/item-instances/:id
// TODO: Add GetWithFilters method to EquipmentRepo
func getItemInstance(repos *repository.Container, client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item instance id"})
			return
		}

		eq, err := client.Equipment.Query().
			Where(equipment.IDEQ(id), equipment.EquipmentTemplateIDNotNil()).
			Only(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "item instance not found"})
			return
		}

		c.JSON(http.StatusOK, toItemInstanceView(eq))
	}
}

// PUT /api/item-instances/:id — update instance fields
// TODO: Move to EquipmentRepo.Update once EquipmentUpdates supports all fields
func updateItemInstance(repos *repository.Container, client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item instance id"})
			return
		}

		var req struct {
			Name              *string `json:"name"`
			Description       *string `json:"description"`
			Slot              *string `json:"slot"`
			Level             *int    `json:"level"`
			Weight            *int    `json:"weight"`
			IsEquipped        *bool   `json:"isEquipped"`
			IsImmovable       *bool   `json:"isImmovable"`
			Color             *string `json:"color"`
			IsVisible         *bool   `json:"isVisible"`
			ItemType          *string `json:"itemType"`
			OwnerID           *int    `json:"ownerId"`
			EffectType        *string `json:"effect_type"`
			EffectValue       *int    `json:"effect_value"`
			EffectDuration    *int    `json:"effect_duration"`
			Healing           *int    `json:"healing"`
			Effect            *string `json:"effect"`
			IsContainer       *bool   `json:"isContainer"`
			ContainerCapacity *int    `json:"containerCapacity"`
			IsLocked          *bool   `json:"isLocked"`
			KeyItemID         *string `json:"keyItemID"`
			ContainedItems    *string `json:"containedItems"`
			RevealCondition   *string         `json:"revealCondition"`
			// Combat fields (EQUIP-002)
			ArmorRating           *int            `json:"armor_rating"`
			ArmorType             *string         `json:"armor_type"`
			Stats                 map[string]int `json:"stats"`
			Rarity                *string         `json:"rarity"`
			SkillRequirement      *string         `json:"skill_requirement"`
			SkillRequirementLevel *int            `json:"skill_requirement_level"`
			DamageDiceCount       *int            `json:"damage_dice_count"`
			DamageDiceSides       *int            `json:"damage_dice_sides"`
			DamageBonus           *int            `json:"damage_bonus"`
			DamageType            *string         `json:"damage_type"`
			WeaponType            *string         `json:"weapon_type"`
			IsTwoHanded           *bool           `json:"is_two_handed"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updater := client.Equipment.UpdateOneID(id)
		if req.Name != nil {
			updater.SetName(*req.Name)
		}
		if req.Description != nil {
			updater.SetDescription(*req.Description)
		}
		if req.Slot != nil {
			updater.SetSlot(*req.Slot)
		}
		if req.Level != nil {
			updater.SetLevel(*req.Level)
		}
		if req.Weight != nil {
			updater.SetWeight(*req.Weight)
		}
		if req.IsEquipped != nil {
			updater.SetIsEquipped(*req.IsEquipped)
		}
		if req.IsImmovable != nil {
			updater.SetIsImmovable(*req.IsImmovable)
		}
		if req.Color != nil {
			updater.SetColor(*req.Color)
		}
		if req.IsVisible != nil {
			updater.SetIsVisible(*req.IsVisible)
		}
		if req.ItemType != nil {
			updater.SetItemType(*req.ItemType)
		}
		if req.OwnerID != nil {
			updater.SetOwnerId(*req.OwnerID)
		}
		if req.EffectType != nil {
			updater.SetEffectType(*req.EffectType)
		}
		if req.EffectValue != nil {
			updater.SetEffectValue(*req.EffectValue)
		}
		if req.EffectDuration != nil {
			updater.SetEffectDuration(*req.EffectDuration)
		}
		if req.Healing != nil {
			updater.SetHealing(*req.Healing)
		}
		if req.Effect != nil {
			updater.SetEffect(*req.Effect)
		}
		if req.IsContainer != nil {
			updater.SetIsContainer(*req.IsContainer)
		}
		if req.ContainerCapacity != nil {
			updater.SetContainerCapacity(*req.ContainerCapacity)
		}
		if req.IsLocked != nil {
			updater.SetIsLocked(*req.IsLocked)
		}
		if req.KeyItemID != nil {
			updater.SetKeyItemID(*req.KeyItemID)
		}
		if req.ContainedItems != nil {
			updater.SetContainedItems(*req.ContainedItems)
		}
		if req.RevealCondition != nil {
			updater.SetRevealCondition(*req.RevealCondition)
		}
		if req.ArmorRating != nil {
			updater.SetArmorRating(*req.ArmorRating)
		}
		if req.ArmorType != nil {
			updater.SetArmorType(*req.ArmorType)
		}
		if req.Stats != nil {
			updater.SetStats(req.Stats)
		}
		if req.Rarity != nil {
			updater.SetRarity(*req.Rarity)
		}
		if req.SkillRequirement != nil {
			updater.SetSkillRequirement(*req.SkillRequirement)
		}
		if req.SkillRequirementLevel != nil {
			updater.SetSkillRequirementLevel(*req.SkillRequirementLevel)
		}
		if req.DamageDiceCount != nil {
			updater.SetDamageDiceCount(*req.DamageDiceCount)
		}
		if req.DamageDiceSides != nil {
			updater.SetDamageDiceSides(*req.DamageDiceSides)
		}
		if req.DamageBonus != nil {
			updater.SetDamageBonus(*req.DamageBonus)
		}
		if req.DamageType != nil {
			updater.SetDamageType(*req.DamageType)
		}
		if req.WeaponType != nil {
			updater.SetWeaponType(*req.WeaponType)
		}
		if req.IsTwoHanded != nil {
			updater.SetIsTwoHanded(*req.IsTwoHanded)
		}

		updated, err := updater.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "item instance not found"})
			return
		}

		c.JSON(http.StatusOK, toItemInstanceView(updated))
	}
}

// DELETE /api/item-instances/:id — hard delete from DB
func deleteItemInstance(repos *repository.Container, client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item instance id"})
			return
		}

		// TODO: Add Delete method to EquipmentRepo
		err = client.Equipment.DeleteOneID(id).Exec(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "item instance not found"})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

// GET /api/equipment-templates — list all templates for admin spawn UI