Feature: Admin Factions and Items — world-scoped create + list

  The faction, faction-category, and equipment-template surfaces must all
  respect the active world. Creating an entity in world N must:
  (1) POST the world_id in the request body,
  (2) persist the row with world_id = N,
  (3) only appear in the list view when the WorldStore is on N,
  (4) be invisible from other worlds.

  Background:
    Given I am authenticated as an admin
    And the active world is "Ooze Surfers"
    And the categories list is empty in "Ooze Surfers" (no prior "Classes")
    And the categories list is empty in "herbst-mud" (no prior "Classes")
    And the factions list is empty in "Ooze Surfers" (no prior "Pizza Chef")

  Scenario: Create a faction category "Classes" in Ooze Surfers — row stored in world 2 only
    Given I navigate to "/factions" and click the "Categories" tab
    When I click "+ New"
    And I fill the category form with name "Classes" and description "Character classes — defines role and abilities."
    And I click "Create Category"
    Then the server response status should be 201
    And the POST body should include world_id "2"
    And the row in faction_categories should have id, name="Classes", world_id="2"
    And "Classes" should appear in the Categories tab list under Ooze Surfers

  Scenario: Switch to herbst-mud — Classes category is not visible
    Given I switch the WorldStore to "herbst-mud"
    And I navigate to "/factions" and click the "Categories" tab
    Then the list should be empty
    And "Classes" should not appear in the list
    And the GET /api/faction-categories request should include world_id "1"

  Scenario: Switch back to Ooze Surfers and create Pizza Chef faction under Classes
    Given I switch the WorldStore to "Ooze Surfers"
    And I navigate to "/factions" and click "+ New"
    And I fill the faction form with:
      | name        | pizza_chef   |
      | display_name| Pizza Chef   |
      | category_id | <Classes id> |
    And I click "Create Faction"
    Then the server response status should be 201
    And the POST body should include world_id "2" and category_id <Classes id>
    And the row in factions should have world_id="2" and faction_category_factions=<Classes id>
    And "Pizza Chef" should appear in the Factions list

  Scenario: Switch to herbst-mud — Pizza Chef is not visible
    Given I switch the WorldStore to "herbst-mud"
    And I navigate to "/factions"
    Then "Pizza Chef" should not appear in the list
    And the GET /api/factions request should include world_id "1"

  Scenario: Create equipment template "Rolling Pin" in Ooze Surfers — weapon, 1d6+1 bludgeoning
    Given I switch the WorldStore to "Ooze Surfers"
    And I navigate to "/items/new"
    And I fill the item form with:
      | name         | Rolling Pin                       |
      | slot         | main_hand                         |
      | item_type    | weapon                            |
      | weapon_type  | staff                             |
      | damage_type  | bludgeoning                       |
      | dice_count   | 1                                 |
      | dice_sides   | 6                                 |
      | damage_bonus | 1                                 |
    And I click "Create"
    Then the server response status should be 201
    And the POST body should include world_id "2"
    And the row in equipment_templates should have world_id="2"
    And the row in equipment_templates should have damage_dice_count=1, damage_dice_sides=6, damage_bonus=1, damage_type="bludgeoning", weapon_type="staff"

  Scenario: Spawn Rolling Pin instance in room 1 of Ooze Surfers
    Given I navigate to "/items/<Rolling Pin template id>/spawn"
    And I select room "1" in the Room field
    And I click "Spawn Instance"
    Then the server response status should be 201
    And the row in equipment should have equipment_template_id=<template id>
    And the row in equipment should have room_equipment=1

  Scenario: Switch to herbst-mud — Rolling Pin item is not visible
    Given I switch the WorldStore to "herbst-mud"
    And I navigate to "/items"
    Then "Rolling Pin" should not appear in the list
    And the GET /api/equipment-templates request should include world_id "1"

  Scenario: Cross-check — switching back to Ooze Surfers shows all three entities
    Given I switch the WorldStore to "Ooze Surfers"
    And I navigate to "/factions" and click the "Categories" tab
    Then "Classes" should appear in the list
    When I navigate to "/factions" and click the "Factions" tab
    Then "Pizza Chef" should appear in the list
    When I navigate to "/items"
    Then "Rolling Pin" should appear in the list

  Scenario: Regression guard — the useResourceSearch query enables at length 1
    Given I am on the spawn-item page for a template
    When I type a single character into the Room search field
    Then the GET /api/rooms?world_id=2 request should fire
    And the dropdown should show matching rooms or a "no results" state
