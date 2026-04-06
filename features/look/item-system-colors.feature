Feature: Item System with Improvable Flag and Colors
  Items have immovable flag and custom colors for visual distinction

  Background:
    Given the game is running
    And items exist in the world

  Scenario: Immovable items cannot be picked up
    Given item "fountain" has is_immovable flag set to true
    When player types "take fountain"
    Then message "You cannot take that." is shown
    And item remains in room

  Scenario: Regular items can be picked up
    Given item "rusty_pipe" has is_immovable flag set to false
    When player types "take rusty_pipe"
    Then item is added to inventory
    And item is removed from room

  Scenario: Immovable items shown in gold color
    Given item "gold_chest" has is_immovable true
    When room is displayed
    Then "gold_chest" appears in GOLD colored text
    And visual distinction from regular items

  Scenario: Regular items use default color
    Given item "pipe" is movable
    When room is displayed
    Then "pipe" appears in default text color
    And color indicates it can be taken

  Scenario: Item has custom color field
    Given item has color field
    When room renders item
    Then custom color is applied
    And overrides default behavior

  Scenario: Immovable items show with gold color by default
    Given item "statue" has is_immovable true
    And item has no custom color set
    When room is displayed
    Then statue appears in gold text
    And gold indicates immovable

  Scenario: Custom color overrides immovable default
    Given item "blue_statue" has is_immovable true
    And item has custom color "blue"
    When room is displayed
    Then "blue_statue" appears in blue text
    And custom color takes precedence

  Scenario: Visible items show in room
    Given item has is_visible flag set to true
    When player looks at room
    Then item appears in room description

  Scenario: Hidden items not shown in room
    Given item has is_visible flag set to false
    When player looks at room
    Then item does not appear
    And item is hidden from view

  Scenario: Hidden items can be revealed
    Given item is hidden (is_visible false)
    When reveal condition is met
    Then item becomes visible
    And appears in room listing

  Scenario: Item color affects inventory display
    Given player has items with different colors
    When player types "inventory"
    Then colors are preserved in inventory list
    And visual distinction maintained

  Scenario: Improvable flag on quest items
    Given item is quest-related and immovable
    When player examines item
    Then quest hint may be revealed
    And item behavior is consistent

  Scenario: Container items can be immovable
    Given container "large_chest" has is_immovable true
    And container "small_box" has is_immovable false
    When player looks at room
    Then large_chest is gold and immovable
    And small_box is regular and takeable

  Scenario: Item schema has all required fields
    Given equipment schema is examined
    Then isImmovable field exists (bool)
    And color field exists (string)
    And isVisible field exists (bool)
    And all fields have appropriate defaults

  Scenario: Take command checks immovable flag
    Given player types "take"
    When take command processes
    Then is_immovable is checked before removal
    And appropriate message shown if immovable

  Scenario: Room description shows item colors
    Given room has items:
      | name       | color | immovable |
      | gold_statue| gold  | true      |
      | rusty_pipe |       | false     |
    When room is described
    Then gold_statue appears in gold
    And rusty_pipe appears in default color
