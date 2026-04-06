🔴 Feature: Status Effects (Buffs/Debuffs/DoT) - Issue #8
  As a combat system
  I need status effects that process each tick
  So that combat has strategic depth and lasting consequences

  Background:
    Given combat is in progress
    And status effect system is initialized
    And effect types are defined

  # DAMAGE OVER TIME
  Scenario: Bleeding deals damage each tick
    Given a target has BLEEDING effect for 3 ticks
    When a combat tick occurs
    Then the target takes 1 bleeding damage
    And bleed tick counter decreases by 1
    And effect expires after 3 ticks

  Scenario: Poison deals damage each tick
    Given a target has POISON effect for 5 ticks
    When a combat tick occurs
    Then the target takes 2 poison damage
    And poison tick counter decreases by 1
    And effect expires after 5 ticks

  Scenario: Burning deals damage and reduces accuracy
    Given a target has BURNING effect for 4 ticks
    When a combat tick occurs
    Then the target takes 1 burning damage
    And the target has -10% accuracy
    And burn tick counter decreases by 1

  Scenario: Multiple DoT effects stack
    Given a target has BLEEDING and POISON effects
    When a combat tick occurs
    Then the target takes bleeding damage
    And the target takes poison damage
    And both effects tick down independently

  # CONTROL EFFECTS
  Scenario: Stunned target cannot act
    Given a target has STUNNED effect for 2 ticks
    When the target's turn comes in combat
    Then the target should be skipped
    And no action should be performed
    And stun tick counter decreases by 1

  Scenario: Blinded target has reduced accuracy
    Given a target has BLINDED effect for 3 ticks
    When the target attempts an attack
    Then accuracy should be reduced by 50%
    And the attack may miss
    And blind tick counter decreases by 1

  # BUFF EFFECTS
  Scenario: Strength buff increases damage
    Given a target has BUFF_STRENGTH effect for 5 ticks
    When the target attacks
    Then damage should be increased by 25%
    And buff tick counter decreases by 1

  Scenario: Shield buff reduces incoming damage
    Given a target has BUFF_SHIELD effect for 3 ticks
    When the target receives an attack
    Then incoming damage is reduced by 25%
    And buff tick counter decreases by 1

  # EFFECT APPLICATION
  Scenario: Effect can be applied via ability
    Given a player uses "slash" with bleeding weapon
    When the attack hits an enemy
    Then the enemy should gain BLEEDING effect
    And the effect should have correct duration

  Scenario: Same effect refreshes duration
    Given a target already has BLEEDING for 2 ticks remaining
    When a new BLEEDING effect is applied
    Then the duration should be refreshed to full
    And tick counter should reset

  Scenario: Effect does not stack with itself
    Given a target has BLEEDING effect
    When another BLEEDING source is applied
    Then the effect should not duplicate
    And duration should refresh, not stack
