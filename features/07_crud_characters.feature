Feature: CRUD Characters
  As a game administrator
  I want to manage characters in the system
  So that I can create, read, update, and delete character records

  Background:
    Given the admin backoffice is accessible
    And I am authenticated as an admin
    And a test user "admin_test_user" exists

  Scenario: Create a new character with valid data
    When I submit a create character request with name "Aragorn" class "Warrior" race "Human" gender "Male"
    Then the response status should be 201 Created
    And the character record should exist in the database
    And the character should have a unique UUID
    And the character's class should be "Warrior"
    And the character's level should be 1
    And the character should start with default stats for "Warrior"

  Scenario: Create character with duplicate name fails
    Given a character named "DupName" exists
    When I submit a create character request with name "DupName" class "Mage" race "Elf" gender "Male"
    Then the response status should be 409 Conflict
    And the error message should contain "name already taken"

  Scenario: Create character with invalid class fails
    Given a valid user exists
    When I submit a create character request with name "NewChar" class "NotAClass" race "Human" gender "Male"
    Then the response status should be 400 Bad Request
    And the error message should contain "invalid class"

  Scenario: Create character with name too short
    Given a valid user exists
    When I submit a create character request with name "A" class "Warrior" race "Human" gender "Male"
    Then the response status should be 400 Bad Request
    And the error message should contain "name must be"

  Scenario: Read character details by ID
    Given a character "ReadTestChar" exists in the database
    When I request character details by ID
    Then the response should include: id, name, class, race, gender, level, experience, health, mana, stats, inventory, locationId
    And the name field should match "ReadTestChar"

  Scenario: Read character that does not exist
    When I request character details for a non-existent UUID "00000000-0000-0000-0000-000000000000"
    Then the response status should be 404 Not Found

  Scenario: Update character class
    Given a character "UpdateTestChar" with class "Warrior" exists
    When I update the character's class to "Mage"
    Then the response status should be 200 OK
    And the character's class should now be "Mage"
    And the character's stats should reflect Mage base stats

  Scenario: Update character level and experience
    Given a character "ExpTestChar" at level 1 exists
    When I update the character's experience to 1000
    Then the character should level up to level 2
    And the character's stats should increase accordingly

  Scenario: Update character location
    Given a character "LocTestChar" exists at the starting room
    And a room "Tavern" exists in the world
    When I update the character's location to the "Tavern" room
    Then the response status should be 200 OK
    And the character's locationId should match the Tavern room ID

  Scenario: Update character inventory
    Given a character "InvTestChar" exists with empty inventory
    And an item "Iron Sword" exists in the world
    When I add the item "Iron Sword" to the character's inventory
    Then the response status should be 200 OK
    And the character's inventory should contain "Iron Sword"

  Scenario: Delete a character
    Given a character "DeleteTestChar" exists in the database
    When I delete the character "DeleteTestChar"
    Then the response status should be 204 No Content
    And the character should no longer exist in the database

  Scenario: Delete character that does not exist
    When I attempt to delete a non-existent character UUID "00000000-0000-0000-0000-000000000000"
    Then the response status should be 404 Not Found

  Scenario: List all characters for a user
    Given a user has 3 characters: "Char1" "Char2" "Char3"
    When I request all characters for that user
    Then the response should contain exactly 3 character records
    And each record should include name and class

  Scenario: Character stats validation on creation
    When I create a character "StatChar" with class "Priest" race "Dwarf" gender "Male"
    Then the character should have valid stats for Priest class
    And health should be greater than mana
    And strength should be less than wisdom
