🔴 Feature: Character CRUD Operations
  As an administrator
  I want to manage characters through API endpoints
  So that I can create, read, update, and delete characters

  Background:
    Given the API server is running
    And the database is connected
    And a user "admin@example.com" exists with is_admin: true

  Scenario: Create a new character via POST endpoint
    Given a user "player@example.com" exists
    When I send a POST request to "/api/characters" with:
      | name | "TestHero" |
      | isNPC | false |
      | userId | "{user_id}" |
      | currentRoomId | 1 |
      | startingRoomId | 1 |
    Then the response status should be 201
    And the character "TestHero" should exist in the database

  Scenario: Get all characters via GET endpoint
    Given characters exist:
      | name | "Hero1" |
      | name | "Hero2" |
      | name | "NPCGuard" |
    When I send a GET request to "/api/characters"
    Then the response status should be 200
    And the response should contain 3 characters
    And the response should include character "Hero1"
    And the response should include character "Hero2"

  Scenario: Get a single character by ID
    Given a character "UniqueHero" exists with isNPC: false
    When I send a GET request to "/api/characters/{character_id}"
    Then the response status should be 200
    And the response should contain:
      | name | "UniqueHero" |
      | isNPC | false |

  Scenario: Update a character via PUT endpoint
    Given a character "OldName" exists
    When I send a PUT request to "/api/characters/{character_id}" with:
      | name | "NewName" |
    Then the response status should be 200
    And the character should be updated to "NewName"

  Scenario: Delete a character via DELETE endpoint
    Given a character "ToDelete" exists
    When I send a DELETE request to "/api/characters/{character_id}"
    Then the response status should be 204
    And the character "ToDelete" should no longer exist

  Scenario: Character belongs to user (many-to-one relationship)
    Given a user "owner@example.com" exists
    And a character "OwnedCharacter" belongs to that user
    When I get the character details
    Then the character should have userId matching the owner

  Scenario: Character has currentRoomId field
    Given a character "RoomHero" exists
    When I get the character details
    Then the character should have a currentRoomId field

  Scenario: Seed file contains Gandalf NPC in the "hole" room
    Given the database is seeded
    When I query for character "Gandalf"
    Then the character should exist
    And the character should have isNPC: true
    And the character should be in the "hole" room
    And the character should have _is_admin: true

  Scenario: Create character with invalid userId
    Given no user exists with id "99999"
    When I send a POST request to "/api/characters" with:
      | name | "OrphanHero" |
      | userId | 99999 |
    Then the response status should be 400 or 404
    And an error message about invalid user should be returned

  Scenario: Create character with duplicate name
    Given a character "DuplicateName" exists
    When I send a POST request to "/api/characters" with:
      | name | "DuplicateName" |
    Then the response status should be 409
    And an error message about duplicate name should be returned

  Scenario: Get non-existent character
    Given no character exists with id "99999"
    When I send a GET request to "/api/characters/99999"
    Then the response status should be 404

  Scenario: Update non-existent character
    Given no character exists with id "99999"
    When I send a PUT request to "/api/characters/99999" with:
      | name | "Ghost" |
    Then the response status should be 404

  Scenario: Delete non-existent character
    Given no character exists with id "99999"
    When I send a DELETE request to "/api/characters/99999"
    Then the response status should be 404

  Scenario: Character pagination
    Given 25 characters exist in the database
    When I send a GET request to "/api/characters?page=1&limit=10"
    Then the response status should be 200
    And the response should contain 10 characters
    And the response should include pagination metadata

  Scenario: Character filtering by isNPC
    Given characters exist:
      | name | "HeroPlayer" | isNPC | false |
      | name | "GuardNPC" | isNPC | true |
    When I send a GET request to "/api/characters?isNPC=true"
    Then the response status should be 200
    And all returned characters should have isNPC: true