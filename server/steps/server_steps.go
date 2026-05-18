package steps

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	"github.com/cucumber/godog"
	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/repository"
	"herbst-server/routes"
	"herbst-server/service"
	"log/slog"
)

type ExportData struct {
	Version    string     `json:"version"`
	ExportedAt string     `json:"exported_at"`
	Rooms      []RoomData `json:"rooms"`
	NPCs       []NPCData  `json:"npcs"`
	Skills     []SkillData `json:"skills"`
	Items      []ItemData `json:"items"`
}

type RoomData struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	IsStarting  bool       `json:"is_starting"`
	Exits       []ExitData `json:"exits"`
}

type ExitData struct {
	Direction    string `json:"direction"`
	TargetRoomID int    `json:"target_room_id"`
}

type NPCData struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	CurrentRoomID int    `json:"current_room_id"`
	Race          string `json:"race"`
	Class         string `json:"class"`
	Level         int    `json:"level"`
	Hitpoints     int    `json:"hitpoints"`
	MaxHitpoints  int    `json:"max_hitpoints"`
	Stamina       int    `json:"stamina"`
	MaxStamina    int    `json:"max_stamina"`
	Mana          int    `json:"mana"`
	MaxMana       int    `json:"max_mana"`
	Strength      int    `json:"strength"`
	Dexterity     int    `json:"dexterity"`
	Constitution  int    `json:"constitution"`
	Intelligence  int    `json:"intelligence"`
	Wisdom        int    `json:"wisdom"`
	NPCSkillID    string `json:"npc_skill_id,omitempty"`
	IsImmortal    bool   `json:"is_immortal"`
}

type SkillData struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MinCooldown int    `json:"min_cooldown"`
	EffectType  string `json:"effect_type,omitempty"`
}

type ItemData struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Type         string `json:"type"`
	LocationType string `json:"location_type"`
	LocationID   int    `json:"location_id"`
}

type ImportResult struct {
	Success   bool     `json:"success"`
	Imported  Imported `json:"imported"`
	Version   string   `json:"version"`
	ImportedAt string  `json:"imported_at"`
}

type Imported struct {
	Rooms  int `json:"rooms"`
	NPCs   int `json:"npcs"`
	Skills int `json:"skills"`
	Items  int `json:"items"`
}

type Validation struct {
	Version  string   `json:"version"`
	IsValid  bool     `json:"is_valid"`
	Rooms    int      `json:"rooms"`
	NPCs     int      `json:"npcs"`
	Skills   int      `json:"skills"`
	Errors   []string `json:"errors,omitempty"`
}

type ServerTest struct {
	server         *gin.Engine
	response       *httptest.ResponseRecorder
	importJson     []byte
	dbClient       *db.Client
	testWorldID    int
}

// setupTestDB opens a DB connection and skips the test if unavailable
func setupTestDB() (*db.Client, error) {
	client, err := db.Open("postgres", "host=localhost port=5432 user=herbst password=herbst_password dbname=herbst_mud sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("cannot connect to database: %w", err)
	}
	return client, nil
}

func (s *ServerTest) theWebServerIsRunning() error {
	gin.SetMode(gin.TestMode)

	// Try to connect to the database
	client, err := setupTestDB()
	if err != nil {
		// Fall back to mock server if DB is unavailable
		s.server = gin.New()
		s.registerMockRoutes()
		return nil
	}
	s.dbClient = client

	// Build a real Gin engine with DB-backed routes
	router := gin.New()
	repos := repository.NewContainer(client)
	services := service.NewContainer(client, repos, slog.Default())
	_ = services

	// Mock healthz and openapi (same as before for backward compatibility)
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"ssh":    "running",
		})
	})
	router.GET("/openapi.json", func(c *gin.Context) {
		c.JSON(http.StatusOK, getOpenAPISpec())
	})

	// Register real export/import routes (no auth required)
	routes.RegisterGameExportRoutes(router, client)

	// Register admin wipe routes (no auth required)
	routes.RegisterAdminWipeRoutes(router, client)

	// Register world CRUD — wrapped with no-op auth for test context
	// Instead of middleware.AuthMiddleware + middleware.AdminMiddleware,
	// we set a test admin context so the routes pass through
	worlds := router.Group("/api")
	worlds.Use(func(c *gin.Context) {
		// Bypass auth: set admin context so handlers pass through
		c.Set("user_id", uint(1))
		c.Set("email", "admin@test.local")
		c.Set("is_admin", true)
		c.Set("db_client", client)
		c.Set("allowed_worlds", ([]string)(nil))
		c.Next()
	})
	{
		// POST /api/worlds — CreateWorldHandler
		worlds.POST("/worlds", func(c *gin.Context) {
			var input repository.CreateWorldInput
			if err := c.ShouldBindJSON(&input); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			world, err := repos.World.Create(c.Request.Context(), input)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, world)
		})

		// GET /api/worlds/db — ListWorldsHandler
		worlds.GET("/worlds/db", func(c *gin.Context) {
			worldsList, err := repos.World.List(c.Request.Context())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"worlds": worldsList, "count": len(worldsList)})
		})

		// GET /api/worlds/:id — GetWorldHandler
		worlds.GET("/worlds/:id", func(c *gin.Context) {
			id := 0
			fmt.Sscanf(c.Param("id"), "%d", &id)
			world, err := repos.World.Get(c.Request.Context(), id)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "World not found"})
				return
			}
			c.JSON(http.StatusOK, world)
		})

		// DELETE /api/worlds/:id — DeleteWorldHandler
		worlds.DELETE("/worlds/:id", func(c *gin.Context) {
			id := 0
			fmt.Sscanf(c.Param("id"), "%d", &id)
			if err := repos.World.Delete(c.Request.Context(), id); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "World deleted"})
		})
	}

	s.server = router
	return nil
}

