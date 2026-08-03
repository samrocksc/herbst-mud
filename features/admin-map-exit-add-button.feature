Feature: Admin Map Per-Direction "+ Add Room" Button
  As a game administrator
  I want a "+" button next to every empty direction in the right sidebar
  So that I can create new rooms with automatic bidirectional exits and
  BFS-laid positions, without leaving the room detail view

  Background:
    Given I am authenticated as an admin
    And the active world is "Ooze Surfers"

  @MAP-EXIT-ADD-BUTTON-001 @reproduction
  Scenario: Clicking "+" next to a direction opens NewRoomModal
    Given I have a room "A Dank Ass Closet" with no exits
    When I select the room in the right sidebar
    And I click the "+" button next to "north"
    Then the "NewRoomModal" should be visible
    And the modal title should be "Add Room to the north"
    And the modal should pre-fill "parentRoom" with "A Dank Ass Closet"
    And the modal should pre-fill "direction" with "north"

  @MAP-EXIT-ADD-BUTTON-001 @reproduction
  Scenario: Confirming the modal creates a new room with a wired exit
    Given the "NewRoomModal" is open with parent "A Dank Ass Closet" and direction "north"
    When I type "Northern Passage" in the name field
    And I click "Create"
    Then a new room named "Northern Passage" should be created
    And the new room should have posZ equal to the parent's posZ
    And a bidirectional exit from "A Dank Ass Closet" should point north to "Northern Passage"
    And a bidirectional exit from "Northern Passage" should point south to "A Dank Ass Closet"
    And the new room should be auto-selected in the right sidebar
    And the React Flow canvas should render the new room one step north of the parent

  @MAP-EXIT-ADD-BUTTON-001 @regression
  Scenario: A "+" button appears for every empty direction in the sidebar
    Given I have a room with no exits
    When I select the room in the right sidebar
    Then the right sidebar should show 10 "+" buttons
    And the buttons should be labeled "Add room to the (north|northeast|east|southeast|south|southwest|west|northwest|up|down)"

  @MAP-EXIT-ADD-BUTTON-001 @regression
  Scenario: Creating a room "up" places the new room on the floor above
    Given I have a room "Basement" with posZ -1
    When I select the room in the right sidebar
    And I click the "+" button next to "up"
    And I confirm the modal with name "Ground Floor"
    Then the new room should have posZ 0
    And a bidirectional exit from "Basement" should point up to "Ground Floor"
    And a bidirectional exit from "Ground Floor" should point down to "Basement"

  @MAP-EXIT-ADD-BUTTON-001 @regression
  Scenario: Existing exits do not show a "+" button
    Given I have a room "A Dank Ass Closet" with exit "north" pointing to "Arch of Smashed Cars"
    When I select "A Dank Ass Closet" in the right sidebar
    Then no "+" button should appear next to "north"
    And a "+" button should still appear next to "south"
