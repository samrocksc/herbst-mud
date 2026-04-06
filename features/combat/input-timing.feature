Feature: Combat Input Timing
  As a player
  I need a clear input window during combat
  So that I can choose my actions before the tick resolves

  Background:
    Given a combat session is active
    And it is the player's turn

  Scenario: Input window timer is displayed
    Given combat tick has started
    When player looks at combat screen
    Then input timer countdown should be visible

  Scenario: Input window defaults to 1.5 seconds
    Given combat is active
    When tick starts
    Then player should have 1.5 seconds to input

  Scenario: Key 1 triggers first talent
    Given player has talent "slash" in slot 1
    When player presses "1"
    Then slash action should be queued

  Scenario: Key 2 triggers second talent
    Given player has talent "parry" in slot 2
    When player presses "2"
    Then parry action should be queued

  Scenario: Key 3 triggers third talent
    Given player has talent "heavy_strike" in slot 3
    When player presses "3"
    Then heavy_strike action should start

  Scenario: Key 4 triggers fourth talent
    Given player has talent "second_wind" in slot 4
    When player presses "4"
    Then second_wind should start channeling

  Scenario: No input results in auto-attack
    Given player provides no input
    When input window expires
    Then basic attack should be executed automatically

  Scenario: Auto-defend triggers at low HP
    Given player HP is below 25%
    And player provides no input
    When input window expires
    Then defend action should be auto-selected

  Scenario: Late input is queued for next tick
    Given player presses key after tick boundary
    When next tick starts
    Then the key press should be processed
    And input should not be lost

  Scenario: Visual countdown shows time remaining
    Given combat is active
    When player looks at input window
    Then remaining time should be shown
    And time should count down

  Scenario: Input works at various terminal sizes
    Given terminal is 80x24
    When combat starts
    Then input should work correctly

  Scenario: Input works at minimum terminal size
    Given terminal is minimum supported size
    When combat starts
    Then input functionality should be preserved

  Scenario: Key presses are responsive
    Given player presses key during input window
    When key is received
    Then action should trigger immediately

  Scenario Outline: Input timer respects configured duration
    Given input window is set to <duration> seconds
    When tick starts
    Then player should have <duration> seconds to respond

    Examples:
      | duration |
      | 1.0      |
      | 1.5      |
      | 2.0      |