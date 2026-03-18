Feature: CRUD Characters Operations (Issue #7)
  As a game developer
  I want full CRUD operations for characters
  So that I can manage player characters in the database

  Background:
    Given the database is connected
    And I have valid API credentials

  Scenario: Create a new character
    Given I have character data with name "Leo" and race "Mutant"
    When I send a POST request to /api/characters
    Then a new character should be created
    And the response should include the character ID
    And the character should have the correct name and race

  Scenario: Retrieve all characters
    Given multiple characters exist in the database
    When I send a GET request to /api/characters
    Then I should receive a list of all characters
    And each character should have name, isNPC, and currentRoomId

  Scenario: Retrieve a character by ID
    Given a character exists with ID "char-001"
    When I send a GET request to /api/characters/char-001
    Then I should receive the character details
    And the response should include name, isNPC, currentRoomId, and startingRoomId

  Scenario: Update a character
    Given a character exists with ID "char-001"
    When I send a PUT request to /api/characters/char-001 with updated data
    Then the character should be updated
    And the response should reflect the changes

  Scenario: Delete a character
    Given a character exists with ID "char-001"
    When I send a DELETE request to /api/characters/char-001
    Then the character should be removed
    And subsequent GET requests should return 404

  Scenario: Character linked to user
    Given a user exists with ID "user-001"
    When I create a character for that user
    Then the character should have a user_id foreign key
    And the user can have multiple characters

  Scenario: Create NPC character
    Given I want to create an NPC
    When I create a character with isNPC set to true
    Then the character should be marked as NPC
    And it should not require a user association

  Scenario: Gandalf NPC admin character
    Given I want to create the Gandalf NPC
    When I create a character with name "Gandalf"
    Then the character should have isNPC set to true
    And the character should have is_admin set to true
    And the character should start in the "hole" room