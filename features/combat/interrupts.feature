🔴 Feature: Combat Interrupts (Parry, Shield Bash, Stun) - Issue #9
  As a combat system
  I need interrupt mechanics that cancel or redirect enemy actions
  So that timing and skill matter in combat

  Background:
    Given combat is in progress
    And interrupt abilities are defined

  # PARRY
  Scenario: Successful parry cancels incoming attack
    Given a defender has PARRY ability available
    And an enemy is attacking
    When the defender uses PARRY at the right moment
    Then the enemy's attack should be cancelled
    And the defender should not take damage
    And the defender may gain a counter-attack opportunity

  Scenario: Failed parry leaves defender vulnerable
    Given a defender attempts to PARRY
    But the timing is incorrect
    When the enemy's attack lands
    Then the defender takes full damage
    And parry ability may enter cooldown

  Scenario: Parry has cooldown
    Given a defender uses PARRY
    When the parry is successful
    Then PARRY should enter cooldown
    And the defender cannot parry again until cooldown ends

  # SHIELD BASH
  Scenario: Shield bash stuns enemy
    Given a defender has SHIELD_BASH ability
    And an enemy is in melee range
    When the defender uses SHIELD_BASH
    Then the enemy should be stunned for 1 tick
    And the enemy should take shield bash damage
    And the enemy's action for that tick is cancelled

  Scenario: Shield bash interrupts channeled ability
    Given an enemy is channeling an ability
    And the defender uses SHIELD_BASH
    When the shield bash lands
    Then the enemy's channel should be interrupted
    And the channeled ability should not complete
    And any partial progress should be lost

  # STUN INTERRUPT
  Scenario: Stun prevents enemy action
    Given an enemy is STUNNED
    When the enemy's turn in combat tick arrives
    Then the enemy should take no action
    And the enemy should be skipped in turn order
    And stun duration should decrease by 1

  Scenario: Stun clears on action tick
    Given an enemy has STUNNED for 1 tick remaining
    When the enemy's turn is skipped
    Then the stun should be cleared
    And the enemy can act normally on next tick

  # COUNTER-ATTACK
  Scenario: Successful interrupt grants counter-attack
    Given a defender successfully PARRYs an attack
    When the parry is resolved
    Then the defender should have a counter-attack window
    And the defender can use a bonus attack
    And the bonus attack does not consume the defender's main action

  # COMBAT LOG
  Scenario: Interrupt is logged in combat output
    Given an interrupt occurs in combat
    When the interrupt is resolved
    Then the combat log should show the interrupt
    And the log should describe what was interrupted
