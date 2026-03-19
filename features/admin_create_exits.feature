Feature: Create Exits by Dragging
  As an admin map builder
  I want to create exits by dragging between rooms
  So that I can easily connect rooms

  Background:
    Given I am in the admin map builder with multiple rooms

  Scenario: Create one-way exit by dragging
    Given "Room A" exists
    And "Room B" exists
    When I drag from the edge of "Room A" to "Room B"
    Then a one-way exit should be created from "Room A" to "Room B"
    And the direction should be determined by positions

  Scenario: Show preview line while dragging
    When I start dragging to create an exit
    Then I should see a preview line from the source room

  Scenario: Create two-way exit with Ctrl+drag
    Given "Room A" and "Room B" exist
    When I Ctrl+drag from "Room A" to "Room B"
    Then a two-way exit should be created

  Scenario: Cancel exit creation
    When I start dragging to create an exit
    And I press Escape
    Then the exit creation should be cancelled