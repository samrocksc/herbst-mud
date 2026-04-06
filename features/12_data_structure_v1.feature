Feature: Data Structure V1
  As a developer
  I want well-defined and validated data structures
  So that the system can store and manage game data consistently

  Background:
    Given the data structures are defined in the codebase
    And validation rules are enforced at the API layer

  # Character Data Structure
  Scenario: Character data structure contains all required fields
    When I define or receive a character data structure
    Then it MUST contain the following fields:
      | field      | expectedType |
      | id         | UUID         |
      | userId     | UUID         |
      | name       | string       |
      | class      | string       |
      | race       | string       |
      | gender     | string       |
      | level      | integer      |
      | experience | integer      |
      | health     | integer      |
      | mana       | integer      |
      | stats      | JSON object  |
      | inventory  | JSON array   |
      | locationId | UUID         |
      | createdAt  | timestamp    |
      | updatedAt  | timestamp    |

  Scenario: Character UUID format is valid
    Given I create or retrieve a character
    Then the id field should be a valid UUID v4 format
    And the userId field should be a valid UUID v4 format
    And the locationId field should be a valid UUID v4 format

  Scenario: Character stats structure is valid JSON
    Given a character exists in the system
    When I retrieve the character's stats
    Then the stats field should be valid JSON
    And it should contain at minimum: strength, agility, intelligence, wisdom, constitution
    And all stat values should be integers

  Scenario: Character inventory structure is valid JSON array
    Given a character exists in the system
    When I retrieve the character's inventory
    Then the inventory field should be a valid JSON array
    And each item in the array should have: id, name, quantity

  # Room Data Structure
  Scenario: Room data structure contains all required fields
    When I define or receive a room data structure
    Then it MUST contain the following fields:
      | field       | expectedType |
      | id          | UUID         |
      | name        | string       |
      | description | text         |
      | x           | integer      |
      | y           | integer      |
      | z           | integer      |
      | exits       | JSON object  |
      | items       | JSON array   |
      | npcs        | JSON array   |
      | createdAt   | timestamp    |
      | updatedAt   | timestamp    |

  Scenario: Room coordinates are integers
    Given I create or retrieve a room
    Then the x coordinate should be an integer
    And the y coordinate should be an integer
    And the z coordinate should be an integer
    And z represents the floor/level number

  Scenario: Room exits structure is valid
    Given a room with defined exits exists
    When I examine the room's exits
    Then the exits field should be valid JSON object
    And each exit key should be a valid direction: north, south, east, west, up, down
    And each exit value should be the target room's UUID

  Scenario: Room exits reference valid room IDs
    Given a room has exits defined
    When I follow any exit
    Then the target room ID should reference an existing room
    And navigating through an exit should place the character in the correct room

  Scenario: Room items structure is valid
    Given a room contains items
    When I examine the room's items array
    Then each item should have: id, name, description, quantity
    And the items array should be a valid JSON array

  Scenario: Room NPCs structure is valid
    Given a room contains NPCs
    When I examine the room's NPCs array
    Then each NPC should have: id, name, description
    And the NPCs array should be a valid JSON array

  # User Data Structure
  Scenario: User data structure contains all required fields
    When I define or receive a user data structure
    Then it MUST contain the following fields:
      | field         | expectedType |
      | id            | UUID         |
      | username      | string       |
      | email         | string       |
      | passwordHash  | string       |
      | createdAt     | timestamp    |
      | updatedAt     | timestamp    |

  Scenario: User password is never returned in responses
    When I retrieve user data via the API
    Then the passwordHash field should NOT be included in the response
    And the actual password should NEVER be stored in plaintext

  # Item Data Structure
  Scenario: Item data structure contains all required fields
    When I define or receive an item data structure
    Then it MUST contain the following fields:
      | field       | expectedType |
      | id          | UUID         |
      | name        | string       |
      | description | text         |
      | type        | string       |
      | stats       | JSON object  |
      | container   | string       |
      | quantity    | integer      |

  Scenario: Item types are enumerated
    Given an item exists in the system
    When I check the item type
    Then it should be one of: weapon, armor, consumable, quest, misc

  Scenario: Item container field references valid location
    Given an item exists in a container or room
    When I check the item's container field
    Then the container value should be a room UUID, character UUID, or "world"

  # Timestamp Consistency
  Scenario: createdAt is set on record creation and never modified
    When I create a new record
    Then the createdAt timestamp should be set automatically
    And the createdAt value should not change on subsequent updates

  Scenario: updatedAt changes on every record modification
    Given a record exists with an initial updatedAt timestamp
    When I update the record
    Then the updatedAt timestamp should be updated to the current time
    And updatedAt should be different from createdAt after an update

  # Data Integrity
  Scenario: All UUID references point to existing records
    Given I validate a character's locationId
    Then a room with that UUID should exist in the database
    And a character's userId should reference an existing user
