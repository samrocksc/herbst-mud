Feature: User CRUD Operations
  As a system administrator
  I want to manage user accounts
  So that I can create, read, update, and delete users

  Background:
    Given the user management API is accessible
    And I am authenticated as a system admin

  Scenario: Create a new user with valid data
    When I submit a create user request with username "newplayer" email "newplayer@test.com" password "SecurePass123"
    Then the response status should be 201 Created
    And the user record should exist in the database
    And the user should have a unique UUID
    And the password should be hashed (not stored in plaintext)

  Scenario: Create user with duplicate username fails
    Given a user "takenuser" already exists
    When I submit a create user request with username "takenuser" email "another@test.com" password "SecurePass123"
    Then the response status should be 409 Conflict
    And the error message should contain "username already taken"

  Scenario: Create user with duplicate email fails
    Given a user with email "duplicate@test.com" already exists
    When I submit a create user request with username "differentuser" email "duplicate@test.com" password "SecurePass123"
    Then the response status should be 409 Conflict
    And the error message should contain "email already in use"

  Scenario: Create user with invalid email format
    When I submit a create user request with username "validuser" email "notanemail" password "SecurePass123"
    Then the response status should be 400 Bad Request
    And the error message should contain "invalid email"

  Scenario: Create user with short password fails
    When I submit a create user request with username "validuser" email "valid@test.com" password "short"
    Then the response status should be 400 Bad Request
    And the error message should contain "password must be at least"

  Scenario: Create user with username too short
    When I submit a create user request with username "ab" email "valid@test.com" password "SecurePass123"
    Then the response status should be 400 Bad Request
    And the error message should contain "username must be at least"

  Scenario: Create user with username too long
    When I submit a create user request with username "thisusernameiswaytoolongtobevalid" email "valid@test.com" password "SecurePass123"
    Then the response status should be 400 Bad Request
    And the error message should contain "username must be at most"

  Scenario: Read user by ID
    Given a user "ReadUser" exists in the database
    When I request user details by ID
    Then the response should include: id, username, email, createdAt, updatedAt
    And the password hash should NOT be included in the response

  Scenario: Read user that does not exist
    When I request user details for a non-existent UUID "00000000-0000-0000-0000-000000000000"
    Then the response status should be 404 Not Found

  Scenario: Update user email
    Given a user "UpdateEmailUser" with email "old@test.com" exists
    When I update the user's email to "new@test.com"
    Then the response status should be 200 OK
    And the user's email should now be "new@test.com"
    And the updatedAt timestamp should be updated

  Scenario: Update user password
    Given a user "UpdatePassUser" exists with password "OldPass123"
    When I update the user's password to "NewPass456"
    Then the response status should be 200 OK
    And the old password "OldPass123" should no longer work
    And the new password "NewPass456" should work

  Scenario: Update username
    Given a user "OldUsername" exists
    When I update the username to "NewUsername"
    Then the response status should be 200 OK
    And the username should now be "NewUsername"

  Scenario: Update username to one that is taken
    Given a user "User1" exists
    And a user "User2" exists
    When I update "User1" username to "User2"
    Then the response status should be 409 Conflict
    And the error should contain "username already taken"

  Scenario: Delete a user
    Given a user "DeleteUser" exists in the database
    When I delete the user "DeleteUser"
    Then the response status should be 204 No Content
    And the user should no longer exist in the database

  Scenario: Delete user that does not exist
    When I attempt to delete a non-existent user UUID "00000000-0000-0000-0000-000000000000"
    Then the response status should be 404 Not Found

  Scenario: List all users (admin)
    Given 5 users exist in the system
    When I request a list of all users
    Then the response should contain 5 user records
    And each record should NOT include password hashes

  Scenario: User createdAt timestamp is set automatically
    When I create a new user "TimestampUser" with email "time@test.com" password "SecurePass123"
    Then the createdAt timestamp should be set to the current time
    And the updatedAt timestamp should match createdAt
