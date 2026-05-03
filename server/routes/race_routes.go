package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/db/race"
	"herbst-server/middleware"
)

// RegisterRaceRoutes registers REST endpoints for races.
func RegisterRaceRoutes(r *gin.Engine, client *db.Client) {
	races := r.Group("/api/races")
	races.Use(middleware.AuthMiddleware())
	races.Use(middleware.AdminMiddleware())
	{
		races.GET("", listRaces(client))
		races.GET("/:id", getRace(client))
		races.POST("", createRace(client))
		races.PUT("/:id", updateRace(client))
		races.DELETE("/:id", deleteRace(client))
	}
}

// raceView is the JSON shape returned by the API.
type raceView struct {
	ID             int      `json:"id"`
	Name           string   `json:"name"`
	DisplayName    string   `json:"display_name"`
	Description    string   `json:"description,omitempty"`
	StatModifiers  any      `json:"stat_modifiers,omitempty"`
	SkillGrants   []string `json:"skill_grants,omitempty"`
	AbilityModifiers []string `json:"ability_modifiers,omitempty"`
	IsPlayable    bool     `json:"is_playable"`
	Color         string   `json:"color,omitempty"`
}

// parseJSON safely parses a JSON string into a value. Returns nil on error.
func parseJSON[T any](v *T) any {
	if v == nil {
		return nil
	}
	var out any
	data, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	if err := json.Unmarshal(data, &out); err != nil {
		return nil
	}
	return out
}

// listRaces returns all races.
func listRaces(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		races, err := client.Race.Query().Order(race.ByDisplayName()).All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		views := make([]raceView, len(races))
		for i, r := range races {
			views[i] = raceToView(r)
		}
		c.JSON(http.StatusOK, gin.H{"races": views})
	}
}

// getRace returns a single race by ID.
func getRace(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid race id"})
			return
		}
		r, err := client.Race.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "race not found"})
			return
		}
		c.JSON(http.StatusOK, raceToView(r))
	}
}

// createRace creates a new race.
func createRace(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name             string   `json:"name" binding:"required"`
			DisplayName      string   `json:"display_name"`
			Description      string   `json:"description"`
			StatModifiers    *string  `json:"stat_modifiers"`
			SkillGrants     []string `json:"skill_grants"`
			AbilityModifiers []string `json:"ability_modifiers"`
			IsPlayable      *bool    `json:"is_playable"`
			Color           string   `json:"color"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
			return
		}

		// Check for duplicate name
		existing, err := client.Race.Query().Where(race.NameEQ(req.Name)).Only(c.Request.Context())
		if err == nil && existing != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "a race with this name already exists"})
			return
		}

		isPlayable := true
		if req.IsPlayable != nil {
			isPlayable = *req.IsPlayable
		}
		displayName := req.DisplayName
		if displayName == "" {
			displayName = req.Name
		}

		mut := client.Race.Create().
			SetName(req.Name).
			SetDisplayName(displayName).
			SetDescription(req.Description).
			SetIsPlayable(isPlayable).
			SetColor(req.Color)

		if req.StatModifiers != nil {
			mut = mut.SetStatModifiers(*req.StatModifiers)
		}

		r, err := mut.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, raceToView(r))
	}
}

// updateRace updates an existing race.
func updateRace(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid race id"})
			return
		}

		var req struct {
			Name             *string  `json:"name"`
			DisplayName      *string  `json:"display_name"`
			Description      *string  `json:"description"`
			StatModifiers    *string  `json:"stat_modifiers"`
			SkillGrants     []string `json:"skill_grants"`
			AbilityModifiers []string `json:"ability_modifiers"`
			IsPlayable      *bool    `json:"is_playable"`
			Color           *string  `json:"color"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		// Check name uniqueness if name is being changed
		if req.Name != nil {
			existing, err := client.Race.Query().Where(race.NameEQ(*req.Name)).Only(c.Request.Context())
			if err == nil && existing != nil && existing.ID != id {
				c.JSON(http.StatusConflict, gin.H{"error": "a race with this name already exists"})
				return
			}
		}

		mut := client.Race.UpdateOneID(id)
		if req.Name != nil {
			mut = mut.SetName(*req.Name)
		}
		if req.DisplayName != nil {
			mut = mut.SetDisplayName(*req.DisplayName)
		}
		if req.Description != nil {
			mut = mut.SetDescription(*req.Description)
		}
		if req.StatModifiers != nil {
			mut = mut.SetStatModifiers(*req.StatModifiers)
		}
		if req.IsPlayable != nil {
			mut = mut.SetIsPlayable(*req.IsPlayable)
		}
		if req.Color != nil {
			mut = mut.SetColor(*req.Color)
		}

		r, err := mut.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "race not found or update failed"})
			return
		}
		c.JSON(http.StatusOK, raceToView(r))
	}
}

// deleteRace deletes a race by ID.
func deleteRace(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid race id"})
			return
		}

		// Check if any characters are using this race by name
		raceName := client.Race.Query().Where(race.ID(id)).OnlyX(c.Request.Context()).Name
		count, err := client.Character.Query().
			Where(character.RaceEQ(raceName)).
			Count(c.Request.Context())
		if err == nil && count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "cannot delete: race is in use by characters"})
			return
		}

		err = client.Race.DeleteOneID(id).Exec(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "race not found"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

// raceToView converts a Race ent model to a raceView.
func raceToView(r *db.Race) raceView {
	var statMod any
	if r.StatModifiers != "" {
		_ = json.Unmarshal([]byte(r.StatModifiers), &statMod)
	}

	// Parse skill_grants JSON array
	var skillGrants []string
	if r.SkillGrants != "" {
		_ = json.Unmarshal([]byte(r.SkillGrants), &skillGrants)
	}

	return raceView{
		ID:              r.ID,
		Name:            r.Name,
		DisplayName:     r.DisplayName,
		Description:     r.Description,
		StatModifiers:   statMod,
		SkillGrants:    skillGrants,
		AbilityModifiers: nil,
		IsPlayable:     r.IsPlayable,
		Color:          r.Color,
	}
}
