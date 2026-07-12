package routes

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/dblog"
	"herbst-server/repository"
	"herbst-server/service"
)

// getCharacterClass handles GET /characters/:id/class.
func getCharacterClass(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := getIDParam(c)
		if !ok {
			return
		}
		char, err := repos.Character.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": char.ID, "name": char.Name, "class": char.Class})
	}
}

// updateCharacterClass handles PUT /characters/:id/class.
func updateCharacterClass(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := getIDParam(c)
		if !ok {
			return
		}
		var req struct {
			Class string `json:"class" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Validate class against DB: check if a faction with name=req.Class
		// exists in a "class" category for the character's world.
		char, err := repos.Character.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}
		client := c.MustGet("db_client").(*db.Client)
		if !isValidClassDB(c.Request.Context(), client, req.Class, char.CurrentWorld) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid class"})
			return
		}
		char, err = repos.Character.Update(c.Request.Context(), id, repository.CharacterUpdates{Class: &req.Class})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": char.ID, "name": char.Name, "class": char.Class})
	}
}

// getCharacterSpecialty handles GET /characters/:id/specialty.
func getCharacterSpecialty(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := getIDParam(c)
		if !ok {
			return
		}
		char, err := repos.Character.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": char.ID, "name": char.Name, "class": char.Class, "specialty": char.Specialty})
	}
}

// updateCharacterSpecialty handles PUT /characters/:id/specialty.
func updateCharacterSpecialty(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := getIDParam(c)
		if !ok {
			return
		}
		var req struct {
			Specialty string `json:"specialty" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		char, err := repos.Character.Update(c.Request.Context(), id, repository.CharacterUpdates{Specialty: &req.Specialty})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": char.ID, "name": char.Name, "class": char.Class, "specialty": char.Specialty})
	}
}

// getSpecialtiesForClass handles GET /classes/:class/specialties.
func getSpecialtiesForClass(c *gin.Context) {
	class := c.Param("class")
	worldID := c.Query("world_id")
	client := c.MustGet("db_client").(*db.Client)
	ctx := context.Background()

	// Query the class faction from DB, then read its specialties JSON field.
	f, err := getClassFactionByName(ctx, client, class, worldID)
	if err != nil || f == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Class not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"class": class, "specialties": f.Specialties})
}

// getCharacterStats handles GET /characters/:id/stats.
func getCharacterStats(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := getIDParam(c)
		if !ok {
			return
		}
		char, err := repos.Character.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}
		derivedStats := gin.H{
			"max_hp":        char.Constitution*10 + char.Level*10,
			"max_stamina":   char.Constitution*5 + char.Level*5,
			"max_mana":      char.Intelligence*5 + char.Level*5,
			"carry_weight":  char.Strength * 10,
			"dodge_chance":  char.Dexterity,
			"crit_chance":   char.Dexterity * 5 / 10,
		}
		c.JSON(http.StatusOK, gin.H{
			"id":            char.ID,
			"name":          char.Name,
			"strength":      char.Strength,
			"dexterity":     char.Dexterity,
			"constitution":  char.Constitution,
			"intelligence":  char.Intelligence,
			"wisdom":        char.Wisdom,
			"derived":       derivedStats,
		})
	}
}

