Feature: Exit Edge Component
  As an admin map builder
  I want to see exit connections between rooms
  So that I can visualize room navigation

  Background:
    Given I am in the admin map builder

  Scenario: Display directed edge between rooms
    Given "Room A" has an exit to "Room B" in direction "north"
    When the map is rendered
    Then I should see an edge connecting Room A to Room B

  Scenario: Edge shows exit direction
    Given "Room A" has an exit to "Room B" in direction "north"
    When I look at the edge
    Then I should see a direction arrow pointing north
    And the label should show "N" or "north"

  Scenario: Click edge to edit properties
    Given there is an exit edge between rooms
    When I click on the edge
    Then an exit properties panel should open
    And I can edit the exit details

  Scenario: Two-way exit shows bidirectional arrows
    Given "Room A" has a two-way exit to "Room B"
    When the map is rendered
    Then I should see arrows pointing both directions