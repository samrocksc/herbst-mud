Feature: Look Command - Room Display
  Players can view their current location and its contents.

  Background:
    Given player is logged in
    And player is in a room

  Scenario: Look shows room name and description
    When player enters "look"
    Then room title is displayed
    And room description is shown

  Scenario: Look shows exits
    Given room has exits north and south
    When player enters "look"
    Then exits section shows "[N]orth" and "[S]outh"

  Scenario: Look shows items in room
    Given room contains "rusty_pipe" and "old_sign"
    When player enters "look"
    Then items listed in "HERE:" section

  Scenario: Look shows NPCs in room
    Given room contains NPC "Guard Marco"
    When player enters "look"
    Then NPC is displayed in "HERE:" section

  Scenario: Look alias 'l' works
    When player enters "l"
    Then same output as "look" command

  Scenario: Look at specific item shows details
    Given room contains "old_sign"
    When player enters "look at old_sign"
    Then item description is displayed

  Scenario: Look at NPC shows details
    Given NPC "Guard Marco" is in room
    When player enters "look at Guard Marco"
    Then NPC description is shown

  Scenario: Look at self shows character
    When player enters "look at me"
    Then character description is displayed

  Scenario: Room uses box-drawing characters
    When room is rendered
    Then uses "═" for top/bottom borders
    And uses "─" for section dividers

  Scenario: Look at direction shows destination
    Given room has exit "north" to "Foggy Gate"
    When player enters "look north"
    Then destination "Foggy Gate" is shown
    And exit description is displayed if available

  Scenario: Look shows player health status
    Given the player has 75/100 HP
    When player enters "look"
    Then health bar should show current HP
    And percentage should be displayed

  Scenario: Look hides invisible entities
    Given room contains invisible NPC "hidden_spy"
    When player enters "look"
    Then "hidden_spy" should not appear
    But the spy can still be seen with "look at hidden_spy" if they have invisibility skill

  Scenario: Look updates after combat
    Given combat ends in victory
    When player enters "look"
    Then the defeated enemy should not appear
    And any new loot on ground should be shown