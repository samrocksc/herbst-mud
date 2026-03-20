Feature: Container System
  As a player
  I want to open containers and take items from them
  So that I can access loot and stored items in the game world

  Background:
    Given I am at a location with a "dusty_crate"
    And the crate contains "rusty_knife", "torn_cloth", "broken_watch"

  Scenario: Look in container shows contents
    When I type "look in dusty_crate"
    Then I see the container contents:
      | item          | type  |
      | rusty_knife   | weapon |
      | torn_cloth    | misc  |
      | broken_watch  | misc  |

  Scenario: Look in empty container
    Given the crate is empty
    When I type "look in dusty_crate"
    Then I see "The dusty_crate is empty"

  Scenario: Open a closed container
    Given the crate is closed
    When I type "open dusty_crate"
    Then the crate is now open
    And I see "You open the dusty_crate"

  Scenario: Close an open container
    Given the crate is open
    When I type "close dusty_crate"
    Then the crate is now closed
    And I see "You close the dusty_crate"

  Scenario: Take item from container
    When I type "take rusty_knife from dusty_crate"
    Then rusty_knife is in my inventory
    And rusty_knife is removed from the crate
    And I see "You take the rusty_knife from the dusty_crate"

  Scenario: Take item not in container
    When I type "take nonexistent from dusty_crate"
    Then I see error: "rusty_nail is not in the dusty_crate"

  Scenario: Take from empty container
    Given the crate is empty
    When I type "take rusty_knife from dusty_crate"
    Then I see error: "nothing to take from dusty_crate"

  Scenario: Take from closed container
    Given the crate is closed
    When I type "take rusty_knife from dusty_crate"
    Then I see error: "the dusty_crate is closed"
    And the item remains in the container

  Scenario: Container with capacity limit
    Given a container with capacity 2
    And the container already holds 2 items
    When I try to add an item
    Then the item does not fit
    And I see "the dusty_crate is full"

  Scenario: Look in closed container shows locked message
    Given the crate is closed
    When I type "look in dusty_crate"
    Then I see "The dusty_crate is closed"
    And I cannot see contents

  Scenario: Container types identified by schema or name
    Given an item has "is_container: true" in schema
    When I interact with it
    Then it behaves as a container
    And name patterns (crate/chest/bag) also work as fallback

  Scenario: Put item in container
    Given I have "rusty_knife" in inventory
    And the crate is open with capacity
    When I type "put rusty_knife in dusty_crate"
    Then the knife is in the crate
    And the knife is removed from inventory
    And I see "You put the rusty_knife in the dusty_crate"

  Scenario: Container weight limits
    Given a container has weight limit 10
    And an item weighs 8
    When I try to add another item weighing 5
    Then the item does not fit
    And I see "the dusty_crate cannot hold that much weight"

  Scenario: Locked container requires key
    Given the crate is locked
    When I type "open dusty_crate"
    Then I see "The dusty_crate is locked"
    And when I use "old_key" on dusty_crate
    Then the crate unlocks
    And I can open it
