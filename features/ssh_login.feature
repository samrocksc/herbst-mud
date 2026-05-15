Feature: SSH Login Authentication
  As an admin user
  I want to log into the SSH server
  So that I can access the game server

  Scenario: Admin user can log into SSH server
    Given the admin user exists with credentials "admin@herbstmud.local" / "herb5t2026!"
    When I attempt to authenticate via SSH
    Then I should be granted access to the game
    And I should be placed in the starting room

  Scenario: Non-admin user cannot log in
    Given the admin user exists with credentials "admin@herbstmud.local" / "herb5t2026!"
    And a regular user "player@example.com" exists
    When the regular user attempts to authenticate via SSH
    Then I should receive a 403 Forbidden
    And the error message should indicate admin access required

  Scenario: Failed SSH authentication shows error
    Given the admin user exists with credentials "admin@herbstmud.local" / "herb5t2026!"
    When I attempt to authenticate with wrong password "wrongpass"
    Then I should receive a 401 Unauthorized
    And the error message should indicate invalid credentials
