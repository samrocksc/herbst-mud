Feature: Data Structure V1 - Core Entities
  As a game developer
  I want to understand the core data structures
  So that I can build on top of them correctly

  Background:
    Given the database schema is properly initialized
    And the ent ORM is configured

  Scenario: Character entity has all required fields
    When I examine the Character entity/model definition
    Then it should have the following fields:
      | field        | type       |
      | id           | UUID       |
      | name         | string     |
      | class        | string     |
      | race         | string     |
      | gender       | string     |
      | level        | integer    |
      | health       | integer    |
      | max_health   | integer    |
      | mana         | integer    |
      | max_mana     | integer    |
      | strength     | integer    |
      | dexterity    | integer    |
      | constitution | integer    |
      | intelligence | integer    |
      | wisdom       | integer    |

  Scenario: Character has skill proficiencies
    Given a character with class "Warrior"
    When I examine the character's skills
    Then it should have: blades, staves, unarmed, heavy_armor

  Scenario: Character has stats
    Given a character with class "Mage"
    When I examine the character's stats
    Then it should have: strength, dexterity, constitution, intelligence, wisdom

  Scenario: Character inventory is a JSON array
    When I examine the Character entity
    Then the inventory field should be a JSON array
    And each item should have: id, name, quantity, equipped

  Scenario: Character location is a UUID
    When I examine the Character entity
    Then location_id should be a UUID referencing a room

  Scenario: Character can be an NPC
    Given a character with isNPC set to true
    When I query the character
    Then the isNPC flag should be true
    And the character does not have a user_id

  Scenario: Character can be an admin
    Given a character with is_admin set to true
    When I query the character
    Then the is_admin flag should be true

  Scenario: Character name is unique
    When I create two characters with the same name
    Then the second creation should fail with 409 Conflict
