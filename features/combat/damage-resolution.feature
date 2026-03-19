Feature: Damage Resolution Formula
  As a player
  I want transparent damage calculations
  So that I understand how my skills and gear affect combat

  Background:
    Given a combat is active
    And the player has a "Scrap Machete" (base damage: 6)

  Scenario: Base damage applies correctly
    Given the player has no skill bonuses
    And the player has no active buffs
    When the player attacks with base damage 6
    Then the final damage is 6

  Scenario Outline: Skill bonus applies at correct level thresholds
    Given the player has blades skill at level <skill_level>
    And the weapon base damage is 6
    When damage is calculated
    Then the skill bonus is <expected_bonus>

    Examples:
      | skill_level | expected_bonus |
      | 25          | +0%           |
      | 26          | +10%          |
      | 50          | +10%          |
      | 51          | +25%          |
      | 75          | +25%          |
      | 76          | +50%          |
      | 90          | +50%          |
      | 91          | +75%          |
      | 99          | +75%          |
      | 100         | +100%         |

  Scenario: Buff bonuses stack multiplicatively with skill
    Given the player has blades skill at level 50 (+10%)
    And the player has Strength buff (+25%)
    And the weapon base damage is 6
    When damage is calculated
    Then the formula is: 6 × 1.10 × 1.25
    And the final damage is 8 (floor)

  Scenario: Defense reduces incoming damage
    Given the enemy has 5 defense
    And the enemy has 0 armor
    When raw damage 10 is applied
    Then the final damage is 5 (10 - 5)

  Scenario: Armor reduces damage multiplicatively
    Given the enemy has 10 defense
    And the enemy has 50% armor
    When raw damage 20 is applied
    Then the reduction is 10 × 0.50 = 5
    And the final damage is 15

  Scenario: Minimum damage is always 1
    Given the enemy has very high defense
    And raw damage is less than defense
    When damage is calculated
    Then the minimum damage is 1

  Scenario: Rounding uses floor
    Given base damage is 7
    And skill bonus is 25%
    And buff bonus is 25%
    When damage is calculated: 7 × 1.25 × 1.25 = 10.9375
    Then the result is floored to 10

  Scenario: Different weapon types use correct skill
    Given the player has brawling skill at level 60
    And the player uses a fist weapon
    When damage is calculated
    Then the brawling bonus applies (not blades)

  Scenario: Damage calculation order
    Given the player has skill bonus and buff bonus
    When damage is calculated
    Then skill bonus applies first (multiplicative)
    Then buff bonus applies (multiplicative)
    Then defense/armor reduction applies (multiplicative)

  Scenario: Skill level 0 grants no bonus
    Given the player has blades skill at level 0
    When attacking with a blade weapon
    Then the damage bonus is +0%

  Scenario: Full damage formula verification
    Given player: blades 45 (+10%), weapon 6 dmg, Strength buff (+25%)
    And enemy: defense 3, armor 0%
    When damage = 6 × 1.10 × 1.25 - 3
    Then damage = 8.25 - 3 = 5.25
    And final damage is 5 (floored)
