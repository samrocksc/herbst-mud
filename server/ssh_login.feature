Feature: SSH Login Authentication
  As an admin user
  I want to log into the SSH server
  So that I can access the game server

  Scenario: Admin user can log into SSH server
    Given the admin user exists with credentials "admin@herbstmud.local" / "herb5t2026!"
    When I attempt to authenticate via SSH
    Then I should be granted access to the game
    And I should be placed in the starting room