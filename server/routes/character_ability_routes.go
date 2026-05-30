package routes

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/dblog"
	"herbst-server/service"
)

// getCharacterAbilities handles GET /characters/:id/abilities.
func getCharacterAbilities(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := getIDParam(c)
		if !ok {
			return
		}
		charAbilities, err := svc.Ability.GetAbilitiesWithDetails(c.Request.Context(), id)
		if err != nil {
			dblog.Error("failed to get character abilities", err, slog.String("service", "characters"), slog.Int("character_id", id))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		slots := make([]interface{}, 6)
		for i := range slots {
			slots[i] = nil
		}
		for _, ca := range charAbilities {
			entry := service.FormatAbilitySlot(ca)
			slot := ca.Slot
			if slot >= 0 && slot < 6 {
				slots[slot] = entry
			}
		}
		c.JSON(http.StatusOK, gin.H{"character_id": id, "slots": slots})
	}
}

// equipAbility handles POST /characters/:id/abilities.
func equipAbility(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := getIDParam(c)
		if !ok {
			return
		}
		var req struct {
			AbilityID int `json:"ability_id"`
			Slot      int `json:"slot"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			slog.Warn("bad request: invalid equip ability request", slog.String("service", "characters"), slog.Int("character_id", id), slog.String("error", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := svc.Ability.EquipAbility(c.Request.Context(), id, req.AbilityID, req.Slot); err != nil {
			switch {
			case isErr(err, service.ErrSlotOutOfRange):
				slog.Warn("bad request: slot out of range", slog.String("service", "characters"), slog.Int("character_id", id), slog.String("error", err.Error()))
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			case isErr(err, service.ErrAbilityNotFound):
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			case isErr(err, service.ErrMaxAbilities):
				slog.Warn("bad request: max abilities reached", slog.String("service", "characters"), slog.Int("character_id", id), slog.String("error", err.Error()))
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			default:
				dblog.Error("failed to equip ability", err, slog.String("service", "characters"), slog.Int("character_id", id), slog.Int("ability_id", req.AbilityID), slog.Int("slot", req.Slot))
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
		abilityObj, _ := svc.Ability.GetAbility(c.Request.Context(), req.AbilityID)
		abilityName := ""
		if abilityObj != nil {
			abilityName = abilityObj.Name
		}
		slog.Info("ability equipped", slog.String("service", "characters"), slog.Int("character_id", id), slog.Int("ability_id", req.AbilityID), slog.Int("slot", req.Slot))
		c.JSON(http.StatusCreated, gin.H{
			"success":      true,
			"slot":         req.Slot,
			"ability_id":   req.AbilityID,
			"ability_name": abilityName,
		})
	}
}

// unequipAbility handles DELETE /characters/:id/abilities/:slot.
func unequipAbility(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := getIDParam(c)
		if !ok {
			return
		}
		slot, err := strconv.Atoi(c.Param("slot"))
		if err != nil || slot < 1 || slot > 5 {
			slog.Warn("bad request: invalid slot", slog.String("service", "characters"), slog.Int("character_id", id), slog.String("slot", c.Param("slot")))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid slot"})
			return
		}
		if err := svc.Ability.UnequipAbility(c.Request.Context(), id, slot); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		slog.Info("ability unequipped", slog.String("service", "characters"), slog.Int("character_id", id), slog.Int("slot", slot))
		c.JSON(http.StatusOK, gin.H{"success": true, "slot": slot})
	}
}

// swapAbilities handles PUT /characters/:id/abilities/swap.
func swapAbilities(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := getIDParam(c)
		if !ok {
			return
		}
		var req struct {
			Slot1 int `json:"slot1"`
			Slot2 int `json:"slot2"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			slog.Warn("bad request: invalid swap request", slog.String("service", "characters"), slog.Int("character_id", id), slog.String("error", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := svc.Ability.SwapAbilities(c.Request.Context(), id, req.Slot1, req.Slot2)
		if err != nil {
			slog.Warn("bad request: failed to swap abilities", slog.String("service", "characters"), slog.Int("character_id", id), slog.String("error", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		slog.Info("abilities swapped", slog.String("service", "characters"), slog.Int("character_id", id), slog.Int("slot1", req.Slot1), slog.Int("slot2", req.Slot2))
		c.JSON(http.StatusOK, gin.H{"success": true, "slot1": req.Slot1, "slot2": req.Slot2, "result": result})
	}
}

// listPassiveAbilities handles GET /characters/:id/passive-abilities.
func listPassiveAbilities(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := getIDParam(c)
		if !ok {
			return
		}
		// Verify character exists via repos (service doesn't have simple Get)
		passives, err := svc.Ability.ListPassiveAbilities(c.Request.Context(), "")
		if err != nil {
			dblog.Error("failed to list passive abilities", err, slog.String("service", "characters"), slog.Int("character_id", id))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result := make([]gin.H, len(passives))
		for i, a := range passives {
			result[i] = gin.H{"id": a.ID, "name": a.Name, "description": a.Description}
		}
		c.JSON(http.StatusOK, gin.H{"character_id": id, "passive_abilities": result, "count": len(result)})
	}
}

// unlockPassiveAbility handles POST /characters/:id/passive-abilities.
func unlockPassiveAbility(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := getIDParam(c)
		if !ok {
			return
		}
		var req struct {
			AbilityID int `json:"ability_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			slog.Warn("bad request: invalid unlock passive request", slog.String("service", "characters"), slog.Int("character_id", id), slog.String("error", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		charAbility, err := svc.Ability.UnlockPassiveAbility(c.Request.Context(), id, req.AbilityID)
		if err != nil {
			switch {
			case isErr(err, service.ErrAbilityNotFound):
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			case isErr(err, service.ErrNotPassive):
				slog.Warn("bad request: ability is not passive", slog.String("service", "characters"), slog.Int("character_id", id), slog.String("error", err.Error()))
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			case isErr(err, service.ErrAlreadyEquipped):
				c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			case isErr(err, service.ErrNoAvailableSlots):
				slog.Warn("bad request: no available slots", slog.String("service", "characters"), slog.Int("character_id", id), slog.String("error", err.Error()))
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			default:
				dblog.Error("failed to unlock passive ability", err, slog.String("service", "characters"), slog.Int("character_id", id), slog.Int("ability_id", req.AbilityID))
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
		abilityObj, _ := svc.Ability.GetAbility(c.Request.Context(), req.AbilityID)
		abilityName := ""
		if abilityObj != nil {
			abilityName = abilityObj.Name
		}
		slog.Info("passive ability unlocked", slog.String("service", "characters"), slog.Int("character_id", id), slog.Int("ability_id", req.AbilityID), slog.Int("slot", charAbility.Slot))
		c.JSON(http.StatusCreated, gin.H{
			"success":      true,
			"id":           charAbility.ID,
			"ability_id":   req.AbilityID,
			"ability_name": abilityName,
			"slot":         charAbility.Slot,
		})
	}
}

// removePassiveAbility handles DELETE /characters/:id/passive-abilities/:abilityId.
func removePassiveAbility(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := getIDParam(c)
		if !ok {
			return
		}
		abilityId, err := strconv.Atoi(c.Param("abilityId"))
		if err != nil {
			slog.Warn("bad request: invalid ability ID", slog.String("service", "characters"), slog.Int("character_id", id), slog.String("ability_id", c.Param("abilityId")))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ability ID"})
			return
		}
		if err := svc.Ability.RemovePassiveAbility(c.Request.Context(), id, abilityId); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		slog.Info("passive ability removed", slog.String("service", "characters"), slog.Int("character_id", id), slog.Int("ability_id", abilityId))
		c.JSON(http.StatusOK, gin.H{"success": true, "ability_id": abilityId})
	}
}

// getClasslessSkills handles GET /characters/:id/classless-skills.
func getClasslessSkills(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := getIDParam(c)
		if !ok {
			return
		}
		charAbilities, err := svc.Ability.GetAbilitiesWithDetails(c.Request.Context(), id)
		if err != nil {
			dblog.Error("failed to get classless skills", err, slog.String("service", "characters"), slog.Int("character_id", id))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		skills := make([]gin.H, 0)
		for _, ca := range charAbilities {
			ab := ca.Edges.Ability
			if ab == nil || ab.AbilityClass != "active" {
				continue
			}
			effectList := make([]gin.H, 0)
			for _, e := range ab.Edges.Effects {
				effectList = append(effectList, gin.H{
					"effectType":    e.EffectType,
					"damageSubtype": e.DamageSubtype,
					"target":        e.Target,
					"value":         e.Value,
					"duration":      e.Duration,
					"scalingStat":   e.ScalingStat,
					"scalingRatio":  e.ScalingRatio,
					"sortOrder":     e.SortOrder,
				})
			}
			skills = append(skills, gin.H{
				"id":          ab.ID,
				"name":        ab.Name,
				"description": ab.Description,
				"slot":        ca.Slot,
				"cooldown":    ab.Cooldown,
				"manaCost":    ab.ManaCost,
				"staminaCost": ab.StaminaCost,
				"effects":     effectList,
			})
		}
		c.JSON(http.StatusOK, gin.H{"skills": skills})
	}
}

// equipClasslessSkill handles POST /characters/:id/classless-skills.
func equipClasslessSkill(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := getIDParam(c)
		if !ok {
			return
		}
		var req struct {
			SkillID int `json:"skill_id" binding:"required"`
			Slot    int `json:"slot" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			slog.Warn("bad request: invalid classless skill request", slog.String("service", "characters"), slog.Int("character_id", id), slog.String("error", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := svc.Ability.EquipClasslessSkill(c.Request.Context(), id, req.SkillID, req.Slot); err != nil {
			status := http.StatusBadRequest
			if isErr(err, service.ErrAbilityNotFound) {
				status = http.StatusNotFound
			}
			slog.Warn("bad request: failed to equip classless skill", slog.String("service", "characters"), slog.Int("character_id", id), slog.String("error", err.Error()))
			c.JSON(status, gin.H{"error": err.Error()})
			return
		}
		slog.Info("classless skill equipped", slog.String("service", "characters"), slog.Int("character_id", id), slog.Int("skill_id", req.SkillID), slog.Int("slot", req.Slot))
		c.JSON(http.StatusOK, gin.H{
			"message":      "Skill equipped",
			"skill_id":     req.SkillID,
			"slot":         req.Slot,
			"character_id": id,
		})
	}
}

// swapClasslessSkills handles PUT /characters/:id/classless-skills/swap.
func swapClasslessSkills(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := getIDParam(c)
		if !ok {
			return
		}
		var req struct {
			Slot1 int `json:"slot1" binding:"required"`
			Slot2 int `json:"slot2" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			slog.Warn("bad request: invalid classless skill swap request", slog.String("service", "characters"), slog.Int("character_id", id), slog.String("error", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := svc.Ability.SwapClasslessSkills(c.Request.Context(), id, req.Slot1, req.Slot2); err != nil {
			slog.Warn("bad request: failed to swap classless skills", slog.String("service", "characters"), slog.Int("character_id", id), slog.String("error", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		slog.Info("classless skills swapped", slog.String("service", "characters"), slog.Int("character_id", id), slog.Int("slot1", req.Slot1), slog.Int("slot2", req.Slot2))
		c.JSON(http.StatusOK, gin.H{
			"message":      "Skills swapped",
			"slot1":        req.Slot1,
			"slot2":        req.Slot2,
			"character_id": id,
		})
	}
}

func isErr(err, target error) bool {
	return err == target
}
