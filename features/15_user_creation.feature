🔴 Feature: User Creation
  As a new player
  I want to create a user account when I SSH into the MUD
  So that I can start playing

  Background:
    Given the SSH server is running on port 4444
    And the database is connected

  Scenario: New user sees login or create account screen
    When I SSH to herbst-mud
    Then I should see turtle ASCII art
    And I should see options for "Log In" or "Create Account"

  Scenario: Create account flow asks for email
    Given I am on the welcome screen
    When I select "Create Account"
    Then I should see a prompt for my email

  Scenario: Create account flow asks for password
    Given I am creating a new account
    And I entered my email "newuser@example.com"
    Then I should see a prompt for my password
    And the password should be obfuscated as "****"

  Scenario: Password verification during account creation
    Given I am creating a new account
    And I entered my email "newuser@example.com"
    And I entered my password "password123"
    When I enter a different password "password456" for confirmation
    Then I should see an error "Passwords do not match"
    And I should be asked to re-enter my password

  Scenario: Successful account creation
    Given I am creating a new account
    And I entered my email "newuser@example.com"
    And I entered my password "password123"
    And I confirmed my password "password123"
    Then my account should be created via REST API
    And I should be logged in automatically
    And I should see the character creation screen

  Scenario: Duplicate email shows error
    Given a user "existing@example.com" already exists
    When I try to create an account with "existing@example.com"
    Then I should see an error "Email already in use"
    And I should be asked to use a different email

  Scenario: Invalid email format rejected
    Given I am creating a new account
    When I enter an invalid email "not-an-email"
    Then I should see an error "Please enter a valid email address"

  Scenario: Password minimum length requirement
    Given I am creating a new account
    When I enter my email "newuser@example.com"
    And I enter a password "12345" (less than 8 characters)
    Then I should see an error "Password must be at least 8 characters"

  Scenario: Password maximum length
    Given I am creating a new account
    When I enter my email "newuser@example.com"
    And I enter a password exceeding 128 characters
    Then I should see an error "Password is too long"

  Scenario: Empty email rejected
    Given I am on the create account screen
    When I press enter without entering an email
    Then I should see an error "Email is required"

  Scenario: Empty password rejected
    Given I am creating a new account
    And I entered my email "newuser@example.com"
    When I press enter without entering a password
    Then I should see an error "Password is required"

  Scenario: New user starts with no characters
    Given I successfully create an account
    When I check my character list
    Then it should be empty
    And I should be prompted to create a character

  Scenario: New user has is_admin: false by default
    Given I successfully create an account
    When I query my user record
    Then is_admin should be false

  Scenario: Cannot create account while already logged in
    Given I am already logged in as a user
    When I attempt to access the create account flow
    Then I should not see the create account option
    Or I should be redirected to the game

  Scenario: Account creation rate limiting
    Given I attempt to create multiple accounts rapidly
    When I exceed the rate limit
    Then I should see "Too many requests, please try again later"
    And subsequent requests should be blocked