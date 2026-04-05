🔴 Feature: Race System - Implementation
  As a player
  I want to choose a race for my character
  So that I can have unique racial traits and abilities

  Background:
    Given the database is connected
    And the race system is implemented

  Scenario: Character can have Mutant race
    Given I create a character
    When I set the race to "Mutant"
    Then the character should be a Mutant
    And the character can have a mix (turtle, rat, rhino, etc.)
    And the character should have Mutant-specific traits

  Scenario: Character can have Human race
    Given I create a character
    When I set the race to "Human"
    Then the character should be Human
    And the character should have Human-specific traits
    And the character should have balanced stats

  Scenario: Character can have Animal race
    Given I create a character
    When I set the race to "Animal"
    Then the character should be an Animal
    And the character should have a species field
    And the species should be specified

  Scenario: Animal race requires species
    Given I create a character with race "Animal"
    When I do not specify a species
    Then I should see an error "Species is required for Animal race"
    When I specify species "Turtle"
    Then the character should be an Animal with species "Turtle"

  Scenario: Race selection during character creation
    Given I am creating a new character
    When I reach the race selection screen
    Then I should see all available races
    And I should be able to select one race
    And the race should be saved to my character

  Scenario: Race affects starting stats
    Given I create a Mutant character
    Then the character should have Mutant stat modifiers
    When I create a Human character
    Then the character should have Human stat modifiers
    When I create an Animal character
    Then the character should have Animal stat modifiers

  Scenario: Race affects available classes
    Given I create a Mutant character
    Then the character should be able to choose any class
    And certain race/class combinations may have bonuses

  Scenario: Race cannot be changed after character creation
    Given a character "Hero" has race "Human"
    When I attempt to change the race to "Mutant"
    Then the race change should be denied
    Or the character should be able to change at special locations

  Scenario: Invalid race is rejected
    Given I am creating a new character
    When I select an invalid race "Alien"
    Then I should see an error "Invalid race selection"
    And the character should not be created

  Scenario: Race is required for character
    Given I am creating a new character
    When I proceed without selecting a race
    Then I should see an error "Race selection is required"
    And the character should not be created

  Scenario: Mutant race can have specific mutations
    Given a character has race "Mutant"
    When I specify mutation "turtle"
    Then the character should have the "turtle" mutation
    And the mutation should affect gameplay

  Scenario: Animal species affects character name
    Given a character has race "Animal" and species "Wolf"
    When other players see the character
    Then the character should be identified as a "Wolf" creature

  Scenario: Race affects starting HP
    Given I create a Mutant character
    Then the character should have Mutant-specific base HP
    When I create a Human character
    Then the character should have Human base HP

  Scenario: Race affects base stats differently
    Given I create characters of each race with same class
    Then Mutants should have different stat base values
    And Humans should have balanced stats
    And Animals should have species-specific stats

  Scenario: Race affects physical appearance in look command
    Given I look at a Mutant character
    Then the character description should mention their mutant nature
    Given I look at an Animal character
    Then the character description should mention their species