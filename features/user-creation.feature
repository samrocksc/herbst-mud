Feature: User Creation (Issue #15)
  As a new player
  I want to create a user account when first connecting
  So that I can save my progress and play the game

  Background:
    Given the SSH server is running
    And I am connecting for the first time

  Scenario: Welcome screen with options
    When I SSH to herbst-mud
    Then I should see turtle ASCII art
    And I should see "Log In" option
    And I should see "Create Account" option

  Scenario: Choose create account
    Given I am at the welcome screen
    When I select "Create Account"
    Then I should be prompted for my email
    And after email, I should be prompted for password

  Scenario: Create account with email
    Given I am creating an account
    When I enter email "newplayer@example.com"
    And I enter a valid password
    Then my account should be created
    And I should see a confirmation message

  Scenario: Password obfuscation during creation
    Given I am at the password prompt
    When I type my password
    Then the characters should appear as ****
    And my actual password should not be visible

  Scenario: Password confirmation
    Given I am creating an account
    When I enter password "secret123"
    Then I should be asked to confirm the password
    And if passwords don't match, I should see an error
    And I should be allowed to re-enter

  Scenario: Successful account creation
    Given I enter valid email and matching passwords
    When I submit the account creation form
    Then the account should be created
    And I should be logged in automatically
    And I should see the character creation or selection screen

  Scenario: REST API account creation
    Given I want to create an account via API
    When I send POST to /api/users with email and password
    Then the response should be 201 Created
    And REST API should follow proper HTTP conventions