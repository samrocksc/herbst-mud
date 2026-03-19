Feature: Room Edit Panel
  As an admin map builder
  I want to edit room properties in a side panel
  So that I can modify room details

  Scenario: Panel opens on room selection
    Given a room exists on the map
    When I click on the room
    Then a side panel should slide in
    And the panel should show room details

  Scenario: Edit room name and description
    Given the room edit panel is open
    When I modify the room name to "Town Square"
    And I modify the description
    And I click Save
    Then the room should be updated

  Scenario: Set room Z-level
    Given the room edit panel is open
    When I select Z-level from dropdown
    And I click Save
    Then the room's Z-level should change

  Scenario: View exits list
    Given the room edit panel is open
    Then I should see a list of exits
    And each exit should show direction and destination

  Scenario: Delete room
    Given the room edit panel is open
    When I click Delete
    And I confirm the deletion
    Then the room should be removed from the map

  Scenario: Close panel
    Given the room edit panel is open
    When I click X button or deselect the room
    Then the panel should close