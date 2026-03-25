Feature: Character Authentication
  As a player
  I want to log in during SSH connection initiation
  So that I can securely access my character account

  Background:
    Given the MUD server is running and accepting SSH connections
    And I have an existing user account with credentials

  Scenario: Login prompt appears on SSH connection
    Given I initiate an SSH connection to the MUD server
    When the connection is established
    Then I should see the welcome ASCII art
    And I should be presented with "Log In" or "Create Account" options

  Scenario: Successful login with valid credentials
    Given I am at the SSH login prompt
    When I enter my username "testplayer"
    And I enter my password "validpassword"
    Then I should be authenticated successfully
    And I should see my character selection screen
    And I should see my existing characters listed

  Scenario: Login with invalid password shows error
    Given I am at the SSH login prompt
    When I enter my username "testplayer"
    And I enter an incorrect password "wrongpassword"
    Then I should see an authentication error message
    And I should be returned to the login prompt
    And I should be able to retry my login

  Scenario: Login with non-existent user
    Given I am at the SSH login prompt
    When I enter username "ghostuser"
    And I enter password "anypassword"
    Then I should see a "user not found" message
    And I should be returned to the login prompt

  Scenario: Password is obfuscated during input
    Given I am at the password prompt
    When I type my password
    Then the password characters should be masked as "****"
    And the actual password should not be visible on screen

  Scenario: Admin user gains admin access on login
    Given I have an admin user "gamemaster" with password "adminpass"
    When I log in with username "gamemaster" and password "adminpass"
    Then I should be authenticated as admin
    And I should have access to admin commands

  Scenario: User with no characters sees character creation
    Given I am a logged-in user with no characters
    When I reach the character selection screen
    Then I should be prompted to create a new character
    And I should see the fountain awakening flow

  Scenario: User with characters sees character list
    Given I am a logged-in user with existing characters
    When I reach the character selection screen
    Then I should see a list of my characters
    And I should be able to select which character to play

  Scenario: Username is separate from character name
    Given I have a user account "adminuser" with character "MyWizard"
    When I log in with username "adminuser"
    Then the username "adminuser" should be different from the character name "MyWizard"

  Scenario: Multiple characters displayed on successful login
    Given I have a user account with 3 characters: "Warrior", "Mage", "Rogue"
    When I log in successfully
    Then I should see all 3 characters in the selection list
    And I should be able to choose which character to play

  Scenario: Login screen uses full terminal width
    Given I initiate an SSH connection
    When the login screen is displayed
    Then it should use the full terminal width
    And the I/O style should match the character creation screen
