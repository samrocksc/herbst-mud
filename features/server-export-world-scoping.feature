Feature: Server — World export scopes abilities to the requested world

  Background:
    Given I am authenticated as an admin
    And the abilities table has at least 5 rows in world "1"
    And the abilities table has at least 5 rows in world "2"

  Scenario: Exporting world 1 returns only world-1 abilities
    Given I call GET /admin/export with world "1"
    Then the response snapshot should have an "abilities" array
    And every ability in the array should have world_id "1"
    And the count of abilities in the payload should match
        `SELECT count(*) FROM abilities WHERE world_id = '1'`

  Scenario: Exporting world 2 returns only world-2 abilities
    Given I call GET /admin/export with world "2"
    Then the response snapshot should have an "abilities" array
    And every ability in the array should have world_id "2"
    And the count of abilities in the payload should match
        `SELECT count(*) FROM abilities WHERE world_id = '2'`

  Scenario: Reproduce the regression (before fix)
    Given the export does not scope abilities by world
    When I call GET /admin/export with world "1"
    And I call GET /admin/export with world "2"
    Then both payloads should be byte-identical for the abilities array
    And every ability in the world-2 payload should have world_id "1"
    And the world-2 payload should not contain any world-2-only abilities
