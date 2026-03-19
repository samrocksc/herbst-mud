Feature: Enemy AI Decision Making
  As a player
  I want enemies that make intelligent combat decisions
  So that combat feels dynamic and challenging

  Background:
    Given a combat is active with player and enemy
    And the enemy is a "Junk Dog" with 25 HP
    And the tick loop is running

  Scenario: Enemy attacks normally when at full health
    Given the enemy HP is above 25%
    And no special abilities are ready
    When a combat tick occurs
    Then the enemy should select basic attack

  Scenario: Enemy considers fleeing at low health
    Given the enemy HP drops below 25%
    When the enemy AI decides an action
    Then the AI should check flee chance
    Or the AI should check for heal ability

  Scenario: Enemy uses special ability when ready
    Given the enemy has a "Bite" ability
    And the ability cooldown is 0 ticks
    When the enemy AI decides an action
    Then the enemy should use Bite

  Scenario: Enemy respects cast time for abilities
    Given the enemy is channeling "Crush" ability
    And the cast time is 2 ticks
    When tick 1 passes
    Then the enemy should still be channeling
    When tick 2 completes
    Then the enemy ability should resolve

  Scenario: Scrap Rat uses only basic attacks
    Given the enemy is a "Scrap Rat"
    And it has no special abilities
    When the enemy AI decides an action
    Then the action should be basic attack

  Scenario: Old Scrap has multiple abilities
    Given the enemy is "Old Scrap"
    And it has "Crush" ability (2 tick cast)
    And it has "Scavenge" ability (2 tick cast)
    When the AI evaluates abilities
    Then it should consider both abilities
    And select based on situation

  Scenario: Ooze Spawn explodes on death
    Given the enemy is an "Ooze Spawn"
    When the enemy HP reaches 0
    Then the enemy should trigger death explosion
    And the explosion should deal damage to the player

  Scenario: Enemy AI decision tree order
    Given the enemy is in combat
    When the AI processes a decision
    Then it should check health first (flee/heal)
    Then it should check ability availability
    Then it should default to basic attack
