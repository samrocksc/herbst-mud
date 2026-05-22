package routes

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/middleware"
	"herbst-server/repository"
)

// RegisterNPCInstanceRoutes registers REST endpoints for NPC instances.
// NPC instances are Character rows with isNPC=true and is_instance=true.
func RegisterNPCInstanceRoutes(r *gin.Engine, repos *repository.Container, client *db.Client) {
	g := r.Group("/api")
	g.Use(middleware.AuthMiddleware(nil))
	g.Use(middleware.AdminMiddleware())
	g.Use(middleware.WorldAccessMiddleware())
	{
		g.GET("/npc-instances", listNPCInstances(repos, client))
		g.POST("/npc-instances", createNPCInstance(repos, client))
		g.GET("/npc-instances/:id", getNPCInstance(repos))
		g.PUT("/npc-instances/:id", updateNPCInstance(client))
		g.DELETE("/npc-instances/:id", deleteNPCInstance(repos))

		// Equipment management for NPC instances
		g.GET("/npc-instances/:id/equipment", listNPCInstanceEquipment(repos))
		g.POST("/npc-instances/:id/equipment", addNPCInstanceEquipment(repos))
		g.DELETE("/npc-instances/:id/equipment/:eqid", removeNPCInstanceEquipment(repos))
	}
}

// ─── JSON views ─────────────────────────────────────────────────────────────

// npcInstanceView is the JSON shape returned by the API.
type npcInstanceView struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	NpcTemplateID  string `json:"npc_template_id"`
	InstanceNumber int    `json:"instance_number"`
	RoomID         int    `json:"room_id"`
	StartingRoomID int    `json:"starting_room_id"`
	Level          int    `json:"level"`
	Race           string `json:"race"`
	Hitpoints      int    `json:"hitpoints"`
	MaxHitpoints   int    `json:"max_hitpoints"`
	Stamina        int    `json:"stamina"`
	MaxStamina     int    `json:"max_stamina"`
	Mana           int    `json:"mana"`
	MaxMana        int    `json:"max_mana"`
	IsNPC          bool   `json:"isNPC"`
	IsInstance     bool   `json:"is_instance"`
	WorldID        string `json:"world_id"`
}

func toView(c *db.Character) npcInstanceView {
	return npcInstanceView{
		ID:             c.ID,
		Name:           c.Name,
		NpcTemplateID:  c.NpcTemplateID,
		InstanceNumber: c.InstanceNumber,
		RoomID:         c.CurrentRoomId,
		StartingRoomID: c.StartingRoomId,
		Level:          c.Level,
		Race:           c.Race,
		Hitpoints:      c.Hitpoints,
		MaxHitpoints:   c.MaxHitpoints,
		Stamina:        c.Stamina,
		MaxStamina:     c.MaxStamina,
		Mana:           c.Mana,
		MaxMana:        c.MaxMana,
		IsNPC:          c.IsNPC,
		IsInstance:     c.IsInstance,
		WorldID:        c.CurrentWorld,
	}
}

// ─── Handlers ───────────────────────────────────────────────────────────────

// GET /api/npc-instances?roomId=X&templateId=X&active=true&world_id=X
func listNPCInstances(repos *repository.Container, client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Remove client dependency once repos support is_instance and roomId/templateId filters
		query := client.Character.Query().
			Where(character.IsNPCEQ(true), character.IsInstanceEQ(true))

		// Filter by world_id
		worldID := c.Query("world_id")
		if worldID != "" {
			query = query.Where(character.CurrentWorldEQ(worldID))
		}

		// Optional filters
		if roomIDStr := c.Query("roomId"); roomIDStr != "" {
			roomID, err := strconv.Atoi(roomIDStr)
			if err == nil {
				query = query.Where(character.CurrentRoomIdEQ(roomID))
			}
		}
		if templateID := c.Query("templateId"); templateID != "" {
			query = query.Where(character.NpcTemplateIDEQ(templateID))
		}

		instances, err := query.All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Filter active (hp > 0) when requested
		if active := c.Query("active"); active == "true" {
			filtered := make([]*db.Character, 0, len(instances))
			for _, inst := range instances {
				if inst.Hitpoints > 0 {
					filtered = append(filtered, inst)
				}
			}
			instances = filtered
		}

		result := make([]npcInstanceView, len(instances))
		for i, inst := range instances {
			result[i] = toView(inst)
		}

		c.JSON(http.StatusOK, result)
	}
}

