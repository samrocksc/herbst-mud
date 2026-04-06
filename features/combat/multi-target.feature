🔴 Feature: Multi-Target Combat & Targeting - Issue #11
  As a combat system
  I need targeting mechanics for fighting multiple enemies and AoE attacks
  So that combat scales with encounter size

  Background:
    Given combat involves multiple combatants
    And targeting system is active

  # SINGLE TARGETING
  Scenario: Player can target specific enemy by position
    Given multiple enemies are in combat
    And enemies are numbered 1, 2, 3
    When the player types "attack 2"
    Then the attack should target enemy #2
    And other enemies should not be affected

  Scenario: Player can target enemy by name
    Given multiple enemies are in combat
    And an enemy is named "Scrap Rat"
    When the player types "attack rat"
    Then the attack should target "Scrap Rat"
    And if multiple matches exist, a disambiguation prompt appears

  Scenario: Targeting invalid position shows error
    Given 2 enemies are in combat
    When the player types "attack 5"
    Then an error message should indicate invalid target
    And no attack should be executed

  # ALL-TARGET / AOE
  Scenario: AoE attack hits all enemies
    Given multiple enemies are in combat
    And the player has an AoE ability
    When the player uses AoE attack "attack all"
    Then all enemies should take damage
    And the damage should be reduced compared to single-target

  Scenario: AoE ability has limited uses
    Given a player uses an AoE ability
    When the ability is used
    Then the ability should enter cooldown
    And the ability should not be available immediately

  # TURN ORDER
  Scenario: Turn order displays for multiple combatants
    Given 3 enemies and 1 player are in combat
    When combat turn order is displayed
    Then the order should be shown: Player, Enemy1, Enemy2, Enemy3
    Or sorted by initiative/DEX

  Scenario: Player sees which enemy is currently targeted
    Given multiple enemies are in combat
    When a specific enemy's turn arrives
    Then the current enemy should be highlighted
    And player should see whose turn it is

  # TARGETING WHILE ENEMIES REMAIN
  Scenario: Cannot target dead enemy
    Given an enemy has been killed
    When the player tries to target that enemy
    Then an error message should indicate the enemy is dead
    And the player should be able to target remaining enemies

  Scenario: Combat ends when all enemies defeated
    Given the last enemy is killed
    When the killing blow lands
    Then combat should end
    And victory screen should be displayed
    And loot should be distributed

  # TURN MANAGEMENT
  Scenario: Combat continues until all targets resolved
    Given combat has multiple targets
    When each combatant takes their turn in order
    Then combat continues until one side is eliminated
    Or until a flee/escape condition is met

  Scenario: Dead combatants are skipped in turn order
    Given a combatant has 0 HP
    When their turn in the order arrives
    Then they should be skipped
    And the next living combatant should act
