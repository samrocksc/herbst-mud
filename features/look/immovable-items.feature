Feature: Movable vs Immovable Items
  As a player
  I want to know which items can be picked up
  So that I understand what I can interact with

  Background:
    Given the game is running
    And the player is in a room with items

  Scenario: Movable items can be taken
    Given the room contains "rusty_pipe" which is movable
    When the player types "take rusty_pipe"
    Then the item should be added to inventory
    And success message "You pick up the rusty_pipe" should be shown

  Scenario: Immovable items cannot be taken
    Given the room contains "fountain" which is immovable
    When the player types "take fountain"
    Then the item should remain in the room
    And error "You can't pick up the fountain" should be shown

  Scenario: Immovable items show gold color
    Given the room contains immovable item "fountain"
    When the player types "look"
    Then the fountain should be displayed in gold color
    And the item should have the ⬥ marker prefix

  Scenario: Movable items show no special marker
    Given the room contains movable item "rusty_pipe"
    When the player types "look"
    Then the pipe should have no marker prefix

  Scenario: Color coding by item type
    Given the room contains:
      | item         | type    | expected_color |
      | rusty_pipe   | weapon  | red            |
      | salvaged_plate | armor | blue           |
      | health_potion | consumable | green      |
      | torn_cloth   | misc    | white          |
    When the player types "look"
    Then each item should display in its correct color

  Scenario: Quest items show purple
    Given the room contains quest item "pre_ooze_key"
    When the player types "look"
    Then the key should be displayed in purple

  Scenario: Look at immovable shows description
    Given the room contains immovable "fountain"
    When the player types "look at fountain"
    Then the fountain's description should be shown

  Scenario: Examine immovable shows details
    Given the room contains immovable "fountain"
    When the player types "examine fountain"
    Then the fountain's detailed description should be shown

  Scenario: Invisible items not shown in room
    Given the room contains invisible item "secret_note"
    When the player types "look"
    Then "secret_note" should not appear in the output

  Scenario: Drop works for movable items
    Given the player has "rusty_pipe" in inventory
    When the player types "drop rusty_pipe"
    Then the item should be in the room
    And the item should be removed from inventory