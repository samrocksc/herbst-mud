Feature: Room Navigation (Issue #13)
  As a player
  I want to navigate between rooms
  So that I can explore the game world

  Background:
    Given I am logged into the game
    And I am in a room with exits

  Scenario: Move through exit
    Given I am in room "entrance" with north exit
    When I type "north"
    Then I should move to the connected room
    And I should see the new room description

  Scenario: See available exits
    Given I am in room "courtyard"
    When I type "look"
    Then I should see available exits listed
    And exits should show direction names (N, S, E, W, U, D)

  Scenario: Invalid direction
    Given I am in room "closet" with only south exit
    When I type "north"
    Then I should see "You cannot go that way"
    And I should remain in the current room

  Scenario: Move between connected rooms
    Given room A connects north to room B
    And I am in room A
    When I move north
    Then I should be in room B
    And room B should have a south exit back to room A

  Scenario: Multi-direction room
    Given I am in room "hub" with exits N, S, E, W
    When I look at the room
    Then I should see all four exits
    And I can move in any valid direction

  Scenario: Up and down navigation
    Given I am on ground floor with stairs up
    When I type "up"
    Then I should move to the upper floor
    And I can type "down" to return

  Scenario: Room description on entry
    Given I move to a new room
    When I enter the room
    Then I should see the room name
    And I should see the room description
    And I should see items and characters in the room
    And I should see visible exits