Feature: Admin Room API Integration
  As an admin map builder
  I want rooms to sync with the backend API
  So that room data persists between sessions

  Scenario: Fetch all rooms on map load
    Given the API has existing rooms
    When I open the map builder
    Then all rooms should load from the API

  Scenario: Create new room via API
    Given I am in the map builder
    When I create a new room
    Then the room should be saved to the API
    And the room should have a server-assigned ID

  Scenario: Update room via API
    Given a room exists in the map
    When I modify room properties
    Then changes should sync to the API

  Scenario: Delete room with confirmation
    Given a room exists in the map
    When I delete the room
    Then I should see a confirmation dialog
    When I confirm deletion
    Then the room should be removed from the API