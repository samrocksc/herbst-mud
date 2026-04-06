Feature: Combat Tick System
  As a combat system
  I need a reliable tick clock that drives all combat timing
  So that actions resolve in the correct order and at the right intervals

  Background:
    Given the game is running in combat mode

  Scenario: Tick loop runs at default interval
    Given combat session is active
    When 1500 milliseconds pass
    Then a new tick should occur
    And tick counter should increment by 1

  Scenario: Tick loop can be configured
    Given tick interval is set to 1000ms (fast mode)
    When combat starts
    Then each tick should occur after 1000ms

  Scenario: Tick loop can run in slow mode
    Given tick interval is set to 2000ms (slow mode)
    When combat starts
    Then each tick should occur after 2000ms

  Scenario: Combat manager tracks active combats
    Given multiple combats are in progress
    When a new combat starts
    Then combat manager should assign unique ID
    And combat should be added to active combats list

  Scenario: Tick counter increments each cycle
    Given combat has been running for 3 ticks
    When tick resolution completes
    Then tick counter should equal 3

  Scenario: Tick loop can be stopped cleanly
    Given combat tick loop is running
    When stop signal is sent
    Then tick loop should stop
    And no more ticks should occur

  Scenario: Multiple combats tick independently
    Given combat A and combat B are both active
    When tick occurs in combat A
    Then combat B should continue independently
    And combat A should not affect combat B's state

  Scenario Outline: Tick interval configuration
    Given tick interval is set to <interval>ms
    When combat starts
    Then each tick should complete in approximately <interval>ms

    Examples:
      | interval |
      | 1000     |
      | 1500     |
      | 2000     |

  Scenario: Tick is suspended during interrupt
    Given enemy is stunned for 1 tick
    When enemy is stunned
    Then tick processing is suspended for that enemy
    And other combatants continue ticking normally

  Scenario: Tick resumes after interrupt ends
    Given enemy was stunned and tick was skipped
    When stun effect expires
    Then enemy is included in the next tick
    And enemy can decide and execute an action

  Scenario: Tick timing does not drift over many cycles
    Given combat has been running for 50 ticks
    When timing is measured across all ticks
    Then each tick fires within 5% of the configured interval
    And cumulative drift is minimal

  Scenario: Tick resolution order is correct
    Given combat includes player and enemy
    When tick occurs
    Then player acts first
    And enemy acts second
    And status effects apply last

  Scenario: Tick interacts with initiative system
    Given player has higher initiative than enemy
    When combat tick occurs
    Then player action resolves before enemy action
    And initiative does not reset mid-tick

  Scenario: Combat expires after maximum ticks
    Given maximum combat duration is set to 100 ticks
    When combat reaches 100 ticks without resolution
    Then combat is force-ended
    And warning message is shown to player

  Scenario: Concurrent combat ticks are isolated
    Given player is in combat A
    And player is in combat B
    When tick occurs in combat A
    Then combat B tick counter is unaffected
    And combat A state does not bleed into combat B

  Scenario: Tick pause and resume mechanics
    Given combat tick loop is running
    When pause command is issued
    Then tick loop pauses
    And no ticks occur
    When resume command is issued
    Then tick loop resumes from where it stopped
    And tick counter continues incrementing

  Scenario: Fast-forward tick processing
    Given combat is paused
    And 10 ticks have passed while paused
    When fast-forward is requested for 5 ticks
    Then 5 ticks are processed immediately
    And status effects apply 5 times
    And no output is shown per tick

  Scenario: Tick with no combatants still fires
    Given all combatants have been defeated
    And combat manager still tracks the combat
    When tick fires
    Then tick counter increments
    And combat is flagged for cleanup
    And no actions are processed