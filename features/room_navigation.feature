Feature: Room Navigation
  As a player
  I want to navigate between rooms using directional commands
  So that I can explore the MUD world

  Scenario: Move north from starting room
    Given I am in "The Hole" starting room
    When I type "north"
    Then I should be in the "North Room"
    And I should see the room description

  Scenario: Move south to return to previous room
    Given I am in "The North Room"
    When I type "south"
    Then I should be in "The Hole" starting room
    And I should see the room description

  Scenario: Attempt to move in invalid direction
    Given I am in a room with no exit to the east
    When I type "east"
    Then I should see "You can't go that way"
    And I should remain in the current room

  Scenario: Move up to upper floor
    Given I am in "Ground Floor Room"
    When I type "up"
    Then I should be in "Second Floor Room"
    And the z-coordinate should increase by 1

  Scenario: Move down to lower floor
    Given I am in "Second Floor Room"
    When I type "down"
    Then I should be in "Ground Floor Room"
    And the z-coordinate should decrease by 1

  Scenario: View available exits
    Given I am in a room with multiple exits
    When I type "exits" or "look"
    Then I should see a list of available directions
    And each exit should show the destination room name

  Scenario: Movement blocked by locked door
    Given there is a locked door to the north
    When I type "north"
    Then I should see "The door is locked"
    And I should remain in the current room

  Scenario: Successful navigation after unlock
    Given there is a locked door to the north
    And I have unlocked it with the key
    When I type "north"
    Then I should be in the next room