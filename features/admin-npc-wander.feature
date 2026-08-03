Feature: Admin NPC Templates — Behavior section (wander / patrol / return_home)

  The NPC template form must expose a "Behavior" section that lets admins
  configure how the NPC moves between rooms. The default pattern is
  "static" (NPC never moves). The other three patterns (wander, patrol,
  return_home) drive the server-side RoamingService that ticks NPCs
  according to their per-template cadence and zone constraints.

  Background:
    Given I am authenticated as an admin
    And the active world is "Ooze Surfers"
    And the NPC templates list is empty in "Ooze Surfers"

  Scenario: Create NPC template "Wander Goblin" with pattern=wander and zone restriction
    Given I navigate to "/npcs"
    When I click "+ Add Template"
    And I fill the NPC form with:
      | name                    | Wander Goblin                                            |
      | description             | A twitchy scout that roams the Junkyard.                 |
      | level                   | 3                                                        |
      | xp_value                | 25                                                       |
      | disposition             | hostile                                                  |
      | roam_pattern            | wander                                                   |
      | roam_interval_seconds   | 90                                                       |
      | roam_pause_min_seconds  | 10                                                       |
      | roam_pause_max_seconds  | 60                                                       |
      | roam_zone_ids           | zone_junkyard                                            |
      | notify_on_enter         | true                                                     |
    And I click "Create Template"
    Then the server response status should be 201
    And the row in npc_templates should have name="Wander Goblin", roam_pattern="wander", roam_interval_seconds=90, roam_zone_ids=["zone_junkyard"], notify_on_enter=true
    And the row in npc_templates should have defaults roam_pause_min_seconds=10, roam_pause_max_seconds=60

  Scenario: New NPC template defaults to roam_pattern=static when section is left blank
    Given I navigate to "/npcs"
    When I click "+ Add Template"
    And I fill the NPC form with:
      | name        | Stationary Guard  |
      | description | Stands at the gate. |
      | level       | 5                  |
      | disposition | neutral            |
    And I click "Create Template"
    Then the server response status should be 201
    And the row in npc_templates should have name="Stationary Guard", roam_pattern="static", roam_interval_seconds=60, roam_pause_min_seconds=15, roam_pause_max_seconds=120, notify_on_enter=true

  Scenario: Edit existing NPC template — change pattern from static to wander
    Given an NPC template exists with name="Stationary Guard" and roam_pattern="static"
    And I navigate to "/npcs/<Stationary Guard template id>"
    When I click "Edit"
    And I change roam_pattern to "wander"
    And I change roam_interval_seconds to 30
    And I change roam_zone_ids to "zone_junkyard,zone_junkyard_inner"
    And I click "Save"
    Then the server response status should be 200
    And the row in npc_templates should have name="Stationary Guard", roam_pattern="wander", roam_interval_seconds=30, roam_zone_ids=["zone_junkyard","zone_junkyard_inner"]
    And the read-only "Behavior" section should display "Wander (random within zones)"

  Scenario: Invalid roam_pattern is rejected by the server
    Given I navigate to "/npcs"
    When I click "+ Add Template"
    And I fill the NPC form with:
      | name          | Bad NPC     |
      | level         | 1           |
      | disposition   | neutral     |
      | roam_pattern  | teleport    |
    And I click "Create Template"
    Then the server response status should be 400
    And the response error should mention "invalid roam_pattern"

  Scenario: Behavior section appears on the NPC template detail view
    Given an NPC template exists with name="Wander Goblin" and roam_pattern="wander"
    And I navigate to "/npcs/<Wander Goblin template id>"
    Then the read-only "Template Stats" view should include a "Behavior" sub-block
    And the Behavior sub-block should show:
      | label             | value                          |
      | Roam Pattern      | Wander (random within zones)   |
      | Roam Interval     | 90s                            |
      | Roam Pause Min    | 10s                            |
      | Roam Pause Max    | 60s                            |
      | Roam Zone IDs     | zone_junkyard                  |
      | Notify on Enter   | Yes                            |

  Scenario: Regression guard — roam fields are persisted and re-loaded after refresh
    Given an NPC template exists with name="Wander Goblin" and roam_pattern="wander" and roam_zone_ids=["zone_junkyard"]
    And I navigate to "/npcs/<Wander Goblin template id>"
    And I click "Edit"
    Then the Behavior section should pre-fill with roam_pattern="wander", roam_interval_seconds=90, roam_zone_ids="zone_junkyard", notify_on_enter=true
    And the server should not return HTTP 500
