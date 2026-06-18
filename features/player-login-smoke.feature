Feature: Player web client login and character selection

  Background:
    Given the web client is running
    And the backend is running
    And the user "sma" exists with password "sma"
    And user "sma" has a character in world "Ooze Surfers"

  Scenario: Login page renders without module import errors
    When I navigate to the player web client
    Then the login form should be visible
    And the browser console should not show import errors

  Scenario: Login and reach character selection
    Given I am on the player login page
    When I enter username "sma" and password "sma"
    And I click "ENTER THE WORLD"
    Then I should see the world selection screen
    When I select "Ooze Surfers"
    Then I should see the character selection screen

  Scenario: Select character and enter game
    Given I am on the character selection screen for "Ooze Surfers"
    When I select a character
    Then the game screen should load
    And the room title should be visible