// updateCharacterStats handles PUT /characters/:id/stats.
func updateCharacterStats(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := getIDParam(c)
		if !ok {
			return
		}
		var req struct {
			Strength     *int `json:"strength"`
			Dexterity    *int `json:"dexterity"`
			Constitution *int `json:"constitution"`
			Intelligence *int `json:"intelligence"`
			Wisdom       *int `json:"wisdom"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validateStat := func(stat *int, name string) error {
			if stat != nil && (*stat < 1 || *stat > 30) {
				return fmt.Errorf("%s must be between 1 and 30", name)
			}
			return nil
		}
		for name, stat := range map[string]*int{
			"strength": req.Strength, "dexterity": req.Dexterity,
			"constitution": req.Constitution, "intelligence": req.Intelligence,
			"wisdom": req.Wisdom,
		} {
			if err := validateStat(stat, name); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}
		char, err := repos.Character.Update(c.Request.Context(), id, repository.CharacterUpdates{
			Strength:     req.Strength,
			Dexterity:    req.Dexterity,
			Constitution: req.Constitution,
			Intelligence: req.Intelligence,
			Wisdom:       req.Wisdom,
		})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"id":            char.ID,
			"name":          char.Name,
			"strength":      char.Strength,
			"dexterity":     char.Dexterity,
			"constitution":  char.Constitution,
			"intelligence":  char.Intelligence,
			"wisdom":        char.Wisdom,
		})
	}
}

// getCharacterSkills handles GET /characters/:id/skills.
func getCharacterSkills(repos *repository.Container, svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := getIDParam(c)
		if !ok {
			return
		}
		char, err := repos.Character.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}
		abilitiesWithElig, err := svc.AbilityEligibility.AbilitiesForCharacterWithEligibility(c.Request.Context(), id)
		if err != nil {
			dblog.Error("eligibility check failed", err, slog.Int("character_id", id))
		}
		factionAbilities := make([]gin.H, 0)
		if err == nil {
			for _, swe := range abilitiesWithElig {
				sk := swe.Ability
				el := swe.Eligibility
				entry := gin.H{
					"id":               sk.ID,
					"name":             sk.Name,
					"slug":             sk.Slug,
					"ability_type":     sk.AbilityType,
					"ability_class":    sk.AbilityClass,
					"proc_chance":      sk.ProcChance,
					"proc_event":       sk.ProcEvent,
					"cooldown_seconds": sk.CooldownSeconds,
					"mana_cost":        sk.ManaCost,
					"stamina_cost":     sk.StaminaCost,
					"hp_cost":          sk.HpCost,
					"required_tag":     sk.RequiredTag,
					"eligible":         el.Eligible,
					"reason":           el.Reason,
				}
				if sk.Edges.Faction != nil {
					entry["faction_id"] = sk.Edges.Faction.ID
					entry["faction_name"] = sk.Edges.Faction.Name
				}
				factionAbilities = append(factionAbilities, entry)
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"id":   char.ID,
			"name": char.Name,
			"skills": gin.H{
				"blades":      gin.H{"level": char.SkillBlades, "bonus": service.CalcSkillBonus(char.SkillBlades)},
				"staves":      gin.H{"level": char.SkillStaves, "bonus": service.CalcSkillBonus(char.SkillStaves)},
				"knives":      gin.H{"level": char.SkillKnives, "bonus": service.CalcSkillBonus(char.SkillKnives)},
				"martial":     gin.H{"level": char.SkillMartial, "bonus": service.CalcSkillBonus(char.SkillMartial)},
				"brawling":    gin.H{"level": char.SkillBrawling, "bonus": service.CalcSkillBonus(char.SkillBrawling)},
				"tech":        gin.H{"level": char.SkillTech, "bonus": service.CalcSkillBonus(char.SkillTech)},
				"light_armor": gin.H{"level": char.SkillLightArmor, "bonus": service.CalcSkillBonus(char.SkillLightArmor)},
				"cloth_armor": gin.H{"level": char.SkillClothArmor, "bonus": service.CalcSkillBonus(char.SkillClothArmor)},
				"heavy_armor": gin.H{"level": char.SkillHeavyArmor, "bonus": service.CalcSkillBonus(char.SkillHeavyArmor)},
			},
			"faction_abilities": factionAbilities,
		})
	}
}

// updateCharacterSkills handles PUT /characters/:id/skills.
func updateCharacterSkills(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := getIDParam(c)
		if !ok {
			return
		}
		var req struct {
			Blades      *int `json:"blades"`
			Staves      *int `json:"staves"`
			Knives      *int `json:"knives"`
			Martial     *int `json:"martial"`
			Brawling    *int `json:"brawling"`
			Tech        *int `json:"tech"`
			LightArmor  *int `json:"light_armor"`
			ClothArmor  *int `json:"cloth_armor"`
			HeavyArmor  *int `json:"heavy_armor"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validateSkill := func(skill *int, name string) error {
			if skill != nil && (*skill < 0 || *skill > 100) {
				return fmt.Errorf("%s must be between 0 and 100", name)
			}
			return nil
		}
		for name, stat := range map[string]*int{
			"blades": req.Blades, "staves": req.Staves, "knives": req.Knives,
			"martial": req.Martial, "brawling": req.Brawling, "tech": req.Tech,
			"light_armor": req.LightArmor, "cloth_armor": req.ClothArmor, "heavy_armor": req.HeavyArmor,
		} {
			if err := validateSkill(stat, name); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}
		char, err := repos.Character.Update(c.Request.Context(), id, repository.CharacterUpdates{
			SkillBlades:      req.Blades,
			SkillStaves:      req.Staves,
			SkillKnives:      req.Knives,
			SkillMartial:     req.Martial,
			SkillBrawling:    req.Brawling,
			SkillTech:        req.Tech,
			SkillLightArmor:  req.LightArmor,
			SkillClothArmor:  req.ClothArmor,
			SkillHeavyArmor:  req.HeavyArmor,
		})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"id":           char.ID,
			"name":         char.Name,
			"blades":       char.SkillBlades,
			"staves":       char.SkillStaves,
			"knives":       char.SkillKnives,
			"martial":      char.SkillMartial,
			"brawling":     char.SkillBrawling,
			"tech":         char.SkillTech,
			"light_armor":  char.SkillLightArmor,
			"cloth_armor":  char.SkillClothArmor,
			"heavy_armor":  char.SkillHeavyArmor,
		})
	}
}