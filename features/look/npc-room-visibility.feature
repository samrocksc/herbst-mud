Feature: NPC Room Visibility
  NPCs are visible in rooms based on room association and can be displayed in room descriptions

  Background:
    Given the game is running
    And NPCs exist in the database

  Scenario: Room displays NPCs present
    Given room "Town Square" has NPCs:
      | name              | disposition |
      | Mutant Raccoon    | hostile     |
      | Old Man Jenkins   | friendly    |
    When player enters "Town Square"
    Then room shows "NPCs here:"
    And "Mutant Raccoon" is listed
    And "Old Man Jenkins" is listed

  Scenario: NPC disposition shown in room
    Given NPC "Mutant Raccoon" has disposition "hostile"
    When player views room with NPC
    Then NPC is marked as hostile
    And player can see disposition indicator

  Scenario: Friendly NPCs are marked
    Given NPC "Old Man Jenkins" has disposition "friendly"
    When player views room with NPC
    Then NPC is marked as friendly
    And color coding reflects friendliness

  Scenario: Player characters shown separately from NPCs
    Given room has NPCs and player characters
    When player views room
    Then NPCs are listed in "NPCs here:" section
    And player characters are listed in "Players here:" section

  Scenario: NPCs have is_npc flag
    Given character record exists
    When character is an NPC
    Then is_npc field is true
    When character is a player
    Then is_npc field is false

  Scenario: NPCs have npc_template_id
    Given NPC is based on "gizmo" template
    When character is created
    Then npc_template_id is set to "gizmo"
    And player characters have null npc_template_id

  Scenario: NPC greeting message
    Given NPC "Gizmo" has greeting "Welcome, new traveler!"
    When player enters room with Gizmo
    Then greeting message is shown
    And message appears automatically

  Scenario: NPC level displayed
    Given NPC "Mutant Raccoon" has level 3
    When player views room
    Then NPC shows level indicator
    And format shows "(Level 3)"

  Scenario: No NPCs in room
    Given room "Empty Alley" has no NPCs
    When player enters room
    Then "NPCs here:" section is not shown
    And no NPC listing appears

  Scenario: Multiple NPCs are listed
    Given room has 5 NPCs
    When player views room
    Then all 5 NPCs are listed
    And list is formatted clearly

  Scenario: NPC list updates on NPC spawn
    Given room has no NPCs
    When NPC spawns into room
    Then NPC immediately appears in room listing
    And player sees the new NPC

  Scenario: NPC list updates on NPC leave
    Given room has NPCs
    When NPC leaves the room
    Then NPC is removed from listing
    And room description updates

  Scenario: NPC room_id tracks location
    Given NPC "Guard" is in "Town Square"
    When NPC moves to "Market"
    Then NPC room_id updates to Market
    And Town Square no longer shows Guard

  Scenario: Player cannot see hidden NPCs
    Given NPC "Secret Agent" is hidden
    When player enters room
    Then "Secret Agent" is not shown
    And only visible NPCs are listed

  Scenario: Admin can see all NPCs including hidden
    Given admin user views room
    When room has hidden NPCs
    Then admin sees all NPCs
    And hidden NPCs are marked as hidden