// POST /api/npc-instances — create instance from template
func createNPCInstance(repos *repository.Container, client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			TemplateID     string `json:"template_id"`
			RoomID         int    `json:"room_id"`
			InstanceNumber int    `json:"instance_number"`
			WorldID        string `json:"world_id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if req.TemplateID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "template_id is required"})
			return
		}
		if req.RoomID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "room_id is required"})
			return
		}

		// Use world_id from request or default to "default"
		worldID := req.WorldID
		if worldID == "" {
			worldID = "default"
		}

		// Fetch the NPC template via repo
		tmpl, err := repos.NPCTemplate.Get(c.Request.Context(), req.TemplateID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "npc template not found: " + req.TemplateID})
			return
		}

		// Determine instance_number: if 0, auto-assign next available
		instanceNum := req.InstanceNumber
		if instanceNum == 0 {
			instanceNum, err = nextInstanceNumber(client, c.Request.Context(), tmpl.ID, worldID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		// TODO: Use repos.Character.Create once CreateCharacterInput supports IsInstance and InstanceNumber
		builder := client.Character.Create().
			SetName(tmpl.Name).
			SetLevel(tmpl.Level).
			SetIsNPC(true).
			SetIsInstance(true).
			SetInstanceNumber(instanceNum).
			SetNpcTemplateID(tmpl.ID).
			SetCurrentRoomId(req.RoomID).
			SetStartingRoomId(req.RoomID).
			SetRace(tmpl.Race).
			SetHitpoints(100).
			SetMaxHitpoints(100).
			SetStamina(50).
			SetMaxStamina(50).
			SetMana(25).
			SetMaxMana(25).
			SetCurrentWorld(worldID)

		created, err := builder.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, toView(created))
	}
}

// GET /api/npc-instances/:id
func getNPCInstance(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid instance id"})
			return
		}

		inst, err := repos.Character.Get(c.Request.Context(), id)
		if err != nil || !inst.IsNPC || !inst.IsInstance {
			c.JSON(http.StatusNotFound, gin.H{"error": "npc instance not found"})
			return
		}

		c.JSON(http.StatusOK, toView(inst))
	}
}

// PUT /api/npc-instances/:id — update instance fields (admin)
// TODO: Use repos.Character.Update once CharacterUpdates supports IsInstance and InstanceNumber
func updateNPCInstance(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid instance id"})
			return
		}

		var req struct {
			RoomID         *int  `json:"room_id"`
			StartingRoomID *int  `json:"starting_room_id"`
			Hitpoints      *int  `json:"hitpoints"`
			IsInstance     *bool `json:"is_instance"`
			InstanceNumber *int  `json:"instance_number"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// TODO: Use repos.Character.Update once CharacterUpdates supports IsInstance and InstanceNumber
		updater := client.Character.UpdateOneID(id)

		if req.RoomID != nil {
			updater.SetCurrentRoomId(*req.RoomID)
		}
		if req.StartingRoomID != nil {
			updater.SetStartingRoomId(*req.StartingRoomID)
		}
		if req.Hitpoints != nil {
			updater.SetHitpoints(*req.Hitpoints)
		}
		if req.IsInstance != nil {
			updater.SetIsInstance(*req.IsInstance)
		}
		if req.InstanceNumber != nil {
			updater.SetInstanceNumber(*req.InstanceNumber)
		}

		updated, err := updater.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "npc instance not found"})
			return
		}

		c.JSON(http.StatusOK, toView(updated))
	}
}

// DELETE /api/npc-instances/:id — hard delete from DB
func deleteNPCInstance(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid instance id"})
			return
		}

		if err := repos.Character.Delete(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "npc instance not found"})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

// ─── Equipment sub-endpoints ────────────────────────────────────────────────

// equipmentItemView is the lightweight representation of equipment assigned to an NPC.
type equipmentItemView struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Slot        string `json:"slot"`
	ItemType    string `json:"item_type"`
	IsEquipped  bool   `json:"is_equipped"`
}

