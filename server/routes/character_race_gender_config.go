package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/repository"
)

// listPlayableRaces handles GET /races (public, playable races only).
func listPlayableRaces(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		races, err := repos.Race.List(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		playable := make([]*db.Race, 0, len(races))
		for _, r := range races {
			if r.IsPlayable {
				playable = append(playable, r)
			}
		}
		result := make([]gin.H, len(playable))
		for i, r := range playable {
			result[i] = gin.H{
				"name":            r.Name,
				"display_name":    r.DisplayName,
				"description":     r.Description,
				"stat_modifiers":  r.StatModifiers,
				"skill_grants":    r.SkillGrants,
				"equipment_slots": r.EquipmentSlots,
			}
		}
		c.JSON(http.StatusOK, result)
	}
}

// listGenders handles GET /genders.
func listGenders(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		genders, err := repos.Gender.List(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result := make([]gin.H, len(genders))
		for i, g := range genders {
			result[i] = gin.H{
				"name":               g.Name,
				"display_name":       g.DisplayName,
				"subject_pronoun":    g.SubjectPronoun,
				"object_pronoun":     g.ObjectPronoun,
				"possessive_pronoun": g.PossessivePronoun,
			}
		}
		c.JSON(http.StatusOK, result)
	}
}

// getCharacterRace handles GET /characters/:id/race.
func getCharacterRace(repos *repository.Container) gin.HandlerFunc {
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
		c.JSON(http.StatusOK, gin.H{"id": char.ID, "race": char.Race})
	}
}

// updateCharacterRace handles PUT /characters/:id/race.
func updateCharacterRace(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := getIDParam(c)
		if !ok {
			return
		}
		var req struct {
			Race string `json:"race" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		existingRace, err := repos.Race.GetByName(c.Request.Context(), req.Race)
		if err != nil || !existingRace.IsPlayable {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or non-playable race"})
			return
		}
		char, err := repos.Character.Update(c.Request.Context(), id, repository.CharacterUpdates{Race: &req.Race})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": char.ID, "name": char.Name, "race": char.Race})
	}
}

// getGameConfig handles GET /game-config/:key.
func getGameConfig(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Param("key")
		cfg, err := repos.GameConfig.Get(c.Request.Context(), key)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Config key not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"key": cfg.Key, "value": cfg.Value})
	}
}

// setGameConfig handles PUT /game-config/:key.
func setGameConfig(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Param("key")
		var req struct {
			Value string `json:"value" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		cfg, err := repos.GameConfig.Set(c.Request.Context(), key, req.Value)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"key": cfg.Key, "value": cfg.Value})
	}
}