Feature: Character Authentication (Issue #9)
  As a player
  I want to authenticate via SSH
  So that I can securely access the game

  Background:
    Given the SSH server is running on port 4444
    And I have a registered user account

  Scenario: SSH connection prompts for credentials
    When I connect to the SSH server
    Then I should see a username prompt
    And after entering username, I should see a password prompt

  Scenario: Successful authentication
    Given I have valid credentials
    When I enter correct username and password
    Then I should be authenticated
    And I should see the game welcome screen

  Scenario: Failed authentication with wrong password
    Given I have a registered username "testuser"
    When I enter correct username but wrong password
    Then authentication should fail
    And I should see an error message
    And I should be prompted to try again

  Scenario: Failed authentication with unknown user
    Given I have no registered account
    When I enter unknown username "nobody"
    Then authentication should fail
    And I should see an error message

  Scenario: Password is obfuscated during entry
    Given I am at the password prompt
    When I type my password
    Then each character should appear as *
    And the actual password should not be visible

  Scenario: Admin user authentication
    Given I have an admin account with is_admin true
    When I authenticate successfully
    Then I should have admin privileges
    And I can access admin commands