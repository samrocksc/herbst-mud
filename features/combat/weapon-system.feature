Feature: Weapon System - First Weapon Drops
  As a player
  I want weapons to drop from enemies and be equippable
  So that I can fight more effectively

  Background:
    Given the game is running
    And the weapon system is implemented

  Scenario: Weapon has damage range
    Given weapon "rusty_sword" has:
      | field       | value |
      | minDamage   | 1     |
      | maxDamage   | 3     |
    When weapon damage is rolled
    Then damage is between 1 and 3 inclusive

  Scenario: Weapon has weapon type
    Given weapon "rusty_sword" has type "sword"
    When weapon is identified
    Then weapon type is "sword"
    And appropriate skill bonus applies

  Scenario: Weapon drop from Rust Bucket Golem
    Given enemy "Rust Bucket Golem" exists
    And enemy has guaranteed drop of "rusty_sword"
    When enemy is defeated
    Then "rusty_sword" is added to loot table
    And player can receive the weapon

  Scenario: Weapon drops are class-specific
    Given warrior character defeats enemy
    When weapon drop occurs
    Then warrior receives class-appropriate weapon
    And chef character would receive different weapon

  Scenario: Chef receives twisted pipe weapon
    Given chef character defeats "Rust Bucket Golem"
    When loot is generated
    Then chef receives "twisted_pipe" weapon
    And weapon type is "pipe"
    And damage range is 1-2

  Scenario: Warrior receives rusty sword
    Given warrior character defeats "Rust Bucket Golem"
    When loot is generated
    Then warrior receives "rusty_sword" weapon
    And weapon type is "sword"
    And damage range is 1-3

  Scenario: Weapon can be equipped
    Given player has "rusty_sword" in inventory
    When player types "equip rusty_sword"
    Then weapon is equipped
    And weapon appears in equipment slot
    And combat damage uses weapon stats

  Scenario: Weapon affects combat damage
    Given player has rusty_sword equipped with damage 1-3
    When player attacks enemy
    Then damage is calculated using weapon range
    And damage is within weapon minDamage to maxDamage

  Scenario: Unequipped player uses base damage
    Given player has no weapon equipped
    When player attacks
    Then base brawling damage is used
    And damage is lower than with weapon

  Scenario: Weapon select function for class
    Given character class is "warrior"
    When SelectWeaponForClass is called
    Then "rusty_sword" is returned

  Scenario: Weapon select function for chef
    Given character class is "chef"
    When SelectWeaponForClass is called
    Then "twisted_pipe" is returned

  Scenario: Weapon database seeding
    Given InitWeapons is called
    Then "rusty_sword" exists in database
    And "twisted_pipe" exists in database
    And both weapons have correct damage values

  Scenario: Weapon has class restriction
    Given weapon "rusty_sword" has classRestriction "warrior"
    When chef character tries to equip rusty_sword
    Then error message is shown
    And weapon is not equipped

  Scenario: Weapon damage display in examine
    Given player examines weapon "rusty_sword"
    Then weapon stats are displayed
    And damage range "1-3" is shown
    And weapon type "sword" is shown

  Scenario: First weapon acquisition tutorial
    Given new character just defeated first enemy
    When loot is generated
    Then first weapon is guaranteed
    And tutorial message is shown

  Scenario Outline: Weapon damage calculation
    Given weapon with minDamage <min> and maxDamage <max>
    When damage roll occurs
    Then result is >= <min>
    And result is <= <max>

    Examples:
      | min | max |
      | 1   | 3   |
      | 2   | 5   |
      | 3   | 8   |

  Scenario: Weapons persist after logout
    Given player has weapon equipped
    When player logs out
    And player logs back in
    Then weapon is still equipped
    And weapon stats are retained

  Scenario: Multiple weapons in inventory
    Given player has multiple weapons in inventory
    When player types "weapons" or "inventory weapons"
    Then all weapons are listed
    And stats are shown for each

  Scenario: Weapon drop rate
    Given enemy "Scrap Rat" does not have guaranteed drop
    When enemy is defeated
    Then weapon drop has random chance
    And drop rate is lower than guaranteed rate
