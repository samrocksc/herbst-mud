package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/talent"
)

// RegisterTalentRoutes registers all talent-related routes
func RegisterTalentRoutes(router *gin.Engine, client *db.Client) {
	// Get all talents
	router.GET("/talents", func(c *gin.Context) {
		talents, err := client.Talent.Query().All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"talents": talents,
			"count":  len(talents),
		})
	})

	// Get a single talent by ID
	router.GET("/talents/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid talent ID"})
			return
		}

		t, err := client.Talent.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Talent not found"})
			return
		}

		c.JSON(http.StatusOK, t)
	})

	// Create a new talent
	router.POST("/talents", func(c *gin.Context) {
		var req struct {
			ID              int    `json:"id"`
			Name            string `json:"name" binding:"required"`
			Description     string `json:"description"`
			Requirements    string `json:"requirements"`
			EffectType      string `json:"effect_type"`
			EffectValue     int    `json:"effect_value"`
			EffectDuration  int    `json:"effect_duration"`
			Cooldown        int    `json:"cooldown"`
			ManaCost        int    `json:"mana_cost"`
			StaminaCost     int    `json:"stamina_cost"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		creator := client.Talent.Create().
			SetName(req.Name).
			SetDescription(req.Description).
			SetRequirements(req.Requirements).
			SetEffectType(req.EffectType).
			SetEffectValue(req.EffectValue).
			SetEffectDuration(req.EffectDuration).
			SetCooldown(req.Cooldown).
			SetManaCost(req.ManaCost).
			SetStaminaCost(req.StaminaCost)

		t, err := creator.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, t)
	})

	// Update a talent by ID
	router.PUT("/talents/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid talent ID"})
			return
		}

		var req struct {
			Name            string `json:"name"`
			Description     string `json:"description"`
			Requirements    string `json:"requirements"`
			EffectType      string `json:"effect_type"`
			EffectValue     *int   `json:"effect_value"`
			EffectDuration  *int   `json:"effect_duration"`
			Cooldown        *int   `json:"cooldown"`
			ManaCost        *int   `json:"mana_cost"`
			StaminaCost     *int   `json:"stamina_cost"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updater := client.Talent.UpdateOneID(id)

		if req.Name != "" {
			updater.SetName(req.Name)
		}
		if req.Description != "" {
			updater.SetDescription(req.Description)
		}
		if req.Requirements != "" {
			updater.SetRequirements(req.Requirements)
		}
		if req.EffectType != "" {
			updater.SetEffectType(req.EffectType)
		}
		if req.EffectValue != nil {
			updater.SetEffectValue(*req.EffectValue)
		}
		if req.EffectDuration != nil {
			updater.SetEffectDuration(*req.EffectDuration)
		}
		if req.Cooldown != nil {
			updater.SetCooldown(*req.Cooldown)
		}
		if req.ManaCost != nil {
			updater.SetManaCost(*req.ManaCost)
		}
		if req.StaminaCost != nil {
			updater.SetStaminaCost(*req.StaminaCost)
		}

		t, err := updater.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Talent not found"})
			return
		}

		c.JSON(http.StatusOK, t)
	})

	// Delete a talent by ID
	router.DELETE("/talents/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid talent ID"})
			return
		}

		ctx := c.Request.Context()

		err = client.Talent.DeleteOneID(id).Exec(ctx)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Talent not found"})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	})

	// Get talents by effect type
	router.GET("/talents/effect/:effectType", func(c *gin.Context) {
		effectType := c.Param("effectType")
		talents, err := client.Talent.Query().
			Where(talent.EffectType(effectType)).
			All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"talents": talents,
			"count":  len(talents),
		})
	})
}
