Feature: Look/Examine API Endpoints
  As a client developer
  I want REST API endpoints for look and examine
  So that the web admin and other clients can access room data

  Background:
    Given the server is running on port 8080
    And a room exists with items and NPCs

  Scenario: GET /rooms/:id returns room with contents
    When I call GET /rooms/:id
    Then the response includes room name and description
    And the response includes exits
    And the response includes items in the room
    And the response includes NPCs in the room

  Scenario: GET /items/:id returns item details
    When I call GET /items/:id
    Then the response includes item name
    And the response includes short_desc
    And the response includes type
    And the response includes is_immovable flag

  Scenario: GET /items/:id/examine returns detailed info
    Given the item has hidden details
    When I call GET /items/:id/examine
    Then the response includes examine_desc
    And the response includes hidden_details array
    And each hidden detail shows revealed status

  Scenario: Examine endpoint performs skill check
    Given a character's examine skill level is known
    When GET /items/:id/examine is called with character context
    Then hidden details are filtered by skill level
    And revealed details are included

  Scenario: Examine response includes XP
    When I call GET /items/:id/examine
    Then the response includes examine_xp field
    And XP granted reflects examine action

  Scenario: GET /npcs/:id returns NPC details
    When I call GET /npcs/:id
    Then the response includes NPC name
    And the response includes description

  Scenario: GET /npcs/:id/examine returns detailed NPC info
    When I call GET /npcs/:id/examine
    Then the response includes detailed description
    And the response includes stats if applicable

  Scenario: GET /rooms/:id/look is an alias for room details
    When I call GET /rooms/:id/look
    Then the response is equivalent to GET /rooms/:id
    And it includes current room state

  Scenario: GET /equipment/:id/examine endpoint
    Given equipment item exists
    When I call GET /equipment/:id/examine
    Then examine details are returned
    And skill-gated content respects character context

  Scenario: Invalid room ID returns 404
    When I call GET /rooms/invalid-id
    Then the response is 404 Not Found
    And error message is returned

  Scenario: Item not in room returns appropriate error
    When I call GET /items/:id for item not accessible to character
    Then 404 or 403 is returned
