package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/ability"
	"herbst-server/middleware"
)

func RegisterAbilityRoutes(r *gin.Engine, client *db.Client) {
	r.GET("/abilities/classless", listClasslessAbilities(client))

	abilities := r.Group("/api")
	abilities.Use(middleware.AuthMiddleware())
	abilities.Use(middleware.AdminMiddleware())
	{
		abilities.GET("/abilities", listAbilities(client))
		abilities.POST("/abilities", createAbility(client))
		abilities.GET("/abilities/:id", getAbility(client))
		abilities.PUT("/abilities/:id", updateAbility(client))
		abilities.DELETE("/abilities/:id", deleteAbility(client))
	}
}

type abilityView struct {
	ID              int      `json:"id"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	AbilityType     string   `json:"ability_type"`
	Cost            int      `json:"cost"`
	Cooldown        int      `json:"cooldown"`
	Requirements    string   `json:"requirements"`
	ManaCost        int      `json:"mana_cost"`
	StaminaCost     int      `json:"stamina_cost"`
	HpCost          int      `json:"hp_cost"`
	Slug            string   `json:"slug"`
	RequiredTag     string   `json:"required_tag"`
	AbilityClass    string   `json:"ability_class"`
	ProcChance      float64  `json:"proc_chance"`
	ProcEvent       string   `json:"proc_event"`
	CooldownSeconds int      `json:"cooldown_seconds"`
	Tags            []string `json:"tags"`
	FactionID       *int     `json:"faction_id,omitempty"`
	FactionName     string   `json:"faction_name,omitempty"`
}

type abilityInput struct {
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	AbilityType     string   `json:"ability_type"`
	Cost            *int     `json:"cost"`
	Cooldown        *int     `json:"cooldown"`
	Requirements    string   `json:"requirements"`
	ManaCost        *int     `json:"mana_cost"`
	StaminaCost     *int     `json:"stamina_cost"`
	HpCost          *int     `json:"hp_cost"`
	Slug            string   `json:"slug"`
	RequiredTag     string   `json:"required_tag"`
	AbilityClass    string   `json:"ability_class"`
	ProcChance      *float64 `json:"proc_chance"`
	ProcEvent       string   `json:"proc_event"`
	CooldownSeconds *int     `json:"cooldown_seconds"`
	Tags            []string `json:"tags"`
	FactionID       *int     `json:"faction_id"`
}

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

func buildRequirementsJSON(tags []string) string {
	if len(tags) == 0 {
		return `{}`
	}
	m := map[string][]string{"tags": tags}
	b, _ := json.Marshal(m)
	return string(b)
}

func abilityToView(s *db.Ability) abilityView {
	v := abilityView{
		ID:              s.ID,
		Name:            s.Name,
		Description:     s.Description,
		AbilityType:     s.AbilityType,
		Cost:            s.Cost,
		Cooldown:        s.Cooldown,
		Requirements:    s.Requirements,
		ManaCost:        s.ManaCost,
		StaminaCost:     s.StaminaCost,
		HpCost:          s.HpCost,
		Slug:            s.Slug,
		RequiredTag:     s.RequiredTag,
		AbilityClass:    s.AbilityClass,
		ProcChance:      s.ProcChance,
		ProcEvent:       s.ProcEvent,
		CooldownSeconds: s.CooldownSeconds,
		Tags:            parseTagsFromRequirements(s.Requirements),
	}
	if s.Edges.Faction != nil {
		v.FactionID = &s.Edges.Faction.ID
		v.FactionName = s.Edges.Faction.Name
	}
	return v
}

func listAbilities(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := client.Ability.Query().WithFaction()
		if t := c.Query("type"); t != "" {
			query = query.Where(ability.AbilityType(t))
		}
		if ac := c.Query("ability_class"); ac != "" {
			query = query.Where(ability.AbilityClass(ac))
		}
		abilities, err := query.Order(ability.ByName()).All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result := make([]abilityView, len(abilities))
		for i, s := range abilities {
			result[i] = abilityToView(s)
		}
		c.JSON(http.StatusOK, gin.H{"abilities": result})
	}
}

func getAbility(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ability id"})
			return
		}
		s, err := client.Ability.Query().Where(ability.ID(id)).WithFaction().Only(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "ability not found"})
			return
		}
		c.JSON(http.StatusOK, abilityToView(s))
	}
}

func createAbility(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input abilityInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if input.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
			return
		}
		mut := client.Ability.Create().
			SetName(input.Name).
			SetDescription(input.Description).
			SetAbilityType(input.AbilityType).
			SetRequirements(buildRequirementsJSON(input.Tags)).
			SetSlug(input.Slug).
			SetRequiredTag(input.RequiredTag).
			SetAbilityClass(input.AbilityClass).
			SetProcEvent(input.ProcEvent)
		if input.Cost != nil {
			mut.SetCost(*input.Cost)
		}
		if input.Cooldown != nil {
			mut.SetCooldown(*input.Cooldown)
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
		s, _ = client.Ability.Query().WithFaction().Where(ability.ID(s.ID)).Only(c.Request.Context())
		c.JSON(http.StatusCreated, abilityToView(s))
	}
}

func updateAbility(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ability id"})
			return
		}
		s, err := client.Ability.Query().Where(ability.ID(id)).Only(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "ability not found"})
			return
		}
		var input abilityInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		mut := client.Ability.UpdateOne(s).
			SetName(input.Name).
			SetDescription(input.Description).
			SetAbilityType(input.AbilityType).
			SetRequirements(buildRequirementsJSON(input.Tags)).
			SetSlug(input.Slug).
			SetRequiredTag(input.RequiredTag).
			SetAbilityClass(input.AbilityClass).
			SetProcEvent(input.ProcEvent)
		if input.Cost != nil {
			mut.SetCost(*input.Cost)
		}
		if input.Cooldown != nil {
			mut.SetCooldown(*input.Cooldown)
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
		updated, _ = client.Ability.Query().WithFaction().Where(ability.ID(updated.ID)).Only(c.Request.Context())
		c.JSON(http.StatusOK, abilityToView(updated))
	}
}

func deleteAbility(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ability id"})
			return
		}
		err = client.Ability.DeleteOneID(id).Exec(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "ability not found"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func listClasslessAbilities(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		abilities, err := client.Ability.Query().
			Where(ability.AbilityClassEQ("active")).
			Order(ability.ByName()).
			All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result := make([]gin.H, len(abilities))
		for i, a := range abilities {
			result[i] = gin.H{
				"id":           a.ID,
				"name":         a.Name,
				"description":  a.Description,
				"ability_type": a.AbilityType,
				"mana_cost":    a.ManaCost,
				"stamina_cost": a.StaminaCost,
				"hp_cost":      a.HpCost,
				"cooldown":     a.Cooldown,
			}
		}
		c.JSON(http.StatusOK, gin.H{"abilities": result})
	}
}