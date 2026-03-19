package routes

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/equipment"
	"herbst-server/db/room"
)

// revealConditions stores reveal conditions in memory (GitHub #12)
// In production, this would be in the database
var (
	revealConditions = make(map[int]map[string]any)
	revealMutex      sync.RWMutex
)

// RegisterEquipmentRoutes registers all equipment-related routes
func RegisterEquipmentRoutes(router *gin.Engine, client *db.Client) {
	// Create a new equipment item
	router.POST("/equipment", func(c *gin.Context) {
		var req struct {
			Name        string `json:"name" binding:"required"`
			Description string `json:"description"`
			Slot        string `json:"slot"`
			Level       int    `json:"level"`
			Weight      int    `json:"weight"`
			IsEquipped  bool   `json:"isEquipped"`
			// Item system fields (GitHub #89)
			IsImmovable bool   `json:"isImmovable"`
			Color       string `json:"color"`
			IsVisible   bool   `json:"isVisible"`
			ItemType    string `json:"itemType"`
			RoomID      int    `json:"roomId"`
			// Hidden items (GitHub #12 - Look System)
			RevealCondition map[string]any `json:"revealCondition"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
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
			SetItemType(req.ItemType)

		if req.Color != "" {
			builder.SetColor(req.Color)
		}

		// Set room if provided
		if req.RoomID > 0 {
			builder.SetRoomID(req.RoomID)
		}

		eq, err := builder.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Store reveal condition in memory (GitHub #12)
		if req.RevealCondition != nil && len(req.RevealCondition) > 0 {
			revealMutex.Lock()
			revealConditions[eq.ID] = req.RevealCondition
			revealMutex.Unlock()
		}

		c.JSON(http.StatusCreated, eq)
	})

	// Get all equipment
	router.GET("/equipment", func(c *gin.Context) {
		items, err := client.Equipment.Query().All(c.Request.Context())
		if err != nil {
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

	// Get equipment in a room
	// Query param "includeHidden" can be set to "true" to include hidden items
	router.GET("/rooms/:id/equipment", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
			return
		}

		includeHidden := c.Query("includeHidden") == "true"

		query := client.Equipment.Query().
			Where(equipment.HasRoomWith(room.IDEQ(id)))

		// By default, only show visible items (GitHub #12)
		if !includeHidden {
			query = query.Where(equipment.IsVisible(true))
		}

		items, err := query.All(c.Request.Context())
		if err != nil {
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid equipment ID"})
			return
		}

		eq, err := client.Equipment.Get(c.Request.Context(), id)
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
			RevealCondition map[string]any `json:"revealCondition,omitempty"`
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
		}

		revealMutex.RLock()
		if cond, ok := revealConditions[eq.ID]; ok {
			result.RevealCondition = cond
		}
		revealMutex.RUnlock()

		c.JSON(http.StatusOK, result)
	})

	// Update equipment by ID
	router.PUT("/equipment/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
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
			RevealCondition map[string]any `json:"revealCondition"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
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

		c.JSON(http.StatusOK, eq)
	})

	// Delete equipment by ID
	router.DELETE("/equipment/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid equipment ID"})
			return
		}

		err = client.Equipment.DeleteOneID(id).Exec(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Equipment not found"})
			return
		}

		// Clean up reveal condition
		revealMutex.Lock()
		delete(revealConditions, id)
		revealMutex.Unlock()

		c.JSON(http.StatusNoContent, nil)
	})

	// Examine equipment endpoint (look-10)
	router.GET("/equipment/:id/examine", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid equipment ID"})
			return
		}

		eq, err := client.Equipment.Get(c.Request.Context(), id)
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
			"id":             eq.ID,
			"name":           eq.Name,
			"description":    eq.Description,
			"type":           eq.ItemType,
			"is_immovable":   eq.IsImmovable,
			"is_visible":     eq.IsVisible,
			"color":          eq.Color,
			"examine_level":  examineLevel,
			"hidden_details": hiddenDetails,
			"examine_xp":      1, // Grant XP for examining
		})
	})

	// Reveal hidden item (GitHub #12 - Look System)
	router.POST("/equipment/:id/reveal", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid equipment ID"})
			return
		}

		var req struct {
			RevealType string `json:"revealType"` // "examine", "perception_check", "use_item", "event"
			Target     string `json:"target"`
			SkillLevel int    `json:"skillLevel"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get the equipment item
		eq, err := client.Equipment.Get(c.Request.Context(), id)
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Item has no reveal condition"})
			return
		}

		// Check if the reveal type matches
		condType, _ := revealCond["type"].(string)
		if condType != req.RevealType {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reveal type for this item"})
			return
		}

		// Check target if required
		if target, ok := revealCond["target"].(string); ok && target != "" {
			if req.Target != target {
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
		visible := true
		updated, err := client.Equipment.UpdateOneID(id).
			SetIsVisible(visible).
			Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Remove the reveal condition (it's been used)
		revealMutex.Lock()
		delete(revealConditions, id)
		revealMutex.Unlock()

		c.JSON(http.StatusOK, updated)
	})

	// Get reveal condition for an item (GitHub #12)
	router.GET("/equipment/:id/reveal", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
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