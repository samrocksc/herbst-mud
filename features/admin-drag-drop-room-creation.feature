Feature: Admin - Drag-and-Drop Room Creation (admin-05)
  As an admin/builder
  I want to create rooms by dragging on the map
  So that I can quickly build the game world

  Background:
    Given I am logged into the admin interface
    And I am in the map builder view

  @mikey
  Scenario: Create room by dragging
    When I drag from the "Add Room" tool onto the canvas
    Then a new room node should appear
    And the room should have default properties
    And I should be able to edit the room details

  @mikey
  Scenario: Drag creates room at drop position
    When I create a room by dragging to position (100, 200)
    Then the room should be created at those coordinates
    And it should appear in the correct position on the map

  @mikey
  Scenario: New room has default name
    Given I create a new room
    Then it should have a default name like "New Room"
    And I should be able to rename it immediately

  @mikey
  Scenario: Room created in correct z-layer
    Given I am on z-layer 2
    When I create a room
    Then the room should be created on z-layer 2

  @mikey
  Scenario: Multiple rooms can be created
    When I drag to create multiple rooms
    Then each should have unique IDs
    And all should appear on the map

  @mikey
  Scenario: Cancel room creation
    When I start dragging to create a room
    And I press Escape or cancel
    Then no room should be created
    And the canvas should remain unchanged