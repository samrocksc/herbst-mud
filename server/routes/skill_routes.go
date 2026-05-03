package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/skill"
	"herbst-server/middleware"
)

// RegisterSkillRoutes registers REST endpoints for skills.
// Protected /api routes — all require JWT auth + admin check
func RegisterSkillRoutes(r *gin.Engine, client *db.Client) {
	skills := r.Group("/api")
	skills.Use(middleware.AuthMiddleware())
	skills.Use(middleware.AdminMiddleware())
	{
		skills.GET("/skills", listSkills(client))
		skills.POST("/skills", createSkill(client))
		skills.PUT("/skills/:id", updateSkill(client))
	}
}

// skillView is the JSON shape returned by the API — all 20+ fields.
type skillView struct {
	ID                     int     `json:"id"`
	Name                   string  `json:"name"`
	Description            string  `json:"description"`
	SkillType              string  `json:"skill_type"`
	Cost                   int     `json:"cost"`
	Cooldown               int     `json:"cooldown"`
	Requirements           string  `json:"requirements"`
	EffectType             string  `json:"effect_type"`
	EffectValue            int     `json:"effect_value"`
	EffectDuration         int     `json:"effect_duration"`
	ScalingStat            string  `json:"scaling_stat"`
	ScalingPercentPerPoint float64 `json:"scaling_percent_per_point"`
	ManaCost               int     `json:"mana_cost"`
	StaminaCost            int     `json:"stamina_cost"`
	HpCost                 int     `json:"hp_cost"`
	Slug                   string  `json:"slug"`
	RequiredTag            string  `json:"required_tag"`
	SkillClass             string  `json:"skill_class"`
	ProcChance             float64 `json:"proc_chance"`
	ProcEvent              string  `json:"proc_event"`
	CooldownSeconds        int     `json:"cooldown_seconds"`
	Tags                   []string `json:"tags"`
	FactionID              *int    `json:"faction_id,omitempty"`
	FactionName            string  `json:"faction_name,omitempty"`
}

// skillInput is the request body for create and update.
// Pointer types (*) mean "not set = don't change this field".
type skillInput struct {
	Name                   string   `json:"name"`
	Description            string   `json:"description"`
	SkillType              string   `json:"skill_type"`
	Cost                   *int     `json:"cost"`
	Cooldown               *int     `json:"cooldown"`
	Requirements           string   `json:"requirements"`
	EffectType             string   `json:"effect_type"`
	EffectValue            *int     `json:"effect_value"`
	EffectDuration         *int     `json:"effect_duration"`
	ScalingStat            string   `json:"scaling_stat"`
	ScalingPercentPerPoint *float64 `json:"scaling_percent_per_point"`
	ManaCost               *int     `json:"mana_cost"`
	StaminaCost            *int     `json:"stamina_cost"`
	HpCost                 *int     `json:"hp_cost"`
	Slug                   string   `json:"slug"`
	RequiredTag            string   `json:"required_tag"`
	SkillClass             string   `json:"skill_class"`
	ProcChance             *float64 `json:"proc_chance"`
	ProcEvent              string   `json:"proc_event"`
	CooldownSeconds        *int     `json:"cooldown_seconds"`
	Tags                   []string `json:"tags"`
	FactionID              *int     `json:"faction_id"`
}

// parseTagsFromRequirements extracts the tags array stored in the requirements JSON column.
// Requirements format: {"tags":["tag1","tag2"]} or {}
func parseTagsFromRequirements(req string) []string {
	if req == "" {
		return nil
	}
	var parsed struct {
		Tags []string `json:"tags"`
	}
	if err := json.Unmarshal([]byte(req), &parsed); err != nil {
		return nil
	}
	return parsed.Tags
}

// buildRequirementsJSON encodes a tags slice into the requirements JSON column.
func buildRequirementsJSON(tags []string) string {
	if len(tags) == 0 {
		return `{}`
	}
	m := map[string][]string{"tags": tags}
	b, _ := json.Marshal(m)
	return string(b)
}

func skillToView(s *db.Skill) skillView {
	v := skillView{
		ID:                     s.ID,
		Name:                   s.Name,
		Description:            s.Description,
		SkillType:              s.SkillType,
		Cost:                   s.Cost,
		Cooldown:               s.Cooldown,
		Requirements:           s.Requirements,
		EffectType:             s.EffectType,
		EffectValue:            s.EffectValue,
		EffectDuration:         s.EffectDuration,
		ScalingStat:            s.ScalingStat,
		ScalingPercentPerPoint: s.ScalingPercentPerPoint,
		ManaCost:               s.ManaCost,
		StaminaCost:            s.StaminaCost,
		HpCost:                 s.HpCost,
		Slug:                   s.Slug,
		RequiredTag:            s.RequiredTag,
		SkillClass:             s.SkillClass,
		ProcChance:             s.ProcChance,
		ProcEvent:              s.ProcEvent,
		CooldownSeconds:        s.CooldownSeconds,
		Tags:                   parseTagsFromRequirements(s.Requirements),
	}
	if s.Edges.Faction != nil {
		v.FactionID = &s.Edges.Faction.ID
		v.FactionName = s.Edges.Faction.Name
	}
	return v
}

