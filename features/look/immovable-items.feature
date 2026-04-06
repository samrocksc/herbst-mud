Feature: Immovable Items
  As a player
  I want to distinguish items I can pick up from environmental objects
  So that I understand what I can interact with

  Background:
    Given the game is running

  Scenario: Movable items can be taken
    Given item "rusty_pipe" has is_immovable: false
    When player types "take rusty_pipe"
    Then item is added to player inventory

  Scenario: Immovable items cannot be taken
    Given item "fountain" has is_immovable: true
    When player types "take fountain"
    Then error message "You can't pick that up" is shown

  Scenario: Immovable items show in gold color
    Given item "fountain" is immovable
    When player types "look"
    Then fountain appears in gold color

  Scenario: Immovable items can still be examined
    Given item "fountain" is immovable
    When player types "examine fountain"
    Then item description is shown

  Scenario: Immovable items show diamond marker
    Given item "old_sign" is immovable
    When player types "look"
    Then item shows "⬥" marker

  Scenario: Color coding by item type
    Given items of different types exist
    When room is displayed
    Then weapons show in red
    And armor shows in blue
    And consumables show in green
    And quest items show in purple
    And misc items show in white

  Scenario: Inventory shows no marker for movable items
    Given player has movable items in inventory
    When inventory is displayed
    Then items have no diamond marker