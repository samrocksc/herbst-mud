package routes

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/db/characterskill"
	"herbst-server/db/skill"
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
			"charisma":      char.Charisma,
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
			Charisma     *int `json:"charisma"`
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
			"wisdom": req.Wisdom, "charisma": req.Charisma,
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
		Charisma:     req.Charisma,
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
		"charisma":      char.Charisma,
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
		// Query character skills from the character_skills join table
		client := c.MustGet("db_client").(*db.Client)
		charSkills, err := client.CharacterSkill.Query().
			Where(characterskill.HasCharacterWith(character.IDEQ(id))).
			WithSkill().
			All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load skills"})
			return
		}
		skillsMap := gin.H{}
		for _, cs := range charSkills {
			if cs.Edges.Skill != nil {
				skillsMap[cs.Edges.Skill.Name] = gin.H{
					"level": cs.Level,
					"bonus": service.CalcSkillBonus(cs.Level),
				}
			}
		}
		c.JSON(http.StatusOK, gin.H{
		"id":               char.ID,
		"name":             char.Name,
		"skills":           skillsMap,
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
			Skills map[string]int `json:"skills"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Validate all skill levels
		for name, level := range req.Skills {
			if level < 0 || level > 100 {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s must be between 0 and 100", name)})
				return
			}
		}
		client := c.MustGet("db_client").(*db.Client)
		ctx := c.Request.Context()
		// Verify character exists
		char, err := repos.Character.Get(ctx, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}
		// For each skill name, find or create the character_skills record
		for skillName, level := range req.Skills {
			// Find the skill by name in this world
			sk, err := client.Skill.Query().
				Where(skill.NameEQ(skillName), skill.WorldIDEQ(char.WorldID)).
				Only(ctx)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("skill not found: %s", skillName)})
				return
			}
			// Check if character_skill record already exists
			existing, err := client.CharacterSkill.Query().
				Where(
					characterskill.HasCharacterWith(character.IDEQ(id)),
					characterskill.HasSkillWith(skill.IDEQ(sk.ID)),
				).
				Only(ctx)
			if err != nil && !db.IsNotFound(err) {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query character skills"})
				return
			}
			if existing != nil {
				// Update existing record
				_, err = client.CharacterSkill.UpdateOneID(existing.ID).
					SetLevel(level).
					Save(ctx)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update skill"})
					return
				}
			} else {
				// Create new record
				err = client.CharacterSkill.Create().
					SetCharacterID(id).
					SetSkillID(sk.ID).
					SetLevel(level).
					Exec(ctx)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create skill"})
					return
				}
			}
		}
		// Return updated skills
		charSkills, err := client.CharacterSkill.Query().
			Where(characterskill.HasCharacterWith(character.IDEQ(id))).
			WithSkill().
			All(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load skills"})
			return
		}
		skillsMap := gin.H{}
		for _, cs := range charSkills {
			if cs.Edges.Skill != nil {
				skillsMap[cs.Edges.Skill.Name] = cs.Level
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"id":     char.ID,
			"name":   char.Name,
			"skills": skillsMap,
		})
	}
}