Feature: Admin - Create Exits by Dragging (admin-06)
  As an admin/builder
  I want to create room exits by dragging between nodes
  So that I can visually connect rooms in the map

  Background:
    Given I am logged into the admin interface
    And I am in the map builder view
    And I have at least two rooms on the canvas

  @mikey
  Scenario: Create exit by dragging from room edge
    Given Room A exists on the map
    When I drag from Room A's edge toward Room B
    Then an exit should be created from A to B
    And I should be able to set the direction

  @mikey
  Scenario: Drag creates exit in correct direction
    Given Room A is at position (0, 0)
    And Room B is at position (100, 0)
    When I create an exit by dragging east
    Then the exit direction should be "east"

  @mikey
  Scenario: Exit appears as visual edge
    Given I create an exit between rooms
    When I view the map
    Then a line/edge should connect the two rooms
    And it should visually indicate the direction

  @mikey
  Scenario: Two-way exit creation
    Given I create a bidirectional exit
    When the connection is made
    Then exits should exist in both directions
    And both edges should be visible

  @mikey
  Scenario: Locked exit creation
    Given I create an exit that is locked
    Then the edge should indicate locked status
    And I should be able to set the lock key/item

  @mikey
  Scenario: Delete exit by selecting and removing
    Given an exit exists between rooms
    When I select the exit edge
    And I press Delete or use context menu
    Then the exit should be removed

  @mikey
  Scenario: Exit shows required level for locked doors
    Given I create a locked exit
    Then the edge should show required level
    Or the exit properties should list the requirement