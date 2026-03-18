Feature: Exit Edge Component (admin-03)
  As a game designer
  I want to see directed edges between room nodes in the map builder
  So that I can visualize and manage exits between rooms

  Background:
    Given I am logged into the admin panel
    And I navigate to the Map Builder page

  Scenario: View exit direction indicators
    Given there are connected rooms in the map
    When I view the map
    Then I should see arrows indicating exit directions
    And each exit should show its direction label (N, S, E, W, U, D)

  Scenario: Z-exit visual distinction
    Given there are exits between different Z-levels
    When I view the map
    Then Z-exits should have orange styling for "up" exits
    And Z-exits should have purple styling for "down" exits
    And Z-exits should display a "Z-Exit" badge

  Scenario: Click exit to edit properties
    Given there is an exit edge in the map
    When I click on the exit edge
    Then the exit should be highlighted
    And I should be able to edit the exit direction

  Scenario: Hover state for exit edges
    Given there are exit edges in the map
    When I hover over an exit edge
    Then the edge should show a hover effect
    And the edge should be more prominent

  Scenario: Bidirectional exit display
    Given room A has a north exit to room B
    When I view the connection
    Then the edge from A should show "N" direction
    And the edge from B should show "S" direction