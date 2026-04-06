Feature: Room Display (Look Command)
  As a player
  I want to see the current room's description
  So that I know my surroundings

  Background:
    Given the game is running

  Scenario: Look shows room name and description
    When player types "look"
    Then room name is displayed
    And room description is shown

  Scenario: Look shows available exits
    Given current room has exits north and south
    When player types "look"
    Then exits are displayed with directions
    And destinations are shown

  Scenario: Look shows items in room
    Given room contains items
    When player types "look"
    Then items in room are listed

  Scenario: Look shows NPCs in room
    Given room contains NPCs
    When player types "look"
    Then NPCs are listed

  Scenario: Look command has alias "l"
    Given the game is running
    When player types "l"
    Then same output as "look" command

  Scenario: Look shows exits with destination names
    Given room has exit north to "Piles of Rust"
    When player types "look"
    Then "[N]orth to Piles of Rust" is displayed