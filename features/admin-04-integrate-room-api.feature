Feature: Admin Integrate Room API
  As an admin user
  I want to manage rooms through the API
  So that I can create, update, and delete rooms in the map builder

  Background:
    Given I am logged in as an admin
    And the admin backoffice is loaded

  Scenario: Fetch all rooms on mount
    Given the Room API is available
    When the admin panel loads
    Then all rooms should be fetched from the API
    And each room should be displayed as a node

  Scenario: Create new room via API
    Given I am on the map builder
    When I create a new room with name "Tavern"
    And I set description "A dimly lit tavern"
    Then the room should be saved via the API
    And I should see the new room on the map

  Scenario: Update room properties
    Given a room "Old Room" exists
    When I update the room name to "New Room"
    And I update the description to "Updated description"
    Then the room should be updated via the API
    And I should see the updated values

  Scenario: Delete room with confirmation
    Given a room "Test Room" exists
    When I select "Delete" for the room
    Then I should see a confirmation dialog
    And when I confirm deletion
    Then the room should be deleted via the API
    And the room should be removed from the map

  Scenario: Cancel room deletion
    Given a room "Test Room" exists
    When I select "Delete" for the room
    And I cancel the confirmation
    Then the room should not be deleted
    And the room should remain on the map

  Scenario: Handle API error on fetch
    Given the Room API is unavailable
    When the admin panel loads
    Then an error message should be displayed
    And the map should show empty state or cached data

  Scenario: Handle API error on create
    Given the Room API returns an error
    When I try to create a new room
    Then an error message should be displayed
    And the room should not appear on the map

  Scenario: Handle API error on update
    Given the Room API returns an error
    And a room "Test Room" exists
    When I try to update the room
    Then an error message should be displayed
    And the original values should be preserved

  Scenario: Handle API error on delete
    Given the Room API returns an error
    And a room "Test Room" exists
    When I confirm deletion
    Then an error message should be displayed
    And the room should remain on the map

  Scenario: Room properties editable
    Given a room exists
    When I edit the room
    Then I should be able to modify:
      | Property      |
      | name          |
      | description   |
      | coordinates   |
      | exits         |

  Scenario: Multiple room operations
    Given multiple rooms exist
    When I perform operations on different rooms
    Then each operation should be independent
    And changes to one room should not affect others

  Scenario: Refresh rooms from API
    Given rooms have been modified externally
    When I trigger a refresh
    Then rooms should be re-fetched from the API
    And the display should reflect current data