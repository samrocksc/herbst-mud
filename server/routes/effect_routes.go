package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/ability"
	"herbst-server/db/abilityeffect"
	"herbst-server/middleware"
)

// RegisterEffectRoutes registers REST endpoints for ability effects.
// Protected /api routes — all require JWT auth + admin check
func RegisterEffectRoutes(r *gin.Engine, client *db.Client) {
	effects := r.Group("/api")
	effects.Use(middleware.AuthMiddleware())
	effects.Use(middleware.AdminMiddleware())
	{
		effects.GET("/abilities/:id/effects", listEffects(client))
		effects.POST("/abilities/:id/effects", createEffect(client))
		effects.PUT("/ability-effects/:id", updateEffect(client))
		effects.DELETE("/ability-effects/:id", deleteEffect(client))
	}
}

// effectView is the JSON shape returned by the API.
type effectView struct {
	ID            int     `json:"id"`
	AbilityID     int     `json:"ability_id"`
	EffectType    string  `json:"effect_type"`
	DamageSubtype string  `json:"damage_subtype"`
	Target        string  `json:"target"`
	Value         int     `json:"value"`
	Duration      int     `json:"duration"`
	ScalingStat   string  `json:"scaling_stat"`
	ScalingRatio  float64 `json:"scaling_ratio"`
	SortOrder     int     `json:"sort_order"`
}

// effectInput is the request body for create and update.
type effectInput struct {
	EffectType    string   `json:"effect_type"`
	DamageSubtype string   `json:"damage_subtype"`
	Target        string   `json:"target"`
	Value         *int     `json:"value"`
	Duration      *int     `json:"duration"`
	ScalingStat   string   `json:"scaling_stat"`
	ScalingRatio  *float64 `json:"scaling_ratio"`
	SortOrder     *int     `json:"sort_order"`
}

func effectToView(e *db.AbilityEffect) effectView {
	return effectView{
		ID:            e.ID,
		AbilityID:     e.Edges.Ability.ID,
		EffectType:    e.EffectType,
		DamageSubtype: e.DamageSubtype,
		Target:        e.Target,
		Value:         e.Value,
		Duration:      e.Duration,
		ScalingStat:   e.ScalingStat,
		ScalingRatio:  e.ScalingRatio,
		SortOrder:     e.SortOrder,
	}
}

func listEffects(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		abilityID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ability id"})
			return
		}

		// Verify ability exists
		_, err = client.Ability.Get(c.Request.Context(), abilityID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "ability not found"})
			return
		}

		effects, err := client.AbilityEffect.Query().
			Where(abilityeffect.HasAbilityWith(ability.ID(abilityID))).
			Order(abilityeffect.BySortOrder()).
			All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		result := make([]effectView, len(effects))
		for i, e := range effects {
			result[i] = effectToView(e)
		}
		c.JSON(http.StatusOK, gin.H{"effects": result})
	}
}

func createEffect(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		abilityID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ability id"})
			return
		}

		// Verify ability exists
		_, err = client.Ability.Get(c.Request.Context(), abilityID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "ability not found"})
			return
		}

		var input effectInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if input.EffectType == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "effect_type is required"})
			return
		}

		mut := client.AbilityEffect.Create().
			SetEffectType(input.EffectType).
			SetAbilityID(abilityID).
			SetDamageSubtype(input.DamageSubtype).
			SetTarget(input.Target)

		if input.Value != nil {
			mut.SetValue(*input.Value)
		}
		if input.Duration != nil {
			mut.SetDuration(*input.Duration)
		}
		if input.ScalingStat != "" {
			mut.SetScalingStat(input.ScalingStat)
		}
		if input.ScalingRatio != nil {
			mut.SetScalingRatio(*input.ScalingRatio)
		}
		if input.SortOrder != nil {
			mut.SetSortOrder(*input.SortOrder)
		}

		e, err := mut.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		e, _ = client.AbilityEffect.Query().
			Where(abilityeffect.ID(e.ID)).
			WithAbility().
			Only(c.Request.Context())
		c.JSON(http.StatusCreated, effectToView(e))
	}
}

func updateEffect(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid effect id"})
			return
		}

		e, err := client.AbilityEffect.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "effect not found"})
			return
		}

		var input effectInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		mut := client.AbilityEffect.UpdateOne(e).
			SetEffectType(input.EffectType).
			SetDamageSubtype(input.DamageSubtype).
			SetTarget(input.Target)

		if input.Value != nil {
			mut.SetValue(*input.Value)
		}
		if input.Duration != nil {
			mut.SetDuration(*input.Duration)
		}
		if input.ScalingStat != "" {
			mut.SetScalingStat(input.ScalingStat)
		}
		if input.ScalingRatio != nil {
			mut.SetScalingRatio(*input.ScalingRatio)
		}
		if input.SortOrder != nil {
			mut.SetSortOrder(*input.SortOrder)
		}

		updated, err := mut.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		updated, _ = client.AbilityEffect.Query().
			Where(abilityeffect.ID(updated.ID)).
			WithAbility().
			Only(c.Request.Context())
		c.JSON(http.StatusOK, effectToView(updated))
	}
}

func deleteEffect(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid effect id"})
			return
		}

		err = client.AbilityEffect.DeleteOneID(id).Exec(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "effect not found"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}