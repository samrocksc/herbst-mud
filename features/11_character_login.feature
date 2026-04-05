🔴 Feature: Character Login
  As a player
  I want to log in with my username and password
  So that I can play with my characters

  Background:
    Given the SSH server is running on port 4444
    And the database is connected

  Scenario: User login screen displays after SSH connection
    Given I have a user account already
    When I SSH to herbst-mud
    Then I should see a screen asking for my username
    And I should see the welcome banner

  Scenario: Successful login with valid credentials
    Given a user "player@example.com" exists with password "password123"
    When I SSH to herbst-mud
    And I enter my username "player@example.com"
    And I enter my password "password123"
    Then I should be logged in successfully
    And I should see my character selection screen

  Scenario: Login with multiple characters shows selection
    Given a user "player@example.com" exists with password "password123"
    And the user has 3 characters:
      | name | "Warrior" |
      | name | "Mage" |
      | name | "Thief" |
    When I log in successfully
    Then I should see a character selection menu
    And I should see all 3 characters listed
    When I select "Warrior"
    Then I should enter the game as "Warrior"

  Scenario: Login redirects to character creation if no characters exist
    Given a user "newplayer@example.com" exists with password "password123"
    And the user has no characters
    When I log in successfully
    Then I should be prompted to create a character
    And I should not see the game screen until character is created

  Scenario: Login with invalid password
    Given a user "player@example.com" exists with password "correctpassword"
    When I SSH to herbst-mud
    And I enter my username "player@example.com"
    And I enter my password "wrongpassword"
    Then I should see an error "Invalid credentials"
    And I should remain at the login prompt

  Scenario: Login with non-existent username
    When I SSH to herbst-mud
    And I enter my username "nonexistent@example.com"
    And I enter my password "anypassword"
    Then I should see an error "User not found"
    And I should be able to try again

  Scenario: Login case sensitivity for username
    Given a user "Player@example.com" exists with password "password123"
    When I SSH to herbst-mud
    And I enter my username "player@example.com"
    And I enter my password "password123"
    Then the behavior should match login system design (case-sensitive or normalize)

  Scenario: Password input is masked during login
    Given I SSH to herbst-mud
    When I enter my username "player"
    Then the password prompt should mask input
    And the password should not be visible on screen

  Scenario: Successful character selection enters the game
    Given a user "player@example.com" exists with password "password123"
    And the user has character "Warrior" in room "The Hole"
    When I log in and select "Warrior"
    Then I should see the room "The Hole"
    And my prompt should show "Warrior"

  Scenario: Login session persists until logout
    Given I am logged in as "player" with character "Warrior"
    When I stay active in the game
    Then my session should remain active
    And I should not be logged out unexpectedly

  Scenario: Logout returns to login screen
    Given I am logged in as "player" with character "Warrior"
    When I type the logout command
    Then I should see the login prompt
    And my character should no longer be in an active session