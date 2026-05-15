package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/repository"
	"herbst-server/service"
)

// applyDamage handles POST /characters/:id/damage.
func applyDamage(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := getIDParam(c)
		if !ok {
			return
		}
		var req struct {
			Damage     int `json:"damage" binding:"required"`
			AttackerID int `json:"attacker_id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := svc.Combat.ApplyDamage(c.Request.Context(), id, req.Damage)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if req.AttackerID > 0 {
			svc.Combat.LogDamage(c.Request.Context(), req.AttackerID, id, req.Damage)
		}
		resp := gin.H{
			"id":       result.ID,
			"hp":       result.HP,
			"maxHp":    result.MaxHP,
			"defeated": result.Defeated,
		}
		if result.Immortal {
			resp["immortal"] = true
			resp["message"] = result.Message
		}
		c.JSON(http.StatusOK, resp)
	}
}

// healCharacter handles POST /characters/:id/heal.
func healCharacter(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := getIDParam(c)
		if !ok {
			return
		}
		var req struct {
			Amount int `json:"amount" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := svc.Combat.HealCharacter(c.Request.Context(), id, req.Amount)
		if err != nil {
			status := http.StatusInternalServerError
			if err.Error() == "character not found" {
				status = http.StatusNotFound
			}
			c.JSON(status, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": result.ID, "hp": result.HP, "maxHp": result.MaxHP})
	}
}

// adjustStamina handles POST /characters/:id/stamina.
func adjustStamina(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := getIDParam(c)
		if !ok {
			return
		}
		var req struct {
			Amount int `json:"amount" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := svc.Combat.AdjustStamina(c.Request.Context(), id, req.Amount)
		if err != nil {
			status := http.StatusInternalServerError
			if err.Error() == "character not found" {
				status = http.StatusNotFound
			}
			c.JSON(status, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": result.ID, "stamina": result.Current, "maxStamina": result.Max})
	}
}

// adjustMana handles POST /characters/:id/mana.
func adjustMana(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := getIDParam(c)
		if !ok {
			return
		}
		var req struct {
			Amount int `json:"amount" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := svc.Combat.AdjustMana(c.Request.Context(), id, req.Amount)
		if err != nil {
			status := http.StatusInternalServerError
			if err.Error() == "character not found" {
				status = http.StatusNotFound
			}
			c.JSON(status, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": result.ID, "mana": result.Current, "maxMana": result.Max})
	}
}

// getCombatStatus handles GET /characters/:id/combat-status.
func getCombatStatus(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := getIDParam(c)
		if !ok {
			return
		}
		result, err := svc.Combat.GetCombatStatus(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": result.ID, "hp": result.HP, "maxHp": result.MaxHP, "isNPC": result.IsNPC})
	}
}

// healNPCsInRoom handles POST /rooms/:id/npcs/heal.
func healNPCsInRoom(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		roomID, err := parseIntParam(c, "id")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
			return
		}
		var req struct {
			Amount int `json:"amount" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		healed, err := svc.Combat.HealNPCsInRoom(c.Request.Context(), roomID, req.Amount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"healed": healed, "amount": req.Amount})
	}
}

// passiveHealNPCsInRoom handles POST /rooms/:id/npcs/passive-heal.
func passiveHealNPCsInRoom(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		roomID, err := parseIntParam(c, "id")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
			return
		}
		result, err := svc.Combat.PassiveHealNPCsInRoom(c.Request.Context(), roomID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"healed": result.Healed, "room": roomID})
	}
}

// NPC routes using repos directly.
func getNPCsByRoom(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		roomID, err := parseIntParam(c, "id")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
			return
		}
		npcs, err := repos.Character.ListNPCsByRoom(c.Request.Context(), roomID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result := make([]gin.H, len(npcs))
		for i, npc := range npcs {
			xpValue := npc.Level * 10
			if npc.NpcTemplateID != "" {
				tmpl, err := repos.NPCTemplate.Get(c.Request.Context(), npc.NpcTemplateID)
				if err == nil && tmpl.XpValue > 0 {
					xpValue = tmpl.XpValue
				}
			}
			result[i] = gin.H{
				"id":            npc.ID,
				"name":          npc.Name,
				"isNPC":         npc.IsNPC,
				"currentRoomId": npc.CurrentRoomId,
				"race":          npc.Race,
				"class":         npc.Class,
				"level":         npc.Level,
				"hitpoints":     npc.Hitpoints,
				"max_hitpoints": npc.MaxHitpoints,
				"xpValue":       xpValue,
			}
		}
		c.JSON(http.StatusOK, gin.H{"roomId": roomID, "npcs": result, "count": len(result)})
	}
}

func listAllNPCs(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Filter by world_id if provided
		worldID := c.Query("world_id")
		var npcs []*db.Character
		var err error
		if worldID != "" {
			npcs, err = repos.Character.ListAllByWorld(c.Request.Context(), worldID)
		} else {
			npcs, err = repos.Character.ListAllNPCs(c.Request.Context())
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result := make([]gin.H, len(npcs))
		for i, npc := range npcs {
			result[i] = gin.H{
				"id":            npc.ID,
				"name":          npc.Name,
				"isNPC":         npc.IsNPC,
				"currentRoomId": npc.CurrentRoomId,
				"race":          npc.Race,
				"class":         npc.Class,
				"level":         npc.Level,
				"hitpoints":     npc.Hitpoints,
				"max_hitpoints": npc.MaxHitpoints,
				"stamina":       npc.Stamina,
				"max_stamina":   npc.MaxStamina,
				"mana":          npc.Mana,
				"max_mana":      npc.MaxMana,
				"constitution":  npc.Constitution,
				"strength":      npc.Strength,
				"dexterity":     npc.Dexterity,
				"intelligence":  npc.Intelligence,
				"wisdom":        npc.Wisdom,
			}
		}
		c.JSON(http.StatusOK, gin.H{"npcs": result, "count": len(result)})
	}
}