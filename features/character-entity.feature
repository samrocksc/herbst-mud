Feature: Character Entity Data Structure (Issue #17)
  As a game architect
  I want the Character entity properly implemented
  So that player characters are stored correctly

  Background:
    Given the database schema is defined
    And the Character entity follows DATA_STRUCTURE_V1

  Scenario: Character gender
    Given a Character entity
    Then it should have a gender field
    And valid options are: he/him, she/her, it/its, they/them
    And the system should be extensible for other pronouns

  Scenario: Character class
    Given a Character entity
    Then it should have a class field
    And valid classes are: Warrior, Magician, Thief, Charlatan, Vigilante
    And class affects starting stats and abilities

  Scenario: Character race
    Given a Character entity
    Then it should have a race field
    And valid races are: Mutant, Human, Animal
    And race provides unique bonuses

  Scenario: Character NPC flag
    Given a Character entity
    Then it should have is_npc boolean
    And NPCs are controlled by AI
    And player characters are controlled by users

  Scenario: Character equipment
    Given a Character entity
    Then it should have equipment slots
    And equipment affects stats
    And items can be equipped/unequipped

  Scenario: Character stats
    Given a Character entity
    Then it should have five primary stats:
    | Stat      | Description          |
    | Strength  | Physical power       |
    | Intellect | Mental capacity      |
    | Wisdom    | Insight and will    |
    | Dexterity | Agility and speed   |
    | Fortitude | Health and endurance|

  Scenario: Character belongs to user
    Given a Character entity
    Then it should have user_id foreign key
    And a character belongs to exactly one user