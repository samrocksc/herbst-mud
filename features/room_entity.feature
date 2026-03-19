Feature: Room Entity
  As a game developer
  I want a Room entity
  So that I can manage game rooms in the system

  Scenario: Room entity fields
    Given I need a Room entity
    When I define the Room model
    Then it should have the following fields:
      | field       | type      |
      | id          | UUID      |
      | name        | string    |
      | description | text      |
      | x           | int       |
      | y           | int       |
      | z           | int       |
      | exits       | JSON      |
      | items       | JSON      |
      | npcs        | JSON      |
      | createdAt   | timestamp |
      | updatedAt   | timestamp |

  Scenario: Room exits structure
    Given a room with exits
    When I examine the exits
    Then each exit should have: direction, targetRoomId
    And direction should be: north, south, east, west, up, or down