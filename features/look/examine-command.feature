🔴 Feature: Examine Command - Issue Look-02
  As a player
  I want to examine specific things in the world
  So that I can learn details about items, NPCs, and surroundings

  Background:
    Given the player is in a room
    And various examinable objects exist

  Scenario: Examine shows detailed item information
    Given an item exists that can be examined
    When the player types "examine <item>"
    Then detailed information should be displayed
    And the description should be longer than the room description
    And stats or properties of the item should be shown

  Scenario: Examine an NPC shows their details
    Given an NPC exists in the room
    When the player types "examine <npc>"
    Then the NPC's detailed description should appear
    And their level, class, or role should be indicated
    And any notable equipment should be listed

  Scenario: Examine self shows character stats
    Given the player wants to review their character
    When the player types "examine self" or "score"
    Then character stats should be displayed
    And HP, mana, level, XP should be shown
    And equipment and inventory summary should appear

  Scenario: Examine a direction shows what's in that way
    Given the player is deciding which way to go
    When the player types "examine north" (or just "n")
    Then a brief description of the adjacent room should be shown
    And any obvious dangers or points of interest should be hinted

  Scenario: Examine unknown object shows error
    Given the player types "examine nonexistent_thing"
    When the examine command runs
    Then an error message should indicate the thing is not found
    And the player should not crash or freeze

  Scenario: Examine partial name matches if unique
    Given an item named "Rusty Iron Sword"
    When the player types "examine sword"
    And only one item matches
    Then the item should be examined successfully
