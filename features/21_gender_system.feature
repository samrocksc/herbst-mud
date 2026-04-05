🔴 Feature: Gender System - Implementation
  As a player
  I want to choose a gender/pronoun for my character
  So that my character can be properly referred to in-game

  Background:
    Given the database is connected
    And the gender system is implemented

  Scenario: Character can have he/him pronouns
    Given I create a character
    When I set the gender to "he/him"
    Then the character should have gender "he/him"
    And the character should be referred to with he/him pronouns
    And the character profile should display "he/him"

  Scenario: Character can have she/her pronouns
    Given I create a character
    When I set the gender to "she/her"
    Then the character should have gender "she/her"
    And the character should be referred to with she/her pronouns
    And the character profile should display "she/her"

  Scenario: Character can have it/its pronouns
    Given I create a character
    When I set the gender to "it/its"
    Then the character should have gender "it/its"
    And the character should be referred to with it/its pronouns
    And the character profile should display "it/its"

  Scenario: Character can have they/them pronouns
    Given I create a character
    When I set the gender to "they/them"
    Then the character should have gender "they/them"
    And the character should be referred to with they/them pronouns
    And the character profile should display "they/them"

  Scenario: Gender system is extensible
    Given I create a character
    When I set the gender to a custom value "xe/xem"
    Then the character should have gender "xe/xem"
    And the system should accept custom gender values
    And the character should be referred to with the specified pronouns

  Scenario: Gender can be edited via profile command
    Given I have a character with gender "he/him"
    When I use the "profile" command
    And I select "Edit Gender"
    And I enter "they/them"
    Then my character gender should be updated to "they/them"
    And the change should persist in the database

  Scenario: Gender displayed in character info
    Given I have a character with gender "she/her"
    When I use the "whoami" command
    Then I should see "Gender: she/her" in the output

  Scenario: Gender affects room descriptions
    Given a character named "Alice" with gender "she/her"
    And the character is in a room with other players
    When another player looks at the room
    Then they should see "Alice (she/her)" in the character list

  Scenario: Gender can be changed after character creation
    Given a character has gender "he/him"
    When I update the gender to "she/her"
    Then the character should have gender "she/her"
    And all references should use the new pronouns

  Scenario: Empty gender is rejected
    Given I am creating a new character
    When I do not select a gender
    Then I should see an error "Gender is required"
    Or a default gender should be assigned

  Scenario: Gender affects examination output
    Given a character with gender "they/them" attacks an enemy
    When another player examines the combat
    Then they should see "they" and "them" pronouns

  Scenario: Gender affects NPC dialogue if player is present
    Given a male NPC "Guard" is in a room
    And a character with gender "she/her" enters
    When the Guard speaks
    Then the dialogue should reference "her" correctly

  Scenario: Gender field accepts empty string for none
    Given I create a character
    When I set gender to an empty string ""
    Then the character should accept the empty gender
    Or should default to a neutral gender

  Scenario: Gender with special characters is handled
    Given I attempt to set gender with special characters "<he>"
    Then the system should sanitize or reject the input
    Or accept it if the system allows custom values