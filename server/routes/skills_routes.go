package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/skill"
)

// RegisterSkillRoutes registers all skill-related routes
func RegisterSkillRoutes(router *gin.Engine, client *db.Client) {
	// Create a new skill
	router.POST("/skills", func(c *gin.Context) {
		var req struct {
			ID             int    `json:"id"`
			Name           string `json:"name" binding:"required"`
			Description    string `json:"description"`
			SkillType      string `json:"skill_type"`
			Cost           int    `json:"cost"`
			Cooldown       int    `json:"cooldown"`
			Requirements   string `json:"requirements"`
			EffectType     string `json:"effect_type"`
			EffectValue    int    `json:"effect_value"`
			EffectDuration int    `json:"effect_duration"`
			ManaCost       int    `json:"mana_cost"`
			StaminaCost    int    `json:"stamina_cost"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		creator := client.Skill.Create().
			SetName(req.Name).
			SetDescription(req.Description).
			SetSkillType(req.SkillType).
			SetCost(req.Cost).
			SetCooldown(req.Cooldown).
			SetRequirements(req.Requirements).
			SetEffectType(req.EffectType).
			SetEffectValue(req.EffectValue).
			SetEffectDuration(req.EffectDuration).
			SetManaCost(req.ManaCost).
			SetStaminaCost(req.StaminaCost)

		sk, err := creator.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, sk)
	})

	// Get all skills
	router.GET("/skills", func(c *gin.Context) {
		skills, err := client.Skill.Query().All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"skills": skills,
			"count":  len(skills),
		})
	})

	// Get a single skill by ID
	router.GET("/skills/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skill ID"})
			return
		}

		sk, err := client.Skill.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Skill not found"})
			return
		}

		c.JSON(http.StatusOK, sk)
	})

	// Update a skill by ID
	router.PUT("/skills/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skill ID"})
			return
		}

		var req struct {
			Name           string `json:"name"`
			Description    string `json:"description"`
			SkillType      string `json:"skill_type"`
			Cost           *int   `json:"cost"`
			Cooldown       *int   `json:"cooldown"`
			Requirements   string `json:"requirements"`
			EffectType     string `json:"effect_type"`
			EffectValue    *int   `json:"effect_value"`
			EffectDuration *int   `json:"effect_duration"`
			ManaCost       *int   `json:"mana_cost"`
			StaminaCost    *int   `json:"stamina_cost"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updater := client.Skill.UpdateOneID(id)

		if req.Name != "" {
			updater.SetName(req.Name)
		}
		if req.Description != "" {
			updater.SetDescription(req.Description)
		}
		if req.SkillType != "" {
			updater.SetSkillType(req.SkillType)
		}
		if req.Cost != nil {
			updater.SetCost(*req.Cost)
		}
		if req.Cooldown != nil {
			updater.SetCooldown(*req.Cooldown)
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
		if req.ManaCost != nil {
			updater.SetManaCost(*req.ManaCost)
		}
		if req.StaminaCost != nil {
			updater.SetStaminaCost(*req.StaminaCost)
		}

		sk, err := updater.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Skill not found"})
			return
		}

		c.JSON(http.StatusOK, sk)
	})

	// Delete a skill by ID
	router.DELETE("/skills/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skill ID"})
			return
		}

		ctx := c.Request.Context()

		// Delete the skill (DB should cascade delete character associations if configured)
		err = client.Skill.DeleteOneID(id).Exec(ctx)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Skill not found"})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	})

	// Get skills by type
	router.GET("/skills/type/:type", func(c *gin.Context) {
		skType := c.Param("type")
		skills, err := client.Skill.Query().
			Where(skill.SkillType(skType)).
			All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"skills": skills,
			"count":  len(skills),
		})
	})

	// Get skills by effect type
	router.GET("/skills/effect/:effectType", func(c *gin.Context) {
		effectType := c.Param("effectType")
		skills, err := client.Skill.Query().
			Where(skill.EffectType(effectType)).
			All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"skills": skills,
			"count":  len(skills),
		})
	})
}
