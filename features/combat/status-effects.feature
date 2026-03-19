Feature: Status Effects (Buffs/Debuffs/DoT)
  As a player
  I want combat status effects that change during battle
  So that combat has strategic depth beyond raw damage

  Background:
    Given a combat is active
    And the player is engaged with a "Junk Dog"

  Scenario: Bleeding deals damage over time
    Given the enemy has bleeding status (3 ticks)
    When a combat tick passes
    Then the enemy takes 1 damage from bleeding
    And the bleed duration decreases by 1

  Scenario: Bleeding expires after duration
    Given the enemy has bleeding status (3 ticks)
    When 3 combat ticks pass
    Then the enemy takes 3 total bleeding damage
    And the bleeding effect is removed

  Scenario: Poison deals damage over time
    Given the enemy has poison status (5 ticks, 2 damage per tick)
    When a combat tick passes
    Then the enemy takes 2 damage from poison
    And the poison duration decreases by 1

  Scenario: Stunned prevents actions
    Given the player is stunned (2 ticks)
    When a combat tick occurs
    Then the player cannot take actions
    And the stun duration decreases by 1

  Scenario: Stunned expires and player can act
    Given the player is stunned (2 ticks)
    When 2 combat ticks pass
    Then the stun effect is removed
    And the player can act normally

  Scenario: Blinded reduces accuracy
    Given the player is blinded (3 ticks)
    When the player attacks
    Then the attack has -50% accuracy
    And the blind duration decreases by 1

  Scenario: Burning deals damage and reduces accuracy
    Given the player is burning (4 ticks)
    When a combat tick passes
    Then the player takes 1 damage from burning
    And the player's attacks have -10% accuracy
    And the burn duration decreases by 1

  Scenario: Strength buff increases damage
    Given the player has Strength buff (5 ticks, +25% damage)
    When the player attacks
    Then the damage is increased by 25%
    And the buff duration decreases by 1

  Scenario: Shield buff reduces incoming damage
    Given the player has Shield buff (3 ticks, -25% incoming damage)
    When the enemy attacks
    Then the incoming damage is reduced by 25%
    And the buff duration decreases by 1

  Scenario: Multiple effects stack correctly
    Given the player has Strength buff
    And the player has Shield buff
    When damage is calculated
    Then buffs apply their multipliers correctly
    And both effects decrease duration by 1 on tick

  Scenario: Effects process in correct order each tick
    Given the enemy has multiple status effects
    When processAllStatusEffects runs
    Then DoT damage is applied first
    Then buff modifiers are calculated
    Then debuff penalties are applied
    And expired effects are removed

  Scenario Outline: Skill level affects status effect damage
    Given the player has the "blades" skill at level <level>
    When the player uses a bleeding weapon
    Then the bleed potency is <potency>

    Examples:
      | level | potency |
      | 1     | 1       |
      | 50    | 1       |
      | 100   | 2       |
