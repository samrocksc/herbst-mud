Feature: Data Structure V1
  As a developer
  I want well-defined data structures
  So that the game has consistent data models

  Scenario: User data structure
    Given I need a user data model
    When I define the user structure
    Then it should include: id, username, email, password_hash, role, created_at

  Scenario: Character data structure
    Given I need a character data model
    When I define the character structure
    Then it should include: id, name, user_id, class, race, gender, level, experience, stats

  Scenario: Room data structure
    Given I need a room data model
    When I define the room structure
    Then it should include: id, name, description, exits, items, npcs

  Scenario: Item data structure
    Given I need an item data model
    When I define the item structure
    Then it should include: id, name, description, type, weight, value