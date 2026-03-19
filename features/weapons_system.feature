Feature: Weapons System - First Weapon Drops
  As a player
  I want to obtain weapons from defeated enemies
  So that I can equip myself for combat

  Background:
    Given I am a "Warrior" class character
    And I am in the Junkyard area

  Scenario: First kill guarantees weapon drop
    Given a "Rust Bucket Golem" exists in the current room
    When I defeat the "Rust Bucket Golem"
    Then I should receive "Rusty Sword"
    And it should be equipped automatically

  Scenario: Chef class receives appropriate weapon
    Given I am a "Chef" class character
    And a "Rust Bucket Golem" exists in the current room
    When I defeat the "Rust Bucket Golem"
    Then I should receive "Twisted Pipe"

  Scenario: Weapon has correct damage values
    Given I have "Rusty Sword" in my inventory
    When I check the weapon stats
    Then I should see damage is "1-3"

  Scenario: Weapon affects combat damage
    Given I have "Rusty Sword" equipped
    When I attack an enemy
    Then my damage should include weapon bonus

  Scenario: Cannot equip wrong class weapon
    Given I am a "Mage" class character
    And I have "Rusty Sword" in my inventory
    When I try to equip "Rusty Sword"
    Then I should see "Warrior class required"
    And the weapon should remain in inventory

  Scenario: Dropped weapon goes to inventory
    Given a "Rust Bucket Golem" drops "Rusty Sword"
    When I collect the weapon
    Then "Rusty Sword" should be in my inventory
    And I can view it with "inventory"