Feature: Room Navigation
  As a player
  I want to move between rooms
  So that I can explore the MUD world

  Scenario: Move north from starting room
    Given I am in "The Hole" starting room
    When I type "north"
    Then I should be in the "North Room"
    And I should see the room description

  Scenario: Move south to return
    Given I am in the "North Room"
    When I type "south"
    Then I should be in "The Hole"

  Scenario: Move east
    Given I am in "The Hole"
    When I type "east"
    Then I should be in the "East Room"

  Scenario: Move west
    Given I am in "The Hole"
    When I type "west"
    Then I should be in the "West Room"

  Scenario: Attempt to move in invalid direction
    Given I am in "The Hole"
    When I type "up"
    Then I should see "You can't go that way"
    And I should remain in "The Hole"

  Scenario: View available exits
    Given I am in "The Hole"
    When I look around
    Then I should see exits: north, east, west