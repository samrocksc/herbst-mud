🔴 Feature: Damage Resolution Formula - Issue #10
  As a combat system
  I need damage calculated with skill bonuses, buffs, and defense
  So that combat stats and abilities have meaningful impact

  Background:
    Given combat is in progress
    And damage resolution system is active

  # BASE DAMAGE
  Scenario: Base damage is applied correctly
    Given an attack has base damage of 10
    When the attack resolves against a target
    Then the target should take 10 damage
    Before any modifiers

  Scenario: Damage can be zero minimum
    Given damage calculation results in negative or zero
    When the damage is applied
    Then the target should take 0 damage minimum
    And no healing occurs from the attack

  # SKILL BONUS
  Scenario Outline: Skill level provides damage bonus
    Given a player has a weapon skill at level <skill_level>
    When the player attacks with that weapon type
    Then a damage bonus of <bonus>% should be applied

    Examples:
      | skill_level | bonus |
      | 0           | 0%    |
      | 26          | 10%   |
      | 51          | 25%   |
      | 76          | 50%   |
      | 91          | 75%   |
      | 100         | 100%  |

  Scenario: Skill bonus applies to weapon-matched attacks
    Given a player has blades skill at level 50
    And the player attacks with a sword (blade weapon)
    When damage is calculated
    Then the skill bonus should apply

  Scenario: Mismatched skill provides no bonus
    Given a player has blades skill at level 50
    But the player attacks with a blunt weapon
    When damage is calculated
    Then no skill bonus should apply

  # BUFF MODIFIER
  Scenario: Buff increases damage output
    Given a player has BUFF_STRENGTH active
    When the player attacks
    Then damage should include +25% from the buff

  Scenario: Multiple buffs stack additively
    Given a player has BUFF_STRENGTH and BUFF_POWER active
    When the player attacks
    Then damage should include both buff bonuses
    And the total buff bonus is the sum of individual buffs

  # DEFENSE REDUCTION
  Scenario: Enemy defense reduces incoming damage
    Given an enemy has defense value of 10
    And armor rating of 20%
    When a player attack of 20 damage lands
    Then damage reduction = 10 × (1 - 0.2) = 8
    And the enemy takes 12 damage

  Scenario: High defense can fully negate low damage
    Given an enemy has very high defense
    When a weak attack is attempted
    Then the enemy may take 0 damage
    And a "glancing blow" or similar message may display

  # FULL FORMULA
  Scenario: Complete damage formula is applied
    Given:
      - Base damage is 20
      - Player has 50% skill bonus (level 51 blades)
      - Player has +25% BUFF_STRENGTH
      - Enemy has 10 defense and 20% armor
    When the attack lands
    Then raw damage = 20 × 1.50 × 1.25 = 37.5
    And final damage = 37.5 - (10 × 0.8) = 29.5
    And the enemy takes 29 or 30 damage (rounded)

  # CRITICAL HIT
  Scenario: Critical hit doubles damage
    Given a player lands a critical hit
    When damage is calculated
    Then the final damage should be doubled
    And the combat log should indicate critical hit

  Scenario: Critical hit chance is affected by stats
    Given a player has high dexterity
    When attacks are made
    Then critical hit chance should be increased
    And critical hits should be possible

  # DAMAGE TYPES
  Scenario: Physical damage is reduced by armor
    Given an attack deals physical damage
    When the attack hits an armored target
    Then armor should reduce the damage

  Scenario: Magical damage bypasses physical armor
    Given an attack deals magical damage
    When the attack hits a target with high physical armor
    Then armor should not reduce magical damage
    And a different defense stat should apply

  # RESISTANCE AND WEAKNESS
  Scenario: Elemental weakness doubles damage
    Given an enemy is weak to fire
    And an attack deals fire damage
    When the attack lands
    Then fire damage should be doubled
    And the combat log should indicate "weakness"

  Scenario: Elemental resistance halves damage
    Given an enemy has fire resistance
    And an attack deals fire damage
    When the attack lands
    Then fire damage should be halved
    And the combat log should indicate "resistance"
