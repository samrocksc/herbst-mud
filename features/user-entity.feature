Feature: User Entity Data Structure (Issue #16)
  As a game architect
  I want the User entity properly implemented
  So that player accounts are stored correctly

  Background:
    Given the database schema is defined
    And the User entity follows DATA_STRUCTURE_V1

  Scenario: User is a person
    Given a User entity
    Then it should represent a real person
    And it should have unique identification

  Scenario: User can have multiple characters
    Given a User entity
    Then it can have up to three characters
    And the characters relationship is one-to-many

  Scenario: User god mode flag
    Given a User entity
    Then it can have a god_mode boolean
    And when god_mode is true, the character is unkillable
    And this is for admin/debug purposes

  Scenario: User admin flag
    Given a User entity
    Then it should have is_admin boolean
    And admins can access special commands
    And admins can see hidden game state

  Scenario: User email unique
    Given multiple users might register
    Then each user email must be unique
    And duplicate emails should be rejected

  Scenario: User timestamps
    Given a User entity
    Then it should have created_at timestamp
    And it should have updated_at timestamp
    And these track account activity