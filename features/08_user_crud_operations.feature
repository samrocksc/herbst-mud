🔴 Feature: CRUD User Operations - Issue #8
  As a game system
  I need full CRUD operations for user management
  So that accounts can be created and managed

  Background:
    Given the game database is initialized

  # CREATE
  Scenario: Create a new user via API
    Given I want to create a new user account
    When I send a POST request to /api/users
    And the request body contains:
      | field    | value              |
      | username | "testuser"          |
      | password | "securepassword123" |
    Then the response status should be 201 Created
    And the user should be saved to the database
    And the password should be hashed (not stored in plaintext)

  Scenario: Create user with duplicate username
    Given a user exists with username "testuser"
    When I try to create another user with username "testuser"
    Then the response status should be 409 Conflict
    And the error message should indicate "username already exists"

  Scenario: Create user with missing username field
    Given I want to create a new user
    When I send a POST request to /api/users
    And the request body is missing username field
    Then the response status should be 400 Bad Request

  Scenario: Create user with short password
    Given I want to create a new user
    When I send a POST request with password "123"
    Then the response status should be 400 Bad Request
    And error should indicate "Password must be at least 6 characters"

  Scenario: Create user with empty JSON body
    Given I want to create a new user
    When I send a POST request to /api/users with empty body
    Then the response status should be 400 Bad Request

  # READ
  Scenario: Get all users
    Given multiple users exist in the database
    When I send a GET request to /api/users
    Then the response status should be 200 OK
    And the response should contain an array of users
    And passwords should NOT be included in the response

  Scenario: Get user by ID
    Given a user exists with ID 1
    When I send a GET request to /api/users/1
    Then the response status should be 200 OK
    And the response should contain user details
    And the password should NOT be included

  Scenario: Get non-existent user
    Given no user exists with ID 9999
    When I send a GET request to /api/users/9999
    Then the response status should be 404 Not Found

  Scenario: List users returns empty array when no users exist
    Given no users exist in the database
    When I send a GET request to /api/users
    Then the response status should be 200 OK
    And the response should be an empty array

  # UPDATE
  Scenario: Update user username
    Given a user exists with ID 1 and username "oldname"
    When I send a PUT request to /api/users/1
    And the request body contains username "newname"
    Then the response status should be 200 OK
    And the user username should be updated

  Scenario: Update user password
    Given a user exists with ID 1
    When I send a PUT request to /api/users/1
    And the request body contains:
      | field    | value              |
      | password | "newsecurepassword" |
    Then the response status should be 200 OK
    And the new password should be hashed
    And the old password hash should be replaced

  Scenario: Update non-existent user
    Given no user exists with ID 9999
    When I send a PUT request to /api/users/9999
    Then the response status should be 404 Not Found

  Scenario: Partial update only changes specified fields
    Given a user exists with ID 1 and username "original"
    When I send a PATCH request to /api/users/1
    And the request body contains username "updated"
    Then the response status should be 200 OK
    And other fields should remain unchanged

  # DELETE
  Scenario: Delete a user
    Given a user exists with ID 1
    When I send a DELETE request to /api/users/1
    Then the response status should be 204 No Content
    And the user should be removed from the database
    And the user's characters should be handled appropriately

  Scenario: Delete non-existent user
    Given no user exists with ID 9999
    When I send a DELETE request to /api/users/9999
    Then the response status should be 404 Not Found

  # ADMIN FLAG
  Scenario: User has admin flag
    Given a user exists with ID 1
    When I retrieve the user
    Then there should be an is_admin field
    And is_admin should default to false
    And setting is_admin to true should grant admin privileges

  Scenario: Admin user can access admin endpoints
    Given a user exists with is_admin true
    When the admin user requests admin endpoints
    Then the requests should be authorized
    And admin-only features should be accessible

  # RELATIONSHIPS
  Scenario: User has one-to-many relationship with characters
    Given a user exists with ID 1
    And the user has 3 characters
    When I query the user with characters
    Then the response should include all 3 characters
    And each character should reference the user

  Scenario: Deleting user handles characters
    Given a user exists with ID 1 and has 2 characters
    When I delete the user
    Then the characters should either be deleted or reassigned
    And no orphaned characters should remain

  # SEED DATA
  Scenario: Seed file creates admin user
    Given the database is seeded
    When I query for the admin user
    Then an admin user should exist
    And the admin should have is_admin true
    And the admin should have a predefined username and password

  # SECURITY
  Scenario: Password never appears in API response
    Given a user exists with ID 1
    When I send a GET request to /api/users/1
    Then the response JSON should not contain "password" field
    And no password hash should be visible in any user endpoint

  Scenario: Cannot set is_admin via regular user endpoint
    Given a regular user exists with ID 1
    When I send a PUT request to /api/users/1
    And the request body contains is_admin true
    Then the response should either reject the request or ignore the field
    And the user's admin status should not change
