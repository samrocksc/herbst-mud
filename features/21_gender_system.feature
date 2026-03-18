Feature: Gender System
  As a player
  I want to choose a gender for my character
  So that I can personalize my character

  Scenario: List available genders
    When I request available genders
    Then I should see: Male, Female, Non-binary, None

  Scenario: Select male gender
    Given I choose gender "Male"
    Then my character should be addressed with male pronouns

  Scenario: Select female gender
    Given I choose gender "Female"
    Then my character should be addressed with female pronouns

  Scenario: Select non-binary gender
    Given I choose gender "Non-binary"
    Then my character should be addressed with neutral pronouns

  Scenario: Select no gender
    Given I choose gender "None"
    Then my character should not have gendered references