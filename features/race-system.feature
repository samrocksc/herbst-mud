Feature: Race System - Implementation (Issue #20)
  As a player
  I want to select a race for my character
  So that my character has unique traits and abilities

  Background:
    Given the race system is implemented in the game

  Scenario: View available races
    Given I am creating a new character
    When I reach the race selection screen
    Then I should see all available races
    And each race should have a description

  Scenario: Select Mutant race
    Given I am on the race selection screen
    When I select "Mutant" race
    Then my character should be a Mutant
    And I should be able to specify a mutation type (turtle, rat, rhino)

  Scenario: Mutant with mixed mutations
    Given I selected the Mutant race
    When I choose "turtle" as my mutation
    Then my character should have turtle-like attributes
    And stats should reflect the turtle mutation

  Scenario: Select Human race
    Given I am on the race selection screen
    When I select "Human" race
    Then my character should be Human
    And I should have standard human stats

  Scenario: Select Animal race
    Given I am on the race selection screen
    When I select "Animal" race
    Then my character should be an Animal type
    And I should be able to specify the species

  Scenario: Animal species specification
    Given I selected Animal race
    When I specify species "wolf"
    Then my character should be a wolf animal
    And stats should reflect the wolf species

  Scenario: Race affects stats
    Given I am creating a character
    When I select different races
    Then each race should have different base stats
    And the stats should be visible during character creation

  Scenario: Race affects appearance
    Given I have created a character with a specific race
    When other players look at my character
    Then they should see visual indicators of my race
    And the character description should reflect the race