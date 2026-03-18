Feature: Admin - Exit Edge Component (admin-03)
  As an admin/builder
  I want to create visual exit connections between rooms
  So that I can visualize how rooms connect in the map

  Background:
    Given I am logged into the admin interface
    And I am in the map builder view

  @mikey
  Scenario: View existing exit edges
    Given two rooms are connected
    When I view the map
    Then I should see a line/edge connecting the room nodes
    And the edge should show the exit direction

  @mikey
  Scenario: Exit edge shows direction label
    Given Room A connects north to Room B
    When I view the connection
    Then I should see "N" or "North" label on the edge
    And the edge should point from Room A toward Room B

  @mikey
  Scenario: Exit edge updates on room move
    Given rooms are connected with an edge
    When I drag Room A to a new position
    Then the edge should follow and stay connected
    And the edge should maintain the correct path

  @mikey
  Scenario: Bidirectional exit shows two edges
    Given Room A connects to Room B and back
    When I view the connections
    Then I should see two edges: A→B and B→A
    And each should have correct direction labels

  @mikey
  Scenario: Edge is styled by exit type
    Given exits can have different types (normal, locked, hidden)
    When I view the edges
    Then normal exits should have solid lines
    And locked exits could have different styling (dashed/locked icon)
    And hidden exits should be invisible or dashed