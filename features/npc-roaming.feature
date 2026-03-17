Feature: NPC Roaming
  As a player
  I want NPCs to move between rooms dynamically
  So that the world feels alive and unpredictable

  Background:
    Given the game world has multiple rooms
    And NPCs have movement behavior defined

  Scenario: NPC spawns in designated room
    Given an NPC template "Gizmo"
    And a spawn room "Fountain Room"
    When the game initializes or respawns NPCs
    Then Gizmo should appear in the Fountain Room
    And Gizmo should have full HP

  Scenario: NPC stays in room when no roam
    Given an NPC with roaming disabled
    When time passes
    Then the NPC should remain in its current room
    And the NPC should not move between rooms

  Scenario: NPC roams to adjacent room
    Given an NPC with roaming enabled
    And the NPC is in a room with exits
    When the roam timer triggers
    Then the NPC should move to an adjacent room
    And I should see the NPC leave if I'm in the source room
    And I should see the NPC arrive if I'm in the destination room

  Scenario: NPC follows valid exits only
    Given an NPC in a room with exits: north, south
    And the NPC roams
    When the NPC moves
    Then it should only use valid exits
    And it should not move through walls or blocked exits

  Scenario: NPC roam interval varies
    Given NPCs with different roam intervals
    When time passes
    Then fast-roaming NPCs should move more frequently
    And slow-roaming NPCs should move less frequently

  Scenario: NPC despawn and respawn
    Given an NPC has been alive for a long time
    When the respawn timer triggers
    Then the NPC should despawn from its current location
    And respawn in its designated spawn room
    And the NPC should have full HP restored

  Scenario: Aggressive NPC follows player
    Given an aggressive NPC
    And I am in combat with the NPC
    When I flee to another room
    Then the aggressive NPC may follow me
    And combat may continue in the new room

  Scenario: NPC aggression states
    Given an NPC with aggression state
    When I observe the NPC
    Then the NPC should have one of: passive, neutral, aggressive
    And aggressive NPCs should attack on sight
    And neutral NPCs should attack only when provoked
    And passive NPCs should never attack

  Scenario: NPC enters room announces
    Given I am in a room
    When an NPC roams into my room
    Then I should see a message like "Gizmo enters from the north"
    And the NPC should appear in my room list

  Scenario: NPC leaves room announces
    Given an NPC is in my room
    When the NPC roams out
    Then I should see a message like "Gizmo leaves to the south"
    And the NPC should disappear from my room list

  Scenario: NPCs do not roam during combat
    Given an NPC is in combat
    When the roam timer triggers
    Then the NPC should not move
    And the NPC should remain in combat

  Scenario: NPC path follows zone boundaries
    Given NPCs are assigned to a zone
    When NPCs roam
    Then they should stay within their zone boundaries
    And they should not cross into other zones

  Scenario: Boss NPCs have restricted roaming
    Given a boss NPC
    When the boss roams
    Then it should have limited movement range
    Or the boss should not roam at all
    And players should find the boss in its designated area

  Scenario: Player sees NPC movement in adjacent rooms
    Given I am in a room
    And an NPC is in an adjacent room
    When I look through the exit
    Then I may see the NPC in the adjacent room description