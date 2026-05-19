package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/repository"
)

func RegisterCraftRoutes(r *gin.Engine, repos *repository.Container) {
	chars := r.Group("/api/characters")
	{
		chars.POST("/:id/craft", craftHandler(repos))
	}
}

type craftRequest struct {
	Recipe string `json:"recipe"`
}

type outputItem struct {
	Name       string `json:"name"`
	InstanceID int    `json:"instance_id"`
}

func craftHandler(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		charID, ok := parseIDParam(c)
		if !ok {
			return
		}

		var req craftRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid request"})
			return
		}

		if req.Recipe == "" {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "recipe name is required"})
			return
		}

		ctx := c.Request.Context()

		recipe, err := repos.CraftingRecipe.Get(ctx, req.Recipe)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "recipe not found"})
			return
		}

		char, err := repos.Character.Get(ctx, charID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "character not found"})
			return
		}

		room, err := repos.Room.Get(ctx, char.CurrentRoomId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "room not found"})
			return
		}

		if !roomHasTag(room, recipe.RequiredStationTag) {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "this room does not have a " + recipe.RequiredStationTag})
			return
		}

		if recipe.RequiredClass != "" && char.Class != recipe.RequiredClass {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "you don't have the required class: " + recipe.RequiredClass})
			return
		}

		if recipe.RequiredSkill != "" {
			comp, err := repos.Competency.GetCharacterCompetency(ctx, charID, recipe.RequiredSkill)
			if err != nil || comp == nil || comp.Level < recipe.RequiredSkillLevel {
				level := 0
				if comp != nil {
					level = comp.Level
				}
				c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": fmt.Sprintf("you don't have enough skill in %s (need %d, have %d)", recipe.RequiredSkill, recipe.RequiredSkillLevel, level)})
				return
			}
		}

		inventory, err := repos.Equipment.ListByOwner(ctx, charID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "failed to check inventory"})
			return
		}

		available := make(map[string]int)
		for _, item := range inventory {
			if item.EquipmentTemplateID != "" {
				available[item.EquipmentTemplateID] += item.Quantity
			}
		}

		for _, input := range recipe.Inputs {
			if available[input.EquipmentTemplateID] < input.Quantity {
				template, _ := repos.EquipmentTemplate.Get(ctx, input.EquipmentTemplateID)
				name := input.EquipmentTemplateID
				if template != nil {
					name = template.Name
				}
				c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "missing ingredient: " + name})
				return
			}
		}

		for _, input := range recipe.Inputs {
			need := input.Quantity
			for _, item := range inventory {
				if item.EquipmentTemplateID != input.EquipmentTemplateID || need <= 0 {
					continue
				}
				if item.Quantity <= need {
					need -= item.Quantity
					repos.Equipment.Delete(ctx, item.ID)
				} else {
					newQty := item.Quantity - need
					need = 0
					repos.Equipment.Update(ctx, item.ID, repository.EquipmentUpdates{Level: &newQty})
				}
			}
		}

		var outputs []outputItem
		for _, output := range recipe.Outputs {
			template, err := repos.EquipmentTemplate.Get(ctx, output.EquipmentTemplateID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "failed to load template: " + output.EquipmentTemplateID})
				return
			}

			for i := 0; i < output.Quantity; i++ {
				ownerID := charID
				created, err := repos.Equipment.Create(ctx, repository.CreateEquipmentInput{
					Name:                  template.Name,
					Description:           template.Description,
					Slot:                  template.Slot,
					Level:                 template.Level,
					ItemType:              template.ItemType,
					ArmorRating:           template.ArmorRating,
					ArmorType:             template.ArmorType,
					DamageDiceCount:       template.DamageDiceCount,
					DamageDiceSides:       template.DamageDiceSides,
					DamageBonus:           template.DamageBonus,
					DamageType:            template.DamageType,
					WeaponType:            template.WeaponType,
					IsTwoHanded:           template.IsTwoHanded,
					Stats:                 template.Stats,
					Rarity:                template.Rarity,
					SkillRequirement:       template.SkillRequirement,
					SkillRequirementLevel:  template.SkillRequirementLevel,
					Weight:                template.Weight,
					IsEquipped:            false,
					IsImmovable:           template.IsImmovable,
					Color:                 template.Color,
					IsVisible:             template.IsVisible,
					EffectType:            template.EffectType,
					EffectValue:           template.EffectValue,
					EffectDuration:        template.EffectDuration,
					EquipmentTemplateID:   &output.EquipmentTemplateID,
					OwnerID:               &ownerID,
				})
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "failed to create item: " + template.Name})
					return
				}
				outputs = append(outputs, outputItem{
					Name:       template.Name,
					InstanceID: created.ID,
				})
			}
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "outputs": outputs})
	}
}

func roomHasTag(room *db.Room, tag string) bool {
	for _, t := range room.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

func parseIDParam(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return 0, false
	}
	return id, true
}