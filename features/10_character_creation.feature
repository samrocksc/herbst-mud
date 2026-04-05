🔴 Feature: Character Creation
  As a new user
  I want to be able to create a character
  So that I can start playing the game

  Background:
    Given the SSH server is running
    And the database is connected

  Scenario: New user with no characters can create a character
    Given a user "newplayer@example.com" exists with password "password123"
    And the user has no characters
    When I log in as the user
    Then I should see a character creation form
    When I create a character "NewHero" with:
      | gender | "he/him" |
      | description | "A brave warrior" |
    Then character creation should be successful
    And the character should belong to the user

  Scenario: User can view their characters
    Given a user "player@example.com" exists
    And the user has characters:
      | name | "Warrior" |
      | name | "Mage" |
    When I request characters for the user
    Then I should see character "Warrior"
    And I should see character "Mage"

  Scenario: User needs to create character after login
    Given an authenticated user with no characters
    When I check if I need to create a character
    Then I should be prompted to create a character
    And I should not be able to enter the game until character is created

  Scenario: Character creation requires name
    Given a user "newplayer@example.com" with no characters
    When I attempt to create a character without a name
    Then I should see an error "Character name is required"

  Scenario: Character name must be unique per user
    Given a user "player@example.com" with character "Hero"
    When I attempt to create another character named "Hero"
    Then I should see an error "Character name already exists"

  Scenario: Character name minimum length validation
    Given a user "newplayer@example.com" exists
    When I attempt to create a character with name "AB"
    Then I should see an error "Name must be at least 3 characters"

  Scenario: Character name maximum length validation
    Given a user "newplayer@example.com" exists
    When I attempt to create a character with name exceeding 20 characters
    Then I should see an error "Name cannot exceed 20 characters"

  Scenario: Character name accepts alphanumeric and underscores
    Given a user "newplayer@example.com" exists
    When I create a character with name "Hero_123"
    Then character creation should be successful

  Scenario: Character name rejects special characters
    Given a user "newplayer@example.com" exists
    When I attempt to create a character with name "Hero@123!"
    Then I should see an error "Name can only contain letters, numbers, and underscores"

  Scenario: Character description optional field
    Given a user "newplayer@example.com" exists
    When I create a character with name "MinimalHero" and no description
    Then character creation should be successful

  Scenario: Character creation assigns default room
    Given a user "newplayer@example.com" exists
    When I create a character "NewHero"
    Then the character should have a starting room assigned
    And the character should be in the starting room

  Scenario: Character creation assigns default stats
    Given a user "newplayer@example.com" exists
    When I create a character "NewHero"
    Then the character should have default stats for the selected class
    And the character should start at level 1

  Scenario: Character creation with invalid gender
    Given a user "newplayer@example.com" exists
    When I create a character with invalid gender "unknown"
    Then I should see an error "Invalid gender selection"

  Scenario: User at character limit cannot create more
    Given a user "maxedplayer@example.com" already has 5 characters
    When I attempt to create another character
    Then I should see an error "Maximum character limit reached"
    And the character count should remain at 5