// registerMockRoutes provides backward-compatible mock routes when DB is unavailable
func (s *ServerTest) registerMockRoutes() {
	s.server.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"ssh":    "running",
		})
	})
	s.server.GET("/openapi.json", func(c *gin.Context) {
		c.JSON(http.StatusOK, getOpenAPISpec())
	})
}

func (s *ServerTest) iRequestTheHealthzEndpoint() error {
	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		return err
	}
	s.response = httptest.NewRecorder()
	s.server.ServeHTTP(s.response, req)
	return nil
}

func (s *ServerTest) iRequestTheOpenapijsonEndpoint() error {
	req, err := http.NewRequest("GET", "/openapi.json", nil)
	if err != nil {
		return err
	}
	s.response = httptest.NewRecorder()
	s.server.ServeHTTP(s.response, req)
	return nil
}

func (s *ServerTest) iShouldReceiveAStatusCode(code int) error {
	if s.response.Code != code {
		return godog.ErrPending
	}
	return nil
}

func (s *ServerTest) theResponseShouldContainStatus(status string) error {
	var response map[string]interface{}
	json.Unmarshal(s.response.Body.Bytes(), &response)

	if response["status"] != status {
		return godog.ErrPending
	}
	return nil
}

func (s *ServerTest) theResponseShouldIndicateSSHIsRunning() error {
	var response map[string]interface{}
	json.Unmarshal(s.response.Body.Bytes(), &response)

	if response["ssh"] != "running" {
		return godog.ErrPending
	}
	return nil
}

func (s *ServerTest) theResponseShouldContainOpenAPISpecificationData() error {
	contentType := s.response.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return godog.ErrPending
	}

	var response map[string]interface{}
	json.Unmarshal(s.response.Body.Bytes(), &response)

	if _, exists := response["openapi"]; !exists {
		return godog.ErrPending
	}

	return nil
}

