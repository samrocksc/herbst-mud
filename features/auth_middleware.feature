Feature: API Authentication Middleware
  As an API consumer
  I want to authenticate requests to protected routes
  So that only authorized users can access the API

  Background:
    Given the API server is running
    And I have a valid JWT token for a regular user
    And I have a valid JWT token for an admin user

  Scenario: Unauthenticated request to protected route returns 401
    When I send a GET request to "/api/rooms" without authentication
    Then the response status should be 401
    And the response should contain "Unauthorized"

  Scenario: Authenticated request passes through correctly
    When I send a GET request to "/api/rooms" with valid user token
    Then the response status should be 200
    And the response should contain room data

  Scenario: Admin-only route returns 403 for regular user
    When I send a GET request to "/api/admin/users" with regular user token
    Then the response status should be 403
    And the response should contain "Forbidden"

  Scenario: Admin user can access admin routes
    When I send a GET request to "/api/admin/users" with admin token
    Then the response status should be 200

  Scenario: Invalid token returns 401
    When I send a GET request to "/api/rooms" with token "invalid_token"
    Then the response status should be 401
    And the response should contain "Invalid token"