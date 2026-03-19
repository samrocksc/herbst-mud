Feature: Character Login
  As a player
  I want to log into the game with my character
  So that I can start playing

  Scenario: Login with valid character
    Given my character "Hero" exists and is not logged in
    When I log in as "Hero"
    Then I should enter the game world
    And I should be placed in my last location or starting room

  Scenario: Login when already logged in
    Given my character "Hero" is already logged in
    When I try to log in as "Hero" again
    Then I should receive an "already logged in" message

  Scenario: Logout from game
    Given I am logged in as "Hero"
    When I type the logout command
    Then I should be disconnected from the game
    And my character position should be saved