Feature: Admin NPC instance equipment management

  Background:
    Given I am authenticated as an admin
    And the active world is "Ooze Surfers"
    And an NPC template "Gizmo" exists with a spawned instance

  Scenario: Add equipment to an NPC instance from a template
    Given I navigate to the NPC instance detail page
    When I click "Equipment"
    And I search for equipment template "Chef Knife"
    And I select "Chef Knife" and click "Add"
    Then "Chef Knife" should appear in the equipment list as equipped

  Scenario: Toggle equipment equipped state
    Given the NPC instance has "Chef Knife" equipped
    And I am on the instance "Equipment" tab
    When I click "Unequip" on "Chef Knife"
    Then "Chef Knife" should show as not equipped
    When I click "Equip" on "Chef Knife"
    Then "Chef Knife" should show as equipped

  Scenario: Remove equipment from an NPC instance
    Given the NPC instance has "Chef Knife" in its inventory
    And I am on the instance "Equipment" tab
    When I click "Remove" on "Chef Knife"
    Then "Chef Knife" should no longer appear in the equipment list
