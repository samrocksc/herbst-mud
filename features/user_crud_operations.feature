Feature: User CRUD Operations
  As a system administrator
  I want to manage user accounts
  So that I can create, read, update, and delete users

  Scenario: Create a new user
    Given I am an admin
    When I create a user with username "newuser" and email "user@example.com"
    Then the user should be created successfully
    And the user should have a unique ID

  Scenario: View all users
    Given multiple users exist in the system
    When I request a list of all users
    Then I should receive a list of user records
    And each record should contain username and email

  Scenario: Update user email
    Given a user "newuser" exists
    When I update the user's email to "newemail@example.com"
    Then the user's email should be "newemail@example.com"

  Scenario: Delete a user
    Given a user "newuser" exists
    When I delete the user "newuser"
    Then the user should no longer exist