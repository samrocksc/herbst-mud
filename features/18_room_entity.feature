Feature: Room Entity - Data Structure
  As a game developer
  I want a properly defined Room entity
  So that game world rooms are stored and managed correctly

  Background:
    Given the database schema is properly initialized
    And the Room model is defined in the codebase

  Scenario: Room entity has all required database fields
    When I examine the Room entity/model definition
    Then it should have the following fields with correct types:
      | field       | type        | constraints              |
      | id          | UUID        | primary key             |
      | name        | string      | not null                |
      | description | text        | not null                |
      | x           | integer     | coordinate              |
      | y           | integer     | coordinate              |
      | z           | integer     | floor/level             |
      | exits       | JSON object | direction -> room UUID  |
      | items       | JSON array  | items in room           |
      | npcs        | JSON array  | NPCs in room            |
      | createdAt   | timestamp   | auto-set                |
      | updatedAt   | timestamp   | auto-updated            |

  Scenario: Room id is a valid UUID
    When I create a new room
    Then the generated id should be a valid UUID v4
    And the id should be unique across all rooms

  Scenario: Room name is required and not empty
    When I attempt to create a room with name ""
    Then the validation should fail
    And the error should indicate name is required

  Scenario: Room description is required and not empty
    When I attempt to create a room with description ""
    Then the validation should fail
    And the error should indicate description is required

  Scenario: Room coordinates are integers
    When I create or retrieve a room
    Then the x coordinate should be an integer
    And the y coordinate should be an integer
    And the z coordinate should be an integer

  Scenario: Room coordinates can be zero
    When I create a room at coordinates (0, 0, 0)
    Then the room should be created successfully

  Scenario: Room coordinates can be negative (representing world map)
    When I create a room at coordinates (-10, 5, -2)
    Then the room should be created successfully

  Scenario: Z-coordinate represents floor/level
    Given I am in a room with z = 0 (ground floor)
    When I move up
    Then my z-coordinate should increase to 1
    And I should be on floor 1

  Scenario: Room exits JSON structure is valid
    When I examine a room's exits JSON
    Then it should be a valid JSON object
    And each key should be a valid direction: north, south, east, west, up, down
    And each value should be a valid UUID referencing another room

  Scenario: Room with no exits is valid
    When I create a room with an empty exits object
    Then the room should be created successfully
    And attempting to move in any direction should fail with "You can't go that way."

  Scenario: Room exit directions are validated
    When I attempt to set an exit direction to "sideways"
    Then the validation should fail
    And the exit direction should be one of: north, south, east, west, up, down

  Scenario: Room exit target must be a valid room UUID
    When I attempt to set an exit to a non-existent room UUID "00000000-0000-0000-0000-000000000000"
    Then the validation should fail
    Or navigating to that exit should result in an error

  Scenario: Room items JSON structure is valid
    When I examine a room's items array
    Then it should be a valid JSON array
    And each item should have at minimum: id, name

  Scenario: Room NPCs JSON structure is valid
    When I examine a room's NPCs array
    Then it should be a valid JSON array
    And each NPC should have at minimum: id, name

  Scenario: Room can contain multiple items
    Given a room exists
    When I add items "Iron Sword" and "Health Potion" to the room
    Then the room's items array should contain both items
    And both items should be visible when entering the room

  Scenario: Room can contain multiple NPCs
    Given a room exists
    When I add NPCs "Guard" and "Merchant" to the room
    Then the room's NPCs array should contain both NPCs
    And both NPCs should be visible when entering the room

  Scenario: createdAt is set automatically on creation
    When I create a new room
    Then the createdAt timestamp should be set automatically

  Scenario: updatedAt is updated on room modification
    Given a room exists
    When I update the room's name, description, or exits
    Then the updatedAt timestamp should be updated

  Scenario: Two rooms at same x,y,z are the same location
    Given I have created a room at (5, 5, 0)
    When I create another room at (5, 5, 0)
    Then this should either be prevented (unique constraint)
    Or both rooms should be considered the same room

  Scenario: Room ID used in navigation is persisted
    Given a room "Tavern" has ID "room-uuid-123"
    When a character navigates to the Tavern
    Then the character's locationId should be set to "room-uuid-123"
