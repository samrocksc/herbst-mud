Feature: Admin Players — character list in player detail modal

  Background:
    Given I am authenticated as an admin
    And the active world is "Ooze Surfers"

  Scenario: Open player detail modal without crashing
    Given I navigate to "/players"
    When I click on a player row
    Then the player detail modal opens
    And the character list renders without a TypeError

  Scenario: Guard against non-array character response
    Given I navigate to "/players"
    When I click on a player row whose character response is not an array
    Then the modal shows "No characters" instead of crashing
