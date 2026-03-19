Feature: Character Authentication
  As a player
  I want to authenticate my character
  So that I can securely access the game

  Scenario: Authenticate with valid credentials
    Given a character "Hero" exists with password "secret123"
    When I authenticate with username "Hero" and password "secret123"
    Then I should receive an authentication token
    And I should be logged in

  Scenario: Authenticate with invalid password
    Given a character "Hero" exists with password "secret123"
    When I authenticate with username "Hero" and password "wrongpass"
    Then I should receive an authentication error
    And I should not be logged in

  Scenario: Authenticate with non-existent character
    Given no character "Unknown" exists
    When I authenticate with username "Unknown" and password "any"
    Then I should receive a "character not found" error