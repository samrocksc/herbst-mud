Feature: Examine Command - Detailed Inspection
  As a player
  I want to examine items and NPCs closely
  So I can discover hidden details and information

  Background:
    Given the game is running
    And the player has examine skill

  Scenario: Examine item shows detailed description
    Given the room contains item "fountain"
    And the fountain has examine_desc "An old stone fountain, cracked and weathered..."
    When the player types "examine fountain"
    Then the detailed description should be shown

  Scenario: Examine falls back to description if no examine_desc
    Given the room contains item "rusty_pipe"
    And the item has description "A rusty metal pipe"
    And the item has no examine_desc
    When the player types "examine rusty_pipe"
    Then the regular description should be shown

  Scenario: Examine with alias "ex"
    Given the room contains "rusty_pipe"
    When the player types "ex rusty_pipe"
    Then the same output as "examine rusty_pipe" should be displayed

  Scenario: Examine with alias "inspect"
    Given the room contains "rusty_pipe"
    When the player types "inspect rusty_pipe"
    Then the same output as "examine rusty_pipe" should be displayed

  Scenario: Examine reveals hidden details based on skill
    Given the player has examine skill at level 30
    And the item "fountain" has a hidden detail requiring level 0
    When the player types "examine fountain"
    Then the hidden detail should be revealed

  Scenario: Examine NPC shows NPC details
    Given the room contains NPC "Guard Marco"
    When the player types "examine Guard Marco"
    Then the NPC's detailed description should be shown
    And the NPC's level should be displayed if known

  Scenario: Examine equipment shows stats
    Given the player has "scrap_machete" in inventory
    When the player types "examine scrap_machete"
    Then weapon damage should be shown
    And weapon type should be shown

  Scenario: Examine grants XP for first time
    Given the player examines item "fountain" for the first time
    When the player types "examine fountain"
    Then the player should gain 1 examine XP

  Scenario: Examine nonexistent item shows error
    When the player types "examine imaginary_item"
    Then error "You don't see that here" should be shown

  Scenario: Examine shows item type and rarity
    Given the room contains "scrap_machete"
    And the item has type "weapon" and rarity "common"
    When the player types "examine scrap_machete"
    Then the type "weapon" should be displayed
    And rarity color should match common (white/default)

  Scenario: Examine quest item shows quest info
    Given the room contains quest item "pre_ooze_key"
    And the item is part of quest "Find the Key"
    When the player types "examine pre_ooze_key"
    Then quest association should be shown
    And the item should display in purple color

  Scenario: Examine container shows contents
    Given the player has "scrap_box" in inventory
    And the box contains "old_key" and "5 coins"
    When the player types "examine scrap_box"
    Then the contents should be listed
    And "old_key" and "5 coins" should be visible