Feature: Player Web Client — Combat System

  Background:
    Given I am authenticated as sma
    And I have a Chef character in Ooze Surfers
    And Gizmo the NPC is in the same room (id=1, "A dank ass closet")

  Scenario: Attack NPC by clicking name → Attack → Confirm
    Given I am in a room with a hostile NPC (Gizmo, L1, 100/100 HP)
    When I click on the NPC's name in the CHARACTERS list
    Then I see "Attack" and "Examine" buttons appear
    When I click "Attack"
    Then I see a "Confirm" button
    When I click "Confirm"
    Then combat starts and the COMBAT — ROUND 1 banner appears
    And the combat log shows "⚔ Combat started with Gizmo!"

  Scenario: Auto-attack ticks resolve every ~1.5s
    Given combat is active with Gizmo as target
    When no ability is queued (action bar slots 1-4 are unassigned)
    Then each tick auto-resolves a player attack via the default action
    And the round counter increments by 1 per tick
    And hit/miss resolution uses d20 + DEX mod vs target AC

  Scenario: Critical hit on natural 20
    Given combat is active
    When the player rolls a natural 20
    Then the combat log shows "⚔ CRITICAL HIT! 6 damage!"
    And the damage is doubled (base 3 → crit 6)

  Scenario: Fumble on natural 1
    Given combat is active
    When the player rolls a natural 1
    Then the combat log shows "🎲 FUMBLE! Natural 1 — You stumble badly!"
    And the player's turn ends with no damage dealt

  Scenario: Flee via F button or Flee action
    Given combat is active
    When I click "F Flee"
    Then the action is queued for the next tick
    And the next tick resolves flee via "d20 + floor(level/2) vs DC 12"
    And on success, combat ends and "🏃 Escape successful!" appears in the log
    And on failure, combat continues and "🏃 Escape failed!" appears

  Scenario: NPC death and respawn
    Given combat is active with an NPC
    When the NPC reaches 0 HP
    Then the NPC is marked dead (died_at populated, current_room_id cleared)
    And after a respawn interval the NPC returns with full HP

  Scenario: Ooze Surfers spawn room is not the NPC's room
    Given I am a newly created Ooze Surfers character
    When I finish character creation
    Then I land in Fountain Plaza (room id=4, world_id=default — the dev world)
    And Gizmo is in "A dank ass closet" (room id=1, world_id=2 — Ooze Surfers)
    And the rooms are NOT connected (no shared exit)
    And a player without direct DB access cannot reach Gizmo without a movement path

  Scenario: Character class field ignored by backend (silent survivor default)
    Given I select class "chef" on the character creation form
    When I submit the form
    Then the API persists the character with class="survivor" silently
    And the character-select card shows "survivor" in the subtitle
