Feature: Room Entity Data Structure (Issue #18)
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

  Scenario: Room description
    Given a Room entity
    Then it should have a description field
    And the description is shown when entering or looking
    And descriptions can be rich text

  Scenario: Room exits
    Given a Room entity
    Then it should have exits
    And exits connect to other rooms
    And exits have directions: N, S, E, W, U, D
    And exits can be one-way or bidirectional

  Scenario: Room atmosphere
    Given a Room entity
    Then it should have atmosphere field
    And valid atmospheres: air, water, wind
    And atmosphere affects gameplay (underwater rooms, etc.)

  Scenario: Room unique ID
    Given a Room entity
    Then it should have a unique identifier
    And rooms can be referenced by ID
    And IDs should be human-readable where possible

  Scenario: Room coordinates optional
    Given a Room entity
    Then it may have x, y, z coordinates
    And coordinates help with map visualization
    And coordinates are for admin tools