Feature: Character Authentication
  As a player
  I want to authenticate my character securely
  So that I can access the game without exposing my credentials

  Background:
    Given the authentication API is available
    And a character "HeroTest" exists with password "heroPass456"

  Scenario: Authenticate with valid credentials
    When I send an authenticate request with username "HeroTest" and password "heroPass456"
    Then the response status should be 200 OK
    And I should receive a valid authentication token
    And the token should not be empty
    And the token should be a non-empty string

  Scenario: Authenticate with invalid password
    When I send an authenticate request with username "HeroTest" and password "wrongpassword"
    Then the response status should be 401 Unauthorized
    And I should not receive an authentication token
    And the error message should indicate invalid credentials

  Scenario: Authenticate with non-existent character
    When I send an authenticate request with username "NobodyCharacter" and password "anypassword"
    Then the response status should be 404 Not Found
    And the error message should contain "character not found"

  Scenario: Authenticate with empty username
    When I send an authenticate request with username "" and password "anypassword"
    Then the response status should be 400 Bad Request

  Scenario: Authenticate with empty password
    When I send an authenticate request with username "HeroTest" and password ""
    Then the response status should be 400 Bad Request

  Scenario: Token is valid for subsequent requests
    When I authenticate successfully as "HeroTest" with password "heroPass456"
    And I use the received token to access a protected endpoint like /api/character/me
    Then the response status should be 200 OK
    And I should receive my character data

  Scenario: Invalid token is rejected
    When I use an invalid token "bad.token.here" to access /api/character/me
    Then the response status should be 401 Unauthorized

  Scenario: Expired token is rejected
    When I use an expired token to access /api/character/me
    Then the response status should be 401 Unauthorized
    And the error should mention token expiration

  Scenario: Authenticate with SQL injection attempt in username
    When I send an authenticate request with username "' OR '1'='1" and password "anything"
    Then the response status should be 400 or 401
    And no characters should be returned improperly

  Scenario: Session persists after authentication
    When I authenticate successfully as "HeroTest" with password "heroPass456"
    And I make another API call within the session timeout
    Then my session should remain active
    And I should not need to re-authenticate
