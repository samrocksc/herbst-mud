Feature: User Login via SSH
  As a player
  I want to log into the MUD via SSH
  So that I can access the game and play

  Background:
    Given the SSH server is running on port 4444
    And the REST API is available

  Scenario: Player sees welcome screen on connect
    Given I am connected to the SSH server
    Then I should see the welcome screen
    And I should see login option
    And I should see register option
    And I should see quit option

  Scenario: Player navigates to login screen
    Given I am on the welcome screen
    When I type "1"
    Then I should be on the login screen
    And I should be prompted for my username

  Scenario: Player enters username
    Given I am on the login screen
    And I am prompted for my username
    When I type "player@example.com"
    Then I should be prompted for my password

  Scenario: Player logs in with valid credentials
    Given I am on the login screen
    And I have entered my username "player@example.com"
    When I enter my password "password123"
    Then the login should succeed
    And I should see the world selection screen
    And I should see a welcome back message

  Scenario: Player logs in with invalid credentials
    Given I am on the login screen
    And I have entered my username "player@example.com"
    When I enter my password "wrongpassword"
    Then the login should fail
    And I should see an error message
    And I should be prompted to try again

  Scenario: Player registers a new account
    Given I am on the welcome screen
    When I type "2"
    Then I should be on the register screen
    And I should be prompted for a username

  Scenario: Player can escape back to welcome screen
    Given I am on the login screen
    When I press the escape key
    Then I should be back on the welcome screen

  Scenario: Player can type 'login' instead of '1'
    Given I am on the welcome screen
    When I type "login"
    Then I should be on the login screen

  Scenario: Player can type 'register' instead of '2'
    Given I am on the welcome screen
    When I type "register"
    Then I should be on the register screen
