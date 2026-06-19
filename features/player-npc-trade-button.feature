Feature: Player web client shows Trade button for shopkeeper NPCs

  Background:
    Given I am authenticated as sma
    And I have a character in Ooze Surfers
    And an NPC "Gearmo Del Toro" with disposition "shopkeeper" is in the same room

  Scenario: Shopkeeper NPC displays Trade button in room actions
    Given I am in the game screen
    When I click on "Gearmo Del Toro" in the Characters list
    Then I should see the action buttons "Attack", "Talk", "Trade", and "Examine"

  Scenario: Trade button sends shop command
    Given I can see "Gearmo Del Toro" in the room
    When I click the "Trade" button
    Then the command input should contain "shop Gearmo Del Toro"
    And the command should be sent to the server
