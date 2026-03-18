Feature: CRUD User Operations (Issue #8)
  As a game developer
  I want full CRUD operations for users
  So that I can manage player accounts in the database

  Background:
    Given the database is connected
    And I have valid API credentials

  Scenario: Create a new user
    Given I have user data with email "player@example.com"
    When I send a POST request to /api/users
    Then a new user should be created
    And the response should include the user ID

  Scenario: Retrieve all users
    Given multiple users exist in the database
    When I send a GET request to /api/users
    Then I should receive a list of all users
    And each user should have email and is_admin fields

  Scenario: Retrieve a user by ID
    Given a user exists with ID "user-001"
    When I send a GET request to /api/users/user-001
    Then I should receive the user details
    And the response should include email and is_admin

  Scenario: Update a user
    Given a user exists with ID "user-001"
    When I send a PUT request to /api/users/user-001 with updated data
    Then the user should be updated
    And the response should reflect the changes

  Scenario: Delete a user
    Given a user exists with ID "user-001"
    When I send a DELETE request to /api/users/user-001
    Then the user should be removed
    And subsequent GET requests should return 404

  Scenario: User admin flag
    Given I create a user with admin privileges
    When I set is_admin to true
    Then the user should have admin access
    And the user can access admin-only endpoints

  Scenario: User password hashing
    Given I create a user with password "secret123"
    When the user is saved to the database
    Then the password should be hashed with bcrypt
    And the plain password should never be stored