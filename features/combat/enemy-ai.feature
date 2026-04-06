🔴 Feature: Enemy AI Decision Making - Issue #7
  As a combat system
  I need enemies to make intelligent decisions each tick
  So that combat feels dynamic and challenging

  Background:
    Given combat is in progress with an enemy
    And the enemy has HP, abilities, and a decision tree

  # HEALTH-BASED DECISIONS
  Scenario: Enemy flees when HP is critically low
    Given enemy HP is below 25%
    And enemy has not yet attempted to flee
    When the enemy AI evaluates its next action
    Then there should be a flee chance based on enemy DEX
    And if flee succeeds, enemy exits combat
    And if flee fails, enemy continues fighting

  Scenario: Enemy heals when HP is low and heal ability available
    Given enemy HP is below 25%
    And enemy has a heal ability available
    When the enemy AI evaluates its next action
    Then the enemy should consider using heal
    And heal should restore a portion of max HP

  Scenario: Enemy attacks normally when HP is healthy
    Given enemy HP is above 50%
    When the enemy AI evaluates its next action
    Then the enemy should choose a basic attack
    And the attack should target the player

  # ABILITY USAGE
  Scenario: Enemy uses special ability based on situation
    Given an enemy has multiple abilities
    And the enemy is in combat
    When the enemy AI evaluates its next action
    Then the enemy should select an appropriate ability
    And ability cooldowns should be respected

  Scenario: Enemy does not use ability on cooldown
    Given an enemy has an ability on cooldown
    When the enemy AI evaluates abilities
    Then the ability should not be selected
    And available abilities should be prioritized

  # TARGETING
  Scenario: Enemy targets player character
    Given an enemy is in combat with a player
    When the enemy selects a target
    Then the target should be the player character
    And enemy attacks should deal damage to the player

  Scenario: Enemy with multi-target ability chooses AoE
    Given an enemy has an AoE ability
    And there are multiple targets in range
    When the enemy AI evaluates its next action
    Then the enemy may use AoE ability
    And all targets in range should be affected

  # INTELLIGENCE TIERS
  Scenario: High-intelligence enemy uses tactical abilities
    Given an enemy with high intelligence tier
    When the enemy AI evaluates combat
    Then the enemy should use abilities strategically
    And the enemy should prioritize high-value targets

  Scenario: Low-intelligence enemy uses basic attacks
    Given an enemy with low intelligence tier
    When the enemy AI evaluates combat
    Then the enemy should prefer basic attacks
    And special abilities should be used rarely

  # TICK-BASED DECISIONS
  Scenario: Enemy makes one decision per combat tick
    Given combat tick occurs
    When the enemy AI processes the tick
    Then exactly one action should be selected
    And the action should be executed

  Scenario: Enemy AI runs each tick until combat ends
    Given combat is ongoing
    When combat ticks progress
    Then the enemy AI should evaluate each tick
    Until combat ends (victory, defeat, or flee)
