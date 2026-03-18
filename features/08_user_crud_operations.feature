Feature: User CRUD Operations
  As a system administrator
  I want to manage user accounts
  So that I can maintain the player registry

  Scenario: Create a new user
    Given no user "player1" exists
    When I create a user "player1" with email "player1@example.com"
    Then the user "player1" should exist
    And the email should be "player1@example.com"

  Scenario: Read user profile
    Given a user "admin" exists with role "admin"
    When I query for user "admin"
    Then I should receive the user's profile
    And the role should be "admin"

  Scenario: Update user email
    Given a user "olduser" exists with email "old@example.com"
    When I update the email to "new@example.com"
    Then the user should have email "new@example.com"

  Scenario: Delete a user
    Given a user "temporary" exists
    When I delete user "temporary"
    Then the user "temporary" should not exist

  Scenario: List all users
    Given the following users exist:
      | username | role   |
      | user1    | player |
      | user2    | player |
      | admin    | admin  |
    When I list all users
    Then I should see 3 users