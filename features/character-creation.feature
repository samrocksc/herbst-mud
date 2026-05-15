Feature: Character Creation
  As a new user
  I want to be able to create a character
  So that I can start playing the game

  Background:
    Given the character creation API is available
    And I am logged in as a valid user

  Scenario: New user with no characters can create a character
    Given a user "newplayer@example.com" exists with password "password123"
    And the user has no characters
    When I create a character "NewHero" with password "heroPass123" for user
    Then character creation should be successful
    And the character should belong to the user

  Scenario: User can view their characters
    Given a user "player@example.com" exists
    And the user has a character "ExistingHero"
    When I request characters for the user
    Then I should see character "ExistingHero"

  Scenario: User needs to create character after login
    Given an authenticated user
    When I check if I need to create a character
    And I have no characters
    Then I should be prompted to create a character

  Scenario: Create character with duplicate name fails
    Given a character "DupName" already exists
    When I try to create a character named "DupName"
    Then I should receive a "name taken" error

  Scenario: Create character with invalid class fails
    Given I am logged in
    When I try to create a character with class "InvalidClass"
    Then I should receive a "invalid class" error

  Scenario: Create character with invalid race fails
    Given I am logged in
    When I try to create a character with race "InvalidRace"
    Then I should receive a "invalid race" error

  Scenario: Create character with name too short
    Given I am logged in
    When I try to create a character with name "A"
    Then I should receive a "name too short" error

  Scenario: Create character with name too long
    Given I am logged in
    When I try to create a character with name "ThisNameIsWayTooLongForTheGame"
    Then I should receive a "name too long" error

  Scenario: Each class has correct base stat distribution
    When I create a character "WarriorChar" class "Warrior"
    Then the character stats should have high strength
    And the character stats should have high health
    And the character stats should have low mana

    When I create a character "MageChar" class "Mage"
    Then the character stats should have high intelligence
    And the character stats should have high mana
    And the character stats should have low health

    When I create a character "RogueChar" class "Rogue"
    Then the character stats should have high agility
    And the character stats should have medium health and mana

    When I create a character "PriestChar" class "Priest"
    Then the character stats should have high wisdom
    And the character stats should have high mana
