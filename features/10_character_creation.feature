Feature: Character Creation
  As a new player
  I want to create a character with custom attributes
  So that I can start playing the game with a unique identity

  Background:
    Given the character creation API is available
    And I am logged in as a valid user with username "newplayer" and password "userPass123"

  Scenario: Create a new character with all required fields
    When I create a character with:
      | field    | value        |
      | name     | BraveHero    |
      | class    | Warrior      |
      | race     | Human        |
      | gender   | Male         |
    Then the response status should be 201 Created
    And the character should be created with name "BraveHero"
    And the character should have Warrior base stats
    And the character should start in the starting room "The Hole"
    And the character level should be 1
    And the character experience should be 0

  Scenario: Create a character with Elf race and Mage class
    When I create a character with name "MagicElf" class "Mage" race "Elf" gender "Male"
    Then the response status should be 201 Created
    And the character should have Elf racial bonuses applied
    And the character should have Mage base stats (high mana, low health)

  Scenario: Create character with duplicate name fails
    Given a character "DuplicateName" already exists
    When I try to create a character with name "DuplicateName" class "Warrior" race "Human" gender "Male"
    Then the response status should be 409 Conflict
    And the error message should contain "name already taken"

  Scenario: Create character with invalid class fails
    When I try to create a character with name "TestChar" class "InvalidClass" race "Human" gender "Male"
    Then the response status should be 400 Bad Request
    And the error message should contain "invalid class"

  Scenario: Create character with invalid race fails
    When I try to create a character with name "TestChar" class "Warrior" race "InvalidRace" gender "Male"
    Then the response status should be 400 Bad Request
    And the error message should contain "invalid race"

  Scenario: Create character with name too short
    When I try to create a character with name "A" class "Warrior" race "Human" gender "Male"
    Then the response status should be 400 Bad Request
    And the error message should contain "name must be at least"

  Scenario: Create character with name too long
    When I try to create a character with name "ThisNameIsWayTooLongForTheGame" class "Warrior" race "Human" gender "Male"
    Then the response status should be 400 Bad Request
    And the error message should contain "name must be at most"

  Scenario: Create character with special characters in name
    When I try to create a character with name "Hero@#$%" class "Warrior" race "Human" gender "Male"
    Then the response status should be 400 Bad Request
    And the error message should contain "name may only contain"

  Scenario: Create character without being logged in
    Given I am not logged in
    When I try to create a character with name "OrphanChar" class "Warrior" race "Human" gender "Male"
    Then the response status should be 401 Unauthorized

  Scenario: Each class has correct base stat distribution
    When I create a character "WarriorChar" class "Warrior" race "Human" gender "Male"
    Then the character stats should have high strength
    And the character stats should have high health
    And the character stats should have low mana

    When I create a character "MageChar" class "Mage" race "Human" gender "Male"
    Then the character stats should have high intelligence
    And the character stats should have high mana
    And the character stats should have low health

    When I create a character "RogueChar" class "Rogue" race "Human" gender "Male"
    Then the character stats should have high agility
    And the character stats should have medium health and mana

    When I create a character "PriestChar" class "Priest" race "Human" gender "Male"
    Then the character stats should have high wisdom
    And the character stats should have high mana
