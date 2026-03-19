Feature: Combat Interrupts (Parry, Shield Bash)
  As a player
  I want interrupt abilities that counter enemy attacks
  So that I can react strategically during combat

  Background:
    Given a combat is active
    And the player has parry and shield_bash available
    And the enemy is channeling an ability

  Scenario: Parry negates same-tick attack
    Given the enemy starts a basic attack
    And the attack costs 1 tick
    When the player uses parry on the same tick
    Then the enemy's attack damage is negated
    And the player counter-attacks
    And the enemy's attack action is cancelled

  Scenario: Parry requires stamina
    Given the player has 0 stamina remaining
    When the player attempts to parry
    Then the parry fails
    And the player's stamina is unchanged
    And the enemy attack proceeds normally

  Scenario: Shield bash cancels enemy channeling
    Given the enemy is channeling "Crush" (2 ticks)
    And the cast is on tick 1
    When the player uses shield_bash
    Then the enemy's channel is cancelled
    And the enemy takes bash damage
    And the enemy is stunned for 2 ticks

  Scenario: Shield bash interrupts at any channel stage
    Given the enemy is channeling a 3-tick ability
    When the player uses shield_bash during any tick
    Then the channel is cancelled
    And the ability does not resolve

  Scenario: Stun interrupts enemy action
    Given the enemy is channeling an ability
    When the enemy receives a stun effect
    Then the channeling is cancelled
    And the enemy's cooldown is reset
    And the enemy cannot act during stun ticks

  Scenario: Instant actions cannot be parried
    Given the enemy uses an instant ability (0 tick cost)
    When the player attempts to parry
    Then the parry does not negate the damage
    And the instant ability resolves normally

  Scenario: Parry shows feedback in combat log
    Given the player parries an enemy attack
    When the parry resolves
    Then the combat log shows "[Player] parries [Enemy]'s attack!"
    And the counter-attack is logged

  Scenario: Shield bash shows feedback in combat log
    Given the player shield bashes a channeling enemy
    When the interrupt resolves
    Then the combat log shows "[Player] interrupts [Enemy]'s action!"
    And the stun application is logged

  Scenario: Interrupt timing validation
    Given the enemy is channeling an ability
    When the player uses an interrupt on the wrong tick
    Then the interrupt is rejected
    And the player receives feedback about timing

  Scenario: Multiple interrupts on same tick
    Given multiple enemies are channeling abilities
    When the player uses parry on one enemy
    And shield_bash on another enemy
    Then both interrupts resolve correctly
    And each enemy's action is handled independently
