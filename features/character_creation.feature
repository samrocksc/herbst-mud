Feature: Character Creation
  As a new player
  I want to create a character
  So that I can start playing the game

  Scenario: Create a new character with all fields
    Given I am logged in as a valid user
    When I create a character with:
      | field    | value        |
      | name     | BraveHero    |
      | class    | Warrior      |
      | race     | Human        |
      | gender   | Male         |
    Then the character should be created
    And the character should have default stats
    And the character should start in the starting room

  Scenario: Create character with duplicate name
    Given a character "DuplicateName" already exists
    When I try to create a character named "DuplicateName"
    Then I should receive a "name taken" error

  Scenario: Create character with invalid class
    Given I am logged in
    When I try to create a character with class "InvalidClass"
    Then I should receive a "invalid class" error