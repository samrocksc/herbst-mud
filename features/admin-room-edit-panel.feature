Feature: Admin - Room Edit Panel (admin-08)
  As an admin/builder
  I want to edit room details through a panel
  So that I can modify room properties efficiently

  Background:
    Given I am logged into the admin interface
    And I am in the map builder view

  @mikey
  Scenario: Open room edit panel
    Given I click on a room node
    When I double-click or right-click
    Then a room edit panel should appear
    And it should show the room's current properties

  @mikey
  Scenario: Edit room name
    Given the room edit panel is open
    When I change the room name
    Then the room should update immediately
    And the node label should reflect the new name

  @mikey
  Scenario: Edit room description
    Given the room edit panel is open
    When I modify the description
    Then the room description should be saved
    And it should appear correctly in-game

  @mikey
  Scenario: Edit room exits
    Given the room edit panel is open
    When I add/remove/edit exits
    Then the changes should persist
    And edges should update on the map

  @mikey
  Scenario: Add item to room
    Given the room edit panel is open
    When I add an item to the room
    Then the item should appear in-game in that room

  @mikey
  Scenario: Add NPC to room
    Given the room edit panel is open
    When I add an NPC to the room
    Then the NPC should spawn in that room in-game

  @mikey
  Scenario: Set room coordinates manually
    Given the room edit panel is open
    When I enter x, y, z coordinates
    Then the room should move to that position

  @mikey
  Scenario: Room edit panel validates input
    Given the room edit panel is open
    When I enter invalid data
    Then I should see validation errors
    And the save should be blocked until fixed

  @mikey
  Scenario: Delete room from edit panel
    Given the room edit panel is open
    When I click "Delete Room"
    Then the room should be removed
    And exits connected to it should also be removed