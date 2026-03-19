Feature: User Creation
  As a new player
  I want to create a user account
  So that I can access the game

  Scenario: Create user with valid data
    Given I am on the registration page
    When I submit:
      | field    | value           |
      | username | newplayer       |
      | email    | new@test.com    |
      | password | securePass123   |
    Then the user account should be created
    And I should receive a confirmation email

  Scenario: Create user with existing username
    Given a user "taken" already exists
    When I try to register with username "taken"
    Then I should receive a "username taken" error

  Scenario: Create user with invalid email
    Given I am on the registration page
    When I submit an invalid email "notanemail"
    Then I should receive an "invalid email" error