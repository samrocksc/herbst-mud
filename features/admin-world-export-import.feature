Feature: Admin - World Export/Import E2E (admin-19)
  As an admin/builder
  I want to export and re-import complete game worlds
  So that I can back up, transfer, or reset world data

  Background:
    Given I am logged into the admin interface
    And I have API access credentials
    And I have at least one world with data (rooms, NPCs)

  @donnie @raph @splinter
  Scenario: Export world returns valid JSON
    Given I have a world with ID "cyberpunk"
    When I call GET /admin/export?world=cyberpunk
    Then I should receive a valid JSON response
    And the response should include: version, exported_at, rooms, npcs, skills, items
    And version should be "1.0"

  @donnie @raph @splinter
  Scenario: Exported world data is re-importable
    Given I have exported world data from "cyberpunk"
    When I call POST /admin/import with the exported JSON
    Then the import should succeed
    And I should receive confirmation with counts
    And the imported data should match the original

  @donnie @raph @splinter
  Scenario: Export with empty world returns skeleton
    Given I have a world "herbst-mud" with no rooms/NPCs
    When I call GET /admin/export?world=herbst-mud
    Then I should receive valid JSON
    And rooms array should be empty (not missing)
    And npcs array should be empty (not missing)
    And skills array should have default skills
    And items array should be empty

  @donnie @raph @splinter
  Scenario: Export filters by world_id correctly
    Given I have rooms with world_id "cyberpunk" and "default"
    When I export "cyberpunk" world
    Then the exported rooms should only include cyberpunk world_id
    And default world rooms should not be included

  @donnie @raph @splinter
  Scenario: Round-trip export/import preserves data integrity
    Given I have a world with rooms and NPCs
    When I export the world to JSON
    And I delete the rooms from the database
    And I re-import the exported JSON
    Then the rooms should be restored
    And the NPCs should be restored
    And exits should be preserved

  @donnie @raph @splinter
  Scenario: Export worlds endpoint returns all DB-backed worlds
    Given I have worlds in the database: "herbst-mud", "cyberpunk", "Ooze Surfers"
    When I call GET /admin/export/worlds
    Then I should receive all 3 worlds
    And each world should have: id, name, description, status
    And the worlds should include my custom worlds

  @donnie @raph @splinter
  Scenario: Export with world filter defaults to "default"
    Given I have not specified a world parameter
    When I call GET /admin/export
    Then the export should default to "default" world
    And rooms should be filtered by world_id="default"

  @donnie @raph @splinter
  Scenario: Import validates JSON structure
    Given I have invalid JSON (missing version, wrong format)
    When I call POST /admin/import with invalid JSON
    Then I should receive a validation error
    And the import should not modify any data
