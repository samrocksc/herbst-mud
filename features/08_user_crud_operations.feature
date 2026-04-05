🔴 Feature: User CRUD Operations
  As an administrator
  I want to manage users through API endpoints
  So that I can create, read, update, and delete users

  Background:
    Given the API server is running
    And the database is connected

  Scenario: Create a new user via POST endpoint
    When I send a POST request to "/api/users" with:
      | username | "newuser" |
      | password | "password123" |
      | is_admin | false |
    Then the response status should be 201
    And the user "newuser" should exist in the database

  Scenario: Get all users via GET endpoint
    Given users exist:
      | username | "user1" |
      | username | "user2" |
      | username | "admin" |
    When I send a GET request to "/api/users"
    Then the response status should be 200
    And the response should contain 3 users
    And the response should include user "user1"
    And the response should include user "user2"

  Scenario: Get a single user by ID
    Given a user "specificuser" exists
    When I send a GET request to "/api/users/{user_id}"
    Then the response status should be 200
    And the response should contain:
      | username | "specificuser" |

  Scenario: Update a user via PUT endpoint
    Given a user "oldname" exists
    When I send a PUT request to "/api/users/{user_id}" with:
      | username | "newname" |
    Then the response status should be 200
    And the user should be updated to "newname"

  Scenario: Delete a user via DELETE endpoint
    Given a user "todelete" exists
    When I send a DELETE request to "/api/users/{user_id}"
    Then the response status should be 204
    And the user "todelete" should no longer exist

  Scenario: User has is_admin boolean field
    Given a user "regularuser" exists with is_admin: false
    And a user "adminuser" exists with is_admin: true
    When I get the user "regularuser"
    Then the user should have is_admin: false
    When I get the user "adminuser"
    Then the user should have is_admin: true

  Scenario: User has one-to-many relationship with characters
    Given a user "owner" exists
    And the user has characters:
      | name | "Char1" |
      | name | "Char2" |
    When I get the user "owner"
    Then the user should have 2 characters

  Scenario: Seed file contains admin user
    Given the database is seeded
    When I query for admin users
    Then at least one user should have is_admin: true

  Scenario: Create user with duplicate username
    Given a user "existinguser" exists
    When I send a POST request to "/api/users" with:
      | username | "existinguser" |
      | password | "password123" |
    Then the response status should be 409
    And an error message about duplicate username should be returned

  Scenario: Create user with invalid email format
    When I send a POST request to "/api/users" with:
      | username | "newuser" |
      | email | "not-an-email" |
    Then the response status should be 400
    And an error message about invalid email should be returned

  Scenario: Create user with short password
    When I send a POST request to "/api/users" with:
      | username | "newuser" |
      | password | "123" |
    Then the response status should be 400
    And an error message about weak password should be returned

  Scenario: Get non-existent user
    Given no user exists with id "99999"
    When I send a GET request to "/api/users/99999"
    Then the response status should be 404

  Scenario: Update non-existent user
    Given no user exists with id "99999"
    When I send a PUT request to "/api/users/99999" with:
      | username | "Ghost" |
    Then the response status should be 404

  Scenario: Delete user with associated characters
    Given a user "owner" exists with characters:
      | name | "Char1" |
      | name | "Char2" |
    When I send a DELETE request to "/api/users/{user_id}"
    Then the response status should be 409 or the characters should be deleted first

  Scenario: User password is hashed
    Given a user "secureuser" exists with password "MyPassword123"
    When I query the database directly for the user
    Then the stored password should be hashed
    And the stored password should not be "MyPassword123"