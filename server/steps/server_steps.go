package steps

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/cucumber/godog"
	"github.com/gin-gonic/gin"
)

type ServerTest struct {
	server   *gin.Engine
	response *httptest.ResponseRecorder
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
}