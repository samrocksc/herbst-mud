Feature: Database Schema for Skills and Talents
  Skills and talents are stored in the database with proper schema

  Background:
    Given database migration is complete
    And InitSkills has been called
    And InitTalents has been called

  Scenario: Skills table exists with correct fields
    Given skill record in database
    Then skill has "id" field
    And skill has "name" field
    And skill has "type" field (weapon/magic/armor)
    And skill has "primary_stat" field (STR/DEX/INT/WIS)
    And skill has "description" field

  Scenario: Talents table exists with correct fields
    Given talent record in database
    Then talent has "id" field
    And talent has "name" field
    And talent has "type" field (attack/defense/buff/heal/passive)
    And talent has "mana_cost" field
    And talent has "cooldown" field
    And talent has "required_skill" field (nullable)

  Scenario: Character skills table exists
    Given character_skills table
    Then table has "character_id" foreign key
    And table has "skill_id" foreign key
    And table has "level" field
    And table tracks character-skill relationship

  Scenario: Character talents table with slot limit
    Given character_talents table
    Then table has "character_id" foreign key
    And table has "talent_id" foreign key
    And table has "slot" field (1-4)
    And maximum 4 talents per character is enforced

  Scenario: Available talents table
    Given available_talents table
    Then table tracks which talents character can equip
    And table has "character_id" foreign key
    And table has "talent_id" foreign key
    And talents must be unlocked before equipping

  Scenario: InitSkills seeds 12 skills
    When InitSkills is called
    Then database contains 12 skills
    And skills include "blades", "staves", "knives", "martial", "brawling"
    And skills include "tech", "fire_magic", "water_magic", "wind_magic"
    And skills include "light_armor", "cloth_armor", "heavy_armor"

  Scenario: InitTalents seeds 14 talents
    When InitTalents is called
    Then database contains 14 talents
    And talents include "slash", "parry", "smash", "crash"
    And talents include "shield_bash", "battle_cry", "second_wind"
    And talents include "hail_storm", "iron_will", "heavy_strike"
    And talents include "dodge", "quick_slash", "shield_wall", "focus"

  Scenario: Skill has correct type classification
    Given skill "blades" has type "weapon"
    Then skill "staves" has type "weapon"
    And skill "fire_magic" has type "magic"
    And skill "cloth_armor" has type "armor"

  Scenario: Skill has correct primary stat
    Given skill "blades" primary stat is "STR"
    Then skill "knives" primary stat is "DEX"
    And skill "tech" primary stat is "INT"
    And skill "wind_magic" primary stat is "WIS"

  Scenario: Talent has mana cost
    Given talent "slash" has mana cost
    Then talent "hail_storm" has mana cost
    And talent "second_wind" has mana cost
    And mana cost is deducted when talent is used

  Scenario: Talent has required skill
    Given talent "slash" requires skill "blades"
    When player tries to equip slash
    Then player must have blades skill at level 1
    And if skill missing, error is shown

  Scenario: Character starts with default skills
    Given new fighter character is created
    When InitAvailableTalentsForCharacter is called
    Then character has "blades" skill at level 1
    And character has "brawling" skill at level 1
    And character has default talents available

  Scenario: Character skills table enforces uniqueness
    Given character already has skill "blades"
    When duplicate skill assignment is attempted
    Then database prevents duplicate
    Or existing skill level is updated

  Scenario: Character talents slot uniqueness
    Given character has talent in slot 1
    When new talent is assigned to slot 1
    Then previous talent is replaced
    Or error is shown if replacement not allowed

  Scenario: Available talents API endpoint
    Given character has ID
    When GET /characters/:id/available-talents is called
    Then list of available talents is returned
    And POST /characters/:id/available-talents adds talent
    And DELETE /characters/:id/available-talents/:talentId removes talent

  Scenario: Skills API endpoint
    Given character has ID
    When GET /characters/:id/skills is called
    Then list of character skills with levels is returned

  Scenario: Talents API endpoint
    Given character has ID
    When GET /characters/:id/talents is called
    Then list of equipped talents with slots is returned

  Scenario: Skill level affects talent availability
    Given talent "heavy_strike" requires skill "blades" level 5
    And character has blades level 3
    When character tries to unlock heavy_strike
    Then unlock fails
    And message shows required skill level

  Scenario: Skill level met unlocks talent
    Given talent "heavy_strike" requires skill "blades" level 5
    And character has blades level 5
    When character tries to unlock heavy_strike
    Then unlock succeeds
    And talent appears in available talents

  Scenario Outline: Default fighter talents
    Given new fighter character is created
    Then default available talents include:
      | talent      |
      | slash       |
      | parry       |
      | smash       |
      | crash       |

  Scenario: Skills and talents are separate systems
    Given character has "blades" skill at level 10
    And character has "slash" talent equipped
    When skill level increases
    Then talent damage bonus increases
    But talent remains separately managed

  Scenario: Talent required skill is nullable
    Given talent "parry" has no required skill
    When any character tries to unlock parry
    Then unlock succeeds regardless of skills
    And parry is available to all classes
