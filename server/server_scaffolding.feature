Feature: Initial Backend Server Scaffolding
  As a developer
  I want to have a web server with health check endpoints
  So that I can verify the server is running properly

  Scenario: Health check endpoint returns success status
    Given the web server is running
    When I request the healthz endpoint
    Then I should receive a 200 status code
    And the response should contain status "ok"
    And the response should indicate SSH is running

  Scenario: OpenAPI specification endpoint is accessible
    Given the web server is running
    When I request the openapi.json endpoint
    Then I should receive a 200 status code
    And the response should contain OpenAPI specification data