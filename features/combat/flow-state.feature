Feature: Combat Flow and State Machine
  As a combat system
  I need to manage combat states and transitions
  So that combat flows correctly from start to finish

  Background:
    Given the game is running

  Scenario: Combat starts when player attacks creature
    Given player is in exploration mode
    And room contains enemy "Scrap Rat"
    When player attacks the enemy
    Then combat mode should activate
    And combat state should be "COMBAT_INIT"

  Scenario: Combat starts when creature aggroes
    Given player is in exploration mode
    And enemy "Junk Dog" has aggro radius of 5
    When player moves within 5 tiles of enemy
    Then combat should start automatically
    And player should see "Combat started!"

  Scenario: Initiative is rolled at combat start
    Given combat has just started
    When initiative is calculated
    Then all participants should have initiative values
    And participants should be sorted by initiative

  Scenario: Combat enters active state after initiative
    Given combat state is "COMBAT_INIT"
    When initiative is resolved
    Then combat state should become "COMBAT_ACTIVE"

  Scenario: Combat tick follows correct resolution flow
    Given combat is in active state
    When a tick occurs
    Then input phase should process player actions
    And action queue should process in initiative order
    And status effects should process
    And enemy AI should decide actions
    And status checks should occur
    And tick counter should increment

  Scenario: Combat ends when all enemies are defeated
    Given combat has enemies remaining
    When final enemy HP reaches 0
    Then combat state should become "COMBAT_END"
    And victory screen should display

  Scenario: Combat ends when player flees successfully
    Given player attempts to flee
    And flee roll succeeds
    Then combat state should become "COMBAT_END"
    And player should return to exploration

  Scenario: Combat ends when player is defeated
    Given player HP reaches 0
    Then combat state should become "COMBAT_END"
    And defeat screen should display

  Scenario: Returning to exploration after combat ends
    Given combat state is "COMBAT_END"
    When loot is distributed
    And XP is awarded
    Then state should transition to "IDLE"
    And normal exploration should resume

  Scenario: Player has no input results in auto-attack
    Given combat is active
    And it is player's turn
    When player provides no input within tick window
    Then player should auto-attack

  Scenario: Player auto-defends at low HP
    Given player HP is below 25%
    And player provides no input within tick window
    Then player should auto-defend
    And defense bonus should apply

  Scenario: Late input is queued for next tick
    Given player presses key after tick boundary
    When next tick occurs
    Then input should be processed
    And no input should be lost

  Scenario: DoT effects process each tick
    Given player has bleeding effect active
    When tick occurs
    Then bleeding should deal damage
    And effect duration should decrease

  Scenario: Buff effects apply modifiers
    Given player has strength buff active
    When player attacks
    Then damage should be increased by buff amount

  Scenario: Debuff effects apply penalties
    Given player has blinded effect
    When player attacks
    Then accuracy should be reduced by 50%

  Scenario: Effects expire after duration
    Given effect has 3 tick duration
    When 3 ticks pass
    Then effect should be removed

  Scenario Outline: Combat state transitions correctly
    Given current state is "<from_state>"
    When "<trigger>" occurs
    Then state should become "<to_state>"

    Examples:
      | from_state   | trigger                  | to_state      |
      | IDLE         | player attacks           | COMBAT_INIT   |
      | COMBAT_INIT  | initiative resolved      | COMBAT_ACTIVE |
      | COMBAT_ACTIVE| tick resolution         | COMBAT_ACTIVE |
      | COMBAT_ACTIVE| victory condition       | COMBAT_END    |
      | COMBAT_ACTIVE| defeat condition        | COMBAT_END    |
      | COMBAT_END   | loot distributed        | IDLE          |