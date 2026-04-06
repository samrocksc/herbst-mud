🔴 Feature: Item Schema - Issue Look-11 / Issue #19
  As a game system
  I need a consistent item data schema
  So that all items have proper structure and behavior

  Background:
    Given the item system is initialized
    And items can exist in various contexts (room, inventory, container)

  Scenario: Item has required fields
    Given an item is created
    Then the item should have:
      | field       | description                    |
      | id          | unique identifier              |
      | name        | display name                   |
      | description | long description                |
      | item_type   | category (weapon, armor, etc.)  |
      | is_movable  | can player take it?            |

  Scenario: Item has optional fields
    Given an item exists
    Then the following optional fields may exist:
      | field         | description                        |
      | stats         | stat bonuses when equipped         |
      | skill_bonus   | skill level bonuses                |
      | damage        | damage values (for weapons)         |
      | armor         | armor values (for armor)           |
      | skill_required | minimum skill to use               |
      | level_required | minimum level to use               |
      | quest_id      | associated quest (if any)          |
      | hidden_detail | detail revealed by skill check     |

  Scenario: Item types are defined
    Given the item system is initialized
    Then item types should include:
      | type      | examples                      |
      | weapon    | sword, dagger, bow            |
      | armor     | helmet, chest, boots          |
      | consumable| potion, scroll, food           |
      | quest     | quest key, quest item         |
      | treasure  | gold, gems                    |
      | misc      | books, keys, containers       |

  Scenario: Items can be stacked
    Given an item is consumable or stackable
    When multiple of the same item are in inventory
    Then they should stack with a quantity
    And using one reduces the stack by 1

  Scenario: Non-stackable items do not merge
    Given two non-stackable items of same type
    When they are in inventory
    Then they should appear as separate entries
    And each has quantity 1

  Scenario: Items have slot requirements
    Given an item is equipment
    When the player tries to equip it
    Then the item should require a valid slot (main hand, off hand, etc.)
    And if slot is occupied, the player must unequip first

  Scenario: Item durability (if applicable)
    Given an item has durability
    When the item is used in combat
    Then durability should decrease
    And at 0 durability, the item may break or need repair

  Scenario: Item quality affects stats
    Given items can have quality levels (common, uncommon, rare, epic)
    When items of same type have different qualities
    Then higher quality should provide better stats
    And quality should be visible in item description

  Scenario: Items can be traded/sold
    Given a player has an item
    When the player visits a shop or trading post
    Then the item can be sold for gold
    Or the item can be traded to another player
