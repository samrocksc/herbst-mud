Feature: Character Creation
  As a new player
  I want to create a character
  So that I can start playing the game

  Scenario: Create a warrior character
    When I create a character with name "Warrior1", class "Warrior", race "Human"
    Then the character "Warrior1" should be created
    And the class should be "Warrior"
    And the race should be "Human"

  Scenario: Create character with invalid class
    When I create a character with name "Bad1", class "InvalidClass", race "Human"
    Then character creation should fail
    And I should receive an "invalid class" error

  Scenario: Create character with duplicate name
    Given a character "Duplicate" already exists
    When I create a character named "Duplicate"
    Then character creation should fail
    And I should receive a "name taken" error