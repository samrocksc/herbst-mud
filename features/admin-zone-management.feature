@admin @zone @map-editor
Feature: Zone management in the map editor sidebar

  Zones group rooms into named geographic or quest areas. The map editor
  sidebar shows a Zones section where the admin can:
    1. Create zones manually (id, name, color, min_level)
    2. Add rooms to a zone in chip fashion (search and add; only existing
       rooms are addable)
    3. See "removed" rooms as red chips (the room was deleted but the
       zone still references it)
    4. Multi-select zones when editing a room

  Background:
    Given the dev stack is running
    And I am logged in as the admin user "sma" in world 2 (Ooze Surfers)

  Scenario: Create a new zone in the sidebar
    When I open the map editor
    And I expand the "Zones" section
    And I click "+ New zone"
    And I fill in id="junkyard_test", name="Junkyard Test", color="#8b4513", min_level=3
    And I click "Create zone"
    Then the zone "Junkyard Test" appears in the Zones list
    And the API has a zone with id "junkyard_test", world_id "2", min_level 3

  Scenario: Add an existing room to a zone (chip fashion)
    Given a zone "Junkyard Test" exists in world 2
    And room id 90 exists in world 2
    When I expand the zone "Junkyard Test" in the sidebar
    And I search for "90" in the room search
    And I click on the result "90 · Archway Approach"
    Then a green chip for "90 · Archway Approach" appears in the zone
    And the API has zone "junkyard_test" with room_ids including 90

  Scenario: Removed room shows as a red chip
    Given a zone "Junkyard Test" exists with room_ids [90, 91, 92, 9999]
    And room 9999 does not exist in world 2
    When I expand the zone "Junkyard Test" in the sidebar
    Then the chip for room 9999 is rendered in red
    And the chip tooltip says "room has been removed from this world"
    And the chips for 90, 91, 92 are rendered in normal (non-red) style

  Scenario: Non-existent room cannot be added
    Given a zone "Junkyard Test" exists in world 2
    When I expand the zone in the sidebar
    And I type "999999" in the room search
    Then no search results are shown
    And there is no "+ add" button for a non-existent room

  Scenario: Multi-select zones in the room editor
    Given room id 90 exists in world 2
    And zones "Zone A" and "Zone B" exist in world 2
    When I click on room 90 in the map
    And I click "Edit"
    And I check the boxes for "Zone A" and "Zone B" in the Zones section
    And I click "Save Changes"
    Then the API has room 90 with zone_ids including "Zone A" and "Zone B"
    And zone "Zone A" has room_ids including 90
    And zone "Zone B" has room_ids including 90

  Scenario: Removing a room from a zone (chip ×)
    Given room 90 is in zone "Zone A"
    When I expand the "Zone A" zone in the sidebar
    And I click the "×" on the chip for room 90
    Then the chip is removed
    And the API has zone "Zone A" without 90 in its room_ids

  Scenario: Deleting a zone cleans up room memberships
    Given room 90 is in zone "ToDelete" in world 2
    When I click the "×" next to "ToDelete" in the Zones section
    And I confirm the deletion
    Then the zone is removed from the list
    And room 90's zone_ids no longer includes "ToDelete"
    And the API returns 404 for GET /api/zones/ToDelete