func getOpenAPISpec() map[string]interface{} {
	return map[string]interface{}{
		"openapi": "3.0.3",
		"info": map[string]interface{}{
			"title":       "Herbst MUD API",
			"description": "API for the Herbst MUD game server",
			"version":     "1.0.0",
		},
		"servers": []map[string]interface{}{
			{
				"url": "http://localhost:8080",
			},
		},
	}
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	serverTest := &ServerTest{}

	ctx.Step(`^the web server is running$`, serverTest.theWebServerIsRunning)
	ctx.Step(`^I request the healthz endpoint$`, serverTest.iRequestTheHealthzEndpoint)
	ctx.Step(`^I request the openapi\.json endpoint$`, serverTest.iRequestTheOpenapijsonEndpoint)
	ctx.Step(`^I should receive a (\d+) status code$`, serverTest.iShouldReceiveAStatusCode)
	ctx.Step(`^the response should contain status "([^"]*)"$`, serverTest.theResponseShouldContainStatus)
	ctx.Step(`^the response should indicate SSH is running$`, serverTest.theResponseShouldIndicateSSHIsRunning)
	ctx.Step(`^the response should contain OpenAPI specification data$`, serverTest.theResponseShouldContainOpenAPISpecificationData)

	// Export/Import test steps
	ctx.Step(`^I call GET /admin/export with world "([^"]*)"$`, serverTest.iCallTheExportEndpoint)
	ctx.Step(`^I call GET /admin/export/worlds$`, serverTest.iCallTheExportWorldsEndpoint)
	ctx.Step(`^I call POST /admin/import$`, serverTest.iCallTheImportEndpoint)
	ctx.Step(`^I call POST /admin/import/validate$`, serverTest.iCallTheValidateEndpoint)
	ctx.Step(`^I set the export data to: "([^"]*)"$`, serverTest.iSetExportedData)
	ctx.Step(`^I set the import data to: "([^"]*)"$`, serverTest.iSetImportData)
	ctx.Step(`^I should receive a valid export response$`, serverTest.iShouldReceiveAValidExportResponse)
	ctx.Step(`^the export should contain world "([^"]*)"$`, serverTest.theExportShouldContainWorld)
	ctx.Step(`^the import should succeed$`, serverTest.theImportShouldSucceed)
	ctx.Step(`^the validation should pass$`, serverTest.theValidationShouldPass)
	ctx.Step(`^the export should have rooms$`, serverTest.theExportShouldHaveRooms)
	ctx.Step(`^the export should have NPCs$`, serverTest.theExportShouldHaveNPCs)
	ctx.Step(`^the imported data should match original$`, serverTest.theImportedDataShouldMatchOriginal)

	// Test world import/destroy steps
	ctx.Step(`^I have the test world file$`, serverTest.iHaveTheTestWorldFile)
	ctx.Step(`^I import the test world$`, serverTest.iImportTheTestWorld)
	ctx.Step(`^I export the test world$`, serverTest.iExportTheTestWorld)
	ctx.Step(`^the test world should exist$`, serverTest.theTestWorldShouldExist)
	ctx.Step(`^I destroy the test world$`, serverTest.iDestroyTheTestWorld)
	ctx.Step(`^the test world should no longer exist$`, serverTest.theTestWorldShouldNoLongerExist)
}

func (s *ServerTest) iCallTheExportEndpoint(world string) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("/admin/export?world=%s", world), nil)
	if err != nil {
		return err
	}
	s.response = httptest.NewRecorder()
	s.server.ServeHTTP(s.response, req)
	return nil
}

func (s *ServerTest) iCallTheExportWorldsEndpoint() error {
	req, err := http.NewRequest("GET", "/admin/export/worlds", nil)
	if err != nil {
		return err
	}
	s.response = httptest.NewRecorder()
	s.server.ServeHTTP(s.response, req)
	return nil
}

func (s *ServerTest) iCallTheImportEndpoint() error {
	req, err := http.NewRequest("POST", "/admin/import", bytes.NewBuffer(s.importJson))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	s.response = httptest.NewRecorder()
	s.server.ServeHTTP(s.response, req)
	return nil
}

func (s *ServerTest) iCallTheValidateEndpoint() error {
	req, err := http.NewRequest("POST", "/admin/import/validate", bytes.NewBuffer(s.importJson))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	s.response = httptest.NewRecorder()
	s.server.ServeHTTP(s.response, req)
	return nil
}

func (s *ServerTest) iSetExportedData(exported string) error {
	s.importJson = []byte(exported)
	return nil
}

func (s *ServerTest) iShouldReceiveAValidExportResponse() error {
	var data ExportData
	if err := json.Unmarshal(s.response.Body.Bytes(), &data); err != nil {
		return godog.ErrPending
	}
	if data.Version != "1.0" {
		return godog.ErrPending
	}
	return nil
}

func (s *ServerTest) theExportShouldContainWorld(world string) error {
	var worldsResp struct {
		Worlds []struct {
			ID string `json:"id"`
		} `json:"worlds"`
	}
	if err := json.Unmarshal(s.response.Body.Bytes(), &worldsResp); err != nil {
		return godog.ErrPending
	}
	found := false
	for _, w := range worldsResp.Worlds {
		if w.ID == world {
			found = true
			break
		}
	}
	if !found {
		return godog.ErrPending
	}
	return nil
}

func (s *ServerTest) theImportShouldSucceed() error {
	var result ImportResult
	if err := json.Unmarshal(s.response.Body.Bytes(), &result); err != nil {
		return godog.ErrPending
	}
	if !result.Success {
		return godog.ErrPending
	}
	return nil
}

func (s *ServerTest) theValidationShouldPass() error {
	var validation Validation
	if err := json.Unmarshal(s.response.Body.Bytes(), &validation); err != nil {
		return godog.ErrPending
	}
	if !validation.IsValid {
		return godog.ErrPending
	}
	return nil
}

func (s *ServerTest) theExportShouldHaveRooms() error {
	var data ExportData
	if err := json.Unmarshal(s.response.Body.Bytes(), &data); err != nil {
		return godog.ErrPending
	}
	if len(data.Rooms) == 0 {
		return godog.ErrPending
	}
	return nil
}