// GET /api/npc-instances/:id/equipment — list equipment owned by this instance
func listNPCInstanceEquipment(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid instance id"})
			return
		}

		items, err := repos.Equipment.ListByOwner(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		result := make([]equipmentItemView, len(items))
		for i, item := range items {
			result[i] = equipmentItemView{
				ID:         item.ID,
				Name:       item.Name,
				Slot:       item.Slot,
				ItemType:   item.ItemType,
				IsEquipped: item.IsEquipped,
			}
		}

		c.JSON(http.StatusOK, result)
	}
}

// POST /api/npc-instances/:id/equipment — add equipment to instance
func addNPCInstanceEquipment(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid instance id"})
			return
		}

		// Verify the character exists and is an NPC instance
		char, err := repos.Character.Get(c.Request.Context(), id)
		if err != nil || !char.IsNPC || !char.IsInstance {
			c.JSON(http.StatusNotFound, gin.H{"error": "npc instance not found"})
			return
		}

		var req struct {
			EquipmentTemplateID int    `json:"equipment_template_id" binding:"required"`
			SlotOverride         string `json:"slot,omitempty"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Fetch the equipment template
		eqTmpl, err := repos.EquipmentTemplate.Get(c.Request.Context(), req.EquipmentTemplateID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "equipment template not found"})
			return
		}

		slot := eqTmpl.Slot
		if req.SlotOverride != "" {
			slot = req.SlotOverride
		}

		eqItem, err := repos.Equipment.Create(c.Request.Context(), repository.CreateEquipmentInput{
			Name:                  eqTmpl.Name,
			Description:           eqTmpl.Description,
			Slot:                  slot,
			OwnerID:               &id,
			Level:                 eqTmpl.Level,
			Weight:                 eqTmpl.Weight,
			ItemType:              eqTmpl.ItemType,
			Color:                  eqTmpl.Color,
			IsVisible:             eqTmpl.IsVisible,
			IsImmovable:           eqTmpl.IsImmovable,
			EffectType:            eqTmpl.EffectType,
			EffectValue:           eqTmpl.EffectValue,
			EffectDuration:        eqTmpl.EffectDuration,
			IsContainer:           eqTmpl.IsContainer,
			ContainerCapacity:     eqTmpl.ContainerCapacity,
			IsLocked:              eqTmpl.IsLocked,
			IsEquipped:            true,
			EquipmentTemplateID:   &eqTmpl.ID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, equipmentItemView{
			ID:         eqItem.ID,
			Name:       eqItem.Name,
			Slot:       eqItem.Slot,
			ItemType:   eqItem.ItemType,
			IsEquipped: eqItem.IsEquipped,
		})
	}
}

// DELETE /api/npc-instances/:id/equipment/:eqid — remove equipment from instance
func removeNPCInstanceEquipment(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		instID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid instance id"})
			return
		}
		eqID, err := strconv.Atoi(c.Param("eqid"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid equipment id"})
			return
		}

		// Verify this equipment belongs to this instance
		item, err := repos.Equipment.Get(c.Request.Context(), eqID)
		if err != nil || item.OwnerId == nil || *item.OwnerId != instID {
			c.JSON(http.StatusNotFound, gin.H{"error": "equipment not found for this instance"})
			return
		}

		if err := repos.Equipment.Delete(c.Request.Context(), eqID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

// ─── Helpers ────────────────────────────────────────────────────────────────

// nextInstanceNumber returns the next available instance_number for a given template.
func nextInstanceNumber(client *db.Client, ctx context.Context, templateID string, worldID string) (int, error) {
	query := client.Character.Query().
		Where(character.IsNPCEQ(true), character.IsInstanceEQ(true), character.NpcTemplateIDEQ(templateID))
	if worldID != "" {
		query = query.Where(character.CurrentWorldEQ(worldID))
	}
	instances, err := query.All(ctx)
	if err != nil {
		return 0, err
	}

	maxNum := 0
	for _, inst := range instances {
		if inst.InstanceNumber > maxNum {
			maxNum = inst.InstanceNumber
		}
	}

	return maxNum + 1, nil
}