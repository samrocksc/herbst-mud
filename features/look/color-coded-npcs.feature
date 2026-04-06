Feature: Color-Coded NPCs and Players
  NPCs and player characters are visually distinguished by color in room descriptions

  Background:
    Given player is in a room with both NPCs and player characters

  Scenario: NPCs display in red color
    Given NPC "Mutant Raccoon" is in the room
    When room is rendered
    Then "Mutant Raccoon" appears in RED text

  Scenario: Player characters display in green color
    Given player character "Sam123" is in the room
    When room is rendered
    Then "Sam123" appears in GREEN text

  Scenario: NPC indicator shows "(NPC)"
    Given NPC is displayed
    Then NPC name is followed by "(NPC)" indicator
    And indicator is also in red

  Scenario: Player indicator shows "(Player)"
    Given player character is displayed
    Then name is followed by "(Player)" indicator
    And indicator is in green

  Scenario: Multiple NPCs all red
    Given room has NPCs:
      | name              |
      | Mutant Raccoon    |
      | Old Man Jenkins   |
    When room is rendered
    Then all NPCs are in red text
    And distinction from players is clear

  Scenario: Multiple players all green
    Given room has players:
      | name     |
      | Sam123   |
      | Player2  |
    When room is rendered
    Then all players are in green text

  Scenario: NPC and player list separated
    Given room has NPCs and players
    When room description is shown
    Then NPCs are in their own section
    And players are in separate section

  Scenario: Example room output format
    Given room has NPC "Combat Dummy" and player "Sam"
    When room is displayed
    Then output looks like:
      """
      NPCs here:
        • Combat Dummy (NPC) [RED]
      Players here:
        • Sam (Player) [GREEN]
      """

  Scenario: Color coding works in look command
    Given room has NPCs and players
    When player types "look"
    Then color coding is preserved
    And NPCs red, players green

  Scenario: Color coding in examine
    Given NPC "Guard" is examined
    When examine output is shown
    Then name is colored red
    And (NPC) indicator shown

  Scenario: Color coding in combat target list
    Given player is in combat with NPCs
    When combat UI shows enemy list
    Then NPC names are in red

  Scenario: Color coding in who command
    Given player types "who"
    When player list is shown
    Then all players are shown in green

  Scenario: NPCs at distance shown with colors
    Given NPCs are visible at distance
    When room is described
    Then color coding applies at all distances

  Scenario: Color coding accessibility - redundancy
    Given color coding is enabled
    Then (NPC) and (Player) indicators provide redundancy
    And users who cannot see colors can still distinguish

  Scenario: Clear distinction between NPC and Player
    Given player "TestPlayer" and NPC "Test NPC" exist
    When room is rendered
    Then "TestPlayer" is green
    And "Test NPC" is red
    And confusion is impossible

  Scenario: Admin NPCs still red
    Given admin NPC "Guard Captain" exists
    When room is rendered
    Then admin NPCs are still red
    And distinction from regular NPCs is by indicator, not color
