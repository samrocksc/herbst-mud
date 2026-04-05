🔴 Feature: Room Navigation
  As a player
  I want to move between rooms
  So that I can explore the MUD world

  Background:
    Given I am logged in as a player
    And I am in the "The Hole" starting room
    And the rooms exist:
      | name | exits |
      | The Hole | north:North Room, south:South Room, east:East Room, west:West Room |
      | North Room | south:The Hole |
      | South Room | north:The Hole |
      | East Room | west:The Hole |
      | West Room | east:The Hole |

  Scenario: Move north from starting room
    Given I am in "The Hole" starting room
    When I type "north"
    Then I should be in the "North Room"
    And I should see the room description
    And I should see the exits: "south"

  Scenario: Move in all cardinal directions
    Given I am in "The Hole" starting room
    When I type "north"
    Then I should be in "North Room"
    When I type "south"
    Then I should be back in "The Hole"
    When I type "east"
    Then I should be in "East Room"
    When I type "west"
    Then I should be back in "The Hole"

  Scenario: Cannot move in invalid direction
    Given I am in "North Room"
    When I type "north"
    Then I should see "You can't go that way."

  Scenario: Abbreviated movement commands
    Given I am in "The Hole" starting room
    When I type "n"
    Then I should be in "North Room"
    When I type "s"
    Then I should be back in "The Hole"
    When I type "e"
    Then I should be in "East Room"
    When I type "w"
    Then I should be back in "The Hole"

  Scenario: View current room
    Given I am in "The Hole" starting room
    When I type "look"
    Then I should see "The Hole"
    And I should see the room description
    And I should see the exits

  Scenario: Check exits
    Given I am in "The Hole" starting room
    When I type "exits"
    Then I should see "north, south, east, west"

  Scenario: Exits are color-coded by visited status
    Given I have not visited any rooms
    When I check exits in "The Hole"
    Then all exits should be displayed in white (new)
    When I move north to "North Room"
    And I return to "The Hole"
    And I check exits
    Then the "north" exit should be green (visited)
    And other exits should still be white

  Scenario: Peer into adjacent room
    Given I am in "The Hole" starting room
    When I type "peer north"
    Then I should see "North Room" name
    And I should see the room description
    And I should NOT move to that room
    And I should still be in "The Hole"

  Scenario: Movement is blocked during combat
    Given I am in "The Hole" starting room
    And I am in combat with an enemy
    When I type "north"
    Then I should see "You can't flee during combat"
    And I should remain in "The Hole"

  Scenario: Cannot move to non-existent room
    Given I am in "The Hole" starting room
    When I attempt to move to a room with invalid ID
    Then an error should be returned

  Scenario: Room description shows when entering
    Given I am in "The Hole" starting room
    When I move to "North Room"
    Then I should see the room description for "North Room"
    And I should see any NPCs in the room
    And I should see any items in the room

  Scenario: Movement is logged for history
    Given I am in "The Hole" starting room
    When I move to "North Room"
    And I return to "The Hole"
    Then my movement history should show both rooms

  Scenario: Auto-look on room entry
    Given I am in "The Hole" starting room
    When I move to "North Room"
    Then the room should automatically be looked at
    And I should see the full room description

  Scenario: Character position updates after movement
    Given my character "Hero" is in "The Hole"
    When I move to "North Room"
    Then my character's currentRoomId should be updated
    And querying my character should show the new room

  Scenario: Case-insensitive direction commands
    Given I am in "The Hole" starting room
    When I type "NORTH"
    Then I should be in "North Room"
    When I type "South"
    Then I should be back in "The Hole"

  Scenario: Full direction names work
    Given I am in "The Hole" starting room
    When I type "go north"
    Then I should be in "North Room"
    When I type "go south"
    Then I should be back in "The Hole"