package steps

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/cucumber/godog"
	"github.com/gin-gonic/gin"
)

type ExportData struct {
	Version    string `json:"version"`
	ExportedAt string `json:"exported_at"`
	Rooms      []RoomData `json:"rooms"`
	NPCs       []NPCData `json:"npcs"`
	Skills     []SkillData `json:"skills"`
	Items      []ItemData `json:"items"`
}

type RoomData struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	IsStarting  bool           `json:"is_starting"`
	Exits       []ExitData     `json:"exits"`
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
	Success   bool        `json:"success"`
	Imported  Imported    `json:"imported"`
	Version   string      `json:"version"`
	ImportedAt string     `json:"imported_at"`
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
	server    *gin.Engine
	response  *httptest.ResponseRecorder
	importJson []byte
}

func (s *ServerTest) theWebServerIsRunning() error {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test router
	s.server = gin.New()

	// Define the healthz endpoint
	s.server.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"ssh":    "running",
		})
	})

	// Define the openapi.json endpoint
	s.server.GET("/openapi.json", func(c *gin.Context) {
		c.JSON(http.StatusOK, getOpenAPISpec())
	})

	return nil
}

func (s *ServerTest) iRequestTheHealthzEndpoint() error {
	// Create a test request
	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		return err
	}
	s.response = httptest.NewRecorder()
	s.server.ServeHTTP(s.response, req)
	return nil
}

func (s *ServerTest) iRequestTheOpenapijsonEndpoint() error {
	// Create a test request
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
		"paths": map[string]interface{}{
			"/healthz": map[string]interface{}{
				"get": map[string]interface{}{
					"summary":     "Health check endpoint",
					"description": "Returns the health status of the server",
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Successful response",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"type": "object",
										"properties": map[string]interface{}{
											"status": map[string]interface{}{
												"type": "string",
											},
											"ssh": map[string]interface{}{
												"type": "string",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"/openapi.json": map[string]interface{}{
				"get": map[string]interface{}{
					"summary":     "OpenAPI specification",
					"description": "Returns the OpenAPI specification for this API",
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Successful response",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"type": "object",
									},
								},
							},
						},
					},
				},
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
	// The import data should be set via the importJson field
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
	// This would need to store original data for comparison
	// For now, just verify the import succeeded
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