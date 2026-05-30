package routes

import (
	"log/slog"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/equipment"
	"herbst-server/db/room"
	"herbst-server/dblog"
	"herbst-server/repository"
)

// revealConditions stores reveal conditions in memory (GitHub #12)
// In production, this would be in the database
var (
	revealConditions = make(map[int]map[string]any)
	revealMutex      sync.RWMutex
)

// RegisterEquipmentRoutes registers all equipment-related routes
func RegisterEquipmentRoutes(router *gin.Engine, repos *repository.Container, client *db.Client) {
	// Create a new equipment item
	// TODO: migrate to repos.Equipment.Create() when repo Create is implemented
	router.POST("/equipment", func(c *gin.Context) {
		var req struct {
			Name        string `json:"name" binding:"required"`
			Description string `json:"description"`
			Slot        string `json:"slot"`
			Level       int    `json:"level"`
			Weight      int    `json:"weight"`
			IsEquipped  bool   `json:"isEquipped"`
			// Item system fields (GitHub #89)
			IsImmovable bool `json:"isImmovable"`
			Color       string `json:"color"`
			IsVisible   bool   `json:"isVisible"`
			ItemType    string `json:"itemType"`
			RoomID      int    `json:"roomId"`
			// Owner system
			OwnerID *int `json:"ownerId"`
			// Consumable effects
			Healing int    `json:"healing"`
			Effect  string `json:"effect"`
			// Hidden items (GitHub #12 - Look System)
			RevealCondition map[string]any `json:"revealCondition"`
			// Corpse rotting (GitHub #22)
			ExpiresAt *time.Time `json:"expiresAt,omitempty"`
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
			slog.Warn("bad request", slog.String("service", "equipment"), slog.String("error", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Set defaults
		if req.ItemType == "" {
			req.ItemType = "misc"
		}
		if req.Slot == "" {
			req.Slot = "none"
		}

		builder := client.Equipment.
			Create().
			SetName(req.Name).
			SetDescription(req.Description).
			SetSlot(req.Slot).
			SetLevel(req.Level).
			SetWeight(req.Weight).
			SetIsEquipped(req.IsEquipped).
			SetIsImmovable(req.IsImmovable).
			SetIsVisible(req.IsVisible).
			SetItemType(req.ItemType).
			SetHealing(req.Healing).
			SetEffect(req.Effect)

		if req.Color != "" {
			builder.SetColor(req.Color)
		}

		// Set room if provided
		if req.RoomID > 0 {
			builder.SetRoomID(req.RoomID)
		}

		// Set owner if provided
		if req.OwnerID != nil {
			builder.SetOwnerId(*req.OwnerID)
		}

		// Set expiry time for corpses and other transient items (GitHub #22)
		if req.ExpiresAt != nil {
			builder.SetExpiresAt(*req.ExpiresAt)
		}

		// Set combat fields (EQUIP-002)
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

		eq, err := builder.Save(c.Request.Context())
		if err != nil {
			dblog.Error("failed to create equipment", err, slog.String("service", "equipment"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Store reveal condition in memory (GitHub #12)
		if req.RevealCondition != nil && len(req.RevealCondition) > 0 {
			revealMutex.Lock()
			revealConditions[eq.ID] = req.RevealCondition
			revealMutex.Unlock()
		}

		slog.Info("equipment created", slog.String("service", "equipment"), slog.Int("equipment_id", eq.ID))
		c.JSON(http.StatusCreated, eq)
	})

	// Get all equipment
	// TODO: migrate to repo methods when filtered list is supported
	router.GET("/equipment", func(c *gin.Context) {
		query := client.Equipment.Query()

		// Filter by room if roomId query param is provided
		if roomIDStr := c.Query("roomId"); roomIDStr != "" {
			roomID, err := strconv.Atoi(roomIDStr)
			if err == nil {
				query = query.Where(equipment.HasRoomWith(room.IDEQ(roomID)))
			}
		}

		// Filter by owner if ownerId query param is provided
		if ownerID := c.Query("ownerId"); ownerID != "" {
			id, err := strconv.Atoi(ownerID)
			if err == nil {
				query = query.Where(equipment.OwnerIdEQ(id))
			}
		}

		// Filter by type if type query param is provided
		if itemType := c.Query("type"); itemType != "" {
			query = query.Where(equipment.ItemTypeEQ(itemType))
		}

		items, err := query.All(c.Request.Context())
		if err != nil {
			dblog.Error("failed to list equipment", err, slog.String("service", "equipment"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Add reveal conditions to response
		type EqWithReveal struct {
			ID              int            `json:"id"`
			Name            string         `json:"name"`
			Description     string         `json:"description"`
			Slot            string         `json:"slot"`
			Level           int            `json:"level"`
			Weight          int            `json:"weight"`
			IsEquipped      bool           `json:"isEquipped"`
			IsImmovable     bool           `json:"isImmovable"`
			Color           string         `json:"color"`
			IsVisible       bool           `json:"isVisible"`
			ItemType        string         `json:"itemType"`
			OwnerID         *int           `json:"ownerId"`
			Healing         int            `json:"healing"`
			Effect          string         `json:"effect"`
			RevealCondition map[string]any `json:"revealCondition,omitempty"`
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

		result := make([]EqWithReveal, len(items))
		revealMutex.RLock()
		for i, item := range items {
			result[i] = EqWithReveal{
				ID:          item.ID,
				Name:        item.Name,
				Description: item.Description,
				Slot:        item.Slot,
				Level:       item.Level,
				Weight:      item.Weight,
				IsEquipped:  item.IsEquipped,
				IsImmovable: item.IsImmovable,
				Color:       item.Color,
				IsVisible:   item.IsVisible,
				ItemType:    item.ItemType,
				OwnerID:     item.OwnerId,
				Healing:     item.Healing,
				Effect:      item.Effect,
				ArmorRating:           item.ArmorRating,
				ArmorType:             item.ArmorType,
				Stats:                 item.Stats,
				DamageDiceCount:       item.DamageDiceCount,
				DamageDiceSides:       item.DamageDiceSides,
				DamageBonus:           item.DamageBonus,
				DamageType:            item.DamageType,
				WeaponType:            item.WeaponType,
				IsTwoHanded:           item.IsTwoHanded,
				Rarity:                item.Rarity,
			}
			if cond, ok := revealConditions[item.ID]; ok {
				result[i].RevealCondition = cond
			}
		}
		revealMutex.RUnlock()

		c.JSON(http.StatusOK, result)
	})

	// Get equipment in a room
	// Query param "includeHidden" can be set to "true" to include hidden items
	router.GET("/rooms/:id/equipment", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			slog.Warn("bad request", slog.String("service", "equipment"), slog.String("reason", "invalid room id"))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
			return
		}

		includeHidden := c.Query("includeHidden") == "true"

		allItems, err := repos.Equipment.ListByRoom(c.Request.Context(), id)
		if err != nil {
			dblog.Error("failed to list room equipment", err, slog.String("service", "equipment"), slog.Int("room_id", id))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// By default, only show visible items (GitHub #12)
		items := allItems
		if !includeHidden {
			items = make([]*db.Equipment, 0)
			for _, item := range allItems {
				if item.IsVisible {
					items = append(items, item)
				}
			}
		}

		// Add reveal conditions to response
		type EqWithReveal struct {
			ID              int            `json:"id"`
			Name            string         `json:"name"`
			Description     string         `json:"description"`
			Slot            string         `json:"slot"`
			Level           int            `json:"level"`
			Weight          int            `json:"weight"`
			IsEquipped      bool           `json:"isEquipped"`
			IsImmovable     bool           `json:"isImmovable"`
			Color           string         `json:"color"`
			IsVisible       bool           `json:"isVisible"`
			ItemType        string         `json:"itemType"`
			RevealCondition map[string]any `json:"revealCondition,omitempty"`
		}

		result := make([]EqWithReveal, len(items))
		revealMutex.RLock()
		for i, item := range items {
			result[i] = EqWithReveal{
				ID:          item.ID,
				Name:        item.Name,
				Description: item.Description,
				Slot:        item.Slot,
				Level:       item.Level,
				Weight:      item.Weight,
				IsEquipped:  item.IsEquipped,
				IsImmovable: item.IsImmovable,
				Color:       item.Color,
				IsVisible:   item.IsVisible,
				ItemType:    item.ItemType,
			}
			if cond, ok := revealConditions[item.ID]; ok {
				result[i].RevealCondition = cond
			}
		}
		revealMutex.RUnlock()

		c.JSON(http.StatusOK, result)
	})

	// Get a single equipment item by ID
	router.GET("/equipment/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			slog.Warn("bad request", slog.String("service", "equipment"), slog.String("reason", "invalid equipment id"))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid equipment ID"})
			return
		}

		eq, err := repos.Equipment.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Equipment not found"})
			return
		}

		// Add reveal condition to response
		type EqWithReveal struct {
			ID              int            `json:"id"`
			Name            string         `json:"name"`
			Description     string         `json:"description"`
			Slot            string         `json:"slot"`
			Level           int            `json:"level"`
			Weight          int            `json:"weight"`
			IsEquipped      bool           `json:"isEquipped"`
			IsImmovable     bool           `json:"isImmovable"`
			Color           string         `json:"color"`
			IsVisible       bool           `json:"isVisible"`
			ItemType        string         `json:"itemType"`
			OwnerID         *int           `json:"ownerId,omitempty"`
			RevealCondition map[string]any `json:"revealCondition,omitempty"`
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

		result := EqWithReveal{
			ID:          eq.ID,
			Name:        eq.Name,
			Description: eq.Description,
			Slot:        eq.Slot,
			Level:       eq.Level,
			Weight:      eq.Weight,
			IsEquipped:  eq.IsEquipped,
			IsImmovable: eq.IsImmovable,
			Color:       eq.Color,
			IsVisible:   eq.IsVisible,
			ItemType:    eq.ItemType,
			OwnerID:     eq.OwnerId,
			ArmorRating:           eq.ArmorRating,
			ArmorType:             eq.ArmorType,
			Stats:                 eq.Stats,
			Rarity:                eq.Rarity,
			SkillRequirement:      eq.SkillRequirement,
			SkillRequirementLevel: eq.SkillRequirementLevel,
			DamageDiceCount:       eq.DamageDiceCount,
			DamageDiceSides:       eq.DamageDiceSides,
			DamageBonus:           eq.DamageBonus,
			DamageType:            eq.DamageType,
			WeaponType:            eq.WeaponType,
			IsTwoHanded:           eq.IsTwoHanded,
		}

		revealMutex.RLock()
		if cond, ok := revealConditions[eq.ID]; ok {
			result.RevealCondition = cond
		}
		revealMutex.RUnlock()

		c.JSON(http.StatusOK, result)
	})

	// Update equipment by ID
	// TODO: migrate to repos.Equipment.Update() when EquipmentUpdates covers all fields
	router.PUT("/equipment/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			slog.Warn("bad request", slog.String("service", "equipment"), slog.String("reason", "invalid equipment id"))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid equipment ID"})
			return
		}

		var req struct {
			Name            string         `json:"name"`
			Description     string         `json:"description"`
			Slot            string         `json:"slot"`
			Level           *int           `json:"level"`
			Weight          *int           `json:"weight"`
			IsEquipped      *bool          `json:"isEquipped"`
			IsImmovable     *bool          `json:"isImmovable"`
			Color           string         `json:"color"`
			IsVisible       *bool          `json:"isVisible"`
			ItemType        string         `json:"itemType"`
			RoomID          *int           `json:"roomId"`
			OwnerID         *int           `json:"ownerId"`
			RevealCondition map[string]any `json:"revealCondition"`
			// Corpse rotting (GitHub #22)
			ExpiresAt *time.Time `json:"expiresAt,omitempty"`
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
			slog.Warn("bad request", slog.String("service", "equipment"), slog.String("error", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updater := client.Equipment.UpdateOneID(id)

		if req.Name != "" {
			updater.SetName(req.Name)
		}
		if req.Description != "" {
			updater.SetDescription(req.Description)
		}
		if req.Slot != "" {
			updater.SetSlot(req.Slot)
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
		if req.Color != "" {
			updater.SetColor(req.Color)
		}
		if req.IsVisible != nil {
			updater.SetIsVisible(*req.IsVisible)
		}
		if req.ItemType != "" {
			updater.SetItemType(req.ItemType)
		}
		if req.RoomID != nil {
			if *req.RoomID == 0 {
				updater.ClearRoom()
			} else {
				updater.SetRoomID(*req.RoomID)
			}
		}
		if req.OwnerID != nil {
			updater.SetOwnerId(*req.OwnerID)
		}

		// Update expiry time (GitHub #22)
		if req.ExpiresAt != nil {
			updater.SetExpiresAt(*req.ExpiresAt)
		}

		// Update combat fields (EQUIP-002)
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

		// Update reveal condition in memory (GitHub #12)
		if req.RevealCondition != nil {
			revealMutex.Lock()
			if len(req.RevealCondition) > 0 {
				revealConditions[id] = req.RevealCondition
			} else {
				delete(revealConditions, id)
			}
			revealMutex.Unlock()
		}

		eq, err := updater.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Equipment not found"})
			return
		}

		slog.Info("equipment updated", slog.String("service", "equipment"), slog.Int("equipment_id", eq.ID))
		c.JSON(http.StatusOK, eq)
	})

	// Delete equipment by ID
	router.DELETE("/equipment/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			slog.Warn("bad request", slog.String("service", "equipment"), slog.String("reason", "invalid equipment id"))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid equipment ID"})
			return
		}

		err = repos.Equipment.Delete(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Equipment not found"})
			return
		}

		// Clean up reveal condition
		revealMutex.Lock()
		delete(revealConditions, id)
		revealMutex.Unlock()

		slog.Info("equipment deleted", slog.String("service", "equipment"), slog.Int("equipment_id", id))
		c.JSON(http.StatusNoContent, nil)
	})

	// Examine equipment endpoint (look-10)
	router.GET("/equipment/:id/examine", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			slog.Warn("bad request", slog.String("service", "equipment"), slog.String("reason", "invalid equipment id"))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid equipment ID"})
			return
		}

		eq, err := repos.Equipment.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Equipment not found"})
			return
		}

		// Build examine response
		examineLevel := 1 // Default examine level
		if level := c.Query("examineLevel"); level != "" {
			if lvl, err := strconv.Atoi(level); err == nil {
				examineLevel = lvl
			}
		}

		// Build examine response with hidden details logic
		hiddenDetails := []map[string]interface{}{}

		// For immovable items like fountain, show hidden details at higher examine levels
		if eq.IsImmovable && examineLevel >= 50 {
			hiddenDetails = append(hiddenDetails, map[string]interface{}{
				"text":         "The item appears worn with age",
				"revealed":     true,
				"reveal_level": 50,
			})
		}

		// Return examine response
		c.JSON(http.StatusOK, gin.H{
			"id":              eq.ID,
			"name":            eq.Name,
			"description":     eq.Description,
			"type":            eq.ItemType,
			"is_immovable":    eq.IsImmovable,
			"is_visible":      eq.IsVisible,
			"color":           eq.Color,
			"examine_level":   examineLevel,
			"hidden_details":  hiddenDetails,
			"examine_xp":      1, // Grant XP for examining
		})
	})

	// Reveal hidden item (GitHub #12 - Look System)
	router.POST("/equipment/:id/reveal", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			slog.Warn("bad request", slog.String("service", "equipment"), slog.String("reason", "invalid equipment id"))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid equipment ID"})
			return
		}

		var req struct {
			RevealType string `json:"revealType"` // "examine", "perception_check", "use_item", "event"
			Target     string `json:"target"`
			SkillLevel int    `json:"skillLevel"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			slog.Warn("bad request", slog.String("service", "equipment"), slog.String("error", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get the equipment item
		eq, err := repos.Equipment.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Equipment not found"})
			return
		}

		// Check if item is already visible
		if eq.IsVisible {
			c.JSON(http.StatusOK, gin.H{"message": "Item is already visible", "item": eq})
			return
		}

		// Get reveal condition from memory
		revealMutex.RLock()
		revealCond, exists := revealConditions[id]
		revealMutex.RUnlock()

		if !exists || revealCond == nil {
			slog.Warn("bad request", slog.String("service", "equipment"), slog.String("reason", "item has no reveal condition"))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Item has no reveal condition"})
			return
		}

		// Check if the reveal type matches
		condType, _ := revealCond["type"].(string)
		if condType != req.RevealType {
			slog.Warn("bad request", slog.String("service", "equipment"), slog.String("reason", "invalid reveal type"))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reveal type for this item"})
			return
		}

		// Check target if required
		if target, ok := revealCond["target"].(string); ok && target != "" {
			if req.Target != target {
				slog.Warn("bad request", slog.String("service", "equipment"), slog.String("reason", "invalid target"))
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target"})
				return
			}
		}

		// Check skill level if required
		if minLevel, ok := revealCond["minLevel"].(float64); ok && minLevel > 0 {
			if req.SkillLevel < int(minLevel) {
				c.JSON(http.StatusForbidden, gin.H{"error": "Skill level too low"})
				return
			}
		}

		// All checks passed - reveal the item
		// TODO: migrate to repos.Equipment.Update() when EquipmentUpdates covers IsVisible
		visible := true
		updated, err := client.Equipment.UpdateOneID(id).
			SetIsVisible(visible).
			Save(c.Request.Context())
		if err != nil {
			dblog.Error("failed to reveal equipment", err, slog.String("service", "equipment"), slog.Int("equipment_id", id))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Remove the reveal condition (it's been used)
		revealMutex.Lock()
		delete(revealConditions, id)
		revealMutex.Unlock()

		slog.Info("equipment revealed", slog.String("service", "equipment"), slog.Int("equipment_id", updated.ID))
		c.JSON(http.StatusOK, updated)
	})

	// Get reveal condition for an item (GitHub #12)
	router.GET("/equipment/:id/reveal", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			slog.Warn("bad request", slog.String("service", "equipment"), slog.String("reason", "invalid equipment id"))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid equipment ID"})
			return
		}

		revealMutex.RLock()
		cond, exists := revealConditions[id]
		revealMutex.RUnlock()

		if !exists || cond == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No reveal condition found"})
			return
		}

		c.JSON(http.StatusOK, cond)
	})
}
