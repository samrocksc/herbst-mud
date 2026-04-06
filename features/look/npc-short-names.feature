Feature: NPC Short Names and Partial Matching
  Players can target NPCs using partial or short names instead of full names

  Background:
    Given player is in combat
    And multiple NPCs are in the room

  Scenario: Attack NPC by partial name
    Given NPC "Combat Dummy" exists
    When player types "attack dummy"
    Then player attacks "Combat Dummy"
    And full name is not required

  Scenario: Attack NPC by first word of name
    Given NPC "Combat Dummy" exists
    When player types "attack combat"
    Then player attacks "Combat Dummy"

  Scenario: Attack NPC by short abbreviation
    Given NPC "Combat Dummy" exists
    When player types "attack cd"
    Then player attacks "Combat Dummy"
    And "cd" is accepted as valid

  Scenario: Single letter match
    Given room has only one NPC starting with "c"
    When player types "attack c"
    Then player attacks that NPC
    And single letter works for unique match

  Scenario: Case insensitive matching
    Given NPC "Combat Dummy" exists
    When player types "attack COMBAT"
    Then matching is case insensitive
    And attack succeeds

  Scenario: Mixed case matching
    Given NPC "Combat Dummy" exists
    When player types "attack CoMbAt"
    Then matching works regardless of case

  Scenario: Ambiguous partial match asks for clarification
    Given NPCs "Combat Dummy" and "Combat Chef" exist
    When player types "attack combat"
    Then system asks "Which do you mean?"
    And options are listed

  Scenario: Disambiguate with more characters
    Given NPCs "Combat Dummy" and "Combat Chef" exist
    When player types "attack combat d"
    Then "Combat Dummy" is uniquely matched
    And attack proceeds

  Scenario: Number targeting still works
    Given NPCs "Combat Dummy" and "Combat Chef" exist
    When player types "attack 1"
    Then first NPC is targeted

  Scenario: Target command with partial name
    Given player is not in combat
    And NPC "Mutant Raccoon" exists
    When player types "target raccoon"
    Then "Mutant Raccoon" becomes current target

  Scenario: Target command outside combat
    Given player is not in combat
    When player types "target dummy"
    Then NPC is selected
    And confirmation shown

  Scenario: Examine with partial name
    Given NPC "Old Man Jenkins" exists
    When player types "examine jenkins"
    Then NPC details are shown

  Scenario: Talk with partial NPC name
    Given NPC "Gizmo" exists
    When player types "talk giz"
    Then conversation with Gizmo starts

  Scenario: Short names work for all NPC interactions
    Given NPC "Guard Captain" exists
    When player uses any NPC command with partial name
    Then system correctly matches NPC
    And command executes

  Scenario: No match found error
    Given NPC "Target" does not exist
    When player types "attack nonexistent"
    Then error "No NPC found matching 'nonexistent'" is shown

  Scenario: Number out of range
    Given room has 3 NPCs
    When player types "attack 5"
    Then error "Invalid target number" is shown

  Scenario: Short name persists until changed
    Given player targets "Combat Dummy" with "attack dummy"
    When player types "attack" without specifying target
    Then last targeted NPC is attacked

  Scenario: Player characters not matched for NPC commands
    Given player character "Sam123" is in room
    When player types "attack sam"
    Then NPC command handles gracefully
    And player characters are not affected

  Scenario: Whitespace trimmed from partial names
    Given NPC "Combat Dummy" exists
    When player types "attack   dummy   "
    Then whitespace is trimmed
    And matching works correctly
