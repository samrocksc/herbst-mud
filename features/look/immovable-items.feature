🔴 Feature: Immovable Items - Issue Look-05
  As a player
  I want to know which items cannot be taken
  So that I don't waste time trying to pick up the environment

  Background:
    Given the player is in a room
    And some items are immovable

  Scenario: Trying to take an immovable item shows informative message
    Given an item is marked as immovable
    When the player types "take <item>"
    Then a message should indicate the item cannot be taken
    And the message should describe WHY (too heavy, attached, part of room)

  Scenario: Immovable items display in gold/yellow color
    Given immovable items exist in a room
    When the player looks at the room
    Then immovable items should be displayed in gold/yellow color
    And movable items should be in a different color

  Scenario: Immovable item is described as part of surroundings
    Given an immovable object exists (e.g., a tree, boulder)
    When the player looks at the room
    Then the description should mention the object
    And it should be clear it is part of the environment

  Scenario: Can examine immovable items
    Given an immovable item exists
    When the player types "examine <immovable>"
    Then a detailed description should be available
    And the examine output should confirm it cannot be taken

  Scenario: Immovable items do not appear in room items list for taking
    Given the player types "items" or "inventory"
    When the room inventory is displayed
    Then immovable items should be clearly separated
    And the player should not attempt to take them
