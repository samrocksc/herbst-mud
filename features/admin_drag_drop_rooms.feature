Feature: Drag-and-Drop Room Creation and Positioning
  As an admin map builder
  I want to drag rooms to new positions
  So that I can arrange the map visually

  Scenario: Drag room to reposition
    Given a room exists on the map
    When I drag the room to a new position
    Then the room should move visually
    And the new position should persist

  Scenario: Snap to grid option
    Given snap-to-grid is enabled
    When I drag a room
    Then the room should snap to grid points

  Scenario: Snap to grid disabled
    Given snap-to-grid is disabled
    When I drag a room
    Then the room should move freely

  Scenario: Persist position after drag
    Given I dragged a room to a new position
    When I reload the map
    Then the room should be in the new position