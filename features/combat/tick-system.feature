🔴 Feature: Combat Tick System - Issue Combat-13 (Part of Combat System)
  As a combat system
  I need a tick system that processes combat round by round
  So that combat has rhythm and timing matters

  Background:
    Given combat has been initiated
    And all combatants are registered in combat

  # TICK RHYTHM
  Scenario: Combat advances in discrete ticks
    Given combat is in progress
    When each tick timer fires
    Then combat state should advance
    And each combatant should get a turn or partial turn

  Scenario: Ticks occur at consistent intervals
    Given combat tick speed is configured
    When combat is running
    Then ticks should occur at the configured interval
    And ticks should not bunch up or skip

  Scenario: Tick speed can be modified by abilities
    Given a combatant uses a speed-affecting ability
    When the tick system processes the ability
    Then the tick rate may increase or decrease
    And a message should indicate the change

  # TURN EXECUTION
  Scenario: Each combatant acts once per full round
    Given combat round starts
    When tick sequence completes a round
    Then each combatant should have one action opportunity
    And turn order should be determined by initiative/DEX

  Scenario: Faster combatants act first
    Given combatants have different DEX values
    When turn order is determined
    Then higher DEX combatants should act before lower DEX
    And ties should be broken consistently

  Scenario: Delayed action executes at correct time
    Given a combatant chooses to delay their action
    When the delayed turn arrives
    Then the combatant should act at the chosen moment
    And the action should resolve normally

  # TICK EVENTS
  Scenario: Status effects tick down each tick
    Given a combatant has a timed status effect
    When a tick occurs
    Then effect duration should decrease by 1
    And effects at 0 duration should expire

  Scenario: DoT damage applies each tick
    Given a combatant has damage-over-time effect
    When a tick occurs
    Then DoT damage should be applied
    And the effect should tick down

  Scenario: Buff effects refresh when reapplied
    Given a combatant has a buff
    When the same buff is reapplied before expiry
    Then the buff duration should refresh
    And buff effect continues without interruption

  # COMBAT TIMEOUT
  Scenario: Combat times out if too long
    Given combat has exceeded maximum tick limit
    When the timeout threshold is reached
    Then combat should end automatically
    And a draw or forced resolution should occur

  Scenario: Combat ends when one side is eliminated
    Given combat is in progress
    When all enemies or all players are defeated
    Then combat should end
    And victory/defeat should be determined

  Scenario: Combat can be fled from
    Given a combatant attempts to flee
    When the flee check succeeds
    Then combat should end
    And the fleeing combatant should exit combat
