Feature: User Entity - Data Structure
  As a game developer
  I want a properly defined User entity
  So that user accounts are stored and validated correctly

  Background:
    Given the database schema is properly initialized
    And the User model is defined in the codebase

  Scenario: User entity has all required database fields
    When I examine the User entity/model definition
    Then it should have the following fields with correct types:
      | field         | type        | constraints          |
      | id            | UUID        | primary key          |
      | username      | string      | unique, 3-20 chars   |
      | email         | string      | unique, valid format |
      | passwordHash  | string      | not null             |
      | createdAt     | timestamp   | auto-set             |
      | updatedAt     | timestamp   | auto-updated         |

  Scenario: User id is a valid UUID
    When I create a new user
    Then the generated id should be a valid UUID v4
    And the id should be unique across all users

  Scenario: Username length validation - minimum
    When I attempt to create a user with username "ab"
    Then the validation should fail
    And the error should indicate username must be at least 3 characters

  Scenario: Username length validation - maximum
    When I attempt to create a user with username "thisusernameiswaytoolongforthefield"
    Then the validation should fail
    And the error should indicate username must be at most 20 characters

  Scenario: Username uniqueness is enforced
    Given a user with username "UniqueTest" exists
    When I attempt to create another user with username "UniqueTest"
    Then the validation should fail with a uniqueness constraint error

  Scenario: Email format is validated
    When I attempt to create a user with email "notavalidemail"
    Then the validation should fail
    And the error should indicate a valid email is required

  Scenario: Email uniqueness is enforced
    Given a user with email "duplicate@test.com" exists
    When I attempt to create another user with email "duplicate@test.com"
    Then the validation should fail with a uniqueness constraint error

  Scenario: Password is never returned in API responses
    When I retrieve user data via the API
    Then the passwordHash field should be excluded from the response
    And no endpoint should ever return the password in plaintext

  Scenario: Password is hashed using a secure algorithm
    When I create a new user with password "MyPassword123"
    Then the stored password should NOT be "MyPassword123"
    And the stored password should be a bcrypt or argon2 hash
    And verifying with the correct password should succeed
    And verifying with a wrong password should fail

  Scenario: Passwords are case-sensitive
    Given a user exists with password "Password123"
    When I verify with "password123"
    Then the verification should fail

  Scenario: createdAt is set automatically on creation
    When I create a new user account
    Then the createdAt timestamp should be set automatically
    And the value should be approximately the current time (within 1 second)

  Scenario: updatedAt is set automatically on creation and update
    When I create a new user account
    Then updatedAt should be set to the same time as createdAt

    When I update the user's email
    Then updatedAt should be updated to the current time
    And updatedAt should be different from createdAt

  Scenario: User timestamps are in UTC
    When I create a user from any timezone
    Then the createdAt timestamp should be stored in UTC
    And the updatedAt timestamp should be stored in UTC

  Scenario: User soft delete (if implemented)
    Given a user exists in the system
    When I "delete" the user account
    Then the user record should either be removed from the database
    Or the user record should be marked as deleted (soft delete)
    And the username should become available again for registration
