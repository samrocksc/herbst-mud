Feature: Data Structure V1 (Issue #12)
  As a game architect
  I want clear data structures for all entities
  So that the game has a consistent foundation

  Background:
    Given the database schema is defined
    And all entities have proper relationships

  Scenario: User data structure
    Given a User entity
    Then it should have is_admin boolean field
    And it can have up to three characters
    And it can have god_mode boolean for unkillable status

  Scenario: Character data structure
    Given a Character entity
    Then it should have a gender field (he/him, she/her, it/its, they/them)
    And it should have a class field (Warrior, Magician, Thief, Charlatan, Vigilante)
    And it should have a race field (Mutant, Human, Animal)
    And it should have is_npc boolean
    And it should have equipment slots
    And it should have stats: Strength, Intellect, Wisdom, Dexterity, Fortitude

  Scenario: Room data structure
    Given a Room entity
    Then it can hold characters
    And it can hold equipment
    And it has a description
    And it has exits (directions to other rooms)
    And it has atmosphere (air/water/wind)

  Scenario: Class types
    Given a Class entity
    Then valid types are: Warrior, Magician, Thief, Charlatan, Vigilante
    And each class has unique starting stats
    And each class has special abilities

  Scenario: User-Character relationship
    Given a user with multiple characters
    When I query the user
    Then I should see all associated characters
    And the relationship is one-to-many (user -> characters)

  Scenario: Character equipment
    Given a character with equipment
    When I query the character
    Then I should see equipped items
    And items should have slot assignments