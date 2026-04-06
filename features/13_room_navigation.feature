Feature: Room Navigation
  As a player
  I want to navigate between rooms using directional commands
  So that I can explore the MUD world

  Background:
    Given the MUD game server is running
    And I am logged in as character "NavTest" in room "The Hole" (the starting room)
    And "The Hole" has an exit north to "North Room"

  # Basic Movement
  Scenario: Move north from starting room
    When I type the command "north"
    Then the response should indicate I am now in "North Room"
    And I should see the room description for "North Room"
    And my character locationId should be updated to "North Room"

  Scenario: Move south to return to previous room
    Given I am in "North Room"
    When I type the command "south"
    Then I should be returned to "The Hole" starting room
    And I should see the room description for "The Hole"

  Scenario: Move east in a room with east exit
    Given I am in a room with an east exit to "East Chamber"
    When I type the command "east"
    Then I should be in "East Chamber"

  Scenario: Move west in a room with west exit
    Given I am in a room with a west exit to "West Chamber"
    When I type the command "west"
    Then I should be in "West Chamber"

  # Vertical Movement
  Scenario: Move up to upper floor increases z-coordinate
    Given I am in "Ground Floor Room"
    And "Ground Floor Room" has an up exit to "Second Floor Room"
    When I type the command "up"
    Then I should be in "Second Floor Room"
    And the z-coordinate should increase by 1

  Scenario: Move down to lower floor decreases z-coordinate
    Given I am in "Second Floor Room"
    And "Second Floor Room" has a down exit to "Ground Floor Room"
    When I type the command "down"
    Then I should be in "Ground Floor Room"
    And the z-coordinate should decrease by 1

  # Invalid Direction
  Scenario: Attempting to move in a direction with no exit
    Given I am in a room with no east exit
    When I type the command "east"
    Then I should remain in the current room
    And I should see the message "You can't go that way."
    And my character locationId should not change

  Scenario: Attempting to move up when no up exit exists
    Given I am in a room with no up exit
    When I type the command "up"
    Then I should remain in the current room
    And I should see "You can't go that way."

  # Exit Listing
  Scenario: View available exits from current room
    Given I am in a room with exits north, south, and east
    When I type the command "exits"
    Then I should see a list of available directions: north, south, east
    And each direction should show the destination room name

  Scenario: Room description shows available exits
    When I look at the current room
    Then the room description should include a list of visible exits

  # Locked Doors
  Scenario: Movement blocked by locked door
    Given there is a locked door to the north
    And I do not have the key
    When I type the command "north"
    Then I should see the message "The door is locked."
    And I should remain in the current room
    And my character locationId should not change

  Scenario: Movement succeeds after unlocking door
    Given there is a locked door to the north
    And I have the key in my inventory
    And I use the key to unlock the door
    When I type the command "north"
    Then I should be in the next room
    And the door should remain unlocked for my session

  Scenario: Locked door with key but wrong key fails
    Given there is a locked door to the north requiring "Silver Key"
    And I have a "Bronze Key" in my inventory
    When I try to unlock the door
    Then the door should remain locked
    And I should see "You don't have the right key."

  # Multi-hop Navigation
  Scenario: Navigate through multiple rooms
    Given I am in "Room A"
    And "Room A" has an exit north to "Room B"
    And "Room B" has an exit east to "Room C"
    When I type "north" then "east"
    Then I should be in "Room C"

  Scenario: Navigating saves progress automatically
    Given I have successfully navigated from "Start Room" to "End Room"
    When I disconnect and reconnect
    Then I should be in "End Room"
    And my progress through intermediate rooms should not need to be repeated

  # Edge Cases
  Scenario: Cannot navigate while in combat
    Given I am in combat with an enemy
    When I type the command "north"
    Then I should see "You can't flee while in combat!"
    And I should remain in the current room
    And the enemy should still be engaged

  Scenario: Navigate to same room via different paths
    Given room "Junction" connects to both "North Wing" and "East Wing"
    When I navigate from "Junction" to "North Wing" then back through "Junction" to "East Wing"
    Then I should be able to navigate both directions successfully

  Scenario: Room with circular exit (self-referencing)
    Given a room has an exit that leads back to itself
    When I take that exit
    Then I should see the current room description
    And I should receive a message indicating a loop or unusual exit

  # Z-Axis Navigation Details
  Scenario: Stairs connect floors correctly
    Given I am in "Stairwell" on floor 0
    And the up exit leads to "Stairwell" on floor 1
    When I type "up"
    Then I should be in the "Stairwell" room on floor 1

  Scenario: Z-axis affects visibility of room connections
    Given I am on floor 2
    When I look at the room "Tower Top"
    Then I should see that exits to floor 1 are "down"
    And I should not see exits to other floors that are not directly adjacent
