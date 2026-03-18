Feature: Character Entity
  As a developer
  I want a Character entity model
  So that I can manage character data in the system

  Scenario: Character entity fields
    Given I need to model a Character
    When I create the Character entity
    Then it should have the following fields:
      | field       | type   |
      | ID          | int    |
      | Name        | string |
      | UserID      | int    |
      | Class       | string |
      | Race        | string |
      | Gender      | string |
      | Level       | int    |
      | Experience  | int    |
      | Health      | int    |
      | Mana        | int    |

  Scenario: Character stats
    Given a Character entity exists
    Then it should have stats: strength, dexterity, constitution, intelligence, wisdom, charisma

  Scenario: Character inventory
    Given a Character entity exists
    Then it should have an inventory relationship