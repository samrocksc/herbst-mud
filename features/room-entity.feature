Feature: Room Entity Data Structure
  As a game architect
  I want the Room entity properly implemented
  So that the game world spaces are stored correctly

  Background:
    Given the database schema is defined
    And the Room entity follows DATA_STRUCTURE_V1

  Scenario: Room holds characters
    Given a Room entity
    Then it can contain multiple characters
    And characters in a room can interact
    And the room tracks which characters are present

  Scenario: Room holds equipment
    Given a Room entity
    Then it can contain equipment items
    And items on the floor are visible to players
    And items can be picked up by characters

  Scenario: Room has exits structure
    Given a Room entity
    Then it has an exits field (JSON object)
    And each exit maps direction to target room UUID

  Scenario: Room coordinates are stored
    Given a Room entity
    Then it has x, y, z coordinates
    And coordinates can be negative (world map)

  Scenario: Room z-coordinate represents floor
    Given a room at z = 0 (ground floor)
    When I move up
    Then z should increase to 1

  Scenario: Room name is required
    When I attempt to create a room with name ""
    Then the validation should fail
    And the error should indicate name is required

  Scenario: Room description is required
    When I attempt to create a room with description ""
    Then the validation should fail
    And the error should indicate description is required

  Scenario: Room has unique constraint on name
    When I create two rooms with the same name
    Then the second creation should fail

  Scenario: Room coordinates can be zero
    When I create a room at (0, 0, 0)
    Then the room should be created successfully

  Scenario: Room can have empty exits
    When I create a room with no exits
    Then the room should be created successfully
    And attempting to move in any direction should fail

  Scenario: Room has createdAt and updatedAt timestamps
    When I create and update a room
    Then createdAt is set on creation
    And updatedAt is updated on changes
