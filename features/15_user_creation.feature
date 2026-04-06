Feature: User Creation
  As a new player
  I want to create a user account
  So that I can access the game and manage my characters

  Background:
    Given the user registration API is available at /api/users/register
    And the API is reachable and functional

  Scenario: Create user with valid data
    When I submit a registration request with:
      | field    | value              |
      | username | newplayer1         |
      | email    | newplayer@test.com |
      | password | SecurePass123      |
    Then the response status should be 201 Created
    And the user should be created in the database
    And the user should have a unique UUID
    And the password should be hashed (not stored as plaintext)
    And I should receive a confirmation response

  Scenario: Create user with existing username fails
    Given a user with username "taken" already exists
    When I submit a registration request with username "taken" email "new@test.com" password "SecurePass123"
    Then the response status should be 409 Conflict
    And the error message should contain "username already taken"

  Scenario: Create user with existing email fails
    Given a user with email "existing@test.com" already exists
    When I submit a registration request with username "newuser" email "existing@test.com" password "SecurePass123"
    Then the response status should be 409 Conflict
    And the error message should contain "email already in use"

  Scenario: Create user with invalid email format fails
    When I submit a registration request with username "validuser" email "notanemail" password "SecurePass123"
    Then the response status should be 400 Bad Request
    And the error message should contain "invalid email"

  Scenario: Create user with email missing @ symbol
    When I submit a registration request with username "validuser" email "validemail.com" password "SecurePass123"
    Then the response status should be 400 Bad Request

  Scenario: Create user with password too short fails
    When I submit a registration request with username "validuser" email "valid@test.com" password "short"
    Then the response status should be 400 Bad Request
    And the error message should contain "password must be at least"

  Scenario: Create user with username too short fails
    When I submit a registration request with username "ab" email "valid@test.com" password "SecurePass123"
    Then the response status should be 400 Bad Request
    And the error message should contain "username must be at least"

  Scenario: Create user with username too long fails
    When I submit a registration request with username "thisusernameiswaytoolongtobevalid" email "valid@test.com" password "SecurePass123"
    Then the response status should be 400 Bad Request
    And the error message should contain "username must be at most"

  Scenario: Create user with empty username fails
    When I submit a registration request with username "" email "valid@test.com" password "SecurePass123"
    Then the response status should be 400 Bad Request

  Scenario: Create user with empty email fails
    When I submit a registration request with username "validuser" email "" password "SecurePass123"
    Then the response status should be 400 Bad Request

  Scenario: Create user with empty password fails
    When I submit a registration request with username "validuser" email "valid@test.com" password ""
    Then the response status should be 400 Bad Request

  Scenario: Username accepts alphanumeric characters
    When I submit a registration request with username "Player123" email "valid@test.com" password "SecurePass123"
    Then the response status should be 201 Created

  Scenario: Username accepts underscores and hyphens
    When I submit a registration request with username "player_name-123" email "valid@test.com" password "SecurePass123"
    Then the response status should be 201 Created

  Scenario: Username rejects spaces
    When I submit a registration request with username "player name" email "valid@test.com" password "SecurePass123"
    Then the response status should be 400 Bad Request

  Scenario: New user can immediately create a character
    When I register a new user "NewUserTest" email "new@test.com" password "SecurePass123"
    And I log in with those credentials
    Then I should be able to create a character

  Scenario: Password is hashed not stored in plaintext
    When I create a new user
    Then the database should NOT contain the plaintext password
    And the stored password should be a bcrypt or argon2 hash

  Scenario: Duplicate simultaneous registration for same username
    Given two registration requests are submitted at the same time with username "sameuser"
    Then exactly one should succeed with 201
    And the other should fail with 409
