package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/skill"
	"herbst-server/middleware"
)

// RegisterSkillRoutes registers REST endpoints for skills.
func RegisterSkillRoutes(r *gin.Engine, client *db.Client) {
	// Protected /api routes — all require JWT auth + admin check
	skills := r.Group("/api")
	skills.Use(middleware.AuthMiddleware())
	skills.Use(middleware.AdminMiddleware())
	{
		skills.GET("/skills", listSkills(client))
		skills.PUT("/skills/:id", updateSkill(client))
	}
}

// ─── Skill CRUD ─────────────────────────────────────────────────────────────

// skillView is the JSON shape returned by the API.
type skillView struct {
	ID              int     `json:"id"`
	Name            string  `json:"name"`
	Slug            string  `json:"slug"`
	SkillClass      string  `json:"skill_class"`
	RequiredTag     string  `json:"required_tag"`
	ProcChance      float64 `json:"proc_chance"`
	ProcEvent       string  `json:"proc_event"`
	CooldownSeconds int     `json:"cooldown_seconds"`
}

func skillToView(s *db.Skill) skillView {
	return skillView{
		ID:              s.ID,
		Name:            s.Name,
		Slug:            s.Slug,
		SkillClass:      s.SkillClass,
		RequiredTag:     s.RequiredTag,
		ProcChance:      s.ProcChance,
		ProcEvent:       s.ProcEvent,
		CooldownSeconds: s.CooldownSeconds,
	}
}

func listSkills(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		skills, err := client.Skill.Query().
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
		c.JSON(http.StatusOK, result)
	}
}

func updateSkill(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid skill id"})
			return
		}

		var req struct {
			SkillClass      *string  `json:"skill_class"`
			RequiredTag     *string  `json:"required_tag"`
			ProcChance      *float64 `json:"proc_chance"`
			ProcEvent       *string  `json:"proc_event"`
			CooldownSeconds *int     `json:"cooldown_seconds"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		update := client.Skill.UpdateOneID(id)
		if req.SkillClass != nil {
			update.SetNillableSkillClass(req.SkillClass)
		}
		if req.RequiredTag != nil {
			if *req.RequiredTag == "" {
				update.ClearRequiredTag()
			} else {
				update.SetRequiredTag(*req.RequiredTag)
			}
		}
		if req.ProcChance != nil {
			update.SetNillableProcChance(req.ProcChance)
		}
		if req.ProcEvent != nil {
			update.SetNillableProcEvent(req.ProcEvent)
		}
		if req.CooldownSeconds != nil {
			update.SetNillableCooldownSeconds(req.CooldownSeconds)
		}

		updated, err := update.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Reload to get the latest state
		reloaded, err := client.Skill.Query().
			Where(skill.ID(updated.ID)).
			Only(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, skillToView(reloaded))
	}
}