Feature: Character Login (Issue #11)
  As a returning player
  I want to log in with my username and password
  So that I can access my characters and continue playing

  Background:
    Given the SSH server is running
    And I have an existing account with username "player1"

  Scenario: Login screen appears on SSH connect
    When I SSH to herbst-mud
    Then I should see a screen asking for username
    And after entering username, I should see a password prompt

  Scenario: Successful login
    Given I have credentials username "player1" and password "secret"
    When I enter my username
    And I enter my password
    Then I should be logged in
    And I should see my character selection screen

  Scenario: Login with wrong password
    Given I have credentials username "player1" and wrong password
    When I enter my username
    And I enter incorrect password
    Then I should see a login failed message
    And I should be allowed to retry

  Scenario: Multiple characters per account
    Given I am logged in as "player1"
    And I have 3 characters: "Leo", "Donnie", "Raph"
    Then I should see a list of my characters
    And I can select which character to play

  Scenario: Password obfuscation
    Given I am at the password prompt
    When I type my password
    Then the password should show as asterisks
    And the real password should not be visible on screen

  Scenario: REST API authentication
    Given I have valid user credentials
    When I send a POST to /api/auth/login with credentials
    Then I should receive a session token
    And I can use the token for subsequent API calls