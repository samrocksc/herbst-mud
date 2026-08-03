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

  @MAP-RENDER
  Scenario: Map renders rooms and edges correctly for the active world
    Given I am authenticated as an admin
    And the active world is "Ooze Surfers"
    When I navigate to "/map"
    Then I should see room nodes for the current world
    And I should see connecting edges between rooms with exits

  @MAP-EDGE-ORTHOGONAL
  Scenario: Cardinal exits place rooms on a straight horizontal or vertical line
    Given I am authenticated as an admin
    And the active world is "Ooze Surfers"
    And a room "Source" exists at posx=0 posy=0 on floor 0
    And a room "East Target" exists at posx=160 posy=0 on floor 0
    And an exit from "Source" direction "east" points to "East Target"
    When I GET "/api/rooms?world_id=2"
    Then for the "east" exit from "Source" to "East Target" the y-coordinates must be equal
    And the x-distance must be greater than 0

  @MAP-EDGE-ORTHOGONAL
  Scenario: Ordinal exits place rooms on a 45-degree diagonal (equal x and y magnitudes)
    Given I am authenticated as an admin
    And the active world is "Ooze Surfers"
    And a room "Origin" exists at posx=0 posy=0 on floor 0
    And a room "NE" exists at posx=106 posy=-106 on floor 0
    And an exit from "Origin" direction "northeast" points to "NE"
    When I GET "/api/rooms?world_id=2"
    Then for the "northeast" exit from "Origin" to "NE" the x and y distances must be equal in magnitude
    And both distances must be non-zero

  @MAP-EDGE-ORTHOGONAL
  Scenario: Reject (or snap) a room save whose coordinates don't match its exit direction label
    Given I am authenticated as an admin
    And the active world is "Ooze Surfers"
    And a room "Mover" exists at posx=0 posy=0 on floor 0 with exit "east" pointing to a sibling room at posx=160 posy=0
    When I PUT the room "Mover" with posx=0 posy=80 (offset on the wrong axis for an "east" exit)
    Then the response status should be 400 or the server should snap posx/posy so the "east" direction still matches the target
    And the stored "east" exit from "Mover" must still satisfy dy=0

  @MAP-EDGE-CROSSING
  Scenario: Edges between rooms on the same floor must not cross other rooms or each other
    Given I am authenticated as an admin
    And the active world is "Ooze Surfers"
    And the floor 0 / floor +1 / floor -1 has N rooms
    When I GET "/api/rooms?world_id=2"
    Then for every pair of rooms (A,B) and (C,D) on the same floor where {A,B} != {C,D}, the segment A-B must not intersect the segment C-D
    And no segment A-B may pass through the bounding box of a third room E (where E != A and E != B)

  @MAP-EDGE-AUTOFIX
  Scenario: Triggering the realign button snaps every cardinal exit into a horizontal or vertical line
    Given I am authenticated as an admin
    And the active world is "Ooze Surfers"
    And floor 0 has a mix of cardinal and non-cardinal exits
    When I click the "Realign rooms to grid" toolbar button (📐) on floor 0
    Then every cardinal exit (north/south/east/west) on floor 0 must have a corresponding zero delta on the perpendicular axis
    And the change must persist (verified by GET "/api/rooms?world_id=2")
    And no "Cannot read properties of undefined" errors should appear in the console
