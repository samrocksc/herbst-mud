package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/repository"
	"herbst-server/service"
)

// RegisterCharacterRoutes registers all character-related routes.
func RegisterCharacterRoutes(router *gin.Engine, svc *service.Container, repos *repository.Container) {
	chars := router.Group("/characters")
	{
		chars.POST("", createCharacter(svc, repos))
		chars.GET("", listCharacters(repos))
		chars.GET("/:id", getCharacter(repos))
		chars.PUT("/:id", updateCharacter(repos))
		chars.DELETE("/:id", deleteCharacter(repos))
				chars.GET("/:id/class", getCharacterClass(repos))
		chars.PUT("/:id/class", updateCharacterClass(repos))
		chars.GET("/:id/specialty", getCharacterSpecialty(repos))
		chars.PUT("/:id/specialty", updateCharacterSpecialty(repos))
		chars.GET("/:id/stats", getCharacterStats(repos))
		chars.PUT("/:id/stats", updateCharacterStats(repos))
		chars.GET("/:id/skills", getCharacterSkills(repos, svc))
		chars.PUT("/:id/skills", updateCharacterSkills(repos))
		chars.GET("/:id/combat-status", getCombatStatus(svc))
		chars.GET("/:id/abilities", getCharacterAbilities(svc))
		chars.POST("/:id/abilities", equipAbility(svc))
		chars.DELETE("/:id/abilities/:slot", unequipAbility(svc))
		chars.PUT("/:id/abilities/swap", swapAbilities(svc))
		chars.GET("/:id/passive-abilities", listPassiveAbilities(svc))
		chars.POST("/:id/passive-abilities", unlockPassiveAbility(svc))
		chars.DELETE("/:id/passive-abilities/:abilityId", removePassiveAbility(svc))
		chars.GET("/:id/classless-skills", getClasslessSkills(svc))
		chars.POST("/:id/classless-skills", equipClasslessSkill(svc))
		chars.PUT("/:id/classless-skills/swap", swapClasslessSkills(svc))
	}
	// NPC routes
	router.GET("/npcs/room/:id", getNPCsByRoom(repos))
	router.GET("/npcs", listAllNPCs(repos))
	// Character combat routes
	router.POST("/characters/:id/damage", applyDamage(svc))
	router.POST("/characters/:id/heal", healCharacter(svc))
	router.POST("/characters/:id/stamina", adjustStamina(svc))
	router.POST("/characters/:id/mana", adjustMana(svc))
	// NPC heal routes
	router.POST("/rooms/:id/npcs/heal", healNPCsInRoom(svc))
	router.POST("/rooms/:id/npcs/passive-heal", passiveHealNPCsInRoom(svc))
	// User-character routes
	router.GET("/user-characters/:id", getUserCharacters(repos))
	router.POST("/user-characters/:id", createCharacterForUser(svc, repos))
	router.GET("/user-characters/:id/needed", needsCharacter(repos))
	// Class/specialty lookup
	router.GET("/classes/:class/specialties", getSpecialtiesForClass)
	// Race/gender routes (public)
	router.GET("/races", listPlayableRaces(repos))
	router.GET("/genders", listGenders(repos))
	router.GET("/characters/:id/race", getCharacterRace(repos))
	router.PUT("/characters/:id/race", updateCharacterRace(repos))
	// Game config routes
	router.GET("/game-config/:key", getGameConfig(repos))
	router.PUT("/game-config/:key", setGameConfig(repos))
}

func getIDParam(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return 0, false
	}
	return id, true
}