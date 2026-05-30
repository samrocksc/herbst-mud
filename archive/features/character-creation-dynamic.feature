Feature: Dynamic Character Creation
  As a player
  I want to create a character with configurable options from the admin UI
  So that I can only choose valid races and classes

  Background:
    Given I am logged in as a valid user
    And I am on the character selection screen
    And I have selected a world

  Scenario: Player starts character creation
    Given I am on the character selection screen
    When I type "n"
    Then character creation should start
    And I should be prompted for a character name

  Scenario: Player enters a valid character name
    Given I am creating a character
    And I am prompted for a character name
    When I type "HeroName"
    Then I should be prompted to select a race
    And I should see a numbered list of available races

  Scenario: Player selects race from dynamic list by number
    Given I am creating a character
    And the available races include "Human" and "Elf" and "Dwarf"
    When I type "2"
    Then my race should be "Elf"

  Scenario: Player selects race from dynamic list by name
    Given I am creating a character
    And the available races include "Human" and "Elf" and "Dwarf"
    When I type "Dwarf"
    Then my race should be "Dwarf"

  Scenario: Player selects invalid race number
    Given I am creating a character
    And there are 3 available races
    When I type "5"
    Then I should see an error message
    And I should be prompted to select a valid race

  Scenario: Player selects invalid race name
    Given I am creating a character
    And the available races include "Human" and "Elf" and "Dwarf"
    When I type "Orc"
    Then I should see an error message
    And I should be prompted to select a valid race

  Scenario: Player accepts default race
    Given I am creating a character
    And the available races include "Human" and "Elf" and "Dwarf"
    When I press enter without typing
    Then my race should be the first available race

  Scenario: Player selects class
    Given I am creating a character
    And I have selected a race
    When I am prompted to select a class
    Then I should see a list of valid classes
    And "adventurer" should be an option
    And "warrior" should be an option
    And "mage" should be an option
    And "rogue" should be an option
    And "cleric" should be an option

  Scenario: Player accepts default class
    Given I am creating a character
    And I have selected a race
    When I press enter without typing for class
    Then my class should be "adventurer"

  Scenario: Player selects a class
    Given I am creating a character
    And I have selected a race
    When I type "mage" for class
    Then my class should be "mage"

  Scenario: Player selects invalid class
    Given I am creating a character
    And I have selected a race
    When I type "ninja" for class
    Then I should see an error message
    And I should be prompted to select a valid class

  Scenario: Player cancels character creation
    Given I am creating a character
    When I type "cancel"
    Then I should return to the character selection screen

  Scenario: Character name validation rejects numbers
    Given I am creating a character
    When I type "Hero123" as name
    Then I should see an error message
    And I should be prompted for a valid name

  Scenario: Character name validation rejects empty
    Given I am creating a character
    When I type "" as name
    Then I should see an error message

  Scenario: Character name validation rejects too long
    Given I am creating a character
    When I type "ThisNameIsWayTooLong" as name
    Then I should see an error message

  Scenario: Dynamic race list updates from API
    Given the API returns 3 playable races
    Then the race selection should show exactly 3 options
    Given the API returns 5 playable races
    Then the race selection should show exactly 5 options
