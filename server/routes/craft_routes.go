package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/dblog"
	"herbst-server/repository"
	"log/slog"
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
			slog.Warn("invalid craft request", slog.String("service", "crafting"), slog.String("error", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid request"})
			return
		}

		if req.Recipe == "" {
			slog.Warn("craft recipe name missing", slog.String("service", "crafting"))
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "recipe name is required"})
			return
		}

		ctx := c.Request.Context()

		recipe, err := repos.CraftingRecipe.Get(ctx, req.Recipe)
		if err != nil {
			slog.Warn("craft recipe not found", slog.String("service", "crafting"), slog.String("recipe", req.Recipe))
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "recipe not found"})
			return
		}

		char, err := repos.Character.Get(ctx, charID)
		if err != nil {
			slog.Warn("craft character not found", slog.String("service", "crafting"), slog.Int("character_id", charID))
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "character not found"})
			return
		}

		room, err := repos.Room.Get(ctx, char.CurrentRoomId)
		if err != nil {
			slog.Warn("craft room not found", slog.String("service", "crafting"), slog.Int("character_id", charID))
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "room not found"})
			return
		}

		if !roomHasTag(room, recipe.RequiredStationTag) {
			slog.Warn("craft missing station", slog.String("service", "crafting"), slog.String("required_tag", recipe.RequiredStationTag))
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "this room does not have a " + recipe.RequiredStationTag})
			return
		}

		if recipe.RequiredClass != "" && char.Class != recipe.RequiredClass {
			slog.Warn("craft wrong class", slog.String("service", "crafting"), slog.String("required_class", recipe.RequiredClass))
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
				slog.Warn("craft insufficient skill", slog.String("service", "crafting"), slog.String("skill", recipe.RequiredSkill), slog.Int("required_level", recipe.RequiredSkillLevel), slog.Int("current_level", level))
				c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": fmt.Sprintf("you don't have enough skill in %s (need %d, have %d)", recipe.RequiredSkill, recipe.RequiredSkillLevel, level)})
				return
			}
		}

		inventory, err := repos.Equipment.ListByOwner(ctx, charID)
		if err != nil {
			dblog.Error("failed to check inventory for crafting", err, slog.String("service", "crafting"), slog.Int("character_id", charID))
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "failed to check inventory"})
			return
		}

		available := make(map[int]int)
		for _, item := range inventory {
			if item.EquipmentTemplateID > 0 {
				available[item.EquipmentTemplateID] += item.Quantity
			}
		}

		// Resolve input slugs to template IDs
		type inputSlot struct {
			TemplateID int
			Quantity   int
			Consumed   bool
		}
		var resolvedInputs []inputSlot
		for _, input := range recipe.Inputs {
			tmpl, err := repos.EquipmentTemplate.GetBySlug(ctx, input.EquipmentTemplateSlug, "")
			if err != nil {
				slog.Warn("craft missing ingredient template", slog.String("service", "crafting"), slog.String("slug", input.EquipmentTemplateSlug))
				c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "missing ingredient template: " + input.EquipmentTemplateSlug})
				return
			}
			if available[tmpl.ID] < input.Quantity {
				slog.Warn("craft missing ingredient", slog.String("service", "crafting"), slog.String("name", tmpl.Name))
				c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "missing ingredient: " + tmpl.Name})
				return
			}
			resolvedInputs = append(resolvedInputs, inputSlot{TemplateID: tmpl.ID, Quantity: input.Quantity, Consumed: input.Consumed})
		}

		// Consume ingredients from inventory
		for _, input := range resolvedInputs {
			need := input.Quantity
			for _, item := range inventory {
				if item.EquipmentTemplateID != input.TemplateID || need <= 0 {
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
			template, err := repos.EquipmentTemplate.GetBySlug(ctx, output.EquipmentTemplateSlug, "")
			if err != nil {
				dblog.Error("failed to load output template for crafting", err, slog.String("service", "crafting"), slog.String("slug", output.EquipmentTemplateSlug))
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "failed to load template: " + output.EquipmentTemplateSlug})
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
					EquipmentTemplateID:   &template.ID,
					OwnerID:               &ownerID,
				})
				if err != nil {
					dblog.Error("failed to create crafted item", err, slog.String("service", "crafting"), slog.String("name", template.Name))
					c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "failed to create item: " + template.Name})
					return
				}
				outputs = append(outputs, outputItem{
					Name:       template.Name,
					InstanceID: created.ID,
				})
			}
		}

		slog.Info("item crafted", slog.String("service", "crafting"), slog.Int("character_id", charID), slog.Int("output_count", len(outputs)))
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
		slog.Warn("invalid character id", slog.String("service", "crafting"), slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return 0, false
	}
	return id, true
}