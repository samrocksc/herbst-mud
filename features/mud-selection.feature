Feature: MUD/World Selection
  As a player
  I want to select which MUD world to play in
  So that I can choose my adventure

  Background:
    Given I am logged in as a valid user
    And the worlds API returns available worlds
    And I am on the world selection screen

  Scenario: Player sees available worlds after login
    Given I have successfully logged in
    Then I should see a list of available worlds
    And each world should show its name
    And I should be prompted to select a world

  Scenario: Player selects a world by number
    Given the available worlds include "Herbst"
    When I type "1"
    Then the world "Herbst" should be selected
    And I should see the character selection screen

  Scenario: Player selects a world by name
    Given the available worlds include "Herbst"
    When I type "Herbst"
    Then the world "Herbst" should be selected
    And I should see the character selection screen

  Scenario: Player goes back from world selection
    Given I am on the world selection screen
    When I type "b"
    Then I should be back on the welcome screen

  Scenario: Player sees error for invalid world number
    Given there are 3 available worlds
    When I type "9"
    Then I should see an error message
    And I should remain on the world selection screen

  Scenario: Player navigates world list with j/k keys
    Given I am on the world selection screen
    And there are multiple worlds available
    When I press "j"
    Then the cursor should move down
    When I press "k"
    Then the cursor should move up

  Scenario: Player selects highlighted world with enter
    Given I am on the world selection screen
    And a world is highlighted by the cursor
    When I press "enter"
    Then the highlighted world should be selected
    And I should see the character selection screen

  Scenario: Player can quit from world selection
    Given I am on the world selection screen
    When I type "quit"
    Then I should see the welcome screen
