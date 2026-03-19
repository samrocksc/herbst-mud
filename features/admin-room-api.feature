Feature: Admin - Integrate Room API (admin-04)
  As an admin/builder
  I want to interact with rooms via API
  So that I can programmatically manage room data

  Background:
    Given I am logged into the admin interface
    And I have API access credentials

  @donnie @raph @mikey
  Scenario: Fetch room by ID
    Given a room exists in the database
    When I call GET /api/rooms/{id}
    Then I should receive the room data
    And the response should include: id, name, description, exits, items

  @donnie @raph @mikey
  Scenario: Fetch all rooms
    When I call GET /api/rooms
    Then I should receive a list of all rooms
    And each room should have basic properties

  @donnie @raph @mikey
  Scenario: Create new room via API
    When I call POST /api/rooms with room data
    Then a new room should be created
    And I should receive the created room with ID

  @donnie @raph @mikey
  Scenario: Update room via API
    Given a room exists
    When I call PUT /api/rooms/{id} with updates
    Then the room should be updated
    And I should receive the updated room

  @donnie @raph @mikey
  Scenario: Delete room via API
    Given a room exists
    When I call DELETE /api/rooms/{id}
    Then the room should be deleted
    And subsequent calls should return 404

  @donnie @raph @mikey
  Scenario: Room API returns proper JSON
    When I call the rooms API
    Then responses should be valid JSON
    And Content-Type should be application/json

  @donnie @raph @mikey
  Scenario: API requires authentication
    When I call the API without credentials
    Then I should receive a 401 Unauthorized
    Or access should be denied

  @donnie @raph @mikey
  Scenario: React Query hooks for rooms
    Given the frontend uses React Query
    When I use useRooms() hook
    Then it should fetch and cache room data
    And mutations should update the cache