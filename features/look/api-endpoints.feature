🔴 Feature: API Endpoints for Look System - Issue Look-10 / Issue #18
  As a game system
  I need API endpoints that support look and examine functionality
  So that the game client can retrieve world data

  Background:
    Given the API server is running
    And the look/examine system is implemented

  # ROOM ENDPOINTS
  Scenario: GET /api/rooms/{id} returns room data
    Given a room exists with ID 1
    When a client sends GET /api/rooms/1
    Then the response should be 200 OK
    And the body should contain room name, description, exits

  Scenario: GET /api/rooms/{id} returns characters in room
    Given a room has characters inside
    When a client retrieves the room
    Then character list should include visible characters
    And character data should be appropriate (not full stats)

  Scenario: GET /api/rooms/{id} returns items in room
    Given a room has items on the ground
    When a client retrieves the room
    Then items list should include visible items
    And hidden items should not appear (unless criteria met)

  Scenario: GET /api/rooms/{id} returns exits
    Given a room has exits
    When a client retrieves the room
    Then exits should list directions and target room IDs
    And locked exits should indicate lock status

  # EXAMINE ENDPOINTS
  Scenario: GET /api/items/{id} returns item details
    Given an item exists with ID
    When a client sends GET /api/items/{id}
    Then the response should contain item name, description, stats

  Scenario: GET /api/items/{id} respects examine skill
    Given an item has hidden detail fields
    When a client retrieves the item
    Then hidden fields should be filtered based on viewer skill
    Or a skill_required flag should be returned

  Scenario: GET /api/characters/{id} returns character details
    Given a character (NPC or player) exists
    When a client sends GET /api/characters/{id}
    Then the response should contain name, description, level
    And NPCs should show relevant game data

  # AUTHENTICATION
  Scenario: Authenticated users can access look endpoints
    Given a player is authenticated
    When the player requests room data
    Then the request should be authorized
    And player-specific data should be included

  Scenario: Unauthenticated requests are rejected
    Given a client is not authenticated
    When the client requests /api/rooms or /api/characters
    Then the response should be 401 Unauthorized
    And no game data should be returned
