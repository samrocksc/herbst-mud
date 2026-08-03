Feature: Root User Authentication
  As a system administrator
  I want to verify that the root account can authenticate with the correct password
  So that I can ensure administrative access is preserved

  Scenario: Successful root login
    Given the database contains a root user "admin@herbstmud. laocal" with password "herb5t2026!"
    When I attempt to authenticate via POST "/users/auth" with email "admin@herbstmud.local" and password "herb5t2026!"
    Then the response status should be 200
    And the response should contain a valid JWT token
    And the response should indicate "is_admin": true

  Scenario: Failed root login with wrong password
    Given the database contains a root user "admin@herbstmud.local"
    When I attempt to authenticate via POST "/users/auth" with email "admin@herbstmud.local" and password "wrong-password"
    Then the response status should be 401
    And the response error should be "Invalid email or password"
