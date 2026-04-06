Feature: Character Entity - Data Structure
  As a game developer
  I want a properly defined Character entity
  So that player characters are stored and validated correctly

  Background:
    Given the database schema is properly initialized
    And the Character model is defined in the codebase

  Scenario: Character entity has all required database fields
    When I examine the Character entity/model definition
    Then it should have the following fields with correct types:
      | field       | type        | constraints              |
      | id          | UUID        | primary key             |
      | userId      | UUID        | foreign key to User     |
      | name        | string      | unique per user, 2-20   |
      | class       | string      | valid class enum        |
      | race        | string      | valid race enum         |
      | gender      | string      | valid gender enum       |
      | level       | integer     | default 1, min 1        |
      | experience  | integer     | default 0               |
      | health      | integer     | min 0, max is stats     |
      | mana        | integer     | min 0, max is stats     |
      | stats       | JSON object | strength/agi/int/wis/con |
      | inventory   | JSON array  | array of item objects   |
      | locationId  | UUID        | foreign key to Room     |
      | isOnline    | boolean     | default false           |
      | createdAt   | timestamp   | auto-set                |
      | updatedAt   | timestamp   | auto-updated            |

  Scenario: Character id is a valid UUID
    When I create a new character
    Then the generated id should be a valid UUID v4
    And the id should be unique across all characters

  Scenario: Character name length validation - minimum
    When I attempt to create a character with name "A"
    Then the validation should fail
    And the error should indicate name must be at least 2 characters

  Scenario: Character name length validation - maximum
    When I attempt to create a character with name "ThisNameIsWayTooLongForTheLimit"
    Then the validation should fail
    And the error should indicate name must be at most 20 characters

  Scenario: Character name uniqueness is per-user not global
    Given user "UserA" has a character named "Hero"
    When user "UserB" creates a character named "Hero"
    Then the character should be created successfully
    And "UserA" and "UserB" can both have characters named "Hero"

  Scenario: Character must reference a valid userId
    When I attempt to create a character with a non-existent userId
    Then the validation should fail
    And the error should indicate the user does not exist

  Scenario: Character class must be a valid class
    When I attempt to create a character with class "InvalidClass"
    Then the validation should fail
    And the error should indicate class must be one of: Warrior, Mage, Rogue, Priest

  Scenario: Character race must be a valid race
    When I attempt to create a character with race "InvalidRace"
    Then the validation should fail
    And the error should indicate race must be one of: Human, Elf, Dwarf, Orc

  Scenario: Character gender must be a valid gender
    When I attempt to create a character with gender "Unknown"
    Then the validation should fail
    And the error should indicate gender must be one of: Male, Female, Non-binary, Other

  Scenario: Character level defaults to 1
    When I create a character without specifying level
    Then the level should be set to 1
    And the experience should be set to 0

  Scenario: Character health cannot exceed maximum
    Given a character has maxHealth of 100
    When the character's health is set to 150
    Then the validation should cap health at 100
    And the character's health should be 100

  Scenario: Character health cannot be negative
    When I attempt to set a character's health to -10
    Then the validation should fail
    And health should remain at the minimum of 0

  Scenario: Character mana cannot exceed maximum
    Given a character has maxMana of 50
    When the character's mana is set to 100
    Then the validation should cap mana at 50

  Scenario: Character mana cannot be negative
    When I attempt to set a character's mana to -5
    Then the validation should fail
    And mana should remain at the minimum of 0

  Scenario: Character experience must be zero or positive
    When I attempt to set experience to -100
    Then the validation should fail

  Scenario: Character stats JSON has required fields
    When I examine a character's stats JSON
    Then it should contain: strength, agility, intelligence, wisdom, constitution
    And all values should be integers
    And each stat should be within valid range (1-100)

  Scenario: Character inventory is a valid JSON array
    When I examine a character's inventory
    Then it should be a valid JSON array
    And each item should have: id, name, quantity

  Scenario: Character locationId references a valid room
    When I create or update a character's location
    Then the locationId should reference an existing room in the database
    And the character should appear in that room when listed

  Scenario: Character isOnline flag reflects login state
    Given a character is logged into the game
    Then the isOnline flag should be true
    And the character should appear in the "online characters" list

    Given a character has logged out
    Then the isOnline flag should be false

  Scenario: createdAt is set automatically on creation
    When I create a new character
    Then the createdAt timestamp should be set automatically

  Scenario: updatedAt is updated on any character modification
    Given a character exists
    When I update any character field (name, level, health, location, etc.)
    Then the updatedAt timestamp should be updated

  Scenario: Character deleted when associated user is deleted
    Given a user has characters "Char1" and "Char2"
    When the user account is deleted
    Then both characters should be deleted
    Or both characters should be anonymized (disassociated from user)
