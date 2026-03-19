Feature: CRUD Characters
  As a game administrator
  I want to manage characters in the system
  So that I can create, read, update, and delete character records

  Scenario: Create a new character
    Given the admin is authenticated
    When I create a character with name "TestChar" and class "Warrior"
    Then the character should be saved in the database
    And the character should have a unique ID

  Scenario: Read character details
    Given a character "TestChar" exists in the database
    When I request character details for "TestChar"
    Then I should receive the character's name, class, and stats

  Scenario: Update character information
    Given a character "TestChar" exists
    When I update the character's class to "Mage"
    Then the character's class should be "Mage"

  Scenario: Delete a character
    Given a character "TestChar" exists
    When I delete the character "TestChar"
    Then the character should no longer exist in the database