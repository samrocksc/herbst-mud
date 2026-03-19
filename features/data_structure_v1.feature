Feature: Data Structure V1
  As a developer
  I want well-defined data structures
  So that the system can store and manage game data consistently

  Scenario: Character data structure
    Given I need to create a character
    When I define the character data structure
    Then it should include: ID, name, class, race, gender, stats, inventory, location
    And each field should have proper validation rules

  Scenario: Room data structure
    Given I need to create a room
    When I define the room data structure
    Then it should include: ID, name, description, exits, items, NPCs
    And exits should reference valid room IDs

  Scenario: Item data structure
    Given I need to create an item
    When I define the item data structure
    Then it should include: ID, name, description, type, stats, container

  Scenario: User data structure
    Given I need to create a user
    When I define the user data structure
    Then it should include: ID, username, email, password_hash, created_at