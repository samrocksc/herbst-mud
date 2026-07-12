package routes

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/skill"
	"herbst-server/dblog"
	"herbst-server/middleware"
)

// RegisterSkillRoutes registers DB-backed skill CRUD endpoints.
func RegisterSkillRoutes(r *gin.Engine, client *db.Client) {
	skills := r.Group("/api")
	skills.Use(middleware.AuthMiddleware(nil))
	skills.Use(middleware.AdminMiddleware())
	skills.Use(middleware.WorldAccessMiddleware())
	{
		skills.GET("/skills", listSkills(client))
		skills.POST("/skills", createSkill(client))
		skills.GET("/skills/:id", getSkill(client))
		skills.PUT("/skills/:id", updateSkill(client))
		skills.DELETE("/skills/:id", deleteSkill(client))
	}
}

// --- input / output structs ---

type skillInput struct {
	WorldID      int                    `json:"world_id"`
	Name         string                 `json:"name"`
	DisplayName  string                 `json:"display_name"`
	Description  string                 `json:"description"`
	Category     string                 `json:"category"`
	MaxLevel     *int                   `json:"max_level"`
	XpCurveMode  string                 `json:"xp_curve_mode"`
	XpCurveData  map[string]interface{} `json:"xp_curve_data"`
	ParentSkillID *int                  `json:"parent_skill_id"`
}

// --- handlers ---

func listSkills(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := client.Skill.Query()
		if w := c.Query("world_id"); w != "" {
			wid, err := strconv.Atoi(w)
			if err == nil {
				query = query.Where(skill.WorldIDEQ(wid))
			}
		}
		if cat := c.Query("category"); cat != "" {
			query = query.Where(skill.CategoryEQ(cat))
		}
		skills, err := query.Order(skill.ByID()).All(c.Request.Context())
		if err != nil {
			dblog.Error("list skills failed", err, slog.String("service", "skills"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"skills": skills, "count": len(skills)})
	}
}

func getSkill(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid skill id"})
			return
		}
		s, err := client.Skill.Query().Where(skill.IDEQ(id)).Only(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "skill not found"})
			return
		}
		c.JSON(http.StatusOK, s)
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
		if input.DisplayName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "display_name is required"})
			return
		}
		if input.WorldID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "world_id is required"})
			return
		}

		builder := client.Skill.Create().
			SetWorldID(input.WorldID).
			SetName(input.Name).
			SetDisplayName(input.DisplayName).
			SetDescription(input.Description).
			SetCategory(input.Category).
			SetXpCurveMode(input.XpCurveMode).
			SetXpCurveData(input.XpCurveData)

		if input.MaxLevel != nil {
			builder = builder.SetMaxLevel(*input.MaxLevel)
		}
		if input.ParentSkillID != nil {
			builder = builder.SetParentSkillID(*input.ParentSkillID)
		}

		created, err := builder.Save(c.Request.Context())
		if err != nil {
			dblog.Error("create skill failed", err, slog.String("service", "skills"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		slog.Info("skill created", slog.Int("skill_id", created.ID), slog.String("name", created.Name), slog.String("service", "skills"))
		c.JSON(http.StatusCreated, created)
	}
}

func updateSkill(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid skill id"})
			return
		}
		// Check existence first
		existing, err := client.Skill.Query().Where(skill.IDEQ(id)).Only(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "skill not found"})
			return
		}
		var input skillInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		builder := client.Skill.UpdateOneID(existing.ID)
		if input.Name != "" {
			builder = builder.SetName(input.Name)
		}
		if input.DisplayName != "" {
			builder = builder.SetDisplayName(input.DisplayName)
		}
		if input.Description != "" {
			builder = builder.SetDescription(input.Description)
		}
		if input.Category != "" {
			builder = builder.SetCategory(input.Category)
		}
		if input.MaxLevel != nil {
			builder = builder.SetMaxLevel(*input.MaxLevel)
		}
		if input.XpCurveMode != "" {
			builder = builder.SetXpCurveMode(input.XpCurveMode)
		}
		if input.XpCurveData != nil {
			builder = builder.SetXpCurveData(input.XpCurveData)
		}
		if input.WorldID != 0 {
			builder = builder.SetWorldID(input.WorldID)
		}
		if input.ParentSkillID != nil {
			builder = builder.SetParentSkillID(*input.ParentSkillID)
		}

		updated, err := builder.Save(c.Request.Context())
		if err != nil {
			dblog.Error("update skill failed", err, slog.Int("skill_id", id), slog.String("service", "skills"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		slog.Info("skill updated", slog.Int("skill_id", id), slog.String("service", "skills"))
		c.JSON(http.StatusOK, updated)
	}
}

func deleteSkill(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid skill id"})
			return
		}
		err = client.Skill.DeleteOneID(id).Exec(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "skill not found"})
			return
		}
		slog.Info("skill deleted", slog.Int("skill_id", id), slog.String("service", "skills"))
		c.Status(http.StatusNoContent)
	}
}