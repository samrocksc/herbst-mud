🔴 Feature: Character Entity - Data Structure
  As a game developer
  I want a properly structured Character entity
  So that players can have fully-featured characters

  Background:
    Given the database is connected
    And the ent schema is generated

  Scenario: Character has gender field
    Given I create a new character
    Then the character should have a gender field
    And the gender can be "he/him"
    And the gender can be "she/her"
    And the gender can be "it/its"
    And the gender can be "they/them"
    And the gender can be custom-defined

  Scenario: Character has class field
    Given I create a new character with class "Warrior"
    Then the character should have class "Warrior"
    And the character can have one of:
      | class |
      | Warrior |
      | Magician |
      | Thief |
      | Charlatan |
      | Vigilante |

  Scenario: Character has race field
    Given I create a new character with race "Mutant"
    Then the character should have race "Mutant"
    And the character can have one of:
      | race |
      | Mutant |
      | Human |
      | Animal |
    And if race is "Animal", the character can have a species

  Scenario: Character can be NPC or player
    Given I create a character with isNPC = true
    Then the character should be an NPC
    And NPCs should not have a user associated
    When I create a character with isNPC = false
    Then the character should be a player character
    And the character should have a user associated

  Scenario: Character has equipment
    Given I create a new character
    Then the character should have an equipment slot system
    And equipment can be added to the character
    And equipment can be removed from the character

  Scenario: Character has stats
    Given I create a new character
    Then the character should have the following stats:
      | stat |
      | Strength |
      | Intellect |
      | Wisdom |
      | Dexterity |
      | Fortitude |
    And each stat should have a numeric value
    And stats should affect gameplay mechanics

  Scenario: Character stats scale with level
    Given a character at level 1 with Strength 10
    When the character levels up to level 2
    Then the character's Strength should increase
    And HP, Stamina, and Mana should recalculate based on stats

  Scenario: Character name is required
    When I attempt to create a character without a name
    Then the creation should fail
    And an error about required name should be returned

  Scenario: Character name has maximum length
    Given I attempt to create a character with name exceeding 20 characters
    Then the creation should fail
    And an error about name length should be returned

  Scenario: Character has description field
    Given I create a new character
    Then the character should have a description field
    And the description can be empty or contain text

  Scenario: Character has currentRoomId field
    Given I create a character in room "The Hole"
    Then the character should have currentRoomId pointing to "The Hole"
    And when the character moves, currentRoomId should update

  Scenario: Character has startingRoomId field
    Given I create a new character
    Then the character should have a startingRoomId
    And startingRoomId should match the room where the character was created

  Scenario: Character has level field
    Given I create a new character
    Then the character should have a level field
    And new characters should start at level 1

  Scenario: Character has experience field
    Given I create a new character
    Then the character should have an experience field
    And new characters should start with 0 experience

  Scenario: Character belongs to a user (foreign key)
    Given a user "owner@example.com" exists
    When I create a character for that user
    Then the character should have userId matching the owner
    And the character should be retrievable through the user's character list

  Scenario: Character gender affects pronouns in text
    Given a character has gender "he/him"
    When other players examine the character
    Then the description should use "he" and "him" pronouns

  Scenario: Animal race has species field
    Given I create a character with race "Animal"
    Then the character should have a species field
    And species could be "cat", "dog", "wolf", etc.

  Scenario: NPC character does not require userId
    Given I create an NPC character "Guard"
    Then the character should have isNPC: true
    And the character should not require a userId