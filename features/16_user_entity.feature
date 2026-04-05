🔴 Feature: User Entity - Data Structure
  As a game developer
  I want a properly structured User entity
  So that players can have accounts with proper relationships

  Background:
    Given the database is connected
    And the ent schema is generated

  Scenario: User is a person entity
    Given I create a new user
    Then the user should represent a person
    And the user should have a unique ID
    And the user should have an email address

  Scenario: User can have up to three characters
    Given a user exists with ID 1
    When I create character 1 for the user
    And I create character 2 for the user
    And I create character 3 for the user
    Then the user should have 3 characters
    When I attempt to create a 4th character
    Then I should see an error "Maximum character limit reached"
    And the user should still have only 3 characters

  Scenario: User can have god mode (unkillable)
    Given I create a user with god mode enabled
    Then the user should have is_god = true
    And the user should be unkillable in combat
    And the user should have special god privileges

  Scenario: User one-to-many relationship with characters
    Given a user exists with ID 1
    And I create character "Hero" for the user
    And I create character "Sidekick" for the user
    When I query the user's characters
    Then I should see "Hero" and "Sidekick"
    And each character should reference the user

  Scenario: User can be an admin
    Given I create a user with is_admin = true
    Then the user should have admin privileges
    And the user can access admin commands
    And the user can manage other users

  Scenario: User email must be unique
    Given a user "existing@example.com" exists
    When I attempt to create another user with email "existing@example.com"
    Then the creation should fail
    And an error about duplicate email should be returned

  Scenario: User entity has created_at timestamp
    Given I create a new user
    Then the user should have a created_at timestamp
    And the timestamp should be set to creation time

  Scenario: User entity has updated_at timestamp
    Given a user exists
    When I update the user's email
    Then the updated_at timestamp should be updated

  Scenario: Regular user cannot access admin functions
    Given a regular user exists with is_admin: false
    When the user attempts to access admin endpoints
    Then access should be denied
    And an error about insufficient permissions should be returned

  Scenario: User can be deleted (soft or hard delete)
    Given a user exists with characters
    When I delete the user
    Then the user should no longer exist
    And the characters should be handled according to deletion policy

  Scenario: User password is required
    When I attempt to create a user without a password
    Then the creation should fail
    And an error about required password should be returned

  Scenario: User email is required
    When I attempt to create a user without an email
    Then the creation should fail
    And an error about required email should be returned