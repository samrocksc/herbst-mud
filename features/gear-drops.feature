Feature: Gear Drops
  As a player
  I want NPCs to drop gear and items when killed
  So that I can loot equipment and progress my character

  Background:
    Given I have created a character
    And I am in a room with an NPC

  Scenario: NPC drops items on death
    Given I kill an NPC
    When the NPC dies
    Then items should drop to the room floor
    And I should see what items dropped
    And I should be able to pick them up

  Scenario: Loot table determines drops
    Given an NPC has a loot table defined
    When the NPC dies
    Then items should be selected from the loot table
    And drops should follow defined probability weights

  Scenario: Common items drop frequently
    Given an NPC with common loot tier
    When I kill the NPC
    Then common items should drop frequently
    And rare items should rarely drop

  Scenario: Rare items from tough NPCs
    Given a high-level or boss NPC
    When I kill the NPC
    Then rare items should have higher drop chance
    And I should potentially receive better gear

  Scenario: Gold drops from NPCs
    Given I kill an NPC
    When the NPC dies
    Then gold coins should have a chance to drop
    And gold amount should scale with NPC difficulty

  Scenario: Item rarity tiers
    Given items have rarity tiers
    When I view dropped items
    Then I should see rarity indicators
    And rarities should include: common, uncommon, rare, epic

  Scenario: Item quality affects stats
    Given an item has quality rating
    When I equip the item
    Then higher quality should provide better stats
    And item quality should be visible in description

  Scenario: Equip dropped item
    Given an item has dropped to the floor
    And the item is equippable
    When I type "get [item]" and "equip [item]"
    Then I should be wearing the item
    And my stats should update accordingly
    And the item should appear in my equipment slot

  Scenario: Unequip item
    Given I have an item equipped
    When I type "unequip [item]"
    Then the item should move to my inventory
    And my stats should update accordingly

  Scenario: Drop to floor on unequip
    Given I am at inventory capacity
    And I unequip an item
    When my inventory is full
    Then the item should drop to the floor
    And I should see a message about dropping the item

  Scenario: Loot command picks up all
    Given multiple items have dropped
    When I type "loot" or "get all"
    Then all items on the floor should be picked up
    And items should go to my inventory

  Scenario: Item shows on room look
    Given items have dropped to the floor
    When I type "look"
    Then I should see items listed on the ground
    And I should see item names and colors by rarity

  Scenario: Boss NPC drops guaranteed loot
    Given I kill a boss NPC
    When the boss dies
    Then at least one guaranteed item should drop
    And the item should be rare or better quality

  Scenario: Empty loot table - no drops
    Given an NPC with no loot table
    When I kill the NPC
    Then no items should drop
    And I should only receive XP