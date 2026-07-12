Feature: Player Web Client — Offline PC visibility in rooms

  Background:
    Given I am authenticated as sma
    And I have a character "smack" in Ooze Surfers

  Scenario: Offline player characters are not visible in the room list
    Given a player character "ChefHuman" exists in room 1 of Ooze Surfers
    And "ChefHuman" does not have an active WebSocket connection
    And I am connected as "smack" in room 1
    When the server sends the room screen for room 1
    Then the character list should include "smack"
    And the character list should NOT include "ChefHuman"

  Scenario: NPCs are always visible regardless of connection state
    Given an NPC "Theodore Von Rad" exists in room 1 of Ooze Surfers
    And I am connected as "smack" in room 1
    When the server sends the room screen for room 1
    Then the character list should include "Theodore Von Rad"
    And "Theodore Von Rad" should have type "npc"

  Scenario: A player that disconnects disappears from the room list
    Given I am connected as "smack" in room 1
    And another player "ChefHuman" is also connected in room 1
    When "ChefHuman" disconnects their WebSocket
    And the server sends the room screen for room 1
    Then the character list should include "smack"
    And the character list should NOT include "ChefHuman"

  Scenario: Examine command does not find offline PCs
    Given a player character "ChefHuman" exists in room 1 of Ooze Surfers
    And "ChefHuman" does not have an active WebSocket connection
    And I am connected as "smack" in room 1
    When I type "examine ChefHuman"
    Then the response should be "You don't see ChefHuman here."