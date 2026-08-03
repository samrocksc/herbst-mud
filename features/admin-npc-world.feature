Feature: Admin NPCs — world-scoped template + instance create

  The NPC template and NPC instance surfaces must respect the active world.
  Creating an entity in world N must:
  (1) POST the world_id in the request (body or query),
  (2) persist the row with world_id = N,
  (3) only appear in the list view when the WorldStore is on N,
  (4) be invisible from other worlds.

  Background:
    Given I am authenticated as an admin
    And the active world is "Ooze Surfers"
    And the NPC templates list is empty in "Ooze Surfers"
    And the NPC templates list is empty in "herbst-mud"

  Scenario: Create NPC template "Gizmo" in Ooze Surfers — hostile, level 1, respawn room 1
    Given I navigate to "/npcs"
    When I click "+ Add Template"
    And I fill the NPC form with:
      | name               | Gizmo                                                                          |
      | description        | A wiry tinkerer with grease-stained fingers.                                   |
      | level              | 1                                                                              |
      | xp_value           | 10                                                                             |
      | disposition        | hostile                                                                        |
      | respawn_cooldown   | 120                                                                            |
      | respawn_rooms      | 1                                                                              |
      | greeting           | Hey buddy, you got any spare flux capacitors?                                   |
    And I click "Create Template"
    Then the server response status should be 201
    And the POST should include world_id "2"
    And the row in npc_templates should have name="Gizmo", world_id="2", level=1, disposition="hostile", respawn_rooms=["1"]

  Scenario: Spawn NPC instance in room 1 of Ooze Surfers
    Given I navigate to "/npcs/<Gizmo template id>"
    When I click "+ Add Instance"
    And the spawn modal opens with room pre-filled to "A Dank Ass Closet (#1)"
    And I click "Spawn"
    Then the server response status should be 201
    And a row should exist in characters with name="Gizmo", is_npc=true, npc_template_id=<template id>, current_world="2", current_room_id=1
    And the NPC Instances table should show 1 row with location "Room #1"

  Scenario: Switch to herbst-mud — Gizmo template and instance are not visible
    Given I switch the WorldStore to "herbst-mud"
    And I navigate to "/npcs"
    Then "Gizmo" should not appear in the templates list
    And the list should be empty ("No NPC templates found.")

  Scenario: Cross-check — switching back to Ooze Surfers shows Gizmo
    Given I switch the WorldStore to "Ooze Surfers"
    And I navigate to "/npcs"
    Then "Gizmo" should appear in the templates list
    When I navigate to "/npcs/<Gizmo template id>"
    Then the NPC Instances table should show 1 instance of Gizmo

  Scenario: Regression guard — NPC instance creation handles missing race_id without 500
    Given I am on the spawn-instance modal for a template with race_id=NULL
    When I click "Spawn"
    Then the server response status should be 201
    And the character row should have race="human" (default fallback)
    And the server should not return HTTP 500
