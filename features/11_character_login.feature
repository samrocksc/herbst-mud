Feature: Character Login
  As a player
  I want to log into the game with my character
  So that I can enter and play the game world

  Background:
    Given the game login API is available
    And the game world is running
    And my character "HeroLogin" exists and is not logged in

  Scenario: Login with valid character enters the game world
    When I log in as character "HeroLogin"
    Then the response status should be 200 OK
    And I should enter the game world
    And I should be placed in my last saved location or starting room "The Hole"
    And the server should send the room description for my current location
    And the server should send my character stats and status

  Scenario: Login sends correct initial game state
    When I log in as character "HeroLogin"
    Then the server should send:
      - Current room description
      - Available exits
      - Items in the room
      - Character inventory
      - Character stats and status
      - Any active quests

  Scenario: Login when character is already logged in
    Given my character "HeroLogin" is already logged in from another session
    When I try to log in as "HeroLogin" from a new connection
    Then the response status should be 409 Conflict
    And the error message should contain "already logged in"
    And the new session should not gain control of the character

  Scenario: Login with non-existent character fails
    When I try to log in as character "NobodyHome"
    Then the response status should be 404 Not Found
    And the error message should contain "character not found"

  Scenario: Logout disconnects from the game
    Given I am logged in as "HeroLogin"
    When I send the logout command
    Then my connection should be closed gracefully
    And my character "HeroLogin" should be marked as logged out
    And my character position should be saved to the database

  Scenario: Character position is saved on logout
    Given I am logged in as "HeroLogin" in room "The Tavern"
    When I send the logout command
    Then the database should record that "HeroLogin" was last in the "Tavern" room

  Scenario: Re-login restores character to last position
    Given I was logged in as "HeroLogin" and last saved in room "North Corridor"
    When I log in as character "HeroLogin"
    Then I should be placed in the "North Corridor" room
    And I should not be reset to the starting room

  Scenario: Character health and mana are restored on login
    Given my character "HeroLogin" was saved with 50 health and 20 mana
    And my character is currently offline
    When I log in as character "HeroLogin"
    Then my character health should be 50
    And my character mana should be 20

  Scenario: Multiple characters of same user cannot both be logged in
    Given I have two characters "CharA" and "CharB"
    And "CharA" is currently logged in
    When I try to log in as "CharB"
    Then the response status should be 409 Conflict
    Or "CharA" should be automatically logged out
