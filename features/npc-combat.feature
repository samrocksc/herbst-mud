Feature: NPC Combat
  As a player
  I want to fight NPCs that can attack back and have combat behavior
  So that combat is challenging and interactive

  Background:
    Given I have created a character
    And I am in a room with an NPC

  Scenario: NPC has combat stats
    Given an NPC "Gizmo" exists in the room
    When I examine the NPC
    Then the NPC should have HP (hit points)
    And the NPC should have an attack stat
    And the NPC should have combat behavior defined

  Scenario: Attack an NPC
    Given an NPC "Gizmo" is in the room
    And the NPC is hostile or neutral
    When I type "attack gizmo"
    Then combat should begin
    And the NPC should take damage
    And I should see the damage dealt

  Scenario: NPC attacks back in combat
    Given I am in combat with an NPC
    When the combat cycle advances
    Then the NPC should attack me
    And I should take damage
    And I should see the damage received

  Scenario: NPC combat AI - aggressive
    Given an aggressive NPC
    When I enter the room
    Then the NPC should attack me automatically
    And combat should begin without my command

  Scenario: NPC combat AI - neutral
    Given a neutral NPC
    When I enter the room
    Then the NPC should not attack me automatically
    And I should be able to choose to attack or not

  Scenario: NPC death and XP reward
    Given I am in combat with an NPC
    And the NPC has 10 HP remaining
    When I deal damage that reduces NPC HP to 0 or below
    Then the NPC should die
    And I should receive XP reward
    And I should see a death message
    And the NPC should be removed from the room

  Scenario: NPC HP bar display
    Given I am in combat with an NPC
    When I view combat status
    Then I should see the NPC's HP bar
    And the HP bar should show current/max HP
    And the bar should update after each attack

  Scenario: Player death in combat
    Given I am in combat with an NPC
    And my HP is low
    When the NPC deals damage that reduces my HP to 0 or below
    Then I should die
    And I should see a death message
    And I should respawn at an appropriate location

  Scenario: Flee from combat
    Given I am in combat with an NPC
    When I type "flee" or move to another room
    Then combat should end
    And I should escape to the new room
    And the NPC should remain in the original room

  Scenario: NPC respawns after death
    Given I killed an NPC
    And the NPC has a respawn timer
    When the respawn timer expires
    Then the NPC should reappear in its spawn room
    And the NPC should have full HP restored

  Scenario: Multiple NPCs in combat
    Given a room with multiple NPCs
    When I attack one NPC
    Then only the targeted NPC should be in combat with me
    And other NPCs should remain neutral unless aggressive

  Scenario: NPC combat - different attack patterns
    Given NPCs with different combat behaviors
    When I fight different NPCs
    Then some NPCs should attack fast with low damage
    And some NPCs should attack slow with high damage
    And NPCs should have varied combat styles

  Scenario: NPCs drop items on death
    Given I kill an NPC
    When the NPC dies
    Then items should drop to the room floor
    And I should see what items dropped
    And I should be able to pick them up