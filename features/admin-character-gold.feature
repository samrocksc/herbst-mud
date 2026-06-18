Feature: Admin character gold credits management

  Background:
    Given I am authenticated as an admin
    And the active world is "Ooze Surfers"
    And a character "Gizmo" exists with "400" gold credits

  Scenario: Add gold credits to a character
    Given I navigate to the character detail page for "Gizmo"
    When I enter "100" in the gold credits amount field
    And I click the "+" button
    Then the character gold credits should display "500"

  Scenario: Spend gold credits from a character
    Given I navigate to the character detail page for "Gizmo"
    When I enter "100" in the gold credits amount field
    And I click the "−" button
    Then the character gold credits should display "400"

  Scenario: Cannot spend more gold than balance
    Given I navigate to the character detail page for "Gizmo"
    When I enter "1000" in the gold credits amount field
    Then the "−" button should be disabled
