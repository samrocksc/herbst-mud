Feature: Character Entity
  As a game developer
  I want a Character entity
  So that I can manage player characters in the system

  Scenario: Character entity fields
    Given I need a Character entity
    When I define the Character model
    Then it should have the following fields:
      | field      | type      |
      | id         | UUID      |
      | userId     | UUID      |
      | name       | string    |
      | class      | string    |
      | race       | string    |
      | gender     | string    |
      | level      | int       |
      | experience | int       |
      | health     | int       |
      | mana       | int       |
      | stats      | JSON      |
      | inventory  | JSON      |
      | locationId | UUID      |
      | createdAt  | timestamp |
      | updatedAt  | timestamp |

  Scenario: Character validation
    Given I have character data
    When I validate a Character
    Then name should be 2-20 characters
    And class must be a valid class
    And race must be a valid race