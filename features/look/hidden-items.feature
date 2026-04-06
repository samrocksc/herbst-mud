🔴 Feature: Hidden Items - Issue Look-12 / Issue #21
  As a player
  I want to find hidden items in the world
  So that exploring is rewarding and secrets exist

  Background:
    Given the player is in a room
    And some items may be hidden

  Scenario: Hidden items are not visible on initial look
    Given a room has hidden items
    When the player types "look"
    Then hidden items should NOT appear in the room description
    And the room should appear normal

  Scenario: Hidden items are revealed by examining carefully
    Given a room has hidden items
    When the player uses "examine" on a suspicious area
    Then a hidden item may be revealed
    And a message should indicate the discovery

  Scenario: Skill check reveals hidden items
    Given a player has high enough skill (perception/investigation)
    When the player examines the room
    Then hidden items should be revealed automatically
    And a skill check success message should display

  Scenario: Hidden items can be picked up once revealed
    Given a hidden item has been revealed
    When the player types "take <item>"
    Then the item should be added to inventory
    And the hidden state should be cleared

  Scenario: Hidden items persist once found
    Given a hidden item was found in a room
    When the player leaves and returns to the room
    Then the item should still be visible
    And it no longer requires skill to see

  Scenario: Some hidden items require quest progression
    Given a hidden item is quest-locked
    When the player tries to find it before unlocking
    Then the item should not be discoverable
    And no hint should be given until quest progress

  Scenario: GM/Builder can place hidden items
    Given a builder is creating a room
    When the builder marks an item as hidden
    Then the item should be invisible to players initially
    And it should become findable per the above rules
