Feature: Examine Command
  As a player
  I want to examine objects in detail
  So that I can learn about items and NPCs

  Background:
    Given the game is running

  Scenario: Examine shows item details
    Given room contains item "rusty machete"
    When player types "examine machete"
    Then item name should be displayed
    And item description should be shown

  Scenario: Examine shows NPC details
    Given room contains NPC "Junkyard Rat"
    When player types "examine rat"
    Then NPC name should display
    And NPC description should show
    And NPC health should be visible if in combat

  Scenario: Examine unknown object shows error
    Given room contains no "magic sword"
    When player types "examine sword"
    Then message should show "You don't see that here"

  Scenario: Examine shows item stats
    Given item "scrap machete" has damage rating
    When player examines the item
    Then damage rating should be displayed

  Scenario: Examine works with partial name
    Given room contains "rusty machete"
    When player types "examine machete"
    Then item should be identified

  Scenario: Examine shows container contents
    Given player has "dumpster" container
    And container has items inside
    When player types "examine dumpster"
    Then container contents should be listed

  Scenario: Examine shows equipment slot
    Given player is wielding "scrap machete"
    When player examines the weapon
    Then equipment slot "wielded" should be shown

  Scenario: Examine case insensitive
    Given room contains "Rusty Machete"
    When player types "examine RUSTY MACHETE"
    Then item should be found

  Scenario: Examine aliases work
    Given room contains item "rusty machete"
    When player types "ex rusty machete"
    Then item details are shown
    When player types "inspect rusty machete"
    Then item details are shown

  Scenario Outline: Examine different object types
    Given room contains <object_type>
    When player examines the <object_type>
    Then relevant details should be displayed

    Examples:
      | object_type     |
      | item            |
      | NPC             |
      | container       |
      | readable object |