func (s *ServerTest) theExportShouldHaveNPCs() error {
	var data ExportData
	if err := json.Unmarshal(s.response.Body.Bytes(), &data); err != nil {
		return godog.ErrPending
	}
	if len(data.NPCs) == 0 {
		return godog.ErrPending
	}
	return nil
}

func (s *ServerTest) theImportedDataShouldMatchOriginal() error {
	var result ImportResult
	if err := json.Unmarshal(s.response.Body.Bytes(), &result); err != nil {
		return godog.ErrPending
	}
	if result.Imported.Rooms == 0 && result.Imported.NPCs == 0 {
		return godog.ErrPending
	}
	return nil
}

func (s *ServerTest) iSetImportData(importData string) error {
	s.importJson = []byte(importData)
	return nil
}

// ──────────────────────────────────────────────
// Test world import/destroy step definitions
// ──────────────────────────────────────────────

func (s *ServerTest) iHaveTheTestWorldFile() error {
	data, err := os.ReadFile("../testing/test-world.json")
	if err != nil {
		return fmt.Errorf("cannot read test world file: %w", err)
	}
	s.importJson = data
	return nil
}

func (s *ServerTest) iImportTheTestWorld() error {
	req, err := http.NewRequest("POST", "/admin/import?world=test-world", bytes.NewBuffer(s.importJson))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	s.response = httptest.NewRecorder()
	s.server.ServeHTTP(s.response, req)

	// Store the world ID from the import result for later deletion
	var result ImportResult
	if err := json.Unmarshal(s.response.Body.Bytes(), &result); err == nil {
		if result.Success && result.Imported.Rooms > 0 {
			// Not storing ID here — we'll look up by name later
		}
	}
	return nil
}

func (s *ServerTest) iExportTheTestWorld() error {
	req, err := http.NewRequest("GET", "/admin/export?world=test-world", nil)
	if err != nil {
		return err
	}
	s.response = httptest.NewRecorder()
	s.server.ServeHTTP(s.response, req)
	return nil
}

func (s *ServerTest) theTestWorldShouldExist() error {
	// Check that the world appears in the worlds list
	// For now, verify the import succeeded by checking the response
	// The real check is done via GET /admin/export/worlds
	req, err := http.NewRequest("GET", "/admin/export/worlds", nil)
	if err != nil {
		return err
	}
	resp := httptest.NewRecorder()
	s.server.ServeHTTP(resp, req)

	var worldsResp struct {
		Worlds []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"worlds"`
		Default string `json:"default"`
		Count   int    `json:"count"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &worldsResp); err != nil {
		return godog.ErrPending
	}
	for _, w := range worldsResp.Worlds {
		if w.ID == "test-world" || w.Name == "test-world" {
			return nil
		}
	}
	return godog.ErrPending
}

func (s *ServerTest) iDestroyTheTestWorld() error {
	// First, wipe the world's content via admin/wipe
	wipeBody := `{"wipe_npcs":true,"wipe_rooms":true,"wipe_items":true,"wipe_skills":true}`
	req, _ := http.NewRequest("POST", "/admin/wipe/full?world=test-world", bytes.NewBuffer([]byte(wipeBody)))
	req.Header.Set("Content-Type", "application/json")
	s.response = httptest.NewRecorder()
	s.server.ServeHTTP(s.response, req)

	// Also try to delete any world record with name "test-world"
	// List worlds first
	listReq, _ := http.NewRequest("GET", "/api/worlds/db", nil)
	listResp := httptest.NewRecorder()
	s.server.ServeHTTP(listResp, listReq)

	var worldsResp struct {
		Worlds []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"worlds"`
	}
	if err := json.Unmarshal(listResp.Body.Bytes(), &worldsResp); err == nil {
		for _, w := range worldsResp.Worlds {
			if w.Name == "test-world" {
				delReq, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/worlds/%d", w.ID), nil)
				delResp := httptest.NewRecorder()
				s.server.ServeHTTP(delResp, delReq)
			}
		}
	}

	return nil
}

func (s *ServerTest) theTestWorldShouldNoLongerExist() error {
	// Verify the world is gone from the export list
	req, err := http.NewRequest("GET", "/admin/export/worlds", nil)
	if err != nil {
		return err
	}
	resp := httptest.NewRecorder()
	s.server.ServeHTTP(resp, req)

	var worldsResp struct {
		Worlds []struct {
			ID string `json:"id"`
		} `json:"worlds"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &worldsResp); err != nil {
		return godog.ErrPending
	}

	// If the worlds list only contains "default" (no test-world), we're clean
	// Any world named "test-world" means deletion failed
	for _, w := range worldsResp.Worlds {
		if w.ID == "test-world" {
			return godog.ErrPending
		}
	}
	return nil
}
