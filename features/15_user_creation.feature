Feature: User Creation
  As a new player
  I want to create a user account
  So that I can register to play the game

  Scenario: Create new user account
    When I create a user account with username "newplayer" and email "player@test.com"
    Then the account should be created
    And I should receive a confirmation

  Scenario: Create user with duplicate username
    Given a user "taken" already exists
    When I try to create a user with username "taken"
    Then creation should fail
    And I should see "username already exists"

  Scenario: Create user with invalid email
    When I create a user with email "notanemail"
    Then creation should fail
    And I should see "invalid email format"