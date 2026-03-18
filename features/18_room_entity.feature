Feature: Room Entity
  As a developer
  I want a Room entity model
  So that I can manage room data in the game world

  Scenario: Room entity fields
    Given I need to model a Room
    When I create the Room entity
    Then it should have the following fields:
      | field        | type   |
      | ID           | int    |
      | Name         | string |
      | Description  | string |
      | AreaID       | int    |
      | CoordinateX  | int    |
      | CoordinateY  | int    |
      | CoordinateZ  | int    |

  Scenario: Room exits
    Given a Room entity exists
    Then it should have exits: north, south, east, west, up, down

  Scenario: Room items and NPCs
    Given a Room entity exists
    Then it should have a one-to-many relationship with Items
    And it should have a one-to-many relationship with NPCs