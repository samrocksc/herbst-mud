Feature: Gender System
  As a game developer
  I want a Gender system
  So that players can identify their characters

  Scenario: Available genders
    When I list all available genders
    Then I should see: Male, Female, Non-binary, Other

  Scenario: Gender selection
    Given I am creating a character
    When I select gender "Female"
    Then my character should have gender "Female"

  Scenario: Gender display
    Given a character with gender "Non-binary"
    When I look at the character
    Then I should see the gender displayed correctly