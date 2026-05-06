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
)

// RegisterItemInstanceRoutes registers REST endpoints for item instances.
// Item instances are Equipment rows with equipment_template_id set.
func RegisterItemInstanceRoutes(r *gin.Engine, client *db.Client) {
	g := r.Group("/api")
	g.Use(middleware.AuthMiddleware())
	g.Use(middleware.AdminMiddleware())
	{
		g.GET("/equipment-templates", listEquipmentTemplates(client))
		g.GET("/equipment-templates/:id", getEquipmentTemplate(client))
		g.GET("/item-instances", listItemInstances(client))
		g.POST("/item-instances", createItemInstance(client))
		g.GET("/item-instances/:id", getItemInstance(client))
		g.PUT("/item-instances/:id", updateItemInstance(client))
		g.DELETE("/item-instances/:id", deleteItemInstance(client))
	}
}

// ─── JSON views ─────────────────────────────────────────────────────────────

// itemInstanceView is the JSON shape returned by the API.
type itemInstanceView struct {
	ID                  int    `json:"id"`
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
	EquipmentTemplateID string `json:"equipment_template_id"`
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
	}
	if r, err := e.QueryRoom().Only(context.TODO()); err == nil {
		v.RoomID = r.ID
	}
	return v
}

// ─── Handlers ───────────────────────────────────────────────────────────────

// GET /api/item-instances?ownerId=X&templateId=X&type=X
func listItemInstances(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := client.Equipment.Query().
			Where(equipment.EquipmentTemplateIDNEQ(""))

		// Optional filters
		if ownerIDStr := c.Query("ownerId"); ownerIDStr != "" {
			ownerID, err := strconv.Atoi(ownerIDStr)
			if err == nil {
				query = query.Where(equipment.OwnerIdEQ(ownerID))
			}
		}
		if templateID := c.Query("templateId"); templateID != "" {
			query = query.Where(equipment.EquipmentTemplateIDEQ(templateID))
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
func createItemInstance(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			EquipmentTemplateID string `json:"equipment_template_id"`
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
			RevealCondition     string `json:"revealCondition"`
			RoomID              int    `json:"room_id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		builder := client.Equipment.Create()

		// If template provided, auto-fill fields from template
		if req.EquipmentTemplateID != "" {
			tmpl, err := client.EquipmentTemplate.Get(c.Request.Context(), req.EquipmentTemplateID)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "equipment template not found: " + req.EquipmentTemplateID})
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
func getItemInstance(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item instance id"})
			return
		}

		eq, err := client.Equipment.Query().
			Where(equipment.IDEQ(id), equipment.EquipmentTemplateIDNEQ("")).
			Only(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "item instance not found"})
			return
		}

		c.JSON(http.StatusOK, toItemInstanceView(eq))
	}
}

// PUT /api/item-instances/:id — update instance fields
func updateItemInstance(client *db.Client) gin.HandlerFunc {
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
			RevealCondition   *string `json:"revealCondition"`
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

		updated, err := updater.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "item instance not found"})
			return
		}

		c.JSON(http.StatusOK, toItemInstanceView(updated))
	}
}

// DELETE /api/item-instances/:id — hard delete from DB
func deleteItemInstance(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item instance id"})
			return
		}

		err = client.Equipment.DeleteOneID(id).Exec(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "item instance not found"})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

// GET /api/equipment-templates — list all templates for admin spawn UI
func listEquipmentTemplates(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		templates, err := client.EquipmentTemplate.Query().All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		result := make([]gin.H, len(templates))
		for i, t := range templates {
			result[i] = gin.H{
				"id":                  t.ID,
				"name":                t.Name,
				"description":         t.Description,
				"slot":                t.Slot,
				"level":               t.Level,
				"weight":              t.Weight,
				"item_type":           t.ItemType,
				"stats":               t.Stats,
				"color":               t.Color,
				"is_visible":          t.IsVisible,
				"is_immovable":        t.IsImmovable,
				"effect_type":         t.EffectType,
				"effect_value":        t.EffectValue,
				"effect_duration":     t.EffectDuration,
				"is_container":        t.IsContainer,
				"container_capacity":  t.ContainerCapacity,
				"is_locked":           t.IsLocked,
				"key_item_id":         t.KeyItemID,
				"reveal_condition":    t.RevealCondition,
				"expires_at":          t.ExpiresAt,
			}
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

		c.JSON(http.StatusOK, gin.H{
			"id":                  t.ID,
			"name":                t.Name,
			"description":         t.Description,
			"slot":                t.Slot,
			"level":               t.Level,
			"weight":              t.Weight,
			"item_type":           t.ItemType,
			"stats":               t.Stats,
			"color":               t.Color,
			"is_visible":          t.IsVisible,
			"is_immovable":        t.IsImmovable,
			"effect_type":         t.EffectType,
			"effect_value":        t.EffectValue,
			"effect_duration":     t.EffectDuration,
			"is_container":        t.IsContainer,
			"container_capacity":  t.ContainerCapacity,
			"is_locked":           t.IsLocked,
			"key_item_id":         t.KeyItemID,
			"reveal_condition":    t.RevealCondition,
			"expires_at":          t.ExpiresAt,
		})
	}
}
