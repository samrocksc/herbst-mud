Feature: Skills and Talents Commands
  Players can view and manage their skills and talents using in-game commands

  Background:
    Given player is logged into the game
    And character exists

  Scenario: View skills with 'skills' command
    When player types "skills"
    Then skill list is displayed
    And each skill shows name and level
    And primary stat bonus is shown

  Scenario: View talents with 'talents' command
    When player types "talents"
    Then equipped talents are displayed
    And each talent shows name
    And slot position is indicated
    And empty slots are shown

  Scenario: Skills command shows correct format
    Given character has skills:
      | skill       | level |
      | blades      | 5     |
      | brawling    | 3     |
    When player types "skills"
    Then output shows "blades: 5"
    And output shows "brawling: 3"

  Scenario: Talents command shows equipped talents
    Given character has talents equipped:
      | slot | talent    |
      | 1    | slash     |
      | 2    | parry     |
      | 3    | smash     |
      | 4    | crash     |
    When player types "talents"
    Then slot 1 shows "slash"
    And slot 2 shows "parry"
    And slot 3 shows "smash"
    And slot 4 shows "crash"

  Scenario: Talents command shows empty slots
    Given character has only slots 1 and 2 filled
    When player types "talents"
    Then slot 1 and 2 show talent names
    And slot 3 and 4 show "empty"

  Scenario: Equip talent to valid slot
    Given character has talent "slash" available
    And slot 1 is empty
    When player types "talent equip slash 1"
    Then "slash" is equipped in slot 1
    And confirmation message is shown

  Scenario: Equip talent to slot 2
    Given character has talent "parry" available
    And slot 2 is empty
    When player types "talent equip parry 2"
    Then "parry" is equipped in slot 2

  Scenario: Equip talent to slot 3
    Given character has talent "smash" available
    And slot 3 is empty
    When player types "talent equip smash 3"
    Then "smash" is equipped in slot 3

  Scenario: Equip talent to slot 4
    Given character has talent "crash" available
    And slot 4 is empty
    When player types "talent equip crash 4"
    Then "crash" is equipped in slot 4

  Scenario: Unequip talent from slot
    Given slot 1 has "slash" equipped
    When player types "talent unequip 1"
    Then slot 1 becomes empty
    And "slash" returns to available talents

  Scenario: Equip talent replaces existing in slot
    Given slot 1 has "slash" equipped
    And character has "heavy_strike" available
    When player types "talent equip heavy_strike 1"
    Then slot 1 now has "heavy_strike"
    And "slash" returns to available talents

  Scenario: Equip to invalid slot number 0
    Given character has talent available
    When player types "talent equip slash 0"
    Then error "Invalid slot. Use 1-4." is shown
    And no change is made

  Scenario: Equip to invalid slot number 5
    Given character has talent available
    When player types "talent equip slash 5"
    Then error "Invalid slot. Use 1-4." is shown
    And no change is made

  Scenario: Equip talent not in available talents
    Given character does not have "hail_storm" available
    When player types "talent equip hail_storm 1"
    Then error "Talent not available. Unlock it first." is shown

  Scenario: Unequip from empty slot
    Given slot 2 is empty
    When player types "talent unequip 2"
    Then error "No talent in slot 2." is shown

  Scenario: Skill equip command shows info message
    When player types "skill equip blades 1"
    Then message "Skills are passive and auto-equipped based on level." is shown
    And no skill is manually equipped

  Scenario: Skills improve with use
    Given character attacks with blades weapon
    When enough attacks are made
    Then blades skill level increases
    And skill command shows new level

  Scenario: Skills command requires character
    Given player has no character
    When player types "skills"
    Then error "You must create a character first." is shown

  Scenario: Talents command requires character
    Given player has no character
    When player types "talents"
    Then error "You must create a character first." is shown

  Scenario: Talent equip requires character
    Given player has no character
    When player types "talent equip slash 1"
    Then error "You must create a character first." is shown

  Scenario: Talent command shows mana cost
    Given slot 1 has "slash" equipped
    When player types "talents"
    Then mana cost is displayed for each talent

  Scenario: Talent command shows talent type
    Given talents are equipped
    When player types "talents"
    Then each talent shows its type (attack/defense/buff/heal)

  Scenario: Help includes skills and talents
    When player types "help"
    Then "skills" command is listed
    And "talents" command is listed
    And "talent equip" syntax is shown

  Scenario Outline: Valid slot numbers
    Given character has talent available
    And slot <slot> is empty
    When player types "talent equip slash <slot>"
    Then talent is equipped in slot <slot>

    Examples:
      | slot |
      | 1    |
      | 2    |
      | 3    |
      | 4    |

  Scenario: Case insensitive talent names
    Given character has "SLASH" available
    When player types "talent equip SLASH 1"
    Then "slash" is equipped (case insensitive)

  Scenario: Alias works - 't' for talents
    Given character has talents equipped
    When player types "t"
    Then talent list is displayed

  Scenario: Skills show primary stat bonus
    Given character has blades skill level 5
    And blades primary stat is STR
    When player types "skills"
    Then stat bonus is shown for blades
