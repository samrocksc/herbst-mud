Feature: Admin - Room Edit Panel (admin-05)
  As an admin/builder
  I want to edit room properties in a side panel
  So that I can modify room details efficiently

  Background:
    Given I am logged into the admin interface
    And I am in the map builder view
    And a room exists on the canvas

  @mikey
  Scenario: Panel opens on room selection
    Given a room exists on the map
    When I click on the room
    Then a side panel should slide in
    And the panel should show room details

  @mikey
  Scenario: Edit room name and description
    Given the room edit panel is open
    When I modify the room name to "Town Square"
    And I modify the description
    And I click Save
    Then the room should be updated

  @mikey
  Scenario: Set room Z-level
    Given the room edit panel is open
    When I select Z-level from dropdown
    And I click Save
    Then the room's Z-level should change

  @mikey
  Scenario: View exits list
    Given the room edit panel is open
    Then I should see a list of exits
    And each exit should show direction and destination

  @mikey
  Scenario: Delete room
    Given the room edit panel is open
    When I click Delete
    And I confirm the deletion
    Then the room should be removed from the map

  @mikey
  Scenario: Close panel
    Given the room edit panel is open
    When I click X button or deselect the room
    Then the panel should close
