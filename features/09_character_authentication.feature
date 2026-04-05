🔴 Feature: Character Authentication via SSH
  As a player
  I want to log in through SSH with my credentials
  So that I can access my characters and play the game

  Background:
    Given the SSH server is running on port 4444
    And the database is connected

  Scenario: SSH connection shows welcome banner
    When I SSH to herbst-mud
    Then I should see a welcome banner
    And I should see turtle ASCII art

  Scenario: SSH login asks for username
    When I SSH to herbst-mud
    Then I should be prompted for my username
    And the prompt should be "Username:"

  Scenario: SSH login asks for password after username
    Given I have a user account with username "player"
    When I SSH to herbst-mud
    And I enter my username "player"
    Then I should be prompted for my password
    And the password input should be obfuscated

  Scenario: Successful login with valid credentials
    Given a user "player@example.com" exists with password "password123"
    When I SSH to herbst-mud
    And I enter my username "player@example.com"
    And I enter my password "password123"
    Then I should be logged in successfully
    And I should see my character selection screen

  Scenario: Failed login with invalid password
    Given a user "player@example.com" exists with password "correctpassword"
    When I SSH to herbst-mud
    And I enter my username "player@example.com"
    And I enter my password "wrongpassword"
    Then I should see an error "Invalid credentials"
    And I should be prompted to try again

  Scenario: Failed login with non-existent user
    When I SSH to herbst-mud
    And I enter my username "nonexistent@example.com"
    And I enter my password "anypassword"
    Then I should see an error "User not found"
    And I should be prompted to try again

  Scenario: Admin user can log in via SSH
    Given an admin user "admin@example.com" exists with password "adminpass" and is_admin: true
    When I SSH to herbst-mud
    And I enter my username "admin@example.com"
    And I enter my password "adminpass"
    Then I should be logged in as admin
    And I should have access to admin commands

  Scenario: User with no characters redirected to creation
    Given a user "newplayer@example.com" exists with password "password123"
    And the user has no characters
    When I log in successfully
    Then I should be redirected to character creation
    And I should not see the game world until I create a character

  Scenario: User with characters sees selection menu
    Given a user "player@example.com" exists with password "password123"
    And the user has characters:
      | name | "Warrior" |
      | name | "Mage" |
    When I log in successfully
    Then I should see a character selection menu
    And I should see "Warrior" in the list
    And I should see "Mage" in the list

  Scenario: Selecting character enters the game
    Given a user "player@example.com" exists with password "password123"
    And the user has character "Warrior" in room "The Hole"
    When I log in and select character "Warrior"
    Then I should enter the game
    And I should see the room description for "The Hole"

  Scenario: Account lockout after multiple failed login attempts
    Given a user "player@example.com" exists with password "correctpassword"
    When I SSH to herbst-mud
    And I enter my username "player@example.com"
    And I enter my password "wrongpassword" 5 times
    Then my account should be temporarily locked
    And I should see a "Too many failed attempts" message

  Scenario: SSH session timeout due to inactivity
    Given I am logged in as "player"
    When I am inactive for 30 minutes
    Then I should see a warning before disconnect
    And the session should be closed

  Scenario: SSH handles terminal resize during login
    When I SSH to herbst-mud
    And I resize my terminal
    Then the UI should adapt to the new terminal size
    And no display errors should occur

  Scenario: Logged in user can log out
    Given I am logged in as "player"
    When I type the logout command
    Then I should be logged out
    And I should see the login prompt again

  Scenario: Concurrent login from multiple sessions
    Given a user "player@example.com" exists with password "password123"
    When I log in from SSH session A
    And I attempt to log in from SSH session B with same credentials
    Then the behavior should match game design (allow/block/multiple chars)