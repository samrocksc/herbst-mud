Feature: CRUD Characters
  As a game administrator
  I want to create, read, update, and delete characters
  So that I can manage the game world effectively

  Scenario: Create a new character
    Given the game database is empty
    When I create a character named "Hero1" with class "Warrior"
    Then the character "Hero1" should exist in the database
    And the character should have class "Warrior"

  Scenario: Read character details
    Given a character "Mage1" exists with class "Mage"
    When I query for character "Mage1"
    Then I should receive the character's details
    And the class should be "Mage"

  Scenario: Update character attributes
    Given a character "Rogue1" exists with level 1
    When I update "Rogue1" to level 5
    Then the character "Rogue1" should have level 5

  Scenario: Delete a character
    Given a character "ToBeDeleted" exists
    When I delete character "ToBeDeleted"
    Then the character "ToBeDeleted" should not exist

  Scenario: List all characters
    Given the following characters exist:
      | name    | class   |
      | Char1   | Warrior |
      | Char2   | Mage    |
      | Char3   | Rogue   |
    When I list all characters
    Then I should see 3 characters