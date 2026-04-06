🔴 Feature: Room Display (Look Command) - Issue Look-01
  As a player
  I want to see room details when I use the look command
  So that I understand my surroundings

  Background:
    Given the player is in a room
    And the room has been entered or re-described

  Scenario: Look shows room name
    Given the player is in a room
    When the player types "look"
    Then the room name should be displayed prominently
    And the name should be at the top of the output

  Scenario: Look shows room description
    Given the player is in a room
    When the player types "look"
    Then the room description should be shown
    And the description should be multi-line text

  Scenario: Look shows exits
    Given the player is in a room
    When the player types "look"
    Then available exits should be listed
    And exits should show directions: north, south, east, west, up, down
    And locked or special exits should be indicated

  Scenario: Look shows other characters in the room
    Given other players or NPCs are in the same room
    When the player types "look"
    Then the other characters should be listed
    And NPC names should be color-coded (configurable)
    And players should see other player names

  Scenario: Look shows items on the ground
    Given items exist in the room
    When the player types "look"
    Then items should be listed in the room description
    And items should indicate they can be taken (if takeable)

  Scenario: Look updates when room state changes
    Given the player is in a room
    When an item is picked up or dropped
    Then another "look" should show updated item list

  Scenario: Look command can be abbreviated
    Given the player is in a room
    When the player types "l"
    Then the output should be identical to "look"
