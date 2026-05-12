package routes

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/ability"
	"herbst-server/middleware"
	"herbst-server/repository"
)

func RegisterAbilityRoutes(r *gin.Engine, repos *repository.Container, client *db.Client) {
	r.GET("/abilities/classless", listClasslessAbilities(repos))

	abilities := r.Group("/api")
	abilities.Use(middleware.AuthMiddleware())
	abilities.Use(middleware.AdminMiddleware())
	{
		abilities.GET("/abilities", listAbilities(client))
		abilities.POST("/abilities", createAbility(repos, client))
		abilities.GET("/abilities/:id", getAbility(repos, client))
		abilities.PUT("/abilities/:id", updateAbility(repos, client))
		abilities.DELETE("/abilities/:id", deleteAbility(repos))
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
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	AbilityType     string    `json:"ability_type"`
	Cost            *int      `json:"cost"`
	Cooldown        *int      `json:"cooldown"`
	Requirements    string    `json:"requirements"`
	ManaCost        *int      `json:"mana_cost"`
	StaminaCost     *int      `json:"stamina_cost"`
	HpCost          *int      `json:"hp_cost"`
	Slug            string    `json:"slug"`
	RequiredTag     string    `json:"required_tag"`
	AbilityClass    string    `json:"ability_class"`
	ProcChance      *float64  `json:"proc_chance"`
	ProcEvent       string    `json:"proc_event"`
	CooldownSeconds *int      `json:"cooldown_seconds"`
	Tags            []string  `json:"tags"`
	FactionID       *int      `json:"faction_id"`
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

// TODO: Migrate to repos.Ability.List once repo supports filtering and WithFaction edge loading
func listAbilities(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := client.Ability.Query().WithFaction()
		if t := c.Query("type"); t != "" {
			query = query.Where(ability.AbilityType(t))
		}
		if ac := c.Query("ability_class"); ac != "" {
			query = query.Where(ability.AbilityClass(ac))
		}
		if s := c.Query("search"); s != "" {
			query = query.Where(ability.NameContains(s))
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

func getAbility(repos *repository.Container, client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ability id"})
			return
		}
		// TODO: Migrate to repos.Ability.Get once repo supports WithFaction edge loading
		s, err := client.Ability.Query().Where(ability.ID(id)).WithFaction().Only(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "ability not found"})
			return
		}
		c.JSON(http.StatusOK, abilityToView(s))
	}
}

func createAbility(repos *repository.Container, client *db.Client) gin.HandlerFunc {
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
		s, err := repos.Ability.Create(c.Request.Context(), repository.CreateAbilityInput{
			Name:            input.Name,
			Description:     input.Description,
			AbilityType:     input.AbilityType,
			Cost:            derefInt(input.Cost),
			Cooldown:        derefInt(input.Cooldown),
			ManaCost:        derefInt(input.ManaCost),
			StaminaCost:     derefInt(input.StaminaCost),
			HPCost:          derefInt(input.HpCost),
			Requirements:    buildRequirementsJSON(input.Tags),
			RequiredTag:     input.RequiredTag,
			ProcChance:      derefFloat64(input.ProcChance),
			ProcEvent:       input.ProcEvent,
			CooldownSeconds: derefInt(input.CooldownSeconds),
			Slug:            input.Slug,
			AbilityClass:    input.AbilityClass,
			FactionID:       input.FactionID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Reload with faction edge for response.
		s, _ = client.Ability.Query().WithFaction().Where(ability.ID(s.ID)).Only(c.Request.Context())
		c.JSON(http.StatusCreated, abilityToView(s))
	}
}

func derefInt(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}

func derefFloat64(p *float64) float64 {
	if p == nil {
		return 0
	}
	return *p
}

func updateAbility(repos *repository.Container, client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ability id"})
			return
		}
		_, err = repos.Ability.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "ability not found"})
			return
		}
		var input abilityInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		_, err = repos.Ability.Update(c.Request.Context(), id, repository.AbilityUpdates{
			Name:            &input.Name,
			Description:     &input.Description,
			AbilityType:    &input.AbilityType,
			AbilityClass:   &input.AbilityClass,
			Cost:            input.Cost,
			Cooldown:        input.Cooldown,
			ManaCost:        input.ManaCost,
			StaminaCost:     input.StaminaCost,
			HPCost:          input.HpCost,
			Requirements:    strPtr(buildRequirementsJSON(input.Tags)),
			RequiredTag:     &input.RequiredTag,
			ProcChance:      input.ProcChance,
			ProcEvent:       &input.ProcEvent,
			CooldownSeconds: input.CooldownSeconds,
			Slug:            &input.Slug,
			FactionID:       input.FactionID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Reload with faction edge for response.
		updated, _ := client.Ability.Query().WithFaction().Where(ability.ID(id)).Only(c.Request.Context())
		c.JSON(http.StatusOK, abilityToView(updated))
	}
}


func deleteAbility(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ability id"})
			return
		}
		if err := repos.Ability.Delete(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "ability not found"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func listClasslessAbilities(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		abilities, err := repos.Ability.ListByClass(c.Request.Context(), "active")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		sort.Slice(abilities, func(i, j int) bool {
			return abilities[i].Name < abilities[j].Name
		})
		result := make([]gin.H, len(abilities))
		for i, a := range abilities {
			result[i] = gin.H{
				"id":           a.ID,
				"name":         a.Name,
				"description":  a.Description,
				"ability_type": a.AbilityType,
				"mana_cost":    a.ManaCost,
				"stamina_cost":  a.StaminaCost,
				"hp_cost":      a.HpCost,
				"cooldown":     a.Cooldown,
			}
		}
		c.JSON(http.StatusOK, gin.H{"abilities": result})
	}
}