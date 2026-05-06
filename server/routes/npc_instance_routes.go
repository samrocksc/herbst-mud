package routes

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/db/equipment"
	"herbst-server/middleware"
)

// RegisterNPCInstanceRoutes registers REST endpoints for NPC instances.
// NPC instances are Character rows with isNPC=true and is_instance=true.
func RegisterNPCInstanceRoutes(r *gin.Engine, client *db.Client) {
	g := r.Group("/api")
	g.Use(middleware.AuthMiddleware())
	g.Use(middleware.AdminMiddleware())
	{
		g.GET("/npc-instances", listNPCInstances(client))
		g.POST("/npc-instances", createNPCInstance(client))
		g.GET("/npc-instances/:id", getNPCInstance(client))
		g.PUT("/npc-instances/:id", updateNPCInstance(client))
		g.DELETE("/npc-instances/:id", deleteNPCInstance(client))

		// Equipment management for NPC instances
		g.GET("/npc-instances/:id/equipment", listNPCInstanceEquipment(client))
		g.POST("/npc-instances/:id/equipment", addNPCInstanceEquipment(client))
		g.DELETE("/npc-instances/:id/equipment/:eqid", removeNPCInstanceEquipment(client))
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
	}
}

// ─── Handlers ───────────────────────────────────────────────────────────────

// GET /api/npc-instances?roomId=X&templateId=X&active=true
func listNPCInstances(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := client.Character.Query().
			Where(character.IsNPCEQ(true), character.IsInstanceEQ(true))

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
func createNPCInstance(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			TemplateID     string `json:"template_id"`
			RoomID         int    `json:"room_id"`
			InstanceNumber int    `json:"instance_number"`
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

		// Fetch the NPC template
		tmpl, err := client.NPCTemplate.Get(c.Request.Context(), req.TemplateID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "npc template not found: " + req.TemplateID})
			return
		}

		// Determine instance_number: if 0, auto-assign next available
		instanceNum := req.InstanceNumber
		if instanceNum == 0 {
			instanceNum, err = nextInstanceNumber(client, c.Request.Context(), tmpl.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

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
			SetMaxMana(25)

		created, err := builder.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, toView(created))
	}
}

// GET /api/npc-instances/:id
func getNPCInstance(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid instance id"})
			return
		}

		inst, err := client.Character.Query().
			Where(character.IDEQ(id), character.IsNPCEQ(true), character.IsInstanceEQ(true)).
			Only(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "npc instance not found"})
			return
		}

		c.JSON(http.StatusOK, toView(inst))
	}
}

// PUT /api/npc-instances/:id — update instance fields (admin)
//
// Accepts: { room_id?, starting_room_id?, hitpoints?, instance_number? }
// room_id maps to CurrentRoomId; starting_room_id maps to StartingRoomId
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
func deleteNPCInstance(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid instance id"})
			return
		}

		err = client.Character.DeleteOneID(id).Exec(c.Request.Context())
		if err != nil {
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
func listNPCInstanceEquipment(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid instance id"})
			return
		}

		items, err := client.Equipment.Query().
			Where(equipment.OwnerIdEQ(id)).
			All(c.Request.Context())
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
//
// Accepts: { equipment_template_id: string, slot?: string }
// Creates a new Equipment row from the template and assigns it to the character.
func addNPCInstanceEquipment(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid instance id"})
			return
		}

		// Verify the character exists and is an NPC instance
		_, err = client.Character.Query().
			Where(character.IDEQ(id), character.IsNPCEQ(true), character.IsInstanceEQ(true)).
			Only(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "npc instance not found"})
			return
		}

		var req struct {
			EquipmentTemplateID string `json:"equipment_template_id" binding:"required"`
			SlotOverride         string `json:"slot,omitempty"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Fetch the equipment template
		eqTmpl, err := client.EquipmentTemplate.Get(c.Request.Context(), req.EquipmentTemplateID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "equipment template not found: " + req.EquipmentTemplateID})
			return
		}

		slot := eqTmpl.Slot
		if req.SlotOverride != "" {
			slot = req.SlotOverride
		}

		// Create the equipment item and assign it to the character
		eqItem, err := client.Equipment.Create().
			SetEquipmentTemplateID(eqTmpl.ID).
			SetName(eqTmpl.Name).
			SetDescription(eqTmpl.Description).
			SetSlot(slot).
			SetLevel(eqTmpl.Level).
			SetWeight(eqTmpl.Weight).
			SetItemType(eqTmpl.ItemType).
			SetColor(eqTmpl.Color).
			SetIsVisible(eqTmpl.IsVisible).
			SetIsImmovable(eqTmpl.IsImmovable).
			SetEffectType(eqTmpl.EffectType).
			SetEffectValue(eqTmpl.EffectValue).
			SetEffectDuration(eqTmpl.EffectDuration).
			SetIsContainer(eqTmpl.IsContainer).
			SetContainerCapacity(eqTmpl.ContainerCapacity).
			SetIsLocked(eqTmpl.IsLocked).
			SetOwnerId(id).
			SetIsEquipped(true).
			Save(c.Request.Context())
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
func removeNPCInstanceEquipment(client *db.Client) gin.HandlerFunc {
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
		item, err := client.Equipment.Query().
			Where(equipment.IDEQ(eqID), equipment.OwnerIdEQ(instID)).
			Only(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "equipment not found for this instance"})
			return
		}

		err = client.Equipment.DeleteOne(item).Exec(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

// ─── Helpers ────────────────────────────────────────────────────────────────

// nextInstanceNumber returns the next available instance_number for a given template.
// It queries all NPC instances with the given template_id and finds max(instance_number) + 1.
func nextInstanceNumber(client *db.Client, ctx context.Context, templateID string) (int, error) {
	instances, err := client.Character.Query().
		Where(character.IsNPCEQ(true), character.IsInstanceEQ(true), character.NpcTemplateIDEQ(templateID)).
		All(ctx)
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