func listSkills(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		skills, err := client.Skill.Query().
			WithFaction().
			Order(skill.ByName()).
			All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result := make([]skillView, len(skills))
		for i, s := range skills {
			result[i] = skillToView(s)
		}
		c.JSON(http.StatusOK, gin.H{"skills": result})
	}
}

func createSkill(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input skillInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if input.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
			return
		}

		mut := client.Skill.Create().
			SetName(input.Name).
			SetDescription(input.Description).
			SetSkillType(input.SkillType).
			SetRequirements(buildRequirementsJSON(input.Tags)).
			SetEffectType(input.EffectType).
			SetScalingStat(input.ScalingStat).
			SetSlug(input.Slug).
			SetRequiredTag(input.RequiredTag).
			SetSkillClass(input.SkillClass).
			SetProcEvent(input.ProcEvent)

		if input.Cost != nil {
			mut.SetCost(*input.Cost)
		}
		if input.Cooldown != nil {
			mut.SetCooldown(*input.Cooldown)
		}
		if input.EffectValue != nil {
			mut.SetEffectValue(*input.EffectValue)
		}
		if input.EffectDuration != nil {
			mut.SetEffectDuration(*input.EffectDuration)
		}
		if input.ScalingPercentPerPoint != nil {
			mut.SetScalingPercentPerPoint(*input.ScalingPercentPerPoint)
		}
		if input.ManaCost != nil {
			mut.SetManaCost(*input.ManaCost)
		}
		if input.StaminaCost != nil {
			mut.SetStaminaCost(*input.StaminaCost)
		}
		if input.HpCost != nil {
			mut.SetHpCost(*input.HpCost)
		}
		if input.ProcChance != nil {
			mut.SetProcChance(*input.ProcChance)
		}
		if input.CooldownSeconds != nil {
			mut.SetCooldownSeconds(*input.CooldownSeconds)
		}
		if input.FactionID != nil {
			mut.SetFactionID(*input.FactionID)
		}

		s, err := mut.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		s, _ = client.Skill.Query().
			WithFaction().
			Where(skill.ID(s.ID)).
			Only(c.Request.Context())
		c.JSON(http.StatusCreated, skillToView(s))
	}
}

func updateSkill(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid skill id"})
			return
		}

		s, err := client.Skill.Query().
			Where(skill.ID(id)).
			Only(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "skill not found"})
			return
		}

		var input skillInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		mut := client.Skill.UpdateOne(s).
			SetName(input.Name).
			SetDescription(input.Description).
			SetSkillType(input.SkillType).
			SetRequirements(buildRequirementsJSON(input.Tags)).
			SetEffectType(input.EffectType).
			SetScalingStat(input.ScalingStat).
			SetSlug(input.Slug).
			SetRequiredTag(input.RequiredTag).
			SetSkillClass(input.SkillClass).
			SetProcEvent(input.ProcEvent)

		if input.Cost != nil {
			mut.SetCost(*input.Cost)
		}
		if input.Cooldown != nil {
			mut.SetCooldown(*input.Cooldown)
		}
		if input.EffectValue != nil {
			mut.SetEffectValue(*input.EffectValue)
		}
		if input.EffectDuration != nil {
			mut.SetEffectDuration(*input.EffectDuration)
		}
		if input.ScalingPercentPerPoint != nil {
			mut.SetScalingPercentPerPoint(*input.ScalingPercentPerPoint)
		}
		if input.ManaCost != nil {
			mut.SetManaCost(*input.ManaCost)
		}
		if input.StaminaCost != nil {
			mut.SetStaminaCost(*input.StaminaCost)
		}
		if input.HpCost != nil {
			mut.SetHpCost(*input.HpCost)
		}
		if input.ProcChance != nil {
			mut.SetProcChance(*input.ProcChance)
		}
		if input.CooldownSeconds != nil {
			mut.SetCooldownSeconds(*input.CooldownSeconds)
		}
		if input.FactionID != nil {
			mut.SetFactionID(*input.FactionID)
		}

		updated, err := mut.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		updated, _ = client.Skill.Query().
			WithFaction().
			Where(skill.ID(updated.ID)).
			Only(c.Request.Context())
		c.JSON(http.StatusOK, skillToView(updated))
	}
}
