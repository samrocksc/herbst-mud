Feature: Character Creation
  As a new user
  I want to be able to create a character
  So that I can start playing the game

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