Feature: Initiative System
  As a combat system
  I need to determine turn order in combat
  So that participants act in the correct sequence

  Background:
    Given a combat is starting

  Scenario: Initiative is calculated correctly
    Given player has 14 DEX
    When initiative is rolled
    Then initiative should include DEX modifier
    And initiative should include random 1-20 roll

  Scenario: Initiative formula applies DEX bonus
    Given player with 10 DEX
    When initiative is calculated
    Then base initiative should include 10 (from DEX)
    And random roll 1-20 is added

  Scenario: Turn order is sorted by initiative
    Given multiple participants in combat
    When initiative is rolled for all
    Then participants should be sorted descending by initiative

  Scenario: Higher initiative acts first
    Given participant A has initiative 18
    And participant B has initiative 14
    When tick resolves
    Then participant A should act before participant B

  Scenario: Ties broken by DEX
    Given participant A and B have same initiative roll
    And participant A has higher DEX
    When tie is broken
    Then participant A should act first

  Scenario: Turn order displayed in UI
    Given combat is active with multiple enemies
    When player looks at combat screen
    Then turn order should be visible
    And each participant's initiative should be shown

  Scenario: Initiative persists for combat duration
    Given combat has started
    And initiative was rolled
    When ticks pass
    Then initiative order should not change

  Scenario: Player can see initiative value
    Given combat is active
    When player checks status
    Then player's initiative value should be shown

  Scenario Outline: Initiative calculations for different DEX values
    Given participant has <dex> DEX
    When initiative is rolled
    Then initiative should be in range <min> to <max>

    Examples:
      | dex | min | max |
      | 8   | 17  | 36  |
      | 10  | 21  | 40  |
      | 14  | 29  | 48  |
      | 18  | 37  | 56  |