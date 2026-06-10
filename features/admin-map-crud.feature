Feature: Admin Map Room CRUD
  As a game administrator
  I want to create, read, update, and delete rooms on the map via the admin API
  So that I can build the MUD world with multiple Z-levels

  Scenario: Create a new room on floor 0
    Given I am authenticated as an admin
    When I POST to "/api/rooms" with:
      | name           | Central Plaza    |
      | description    | A bustling hub   |
      | posX           | 0                |
      | posY           | 0                |
      | posZ           | 0                |
      | atmosphere     | air              |
    Then the response status should be 201
    And the response should contain a room with posZ 0

  Scenario: Create a room on a non-zero Z-level
    Given I am authenticated as an admin
    When I POST to "/api/rooms" with:
      | name           | Tower Spire     |
      | posZ           | 3               |
    Then the response status should be 201
    And the response should contain a room with posZ 3

  Scenario: Reject a room with empty name
    Given I am authenticated as an admin
    When I POST to "/api/rooms" with:
      | name |     |
    Then the response status should be 400

  Scenario: List all rooms
    Given I am authenticated as an admin
    When I GET "/api/rooms"
    Then the response status should be 200
    And the response should be a JSON array

  Scenario: List rooms filtered by world
    Given I am authenticated as an admin
    When I GET "/api/rooms?world_id=1"
    Then the response status should be 200
    And every room should have world_id "1"

  Scenario: Update a room name
    Given I am authenticated as an admin
    And a room with name "Central Plaza" exists
    When I PUT to "/api/rooms/{id}" with:
      | name | Central Plaza (Renovated) |
    Then the response status should be 200
    And the room name should be "Central Plaza (Renovated)"

  Scenario: Reject update with stale version
    Given I am authenticated as an admin
    And a room with name "Central Plaza" exists at version 1
    When I PUT to "/api/rooms/{id}" with:
      | name    | Should Fail    |
      | version | 999            |
    Then the response status should be 409

  Scenario: Delete a room with no characters
    Given I am authenticated as an admin
    And a room with id 999999 does not exist
    When I DELETE "/api/rooms/999999"
    Then the response status should be 404

  Scenario: Create a bidirectional exit
    Given I am authenticated as an admin
    And two rooms "Source" and "Target" exist
    When I POST to "/api/rooms/{source_id}/exits/bidirectional" with:
      | direction    | north     |
      | targetRoomId | {target_id} |
    Then the response status should be 200
    And the source room should have an exit "north" pointing to the target
    And the target room should have an exit "south" pointing to the source

  Scenario: Delete a bidirectional exit
    Given I am authenticated as an admin
    And a room has an exit "north" to another room
    When I DELETE "/api/rooms/{room_id}/exits/bidirectional?direction=north"
    Then the response status should be 200
    And the source room should not have an exit "north"
    And the target room should not have an exit "south"

  Scenario: Clean up orphan exits
    Given I am authenticated as an admin
    When I POST to "/api/rooms/cleanup-orphan-exits" with:
      | {} |
    Then the response status should be 200
    And the response should include a "cleaned" count

  @MAP-CLEANUP-WORLD
  Scenario: Cleanup orphan exits scoped to a world
    Given I am authenticated as an admin
    And rooms exist in world "1" and world "2"
    And an exit in a world "2" room points to a deleted room
    When I POST to "/api/rooms/cleanup-orphan-exits?world_id=2"
    Then the response status should be 200
    And the orphan exit in world "2" is removed
    And the response should include a "cleaned" count

  Scenario: Delete a room relocates characters to the world's root
    Given I am authenticated as an admin
    And a world has a root room
    And a character is in a non-root room of that world
    When I DELETE the character's room
    Then the character's current_room_id should equal the root room's id
    And the character's current_room_id should be in the same world

  @MAP-SIDEBAR-NAV
  Scenario: Clicking a sidebar room selects it without leaving /map
    Given I am authenticated as an admin
    And I am on the map page at "/map"
    And there is a room "Fountain Plaza" visible in the sidebar
    When I click on "Fountain Plaza" in the sidebar's room list
    Then the URL should still be at "/map"
    And the room detail panel should show "Fountain Plaza"
    And the URL should contain "?room="
