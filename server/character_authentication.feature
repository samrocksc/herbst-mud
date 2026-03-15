Feature: Character Authentication
  As a player
  I want to authenticate my character
  So that I can securely access the game

  Scenario: Successful character authentication
    Given a character "Warrior1" exists with password "securePass123"
    When I authenticate with character name "Warrior1" and password "securePass123"
    Then authentication should be successful
    And I should receive an access token

  Scenario: Failed authentication with wrong password
    Given a character "Mage1" exists with password "magicWord"
    When I authenticate with character name "Mage1" and password "wrongPassword"
    Then authentication should fail
    And I should receive an error message

  Scenario: Failed authentication for non-existent character
    Given no character "UnknownHero" exists
    When I authenticate with character name "UnknownHero" and password "anyPassword"
    Then authentication should fail
    And I should receive a "character not found" message