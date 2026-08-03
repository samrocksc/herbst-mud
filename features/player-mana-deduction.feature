Feature: Player Web Client — Mana Deduction on Ability Cast

  Background:
    Given I am authenticated as sma
    And I have a "smack" character in Ooze Surfers with Trash Bolt equipped in slot 1
    And Gizmo the NPC is in the same room (id=3, "A dank ass closet")

  Scenario: Casting Trash Bolt deducts 5 mana from the player's UI
    Given I am in a room with a hostile NPC (Gizmo, L1, 50/50 HP)
    And my character has 45/50 mana
    When I click on the NPC's name in the CHARACTERS list
    Then I see "Attack" and "Examine" buttons appear
    When I click "Attack"
    Then I see a "Confirm" button
    When I click "Confirm"
    Then combat starts and the COMBAT banner appears
    And the header shows "MANA 45/50"
    When I press "1" to queue Trash Bolt (mana_cost: 5)
    Then the action is queued for the next tick
    When the next combat tick resolves
    Then the header shows "MANA 40/50"
    And the mana decreased by exactly 5
