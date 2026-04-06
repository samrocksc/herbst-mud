Feature: Examine Skill
  As a player
  I want to examine my skills and their descriptions
  So that I understand what my abilities do

  Background:
    Given player has skills learned

  Scenario: Examine shows skill description
    Given player has "slash" skill
    When player types "examine slash"
    Then skill name should display
    And skill description should be shown

  Scenario: Examine shows skill level
    Given player has "blades" at level 45
    When player examines blades skill
    Then skill level "45" should be displayed

  Scenario: Examine shows skill bonus
    Given player has blades skill at level 45
    When player examines the skill
    Then damage bonus should be shown

  Scenario: Examine shows talent tree
    Given player has talents in warrior tree
    When player examines the talent
    Then talent position in tree should be shown

  Scenario: Examine unknown skill shows error
    Given player has not learned "fireball"
    When player examines fireball
    Then message should show "You don't know this skill"

  Scenario: Examine shows skill requirements
    Given skill "heavy_strike" requires blades level 10
    When player examines heavy_strike
    Then requirements should be displayed

  Scenario: Examine shows cooldown
    Given skill "second_wind" has cooldown
    When player examines second_wind
    Then cooldown time should be shown

  Scenario: Examine shows tick cost
    Given skill "slash" costs 1 tick
    When player examines slash
    Then tick cost should be displayed

  Scenario Outline: Skill bonus calculation at different levels
    Given player has blades at level <level>
    When calculating skill bonus
    Then bonus should be <bonus>%

    Examples:
      | level | bonus |
      | 0     | 0     |
      | 25    | 0     |
      | 26    | 10    |
      | 50    | 10    |
      | 51    | 25    |
      | 75    | 25    |
      | 76    | 50    |
      | 90    | 50    |
      | 91    | 75    |
      | 99    | 75    |
      | 100   | 100   |