🔴 Feature: Room Entity - Data Structure
  As a game developer
  I want a properly structured Room entity
  So that players can navigate a rich game world

  Background:
    Given the database is connected
    And the ent schema is generated

  Scenario: Room can hold characters
    Given a room exists with ID 1
    And 3 characters are in the room
    When I query the room's characters
    Then I should see 3 characters
    And each character should have currentRoomId = 1

  Scenario: Room can hold equipment
    Given a room exists with ID 1
    And 2 equipment items are in the room
    When I query the room's items
    Then I should see 2 equipment items

  Scenario: Room has description
    Given I create a room with:
      | name | "Dark Cave" |
      | description | "A dark and damp cave with mysterious sounds echoing..." |
    Then the room should have the description
    And players should see the description when entering

  Scenario: Room has exits
    Given I create a room with exits:
      | direction | targetRoomId |
      | north | 2 |
      | east | 3 |
    Then the room should have exits to rooms 2 and 3
    And players should be able to move north and east

  Scenario: Room has atmosphere
    Given I create a room with atmosphere "water"
    Then the room should have atmosphere "water"
    And underwater rooms should have special mechanics
    And air rooms should be normal
    And wind rooms should have special effects

  Scenario: Room has is_peerable flag
    Given I create a room with is_peerable = true
    Then players should be able to peer into this room
    When I create a room with is_peerable = false
    Then players should NOT be able to peer into this room
    And players should NOT be able to peer FROM this room

  Scenario: Room has isStartingRoom flag
    Given I create a room with isStartingRoom = true
    Then new players should start in this room
    And only one room should have this flag per game world

  Scenario: Room exits are stored as JSON map
    Given a room exists
    Then the exits field should be a map of direction to room ID
    And I should be able to add new exits
    And I should be able to remove exits

  Scenario: Room has unique name
    Given a room "The Hole" exists
    When I attempt to create another room named "The Hole"
    Then the creation should fail or name is not unique (per design)

  Scenario: Room can have multiple exits in same direction
    Given rooms exist and exits can lead to different areas
    When checking room exit structure
    Then each exit direction should map to exactly one target room
    Or multiple exits should be distinguished by some mechanism

  Scenario: Room has unique ID
    Given I create a room
    Then the room should have a unique ID
    And no other room should have the same ID

  Scenario: Room description can be empty
    Given I create a room with no description
    Then the room should have an empty description
    And players should see a default description

  Scenario: Room can have NPCs
    Given a room "Town Square" exists
    When I add an NPC "Guard" to the room
    Then the room should have 1 NPC
    And players entering the room should see the NPC

  Scenario: Room has created_at timestamp
    Given I create a room
    Then the room should have a created_at timestamp
    And the timestamp should reflect creation time

  Scenario: Room can be updated
    Given a room "Old Name" exists
    When I update the room name to "New Name"
    Then the room should have name "New Name"
    And other room properties can be updated