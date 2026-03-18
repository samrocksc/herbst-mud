Feature: Character Login
  As a player
  I want to log into the game
  So that I can start playing

  Scenario: Successful login
    Given a character "Hero1" exists with password "pass123"
    When I log in with character "Hero1" and password "pass123"
    Then login should be successful
    And I should be in the game world

  Scenario: Login with wrong password
    Given a character "Mage1" exists with password "correct"
    When I log in with character "Mage1" and password "wrong"
    Then login should fail
    And I should see an error message