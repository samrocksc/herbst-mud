Feature: Character Creation (Issue #10)
  As a new player
  I want to create a character when I first log in
  So that I can start playing the game

  Background:
    Given I am a new user who has logged into the game
    And I have no existing characters

  Scenario: See character creation screen
    Given I have logged in successfully
    When I have no characters
    Then I should see a character creation form
    And the form should allow me to enter a character name

  Scenario: Create character with valid name
    Given I see the character creation screen
    When I enter a valid character name
    And I select a race
    And I select a class
    And I select a gender
    And I submit the form
    Then a new character should be created
    And I should see my character in the game

  Scenario: Character name validation
    Given I see the character creation screen
    When I enter a character name that is too short
    And I submit the form
    Then I should see an error message about the name length
    And the character should not be created

  Scenario: Create multiple characters
    Given I already have one character
    When I choose to create another character
    And I have fewer than 3 characters
    Then I should be able to create a new character
    And both characters should be associated with my account

  Scenario: Maximum character limit
    Given I already have 3 characters
    When I try to create another character
    Then I should see an error about the character limit
    And I should not be able to create more characters

  Scenario: Select race during creation
    Given I am on the character creation screen
    When I select race "Mutant"
    Then the character should have the Mutant race
    And stats should be adjusted for Mutant race

  Scenario: Select class during creation
    Given I am on the character creation screen
    When I select class "Warrior"
    Then the character should have the Warrior class
    And starting equipment should match the class

  Scenario: Select gender during creation
    Given I am on the character creation screen
    When I select gender "they/them"
    Then the character should have the they/them gender
    And pronouns should be used correctly in game text