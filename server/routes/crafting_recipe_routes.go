package routes

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/dblog"
	"herbst-server/middleware"
	"herbst-server/repository"
)

func RegisterCraftingRecipeRoutes(r *gin.Engine, repos *repository.Container, client *db.Client) {
	recipes := r.Group("/api")
	recipes.Use(middleware.AuthMiddleware(nil))
	recipes.Use(middleware.AdminMiddleware())
	recipes.Use(middleware.WorldAccessMiddleware())
	{
		recipes.GET("/recipes", listRecipes(repos))
		recipes.GET("/recipes/:name", getRecipe(repos))
		recipes.POST("/recipes", createRecipe(repos))
		recipes.PUT("/recipes/:name", updateRecipe(repos))
		recipes.DELETE("/recipes/:name", deleteRecipe(repos))
	}
}

type recipeInput struct {
	Name               string                   `json:"name"`
	DisplayName        string                   `json:"display_name"`
	Description        string                   `json:"description"`
	RequiredStationTag string                   `json:"required_station_tag"`
	RequiredClass      string                   `json:"required_class"`
	RequiredSkillLevel *int                     `json:"required_skill_level"`
	RequiredSkill      string                   `json:"required_skill"`
	Inputs             []map[string]any         `json:"inputs"`
	Outputs            []map[string]any         `json:"outputs"`
	CraftTimeSecs      *int                     `json:"craft_time_secs"`
	WorldID            string                   `json:"world_id"`
}

func recipeToMap(r *db.CraftingRecipe) gin.H {
	return gin.H{
		"name":                 r.Name,
		"display_name":         r.DisplayName,
		"description":          r.Description,
		"required_station_tag": r.RequiredStationTag,
		"required_class":       r.RequiredClass,
		"required_skill_level": r.RequiredSkillLevel,
		"required_skill":       r.RequiredSkill,
		"inputs":               r.Inputs,
		"outputs":              r.Outputs,
		"craft_time_secs":      r.CraftTimeSecs,
		"world_id":             r.WorldID,
	}
}

func listRecipes(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		worldID := c.Query("world_id")
		stationTag := c.Query("station_tag")
		recipes, err := repos.CraftingRecipe.List(c.Request.Context(), worldID, stationTag)
		if err != nil {
			dblog.Error("failed to list recipes", err, slog.String("service", "crafting"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result := make([]gin.H, len(recipes))
		for i, r := range recipes {
			result[i] = recipeToMap(r)
		}
		c.JSON(http.StatusOK, gin.H{"recipes": result})
	}
}

func getRecipe(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		r, err := repos.CraftingRecipe.Get(c.Request.Context(), name)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
			return
		}
		c.JSON(http.StatusOK, recipeToMap(r))
	}
}

func createRecipe(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input recipeInput
		if err := c.ShouldBindJSON(&input); err != nil {
			slog.Warn("bad request", slog.String("service", "crafting"), slog.String("reason", "invalid json"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if input.Name == "" {
			slog.Warn("bad request", slog.String("service", "crafting"), slog.String("reason", "name is required"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
			return
		}
		r, err := repos.CraftingRecipe.Create(c.Request.Context(), repository.CreateCraftingRecipeInput{
			Name:               input.Name,
			DisplayName:        input.DisplayName,
			Description:        input.Description,
			RequiredStationTag: input.RequiredStationTag,
			RequiredClass:      input.RequiredClass,
			RequiredSkillLevel: derefInt(input.RequiredSkillLevel),
			RequiredSkill:      input.RequiredSkill,
			Inputs:             inputsFromInterface(input.Inputs),
			Outputs:            outputsFromInterface(input.Outputs),
			CraftTimeSecs:      derefInt(input.CraftTimeSecs),
			WorldID:            input.WorldID,
		})
		if err != nil {
			dblog.Error("failed to create recipe", err, slog.String("service", "crafting"), slog.String("recipe_name", input.Name))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		slog.Info("recipe created", slog.String("recipe_name", r.Name), slog.String("user_email", c.GetString("email")), slog.String("service", "crafting"))
		c.JSON(http.StatusCreated, recipeToMap(r))
	}
}

func updateRecipe(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		_, err := repos.CraftingRecipe.Get(c.Request.Context(), name)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
			return
		}
		var input recipeInput
		if err := c.ShouldBindJSON(&input); err != nil {
			slog.Warn("bad request", slog.String("service", "crafting"), slog.String("reason", "invalid json"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		updates := repository.CraftingRecipeUpdates{}
		if input.DisplayName != "" {
			updates.DisplayName = &input.DisplayName
		}
		if input.Description != "" {
			updates.Description = &input.Description
		}
		if input.RequiredStationTag != "" {
			updates.RequiredStationTag = &input.RequiredStationTag
		}
		if input.RequiredClass != "" {
			updates.RequiredClass = &input.RequiredClass
		}
		if input.RequiredSkillLevel != nil {
			updates.RequiredSkillLevel = input.RequiredSkillLevel
		}
		if input.RequiredSkill != "" {
			updates.RequiredSkill = &input.RequiredSkill
		}
		if len(input.Inputs) > 0 {
			inputs := inputsFromInterface(input.Inputs)
			updates.Inputs = &inputs
		}
		if len(input.Outputs) > 0 {
			outputs := outputsFromInterface(input.Outputs)
			updates.Outputs = &outputs
		}
		if input.CraftTimeSecs != nil {
			updates.CraftTimeSecs = input.CraftTimeSecs
		}
		if input.WorldID != "" {
			updates.WorldID = &input.WorldID
		}
		updated, err := repos.CraftingRecipe.Update(c.Request.Context(), name, updates)
		if err != nil {
			dblog.Error("failed to update recipe", err, slog.String("service", "crafting"), slog.String("recipe_name", name))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		slog.Info("recipe updated", slog.String("recipe_name", name), slog.String("user_email", c.GetString("email")), slog.String("service", "crafting"))
		c.JSON(http.StatusOK, recipeToMap(updated))
	}
}

func deleteRecipe(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		err := repos.CraftingRecipe.Delete(c.Request.Context(), name)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
			return
		}
		slog.Info("recipe deleted", slog.String("recipe_name", name), slog.String("user_email", c.GetString("email")), slog.String("service", "crafting"))
		c.Status(http.StatusNoContent)
	}
}